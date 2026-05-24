import { apiGet, apiPost, apiPut, apiDelete, apiUpload } from './client'
import { Pet, CreatePetRequest, PetListQuery, ApiResponse } from '../types'

export const listPets = (params?: PetListQuery): Promise<ApiResponse<any>> => {
  return apiGet('/pets', params)
}

export const getPet = (id: number): Promise<ApiResponse<Pet>> => {
  return apiGet<Pet>(`/pets/${id}`)
}

export const createPet = (data: CreatePetRequest): Promise<ApiResponse<Pet>> => {
  return apiPost<Pet>('/pets', data)
}

export const updatePet = (id: number, data: any): Promise<ApiResponse<Pet>> => {
  return apiPut<Pet>(`/pets/${id}`, data)
}

export const deletePet = (id: number): Promise<ApiResponse<any>> => {
  return apiDelete(`/pets/${id}`)
}

export const updatePetStatus = (id: number, status: string): Promise<ApiResponse<any>> => {
  return apiPut(`/pets/${id}/status`, { status })
}

export const uploadPetPhotos = (id: number, formData: FormData): Promise<ApiResponse<any>> => {
  return apiUpload(`/pets/${id}/photos`, formData)
}

export const uploadPetVideos = (id: number, formData: FormData): Promise<ApiResponse<any>> => {
  return apiUpload(`/pets/${id}/videos`, formData)
}

export const getPetAdoptionHistory = (id: number): Promise<ApiResponse<any>> => {
  return apiGet(`/pets/${id}/history`)
}

export const getMyPets = (params?: any): Promise<ApiResponse<any>> => {
  return apiGet('/pets/my', params)
}

export const getMyAdoptedPets = (): Promise<ApiResponse<Pet[]>> => {
  return apiGet<Pet[]>('/pets/adopted')
}
