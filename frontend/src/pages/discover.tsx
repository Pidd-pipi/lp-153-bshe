import { useCallback, useEffect, useState } from "react";
import type { FulfillerStat, Wish } from "@/api/wish";
import { wishApi } from "@/api/wish";
import EmptyState from "@/components/EmptyState";
import WishCard from "@/components/WishCard";

export default function Discover() {
  const [stories, setStories] = useState<Wish[]>([]);
  const [leaderboard, setLeaderboard] = useState<FulfillerStat[]>([]);
  const [loading, setLoading] = useState(true);

  const load = useCallback(async () => {
    setLoading(true);
    try {
      const [storyData, lb] = await Promise.all([
        wishApi.discover({ page: 1, page_size: 6 }),
        wishApi.leaderboard(),
      ]);
      setStories(storyData.items);
      setLeaderboard(lb);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    load();
  }, [load]);

  return (
    <div className="space-y-8">
      <section className="rounded-3xl bg-gradient-to-br from-emerald-400 to-teal-500 p-8 text-white">
        <h1 className="text-3xl font-bold">🌍 发现广场</h1>
        <p className="mt-2 text-emerald-50">看看最近被实现的愿望故事，和那些闪闪发光的圆梦人</p>
      </section>

      <section>
        <h2 className="mb-4 text-xl font-bold text-gray-800">🏆 热门圆梦人排行榜</h2>
        {loading ? (
          <p className="py-8 text-center text-gray-400">加载中...</p>
        ) : leaderboard.length === 0 ? (
          <EmptyState title="榜单虚位以待" description="第一个完成心愿的人即将上榜" icon="🏆" />
        ) : (
          <div className="card divide-y divide-purple-50">
            {leaderboard.map((row, i) => (
              <div key={row.user_id} className="flex items-center gap-4 py-3">
                <span className={`w-8 text-center text-xl ${i === 0 ? "text-amber-400" : i === 1 ? "text-gray-400" : i === 2 ? "text-amber-700" : "text-gray-300"}`}>
                  {["🥇", "🥈", "🥉"][i] || `${i + 1}`}
                </span>
                <div className="flex h-10 w-10 items-center justify-center rounded-full bg-gradient-to-br from-pink-400 to-purple-500 text-white">
                  {row.nickname ? row.nickname.slice(0, 1) : "?"}
                </div>
                <span className="font-medium text-gray-800">{row.nickname}</span>
                <span className="ml-auto text-sm text-gray-500">完成 {row.completed_count} 个心愿</span>
              </div>
            ))}
          </div>
        )}
      </section>

      <section>
        <h2 className="mb-4 text-xl font-bold text-gray-800">🎊 最新完成的心愿故事</h2>
        {loading ? (
          <p className="py-8 text-center text-gray-400">加载中...</p>
        ) : stories.length === 0 ? (
          <EmptyState title="还没有完成的故事" description="完成心愿后，故事会在这里发光" icon="🎊" />
        ) : (
          <div className="grid gap-5 sm:grid-cols-2 lg:grid-cols-3">
            {stories.map((w) => <WishCard key={w.id} wish={w} />)}
          </div>
        )}
      </section>
    </div>
  );
}
