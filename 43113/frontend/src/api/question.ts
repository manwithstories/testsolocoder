import { request } from '@/utils/request'
import type {
  Question,
  QuestionDetail,
  QuestionListQuery,
  CreateQuestionRequest,
  CreateAnswerRequest,
  Answer,
  CreateCommentRequest,
  Comment,
  Category,
  Tag,
  PageResult
} from '@/types'

export const questionApi = {
  getQuestionList: (params?: QuestionListQuery) => {
    return request.get<PageResult<Question>>('/public/questions', { params })
  },

  getQuestion: (id: number) => {
    return request.get<QuestionDetail>(`/public/questions/${id}`)
  },

  createQuestion: (data: CreateQuestionRequest) => {
    return request.post<Question>('/questions', data)
  },

  updateQuestion: (id: number, data: { title?: string; content?: string }) => {
    return request.put(`/questions/${id}`, data)
  },

  deleteQuestion: (id: number) => {
    return request.delete(`/questions/${id}`)
  },

  likeQuestion: (id: number) => {
    return request.post(`/questions/${id}/like`)
  },

  acceptAnswer: (questionId: number, answerId: number) => {
    return request.post(`/questions/${questionId}/accept`, { answerId })
  },

  getCategories: () => {
    return request.get<Category[]>('/public/categories')
  },

  getTags: () => {
    return request.get<Tag[]>('/public/tags')
  },

  createCategory: (data: { name: string; description?: string; icon?: string; sortOrder?: number }) => {
    return request.post<Category>('/admin/categories', data)
  },

  createTag: (data: { name: string; description?: string }) => {
    return request.post<Tag>('/admin/tags', data)
  }
}

export const answerApi = {
  createAnswer: (data: CreateAnswerRequest) => {
    return request.post<Answer>('/answers', data)
  },

  getAnswer: (id: number) => {
    return request.get<Answer>(`/answers/${id}`)
  },

  updateAnswer: (id: number, data: { content: string }) => {
    return request.put(`/answers/${id}`, data)
  },

  deleteAnswer: (id: number) => {
    return request.delete(`/answers/${id}`)
  },

  likeAnswer: (id: number) => {
    return request.post(`/answers/${id}/like`)
  },

  dislikeAnswer: (id: number) => {
    return request.post(`/answers/${id}/dislike`)
  }
}

export const commentApi = {
  createComment: (data: CreateCommentRequest) => {
    return request.post<Comment>('/comments', data)
  },

  deleteComment: (id: number) => {
    return request.delete(`/comments/${id}`)
  },

  likeComment: (id: number) => {
    return request.post(`/comments/${id}/like`)
  }
}
