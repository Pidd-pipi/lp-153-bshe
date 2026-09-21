// 统一请求封装：自动携带 JWT、统一响应解包、401 跳转、错误消息抽取。
export interface ApiResponse<T = unknown> {
  code: number;
  message: string;
  data: T;
}

const TOKEN_KEY = "wishwall_token";

export function getToken(): string {
  if (typeof window === "undefined") return "";
  return localStorage.getItem(TOKEN_KEY) || "";
}

export function setToken(token: string) {
  if (typeof window === "undefined") return;
  localStorage.setItem(TOKEN_KEY, token);
}

export function clearToken() {
  if (typeof window === "undefined") return;
  localStorage.removeItem(TOKEN_KEY);
}

export function redirectToLogin() {
  if (typeof window === "undefined") return;
  if (!window.location.pathname.startsWith("/login")) {
    window.location.href = "/login";
  }
}

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...(options.headers as Record<string, string> | undefined),
  };
  const token = getToken();
  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }
  const res = await fetch(`/api/v1${path}`, { ...options, headers });
  let body: ApiResponse<T>;
  try {
    body = await res.json();
  } catch {
    throw new Error(`请求失败(${res.status})`);
  }
  if (res.status === 401) {
    clearToken();
    redirectToLogin();
    throw new Error(body.message || "登录已过期");
  }
  if (body.code !== 0) {
    throw new Error(body.message || "请求失败");
  }
  return body.data;
}

export const http = {
  get<T>(path: string, params?: Record<string, string | number | undefined>) {
    const qs = new URLSearchParams();
    if (params) {
      Object.entries(params).forEach(([k, v]) => {
        if (v !== undefined && v !== "") qs.set(k, String(v));
      });
    }
    const suffix = qs.toString() ? `?${qs.toString()}` : "";
    return request<T>(`${path}${suffix}`);
  },
  post<T>(path: string, data?: unknown) {
    return request<T>(path, { method: "POST", body: data ? JSON.stringify(data) : undefined });
  },
  put<T>(path: string, data?: unknown) {
    return request<T>(path, { method: "PUT", body: data ? JSON.stringify(data) : undefined });
  },
  del<T>(path: string) {
    return request<T>(path, { method: "DELETE" });
  },
  upload<T>(path: string, formData: FormData) {
    const headers: Record<string, string> = {};
    const token = getToken();
    if (token) headers.Authorization = `Bearer ${token}`;
    return fetch(`/api/v1${path}`, { method: "POST", headers, body: formData })
      .then(async (res) => {
        const body = await res.json();
        if (body.code !== 0) throw new Error(body.message || "上传失败");
        return body.data as T;
      });
  },
};
