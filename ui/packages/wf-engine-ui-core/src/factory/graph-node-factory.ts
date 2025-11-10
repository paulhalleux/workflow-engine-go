import { WorkflowStepDefinition } from "@paulhalleux/wf-engine-api";

import { FactoryContext, WorkflowGraphNode } from "../types/graph.ts";

const createNodeBase = (
  definition: WorkflowStepDefinition,
  ctx: FactoryContext,
): WorkflowGraphNode => {
  const handler = ctx.getStepHandler(definition.type);
  const size = handler.getNodeSize(definition);

  return {
    id: definition.stepDefinitionId,
    type: definition.type,
    ...size,
    data: {
      definition,
      handler,
    },
    position: { x: 0, y: 0 },
  };
};

export const NodeFactory = {
  createNode: createNodeBase,
};
