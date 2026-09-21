import { useCallback, useState } from "react";

// usePagination：分页状态管理，列表页复用。
export function usePagination(initialPage = 1, initialPageSize = 10) {
  const [page, setPage] = useState(initialPage);
  const [pageSize, setPageSize] = useState(initialPageSize);
  const [total, setTotal] = useState(0);

  const reset = useCallback(() => setPage(1), []);
  const next = useCallback(() => setPage((p) => p + 1), []);
  const prev = useCallback(() => setPage((p) => Math.max(1, p - 1)), []);

  return { page, pageSize, total, setTotal, reset, next, prev, setPage, setPageSize };
}
