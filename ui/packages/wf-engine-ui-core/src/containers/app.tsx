import "@xyflow/react/dist/base.css";
import "@xyflow/react/dist/style.css";

import { StepType, WorkflowDefinition } from "@paulhalleux/wf-engine-api";
import { useQuery } from "@tanstack/react-query";
import { useState } from "react";

import {
  createDecisionStepHandler,
  createForkStepHandler,
  createJoinStepHandler,
  createTaskStepHandler,
  createWaitStepHandler,
  createWorkflowStepHandler,
} from "../factory/step-handlers";
import { WorkflowDefinitionsQuery } from "../query";
import styles from "./app.module.css";
import { WorkflowDefinitionGraph } from "./WorkflowDefinitionGraph.tsx";

export function App() {
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
