import { clsx } from "clsx";
import * as React from "react";

import styles from "./GraphNode.module.css";
import { useGraphNode } from "./GraphNodeContext.tsx";

type GraphNodeCircleProps = React.PropsWithChildren;

export function GraphNodeCircle({ children }: GraphNodeCircleProps) {
  const { node } = useGraphNode();

  return (
    <div
      className={clsx(styles.node, styles.circle, {
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
