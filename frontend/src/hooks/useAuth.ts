import { useCallback, useEffect } from "react";
import { userApi } from "@/api/user";
import { useAuthStore } from "@/stores/authStore";

// useAuth：登录态恢复 + 登出，供路由守卫与按钮显隐复用。
export function useAuth() {
  const { user, token, setUser, logout, isAuthed, isAdmin } = useAuthStore();

  const refreshUser = useCallback(async () => {
    if (!isAuthed()) return;
    try {
      const me = await userApi.me();
      setUser(me);
    } catch {
      // 401 已由 request 拦截器处理
    }
  }, [isAuthed, setUser]);

  useEffect(() => {
    if (isAuthed() && !user) {
      refreshUser();
    }
  }, [isAuthed, user, refreshUser]);

  return { user, token, logout, isAuthed, isAdmin, refreshUser };
}
