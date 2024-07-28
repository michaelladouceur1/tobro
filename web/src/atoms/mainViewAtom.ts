import {atom} from 'jotai'

export const mainViewAtom = atom<"config" | "pin-list">("config")