import {Box, styled} from "@mui/material";
import {Main} from "./components/Main/Main";
import {Menu} from "./components/Menu";
import {SideNav} from "./components/SideNav";
import {Terminal} from "./components/Terminal";
import {useWsApi} from "./hooks/useWsApi";

import "./App.css";

function App() {
  useWsApi();

  const App = styled(Box)({
    boxSizing: "border-box",
    display: "grid",
    gridTemplateAreas: `
      "menu menu menu"
      "sidenav main terminal"
    `,
    gridTemplateRows: "75px 1fr",
    gridTemplateColumns: "75px 1fr 500px",
    gap: "10px",
    overflowY: "hidden",
    maxHeight: "100vh",
    height: "100vh",
    padding: "5px",
  });

  return (
    <App>
      <Menu />
      <SideNav />
      <Main />
      <Terminal />
    </App>
  );
}

export default App;
