import dagre from "@dagrejs/dagre";
import { StepType, WorkflowDefinition } from "@paulhalleux/wf-engine-api";
import {
  applyNodeChanges,
  Background,
  EdgeProps,
  Position,
  ReactFlow,
} from "@xyflow/react";
import * as React from "react";

import { GraphEdge } from "../components";
import { GraphFactory } from "../factory/graph-factory.ts";
import { WorkflowGraph, WorkflowGraphNode } from "../types/graph.ts";
import { StepHandlerBaseFactory } from "../types/handlers.ts";
import styles from "./app.module.css";

export type WorkflowDefinitionGraphProps = {
  handlers: Record<StepType, StepHandlerBaseFactory>;
  workflowDefinition: WorkflowDefinition;
};

export function WorkflowDefinitionGraph({
  handlers,
  workflowDefinition,
}: WorkflowDefinitionGraphProps) {
  const factoryContext = React.useMemo(() => {
    return GraphFactory.createFactoryContext(handlers, workflowDefinition);
  }, [handlers, workflowDefinition]);

  const [graph, setGraph] = React.useState(() => {
    return getLayoutedElements(
      GraphFactory.createWorkflowGraph(factoryContext, workflowDefinition),
    );
  });

  return (
    <ReactFlow
      className={styles.reactflow}
      nodes={graph.nodes}
      edges={graph.edges}
      nodeTypes={Object.fromEntries(
        Object.values(StepType).map((stepType) => [
          stepType,
          factoryContext.getStepHandler(stepType).render,
        ]),
      )}
      edgeTypes={{
        default: EdgeRenderer,
      }}
      fitView
      onNodesChange={(changes) =>
        setGraph((prev) => ({
          ...prev,
          nodes: applyNodeChanges(changes, prev.nodes),
        }))
      }
    >
      <Background />
    </ReactFlow>
  );
}

// const NodeRenderer = (node: NodeProps) => {
//   return (
//     <GraphNode.Root node={node}>
//       <GraphNode.Rectangle />
//       <GraphNode.Handle type="target" position={Position.Left} />
//       <GraphNode.Handle type="source" position={Position.Right} />
//     </GraphNode.Root>
//   );
// };

const EdgeRenderer = (edge: EdgeProps) => {
  return (
    <GraphEdge.Root edge={edge}>
      <GraphEdge.Bezier />
    </GraphEdge.Root>
  );
};

const getLayoutedElements = (graph: WorkflowGraph): WorkflowGraph => {
  const dagreGraph = new dagre.graphlib.Graph().setDefaultEdgeLabel(() => ({}));
  const isHorizontal = true;

  dagreGraph.setGraph({ rankdir: "LR" });

  graph.nodes.forEach((node) => {
    dagreGraph.setNode(node.id, { width: node.width, height: node.height });
  });

  graph.edges.forEach((edge) => {
    dagreGraph.setEdge(edge.source, edge.target);
  });

  dagre.layout(dagreGraph);

  const newNodes = graph.nodes.map<WorkflowGraphNode>((node) => {
    const nodeWithPosition = dagreGraph.node(node.id);
    return {
      ...node,
      targetPosition: isHorizontal ? Position.Left : Position.Top,
      sourcePosition: isHorizontal ? Position.Right : Position.Bottom,
      position: {
        x: nodeWithPosition.x - node.width / 2,
        y: nodeWithPosition.y - node.height / 2,
      },
    };
  });

  return { nodes: newNodes, edges: graph.edges };
};
