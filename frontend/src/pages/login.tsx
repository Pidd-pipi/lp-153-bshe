import { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/router";
import { userApi } from "@/api/user";
import { useToast } from "@/components/Toast";
import { useAuthStore } from "@/stores/authStore";

export default function Login() {
  const router = useRouter();
  const toast = useToast();
  const setToken = useAuthStore((s) => s.setToken);
  const setUser = useAuthStore((s) => s.setUser);
  const [account, setAccount] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);

  const submit = async () => {
    if (!account || !password) {
      toast.show("请输入账号和密码", "error");
      return;
    }
    setLoading(true);
    try {
      const result = await userApi.login(account, password);
      setToken(result.token);
      setUser(result.user);
      toast.show("登录成功，欢迎回来 🌟");
      router.push("/");
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
          <div className="text-4xl">🔐</div>
          <h1 className="mt-2 text-2xl font-bold text-purple-700">登录</h1>
          <p className="text-sm text-gray-400">回到心愿墙，继续圆梦之旅</p>
        </div>
        <div className="space-y-3">
          <div>
            <label className="label">用户名 / 邮箱</label>
            <input className="input" value={account} onChange={(e) => setAccount(e.target.value)} placeholder="请输入用户名或邮箱" />
          </div>
          <div>
            <label className="label">密码</label>
            <input className="input" type="password" value={password} onChange={(e) => setPassword(e.target.value)} placeholder="请输入密码" onKeyDown={(e) => e.key === "Enter" && submit()} />
          </div>
          <button className="btn-primary w-full" disabled={loading} onClick={submit}>
            {loading ? "登录中..." : "登 录"}
          </button>
        </div>
        <p className="text-center text-sm text-gray-500">
          还没有账号？<Link href="/register" className="text-purple-600 hover:underline">立即注册</Link>
        </p>
      </div>
    </div>
  );
}
