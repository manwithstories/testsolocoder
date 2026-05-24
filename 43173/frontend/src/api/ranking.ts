import request, { PageData, PageParams } from '@/utils/request'

export interface RankingItem {
  rank: number
  work_id: number
  title: string
  artist_name: string
  cover_url: string
  score: number
  work?: any
}

export const rankingApi = {
  getRanking: (params: { type?: string; category?: string; limit?: number }) => 
    request.get<RankingItem[]>('/ranking', params),
  getDailyRanking: (params: { category?: string; limit?: number }) => 
    request.get<RankingItem[]>('/ranking/daily', params),
  getWeeklyRanking: (params: { category?: string; limit?: number }) => 
    request.get<RankingItem[]>('/ranking/weekly', params),
  getMonthlyRanking: (params: { category?: string; limit?: number }) => 
    request.get<RankingItem[]>('/ranking/monthly', params),
  getHotRanking: (params: { limit?: number }) => 
    request.get<RankingItem[]>('/ranking/hot', params),
  getWorkRanking: (workId: number, params: { type?: string; category?: string }) => 
    request.get<{ rank: number; score: number }>(`/ranking/work/${workId}`, params)
}
