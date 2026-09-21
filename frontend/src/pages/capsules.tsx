import { useCallback, useEffect, useState } from "react";
import type { Capsule } from "@/api/capsule";
import { capsuleApi } from "@/api/capsule";
import { uploadApi } from "@/api/upload";
import ConfirmDialog from "@/components/ConfirmDialog";
import EmptyState from "@/components/EmptyState";
import RequireAuth from "@/components/RequireAuth";
import StatusBadge from "@/components/StatusBadge";
import { useToast } from "@/components/Toast";
import { useAuth } from "@/hooks/useAuth";
import { formatDate } from "@/utils/format";

export default function Capsules() {
  const { isAuthed } = useAuth();
  const toast = useToast();
  const [capsules, setCapsules] = useState<Capsule[]>([]);
  const [loading, setLoading] = useState(true);
  const [showCreate, setShowCreate] = useState(false);
  const [deleteId, setDeleteId] = useState<number | null>(null);
  const [form, setForm] = useState({ title: "", content: "", unlock_at: "" });
  const [audioUrl, setAudioUrl] = useState("");
  const [submitting, setSubmitting] = useState(false);

  const load = useCallback(async () => {
    setLoading(true);
    try {
      const data = await capsuleApi.mine({ page: 1, page_size: 50 });
      setCapsules(data.items);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    if (isAuthed()) load();
  }, [isAuthed, load]);

  const uploadAudio = async (file: File) => {
    try {
      const result = await uploadApi.uploadAudio(file);
      setAudioUrl(result.url);
      toast.show("音频上传成功");
    } catch (e) {
      toast.show((e as Error).message, "error");
    }
  };

  const create = async () => {
    if (!form.title || !form.content || !form.unlock_at) {
      toast.show("请填写标题、内容和解锁时间", "error");
      return;
    }
    setSubmitting(true);
    try {
      await capsuleApi.create({
        title: form.title,
        content: form.content,
        audio_url: audioUrl || undefined,
        unlock_at: `${form.unlock_at}:00+08:00`,
      });
      toast.show("时光胶囊已封存 📦");
      setShowCreate(false);
      setForm({ title: "", content: "", unlock_at: "" });
      setAudioUrl("");
      load();
    } catch (e) {
      toast.show((e as Error).message, "error");
    } finally {
      setSubmitting(false);
    }
  };

  const remove = async () => {
    if (!deleteId) return;
    try {
      await capsuleApi.remove(deleteId);
      toast.show("胶囊已删除");
      setDeleteId(null);
      load();
    } catch (e) {
      toast.show((e as Error).message, "error");
    }
  };

  return (
    <RequireAuth>
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold text-purple-700">📦 时光胶囊</h1>
            <p className="mt-1 text-sm text-gray-500">写给未来的自己，到时间才能打开</p>
          </div>
          <button className="btn-primary" onClick={() => setShowCreate(true)}>封存新胶囊</button>
        </div>

        {loading ? (
          <p className="py-16 text-center text-gray-400">加载中...</p>
        ) : capsules.length === 0 ? (
          <EmptyState title="还没有时光胶囊" description="把此刻的心情封存给未来" icon="📦" />
        ) : (
          <div className="grid gap-5 sm:grid-cols-2 lg:grid-cols-3">
            {capsules.map((c) => (
              <div key={c.id} className={`card space-y-3 ${c.status === "unlocked" ? "border-emerald-200" : "border-amber-100"}`}>
                <div className="flex items-center justify-between">
                  <span className="text-2xl">{c.status === "unlocked" ? "📬" : "🔒"}</span>
                  <StatusBadge status={c.status} kind="capsule" />
                </div>
                <div>
                  <h3 className="font-semibold text-gray-800">{c.title}</h3>
                  <p className="mt-1 text-xs text-gray-400">解锁时间 {formatDate(c.unlock_at)}</p>
                </div>
                {c.status === "unlocked" ? (
                  <div className="animate-unlock">
                    <p className="text-sm text-gray-600">{c.content}</p>
                    {c.audio_url && (
                      <audio controls className="mt-2 w-full" src={c.audio_url}>
                        您的浏览器不支持音频播放
                      </audio>
                    )}
                  </div>
                ) : (
                  <p className="rounded-lg bg-amber-50 px-3 py-2 text-xs text-amber-600">🔐 未到解锁时间，内容暂不可见</p>
                )}
                <div className="flex justify-end">
                  <button className="text-xs text-red-400 hover:text-red-600" onClick={() => setDeleteId(c.id)}>删除</button>
                </div>
              </div>
            ))}
          </div>
        )}

        {showCreate && (
          <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 px-4" onClick={() => setShowCreate(false)}>
            <div className="w-full max-w-lg space-y-4 rounded-2xl bg-white p-6 shadow-xl" onClick={(e) => e.stopPropagation()}>
              <h3 className="text-lg font-semibold text-gray-800">封存时光胶囊</h3>
              <div>
                <label className="label">标题</label>
                <input className="input" value={form.title} onChange={(e) => setForm((f) => ({ ...f, title: e.target.value }))} placeholder="给未来的信" />
              </div>
              <div>
                <label className="label">内容</label>
                <textarea className="input min-h-[100px]" value={form.content} onChange={(e) => setForm((f) => ({ ...f, content: e.target.value }))} placeholder="想对未来的自己说什么..." />
              </div>
              <div>
                <label className="label">解锁时间</label>
                <input className="input" type="datetime-local" value={form.unlock_at} onChange={(e) => setForm((f) => ({ ...f, unlock_at: e.target.value }))} />
              </div>
              <div>
                <label className="label">语音留言（选填）</label>
                <input type="file" accept="audio/*" className="text-sm" onChange={(e) => e.target.files?.[0] && uploadAudio(e.target.files[0])} />
                {audioUrl && <p className="mt-1 text-xs text-emerald-600">✓ 音频已上传</p>}
              </div>
              <div className="flex justify-end gap-3">
                <button className="btn-secondary" onClick={() => setShowCreate(false)}>取消</button>
                <button className="btn-primary" disabled={submitting} onClick={create}>{submitting ? "封存中..." : "封存"}</button>
              </div>
            </div>
          </div>
        )}

        <ConfirmDialog
          open={deleteId !== null}
          title="删除时光胶囊"
          description="删除后无法恢复，确定要删除吗？"
          confirmText="删除"
          danger
          onConfirm={remove}
          onCancel={() => setDeleteId(null)}
        />
      </div>
    </RequireAuth>
  );
}
