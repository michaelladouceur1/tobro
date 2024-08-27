import {Card} from "@mui/material";
import {useAtomValue} from "jotai";
import {mainViewAtom} from "../../atoms/mainViewAtom";
import {Config} from "./Config";
import {Sketch} from "./Sketch";
import {useEffect, useMemo} from "react";

export function Main() {
  const mainView = useAtomValue(mainViewAtom);

  const viewMap = useMemo(
    () => ({
      config: <Config />,
      sketch: <Sketch />,
    }),
    []
  );

  useEffect(() => {
    console.log("Main view changed to", mainView);
  }, [mainView]);

  return (
    <Card variant="outlined" sx={{gridArea: "main", overflowY: "auto"}}>
      {viewMap[mainView]}
    </Card>
  );
}
