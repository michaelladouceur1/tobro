import { useSetAtom } from "jotai";
import { MessageType } from "../api/ws_client/types";
import { WebsocketClient } from "../api/ws_client/wsClient";
import { boardAtom } from "../atoms/boardAtom";
import { portsAtom } from "../atoms/portsAtom";
import { useEffect, useState } from "react";

export function useWsApi() {
  const [ws, setWs] = useState<WebSocket | null>(null);
  const setBoard = useSetAtom(boardAtom);
  const setPorts = useSetAtom(portsAtom);

  useEffect(() => {
    const ws = WebsocketClient("ws://localhost:8080/ws", {
      [MessageType.Board]: (data) => {
        const { pins } = data;
        setBoard({ pins });
      },
      [MessageType.Ports]: (data) => {
        const { ports } = data;
        setPorts({ ports });
      },
    });

    setWs(ws);
  }, []);

  return ws;
}
