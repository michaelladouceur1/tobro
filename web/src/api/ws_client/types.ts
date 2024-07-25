import type { Board, PinState, Ports } from "../../types";

export enum MessageType {
  Board = "board",
  Ports = "ports",
  PinState = "pin_state",
}

export type MessageDataMap = {
  [MessageType.Board]: Board;
  [MessageType.Ports]: Ports;
  [MessageType.PinState]: PinState;
};

export interface BaseResponse {
  type: MessageType;
  data: any;
}
