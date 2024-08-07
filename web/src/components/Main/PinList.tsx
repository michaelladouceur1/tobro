import {Divider, ListItemIcon, Switch} from "@mui/material";
import List from "@mui/material/List";
import ListItem from "@mui/material/ListItem";
import ListItemText from "@mui/material/ListItemText";
import {useAtomValue} from "jotai";
import {PiWaveSineLight, PiWaveSquareLight} from "react-icons/pi";
import {circuitAtom} from "../../atoms/circuitAtom";
import {useHttpApi} from "../../hooks/useHttpApi";
import {PinMode, PinType} from "../../types";

export function PinList() {
  const api = useHttpApi();
  const circuit = useAtomValue(circuitAtom);

  return (
    // <Card variant="outlined" sx={{gridArea: "main", overflowY: "auto"}}>
    <List dense={true}>
      {circuit.pins.map((pin) => (
        <>
          <ListItem key={pin.pinNumber}>
            <ListItemIcon>
              {pin.type === PinType.Digital ? (
                <PiWaveSquareLight size="20px" />
              ) : (
                <PiWaveSineLight size="20px" />
              )}
            </ListItemIcon>
            <ListItemText primary={pin.pinNumber} secondary={pin.type} />
            <Switch
              checked={pin.mode === PinMode.Output}
              onChange={() => api.setupPin(pin)}
            />
            <Switch
              checked={pin.state === pin.max}
              onChange={() => api.digitalWritePin(pin)}
            />
          </ListItem>
          <Divider />
        </>
      ))}
    </List>
    // </Card>
  );
}
