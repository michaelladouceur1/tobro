import type { Board, Ports } from "../../types";

export enum MessageType {
  Board = "board",
  Ports = "ports",
}

export type MessageDataMap = {
  [MessageType.Board]: Board;
  [MessageType.Ports]: Ports;
};

export interface BaseResponse {
  type: MessageType;
  data: any;
}
