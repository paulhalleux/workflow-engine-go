import { StepType, WorkflowStepDefinition } from "@paulhalleux/wf-engine-api";
import invariant from "tiny-invariant";

import { StepHandlerFactory } from "../../types/handlers.ts";
import { createPreviousStepGetter, defaultNodeRenderer } from "./helpers.tsx";

export const createJoinStepHandler: StepHandlerFactory = (overridesFactory) => {
  return (ctx) => {
    const getConfig = (definition: WorkflowStepDefinition) => {
      const config = definition.joinConfig;
      invariant(config, "joinConfig must be defined for JoinStepHandler");
      return config;
    };

    const overrides = overridesFactory?.(ctx);
    return {
      stepType: StepType.StepTypeJoin,
      getNodeSize: () => {
        return { width: 48, height: 48 };
      },
      getNextStepIds: (definition) => {
        const config = getConfig(definition);
        return config.nextStepId ? [config.nextStepId] : [];
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
