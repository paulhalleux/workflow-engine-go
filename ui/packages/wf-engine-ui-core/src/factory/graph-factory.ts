import { StepType, WorkflowDefinition } from "@paulhalleux/wf-engine-api";
import invariant from "tiny-invariant";

import {
  FactoryContext,
  WorkflowGraph,
  WorkflowGraphEdge,
  WorkflowGraphNode,
} from "../types/graph.ts";
import { StepHandler, StepHandlerBaseFactory } from "../types/handlers.ts";
import { EdgeFactory } from "./graph-edge-factory.ts";
import { NodeFactory } from "./graph-node-factory.ts";

const createFactoryContext = (
  handlers: Record<StepType, StepHandlerBaseFactory>,
  workflowDefinition: WorkflowDefinition,
): FactoryContext => {
  const handlersCache: Record<string, StepHandler> = {};
  const stepDefinitions = workflowDefinition.steps || [];

  return {
    getStepDefinitions() {
      return stepDefinitions;
    },
    getStepHandler: (stepType: StepType) => {
      if (!handlersCache[stepType]) {
        const createStepHandler = handlers[stepType];
        invariant(
          createStepHandler,
          `No handler found for step type: ${stepType}`,
        );
        handlersCache[stepType] = createStepHandler(this);
      }
      return handlersCache[stepType];
    },
  };
};

const createWorkflowGraph = (
  factoryContext: FactoryContext,
  workflowDefinition: WorkflowDefinition,
): WorkflowGraph => {
  const nodes: WorkflowGraphNode[] = [];
  const edges: WorkflowGraphEdge[] = [];

  workflowDefinition.steps?.forEach((step) => {
    nodes.push(NodeFactory.createNode(step, factoryContext));
    edges.push(...EdgeFactory.createEdges(step, factoryContext));
  });

  return {
    nodes,
    edges,
  };
};

export const GraphFactory = {
  createFactoryContext,
  createWorkflowGraph,
};
