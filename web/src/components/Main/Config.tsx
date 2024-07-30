import {
  Box,
  Card,
  List,
  ListItem,
  SpeedDial,
  SpeedDialAction,
  styled,
  Switch,
} from "@mui/material";
import ArduinoNanoSVG from "../../assets/arduino-nano.svg";
import {useEffect, useRef} from "react";
import {useAtomValue} from "jotai";
import {boardAtom} from "../../atoms/boardAtom";

export function Config() {
  const board = useAtomValue(boardAtom);
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

  const leftPins = [2, 3, 4, 5, 6, 7, 8, 9, 10, 11];

  const rightPins = [12, 13, 14, 15, 16, 17, 18, 19, 20, 21];

  const Config = styled(Box)({
    width: "100%",
    height: "100%",
    display: "flex",
    justifyContent: "center",
    alignItems: "center",
  });

  const SVG = styled("img")({
    gridArea: "svg",
    width: "600px",
    transform: "rotate(90deg)",
  });

  const actions = [
    {icon: <Switch />, name: "Mode"},
    {icon: <Switch />, name: "State"},
  ];

  return (
    <Config>
      <List sx={{gap: "10px"}}>
        {leftPins.map((pin) => {
          return (
            // <ListItem key={pin}>
            <Card variant="outlined" sx={{width: "200px"}}>
              {pin}
            </Card>
            // </ListItem>
          );
        })}
      </List>
      <SVG src={ArduinoNanoSVG} alt="Arduino Nano" />
    </Config>
  );
}
