import useSWRMutation from "swr/mutation";
import { deleteTodo, type DeleteTodoParams } from "../api";

export function useDeleteTodo() {
  return useSWRMutation(
    "/todos",
    (_url: string, { arg }: { arg: DeleteTodoParams }) => {
      return deleteTodo(arg);
    },
  );
}
