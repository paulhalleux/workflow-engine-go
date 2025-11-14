import { WorkflowStepDefinition } from "@paulhalleux/wf-engine-api";
import { NodeProps, Position } from "@xyflow/react";

import { GraphNode } from "../../components";
import { FactoryContext, WorkflowGraphNode } from "../../types/graph.ts";

/**
 * Creates a function that returns the IDs of the previous steps leading to the given step definition.
 * @param ctx - The factory context containing step definitions and handlers.
 * @returns A function that, when called, returns an array of previous step IDs. The result is memoized for performance.
 */
export const createPreviousStepGetter = (ctx: FactoryContext) => {
  return (definition: WorkflowStepDefinition) => {
    const ids: string[] = [];
    const stepDefinitions = ctx.getStepDefinitions();
    if (stepDefinitions.length === 0) {
      return ids;
    }

    for (const workflowStepDefinition of stepDefinitions) {
      const handler = ctx.getStepHandler(workflowStepDefinition.type);

      const nextStepIds = handler.getNextStepIds(definition);
      if (nextStepIds.includes(definition.stepDefinitionId)) {
        ids.push(workflowStepDefinition.stepDefinitionId);
      }
    }

    return ids;
  };
};

export const defaultNodeRenderer = (node: NodeProps<WorkflowGraphNode>) => {
  return (
    <GraphNode.Root node={node}>
      <GraphNode.Rectangle>{node.data.definition.name}</GraphNode.Rectangle>
      <GraphNode.Handle type="target" position={Position.Left} />
      <GraphNode.Handle type="source" position={Position.Right} />
    </GraphNode.Root>
  );
};
