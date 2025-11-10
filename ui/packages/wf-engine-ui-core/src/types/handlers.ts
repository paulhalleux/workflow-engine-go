import { StepType, WorkflowStepDefinition } from "@paulhalleux/wf-engine-api";
import { NodeProps } from "@xyflow/react";
import * as React from "react";

import { FactoryContext } from "./graph.ts";

export type Size = {
  width: number;
  height: number;
};

export type StepHandler = {
  stepType: StepType;
  getNodeSize: (definition: WorkflowStepDefinition) => Size;
  getNextStepIds: (definition: WorkflowStepDefinition) => string[];
  getPreviousStepIds: (definition: WorkflowStepDefinition) => string[];
  render: React.ComponentType<NodeProps>;
};

export type StepHandlerBaseFactory = (ctx: FactoryContext) => StepHandler;
export type StepHandlerFactory = (
  overrides?: (ctx: FactoryContext) => Partial<StepHandler>,
) => StepHandlerBaseFactory;
