export enum PinMode {
  Input = "input",
  Output = "output",
}

export enum DigitalState {
  Low = 0,
  High = 1,
}

export interface Pin {
  id: number;
  type: string;
  mode: string;
  min: number;
  max: number;
  digitalRead: boolean;
  digitalWrite: boolean;
  analogRead: boolean;
  analogWrite: boolean;
  state: number;
}

export interface PinState {
  id: number;
  state: number;
}

export interface Board {
  pins: Pin[];
}

export interface BoardState {
    [key: number]: number;
}

export interface Ports {
  ports: string[];
}
