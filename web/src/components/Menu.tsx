import { Box, Card, FormControl, InputLabel, MenuItem, Select, SelectChangeEvent, Toolbar } from "@mui/material";
import { useAtomValue } from "jotai";
import { useState } from "react";
import { portsAtom } from "../atoms/portsAtom";
import { useHttpApi } from "../hooks/useHttpApi";

export function Menu() {
  const api = useHttpApi();
  const ports = useAtomValue(portsAtom);
  const [port, setPort] = useState<string | undefined>(undefined);
  const [error, setError] = useState<boolean>(false);

  const handleConnect = async (e: SelectChangeEvent<string>) => {
    e.preventDefault();
    if (!e.target.value) {
      return;
    }
    setError(false);
    const res = await api.connectPost({ connectRequest: { port: e.target.value } });
    console.log("connect response: ", res);
  };

  return (
    <Card variant="outlined" sx={{ gridArea: "menu" }}>
      <Toolbar>
        {/* <h1>Tobro UI</h1> */}
        <Box>
          <FormControl sx={{ width: "200px" }}>
            <InputLabel id="port-select-label">Port</InputLabel>
            <Select labelId="port-select-label" label="Port" value={port} onChange={handleConnect}>
              <MenuItem value={undefined}>Select a port</MenuItem>
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
