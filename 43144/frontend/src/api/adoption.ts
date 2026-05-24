import { apiGet, apiPost, apiPut } from './client'
import {
  AdoptionApplication,
  AdoptionAgreement,
  FollowUpRecord,
  CreateAdoptionRequest,
  ReviewAdoptionRequest,
  CreateFollowUpRequest,
  AdoptionListQuery,
  ApiResponse,
} from '../types'

export const listAdoptionApplications = (params?: AdoptionListQuery): Promise<ApiResponse<any>> => {
  return apiGet('/adoption/applications', params)
}

export const getAdoptionApplication = (id: number): Promise<ApiResponse<AdoptionApplication>> => {
  return apiGet<AdoptionApplication>(`/adoption/applications/${id}`)
}

export const createAdoptionApplication = (data: CreateAdoptionRequest): Promise<ApiResponse<AdoptionApplication>> => {
  return apiPost<AdoptionApplication>('/adoption/applications', data)
}

export const reviewAdoptionApplication = (id: number, data: ReviewAdoptionRequest): Promise<ApiResponse<AdoptionApplication>> => {
  return apiPut<AdoptionApplication>(`/adoption/applications/${id}/review`, data)
}

export const signAdoptionAgreement = (id: number): Promise<ApiResponse<AdoptionAgreement>> => {
  return apiPut<AdoptionAgreement>(`/adoption/applications/${id}/sign`)
}

export const getAdoptionAgreement = (id: number): Promise<ApiResponse<AdoptionAgreement>> => {
  return apiGet<AdoptionAgreement>(`/adoption/applications/${id}/agreement`)
}

export const completeAdoption = (id: number): Promise<ApiResponse<AdoptionApplication>> => {
  return apiPut<AdoptionApplication>(`/adoption/applications/${id}/complete`)
}

export const createFollowUpRecord = (data: CreateFollowUpRequest): Promise<ApiResponse<FollowUpRecord>> => {
  return apiPost<FollowUpRecord>('/adoption/follow-ups', data)
}

export const listFollowUpRecords = (petId: number): Promise<ApiResponse<FollowUpRecord[]>> => {
  return apiGet<FollowUpRecord[]>(`/adoption/pets/${petId}/follow-ups`)
}
