import { StepType, WorkflowStepDefinition } from "@paulhalleux/wf-engine-api";
import invariant from "tiny-invariant";

import { StepHandlerFactory } from "../../types/handlers.ts";
import { createPreviousStepGetter, defaultNodeRenderer } from "./helpers.tsx";

export const createDecisionStepHandler: StepHandlerFactory = (
  overridesFactory,
) => {
  return (ctx) => {
    const getConfig = (definition: WorkflowStepDefinition) => {
      const config = definition.decisionConfig;
      invariant(
        config,
        "decisionConfig must be defined for DecisionStepHandler",
      );
      return config;
    };

    const overrides = overridesFactory?.(ctx);
    return {
      stepType: StepType.StepTypeDecision,
      getNodeSize: () => {
        return { width: 48, height: 48 };
      },
      getNextStepIds: (definition) => {
        const config = getConfig(definition);
        return config.cases.map((branch) => branch.nextStepId);
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
