import { SetStateAction } from "jotai";
import { CircuitResponse, SetupPinResponse } from "../api/http_client";
import { Circuit, PinMode } from "../types";

export function setCircuitFromSetupPinResponse(setCircuit: (update: SetStateAction<Circuit>) => void, data: SetupPinResponse) {
    setCircuit((prev) => {
        const newCircuit = prev.pins.map((pin) => {
          if (pin.pinNumber === data.pinNumber) {
            return { ...pin, mode: data.mode  };
          }
          return pin;
        });
        return { pins: newCircuit };
      });
}

export function setCircuitFromCircuitResponse(setCircuit: (update: SetStateAction<Circuit>) => void, data: CircuitResponse) {
    setCircuit({pins: data.pins.map((pin) => ({...pin, state: 0, mode: pin.mode === 0 ? PinMode.Input : PinMode.Output}))});
}