import { http } from "@/utils/request";

export interface User {
  id: number;
  username: string;
  email: string;
  nickname: string;
  avatar: string;
  bio: string;
  role: string;
  status: string;
  created_at: string;
}

export interface LoginResult {
  token: string;
  user: User;
}

export interface RegisterPayload {
  username: string;
  email: string;
  password: string;
  nickname?: string;
}

export interface UpdateProfilePayload {
  nickname?: string;
  avatar?: string;
  bio?: string;
}

export const userApi = {
  register: (payload: RegisterPayload) => http.post<User>("/auth/register", payload),
  login: (account: string, password: string) =>
    http.post<LoginResult>("/auth/login", { account, password }),
  me: () => http.get<User>("/users/me"),
  getById: (id: number) => http.get<User>(`/users/${id}`),
  updateProfile: (payload: UpdateProfilePayload) =>
    http.put<User>("/users/me", payload),
};
