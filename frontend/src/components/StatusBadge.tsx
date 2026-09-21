import { WISH_STATUS_STYLE, WISH_STATUS_TEXT, CAPSULE_STATUS_TEXT, ACCEPTANCE_STATUS_STYLE, ACCEPTANCE_STATUS_TEXT } from "@/constants";

interface StatusBadgeProps {
  status: string;
  kind?: "wish" | "capsule" | "acceptance";
}

// 状态徽标：心愿状态 / 胶囊状态 / 验收状态跨页复用。
export default function StatusBadge({ status, kind = "wish" }: StatusBadgeProps) {
  const text = kind === "capsule"
    ? CAPSULE_STATUS_TEXT[status] || status
    : kind === "acceptance"
      ? ACCEPTANCE_STATUS_TEXT[status] || status
      : WISH_STATUS_TEXT[status] || status;
  const style = kind === "capsule"
    ? status === "unlocked" ? "bg-emerald-100 text-emerald-700" : "bg-gray-100 text-gray-600"
    : kind === "acceptance"
      ? ACCEPTANCE_STATUS_STYLE[status] || "bg-gray-100 text-gray-600"
      : WISH_STATUS_STYLE[status] || "bg-gray-100 text-gray-600";
  return (
    <span className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium ${style}`}>
      {text}
    </span>
  );
}
