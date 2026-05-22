import { request } from './request'

export const checkIn = (qrCode: string) => {
  return request({
    url: '/checkins',
    method: 'post',
    data: { qrCode }
  })
}

export const getCheckInList = (params: any) => {
  return request({
    url: '/checkins',
    method: 'get',
    params
  })
}

export const getCheckInStatistics = (activityId?: number) => {
  const params = activityId ? { activityId } : {}
  return request({
    url: '/checkins/statistics',
    method: 'get',
    params
  })
}
