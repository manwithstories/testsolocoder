import { get, post, put, del } from '@/utils/request'
import type { Question } from '@/types'

export const questionApi = {
  list(surveyId: number) {
    return get<Question[]>(`/surveys/${surveyId}/questions`)
  },

  create(surveyId: number, data: any) {
    return post<Question>(`/surveys/${surveyId}/questions`, data)
  },

  batchCreate(surveyId: number, data: any) {
    return post(`/surveys/${surveyId}/questions/batch`, data)
  },

  getById(id: number) {
    return get<Question>(`/questions/${id}`)
  },

  update(id: number, data: any) {
    return put(`/questions/${id}`, data)
  },

  remove(id: number) {
    return del(`/questions/${id}`)
  },

  reorder(id: number, orderIndex: number) {
    return put(`/questions/${id}/reorder`, { order_index: orderIndex })
  }
}
