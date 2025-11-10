import { GraphEdgeBezier } from "./GraphEdgeBezier.tsx";
import { GraphEdgeRoot } from "./GraphEdgeRoot.tsx";
import { GraphEdgeSmooth } from "./GraphEdgeSmooth.tsx";

export const GraphEdge = {
  Root: GraphEdgeRoot,
  Smooth: GraphEdgeSmooth,
  Bezier: GraphEdgeBezier,
};
