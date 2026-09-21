// 与后端 internal/constants 对应的共享枚举（README 中列出前后端出现位置）。

export const ROLE = {
  USER: "user",
  ADMIN: "admin",
} as const;
export const ROLE_TEXT: Record<string, string> = {
  user: "普通用户",
  admin: "管理员",
};

export const WISH_STATUS = {
  PENDING: "pending",
  CLAIMED: "claimed",
  IN_PROGRESS: "in_progress",
  COMPLETED: "completed",
} as const;
export const WISH_STATUS_TEXT: Record<string, string> = {
  pending: "待认领",
  claimed: "已被认领",
  in_progress: "圆梦中",
  completed: "已完成",
};
export const WISH_STATUS_STYLE: Record<string, string> = {
  pending: "bg-amber-100 text-amber-700",
  claimed: "bg-sky-100 text-sky-700",
  in_progress: "bg-violet-100 text-violet-700",
  completed: "bg-emerald-100 text-emerald-700",
};

export const VISIBILITY = {
  PUBLIC: "public",
  FRIEND: "friend",
  ANONYMOUS: "anonymous",
} as const;
export const VISIBILITY_TEXT: Record<string, string> = {
  public: "公开",
  friend: "好友可见",
  anonymous: "匿名",
};

export const DIFFICULTY = {
  EASY: "easy",
  MEDIUM: "medium",
  HARD: "hard",
} as const;
export const DIFFICULTY_TEXT: Record<string, string> = {
  easy: "简单",
  medium: "中等",
  hard: "困难",
};

export const CATEGORY = {
  STUDY: "study",
  TRAVEL: "travel",
  EMOTION: "emotion",
  CAREER: "career",
  LIFE: "life",
  OTHER: "other",
} as const;
export const CATEGORY_TEXT: Record<string, string> = {
  study: "学习成长",
  travel: "旅行探险",
  emotion: "情感陪伴",
  career: "职业发展",
  life: "生活小确幸",
  other: "其他",
};

export const CAPSULE_STATUS = {
  LOCKED: "locked",
  UNLOCKED: "unlocked",
} as const;
export const CAPSULE_STATUS_TEXT: Record<string, string> = {
  locked: "未解锁",
  unlocked: "已解锁",
};

export const BADGE_TYPE = {
  FIRST_WISH: "first_wish",
  FIRST_CLAIM: "first_claim",
  FIRST_BLESSING: "first_blessing",
  TEN_COMPLETIONS: "ten_completions",
  WISH_MASTER: "wish_master",
} as const;
export const BADGE_TYPE_TEXT: Record<string, string> = {
  first_wish: "首次许愿",
  first_claim: "首次认领",
  first_blessing: "首次祝福",
  ten_completions: "十次圆梦",
  wish_master: "圆梦大师",
};

export const GIFT_EMOJIS = ["🎁", "🌸", "🌻", "⭐", "🌈", "🕊️", "💖", "🍀"];
