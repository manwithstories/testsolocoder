import { request } from './client'
import {
  RoasterCertification,
  PaginatedData,
  User,
  Product,
  RoastingRecord,
} from '@/types'

export const certificationApi = {
  apply: (data: {
    cert_name: string
    cert_number: string
    org_name: string
    cert_file?: string
    experience: string
    specialty?: string
  }) => request.post<RoasterCertification>('/certification/apply', data),

  updateApplication: (data: {
    cert_name: string
    cert_number: string
    org_name: string
    cert_file?: string
    experience: string
    specialty?: string
  }) => request.put('/certification/apply', data),

  getMyCertification: () =>
    request.get<RoasterCertification>('/certification/my'),

  list: (params?: { page?: number; page_size?: number; status?: string }) =>
    request.get<PaginatedData<RoasterCertification>>('/certification', { params }),

  review: (id: number, data: { status: string; review_comment?: string }) =>
    request.post(`/certification/${id}/review`, data),

  getRoasterProfile: (id: number) =>
    request.get<{
      user: User
      certification?: RoasterCertification
      products: Product[]
      roasting_records: RoastingRecord[]
      total_products: number
      total_roasts: number
      avg_score: number
    }>(`/roasters/${id}`),

  listCertifiedRoasters: (params?: { page?: number; page_size?: number; keyword?: string }) =>
    request.get<PaginatedData<User>>('/roasters', { params }),
}
