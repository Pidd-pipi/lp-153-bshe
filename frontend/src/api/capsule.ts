import { http } from "@/utils/request";
import type { PageResult } from "./wish";

export interface Capsule {
  id: number;
  user_id: number;
  title: string;
  content: string;
  image_urls: string[];
  audio_url?: string;
  unlock_at: string;
  status: string;
  unlocked_at?: string | null;
  created_at: string;
}

export const capsuleApi = {
  create: (payload: { title: string; content: string; image_urls?: string[]; audio_url?: string; unlock_at: string }) =>
    http.post<Capsule>("/capsules", payload),
  mine: (params?: Record<string, string | number | undefined>) =>
    http.get<PageResult<Capsule>>("/capsules/mine", params),
  detail: (id: number) => http.get<Capsule>(`/capsules/${id}`),
  remove: (id: number) => http.del<null>(`/capsules/${id}`),
};
