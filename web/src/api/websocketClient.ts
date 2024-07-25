const BASE_URL = "ws://localhost:8080";
let ws: WebSocket;

export function connect() {
  ws = new WebSocket(`${BASE_URL}`);

  ws.onopen = () => {
    console.log("Connected to server");
  };

  ws.onmessage = (event) => {
    console.log("Received message from server", event.data);
  };
}
