import { BaseEdge, getBezierPath } from "@xyflow/react";

import styles from "./GraphEdge.module.css";
import { useGraphEdge } from "./GraphEdgeContext.tsx";

export function GraphEdgeBezier() {
  const { edge } = useGraphEdge();
  const [path] = getBezierPath(edge);

  return <BaseEdge path={path} className={styles.edge} {...edge} />;
}
