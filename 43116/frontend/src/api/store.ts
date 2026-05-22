import request from './request'
import type { City, Store, PageResult, ApiResponse } from '@/types'

export const storeApi = {
  createCity: (data: Record<string, any>) =>
    request.post<any, ApiResponse<City>>('/cities', data),

  getCities: () =>
    request.get<any, ApiResponse<City[]>>('/cities'),

  getCityById: (id: number) =>
    request.get<any, ApiResponse<City>>(`/cities/${id}`),

  updateCity: (id: number, data: Record<string, any>) =>
    request.put<any, ApiResponse<null>>(`/cities/${id}`, data),

  deleteCity: (id: number) =>
    request.delete<any, ApiResponse<null>>(`/cities/${id}`),

  getStoresByCity: (cityId: number) =>
    request.get<any, ApiResponse<Store[]>>(`/cities/${cityId}/stores`),

  createStore: (data: Record<string, any>) =>
    request.post<any, ApiResponse<Store>>('/stores', data),

  getStores: (params?: { page?: number; page_size?: number; city_id?: number; keyword?: string }) =>
    request.get<any, ApiResponse<PageResult<Store>>>('/stores', { params }),

  getStoreById: (id: number) =>
    request.get<any, ApiResponse<Store>>(`/stores/${id}`),

  updateStore: (id: number, data: Record<string, any>) =>
    request.put<any, ApiResponse<null>>(`/stores/${id}`, data),

  deleteStore: (id: number) =>
    request.delete<any, ApiResponse<null>>(`/stores/${id}`)
}
