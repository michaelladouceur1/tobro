import {Box, Card} from "@mui/material";
import {useAtomValue} from "jotai";
import {mainViewAtom} from "../../atoms/mainViewAtom";
import {Config} from "./Config";
import {PinList} from "./PinList";

export function Main() {
  const mainView = useAtomValue(mainViewAtom);

  const viewMap = {
    config: <Config />,
    "pin-list": <PinList />,
  };

  return (
    <Card variant="outlined" sx={{gridArea: "main", overflowY: "auto"}}>
      {viewMap[mainView]}
    </Card>
  );
}
