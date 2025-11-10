import { NodeProps } from "@xyflow/react";
import * as React from "react";

import { GraphNodeProvider } from "./GraphNodeContext.tsx";

type GraphNodeRootProps = React.PropsWithChildren<{
  node: NodeProps;
}>;

export function GraphNodeRoot({ children, node }: GraphNodeRootProps) {
  return <GraphNodeProvider node={node}>{children}</GraphNodeProvider>;
}
