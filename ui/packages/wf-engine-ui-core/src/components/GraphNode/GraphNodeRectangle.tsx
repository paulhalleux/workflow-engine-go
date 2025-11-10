import { clsx } from "clsx";
import * as React from "react";

import styles from "./GraphNode.module.css";
import { useGraphNode } from "./GraphNodeContext.tsx";

type GraphNodeRectangleProps = React.PropsWithChildren;

export function GraphNodeRectangle({ children }: GraphNodeRectangleProps) {
  const { node } = useGraphNode();

  return (
    <div
      className={clsx(styles.node, {
        [styles.selected]: node.selected,
      })}
      style={{
        width: node.width,
        height: node.height,
      }}
    >
      {children}
    </div>
  );
}
