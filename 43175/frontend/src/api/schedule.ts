import request from './request'

export interface Schedule {
  id: number
  familyId: number
  name: string
  description: string
  deviceId: number
  action: string
  value: string
  cronExpr: string
  isEnabled: boolean
  lastRun?: string
  createdAt: string
  updatedAt: string
}

export interface ScheduleLog {
  id: number
  scheduleId: number
  deviceId: number
  action: string
  value: string
  success: boolean
  message: string
  energyDelta: number
  executedAt: string
}

export function listSchedules(params?: { familyId?: number; deviceId?: number }): Promise<Schedule[]> {
  return request.get('/schedules', { params })
}

export function createSchedule(data: any): Promise<Schedule> {
  return request.post('/schedules', data)
}

export function getSchedule(id: number): Promise<Schedule> {
  return request.get(`/schedules/${id}`)
}

export function updateSchedule(id: number, data: any): Promise<Schedule> {
  return request.put(`/schedules/${id}`, data)
}

export function deleteSchedule(id: number): Promise<void> {
  return request.delete(`/schedules/${id}`)
}

export function listScheduleLogs(id: number): Promise<ScheduleLog[]> {
  return request.get(`/schedules/${id}/logs`)
}

export const scheduleActionOptions = [
  { value: 'on', label: '开启' },
  { value: 'off', label: '关闭' },
  { value: 'toggle', label: '切换' },
  { value: 'dim', label: '调光' },
  { value: 'set_temp', label: '设置温度' }
]

export const cronPresets = [
  { label: '每天 8:00', value: '0 8 * * *' },
  { label: '每天 22:00', value: '0 22 * * *' },
  { label: '工作日 18:00', value: '0 18 * * 1-5' },
  { label: '周末 9:00', value: '0 9 * * 0,6' },
  { label: '每小时', value: '0 * * * *' },
  { label: '每30分钟', value: '*/30 * * * *' }
]
