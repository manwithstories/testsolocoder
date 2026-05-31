import request from '@/utils/request'
import type {
  Exhibition,
  Collection,
  TimeSlot,
  ApiResponse,
  PageResult
} from '@/types'

export interface ExhibitionQuery {
  page?: number
  page_size?: number
  keyword?: string
  status?: string
  start_date?: string
  end_date?: string
}

export const listExhibitions = (params?: ExhibitionQuery) => {
  return request.get<any, ApiResponse<PageResult<Exhibition>>>('/exhibitions', { params })
}

export const getHotExhibitions = () => {
  return request.get<any, ApiResponse<Exhibition[]>>('/exhibitions/hot')
}

export const getExhibition = (id: number) => {
  return request.get<any, ApiResponse<Exhibition>>(`/exhibitions/${id}`)
}

export const createExhibition = (data: Partial<Exhibition> & { collection_ids?: number[] }) => {
  return request.post<any, ApiResponse<Exhibition>>('/exhibitions', data)
}

export const updateExhibition = (id: number, data: Partial<Exhibition> & { collection_ids?: number[] }) => {
  return request.put<any, ApiResponse<void>>(`/exhibitions/${id}`, data)
}

export const deleteExhibition = (id: number) => {
  return request.delete<any, ApiResponse<void>>(`/exhibitions/${id}`)
}

export const addExhibitionCollections = (id: number, collection_ids: number[]) => {
  return request.post<any, ApiResponse<void>>(`/exhibitions/${id}/collections`, { collection_ids })
}

export const removeExhibitionCollections = (id: number, collection_ids: number[]) => {
  return request.delete<any, ApiResponse<void>>(`/exhibitions/${id}/collections`, { data: { collection_ids } })
}

export const getExhibitionCollections = (id: number) => {
  return request.get<any, ApiResponse<Collection[]>>(`/exhibitions/${id}/collections`)
}

export const createTimeSlot = (exhibitionId: number, data: Partial<TimeSlot>) => {
  return request.post<any, ApiResponse<TimeSlot>>(`/exhibitions/${exhibitionId}/time-slots`, data)
}

export const batchCreateTimeSlots = (data: {
  exhibition_id: number
  start_date: string
  end_date: string
  start_time: string
  end_time: string
  interval: number
  max_capacity: number
}) => {
  return request.post<any, ApiResponse<void>>('/exhibitions/time-slots/batch', data)
}

export const listTimeSlots = (exhibitionId: number, date?: string) => {
  return request.get<any, ApiResponse<TimeSlot[]>>(`/exhibitions/${exhibitionId}/time-slots`, {
    params: { date }
  })
}
