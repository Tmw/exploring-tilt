import useSWRMutation from "swr/mutation";
import { toggleTodo, type ToggleTodoParams } from "../api";

export function useToggleTodo() {
  return useSWRMutation(
    "/todos",
    (_url: string, { arg }: { arg: ToggleTodoParams }) => {
      return toggleTodo(arg);
    },
  );
}
