import { create } from "zustand";
import type { User } from "@/api/user";
import { clearToken, getToken, setToken } from "@/utils/request";

interface AuthState {
  user: User | null;
  token: string;
  loading: boolean;
  setUser: (user: User | null) => void;
  setToken: (token: string) => void;
  logout: () => void;
  isAuthed: () => boolean;
  isAdmin: () => boolean;
}

export const useAuthStore = create<AuthState>((set, get) => ({
  user: null,
  token: typeof window !== "undefined" ? getToken() : "",
  loading: false,
  setUser: (user) => set({ user }),
  setToken: (token) => {
    setToken(token);
    set({ token });
  },
  logout: () => {
    clearToken();
    set({ user: null, token: "" });
  },
  isAuthed: () => Boolean(get().token),
  isAdmin: () => get().user?.role === "admin",
}));
