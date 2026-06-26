import { Box, Center, Divider, Loader, TextInput } from "@mantine/core";
import { Todo } from "./Todo";
import {
  CreateTodoSchema,
  type CreateTodoParams,
  type Todo as TodoModel,
} from "../api";
import { useTodoListManager } from "../hooks/useTodoListManager";
import { useForm, schemaResolver } from "@mantine/form";

interface TodoListSectionProps {
  items: TodoModel[];
  sectionName: (items: TodoModel[]) => string;
  render: (todo: TodoModel) => React.ReactNode;
}

function TodoListSection({ items, sectionName, render }: TodoListSectionProps) {
  if (items.length === 0) {
    return;
  }

  const label = sectionName(items);

  return (
    <Box>
      <Divider my="xs" label={label} labelPosition="center" />
      <div style={{ display: "flex", flexDirection: "column", gap: "5px" }}>
        {items?.map(render)}
      </div>
    </Box>
  );
}

export function TodoList() {
  const {
    isLoading,
    loadingError,
    openTodos,
    closedTodos,
    createTodo,
    isCreating,
    toggleTodo,
    deleteTodo,
  } = useTodoListManager();

  const form = useForm({
    mode: "uncontrolled",
    onSubmitPreventDefault: "always",
    initialValues: {
      title: "",
    },
    validate: schemaResolver(CreateTodoSchema, { sync: true }),
  });

  if (isLoading) {
    return (
      <Center h="100vh">
        <Loader color="gray" />
      </Center>
    );
  }

  if (loadingError) {
    console.error(loadingError);
  }

  const handleCreateTodo = async (params: CreateTodoParams) => {
    await createTodo(params);
    form.reset();
  };

  const pluralized = (n: number | null | undefined) =>
    (n ?? 0) === 1 ? "todo" : "todos";

  const handleOnToggle = async (id: string, newState: boolean) => {
    await toggleTodo({ id, newState });
  };

  const handleOnDelete = async (id: string) => {
    if (confirm("Are you sure you want to delete this todo?")) {
      await deleteTodo({ id });
    }
  };

  const makeSectionName = (
    openOrClosed: "open" | "closed",
    items: TodoModel[],
  ) => `${items.length} ${openOrClosed} ${pluralized(items.length)}`;

  return (
    <>
      <Box mt={20}>
        <form onSubmit={form.onSubmit(handleCreateTodo)}>
          <TextInput
            size="xl"
            name="title"
            placeholder="add a new todo"
            disabled={isCreating}
            loading={isCreating}
            key={form.key("title")}
            {...form.getInputProps("title")}
          />
        </form>
      </Box>

      <TodoListSection
        items={openTodos ?? []}
        sectionName={(items) => makeSectionName("open", items)}
        render={(todo) => (
          <Todo
            item={todo}
            key={`item-${todo.id}`}
            onToggle={handleOnToggle}
            onDelete={handleOnDelete}
          />
        )}
      />

      <TodoListSection
        items={closedTodos ?? []}
        sectionName={(items) => makeSectionName("closed", items)}
        render={(todo) => (
          <Todo
            item={todo}
            key={`item-${todo.id}`}
            onToggle={handleOnToggle}
            onDelete={handleOnDelete}
          />
        )}
      />
    </>
  );
}
