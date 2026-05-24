import request from '@/utils/request'
import type { PaginationParams, PaginationResult } from '@/types'

export interface Delivery {
  id: number
  orderId?: number
  orderNo?: string
  ownerName?: string
  address: string
  contactName: string
  contactPhone: string
  deliveryDate?: string
  timeSlot?: string
  installer?: string
  status: number
  remark?: string
  createdAt: string
  updatedAt: string
}

export interface ListDeliveryParams extends PaginationParams {
  orderNo?: string
  status?: number
  keyword?: string
  contactPhone?: string
}

export interface DeliveryFormData {
  id?: number
  orderId?: number
  address: string
  contactName: string
  contactPhone: string
  deliveryDate?: string
  timeSlot?: string
  installer?: string
  remark?: string
  status?: number
}

export interface AddressBook {
  id: number
  name: string
  contactName: string
  contactPhone: string
  address: string
  isDefault?: boolean
}

export const listDeliveries = (params: ListDeliveryParams) => {
  return request.get<any, PaginationResult<Delivery>>('/deliveries', { params })
}

export const getDelivery = (id: number | string) => {
  return request.get<any, Delivery>(`/deliveries/${id}`)
}

export const createDelivery = (data: DeliveryFormData) => {
  return request.post<any, Delivery>('/deliveries', data)
}

export const updateDelivery = (id: number | string, data: DeliveryFormData) => {
  return request.put<any, Delivery>(`/deliveries/${id}`, data)
}

export const deleteDelivery = (id: number | string) => {
  return request.delete<any, void>(`/deliveries/${id}`)
}

export const updateDeliveryStatus = (id: number | string, status: number) => {
  return request.patch<any, void>(`/deliveries/${id}/status`, { status })
}

export const listAddressBook = () => {
  return request.get<any, AddressBook[]>('/deliveries/address-book')
}
