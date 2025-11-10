import { EdgeProps } from "@xyflow/react";
import * as React from "react";

type GraphEdgeContextType = {
  edge: EdgeProps;
};

const GraphEdgeContext = React.createContext<GraphEdgeContextType | null>(null);

export function GraphEdgeProvider({
  edge,
  children,
}: React.PropsWithChildren<GraphEdgeContextType>) {
  return (
    <GraphEdgeContext.Provider value={React.useMemo(() => ({ edge }), [edge])}>
      {children}
    </GraphEdgeContext.Provider>
  );
}

export function useGraphEdge() {
  const context = React.useContext(GraphEdgeContext);
  if (!context) {
    throw new Error("useGraphEdge must be used within a GraphEdgeProvider");
  }
  return context;
}
