import { useEffect } from "react";

/*
 * useReloadNotification uses EventSource (SSE)
 * to get notified when updates are available.
 *
 * The callback will be triggered once for each incoming update,
 * this can then be used to invalidate SWR caches for example.
 */
export function useReloadNotification(cb: () => void) {
  useEffect(() => {
    const es = new EventSource("/events");
    es.onmessage = () => cb();
    return () => es.close();
  }, []);
}
