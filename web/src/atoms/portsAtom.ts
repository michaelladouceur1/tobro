import { atom } from "jotai";
import { Ports } from "../types";

export const portsAtom = atom<Ports>({ ports: [] });
