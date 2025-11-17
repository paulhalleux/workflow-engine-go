# Revue du projet Go `engine`

Ce document liste les points qui pourraient être améliorés dans ton premier projet Go, à la fois côté implémentation et côté architecture. L’objectif est d’identifier ce qui fonctionne, ce qui est risqué à moyen/long terme et des pistes d’évolution concrètes.

> Note : ce document ne signifie pas que le projet est « mauvais » – au contraire, pour un premier projet il est déjà bien structuré – mais il met le focus sur les éléments perfectibles.

---

## 1. Architecture globale du projet

### 1.1. Couplage fort entre `engine.Engine` et l’infrastructure

**Constat**
- `engine.Engine` construit directement :
  - la connexion DB (`gorm.Open` dans `createDatabase`)
  - la couche persistence (`persistence.NewPersistence`)
  - les services (workflow, steps, etc.)
  - les exécuteurs (workflow / step executor)
  - les channels de communication (maps de `chan`)
  - les serveurs HTTP/gRPC (dans `startHttpServer` / `startGrpcServer`).
- Tout est instancié dans `NewEngine`, ce qui rend les tests difficiles et empêche l’injection de dépendances.

**Risques / limitations**
- Difficile de tester séparément chaque couche (par exemple un service métier sans DB réelle).
- Difficile de remplacer la DB, d’ajouter des middlewares, de changer de framework HTTP.
- La construction devient vite complexe à maintenir si tu ajoutes d’autres composants.

**Améliorations possibles**
- Introduire un **layer d’initialisation/composition** (par exemple dans `cmd/engine/main.go` ou un package `internal/app`).
  - `main` construit les dépendances, puis les injecte dans un type `App`/`Engine` plus simple.
  - `Engine` pourrait se concentrer sur **la logique métier** et l’orchestration, pas sur la création des ressources.
- Utiliser une interface pour la persistence, les services, etc., afin de permettre le mock en tests.

---

### 1.2. Organisation des packages `internal`

**Points positifs**
- Bonne séparation par responsabilité : `httpapi`, `grpcapi`, `persistence`, `service`, `models`, `utils`.
- Utilisation de `internal/` pour limiter la surface publique du module.

**Points perfectibles**
- Le package `service` contient plusieurs responsabilités :
  - logique pure (services `WorkflowDefinitionsService`, `WorkflowInstanceService`, etc.)
  - exécution concurrente (executors, queues, semaphores, channels).
- `internal` mélange :
  - config (`config.go`)
  - registry d’agent (`agent_registry.go`)
  - websockets (`websocket.go`)

**Suggestions**
- Créer des sous-packages plus précis, par exemple :
  - `internal/app` ou `internal/bootstrap` : construction et wiring de l’application.
  - `internal/execution` : `WorkflowExecutor`, `StepExecutor`, types `WorkflowExecution`, `StepResult`, etc.
  - `internal/agent` : `AgentRegistry`, connectors, modèles de tâche.
  - `internal/ws` : `WebsocketHub`, `WebsocketClient`, helpers.
- Garder un package par contexte métier (bounded context) pour clarifier où va quelle logique.

---

### 1.3. Gestion de la configuration

**Constat**
- `WorkflowEngineConfig` est une simple struct avec des champs publics, instanciée dans `cmd/engine/main.go` à partir de variables d’environnement.
- Aucune validation n’est faite sur les valeurs lues.
- `GrpcAddress` / `HttpAddress` sont des `*string` alors que le reste est en `string`.

**Problèmes potentiels**
- Démarrage possible avec une config invalide (ports vides, DB incomplète, etc.).
- `*string` pour les adresses est plus compliqué que nécessaire (gestion des `nil`).

**Améliorations proposées**
- Ajouter une fonction de **chargement + validation** de la configuration, par ex. :
  - `func LoadConfigFromEnv() (WorkflowEngineConfig, error)`
  - qui vérifie la présence des champs obligatoires et retourne une erreur claire.
- Remplacer `*string` par un type simple (string) et utiliser une valeur vide/"0.0.0.0" comme défaut, ou un type dédié (`HostPort`, etc.).
- Distinguer les ports HTTP/gRPC pour les environnements dev/prod (flag CLI, fichier de config, etc.).

---

### 1.4. Structure côté HTTP/gRPC

**Côté HTTP**
- `startHttpServer` crée un `gin.Engine` et enregistre les handlers via des types `New*Handlers(...)`.
- Les routes sont toutes montées sous `/api`.

**Points d’amélioration**
- Pas de versionnement d’API (`/api/v1/...`).
- CORS ouvert sur `*` sans possibilité de configuration.
- `corsMiddleware` est défini dans le package `engine`, ce qui mélange infra HTTP et logique moteur.

**Côté gRPC**
- `startGrpcServer` : création très directe sans possibilité d’injecter des options (intercepteurs, logs, auth, etc.).

**Suggestions**
- Ajouter un préfixe de version : `/api/v1/...`.
- Externaliser la création du `gin.Engine` dans un package `internal/httpserver` avec :
  - construction configurée (CORS, middlewares, logging, récupération des panic, etc.).
- Faire de même pour le serveur gRPC (`internal/grpcserver`).
- Ajouter progressivement :
  - intercepteurs (logging, metrics, auth),
  - gestion centralisée des erreurs.

---

### 1.5. Accès concurrent aux maps et channels

**Constat**
- `Engine` possède :
  - `agentTaskChan          map[string]chan *proto.NotifyTaskStatusRequest`
  - `workflowChan           map[string]chan *service.WorkflowExecutionResult`
  - `workflowStepOutputChan map[string]chan *service.StepResult`
- `WorkflowExecutor` manipule `workflowChan` et `workflowStepOutputChan` via des pointeurs sur ces maps.
- Ces maps sont lues/écrites dans des goroutines (par exemple `startWorkflow` qui crée/supprime des entrées pendant que d’autres goroutines les utilisent).

**Problème**
- Ces maps ne sont **pas protégées** par des mutex ou des structures thread-safe.
- En Go, l’accès concurrent non synchronisé à une map peut entraîner des **panic runtime** (`concurrent map read and map write`).

**Améliorations nécessaires**
- Introduire une abstraction thread-safe :
  - utiliser `sync.RWMutex` autour de ces maps,
  - ou utiliser `sync.Map` si l’accès est très concurrent et simple,
  - ou encapsuler l’accès dans un type dédié (par ex. `WorkflowChannelRegistry`) avec des méthodes `Register`, `Get`, `Delete`.
- Éviter le passage de `*map[...]chan` dans plusieurs couches ; exposer des méthodes sur un type dédié.

---

## 2. Qualité du code et implémentation

### 2.1. Gestion des erreurs

**Constats**
- Utilisation fréquente de `log.Fatalf` (par exemple dans `NewEngine` pour les erreurs DB / migrations, dans `main` pour le `.env`).
- Beaucoup de logs d’erreur mais peu de propagation vers l’appelant.
- Peu de wrapping d’erreur (pas d’utilisation de `fmt.Errorf("...: %w", err)`).

**Conséquences**
- `log.Fatalf` termine tout le processus ; c’est acceptable au démarrage mais rend certains chemins d’erreur moins contrôlables.
- Difficile pour un appelant (ou des tests) de vérifier la nature d’une erreur.

**Recommandations**
- Limiter `log.Fatalf` à **main seulement** :
  - `NewEngine` devrait retourner `(*Engine, error)` au lieu de panic/log.Fatalf.
- Propager les erreurs vers le haut avec du wrapping contextuel :
  - `return nil, fmt.Errorf("connect DB: %w", err)`.
- Dans les handlers HTTP/gRPC :
  - retourner des codes d’erreur explicitement,
  - loguer avec le contexte (ID de requête, workflow, etc.).

---

### 2.2. Contexte (`context.Context`)

**Constats**
- `Engine` possède un `Context context.Context` créé avec `context.Background()`.
- Ce contexte est passé à `workflowExecutor.Start`, `stepExecutor.Start`, `wsHub.Run` mais il n’est jamais annulé (pas de `CancelFunc`).
- `main` n’écoute pas les signaux système (SIGINT/SIGTERM) pour arrêter proprement.

**Améliorations**
- Utiliser `context.WithCancel` ou `context.WithTimeout` :
  - dans `NewEngine` ou dans `main`, créer un contexte racine et une fonction `cancel`.
  - exposer une méthode `Shutdown()` sur `Engine` qui appelle `cancel`.
- Dans `main`, écouter les signaux OS (`os/signal.Notify`) et appeler `Shutdown` pour favoriser un arrêt propre.

---

### 2.3. Websocket : logs, scopes et sérialisation

`internal/websocket.go`

**Points positifs**
- Utilisation de protobuf pour sérialiser/désérialiser les messages WS (binary, efficace et typé).
- Séparation claire entre `WebsocketHub` et `WebsocketClient`.

**Points perfectibles**
- Beaucoup de logs de debug verbeux (`log.Printf("Sending message to clients: %v", message)`), y compris en production.
- Pas de gestion de reconnection, ni de timeout sur `BroadcastMessage` si un client est lent.
- `IsInScope` logue tous les scopes pour chaque message, ce qui peut être très verbeux.

**Suggestions**
- Introduire un niveau de logs (ou un logger injecté) pour contrôler la verbosité.
- Ajouter un mécanisme de **back-pressure** ou de limite (déjà un buffer 256, mais vérifier la stratégie pour les clients lents : close vs skip). 
- Envisager un découplage de la gestion du scope dans un composant séparé.

---

### 2.4. Persistence et modèles

`internal/persistence/persistence.go` et `internal/models/...`

**Constats**
- `AutoMigrate` est appelé au démarrage, ce qui est pratique mais :
  - pas de contrôle fin sur les migrations,
  - pas de gestion de version de schéma.
- Les modèles (workflow, steps) semblent couplés directement avec GORM.

**Améliorations possibles**
- À terme, envisager un système de migrations dédié (Goose, golang-migrate, etc.).
- Introduire des interfaces de repository :
  - `type WorkflowDefinitionRepository interface { ... }` etc.
- Séparer les modèles GORM des modèles métiers si le couplage devient gênant.

---

### 2.5. Gestion de la concurrence dans `WorkflowExecutor`

`internal/service/workflow_executor.go`

**Constats**
- Utilisation d’un `taskQueue` (chan) et d’un sémaphore (`sem`) pour limiter le nombre de workflows parallèles.
- Usage d’un `sync.WaitGroup` pour synchroniser la fin des steps.
- Mémorisation des résultats des steps dans une map locale `stepOutputMap`.

**Problèmes/risques**
- `workflowChan` et `workflowStepOutputChan` sont des `*map[...]` partagées (cf. section 1.5). Leur utilisation dans des goroutines est non protégée.
- `startWorkflow` écrit/importe dans ces maps pendant que d’autres goroutines peuvent également y toucher.
- `workflowChan` contient un canal par workflow, mais la logique de fermeture/cleanup est dispersée.

**Suggestions**
- Créer un type dédié, par ex. :
  - `type WorkflowChannelRegistry struct { mu sync.RWMutex; workflows map[string]chan *WorkflowExecutionResult; ... }`
- Donner à `WorkflowExecutor` une dépendance sur cette abstraction plutôt que sur un `*map`.
- Documenter le protocole :
  - quand un canal est créé,
  - quand il est fermé,
  - qui est responsable de la fermeture.

---

### 2.6. Style Go / idiomes

**Points perfectibles**
- Noms assez longs et parfois redondants (`WorkflowExecutionService`, `WorkflowExecutor`, etc.). Ce n’est pas bloquant, mais tu peux parfois simplifier.
- Mélange dans certains endroits de `fmt.Println` et `log.Printf` (par exemple dans `websocket.go`).
- Les fonctions utilitaires comme `joinHostPort` pourraient réutiliser `net.JoinHostPort` (après adaptation type).

**Recommandations**
- Harmoniser l’usage du logger (`log` standard ou un logger structuré), éviter les `fmt.Println` dans le code de prod.
- Profiter des fonctions standard (`net.JoinHostPort`, etc.) pour éviter la duplication.
- Lancer `gofmt`, `go vet` et éventuellement `golangci-lint` pour homogénéiser le style.

---

## 3. Tests et outillage

### 3.1. Absence apparente de tests

**Constat**
- Aucun fichier `*_test.go` dans l’arborescence fournie.

**Impact**
- Difficile d’assurer la non-régression quand tu modifies/ajoutes des fonctionnalités.
- La logique concurrente (executors, websockets, registry) est particulièrement critique et mériterait des tests.

**Recommandations**
- Commencer par des **tests unitaires** sur :
  - les services métiers (workflowDefinitions, workflowInstances, stepInstances) en mockant la persistence,
  - les fonctions pures (`IsInScope`, mapping entre modèles/protobuf, pagination).
- Ajouter quelques **tests d’intégration** :
  - démarrage du moteur en mémoire (DB SQLite in-memory ou un mock DB),
  - appel d’un endpoint HTTP simple,
  - exécution d’un workflow très simple (1 step) et vérification de l’état.

---

### 3.2. Makefile et scripts

**Constat**
- Présence d’un `Makefile`, d’un `docker-compose.yml`, et d’un script `clean-swagger-output.bat`.
- Impossible de juger leur contenu sans les ouvrir, mais on note :
  - pas de standardisation documentée pour le build/test.

**Suggestion**
- Ajouter des cibles standard dans le `Makefile` :
  - `make build`
  - `make test`
  - `make lint`
- Ajouter un petit `README.md` au niveau du module `engine` expliquant :
  - comment lancer la DB (docker compose),
  - comment lancer l’engine,
  - comment exécuter les tests.

---

## 4. API et modèles exposés

### 4.1. HTTP Handlers

(Analyse basée sur les noms de fichiers : `workflow_instances_handlers.go`, etc.)

**Points à vérifier / améliorer**
- Validation des payloads JSON (types, champs requis, etc.).
- Gestion des erreurs HTTP homogène (codes, message d’erreur, structure JSON commune).
- Pagination, filtres et tri exposés de manière cohérente (utilisation de `utils/paginate.go`).
- Utilisation de context (timeouts, cancellations) dans les handlers.

**Améliorations possibles**
- Introduire un **format d’erreur standard** (par ex. `{ "error": { "code": "...", "message": "..." } }`).
- Ajouter une couche de validation (par ex. `validator/v10` déjà en indirect dans le `go.mod`).

---

### 4.2. gRPC & Protobuf

**Constats**
- Utilisation d’un module `github.com/paulhalleux/workflow-engine-go/proto` versionné.
- Les services gRPC `EngineService` et `TaskService` sont implémentés dans `internal/grpcapi` (non détaillé ici, mais structurellement c’est propre).

**Améliorations potentielles**
- S’assurer que la gestion des erreurs gRPC est cohérente (utilisation de `status.Error`, `codes.X` plutôt que de retourner des erreurs génériques).
- Documenter les contrats gRPC (dans la doc ou via commentaires proto).

---

## 5. Points de sécurité et robustesse

### 5.1. CORS et exposition HTTP

**Constat**
- `corsMiddleware` autorise `Access-Control-Allow-Origin: *` et `Access-Control-Allow-Credentials: true`.

**Problèmes**
- En production, cela peut être dangereux (toute origine peut effectuer des requêtes avec cookies/credentials).

**Recommandations**
- Rendre la politique CORS configurable.
- Restreindre `Allow-Origin` aux domaines connus en prod.

---

### 5.2. Gestion des secrets et de la configuration DB

**Constat**
- Les paramètres DB viennent de variables d’environnement (bien), mais pas de mécanisme de rotation, ni de masquage dans les logs.

**Améliorations**
- Attention à ne jamais loguer les mots de passe DB.
- Envisager, plus tard, l’utilisation d’un secret manager (Vault, AWS Secret Manager, etc.).

---

## 6. Synthèse et priorisation des améliorations

Pour t’aider à prioriser, voici une vue rapide des améliorations **par importance** :

### Priorité haute
- [ ] Protéger l’accès concurrent aux maps (`workflowChan`, `workflowStepOutputChan`, `agentTaskChan`) via mutex ou abstraction dédiée.
- [ ] Introduire un mécanisme d’arrêt propre avec `context.WithCancel` et gestion des signaux OS dans `main`.
- [ ] Éviter `log.Fatalf` en profondeur : faire remonter les erreurs au `main`.
- [ ] Ajouter des tests unitaires de base (services, fonctions utilitaires, scope WS).

### Priorité moyenne
- [ ] Structurer davantage les packages (`internal/app`, `internal/execution`, `internal/agent`, `internal/ws`).
- [ ] Améliorer la gestion de configuration (chargement + validation).
- [ ] Factoriser la création des serveurs HTTP/gRPC dans des packages dédiés.
- [ ] Standardiser la gestion des erreurs HTTP/gRPC.

### Priorité basse / évolutions futures
- [ ] Introduire un système de migrations DB plus avancé que `AutoMigrate`.
- [ ] Ajouter un README détaillé pour `engine` (build, run, test).
- [ ] Paramétrer plus finement CORS et la sécurité HTTP.
- [ ] Introduire un logger structuré (zap, zerolog, logrus, etc.) et un niveau de log.

---

## 7. Conclusion

Pour un premier projet Go, la base est solide : séparation par couches, utilisation de `internal`, gRPC, HTTP, WebSockets, et une logique d’exécution concurrente déjà bien avancée.

Les principaux axes d’amélioration sont :
- la **sécurité et la robustesse concurrente** (maps/channels),
- la **testabilité** (découpler la construction, exposer des interfaces, ajouter des tests),
- la **clarté de l’architecture** (packages par domaine, gestion des serveurs et de la config).

Si tu veux, on peut prendre un axe précis (par exemple la sécurisation des maps de channels ou la refactorisation de `NewEngine`) et le retravailler ensemble pas à pas dans le code.
