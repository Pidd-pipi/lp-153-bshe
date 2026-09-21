import { http } from "@/utils/request";

export interface AcceptanceSummary {
  id: number;
  wish_id: number;
  claim_id: number;
  user_id: number;
  fulfiller_name?: string;
  note: string;
  status: string;
  reject_reason?: string;
  submitted_at?: string;
  reviewed_at?: string;
  reviewed_by?: number;
  created_at: string;
  updated_at: string;
}

export const acceptanceApi = {
  submit: (claimId: number, payload: { note: string }) =>
    http.post<AcceptanceSummary>(`/claims/${claimId}/acceptance`, payload),
  confirm: (acceptanceId: number) =>
    http.post<AcceptanceSummary>(`/acceptances/${acceptanceId}/confirm`),
  reject: (acceptanceId: number, payload: { reason: string }) =>
    http.post<AcceptanceSummary>(`/acceptances/${acceptanceId}/reject`, payload),
  getByWish: (wishId: number) => http.get<AcceptanceSummary>(`/wishes/${wishId}/acceptance`),
};
