import { User } from './user'
import { Pet } from './pet'

export interface Appointment {
  id: number
  user_id: number
  user?: User
  pet_id: number
  pet?: Pet
  rescue_id: number
  appointment_type: 'visit' | 'checkup'
  appointment_date: string
  start_time: string
  end_time: string
  status: 'pending' | 'confirmed' | 'cancelled' | 'completed' | 'rescheduled'
  location?: string
  notes?: string
  cancel_reason?: string
  original_id?: number
  created_at: string
  updated_at: string
}

export interface CreateAppointmentRequest {
  pet_id: number
  appointment_type: string
  appointment_date: string
  start_time: string
  end_time: string
  location?: string
  notes?: string
}

export interface UpdateAppointmentRequest {
  appointment_date?: string
  start_time?: string
  end_time?: string
  location?: string
  notes?: string
}

export interface AppointmentListQuery {
  page?: number
  page_size?: number
  status?: string
  pet_id?: number
  rescue_id?: number
  date_from?: string
  date_to?: string
}
