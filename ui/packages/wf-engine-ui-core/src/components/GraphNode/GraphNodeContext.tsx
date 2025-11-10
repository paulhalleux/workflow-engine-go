import { NodeProps } from "@xyflow/react";
import * as React from "react";

type GraphNodeContextType = {
  node: NodeProps;
};

const GraphNodeContext = React.createContext<GraphNodeContextType | null>(null);

export function GraphNodeProvider({
  node,
  children,
}: React.PropsWithChildren<GraphNodeContextType>) {
  return (
    <GraphNodeContext.Provider value={React.useMemo(() => ({ node }), [node])}>
      {children}
    </GraphNodeContext.Provider>
  );
}

export function useGraphNode() {
  const context = React.useContext(GraphNodeContext);
  if (!context) {
    throw new Error("useGraphNode must be used within a GraphNodeProvider");
  }
  return context;
}
