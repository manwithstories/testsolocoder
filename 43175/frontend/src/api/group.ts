import request from './request'

export interface DeviceGroup {
  id: number
  familyId: number
  name: string
  description: string
  type: string
  devices?: any[]
  createdAt: string
  updatedAt: string
}

export function listGroups(params?: { familyId?: number }): Promise<DeviceGroup[]> {
  return request.get('/groups', { params })
}

export function createGroup(data: { familyId: number; name: string; description?: string; type: string }): Promise<DeviceGroup> {
  return request.post('/groups', data)
}

export function getGroup(id: number): Promise<DeviceGroup> {
  return request.get(`/groups/${id}`)
}

export function updateGroup(id: number, data: { name?: string; description?: string; type?: string }): Promise<DeviceGroup> {
  return request.put(`/groups/${id}`, data)
}

export function deleteGroup(id: number): Promise<void> {
  return request.delete(`/groups/${id}`)
}

export function addDeviceToGroup(groupId: number, deviceId: number): Promise<void> {
  return request.post(`/groups/${groupId}/devices`, { deviceId })
}

export function removeDeviceFromGroup(groupId: number, deviceId: number): Promise<void> {
  return request.delete(`/groups/${groupId}/devices/${deviceId}`)
}

export function batchControlGroup(groupId: number, action: string): Promise<any> {
  return request.put(`/groups/${groupId}/control`, { action })
}

export function getGroupEnergy(id: number, period?: string): Promise<any> {
  return request.get(`/groups/${id}/energy`, { params: { period } })
}

export const groupTypeOptions = [
  { value: 'room', label: '房间' },
  { value: 'floor', label: '楼层' },
  { value: 'function', label: '功能' },
  { value: 'custom', label: '自定义' }
]
