import {Divider, ListItemIcon, Switch} from "@mui/material";
import List from "@mui/material/List";
import ListItem from "@mui/material/ListItem";
import ListItemText from "@mui/material/ListItemText";
import {useAtomValue} from "jotai";
import {PiWaveSineLight, PiWaveSquareLight} from "react-icons/pi";
import {circuitAtom} from "../../atoms/circuitAtom";
import {useHttpApi} from "../../hooks/useHttpApi";
import {DigitalState, Pin, PinMode, PinType} from "../../types";

export function PinList() {
  const api = useHttpApi();
  const circuit = useAtomValue(circuitAtom);

  const handleSetupPin = async (pin: Pin) => {
    const {id} = pin;
    const mode = pin.mode === PinMode.Output ? PinMode.Input : PinMode.Output;
    await api.setupPinPost({setupPinRequest: {pin: id, mode}});
  };

  const handleDigitalWrite = async (pin: Pin) => {
    const {id} = pin;
    const value = pin.state === pin.max ? DigitalState.Low : DigitalState.High;
    await api.digitalWritePinPost({
      digitalWritePinRequest: {pin: id, value},
    });
  };

  return (
    // <Card variant="outlined" sx={{gridArea: "main", overflowY: "auto"}}>
    <List dense={true}>
      {circuit.pins.map((pin) => (
        <>
          <ListItem key={pin.id}>
            <ListItemIcon>
              {pin.type === PinType.Digital ? (
                <PiWaveSquareLight size="20px" />
              ) : (
                <PiWaveSineLight size="20px" />
              )}
            </ListItemIcon>
            <ListItemText primary={pin.id} secondary={pin.type} />
            <Switch
              checked={pin.mode === PinMode.Output}
              onChange={() => handleSetupPin(pin)}
            />
            <Switch
              checked={pin.state === pin.max}
              onChange={() => handleDigitalWrite(pin)}
            />
          </ListItem>
          <Divider />
        </>
      ))}
    </List>
    // </Card>
  );
}
