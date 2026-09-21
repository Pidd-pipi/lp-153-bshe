import { create } from "zustand";
import type { PageResult, Wish, WishDetail } from "@/api/wish";
import { wishApi } from "@/api/wish";

interface WishState {
  wishes: Wish[];
  total: number;
  loading: boolean;
  detail: WishDetail | null;
  fetchList: (params?: Record<string, string | number | undefined>) => Promise<void>;
  fetchDetail: (id: number) => Promise<void>;
}

export const useWishStore = create<WishState>((set) => ({
  wishes: [],
  total: 0,
  loading: false,
  detail: null,
  fetchList: async (params) => {
    set({ loading: true });
    try {
      const data: PageResult<Wish> = await wishApi.list(params);
      set({ wishes: data.items, total: data.total });
    } finally {
      set({ loading: false });
    }
  },
  fetchDetail: async (id) => {
    const detail = await wishApi.detail(id);
    set({ detail });
  },
}));
