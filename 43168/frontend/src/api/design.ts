import request from '@/utils/request'
import type { PaginationParams, PaginationResult } from '@/types'

export interface DesignProject {
  id: number
  name: string
  description?: string
  ownerName?: string
  ownerId?: number
  designerId?: number
  designerName?: string
  houseType?: string
  area?: number
  budget?: number
  coverImage?: string
  images?: DesignImage[]
  annotations?: DesignAnnotation[]
  status: number
  createdAt: string
  updatedAt: string
}

export interface DesignImage {
  id: number
  url: string
  name?: string
  sort?: number
}

export interface DesignAnnotation {
  id: number
  imageId?: number
  x?: number
  y?: number
  content: string
  author?: string
  createdAt: string
}

export interface ListDesignParams extends PaginationParams {
  name?: string
  status?: number
  keyword?: string
}

export interface DesignFormData {
  id?: number
  name: string
  description?: string
  ownerId?: number
  houseType?: string
  area?: number
  budget?: number
  coverImage?: string
  images?: string[]
  status?: number
}

export const listDesigns = (params: ListDesignParams) => {
  return request.get<any, PaginationResult<DesignProject>>('/designs', { params })
}

export const getDesign = (id: number | string) => {
  return request.get<any, DesignProject>(`/designs/${id}`)
}

export const createDesign = (data: DesignFormData) => {
  return request.post<any, DesignProject>('/designs', data)
}

export const updateDesign = (id: number | string, data: DesignFormData) => {
  return request.put<any, DesignProject>(`/designs/${id}`, data)
}

export const deleteDesign = (id: number | string) => {
  return request.delete<any, void>(`/designs/${id}`)
}

export const updateDesignStatus = (id: number | string, status: number) => {
  return request.patch<any, void>(`/designs/${id}/status`, { status })
}

export const addAnnotation = (id: number | string, data: { imageId?: number; x?: number; y?: number; content: string }) => {
  return request.post<any, DesignAnnotation>(`/designs/${id}/annotations`, data)
}

export const listAnnotations = (id: number | string) => {
  return request.get<any, DesignAnnotation[]>(`/designs/${id}/annotations`)
}
