import { EdgeProps } from "@xyflow/react";
import * as React from "react";

import { GraphEdgeProvider } from "./GraphEdgeContext.tsx";

type GraphEdgeRootProps = React.PropsWithChildren<{
  edge: EdgeProps;
}>;

export function GraphEdgeRoot({ edge, children }: GraphEdgeRootProps) {
  return <GraphEdgeProvider edge={edge}>{children}</GraphEdgeProvider>;
}
