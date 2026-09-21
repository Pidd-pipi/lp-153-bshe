import { create } from "zustand";
import type { Badge } from "@/api/badge";
import { badgeApi } from "@/api/badge";

interface BadgeState {
  badges: Badge[];
  loading: boolean;
  fetchMine: () => Promise<void>;
}

export const useBadgeStore = create<BadgeState>((set) => ({
  badges: [],
  loading: false,
  fetchMine: async () => {
    set({ loading: true });
    try {
      const items = await badgeApi.mine();
      set({ badges: items });
    } finally {
      set({ loading: false });
    }
  },
}));
