import request from './request'

export interface Device {
  id: number
  familyId: number
  name: string
  deviceType: string
  vendor: string
  location: string
  power: number
  protocol: string
  status: string
  lastOnlineTime?: string
  createdAt: string
  updatedAt: string
}

export interface CreateDeviceRequest {
  familyId: number
  name: string
  deviceType: string
  vendor?: string
  location?: string
  power: number
  protocol: string
}

export interface UpdateDeviceRequest {
  name?: string
  deviceType?: string
  vendor?: string
  location?: string
  power?: number
  protocol?: string
  status?: string
}

export function listDevices(params?: { familyId?: number; deviceType?: string; status?: string; location?: string }): Promise<Device[]> {
  return request.get('/devices', { params })
}

export function createDevice(data: CreateDeviceRequest): Promise<Device> {
  return request.post('/devices', data)
}

export function getDevice(id: number): Promise<Device> {
  return request.get(`/devices/${id}`)
}

export function updateDevice(id: number, data: UpdateDeviceRequest): Promise<Device> {
  return request.put(`/devices/${id}`, data)
}

export function deleteDevice(id: number): Promise<void> {
  return request.delete(`/devices/${id}`)
}

export function updateDeviceStatus(id: number, status: string): Promise<any> {
  return request.put(`/devices/${id}/status`, { status })
}

export function getDeviceEnergy(id: number, period?: string): Promise<any> {
  return request.get(`/devices/${id}/energy`, { params: { period } })
}

export const deviceTypeOptions = [
  { value: 'light', label: '照明灯具' },
  { value: 'ac', label: '空调' },
  { value: 'heater', label: '取暖器' },
  { value: 'curtain', label: '窗帘' },
  { value: 'camera', label: '摄像头' },
  { value: 'sensor', label: '传感器' },
  { value: 'switch', label: '开关' },
  { value: 'plug', label: '智能插座' },
  { value: 'other', label: '其他' }
]

export const protocolOptions = [
  { value: 'wifi', label: 'Wi-Fi' },
  { value: 'zigbee', label: 'Zigbee' },
  { value: 'bluetooth', label: '蓝牙' },
  { value: 'zwave', label: 'Z-Wave' },
  { value: 'mqtt', label: 'MQTT' },
  { value: 'other', label: '其他' }
]

export const deviceStatusMap: Record<string, { label: string; color: string }> = {
  online: { label: '在线', color: 'success' },
  offline: { label: '离线', color: 'info' },
  on: { label: '开启', color: 'success' },
  off: { label: '关闭', color: 'warning' }
}
