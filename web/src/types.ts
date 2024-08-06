export enum PinMode {
  Input = "input",
  Output = "output",
}

export enum PinType {
  Digital = "digital",
  Analog = "analog",
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

export interface Circuit {
  pins: Pin[];
}

export interface CircuitState {
  [key: number]: number;
}

export interface Ports {
  ports: string[];
}
