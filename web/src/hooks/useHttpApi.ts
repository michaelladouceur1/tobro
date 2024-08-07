import { useSetAtom } from "jotai";
import { Configuration, DefaultApi } from "../api/http_client";
import { circuitAtom } from "../atoms/circuitAtom";
import { DigitalState, Pin, PinMode } from "../types";
import { setCircuitFromCircuitResponse, setCircuitFromSetupPinResponse } from "../utils/circuitUtils";

export function useHttpApi() {
  const setCircuit = useSetAtom(circuitAtom);

  const configuration = new Configuration({ basePath: "http://localhost:8080" });
  const api = new DefaultApi(configuration);

  return {
    init: async () => {
      const circuitRes = await api.circuitGet();
      setCircuitFromCircuitResponse(setCircuit, circuitRes);
    },
    getCircuit: async () => {
      const res = await api.circuitGet();
      setCircuitFromCircuitResponse(setCircuit, res);
      return res;
    },
    connect: async (port: string) => {
      if (!port) return
      const res = await api.connectPost({ connectRequest: { port } });
      return res;
    },
    setupPin: async (pin: Pin) => {
      try {
        const {pinNumber} = pin;
        const mode = pin.mode === PinMode.Output ? PinMode.Input : PinMode.Output;
        const res = await api.setupPinPost({setupPinRequest: {pinNumber: pinNumber, mode}});
        setCircuitFromSetupPinResponse(setCircuit, res);
        return res;
      } catch (error) {
        console.log(error);
      }
    },
    digitalWritePin: async (pin: Pin) => {
      const {pinNumber} = pin;
      const value = pin.state === pin.max ? DigitalState.Low : DigitalState.High;
      const res = await api.digitalWritePinPost({digitalWritePinRequest: {pinNumber: pinNumber, value}});
      return res;
    }
  };
}