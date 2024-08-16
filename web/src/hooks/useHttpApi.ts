import { useSetAtom } from "jotai";
import { Configuration, DefaultApi } from "../api/http_client";
import { circuitAtom } from "../atoms/circuitAtom";
import { DigitalState, Pin, PinMode } from "../types";
import { setCircuitFromCircuitResponse, setCircuitFromSetupPinResponse } from "../utils/circuitUtils";
import { boardsAtom } from "../atoms/boardsAtom";

export function useHttpApi() {
  const setBoards = useSetAtom(boardsAtom);
  const setCircuit = useSetAtom(circuitAtom);

  const configuration = new Configuration({ basePath: "http://localhost:8080" });
  const api = new DefaultApi(configuration);

  return {
    init: async () => {
      const reqs = [
        {
          req: api.boardsGet(),
          set: (res: any) => setBoards(res.boards)
        },
        {
          req: api.circuitGet(),
          set: (res: any) => {
            console.log("Got circuit", res);
            setCircuitFromCircuitResponse(setCircuit, res);
          }
        }
      ]
      
      await Promise.all(reqs.map((r) => r.req.then(r.set)));
    },
    getCircuit: async () => {
      const res = await api.circuitGet();
      setCircuitFromCircuitResponse(setCircuit, res);
      return res;
    },
    createCircuit: async (name: string, board: string) => {
      const res = await api.circuitPost({createCircuitRequest: {name, board}});
      setCircuitFromCircuitResponse(setCircuit, res);
    },
    saveCircuit: async (id: number) => {
      const res = await api.saveCircuitPost({saveCircuitRequest: {id}});
      setCircuitFromCircuitResponse(setCircuit, res);
    },
    getBoards: async () => {
      const res = await api.boardsGet();
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