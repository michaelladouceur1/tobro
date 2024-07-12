import React, {useEffect, useState} from "react";
import logo from "./logo.svg";
import {useAtom} from "jotai";
import {portsAtom} from "./atoms/portsAtom";
import {connectPort} from "./api/connectService";
import "./App.css";

function App() {
  const [delay, setDelay] = useState(100);
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
      <h1>Tobro UI</h1>
      <input
        type="number"
        value={delay}
        placeholder="Enter delay in milliseconds"
        onChange={(e) => setDelay(parseInt(e.target.value))}
      />
      <button
        onClick={() => {
          console.log("Sending request");
          fetch("/delay", {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify({
              delay: delay,
            }),
          });
        }}
      >
        Send
      </button>

      <div>
        <h2>Ports</h2>
        <ul>
          {ports.ports.map((port) => (
            <span>
              <li key={port}>{port}</li>
              <button
                onClick={() => {
                  connectPort({port});
                }}
              >
                Connect
              </button>
            </span>
          ))}
        </ul>
      </div>
    </div>
  );
}

export default App;
