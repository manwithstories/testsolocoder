import request from './request'

export interface EnergyData {
  id: number
  deviceId: number
  familyId: number
  power: number
  voltage: number
  current: number
  energyUsed: number
  timestamp: string
  date: string
  hour: number
}

export interface EnergyAlert {
  id: number
  familyId: number
  deviceId: number
  alertType: string
  level: string
  message: string
  value: number
  threshold: number
  resolved: boolean
  createdAt: string
}

export function getRealtimeEnergy(params?: { familyId?: number }): Promise<any> {
  return request.get('/energy/realtime', { params })
}

export function getEnergyStatistics(params?: { familyId?: number; period?: string }): Promise<any> {
  return request.get('/energy/statistics', { params })
}

export function getEnergyTrend(params?: { familyId?: number; days?: number }): Promise<any> {
  return request.get('/energy/trend', { params })
}

export function listEnergyAlerts(params?: { familyId?: number; resolved?: string }): Promise<EnergyAlert[]> {
  return request.get('/energy/alerts', { params })
}

export function exportEnergyReport(params?: { familyId?: number; format?: string; period?: string }): Promise<any> {
  return request.get('/energy/export', {
    params,
    responseType: 'blob'
  })
}
