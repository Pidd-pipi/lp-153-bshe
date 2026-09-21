import Link from "next/link";
import { useRouter } from "next/router";
import { useAuth } from "@/hooks/useAuth";

export default function Navbar() {
  const { user, logout, isAuthed, isAdmin } = useAuth();
  const router = useRouter();

  const linkCls = (path: string) =>
    `px-3 py-2 rounded-lg text-sm font-medium transition ${
      router.pathname === path ? "bg-purple-100 text-purple-700" : "text-gray-600 hover:bg-purple-50 hover:text-purple-700"
    }`;

  return (
    <header className="sticky top-0 z-40 border-b border-purple-100 bg-white/90 backdrop-blur">
      <div className="mx-auto flex max-w-6xl items-center justify-between px-4 py-3">
        <Link href="/" className="flex items-center gap-2">
          <span className="text-2xl">✨</span>
          <span className="text-lg font-bold text-purple-700">虚拟心愿墙</span>
        </Link>

        <nav className="hidden items-center gap-1 md:flex">
          <Link href="/" className={linkCls("/")}>心愿广场</Link>
          <Link href="/discover" className={linkCls("/discover")}>发现广场</Link>
          {isAuthed() && (
            <>
              <Link href="/wishes/create" className={linkCls("/wishes/create")}>发布心愿</Link>
              <Link href="/capsules" className={linkCls("/capsules")}>时光胶囊</Link>
              <Link href="/profile" className={linkCls("/profile")}>个人主页</Link>
              {isAdmin() && <Link href="/audit" className={linkCls("/audit")}>审计日志</Link>}
            </>
          )}
        </nav>

        <div className="flex items-center gap-2">
          {isAuthed() ? (
            <>
              <Link href="/profile" className="hidden items-center gap-1.5 sm:flex">
                <span className="flex h-8 w-8 items-center justify-center rounded-full bg-gradient-to-br from-pink-400 to-purple-500 text-sm text-white">
                  {(user?.nickname || "我").slice(0, 1)}
                </span>
                <span className="text-sm text-gray-700">{user?.nickname}</span>
              </Link>
              <button
                className="btn-secondary !py-1.5"
                onClick={() => {
                  logout();
                  router.push("/");
                }}
              >
                退出
              </button>
            </>
          ) : (
            <>
              <Link href="/login" className="btn-secondary !py-1.5">登录</Link>
              <Link href="/register" className="btn-primary !py-1.5">注册</Link>
            </>
          )}
        </div>
      </div>
    </header>
  );
}
