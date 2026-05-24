import api from './index'
import type { ApiResponse, PaginatedResponse, User, Equipment, SearchRequest } from '@/types'

export const userApi = {
  getProfile: () => {
    return api.get<any, ApiResponse<User>>('/users/profile')
  },

  updateProfile: (data: Partial<User>) => {
    return api.put<any, ApiResponse<User>>('/users/profile', data)
  },

  getAllUsers: () => {
    return api.get<any, ApiResponse<User[]>>('/users/')
  },

  verifyUser: (userId: number, verified: boolean) => {
    return api.put<any, ApiResponse<null>>('/users/verify', { userId, verified })
  }
}

export const equipmentApi = {
  getEquipment: (id: number) => {
    return api.get<any, ApiResponse<Equipment>>(`/equipments/${id}`)
  },

  getMyEquipments: () => {
    return api.get<any, ApiResponse<Equipment[]>>('/equipments/my')
  },

  createEquipment: (data: any) => {
    return api.post<any, ApiResponse<Equipment>>('/equipments/', data)
  },

  updateEquipment: (id: number, data: any) => {
    return api.put<any, ApiResponse<Equipment>>(`/equipments/${id}`, data)
  },

  deleteEquipment: (id: number) => {
    return api.delete<any, ApiResponse<null>>(`/equipments/${id}`)
  },

  uploadImage: (equipmentId: number, file: File) => {
    const formData = new FormData()
    formData.append('image', file)
    return api.post<any, ApiResponse<any>>(`/equipments/${equipmentId}/images`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
  },

  deleteImage: (equipmentId: number, imageId: number) => {
    return api.delete<any, ApiResponse<null>>(`/equipments/${equipmentId}/images/${imageId}`)
  },

  getReservedDates: (id: number) => {
    return api.get<any, ApiResponse<string[]>>(`/equipments/${id}/reserved-dates`)
  },

  getCategories: () => {
    return api.get<any, ApiResponse<string[]>>('/equipments/categories')
  },

  getBrands: () => {
    return api.get<any, ApiResponse<string[]>>('/equipments/brands')
  }
}

export const searchApi = {
  searchEquipments: (params: SearchRequest) => {
    return api.post<any, PaginatedResponse<Equipment>>('/search/equipments', params)
  },

  searchOrders: (params: SearchRequest) => {
    return api.post<any, PaginatedResponse<any>>('/search/orders', params)
  }
}
