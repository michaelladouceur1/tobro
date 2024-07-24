import { Button, TextField, AppBar } from "@mui/material";
import { useAtom } from "jotai";
import { useEffect, useState } from "react";
import { portsAtom } from "./atoms/portsAtom";
import { useHttpApi } from "./hooks/useHttpApi";

import "./App.css";
import { boardAtom } from "./atoms/boardAtom";

function App() {
  const api = useHttpApi();

  const [port, setPort] = useState<number | null>(null);
  const [analogValue, setAnalogValue] = useState(50);

  const [board, setBoard] = useAtom(boardAtom);
  const [ports, setPorts] = useAtom(portsAtom);

  useEffect(() => {
    const ws = new WebSocket("ws://localhost:8080/ws");

    ws.onopen = () => {
      console.log("Connected to server");
    };

    ws.onmessage = (event) => {
      const message = JSON.parse(event.data);

      if (message.type === "board") {
        console.log("board: ", message);
        const { pins } = message.data;
        setBoard({ pins });
      }

      if (message.type === "ports") {
        const { ports } = message.data;
        setPorts({ ports });
      }
    };
  }, []);

  return (
    <div>
      <AppBar position="static" color="primary">
        <h1>Tobro UI</h1>
      </AppBar>
      <div>
        <h2>Ports</h2>
        <ul>
          {ports.ports.map((port) => (
            <span>
              <li key={port}>{port}</li>
              <button
                onClick={() => {
                  console.log("Connecting to port: ", port);
                  api.connectPost({ connectRequest: { port } });
                }}
              >
                Connect
              </button>
            </span>
          ))}
        </ul>
      </div>
      <div>
        <h2>Digital Write Pin</h2>
        <TextField
          label="Port"
          type="number"
          size="small"
          value={port}
          onChange={(event) => {
            setPort(parseInt(event.target.value));
          }}
        />
        <TextField
          label="Value"
          type="number"
          size="small"
          value={analogValue}
          onChange={(event) => {
            setAnalogValue(parseInt(event.target.value));
          }}
        />
        <Button
          variant="outlined"
          onClick={() => {
            if (port === null) {
              return;
            }
            api.setupPinPost({ setupPinRequest: { pin: port, mode: "output" } });
          }}
        >
          Setup
        </Button>
        <Button
          variant="outlined"
          onClick={() => {
            if (port === null) {
              return;
            }
            api.digitalWritePinPost({
              digitalWritePinRequest: { pin: port, value: 1 },
            });
          }}
        >
          High
        </Button>
        <Button
          variant="outlined"
          onClick={() => {
            if (port === null) {
              return;
            }
            api.digitalWritePinPost({
              digitalWritePinRequest: { pin: port, value: 0 },
            });
          }}
        >
          Low
        </Button>
        <Button
          variant="outlined"
          onClick={async () => {
            if (port === null) {
              return;
            }
            const res = await api.analogWritePinPost({
              analogWritePinRequest: { pin: port, value: analogValue },
            });
            console.log(res);
          }}
        >
          Analog
        </Button>
      </div>
    </div>
  );
}

export default App;
