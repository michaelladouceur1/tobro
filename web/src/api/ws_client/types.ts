import type { Board, Ports } from "../../types";

export enum MessageType {
  Board = "board",
  Ports = "ports",
  PinMode = "pin_mode",
  PinState = "pin_state",
}

interface BoardMessage extends Board {}

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
  [MessageType.Board]: BoardMessage;
  [MessageType.Ports]: PortsMessage;
  [MessageType.PinMode]: PinModeMessage;
  [MessageType.PinState]: PinStateMessage;
};

export interface BaseResponse {
  type: MessageType;
  data: any;
}
