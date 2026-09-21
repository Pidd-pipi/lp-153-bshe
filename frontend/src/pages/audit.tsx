import { useCallback, useEffect, useState } from "react";
import type { AuditLog } from "@/api/audit";
import { auditApi } from "@/api/audit";
import EmptyState from "@/components/EmptyState";
import RequireAuth from "@/components/RequireAuth";
import { formatDate } from "@/utils/format";

export default function Audit() {
  const [logs, setLogs] = useState<AuditLog[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(true);

  const load = useCallback(async () => {
    setLoading(true);
    try {
      const data = await auditApi.list({ page, page_size: 20 });
      setLogs(data.items);
      setTotal(data.total);
    } finally {
      setLoading(false);
    }
  }, [page]);

  useEffect(() => {
    load();
  }, [load]);

  const totalPages = Math.max(1, Math.ceil(total / 20));

  return (
    <RequireAuth adminOnly>
      <div className="space-y-5">
        <div>
          <h1 className="text-2xl font-bold text-purple-700">📋 操作审计日志</h1>
          <p className="mt-1 text-sm text-gray-500">记录系统关键写操作（仅管理员可见）</p>
        </div>

        {loading ? (
          <p className="py-16 text-center text-gray-400">加载中...</p>
        ) : logs.length === 0 ? (
          <EmptyState title="暂无审计记录" icon="📋" />
        ) : (
          <div className="card overflow-x-auto">
            <table className="w-full text-left text-sm">
              <thead>
                <tr className="border-b border-purple-100 text-xs text-gray-400">
                  <th className="py-2 pr-4">时间</th>
                  <th className="py-2 pr-4">用户</th>
                  <th className="py-2 pr-4">操作</th>
                  <th className="py-2 pr-4">对象</th>
                  <th className="py-2 pr-4">详情</th>
                  <th className="py-2">IP</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-purple-50">
                {logs.map((log) => (
                  <tr key={log.id}>
                    <td className="py-2.5 pr-4 text-xs text-gray-500">{formatDate(log.created_at)}</td>
                    <td className="py-2.5 pr-4">{log.username || `#${log.user_id}`}</td>
                    <td className="py-2.5 pr-4"><span className="rounded-full bg-purple-50 px-2 py-0.5 text-xs text-purple-600">{log.action}</span></td>
                    <td className="py-2.5 pr-4 text-xs text-gray-500">{log.entity_type}/{log.entity_id}</td>
                    <td className="py-2.5 pr-4 text-gray-600">{log.detail}</td>
                    <td className="py-2.5 text-xs text-gray-400">{log.ip}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}

        <div className="flex items-center justify-center gap-4">
          <button className="btn-secondary" disabled={page <= 1} onClick={() => setPage((p) => p - 1)}>上一页</button>
          <span className="text-sm text-gray-500">{page} / {totalPages}</span>
          <button className="btn-secondary" disabled={page >= totalPages} onClick={() => setPage((p) => p + 1)}>下一页</button>
        </div>
      </div>
    </RequireAuth>
  );
}
