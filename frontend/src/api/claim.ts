import { http } from "@/utils/request";
import type { ClaimSummary, PageResult } from "./wish";

export const claimApi = {
  claim: (wishId: number) => http.post<ClaimSummary>(`/wishes/${wishId}/claim`),
  getByWish: (wishId: number) => http.get<ClaimSummary>(`/wishes/${wishId}/claim`),
  mine: (params?: Record<string, string | number | undefined>) =>
    http.get<PageResult<ClaimSummary>>("/claims/mine", params),
  updateProgress: (claimId: number, payload: { progress: number; note?: string; is_milestone?: boolean }) =>
    http.put<ClaimSummary>(`/claims/${claimId}/progress`, payload),
  complete: (claimId: number, payload: { note?: string }) =>
    http.post<ClaimSummary>(`/claims/${claimId}/complete`, payload),
};
