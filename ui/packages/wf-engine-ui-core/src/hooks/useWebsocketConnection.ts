import { BinaryReader } from "@bufbuild/protobuf/wire";
import * as React from "react";

import { MessageFns } from "../types/ws.ts";

export const useWebsocketConnection = <Command, Message>(
  url: string,
  command: MessageFns<Command>,
  message: MessageFns<Message>,
  onMessage: (data: Message) => void,
) => {
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
      let bytes: Uint8Array;

      if (typeof event.data === "string") {
        console.warn("Received text message, ignoring binary decoder");
        return;
      } else if (event.data instanceof ArrayBuffer) {
        bytes = new Uint8Array(event.data);
      } else if (event.data instanceof Blob) {
        const buffer = await event.data.arrayBuffer();
        bytes = new Uint8Array(buffer);
      } else {
        console.warn("Unknown WebSocket message type", typeof event.data);
        return;
      }

      const reader = new BinaryReader(bytes);
      const decodedMessage = message.decode(reader);
      latestOnMessage.current(decodedMessage);
    };

    return () => {
      ws.current?.close();
    };
  }, [url, message]);

  const sendMessage = React.useCallback(
    (data: Command) => {
      if (ws.current && ws.current.readyState === WebSocket.OPEN) {
        const encodedMessage = command.encode(data).finish();
        ws.current.send(encodedMessage);
        console.log("WebSocket message sent:", data);
      } else {
        console.warn("WebSocket is not open. Unable to send message.");
      }
    },
    [command],
  );

  return { sendMessage };
};
