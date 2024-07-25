import List from "@mui/material/List";
import ListItem from "@mui/material/ListItem";
import ListItemText from "@mui/material/ListItemText";
// import ListItemAvatar from "@mui/material/ListItemAvatar";
import {Button, Divider, Switch} from "@mui/material";
import {useAtomValue} from "jotai";
import {boardAtom} from "../../atoms/boardAtom";
import {useHttpApi} from "../../hooks/useHttpApi";
import {DigitalState, PinMode} from "../../types";

export function PinList() {
  const api = useHttpApi();
  const board = useAtomValue(boardAtom);

  const handleSetupPin = async (id: number) => {
    await api.setupPinPost({
      setupPinRequest: {pin: id, mode: PinMode.Output},
    });
  };

  const handleDigitalWrite = async (id: number, value: number) => {
    await api.digitalWritePinPost({
      digitalWritePinRequest: {pin: id, value},
    });
  };

  return (
    <div>
      <h2>Board</h2>
      <List dense={true}>
        {board.pins.map((pin) => (
          <>
            <ListItem key={pin.id}>
              <ListItemText primary={pin.id} secondary={pin.type} />
              <Button
                variant="contained"
                onClick={async () => {
                  const {id} = pin;
                  await api.setupPinPost({
                    setupPinRequest: {pin: id, mode: PinMode.Output},
                  });
                }}
              >
                Setup
              </Button>
              <Switch
                checked={pin.state === pin.max}
                onChange={async () => {
                  const {id} = pin;
                  const value =
                    pin.state === pin.max
                      ? DigitalState.Low
                      : DigitalState.High;
                  await api.digitalWritePinPost({
                    digitalWritePinRequest: {pin: id, value},
                  });
                }}
              />
              {/* <Button
                variant="contained"
                color="primary"
                onClick={() => {
                  const { id } = pin;
                  console.log("High", id, typeof id);
                  api.digitalWritePinPost({ digitalWritePinRequest: { pin: id, value: 1 } });
                }}
              >
                High
              </Button>
              <Button
                variant="contained"
                color="secondary"
                onClick={() => {
                  const { id } = pin;
                  api.digitalWritePinPost({ digitalWritePinRequest: { pin: id, value: 0 } });
                }}
              >
                Low
              </Button> */}
            </ListItem>
            <Divider />
          </>
        ))}
      </List>
    </div>
  );
}
