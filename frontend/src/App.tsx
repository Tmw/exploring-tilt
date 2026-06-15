import "@mantine/core/styles.css";
import { MantineProvider } from "@mantine/core";
import { Page } from "./components/Page";
import { TodoList } from "./components/TodoList";

function App() {
  return (
    <MantineProvider>
      <Page>
        <TodoList />
      </Page>
    </MantineProvider>
  );
}

export default App;
