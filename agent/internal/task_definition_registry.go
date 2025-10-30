package internal

import (
	"encoding/json"
	"log"

	"github.com/paulhalleux/workflow-engine-go/agent/internal/proto"
	"github.com/swaggest/jsonschema-go"
	"google.golang.org/protobuf/types/known/structpb"
)

type TaskHandler func()
type TaskDefinition struct {
	ID               string
	Name             string
	Description      string
	InputParameters  *jsonschema.Schema
	OutputParameters *jsonschema.Schema
	Handle           TaskHandler
}

type TaskDefinitionRegistry struct {
	definitions map[string]TaskDefinition
}

func NewTaskDefinitionRegistry() *TaskDefinitionRegistry {
	return &TaskDefinitionRegistry{
		definitions: make(map[string]TaskDefinition),
	}
}

func (r *TaskDefinitionRegistry) Register(def TaskDefinition) {
	r.definitions[def.ID] = def
}

func (r *TaskDefinitionRegistry) Get(id string) (TaskDefinition, bool) {
	def, exists := r.definitions[id]
	return def, exists
}

func (r *TaskDefinitionRegistry) List() []TaskDefinition {
	defs := make([]TaskDefinition, 0, len(r.definitions))
	for _, def := range r.definitions {
		defs = append(defs, def)
	}
	return defs
}

func (r *TaskDefinitionRegistry) ToProto() []*proto.TaskDefinition {
	protoDefs := make([]*proto.TaskDefinition, 0, len(r.definitions))
	for _, def := range r.definitions {
		protoDef := &proto.TaskDefinition{
			Id:               def.ID,
			Name:             def.Name,
			Description:      def.Description,
			InputParameters:  schemaToProto(def.InputParameters),
			OutputParameters: schemaToProto(def.OutputParameters),
		}
		protoDefs = append(protoDefs, protoDef)
	}
	return protoDefs
}

func schemaToProto(p *jsonschema.Schema) *structpb.Struct {
	if p == nil {
		return nil
	}

	j, err := json.Marshal(p)
	if err != nil {
		log.Printf("Error marshaling schema to JSON: %v", err)
		return nil
	}

	var m map[string]interface{}
	err = json.Unmarshal(j, &m)
	if err != nil {
		log.Printf("Error unmarshaling JSON to map: %v", err)
		return nil
	}

	s, err := structpb.NewStruct(m)
	if err != nil {
		log.Printf("Error converting map to structpb.Struct: %v", err)
		return nil
	}

	return s
}
