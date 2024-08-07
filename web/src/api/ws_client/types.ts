import type { Circuit, Ports } from "../../types";

export enum MessageType {
  Circuit = "circuit",
  Ports = "ports",
  PinState = "pin_state",
}

interface CircuitMessage extends Circuit {}

interface PortsMessage extends Ports {}


interface PinStateMessage {
  pinNumber: number;
  state: number;
}

export type MessageDataMap = {
  [MessageType.Circuit]: CircuitMessage;
  [MessageType.Ports]: PortsMessage;
  [MessageType.PinState]: PinStateMessage;
};

export interface BaseResponse {
  type: MessageType;
  data: any;
}
