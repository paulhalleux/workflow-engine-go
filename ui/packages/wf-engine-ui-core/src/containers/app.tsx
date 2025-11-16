import "@xyflow/react/dist/base.css";
import "@xyflow/react/dist/style.css";

import { StepType, WorkflowDefinition } from "@paulhalleux/wf-engine-api";
import { websocket } from "@paulhalleux/wf-engine-proto";
import { useQuery } from "@tanstack/react-query";
import { useRef, useState } from "react";

import {
  createDecisionStepHandler,
  createForkStepHandler,
  createJoinStepHandler,
  createTaskStepHandler,
  createWaitStepHandler,
  createWorkflowStepHandler,
} from "../factory/step-handlers";
import { useWebsocketConnection } from "../hooks/useWebsocketConnection.ts";
import { WorkflowDefinitionsQuery } from "../query";
import { createProtobufProtocol } from "../utils/websocket.ts";
import styles from "./app.module.css";
import { WorkflowDefinitionGraph } from "./WorkflowDefinitionGraph.tsx";

const protocol = createProtobufProtocol({
  message: websocket.WebsocketMessage,
  command: websocket.WebsocketCommand,
});

export function App() {
  const clientId = useRef<string>(null);
  const { sendMessage } = useWebsocketConnection({
    url: "ws://localhost:8080/ws",
    protocol,
    onMessage: (data) => {
      if (data.type === websocket.WebsocketMessageType.REGISTERED) {
        clientId.current = data.registeredMessage.clientId;
        sendMessage({
          type: websocket.WebsocketCommandType.SUBSCRIBE,
          clientId: clientId.current,
          subscribeCommand: {
            scopes: [],
          },
        });
      } else {
        console.log("WebSocket message received:", data);
      }
    },
  });

  const { data: workflowDefinitions = [] } = useQuery(
    WorkflowDefinitionsQuery.getAll(),
  );

  const [selectedWorkflowDefinition, setSelectedWorkflowDefinition] =
    useState<WorkflowDefinition>();

  const onDefinitionSelect = (definitionId: string) => {
    setSelectedWorkflowDefinition(
      workflowDefinitions.find((wd) => wd.id === definitionId) || null,
    );
  };

  return (
    <div>
      <ul className={styles["workflow-list"]}>
        {workflowDefinitions.map((wd) => (
          <li key={wd.id} onClick={() => onDefinitionSelect(wd.id)}>
            {wd.name} (ID: {wd.id})
          </li>
        ))}
      </ul>
      <div className={styles.graph}>
        {selectedWorkflowDefinition && (
          <WorkflowDefinitionGraph
            key={selectedWorkflowDefinition.id}
            workflowDefinition={selectedWorkflowDefinition}
            handlers={{
              [StepType.StepTypeTask]: createTaskStepHandler(),
              [StepType.StepTypeWait]: createWaitStepHandler(),
              [StepType.StepTypeWorkflow]: createWorkflowStepHandler(),
              [StepType.StepTypeJoin]: createJoinStepHandler(),
              [StepType.StepTypeDecision]: createDecisionStepHandler(),
              [StepType.StepTypeFork]: createForkStepHandler(),
            }}
          />
        )}
      </div>
    </div>
  );
}
