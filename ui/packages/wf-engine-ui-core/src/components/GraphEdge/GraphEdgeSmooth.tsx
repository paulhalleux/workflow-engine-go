import { BaseEdge, getSmoothStepPath } from "@xyflow/react";

import styles from "./GraphEdge.module.css";
import { useGraphEdge } from "./GraphEdgeContext.tsx";

export function GraphEdgeSmooth() {
  const { edge } = useGraphEdge();
  const [path] = getSmoothStepPath(edge);

  return <BaseEdge path={path} className={styles.edge} {...edge} />;
}
