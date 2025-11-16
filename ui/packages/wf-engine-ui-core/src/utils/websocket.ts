import { BinaryReader } from "@bufbuild/protobuf/wire";

import { MessageFns } from "../types/ws.ts";

export type WebsocketProtocol<Message, Command> = {
  encode: (data: Command) => Uint8Array;
  decode: (data: MessageEvent) => Promise<Message>;
};

type CreateProtobufProtocolOptions<Message, Command> = {
  message: MessageFns<Message>;
  command: MessageFns<Command>;
};

export const createProtobufProtocol = <Message, Command>({
  message,
  command,
}: CreateProtobufProtocolOptions<Message, Command>): WebsocketProtocol<
  Message,
  Command
> => {
  return {
    encode: (data: Command) => {
      return command.encode(data).finish();
    },
    decode: async (data: MessageEvent) => {
      let bytes: Uint8Array;

      if (typeof data.data === "string") {
        throw new Error("Received text message, expected binary data");
      } else if (data.data instanceof ArrayBuffer) {
        bytes = new Uint8Array(data.data);
      } else if (data.data instanceof Blob) {
        const arrayBuffer = await data.data.arrayBuffer();
        bytes = new Uint8Array(arrayBuffer);
      } else {
        throw new Error("Unsupported message data type");
      }

      const reader = new BinaryReader(bytes);
      return message.decode(reader);
    },
  };
};
