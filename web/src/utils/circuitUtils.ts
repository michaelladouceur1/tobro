import { SetStateAction } from "jotai";
import { SetupPinResponse } from "../api/http_client";
import { Circuit } from "../types";

export function setCircuitFromSetupPinResponse(setCircuit: (update: SetStateAction<Circuit>) => void, data: SetupPinResponse) {
    setCircuit((prev) => {
        const newCircuit = prev.pins.map((pin) => {
          if (pin.id === data.pin) {
            return { ...pin, mode: data.mode  };
          }
          return pin;
        });
        return { pins: newCircuit };
      });
}