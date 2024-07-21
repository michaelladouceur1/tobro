export interface Pin {
    id: number;
    pinType: string;
    mode: string;
    state: number;
    min: number;
    max: number;
}

export interface DigitalPin extends Pin {
    pwm: boolean;
}

export interface AnalogPin extends Pin {}

export interface Board {
    digitalPins: DigitalPin[];
    analogPins: AnalogPin[];
}

export interface Ports {
    ports: string[];
}