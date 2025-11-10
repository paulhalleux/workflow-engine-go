import { WorkflowStepDefinition } from "@paulhalleux/wf-engine-api";

import { FactoryContext, WorkflowGraphEdge } from "../types/graph.ts";

const createEdgeBase = (
  stepId: string,
  nextStepId: string,
): WorkflowGraphEdge => ({
  type: "default",
  id: `${stepId}-to-${nextStepId}`,
  source: stepId,
  target: nextStepId,
});

export const EdgeFactory = {
  createEdges: (
    step: WorkflowStepDefinition,
    ctx: FactoryContext,
  ): WorkflowGraphEdge[] => {
    const handler = ctx.getStepHandler(step.type);
    return handler
      .getNextStepIds(step)
      .map((nextStepId) => createEdgeBase(step.stepDefinitionId, nextStepId));
  },
};
