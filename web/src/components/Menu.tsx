import {
  Box,
  Card,
  FormControl,
  InputLabel,
  MenuItem,
  Select,
  Toolbar,
} from "@mui/material";
import {useAtomValue} from "jotai";
import {portConnectionAtom} from "../atoms/portConnection";
import {portsAtom} from "../atoms/portsAtom";
import {useHttpApi} from "../hooks/useHttpApi";

export function Menu() {
  const api = useHttpApi();
  const ports = useAtomValue(portsAtom);
  const portConnection = useAtomValue(portConnectionAtom);

  return (
    <Card variant="outlined" sx={{gridArea: "menu"}}>
      <Toolbar>
        <Box>
          <FormControl sx={{width: "200px"}}>
            <InputLabel id="port-select-label">Port</InputLabel>
            <Select
              labelId="port-select-label"
              label="Port"
              value={portConnection.portName}
              onChange={(e) => api.connect(e.target.value)}
            >
              {/* <MenuItem value={portConnection.portName}>Select a port</MenuItem> */}
              {ports.ports.map((port) => (
                <MenuItem key={port} value={port}>
                  {port}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </Box>
      </Toolbar>
    </Card>
  );
}
