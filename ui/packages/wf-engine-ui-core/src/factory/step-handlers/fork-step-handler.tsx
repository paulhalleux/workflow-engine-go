import { StepType, WorkflowStepDefinition } from "@paulhalleux/wf-engine-api";
import invariant from "tiny-invariant";

import { StepHandlerFactory } from "../../types/handlers.ts";
import { createPreviousStepGetter, defaultNodeRenderer } from "./helpers.tsx";

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
        return { width: 48, height: 48 };
      },
      getNextStepIds: (definition) => {
        const config = getConfig(definition);
        return config.branches.map((branch) => branch.nextStepId);
      },
      getPreviousStepIds: createPreviousStepGetter(ctx),
      ...overrides,
      render: (props) => {
        const Render = overrides?.render ?? defaultNodeRenderer;
        return <Render {...props} />;
      },
    };
  };
};
