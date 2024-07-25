import { useSetAtom } from "jotai";
import { useEffect, useState } from "react";
import { MessageType } from "../api/ws_client/types";
import { WebsocketClient } from "../api/ws_client/wsClient";
import { boardAtom } from "../atoms/boardAtom";
import { portsAtom } from "../atoms/portsAtom";
import { PinMode } from "../types";

export function useWsApi() {
  const [ws, setWs] = useState<WebSocket | null>(null);

  const setBoard = useSetAtom(boardAtom);
  const setPorts = useSetAtom(portsAtom);

  useEffect(() => {
    let initialized = false;
    const ws = WebsocketClient("ws://localhost:8080/ws", {
      [MessageType.Board]: (data) => {
        const { pins } = data;

        if (!initialized) {
          const state = pins.map((pin) => ({ ...pin, state: 0, mode: "input" }));
          setBoard({ pins: state });
          initialized = true;
          return;
        }

        setBoard({ pins });
      },
      [MessageType.Ports]: (data) => {
        const { ports } = data;
        setPorts({ ports });
      },
      [MessageType.PinState]: (data) => {
        const { id, state } = data;
        setBoard((prev) => {
          const newBoard = prev.pins.map((pin) => {
            if (pin.id === id) {
              return { ...pin, state };
            }
            return pin;
          });
          return { pins: newBoard };
        });
      },
      [MessageType.PinMode]: (data) => {
        const { id, mode } = data;
        setBoard((prev) => {
          const newBoard = prev.pins.map((pin) => {
            if (pin.id === id) {
              return { ...pin, mode: mode === 0 ? PinMode.Input : PinMode.Output };
            }
            return pin;
          });
          return { pins: newBoard };
        });
      },
    });

    setWs(ws);
  }, []);

  return ws;
}
