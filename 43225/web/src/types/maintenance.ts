export type MaintenanceType = 'routine' | 'repair' | 'inspection' | 'overhaul'
export type MaintenanceStatus = 'scheduled' | 'in_progress' | 'completed' | 'cancelled'

export interface MaintenanceRecord {
  id: string
  ship_id: string
  ship?: Ship
  maintenance_type: MaintenanceType
  title: string
  description?: string
  status: MaintenanceStatus
  planned_date: string
  start_date?: string
  completed_date?: string
  cost?: number
  currency?: string
  provider?: string
  technician?: string
  next_due_date?: string
  priority?: number
  attachments?: string
  reminder_sent: boolean
  notes?: string
  created_at: string
  updated_at: string
}

export interface MaintenanceSchedule {
  id: string
  ship_id: string
  ship?: Ship
  title: string
  description?: string
  interval_days: number
  last_completed?: string
  next_due: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface CreateMaintenanceRequest {
  ship_id: string
  maintenance_type: MaintenanceType
  title: string
  description?: string
  planned_date: string
  cost?: number
  currency?: string
  provider?: string
  technician?: string
  next_due_date?: string
  priority?: number
  notes?: string
}

export interface CreateMaintenanceScheduleRequest {
  ship_id: string
  title: string
  description?: string
  interval_days: number
  next_due: string
  is_active?: boolean
}
