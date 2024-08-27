import {Box, styled} from "@mui/material";

export function Sketch() {
  const Sketch = styled(Box)({
    width: "100%",
    height: "100%",
    display: "flex",
  });

  return (
    <Sketch className="sketch-view">
      <img src="/images/sketch.png" alt="sketch" />
    </Sketch>
  );
}
