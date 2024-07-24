import Box from "@mui/material/Box";
import List from "@mui/material/List";
import ListItem from "@mui/material/ListItem";
import ListItemText from "@mui/material/ListItemText";
// import ListItemAvatar from "@mui/material/ListItemAvatar";
import { useAtomValue } from "jotai";
import { boardAtom } from "../../atoms/boardAtom";
import { Button, Divider } from "@mui/material";
import { useHttpApi } from "../../hooks/useHttpApi";

export function PinList() {
  const api = useHttpApi();
  const board = useAtomValue(boardAtom);

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
                  const { id } = pin;
                  await api.setupPinPost({ setupPinRequest: { pin: id, mode: "output" } });
                }}
              >
                Setup
              </Button>
              <Button
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
              </Button>
            </ListItem>
            <Divider />
          </>
        ))}
      </List>
    </div>
  );
}
