import { http } from "@/utils/request";
import type { PageResult } from "./wish";

export interface AuditLog {
  id: number;
  user_id: number;
  username?: string;
  action: string;
  entity_type: string;
  entity_id: string;
  detail: string;
  ip: string;
  request_id: string;
  created_at: string;
}

export const auditApi = {
  list: (params?: Record<string, string | number | undefined>) =>
    http.get<PageResult<AuditLog>>("/audit-logs", params),
};
