export type BerthType = 'small' | 'medium' | 'large'
export type BerthStatus = 'available' | 'occupied' | 'reserved' | 'maintenance'

export interface Berth {
  id: string
  dock_id: string
  number: string
  berth_type: BerthType
  max_length?: number
  max_width?: number
  hourly_rate: number
  daily_rate: number
  has_water: boolean
  has_electric: boolean
  has_internet: boolean
  description?: string
  status: BerthStatus
  created_at: string
  updated_at: string
}

export interface Dock {
  id: string
  name: string
  address: string
  city?: string
  country?: string
  latitude?: number
  longitude?: number
  description?: string
  image_url?: string
  amenities?: string
  open_time?: string
  close_time?: string
  is_active: boolean
  average_rating: number
  review_count: number
  created_at: string
  updated_at: string
  berths?: Berth[]
}

export interface BerthReservation {
  id: string
  berth_id: string
  berth?: Berth
  ship_id?: string
  ship?: Ship
  rental_id?: string
  user_id: string
  user?: User
  start_time: string
  end_time: string
  total_amount: number
  status: string
  notes?: string
  created_at: string
  updated_at: string
}

export interface CreateBerthRequest {
  dock_id: string
  number: string
  berth_type: BerthType
  max_length?: number
  max_width?: number
  hourly_rate: number
  daily_rate: number
  has_water?: boolean
  has_electric?: boolean
  has_internet?: boolean
  description?: string
}

export interface CreateReservationRequest {
  berth_id: string
  ship_id?: string
  start_time: string
  end_time: string
  notes?: string
}
