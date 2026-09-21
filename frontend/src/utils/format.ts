import {
  WISH_STATUS_TEXT,
  VISIBILITY_TEXT,
  DIFFICULTY_TEXT,
  CATEGORY_TEXT,
  CAPSULE_STATUS_TEXT,
  BADGE_TYPE_TEXT,
  ROLE_TEXT,
  ACCEPTANCE_STATUS_TEXT,
} from "@/constants";

export function formatDate(input?: string | null): string {
  if (!input) return "-";
  const d = new Date(input);
  if (Number.isNaN(d.getTime())) return input;
  const pad = (n: number) => String(n).padStart(2, "0");
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`;
}

export function formatDeadline(input?: string | null): string {
  if (!input) return "不限时";
  return input.slice(0, 10);
}

export function formatWishStatus(s: string) {
  return WISH_STATUS_TEXT[s] || s;
}

export function formatVisibility(v: string) {
  return VISIBILITY_TEXT[v] || v;
}

export function formatDifficulty(d: string) {
  return DIFFICULTY_TEXT[d] || d;
}

export function formatCategory(c: string) {
  return CATEGORY_TEXT[c] || c;
}

export function formatCapsuleStatus(s: string) {
  return CAPSULE_STATUS_TEXT[s] || s;
}

export function formatAcceptanceStatus(s: string) {
  return ACCEPTANCE_STATUS_TEXT[s] || s;
}

export function formatBadgeType(t: string) {
  return BADGE_TYPE_TEXT[t] || t;
}

export function formatRole(r: string) {
  return ROLE_TEXT[r] || r;
}

export function greet() {
  const h = new Date().getHours();
  if (h < 6) return "夜深了，许个愿吧";
  if (h < 12) return "早上好，今天也有小确幸";
  if (h < 18) return "下午好，来点亮一盏心愿";
  return "晚上好，愿你好梦";
}
