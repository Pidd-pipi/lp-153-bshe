import { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/router";
import { userApi } from "@/api/user";
import { useToast } from "@/components/Toast";

export default function Register() {
  const router = useRouter();
  const toast = useToast();
  const [form, setForm] = useState({ username: "", email: "", password: "", nickname: "" });
  const [loading, setLoading] = useState(false);

  const set = (k: keyof typeof form) => (e: React.ChangeEvent<HTMLInputElement>) =>
    setForm((f) => ({ ...f, [k]: e.target.value }));

  const submit = async () => {
    if (form.password.length < 6) {
      toast.show("密码至少 6 位", "error");
      return;
    }
    setLoading(true);
    try {
      await userApi.register(form);
      toast.show("注册成功，去登录吧 🎉");
      router.push("/login");
    } catch (e) {
      toast.show((e as Error).message, "error");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="mx-auto max-w-md py-10">
      <div className="card space-y-5">
        <div className="text-center">
          <div className="text-4xl">✨</div>
          <h1 className="mt-2 text-2xl font-bold text-purple-700">注册</h1>
          <p className="text-sm text-gray-400">加入心愿墙，许下你的第一个愿望</p>
        </div>
        <div className="space-y-3">
          <div>
            <label className="label">用户名</label>
            <input className="input" value={form.username} onChange={set("username")} placeholder="3-50 个字符" />
          </div>
          <div>
            <label className="label">昵称（选填）</label>
            <input className="input" value={form.nickname} onChange={set("nickname")} placeholder="展示给其他用户的名字" />
          </div>
          <div>
            <label className="label">邮箱</label>
            <input className="input" type="email" value={form.email} onChange={set("email")} placeholder="you@example.com" />
          </div>
          <div>
            <label className="label">密码</label>
            <input className="input" type="password" value={form.password} onChange={set("password")} placeholder="至少 6 位" onKeyDown={(e) => e.key === "Enter" && submit()} />
          </div>
          <button className="btn-primary w-full" disabled={loading} onClick={submit}>
            {loading ? "注册中..." : "注 册"}
          </button>
        </div>
        <p className="text-center text-sm text-gray-500">
          已有账号？<Link href="/login" className="text-purple-600 hover:underline">去登录</Link>
        </p>
      </div>
    </div>
  );
}
