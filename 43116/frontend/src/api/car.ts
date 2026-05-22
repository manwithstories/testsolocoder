import request from './request'
import type { Car, CarImage, PageResult, ApiResponse } from '@/types'

export const carApi = {
  createCar: (data: Record<string, any>) =>
    request.post<any, ApiResponse<Car>>('/cars', data),

  getCars: (params?: {
    page?: number; page_size?: number; keyword?: string;
    status?: string; brand?: string; store_id?: number
  }) =>
    request.get<any, ApiResponse<PageResult<Car>>>('/cars', { params }),

  getAvailableCars: (params?: { page?: number; page_size?: number; store_id?: number }) =>
    request.get<any, ApiResponse<PageResult<Car>>>('/cars/available', { params }),

  getCarById: (id: number) =>
    request.get<any, ApiResponse<Car>>(`/cars/${id}`),

  updateCar: (id: number, data: Record<string, any>) =>
    request.put<any, ApiResponse<null>>(`/cars/${id}`, data),

  updateCarStatus: (id: number, status: string) =>
    request.put<any, ApiResponse<null>>(`/cars/${id}/status`, { status }),

  deleteCar: (id: number) =>
    request.delete<any, ApiResponse<null>>(`/cars/${id}`),

  uploadCarImage: (id: number, file: File) => {
    const formData = new FormData()
    formData.append('file', file)
    return request.post<any, ApiResponse<CarImage>>(`/cars/${id}/upload`, formData)
  },

  batchUploadCarImages: (id: number, files: File[]) => {
    const formData = new FormData()
    files.forEach(file => formData.append('files', file))
    return request.post<any, ApiResponse<CarImage[]>>(`/cars/${id}/batch-upload`, formData)
  },

  deleteCarImage: (id: number, imageId: number) =>
    request.delete<any, ApiResponse<null>>(`/cars/${id}/images/${imageId}`),

  getCarImages: (id: number) =>
    request.get<any, ApiResponse<CarImage[]>>(`/cars/${id}/images`)
}
