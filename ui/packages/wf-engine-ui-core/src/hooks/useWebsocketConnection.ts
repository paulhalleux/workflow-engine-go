import * as React from "react";

import { WebsocketProtocol } from "../utils/websocket.ts";

type UseWebsocketConnectionArguments<Message, Command> = {
  url: string;
  protocol: WebsocketProtocol<Message, Command>;
  onMessage?: (data: Message) => void;
  onOpen?: () => void;
  onClose?: () => void;
  onError?: (error: Event) => void;
};

export const useWebsocketConnection = <Message, Command>({
  url,
  protocol,
  onMessage,
  onOpen,
  onClose,
  onError,
}: UseWebsocketConnectionArguments<Message, Command>) => {
  const ws = React.useRef<WebSocket | null>(null);

  const latestOnMessage = React.useRef(onMessage);
  const latestOnClose = React.useRef(onClose);
  const latestOnOpen = React.useRef(onOpen);
  const latestOnError = React.useRef(onError);

  React.useEffect(() => {
    latestOnMessage.current = onMessage;
    latestOnClose.current = onClose;
    latestOnOpen.current = onOpen;
    latestOnError.current = onError;
  }, [onMessage, onClose, onOpen, onError]);

  React.useEffect(() => {
    ws.current = new WebSocket(url);

    ws.current.onopen = () => {
      latestOnOpen.current?.();
    };

    ws.current.onclose = () => {
      latestOnClose.current?.();
    };

    ws.current.onerror = (event) => {
      latestOnError.current?.(event);
    };

    ws.current.onmessage = async (event) => {
      const decodedMessage = await protocol.decode(event);
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
