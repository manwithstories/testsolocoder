import { request } from './request'

export interface TicketTypeData {
  activityId: number
  name: string
  type: string
  price: number
  stock: number
}

export const createTicketType = (data: TicketTypeData) => {
  return request({
    url: '/ticket-types',
    method: 'post',
    data
  })
}

export const getTicketTypeList = (params: any) => {
  return request({
    url: '/ticket-types',
    method: 'get',
    params
  })
}

export const getTicketType = (id: number) => {
  return request({
    url: `/ticket-types/${id}`,
    method: 'get'
  })
}

export const updateTicketType = (id: number, data: Partial<TicketTypeData> & { status?: string }) => {
  return request({
    url: `/ticket-types/${id}`,
    method: 'put',
    data
  })
}

export const deleteTicketType = (id: number) => {
  return request({
    url: `/ticket-types/${id}`,
    method: 'delete'
  })
}
