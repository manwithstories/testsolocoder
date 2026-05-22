import { request } from './request'

export interface ActivityData {
  title: string
  description: string
  startTime: string
  endTime: string
  location: string
  capacity: number
  poster?: string
}

export const createActivity = (data: ActivityData) => {
  return request({
    url: '/activities',
    method: 'post',
    data
  })
}

export const getActivityList = (params: any) => {
  return request({
    url: '/activities',
    method: 'get',
    params
  })
}

export const getActivity = (id: number) => {
  return request({
    url: `/activities/${id}`,
    method: 'get'
  })
}

export const updateActivity = (id: number, data: Partial<ActivityData>) => {
  return request({
    url: `/activities/${id}`,
    method: 'put',
    data
  })
}

export const updateActivityStatus = (id: number, status: string) => {
  return request({
    url: `/activities/${id}/status`,
    method: 'put',
    data: { status }
  })
}

export const deleteActivity = (id: number) => {
  return request({
    url: `/activities/${id}`,
    method: 'delete'
  })
}

export const uploadImage = (file: File) => {
  const formData = new FormData()
  formData.append('file', file)
  return request({
    url: '/upload/image',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}
