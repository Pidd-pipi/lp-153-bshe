import { http } from "@/utils/request";
import type { PageResult } from "./wish";

export interface Blessing {
  id: number;
  wish_id: number;
  user_id: number;
  sender_name?: string;
  content: string;
  gift_emoji?: string;
  is_celebrating: boolean;
  created_at: string;
}

export const blessingApi = {
  create: (wishId: number, payload: { content: string; gift_emoji?: string }) =>
    http.post<Blessing>(`/wishes/${wishId}/blessings`, payload),
  list: (wishId: number, params?: Record<string, string | number | undefined>) =>
    http.get<PageResult<Blessing>>(`/wishes/${wishId}/blessings`, params),
};
