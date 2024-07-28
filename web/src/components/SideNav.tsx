import {Card, IconButton, List, ListItem, styled} from "@mui/material";
import {useAtom} from "jotai";
import {PiCircuitryLight, PiWaveformLight} from "react-icons/pi";
import {mainViewAtom} from "../atoms/mainViewAtom";

export function SideNav() {
  const [mainView, setMainView] = useAtom(mainViewAtom);

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
          <SnButton color="primary" onClick={() => setMainView("config")}>
            <PiCircuitryLight size="30px" />
          </SnButton>
        </SnListItem>
        <SnListItem>
          <SnButton color="primary" onClick={() => setMainView("pin-list")}>
            <PiWaveformLight size="30px" />
          </SnButton>
        </SnListItem>
      </SnList>
    </SnContainer>
  );
}
