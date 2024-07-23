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
}

export interface Board {
  pins: Pin[];
}

export interface Ports {
  ports: string[];
}
