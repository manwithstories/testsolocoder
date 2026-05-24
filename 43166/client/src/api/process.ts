import request from '@/utils/request'
import {
  ProcessStep
} from '@/types'

export const processApi = {
  getSteps: (applicationId: number) => {
    return request.get<any, ProcessStep[]>(`/applications/${applicationId}/steps`)
  },

  updateStep: (applicationId: number, stepId: number, data: {
    status?: string
    remark?: string
    description?: string
    certificateFile?: string
  }) => {
    return request.put<any, null>(`/applications/${applicationId}/steps/${stepId}`, data)
  },

  startStep: (applicationId: number, stepId: number) => {
    return request.post<any, null>(`/applications/${applicationId}/steps/${stepId}/start`)
  },

  completeStep: (applicationId: number, stepId: number, certificateFile?: string, remark?: string) => {
    return request.post<any, null>(`/applications/${applicationId}/steps/${stepId}/complete`, {
      certificateFile,
      remark
    })
  },

  skipStep: (applicationId: number, stepId: number, remark: string) => {
    return request.post<any, null>(`/applications/${applicationId}/steps/${stepId}/skip`, { remark })
  },

  uploadCertificate: (applicationId: number, stepId: number, file: File) => {
    const formData = new FormData()
    formData.append('file', file)
    return request.post<any, { filePath: string; fileUrl: string; stepId: number }>(
      `/applications/${applicationId}/steps/${stepId}/upload`,
      formData,
      { headers: { 'Content-Type': 'multipart/form-data' } }
    )
  }
}
