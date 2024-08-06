import { useSetAtom } from "jotai";
import { useEffect, useState } from "react";
import { MessageType } from "../api/ws_client/types";
import { WebsocketClient } from "../api/ws_client/wsClient";
import { circuitAtom } from "../atoms/circuitAtom";
import { portsAtom } from "../atoms/portsAtom";
import { PinMode } from "../types";

export function useWsApi() {
  const [ws, setWs] = useState<WebSocket | null>(null);

  const setCircuit = useSetAtom(circuitAtom);
  const setPorts = useSetAtom(portsAtom);

  useEffect(() => {
    let initialized = false;
    const ws = WebsocketClient("ws://localhost:8080/ws", {
      [MessageType.Circuit]: (data) => {
        const { pins } = data;

        if (!initialized) {
          const state = pins.map((pin) => ({ ...pin, state: 0, mode: "input" }));
          setCircuit({ pins: state });
          initialized = true;
          return;
        }

        setCircuit({ pins });
      },
      [MessageType.Ports]: (data) => {
        const { ports } = data;
        setPorts({ ports });
      },
      [MessageType.PinState]: (data) => {
        const { id, state } = data;
        setCircuit((prev) => {
          const newCircuit = prev.pins.map((pin) => {
            if (pin.id === id) {
              return { ...pin, state };
            }
            return pin;
          });
          return { pins: newCircuit };
        });
      },
      [MessageType.PinMode]: (data) => {
        console.log("PinMode", data);
        const { id, mode } = data;
        setCircuit((prev) => {
          const newCircuit = prev.pins.map((pin) => {
            if (pin.id === id) {
              return { ...pin, mode: mode === 0 ? PinMode.Input : PinMode.Output };
            }
            return pin;
          });
          return { pins: newCircuit };
        });
      },
    });

    setWs(ws);
  }, []);

  return ws;
}
