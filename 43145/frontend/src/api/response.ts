import { get, post, del } from '@/utils/request'
import type { SurveyResponse, DistributionLink, Invitation, PaginatedData } from '@/types'

export const responseApi = {
  list(surveyId: number, params: any) {
    return get<PaginatedData<SurveyResponse>>(`/surveys/${surveyId}/responses`, { params })
  },

  getById(id: number) {
    return get(`/responses/${id}`)
  },

  remove(id: number) {
    return del(`/responses/${id}`)
  }
}

export const distributionApi = {
  createLink(surveyId: number, data: any) {
    return post<DistributionLink>(`/surveys/${surveyId}/distribution/link`, data)
  },

  listLinks(surveyId: number) {
    return get<DistributionLink[]>(`/surveys/${surveyId}/distribution/links`)
  },

  sendInvitations(surveyId: number, data: any) {
    return post(`/surveys/${surveyId}/distribution/invitations`, data)
  },

  listInvitations(surveyId: number, params?: any) {
    return get<PaginatedData<Invitation>>(`/surveys/${surveyId}/distribution/invitations`, { params })
  },

  deleteLink(id: number) {
    return del(`/distribution/${id}`)
  }
}

export const publicApi = {
  getSurveyByToken(token: string) {
    return get(`/public/survey/${token}`)
  },

  validateAccess(surveyId: number, data: { password?: string; session_id?: string }) {
    return post(`/public/survey/${surveyId}/validate`, data)
  },

  startResponse(surveyId: number, data?: any, token?: string) {
    return post(`/public/survey/${surveyId}/start${token ? `?token=${token}` : ''}`, data)
  },

  saveProgress(surveyId: number, data: any) {
    return post(`/public/survey/${surveyId}/save`, data)
  },

  submitResponse(surveyId: number, data: any) {
    return post(`/public/survey/${surveyId}/submit`, data)
  }
}
