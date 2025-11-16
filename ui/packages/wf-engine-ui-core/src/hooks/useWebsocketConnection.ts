import * as React from "react";

import { WebsocketProtocol } from "../utils/websocket.ts";

type UseWebsocketConnectionArguments<Message, Command> = {
  url: string;
  protocol: WebsocketProtocol<Message, Command>;
  onMessage: (data: Message) => void;
};

export const useWebsocketConnection = <Message, Command>({
  url,
  protocol,
  onMessage,
}: UseWebsocketConnectionArguments<Message, Command>) => {
  const ws = React.useRef<WebSocket | null>(null);

  const latestOnMessage = React.useRef(onMessage);
  React.useEffect(() => {
    latestOnMessage.current = onMessage;
  }, [onMessage]);

  React.useEffect(() => {
    ws.current = new WebSocket(url);

    ws.current.onopen = () => {
      console.log("WebSocket connected");
    };

    ws.current.onclose = () => {
      console.log("WebSocket disconnected");
    };

    ws.current.onmessage = async (event) => {
      const decodedMessage = protocol.decode(event);
      latestOnMessage.current(decodedMessage);
    };

    return () => {
      ws.current?.close();
    };
  }, [protocol, url]);

  const sendMessage = React.useCallback(
    (data: Command) => {
      if (!(ws.current && ws.current.readyState === WebSocket.OPEN)) {
        return;
      }

      const encodedMessage = protocol.encode(data);
      ws.current.send(encodedMessage);
    },
    [protocol],
  );

  return { sendMessage };
};
