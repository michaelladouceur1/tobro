import {atom} from 'jotai';
import {PortConnection} from '../types';

export const portConnectionAtom = atom<PortConnection>({connected: false, portName: ''});