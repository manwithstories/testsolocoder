import { request } from '@/utils/request'
import type {
  AuditRecord,
  Report,
  SensitiveWord,
  Favorite,
  Follow,
  Notification,
  Reward,
  RewardExchange,
  PointLog,
  DashboardStats,
  ActivityReport,
  AuditReport,
  PageResult,
  FavoriteRequest,
  FollowRequest
} from '@/types'

export const auditApi = {
  getAuditList: (params?: { page?: number; pageSize?: number; targetType?: string; status?: string }) => {
    return request.get<PageResult<AuditRecord>>('/admin/audit', { params })
  },

  auditContent: (data: { id: number; targetType: string; action: string; reason?: string }) => {
    return request.post('/admin/audit', data)
  },

  getPendingAuditCount: () => {
    return request.get<{ questions: number; answers: number; comments: number }>('/admin/audit/pending-count')
  },

  getReportList: (params?: { page?: number; pageSize?: number; status?: string }) => {
    return request.get<PageResult<Report>>('/admin/reports', { params })
  },

  handleReport: (id: number, data: { status: string; result?: string }) => {
    return request.put(`/admin/reports/${id}`, data)
  },

  createReport: (data: { targetType: string; targetId: number; reason: string; description?: string }) => {
    return request.post('/reports', data)
  },

  getSensitiveWords: (params?: { page?: number; pageSize?: number; category?: string; keyword?: string }) => {
    return request.get<PageResult<SensitiveWord>>('/admin/sensitive-words', { params })
  },

  createSensitiveWord: (data: { word: string; category?: string; replaceTo?: string; level?: number }) => {
    return request.post('/admin/sensitive-words', data)
  },

  deleteSensitiveWord: (id: number) => {
    return request.delete(`/admin/sensitive-words/${id}`)
  },

  checkContent: (content: string) => {
    return request.post<{ filteredContent: string; sensitiveWords: string[]; hasSensitive: boolean }>('/admin/sensitive-words/check', { content })
  }
}

export const favoriteApi = {
  addFavorite: (data: FavoriteRequest) => {
    return request.post('/favorites', data)
  },

  removeFavorite: (data: FavoriteRequest) => {
    return request.delete('/favorites', { data })
  },

  getUserFavorites: (params?: { page?: number; pageSize?: number; targetType?: string }) => {
    return request.get<PageResult<Favorite>>('/favorites', { params })
  }
}

export const followApi = {
  follow: (data: FollowRequest) => {
    return request.post('/follows', data)
  },

  unfollow: (data: FollowRequest) => {
    return request.delete('/follows', { data })
  },

  getUserFollows: (params?: { page?: number; pageSize?: number; followingType?: string }) => {
    return request.get<PageResult<Follow>>('/follows', { params })
  },

  getUserFollowers: (params?: { page?: number; pageSize?: number }) => {
    return request.get<PageResult<Follow>>('/follows/followers', { params })
  },

  isFollowing: (params: { followingType: string; followingId: number }) => {
    return request.get<{ isFollowing: boolean }>('/follows/check', { params })
  }
}

export const notificationApi = {
  getNotifications: (params?: { page?: number; pageSize?: number; isRead?: string }) => {
    return request.get<PageResult<Notification>>('/notifications', { params })
  },

  markAsRead: (id: number) => {
    return request.put(`/notifications/${id}/read`)
  },

  markAllAsRead: () => {
    return request.put('/notifications/read-all')
  },

  getUnreadCount: () => {
    return request.get<{ unreadCount: number }>('/notifications/unread-count')
  }
}

export const rewardApi = {
  getRewardList: (params?: { page?: number; pageSize?: number }) => {
    return request.get<PageResult<Reward>>('/rewards', { params })
  },

  createReward: (data: { name: string; description?: string; image?: string; pointsCost: number; stock?: number }) => {
    return request.post<Reward>('/admin/rewards', data)
  },

  updateReward: (id: number, data: Partial<Reward>) => {
    return request.put(`/admin/rewards/${id}`, data)
  },

  deleteReward: (id: number) => {
    return request.delete(`/admin/rewards/${id}`)
  },

  exchangeReward: (id: number) => {
    return request.post(`/rewards/${id}/exchange`)
  },

  getExchangeList: (params?: { page?: number; pageSize?: number }) => {
    return request.get<PageResult<RewardExchange>>('/rewards/exchanges', { params })
  },

  getPointLogs: (params?: { page?: number; pageSize?: number }) => {
    return request.get<PageResult<PointLog>>('/rewards/points/logs', { params })
  }
}

export const searchApi = {
  searchQuestions: (params?: { keyword?: string; categoryId?: number; tagId?: number; page?: number; pageSize?: number; sort?: string }) => {
    return request.get<PageResult<Question>>('/public/search/questions', { params })
  },

  getRecommendations: (params?: { limit?: number }) => {
    return request.get<Question[]>('/search/recommendations', { params })
  }
}

export const statsApi = {
  getDashboardStats: () => {
    return request.get<DashboardStats>('/admin/stats/dashboard')
  },

  getActivityReport: (params: { startDate: string; endDate: string }) => {
    return request.get<ActivityReport[]>('/admin/stats/activity', { params })
  },

  getAuditReport: (params: { startDate: string; endDate: string }) => {
    return request.get<AuditReport[]>('/admin/stats/audit', { params })
  },

  exportActivityReport: (params: { startDate: string; endDate: string }) => {
    return request.get('/admin/stats/activity/export', { params, responseType: 'blob' })
  },

  exportAuditReport: (params: { startDate: string; endDate: string }) => {
    return request.get('/admin/stats/audit/export', { params, responseType: 'blob' })
  }
}
