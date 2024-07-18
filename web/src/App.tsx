import {Button, TextField, AppBar} from "@mui/material";
import {useAtom} from "jotai";
import {useEffect, useState} from "react";
import {portsAtom} from "./atoms/portsAtom";
import {useHttpApi} from "./hooks/useHttpApi";

import "./App.css";

function App() {
  const api = useHttpApi();

  const [delay, setDelay] = useState(100);
  const [port, setPort] = useState<number | null>(null);
  const [pwmDutyCycle, setPwmDutyCycle] = useState(0.5);
  const [pwmPeriod, setPwmPeriod] = useState(10);

  const [ports, setPorts] = useAtom(portsAtom);

  useEffect(() => {
    const ws = new WebSocket("ws://localhost:8080/ws");

    ws.onopen = () => {
      console.log("Connected to server");
    };

    ws.onmessage = (event) => {
      const message = JSON.parse(event.data);

      if (message.type === "ports") {
        const {ports} = message.data;
        setPorts({ports});
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
                  api.connectPost({connectRequest: {port}});
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
          label="Value"
          type="number"
          size="small"
          value={port}
          onChange={(event) => {
            setPort(parseInt(event.target.value));
          }}
        />
        <Button
          variant="outlined"
          onClick={() => {
            if (port === null) {
              return;
            }
            api.setupPinPost({setupPinRequest: {pin: port, mode: "output"}});
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
            api.digitalWritePinPost({writePinRequest: {pin: port, value: 1}});
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
            api.digitalWritePinPost({writePinRequest: {pin: port, value: 0}});
          }}
        >
          Low
        </Button>
        <Button
          variant="outlined"
          onClick={() => {
            if (port === null) {
              return;
            }
            api.pwmPost({
              pWMRequest: {
                pin: port,
                dutyCycle: pwmDutyCycle,
                period: pwmPeriod,
              },
            });
          }}
        >
          PWM
        </Button>
      </div>
      <div>
        <TextField
          label="Duty Cycle"
          type="number"
          size="small"
          value={pwmDutyCycle}
          onChange={(event) => {
            setPwmDutyCycle(parseInt(event.target.value));
          }}
        />
        <TextField
          label="Period"
          type="number"
          size="small"
          value={pwmPeriod}
          onChange={(event) => {
            setPwmPeriod(parseInt(event.target.value));
          }}
        />
      </div>
    </div>
  );
}

export default App;
