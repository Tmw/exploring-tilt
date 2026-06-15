import useSWR from "swr";
import { fetchTodos } from "../api";

export function useGetTodos() {
  return useSWR("/todos", fetchTodos);
}
