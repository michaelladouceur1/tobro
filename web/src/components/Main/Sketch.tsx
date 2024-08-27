import {
  Box,
  SpeedDial,
  SpeedDialAction,
  SpeedDialIcon,
  styled,
} from "@mui/material";
import {FaPlus, FaSave} from "react-icons/fa";
import {useHttpApi} from "../../hooks/useHttpApi";

function SketchSpeedDial() {
  const api = useHttpApi();

  const actions = [
    {
      icon: <FaPlus />,
      name: "New Sketch",
      onclick: () => console.log("New Sketch"),
    },
    {
      icon: <FaSave />,
      name: "Save Sketch",
      onclick: () => console.log("Save Sketch"),
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

export function Sketch() {
  const Sketch = styled(Box)({
    width: "100%",
    height: "100%",
    display: "flex",
  });

  return (
    <Sketch className="sketch-view" sx={{position: "relative"}}>
      <SketchSpeedDial />
      <img src="/images/sketch.png" alt="sketch" />
      <h1>SKETCH</h1>
    </Sketch>
  );
}
