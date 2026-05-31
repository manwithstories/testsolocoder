import request from '@/utils/request'
import type {
  Collection,
  CollectionCategory,
  CollectionTag,
  ApiResponse,
  PageResult
} from '@/types'

export interface CollectionQuery {
  page?: number
  page_size?: number
  keyword?: string
  category_id?: number
  status?: string
  era?: string
  material?: string
  sort_by?: string
  sort_order?: string
}

export const listCollections = (params?: CollectionQuery) => {
  return request.get<any, ApiResponse<PageResult<Collection>>>('/collections', { params })
}

export const getCollection = (id: number) => {
  return request.get<any, ApiResponse<Collection>>(`/collections/${id}`)
}

export const createCollection = (data: Partial<Collection>) => {
  return request.post<any, ApiResponse<Collection>>('/collections', data)
}

export const updateCollection = (id: number, data: Partial<Collection>) => {
  return request.put<any, ApiResponse<void>>(`/collections/${id}`, data)
}

export const deleteCollection = (id: number) => {
  return request.delete<any, ApiResponse<void>>(`/collections/${id}`)
}

export const batchImportCollections = (data: Partial<Collection>[]) => {
  return request.post<any, ApiResponse<{ imported_count: number }>>('/collections/batch-import', data)
}

export const listCategories = () => {
  return request.get<any, ApiResponse<CollectionCategory[]>>('/collections/categories')
}

export const createCategory = (data: { name: string; parent_id?: number; sort_order?: number }) => {
  return request.post<any, ApiResponse<CollectionCategory>>('/collections/categories', data)
}

export const updateCategory = (id: number, data: { name: string; parent_id?: number; sort_order?: number }) => {
  return request.put<any, ApiResponse<void>>(`/collections/categories/${id}`, data)
}

export const deleteCategory = (id: number) => {
  return request.delete<any, ApiResponse<void>>(`/collections/categories/${id}`)
}

export const listTags = () => {
  return request.get<any, ApiResponse<CollectionTag[]>>('/collections/tags')
}

export const createTag = (name: string) => {
  return request.post<any, ApiResponse<CollectionTag>>('/collections/tags', { name })
}
