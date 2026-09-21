import Link from "next/link";
import type { Wish } from "@/api/wish";
import StatusBadge from "./StatusBadge";
import { formatCategory, formatDeadline, formatDifficulty } from "@/utils/format";

interface WishCardProps {
  wish: Wish;
  onLike?: (id: number) => void;
  showActions?: boolean;
}

// 心愿卡片：心愿广场/发现广场/我的心愿跨页复用。
export default function WishCard({ wish, onLike, showActions = false }: WishCardProps) {
  return (
    <div className="card flex flex-col gap-3 transition hover:shadow-md">
      <div className="flex items-start justify-between gap-2">
        <div className="flex items-center gap-2">
          <div className="flex h-9 w-9 items-center justify-center rounded-full bg-gradient-to-br from-pink-400 to-purple-500 text-white">
            {wish.author_nickname ? wish.author_nickname.slice(0, 1) : "心"}
          </div>
          <div>
            <p className="text-sm font-medium text-gray-800">{wish.author_nickname || "神秘人"}</p>
            <p className="text-xs text-gray-400">{wish.created_at?.slice(0, 16)}</p>
          </div>
        </div>
        <StatusBadge status={wish.status} />
      </div>

      <Link href={`/wishes/detail?id=${wish.id}`} className="block">
        <h3 className="text-base font-semibold text-gray-900 hover:text-pink-600">{wish.title}</h3>
        <p className="mt-1 line-clamp-2 text-sm text-gray-600">{wish.content}</p>
      </Link>

      {wish.image_urls && wish.image_urls.length > 0 && (
        <img src={wish.image_urls[0]} alt={wish.title} className="h-36 w-full rounded-xl object-cover" />
      )}

      <div className="flex flex-wrap gap-1.5 text-xs">
        <span className="rounded-full bg-purple-50 px-2 py-0.5 text-purple-600">{formatCategory(wish.category)}</span>
        <span className="rounded-full bg-sky-50 px-2 py-0.5 text-sky-600">{formatDifficulty(wish.difficulty)}</span>
        <span className="rounded-full bg-amber-50 px-2 py-0.5 text-amber-600">截止 {formatDeadline(wish.expected_deadline)}</span>
      </div>

      {showActions && (
        <div className="flex items-center justify-between border-t border-purple-50 pt-3">
          <button
            onClick={() => onLike?.(wish.id)}
            className="text-sm text-gray-500 hover:text-pink-500"
          >
            ❤️ {wish.likes_count || 0}
          </button>
          <Link href={`/wishes/detail?id=${wish.id}`} className="text-sm text-purple-600 hover:text-purple-700">
            查看详情 →
          </Link>
        </div>
      )}
    </div>
  );
}
