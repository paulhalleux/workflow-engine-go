import {
  Configuration,
  WorkflowDefinitionsApi,
} from "@paulhalleux/wf-engine-api";
import { queryOptions } from "@tanstack/react-query";

const config = new Configuration({
  basePath: "http://localhost:8080",
});

const api = new WorkflowDefinitionsApi(config);

const getAll = () => {
  return queryOptions({
    queryKey: ["workflow-definitions", "all"],
    queryFn: async () => {
      return await api
        .getAllWorkflowDefinitions()
        .catch((err) => console.log(err));
    },
  });
};

const getById = (id: string) => {
  return queryOptions({
    queryKey: ["workflow-definitions", "by-id", id],
    queryFn: async () => {
      return await api.getWorkflowDefinitionByID({ id });
    },
  });
};

export const WorkflowDefinitionsQuery = {
  getAll,
  getById,
};
