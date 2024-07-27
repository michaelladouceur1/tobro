import { Box, Button, Card, Drawer, IconButton, List, ListItem, ListItemIcon, styled } from "@mui/material";
import { PiCircuitryLight, PiWaveformLight } from "react-icons/pi";

export function SideNav() {
  const SnContainer = styled(Card)({
    gridArea: "sidenav",
    display: "flex",
    justifyContent: "center",
    alignItems: "flex-start",
  });

  const SnList = styled(List)({});

  const SnListItem = styled(ListItem)({
    display: "flex",
    justifyContent: "center",
    width: "70px",
    height: "70px",
  });

  const SnButton = styled(IconButton)({
    width: "64px",
    height: "64px",
  });

  return (
    <SnContainer variant="outlined">
      <SnList>
        <SnListItem>
          <SnButton color="primary">
            <PiCircuitryLight size="30px" />
          </SnButton>
        </SnListItem>
        <SnListItem>
          <SnButton color="primary">
            <PiWaveformLight size="30px" />
          </SnButton>
        </SnListItem>
      </SnList>
    </SnContainer>
  );
}
