import { useEffect, type ReactNode } from "react";
import { useRouter } from "next/router";
import { useAuth } from "@/hooks/useAuth";

interface RequireAuthProps {
  children: ReactNode;
  adminOnly?: boolean;
}

// 路由守卫：未登录跳转登录页；adminOnly 页面校验管理员角色。
export default function RequireAuth({ children, adminOnly = false }: RequireAuthProps) {
  const router = useRouter();
  const { isAuthed, isAdmin, user } = useAuth();

  useEffect(() => {
    if (!isAuthed()) {
      router.replace("/login");
      return;
    }
    if (adminOnly && user && !isAdmin()) {
      router.replace("/");
    }
  }, [isAuthed, isAdmin, user, router, adminOnly]);

  if (!isAuthed()) return null;
  if (adminOnly && user && !isAdmin()) return null;
  return <>{children}</>;
}
