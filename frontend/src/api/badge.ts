import { http } from "@/utils/request";
import type { FulfillerStat } from "./wish";

export interface Badge {
  id: number;
  user_id: number;
  type: string;
  title: string;
  description: string;
  icon: string;
  earned_at: string;
}

export const badgeApi = {
  mine: () => http.get<Badge[]>("/badges/mine"),
  leaderboard: () => http.get<FulfillerStat[]>("/badges/leaderboard"),
};
