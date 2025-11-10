import { StepType, WorkflowStepDefinition } from "@paulhalleux/wf-engine-api";
import invariant from "tiny-invariant";

import { StepHandlerFactory } from "../../types/handlers.ts";
import { createPreviousStepGetter, defaultNodeRenderer } from "./helpers.tsx";

export const createTaskStepHandler: StepHandlerFactory = (overridesFactory) => {
  return (ctx) => {
    const getConfig = (definition: WorkflowStepDefinition) => {
      const config = definition.taskConfig;
      invariant(config, "taskConfig must be defined for TaskStepHandler");
      return config;
    };

    const overrides = overridesFactory?.(ctx);
    return {
      stepType: StepType.StepTypeTask,
      getNodeSize: () => {
        return { width: 120, height: 60 };
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
