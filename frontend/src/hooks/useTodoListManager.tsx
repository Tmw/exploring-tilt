import { useGetTodos } from "../hooks/useGetTodos";
import { useCreateTodo } from "../hooks/useCreateTodo";
import { useDeleteTodo } from "../hooks/useDeleteTodo";
import { useToggleTodo } from "../hooks/useToggleTodo";
import type { Todo } from "../api";

const isCompleted = (todo: Todo) => todo.completedAt !== null;
const isNotCompleted = (todo: Todo) => todo.completedAt === null;

export function useTodoListManager() {
  const { data: allTodos, error: loadingError, isLoading } = useGetTodos();
  const { trigger: createTodo, isMutating: isCreating } = useCreateTodo();
  const { trigger: deleteTodo, isMutating: isDeleting } = useDeleteTodo();
  const { trigger: toggleTodo, isMutating: isToggling } = useToggleTodo();

  const openTodos = allTodos?.filter(isNotCompleted);
  const closedTodos = allTodos?.filter(isCompleted);

  return {
    // fetching
    loadingError,
    isLoading,
    allTodos,

    openTodos,
    closedTodos,

    // creation
    createTodo,
    isCreating,

    // patching
    toggleTodo,
    isToggling,

    // deleting
    deleteTodo,
    isDeleting,
  };
}
