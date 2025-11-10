import { StepType, WorkflowStepDefinition } from "@paulhalleux/wf-engine-api";
import invariant from "tiny-invariant";

import { StepHandlerFactory } from "../../types/handlers.ts";
import { createPreviousStepGetter, defaultNodeRenderer } from "./helpers.tsx";

export const createWorkflowStepHandler: StepHandlerFactory = (
  overridesFactory,
) => {
  return (ctx) => {
    const getConfig = (definition: WorkflowStepDefinition) => {
      const config = definition.workflowConfig;
      invariant(
        config,
        "workflowConfig must be defined for WorkflowStepHandler",
      );
      return config;
    };

    const overrides = overridesFactory?.(ctx);
    return {
      stepType: StepType.StepTypeWorkflow,
      getNodeSize: () => {
        return { width: 140, height: 80 };
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
