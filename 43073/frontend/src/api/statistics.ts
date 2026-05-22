import { request } from './request'

export const getActivityStatistics = (params: any) => {
  return request({
    url: '/statistics/activities',
    method: 'get',
    params
  })
}

export const getTicketTypeStatistics = (params: any) => {
  return request({
    url: '/statistics/ticket-types',
    method: 'get',
    params
  })
}

export const getDailyStatistics = (params: any) => {
  return request({
    url: '/statistics/daily',
    method: 'get',
    params
  })
}

export const exportStatistics = (params: any) => {
  return request({
    url: '/statistics/export',
    method: 'get',
    params,
    responseType: 'blob'
  })
}
