import { http } from "@/utils/request";

export interface Wish {
  id: number;
  user_id: number;
  author_nickname?: string;
  author_avatar?: string;
  title: string;
  content: string;
  image_urls: string[];
  category: string;
  visibility: string;
  difficulty: string;
  expected_deadline?: string | null;
  status: string;
  likes_count: number;
  completion_note?: string;
  is_anonymous: boolean;
  created_at: string;
  updated_at: string;
}

export interface ClaimSummary {
  id: number;
  wish_id: number;
  wish_title?: string;
  user_id: number;
  fulfiller_name?: string;
  progress: number;
  latest_note?: string;
  status: string;
  milestone_count: number;
  submission_note?: string;
  submitted_at?: string | null;
  reject_reason?: string;
  reviewed_at?: string | null;
  created_at: string;
  updated_at: string;
}

export interface WishDetail extends Wish {
  claim?: ClaimSummary | null;
  blessing_count: number;
}

export interface CreateWishPayload {
  title: string;
  content: string;
  image_urls?: string[];
  category: string;
  visibility: string;
  difficulty: string;
  expected_deadline?: string | null;
  is_anonymous?: boolean;
}

export interface UpdateWishPayload {
  title?: string;
  content?: string;
  image_urls?: string[];
  category?: string;
  visibility?: string;
  difficulty?: string;
  expected_deadline?: string | null;
  is_anonymous?: boolean;
}

export interface PageResult<T> {
  items: T[];
  total: number;
  page: number;
  page_size: number;
}

export const wishApi = {
  create: (payload: CreateWishPayload) => http.post<Wish>("/wishes", payload),
  list: (params?: Record<string, string | number | undefined>) =>
    http.get<PageResult<Wish>>("/wishes", params),
  mine: (params?: Record<string, string | number | undefined>) =>
    http.get<PageResult<Wish>>("/wishes/mine", params),
  detail: (id: number) => http.get<WishDetail>(`/wishes/${id}`),
  update: (id: number, payload: UpdateWishPayload) =>
    http.put<Wish>(`/wishes/${id}`, payload),
  remove: (id: number) => http.del<null>(`/wishes/${id}`),
  like: (id: number) => http.post<{ wish_id: number }>(`/wishes/${id}/like`),
  discover: (params?: Record<string, string | number | undefined>) =>
    http.get<PageResult<Wish>>("/discover", params),
  leaderboard: () => http.get<FulfillerStat[]>("/discover/leaderboard"),
};

export interface FulfillerStat {
  user_id: number;
  nickname: string;
  avatar: string;
  completed_count: number;
}
