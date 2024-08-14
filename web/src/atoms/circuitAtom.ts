import { atom } from "jotai";
import { Circuit } from "../types";

export const circuitAtom = atom<Circuit>({
  id: 0,
  name: "",
  board: "",
  pins: [],
});
