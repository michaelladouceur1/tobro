import { useSetAtom } from "jotai";
import { useEffect, useState } from "react";
import { MessageType } from "../api/ws_client/types";
import { WebsocketClient } from "../api/ws_client/wsClient";
import { circuitAtom } from "../atoms/circuitAtom";
import { portConnectionAtom } from "../atoms/portConnection";
import { portsAtom } from "../atoms/portsAtom";

export function useWsApi() {
  const [ws, setWs] = useState<WebSocket | null>(null);

  const setCircuit = useSetAtom(circuitAtom);
  const setPorts = useSetAtom(portsAtom);
  const setPortConnection = useSetAtom(portConnectionAtom);

  useEffect(() => {
    const ws = WebsocketClient("ws://localhost:8080/ws", {
      [MessageType.Ports]: (data) => {
        const { ports } = data;
        setPorts({ ports });
      },
      [MessageType.PortConnection]: (data) => {
        const { connected, portName } = data;
        setPortConnection({ connected, portName });
      },
      [MessageType.PinState]: (data) => {
        const { pinNumber, state } = data;
        setCircuit((prev) => {
          const newCircuit = prev.pins.map((pin) => {
            if (pin.pinNumber === pinNumber) {
              return { ...pin, state };
            }
            return pin;
          });
          return { ...prev, pins: newCircuit };
        });
      },
    });

    setWs(ws);
  }, []);

  return ws;
}
