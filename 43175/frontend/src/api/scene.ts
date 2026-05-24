import request from './request'

export interface Scene {
  id: number
  familyId: number
  name: string
  description: string
  icon: string
  isActive: boolean
  conditions: SceneCondition[]
  actions: SceneAction[]
  createdAt: string
  updatedAt: string
}

export interface SceneCondition {
  id: number
  sceneId: number
  type: string
  deviceId?: number
  operator: string
  value: string
  timeExpr: string
}

export interface SceneAction {
  id: number
  sceneId: number
  deviceId: number
  action: string
  value: string
}

export function listScenes(params?: { familyId?: number }): Promise<Scene[]> {
  return request.get('/scenes', { params })
}

export function createScene(data: any): Promise<Scene> {
  return request.post('/scenes', data)
}

export function getScene(id: number): Promise<Scene> {
  return request.get(`/scenes/${id}`)
}

export function updateScene(id: number, data: any): Promise<Scene> {
  return request.put(`/scenes/${id}`, data)
}

export function deleteScene(id: number): Promise<void> {
  return request.delete(`/scenes/${id}`)
}

export function executeScene(id: number): Promise<any> {
  return request.post(`/scenes/${id}/execute`)
}

export const sceneIcons = [
  { value: '🏠', label: '离家' },
  { value: '🌙', label: '睡眠' },
  { value: '☀️', label: '起床' },
  { value: '🎬', label: '影院' },
  { value: '🍽️', label: '用餐' },
  { value: '📚', label: '阅读' },
  { value: '🎉', label: '聚会' },
  { value: '🌿', label: '节能' }
]

export const conditionTypeOptions = [
  { value: 'time', label: '时间触发' },
  { value: 'device', label: '设备状态' },
  { value: 'sensor', label: '传感器数据' }
]

export const actionOptions = [
  { value: 'on', label: '开启' },
  { value: 'off', label: '关闭' },
  { value: 'toggle', label: '切换' },
  { value: 'dim', label: '调光' },
  { value: 'set_temp', label: '设置温度' }
]
