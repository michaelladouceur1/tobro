import { BaseResponse, MessageDataMap, MessageType } from "./types";

export function WebsocketClient<T extends MessageType>(
  baseUrl: string,
  messageMap: { [K in T]: (data: MessageDataMap[K]) => void }
) {
  const ws = new WebSocket(`${baseUrl}`);

  ws.onopen = () => {
    console.log("Connected to server");
  };

  ws.onmessage = (event) => {
    const message: BaseResponse = JSON.parse(event.data);

    if (message.type in messageMap) {
      const handler = messageMap[message.type as T];
      handler(message.data);
    }
  };

  return ws;
}
