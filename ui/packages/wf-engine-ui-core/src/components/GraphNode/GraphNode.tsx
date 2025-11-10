import { Handle } from "@xyflow/react";

import { GraphNodeCircle } from "./GraphNodeCircle.tsx";
import { GraphNodeRectangle } from "./GraphNodeRectangle.tsx";
import { GraphNodeRoot } from "./GraphNodeRoot.tsx";

export const GraphNode = {
  Root: GraphNodeRoot,
  Rectangle: GraphNodeRectangle,
  Circle: GraphNodeCircle,
  Handle,
};
