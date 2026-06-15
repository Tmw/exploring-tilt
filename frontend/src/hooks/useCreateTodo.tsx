import useSWRMutation from "swr/mutation";
import { createTodo, type CreateTodoParams } from "../api";

export function useCreateTodo() {
  return useSWRMutation(
    "/todos",
    (_url: string, { arg }: { arg: CreateTodoParams }) => {
      return createTodo(arg);
    },
  );
}
