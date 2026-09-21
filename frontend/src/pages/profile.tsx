import { useCallback, useEffect, useState } from "react";
import type { Badge } from "@/api/badge";
import type { ClaimSummary, Wish } from "@/api/wish";
import { userApi } from "@/api/user";
import { wishApi } from "@/api/wish";
import { claimApi } from "@/api/claim";
import { badgeApi } from "@/api/badge";
import EmptyState from "@/components/EmptyState";
import ProgressBar from "@/components/ProgressBar";
import RequireAuth from "@/components/RequireAuth";
import StatusBadge from "@/components/StatusBadge";
import WishCard from "@/components/WishCard";
import { useToast } from "@/components/Toast";
import { useAuth } from "@/hooks/useAuth";
import { formatBadgeType, formatRole } from "@/utils/format";

export default function Profile() {
  const { user, refreshUser } = useAuth();
  const toast = useToast();
  const [myWishes, setMyWishes] = useState<Wish[]>([]);
  const [myClaims, setMyClaims] = useState<ClaimSummary[]>([]);
  const [badges, setBadges] = useState<Badge[]>([]);
  const [bio, setBio] = useState("");
  const [nickname, setNickname] = useState("");
  const [saving, setSaving] = useState(false);

  const load = useCallback(async () => {
    if (!user) return;
    try {
      const [wishes, claims, badgeItems] = await Promise.all([
        wishApi.mine({ page: 1, page_size: 10 }),
        claimApi.mine({ page: 1, page_size: 20 }),
        badgeApi.mine(),
      ]);
      setMyWishes(wishes.items);
      setMyClaims(claims.items);
      setBadges(badgeItems);
      setBio(user.bio || "");
      setNickname(user.nickname);
    } catch (e) {
      toast.show((e as Error).message, "error");
    }
  }, [user, toast]);

  useEffect(() => {
    load();
  }, [load]);

  const saveProfile = async () => {
    setSaving(true);
    try {
      await userApi.updateProfile({ nickname, bio });
      toast.show("资料已更新");
      refreshUser();
    } catch (e) {
      toast.show((e as Error).message, "error");
    } finally {
      setSaving(false);
    }
  };

  if (!user) return null;

  return (
    <RequireAuth>
      <div className="space-y-6">
        <section className="rounded-3xl bg-gradient-to-br from-purple-600 to-indigo-600 p-8 text-white">
          <div className="flex items-center gap-5">
            <div className="flex h-20 w-20 items-center justify-center rounded-full bg-white/20 text-3xl font-bold">
              {user.nickname.slice(0, 1)}
            </div>
            <div>
              <h1 className="text-2xl font-bold">{user.nickname}</h1>
              <p className="mt-1 text-sm text-purple-100">@{user.username} · {formatRole(user.role)}</p>
              <p className="mt-2 max-w-lg text-purple-50">{user.bio || "这个人很低调，还没写简介"}</p>
            </div>
          </div>
        </section>

        <div className="grid gap-6 lg:grid-cols-2">
          <section className="card space-y-4">
            <h2 className="text-lg font-semibold text-gray-800">🎖️ 我的成就徽章（{badges.length}）</h2>
            {badges.length === 0 ? (
              <EmptyState title="还没有徽章" description="完成心愿、认领心愿都会解锁徽章" icon="🎖️" />
            ) : (
              <div className="grid grid-cols-2 gap-3 sm:grid-cols-3">
                {badges.map((b) => (
                  <div key={b.id} className="flex flex-col items-center rounded-xl bg-gradient-to-b from-amber-50 to-white p-3 text-center">
                    <span className="text-3xl">{b.icon}</span>
                    <p className="mt-1 text-sm font-medium text-gray-800">{b.title}</p>
                    <p className="text-xs text-gray-400">{formatBadgeType(b.type)}</p>
                  </div>
                ))}
              </div>
            )}
          </section>

          <section className="card space-y-4">
            <h2 className="text-lg font-semibold text-gray-800">📝 个人资料</h2>
            <div>
              <label className="label">昵称</label>
              <input className="input" value={nickname} onChange={(e) => setNickname(e.target.value)} />
            </div>
            <div>
              <label className="label">简介</label>
              <textarea className="input min-h-[80px]" value={bio} onChange={(e) => setBio(e.target.value)} placeholder="介绍一下自己" />
            </div>
            <button className="btn-primary" disabled={saving} onClick={saveProfile}>{saving ? "保存中..." : "保存资料"}</button>
          </section>
        </div>

        <div className="grid gap-6 lg:grid-cols-2">
          <section className="space-y-4">
            <h2 className="text-lg font-semibold text-gray-800">✨ 我发布的心愿</h2>
            {myWishes.length === 0 ? (
              <EmptyState title="还没有发布心愿" icon="✨" />
            ) : (
              <div className="space-y-4">
                {myWishes.map((w) => <WishCard key={w.id} wish={w} />)}
              </div>
            )}
          </section>

          <section className="space-y-4">
            <h2 className="text-lg font-semibold text-gray-800">🤝 我认领的心愿</h2>
            {myClaims.length === 0 ? (
              <EmptyState title="还没有认领心愿" description="去广场逛逛，点亮别人的愿望" icon="🤝" />
            ) : (
              <div className="space-y-3">
                {myClaims.map((c) => (
                  <div key={c.id} className="card space-y-2 !py-4">
                    <div className="flex items-center justify-between gap-2">
                      <p className="font-medium text-gray-800">{c.wish_title}</p>
                      <StatusBadge status={c.status} />
                    </div>
                    <ProgressBar progress={c.progress} label="圆梦进度" />
                    {c.status === "pending_acceptance" && c.submission_note && (
                      <p className="text-xs text-orange-600">待验收 · 送交说明：{c.submission_note}</p>
                    )}
                    {c.status === "in_progress" && c.reject_reason && (
                      <p className="text-xs text-rose-600">驳回原因：{c.reject_reason}（可重新送交）</p>
                    )}
                  </div>
                ))}
              </div>
            )}
          </section>
        </div>
      </div>
    </RequireAuth>
  );
}
