import { useCallback, useEffect, useMemo, useState } from "react";
import Link from "next/link";
import type { Wish } from "@/api/wish";
import { wishApi } from "@/api/wish";
import EmptyState from "@/components/EmptyState";
import WishCard from "@/components/WishCard";
import { useToast } from "@/components/Toast";
import { useAuth } from "@/hooks/useAuth";
import { CATEGORY_TEXT, WISH_STATUS_TEXT } from "@/constants";

const CATEGORY_OPTIONS = Object.entries(CATEGORY_TEXT);
const STATUS_OPTIONS = Object.entries(WISH_STATUS_TEXT);

export default function Home() {
  const { isAuthed } = useAuth();
  const toast = useToast();
  const [wishes, setWishes] = useState<Wish[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [keyword, setKeyword] = useState("");
  const [category, setCategory] = useState("");
  const [status, setStatus] = useState("");
  const [sort, setSort] = useState("new");
  const [page, setPage] = useState(1);

  const fetchWishes = useCallback(async () => {
    setLoading(true);
    try {
      const data = await wishApi.list({ keyword, category, status, sort, page, page_size: 9 });
      setWishes(data.items);
      setTotal(data.total);
    } catch (e) {
      toast.show((e as Error).message, "error");
    } finally {
      setLoading(false);
    }
  }, [keyword, category, status, sort, page, toast]);

  useEffect(() => {
    fetchWishes();
  }, [fetchWishes]);

  const like = useCallback(async (id: number) => {
    if (!isAuthed()) {
      toast.show("请先登录后再点赞", "error");
      return;
    }
    try {
      await wishApi.like(id);
      setWishes((prev) => prev.map((w) => (w.id === id ? { ...w, likes_count: w.likes_count + 1 } : w)));
      toast.show("点赞成功 💖");
    } catch (e) {
      toast.show((e as Error).message, "error");
    }
  }, [isAuthed, toast]);

  const totalPages = useMemo(() => Math.max(1, Math.ceil(total / 9)), [total]);

  return (
    <div className="space-y-6">
      <section className="relative overflow-hidden rounded-3xl bg-gradient-to-br from-pink-500 via-purple-600 to-indigo-600 p-8 text-white">
        <div className="absolute right-8 top-6 text-6xl opacity-30 animate-float">🌙</div>
        <div className="absolute right-24 top-14 text-4xl opacity-20 animate-float">⭐</div>
        <h1 className="text-3xl font-bold">心愿广场</h1>
        <p className="mt-2 max-w-xl text-purple-50">
          写下心愿，也去点亮别人的心愿。认领、陪伴、祝福，让每个愿望都被温柔以待。
        </p>
        {isAuthed() && (
          <Link href="/wishes/create" className="mt-4 inline-block rounded-xl bg-white px-4 py-2 text-sm font-medium text-purple-700 hover:bg-purple-50">
            ✏️ 发布我的心愿
          </Link>
        )}
      </section>

      <section className="card space-y-3">
        <div className="flex flex-wrap items-center gap-3">
          <input
            className="input max-w-xs"
            placeholder="搜索心愿关键词..."
            value={keyword}
            onChange={(e) => { setKeyword(e.target.value); setPage(1); }}
          />
          <select className="input max-w-[140px]" value={category} onChange={(e) => { setCategory(e.target.value); setPage(1); }}>
            <option value="">全部分类</option>
            {CATEGORY_OPTIONS.map(([k, v]) => <option key={k} value={k}>{v}</option>)}
          </select>
          <select className="input max-w-[140px]" value={status} onChange={(e) => { setStatus(e.target.value); setPage(1); }}>
            <option value="">全部状态</option>
            {STATUS_OPTIONS.map(([k, v]) => <option key={k} value={k}>{v}</option>)}
          </select>
          <select className="input max-w-[140px]" value={sort} onChange={(e) => { setSort(e.target.value); setPage(1); }}>
            <option value="new">最新发布</option>
            <option value="hot">热度最高</option>
          </select>
        </div>
      </section>

      {loading ? (
        <p className="py-16 text-center text-gray-400">加载中...</p>
      ) : wishes.length === 0 ? (
        <EmptyState title="还没有心愿" description="来发布第一个心愿吧" icon="✨" />
      ) : (
        <>
          <div className="grid gap-5 sm:grid-cols-2 lg:grid-cols-3">
            {wishes.map((w) => (
              <WishCard key={w.id} wish={w} onLike={like} showActions />
            ))}
          </div>
          <div className="flex items-center justify-center gap-4 pt-2">
            <button className="btn-secondary" disabled={page <= 1} onClick={() => setPage((p) => p - 1)}>上一页</button>
            <span className="text-sm text-gray-500">{page} / {totalPages}</span>
            <button className="btn-secondary" disabled={page >= totalPages} onClick={() => setPage((p) => p + 1)}>下一页</button>
          </div>
        </>
      )}
    </div>
  );
}
