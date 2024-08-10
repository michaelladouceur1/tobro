import { createTheme} from "@mui/material/styles";


export const theme = createTheme({
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
      MuiInputLabel: {
        styleOverrides: {
          root: {
            color: "#85827E",
          },
        },
      },
      MuiCard: {
        styleOverrides: {
          root: {
            backgroundColor: "#1B2F33",
            color: "#F5E0B7",
          },
        },
      },
    },
  });