import { StepType, WorkflowStepDefinition } from "@paulhalleux/wf-engine-api";
import { Position } from "@xyflow/system";
import { DynamicIcon } from "lucide-react/dynamic";
import invariant from "tiny-invariant";

import { GraphNode } from "../../components";
import { StepHandlerFactory } from "../../types/handlers.ts";
import { createPreviousStepGetter } from "./helpers.tsx";

export const createForkStepHandler: StepHandlerFactory = (overridesFactory) => {
  return (ctx) => {
    const getConfig = (definition: WorkflowStepDefinition) => {
      const config = definition.forkConfig;
      invariant(config, "forkConfig must be defined for ForkStepHandler");
      return config;
    };

    const overrides = overridesFactory?.(ctx);
    return {
      stepType: StepType.StepTypeFork,
      getNodeSize: () => {
        return { width: 60, height: 60 };
      },
      getNextStepIds: (definition) => {
        const config = getConfig(definition);
        return config.branches.map((branch) => branch.nextStepId);
      },
      getPreviousStepIds: createPreviousStepGetter(ctx),
      ...overrides,
      render: (node) => {
        if (overrides?.render) {
          return <overrides.render {...node} />;
        }

        return (
          <GraphNode.Root node={node}>
            <GraphNode.Circle>
              <DynamicIcon name="git-fork" />
            </GraphNode.Circle>
            <GraphNode.Handle type="target" position={Position.Left} />
            <GraphNode.Handle type="source" position={Position.Right} />
          </GraphNode.Root>
        );
      },
    };
  };
};
