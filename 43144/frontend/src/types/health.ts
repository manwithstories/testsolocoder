import { Pet } from './pet'

export interface HealthRecord {
  id: number
  pet_id: number
  pet?: Pet
  record_type: 'vaccine' | 'deworm' | 'checkup' | 'disease' | 'surgery' | 'other'
  title: string
  description: string
  vaccine_name?: string
  record_date: string
  next_date?: string
  weight?: number
  temperature?: number
  vet_name?: string
  hospital?: string
  report_file?: string
  notes?: string
  recorded_by: number
  rescue_id: number
  created_at: string
  updated_at: string
}

export interface CreateHealthRecordRequest {
  pet_id: number
  record_type: string
  title: string
  description?: string
  vaccine_name?: string
  record_date: string
  next_date?: string
  weight?: number
  temperature?: number
  vet_name?: string
  hospital?: string
  notes?: string
}

export interface HealthRecordListQuery {
  page?: number
  page_size?: number
  pet_id?: number
  record_type?: string
}

export interface HealthReminder {
  id: number
  pet_id: number
  record_id?: number
  title: string
  reminder_date: string
  is_completed: boolean
  notes?: string
  created_at: string
  updated_at: string
}
