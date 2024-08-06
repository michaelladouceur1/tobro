import {
  Avatar,
  Box,
  Divider,
  List,
  ListItem,
  ListItemAvatar,
  ListItemText,
  Modal,
  SpeedDial,
  SpeedDialAction,
  SpeedDialIcon,
  Stack,
  styled,
  Switch,
} from "@mui/material";
import {useAtomValue} from "jotai";
import {FaPlus, FaRegFolder} from "react-icons/fa";
import {PiWaveSineLight, PiWaveSquareLight} from "react-icons/pi";
import {circuitAtom} from "../../atoms/circuitAtom";
import {useHttpApi} from "../../hooks/useHttpApi";
import {theme} from "../../theme";
import {DigitalState, Pin, PinMode, PinType} from "../../types";
import {useState} from "react";

export function Config() {
  const api = useHttpApi();
  const circuit = useAtomValue(circuitAtom);

  const [newCircuitModalOpen, setNewCircuitModalOpen] = useState(false);
  const [openCircuitModalOpen, setOpenCircuitModalOpen] = useState(false);

  const handleNewCircuitOpen = () => {
    setNewCircuitModalOpen(true);
    setOpenCircuitModalOpen(false);
  };
  const handleNewCircuitClose = () => {
    setNewCircuitModalOpen(false);
  };
  //   const canvasRef = useRef<HTMLCanvasElement>(null);

  //   useEffect(() => {
  //     const canvas = canvasRef.current;
  //     if (!canvas) {
  //       return;
  //     }
  //     const ctx = canvas.getContext("2d");
  //     if (!ctx) {
  //       return;
  //     }

  //     const img = new Image();
  //     img.src = ArduinoNanoSVG;
  //     img.onload = () => {
  //       ctx.clearRect(0, 0, canvas.width, canvas.height);
  //       ctx.save();
  //       ctx.translate(canvas.width / 2, canvas.height / 2);
  //       ctx.rotate((90 * Math.PI) / 180); // Rotate 90 degrees
  //       ctx.drawImage(img, -img.width / 2, -img.height / 2);
  //       ctx.restore();

  //       // Define pin positions (example positions, adjust as needed)
  //       const pinPositions = [
  //         {x: -img.width / 2 + 10, y: -img.height / 2 + 20}, // Pin 1
  //         {x: -img.width / 2 + 10, y: -img.height / 2 + 40}, // Pin 2
  //         // Add more pins as needed
  //       ];

  //       // Draw lines from pins to the edge of the canvas
  //       pinPositions.forEach((pin) => {
  //         ctx.beginPath();
  //         ctx.moveTo(canvas.width / 2 + pin.x, canvas.height / 2 + pin.y);
  //         ctx.lineTo(canvas.width, canvas.height / 2 + pin.y); // Line to the right edge
  //         ctx.strokeStyle = "red"; // Line color
  //         ctx.lineWidth = 2; // Line width
  //         ctx.stroke();
  //       });
  //     };
  //   }, []);

  // const handleSetupPin = async (pin: Pin) => {
  //   const {id} = pin;
  //   const mode = pin.mode === PinMode.Output ? PinMode.Input : PinMode.Output;
  //   const res = await api.setupPinPost({setupPinRequest: {pin: id, mode}});
  //   console.log(res);
  // };

  // const handleDigitalWrite = async (pin: Pin) => {
  //   const {id} = pin;
  //   const value = pin.state === pin.max ? DigitalState.Low : DigitalState.High;
  //   await api.digitalWritePinPost({
  //     digitalWritePinRequest: {pin: id, value},
  //   });
  // };

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
    <Config sx={{position: "relative"}}>
      <Modal open={newCircuitModalOpen} onClose={handleNewCircuitClose}>
        <Box
          sx={{
            position: "absolute",
            top: "50%",
            left: "50%",
            transform: "translate(-50%, -50%)",
            width: "80%",
            bgcolor: "background.paper",
            border: "2px solid #000",
            boxShadow: 24,
            p: 4,
          }}
        >
          <h2>New Circuit</h2>
          <p>Content</p>
        </Box>
      </Modal>
      <List dense={true}>
        {circuit.pins.map((pin) => {
          return (
            <>
              <ListItem key={pin.id}>
                <ListItemText primary={pin.id} />
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
      {/* <SVG src={ArduinoNanoSVG} alt="Arduino Nano" /> */}
      <SpeedDial
        ariaLabel="Speedial"
        sx={{position: "absolute", right: "20px", top: "20px"}}
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
    </Config>
  );
}
