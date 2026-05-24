import { apiDownload } from './client'

export const exportAdoptionReport = (params?: { start_date?: string; end_date?: string }): Promise<Blob> => {
  return apiDownload('/export/adoption', params)
}

export const exportHealthReport = (petId: number): Promise<Blob> => {
  return apiDownload(`/export/health/${petId}`)
}
