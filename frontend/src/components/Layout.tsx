import type { ReactNode } from "react";
import Navbar from "./Navbar";

export default function Layout({ children }: { children: ReactNode }) {
  return (
    <div className="min-h-screen">
      <Navbar />
      <main className="mx-auto max-w-6xl px-4 py-6">{children}</main>
      <footer className="border-t border-purple-100 py-6 text-center text-xs text-gray-400">
        ✨ 虚拟心愿墙 · 用善意连接彼此
      </footer>
    </div>
  );
}
