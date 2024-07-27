import React from "react";
import ReactDOM from "react-dom/client";
import { ThemeProvider, createTheme } from "@mui/material/styles";
import App from "./App";

import "./index.css";

const theme = createTheme({
  palette: {
    primary: {
      main: "#549F93",
    },
    secondary: {
      main: "#28536B",
    },
    text: {
      primary: "#F5E0B7",
    },
    error: {
      main: "#FB6107",
    },
    background: {
      default: "#0C1618",
    },
  },
  components: {
    MuiCard: {},
  },
});

const root = ReactDOM.createRoot(document.getElementById("root") as HTMLElement);

root.render(
  <React.StrictMode>
    <ThemeProvider theme={theme}>
      <App />
    </ThemeProvider>
  </React.StrictMode>
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
