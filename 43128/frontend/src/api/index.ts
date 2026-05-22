import { get, post, put } from './request'
import type { LoginResp, Paged, User } from '@/types'

export const userApi = {
  register: (data: any) => post<any>('/auth/register', data),
  login: (data: any) => post<LoginResp>('/auth/login', data),
  profile: () => get<User>('/users/me'),
  verify: (data: any) => put<any>('/users/me/verify', data),
  list: (params?: any) => get<Paged<User>>('/admin/users', params),
}

export const eventApi = {
  list: (params?: any) => get<Paged<any>>('/events', params),
  listAll: (params?: any) => get<Paged<any>>('/admin/events', { ...params, all: 1 }),
  get: (id: number) => get<any>(`/events/${id}`),
  create: (data: any) => post<any>('/admin/events', data),
  update: (id: number, data: any) => put<any>(`/admin/events/${id}`, data),
  publish: (id: number) => put<any>(`/admin/events/${id}/publish`),
  unpublish: (id: number) => put<any>(`/admin/events/${id}/unpublish`),
}

export const regApi = {
  my: (params?: any) => get<Paged<any>>('/registrations/me', params),
  listByEvent: (eventId: number, params?: any) =>
    get<Paged<any>>(`/admin/registrations/event/${eventId}`, params),
  create: (data: any) => post<any>('/registrations', data),
  confirm: (id: number) => put<any>(`/admin/registrations/${id}/confirm`),
}

export const scoreApi = {
  my: (params?: any) => get<Paged<any>>('/scores/me', params),
  listByItem: (itemId: number) => get<any[]>(`/admin/scores/item/${itemId}`),
  entry: (data: any) => post<any>('/admin/scores', data),
  import: (form: FormData) => post<any>('/admin/scores/import', form),
  template: () => `${location.origin}/api/v1/admin/scores/import/template`,
}

export const certApi = {
  my: () => get<any[]>('/certificates/me'),
  generate: (scoreId: number) => post<any>(`/certificates/${scoreId}/generate`),
  download: (id: number) => `${location.origin}/api/v1/certificates/${id}/download`,
}

export const msgApi = {
  list: (params?: any) => get<Paged<any>>('/messages', params),
  unreadCount: () => get<{ count: number }>('/messages/unread-count'),
  markRead: (id: number) => put<any>(`/messages/${id}/read`),
  markAll: () => put<any>('/messages/read-all'),
}

export const statsApi = {
  overview: (params?: any) => get<any>('/stats/overview', params),
  export: () => `${location.origin}/api/v1/admin/stats/export`,
}
