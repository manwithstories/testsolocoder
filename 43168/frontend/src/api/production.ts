import request from '@/utils/request'
import type { PaginationParams, PaginationResult } from '@/types'

export interface ProductionProgress {
  id: number
  orderId?: number
  orderNo?: string
  productName?: string
  quantity: number
  currentStage: number
  stageName?: string
  progress: number
  estimatedDate?: string
  completedDate?: string
  remark?: string
  createdAt: string
  updatedAt: string
}

export interface ListProductionParams extends PaginationParams {
  orderNo?: string
  currentStage?: number
  keyword?: string
}

export const listProduction = (params: ListProductionParams) => {
  return request.get<any, PaginationResult<ProductionProgress>>('/productions', { params })
}

export const getProduction = (id: number | string) => {
  return request.get<any, ProductionProgress>(`/productions/${id}`)
}

export const updateProductionStage = (id: number | string, stage: number, remark?: string) => {
  return request.patch<any, void>(`/productions/${id}/stage`, { stage, remark })
}
