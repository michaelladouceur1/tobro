import { useSetAtom } from "jotai";
import { DefaultApi, Configuration } from "../api/http_client";
import { DigitalState, Pin, PinMode } from "../types";
import { circuitAtom } from "../atoms/circuitAtom";

export function useHttpApi() {
  const setCircuit = useSetAtom(circuitAtom);

  const configuration = new Configuration({ basePath: "http://localhost:8080" });
  const api = new DefaultApi(configuration);

  const handlers = {
    connect: async (port: string) => {
      if (!port) return
      const res = await api.connectPost({ connectRequest: { port } });
    },
    setupPin: async (pin: Pin) => {
      try {
        const {id} = pin;
        const mode = pin.mode === PinMode.Output ? PinMode.Input : PinMode.Output;
        const res = await api.setupPinPost({setupPinRequest: {pin: id, mode}});

        // TODO: move to atom file
        setCircuit((prev) => {
          const newCircuit = prev.pins.map((pin) => {
            if (pin.id === res.pin) {
              return { ...pin, mode: res.mode  };
            }
            return pin;
          });
          return { pins: newCircuit };
        });
      } catch (error) {
        console.log(error);
      }
    },
    digitalWritePin: async (pin: Pin) => {
      const {id} = pin;
      const value = pin.state === pin.max ? DigitalState.Low : DigitalState.High;
      await api.digitalWritePinPost({
        digitalWritePinRequest: {pin: id, value},
      });
    }
  }

  return handlers;
}