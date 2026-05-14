import { Box, render, Text, useInput, useWindowSize } from "ink";

function App() {
  const { columns, rows } = useWindowSize();
  useInput(() => {});
  return (
    <Box width={columns} height={rows}>
      <Box width={30} borderColor={"green"} borderStyle={"single"}>
        <Text>Hello world</Text>
      </Box>
    </Box>
  );
}
render(<App />, { alternateScreen: true });
