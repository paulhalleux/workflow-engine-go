import { StepType, WorkflowStepDefinition } from "@paulhalleux/wf-engine-api";
import { Edge, Node } from "@xyflow/react";

import { StepHandler } from "./handlers.ts";

export type WorkflowEdgeType = "default";
export type WorkflowNodeType = StepType;
export type WorkflowNodeData = {
  definition: WorkflowStepDefinition;
  handler: StepHandler;
};

export type WorkflowGraphNode = Node<WorkflowNodeData, WorkflowNodeType>;
export type WorkflowGraphEdge = Edge<Record<string, unknown>, WorkflowEdgeType>;

export type WorkflowGraph = {
  nodes: WorkflowGraphNode[];
  edges: WorkflowGraphEdge[];
};

export type FactoryContext = {
  getStepDefinitions(): WorkflowStepDefinition[];
  getStepHandler(stepType: StepType): StepHandler;
};
