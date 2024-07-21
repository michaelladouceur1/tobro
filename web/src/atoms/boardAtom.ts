import {atom} from "jotai";
import {Board} from "./types";

export const boardAtom = atom<Board>({
    digitalPins: [],
    analogPins: []
});