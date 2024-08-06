import type { Circuit, Ports } from "../../types";

export enum MessageType {
  Circuit = "circuit",
  Ports = "ports",
  PinMode = "pin_mode",
  PinState = "pin_state",
}

interface CircuitMessage extends Circuit {}

interface PortsMessage extends Ports {}

interface PinModeMessage {
  id: number;
  mode: number;
}

interface PinStateMessage {
  id: number;
  state: number;
}

export type MessageDataMap = {
  [MessageType.Circuit]: CircuitMessage;
  [MessageType.Ports]: PortsMessage;
  [MessageType.PinMode]: PinModeMessage;
  [MessageType.PinState]: PinStateMessage;
};

export interface BaseResponse {
  type: MessageType;
  data: any;
}
