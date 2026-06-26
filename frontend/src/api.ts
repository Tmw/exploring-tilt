import * as z from "zod";

const API_PREFIX = "/api";

export const TodoSchema = z.object({
  id: z.string(),
  title: z.string(),
  createdAt: z.coerce.date(),
  completedAt: z.coerce.date().nullable(),
});

export const TodosSchema = z.array(TodoSchema);
export type Todo = z.infer<typeof TodoSchema>;
export type Todos = z.infer<typeof TodosSchema>;

export async function fetchTodos() {
  const resp = await fetch(API_PREFIX + "/todos");
  if (!resp.ok) {
    throw new Error("error fetching todos");
  }
  const json = await resp.json();
  return TodosSchema.parse(json);
}

export const CreateTodoSchema = z.object({
  title: z.string().min(3),
});

export type CreateTodoParams = z.infer<typeof CreateTodoSchema>;
export async function createTodo(params: CreateTodoParams) {
  const payload = CreateTodoSchema.parse(params);
  const resp = await fetch(API_PREFIX + "/todos", {
    method: "POST",
    headers: {
      "content-type": "application/json",
    },
    body: JSON.stringify(payload),
  });

  if (!resp.ok) {
    throw new Error("error creating todo");
  }

  const json = await resp.json();
  return TodosSchema.parse(json);
}

export type DeleteTodoParams = Pick<Todo, "id">;
export async function deleteTodo(params: DeleteTodoParams) {
  const resp = await fetch(API_PREFIX + "/todos/" + params.id, {
    method: "DELETE",
  });

  if (!resp.ok) {
    throw new Error("error deleting todo");
  }

  const json = await resp.json();
  return TodosSchema.parse(json);
}

export type ToggleTodoParams = {
  id: string;
  newState: boolean;
};

export async function toggleTodo(params: ToggleTodoParams) {
  const resp = await fetch(API_PREFIX + "/todos/" + params.id + "/status", {
    method: "PATCH",
    headers: {
      "content-type": "application/json",
    },
    body: JSON.stringify({ newState: params.newState }),
  });

  if (!resp.ok) {
    throw new Error("error deleting todo");
  }

  const json = await resp.json();
  return TodosSchema.parse(json);
}
