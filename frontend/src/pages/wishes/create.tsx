import { useState } from "react";
import { useRouter } from "next/router";
import type { CreateWishPayload } from "@/api/wish";
import { wishApi } from "@/api/wish";
import { uploadApi } from "@/api/upload";
import RequireAuth from "@/components/RequireAuth";
import { useToast } from "@/components/Toast";
import { CATEGORY_TEXT, DIFFICULTY_TEXT, VISIBILITY_TEXT } from "@/constants";

export default function CreateWish() {
  const router = useRouter();
  const toast = useToast();
  const [form, setForm] = useState<CreateWishPayload>({
    title: "",
    content: "",
    category: "study",
    visibility: "public",
    difficulty: "medium",
    image_urls: [],
    is_anonymous: false,
  });
  const [deadline, setDeadline] = useState("");
  const [uploading, setUploading] = useState(false);
  const [submitting, setSubmitting] = useState(false);

  const set = <K extends keyof CreateWishPayload>(k: K, v: CreateWishPayload[K]) =>
    setForm((f) => ({ ...f, [k]: v }));

  const upload = async (file: File) => {
    setUploading(true);
    try {
      const result = await uploadApi.uploadImage(file);
      set("image_urls", [...(form.image_urls || []), result.url]);
      toast.show("图片上传成功");
    } catch (e) {
      toast.show((e as Error).message, "error");
    } finally {
      setUploading(false);
    }
  };

  const submit = async () => {
    if (form.title.trim().length < 2 || form.content.trim().length < 5) {
      toast.show("标题至少 2 字，内容至少 5 字", "error");
      return;
    }
    setSubmitting(true);
    try {
      const payload = { ...form, expected_deadline: deadline ? `${deadline}T23:59:59+08:00` : null };
      await wishApi.create(payload);
      toast.show("心愿发布成功 🌟");
      router.push("/");
    } catch (e) {
      toast.show((e as Error).message, "error");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <RequireAuth>
      <div className="mx-auto max-w-2xl">
        <h1 className="text-2xl font-bold text-purple-700">✨ 发布心愿</h1>
        <p className="mt-1 text-sm text-gray-500">写下愿望，等待有缘的圆梦人</p>

        <div className="card mt-5 space-y-4">
          <div>
            <label className="label">心愿标题</label>
            <input className="input" value={form.title} onChange={(e) => set("title", e.target.value)} placeholder="例如：去冰岛看一次极光" />
          </div>
          <div>
            <label className="label">心愿内容</label>
            <textarea className="input min-h-[120px]" value={form.content} onChange={(e) => set("content", e.target.value)} placeholder="详细描述你的心愿、背后的故事..." />
          </div>

          <div className="grid gap-4 sm:grid-cols-3">
            <div>
              <label className="label">分类</label>
              <select className="input" value={form.category} onChange={(e) => set("category", e.target.value)}>
                {Object.entries(CATEGORY_TEXT).map(([k, v]) => <option key={k} value={k}>{v}</option>)}
              </select>
            </div>
            <div>
              <label className="label">可见范围</label>
              <select className="input" value={form.visibility} onChange={(e) => set("visibility", e.target.value)}>
                {Object.entries(VISIBILITY_TEXT).map(([k, v]) => <option key={k} value={k}>{v}</option>)}
              </select>
            </div>
            <div>
              <label className="label">难度</label>
              <select className="input" value={form.difficulty} onChange={(e) => set("difficulty", e.target.value)}>
                {Object.entries(DIFFICULTY_TEXT).map(([k, v]) => <option key={k} value={k}>{v}</option>)}
              </select>
            </div>
          </div>

          <div className="grid gap-4 sm:grid-cols-2">
            <div>
              <label className="label">期望完成时间</label>
              <input className="input" type="date" value={deadline} onChange={(e) => setDeadline(e.target.value)} />
            </div>
            <div className="flex items-end">
              <label className="flex items-center gap-2 pb-2 text-sm text-gray-600">
                <input type="checkbox" checked={form.is_anonymous} onChange={(e) => set("is_anonymous", e.target.checked)} />
                匿名发布（隐藏作者）
              </label>
            </div>
          </div>

          <div>
            <label className="label">配图（选填，可传多张）</label>
            <div className="flex flex-wrap items-center gap-3">
              {(form.image_urls || []).map((url) => (
                <img key={url} src={url} className="h-20 w-20 rounded-xl object-cover" alt="wish" />
              ))}
              <label className="flex h-20 w-20 cursor-pointer items-center justify-center rounded-xl border-2 border-dashed border-purple-200 text-2xl text-purple-300 hover:bg-purple-50">
                {uploading ? "..." : "+"}
                <input type="file" accept="image/*" className="hidden" onChange={(e) => e.target.files?.[0] && upload(e.target.files[0])} />
              </label>
            </div>
          </div>

          <button className="btn-primary w-full" disabled={submitting || uploading} onClick={submit}>
            {submitting ? "发布中..." : "发布心愿"}
          </button>
        </div>
      </div>
    </RequireAuth>
  );
}
