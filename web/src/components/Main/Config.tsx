import {
  Avatar,
  Box,
  Button,
  Card,
  CardContent,
  CardHeader,
  Divider,
  FormControl,
  IconButton,
  InputLabel,
  List,
  ListItem,
  ListItemAvatar,
  ListItemText,
  MenuItem,
  Modal,
  Select,
  SpeedDial,
  SpeedDialAction,
  SpeedDialActionProps,
  SpeedDialIcon,
  Stack,
  styled,
  Switch,
  TextField,
} from "@mui/material";
import {useAtomValue} from "jotai";
import {memo, useCallback, useEffect, useMemo, useState} from "react";
import {FaPlus, FaRegFolder} from "react-icons/fa";
import {PiWaveSineLight, PiWaveSquareLight} from "react-icons/pi";
import {boardsAtom} from "../../atoms/boardsAtom";
import {circuitAtom} from "../../atoms/circuitAtom";
import {useHttpApi} from "../../hooks/useHttpApi";
import {theme} from "../../theme";
import {PinMode, PinType} from "../../types";

function AddCircuitModal({
  open,
  close,
  boards,
}: {
  open: boolean;
  close: () => void;
  boards: string[];
}) {
  const api = useHttpApi();
  const [newCircuit, setNewCircuit] = useState({name: "", board: ""});

  const updateCircuitField = (field: string, value: string) => {
    setNewCircuit((prev) => ({...prev, [field]: value}));
  };

  const handleAdd = async () => {
    console.log("Adding new circuit");
    const res = await api.createCircuit(newCircuit.name, newCircuit.board);
    console.log(res);
    close();
  };

  return (
    <Modal open={open} onClose={close}>
      <Card
        sx={{
          position: "absolute",
          top: "50%",
          left: "50%",
          transform: "translate(-50%, -50%)",
          width: "500px",
          height: "fit-content",
          border: "2px solid #000",
          boxShadow: 24,
        }}
      >
        <CardHeader
          title="New Circuit"
          action={
            <IconButton color="primary" onClick={handleAdd}>
              <FaPlus />
            </IconButton>
          }
          sx={{borderBottom: "1px solid #000"}}
        />
        <CardContent>
          <FormControl
            fullWidth
            sx={{display: "flex", flexDirection: "column", rowGap: "10px"}}
          >
            <InputLabel required id="board-select-label">
              Board
            </InputLabel>
            <Select
              labelId="board-select-label"
              label="Board"
              fullWidth
              value={newCircuit.board}
              onChange={(e) => updateCircuitField("board", e.target.value)}
            >
              {boards.map((board) => (
                <MenuItem key={board} value={board}>
                  {board}
                </MenuItem>
              ))}
            </Select>
            <TextField
              required
              label="Circuit Name"
              fullWidth
              value={newCircuit.name}
              onChange={(e) => updateCircuitField("name", e.target.value)}
            />
          </FormControl>
        </CardContent>
      </Card>
    </Modal>
  );
}

function ConfigSpeedDial({
  handleNewCircuitOpen,
}: {
  handleNewCircuitOpen: () => void;
}) {
  const actions = [
    {
      icon: <FaPlus size="50%" />,
      name: "New Circuit",
      onclick: handleNewCircuitOpen,
    },
    {
      icon: <FaRegFolder size="50%" />,
      name: "Open Circuit",
      onclick: () => console.log("Open Circuit"),
    },
  ];

  return (
    <SpeedDial
      ariaLabel="Speedial"
      sx={{position: "absolute", right: "20px", top: "20px"}}
      transitionDuration={0}
      icon={<SpeedDialIcon />}
      direction="left"
    >
      {actions.map((action) => (
        <SpeedDialAction
          key={action.name}
          icon={action.icon}
          tooltipTitle={action.name}
          onClick={action.onclick}
        />
      ))}
    </SpeedDial>
  );
}

export function Config() {
  const api = useHttpApi();
  const circuit = useAtomValue(circuitAtom);
  const boards = useAtomValue(boardsAtom);

  const [newCircuitModalOpen, setNewCircuitModalOpen] = useState(false);
  const [openCircuitModalOpen, setOpenCircuitModalOpen] = useState(false);

  const handleNewCircuitOpen = () => {
    setNewCircuitModalOpen(true);
    setOpenCircuitModalOpen(false);
  };

  const handleNewCircuitClose = () => {
    setNewCircuitModalOpen(false);
  };

  const handleCreateCircuit = async () => {
    // await api.createCircuit()
  };

  useEffect(() => {
    console.log(circuit);
  }, [circuit]);

  const Config = styled(Box)({
    width: "100%",
    height: "100%",
    display: "grid",
    gridTemplateColumns: "200px 1fr",
    gridTemplateRows: "1fr",
    gridTemplateAreas: `
      "list svg"
    `,
  });

  const SVG = styled("img")({
    gridArea: "svg",
    width: "600px",
    transform: "rotate(90deg)",
  });

  return (
    <>
      <AddCircuitModal
        open={newCircuitModalOpen}
        close={handleNewCircuitClose}
        boards={boards}
      />
      <Config sx={{position: "relative"}}>
        <List dense={true}>
          {circuit.pins.map((pin) => {
            return (
              <>
                <ListItem key={pin.pinNumber}>
                  <ListItemText primary={pin.pinNumber} />
                  <ListItemAvatar>
                    <Avatar
                      sx={{
                        cursor: "pointer",
                        width: "28px",
                        height: "28px",
                        bgcolor:
                          pin.state === pin.max
                            ? theme.palette.primary.main
                            : null,
                      }}
                      onClick={() => api.digitalWritePin(pin)}
                    >
                      {pin.type === PinType.Digital ? (
                        <PiWaveSineLight size="20px" />
                      ) : (
                        <PiWaveSquareLight size="20px" />
                      )}
                    </Avatar>
                  </ListItemAvatar>
                  <Stack direction="row" spacing={1} alignItems="center">
                    <p>I</p>
                    <Switch
                      size="small"
                      checked={pin.mode === PinMode.Output}
                      onChange={() => api.setupPin(pin)}
                    />
                    <p>O</p>
                  </Stack>
                </ListItem>
                <Divider />
              </>
            );
          })}
        </List>
        <ConfigSpeedDial handleNewCircuitOpen={handleNewCircuitOpen} />
        {/* <Button onClick={handleNewCircuitOpen}>
          <FaPlus />
        </Button> */}
      </Config>
    </>
  );
}
