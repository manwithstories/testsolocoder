export type ShipType = 'sailboat' | 'motorboat' | 'yacht' | 'fishing' | 'cargo'
export type ShipStatus = 'available' | 'rented' | 'maintenance' | 'inactive'

export interface ShipImage {
  id: string
  ship_id: string
  url: string
  is_primary: boolean
  sort_order: number
  created_at: string
}

export interface Ship {
  id: string
  owner_id: string
  owner?: User
  name: string
  description?: string
  ship_type: ShipType
  capacity: number
  cabin_count?: number
  bathroom_count?: number
  length?: number
  width?: number
  year_built?: number
  equipment?: string
  features?: string
  sailing_area?: string
  home_port?: string
  license_number?: string
  hourly_rate: number
  daily_rate: number
  deposit_amount?: number
  insurance_required: boolean
  cancellation_policy?: string
  status: ShipStatus
  average_rating: number
  review_count: number
  view_count: number
  created_at: string
  updated_at: string
  images?: ShipImage[]
}

export interface CreateShipRequest {
  name: string
  description?: string
  ship_type: ShipType
  capacity: number
  cabin_count?: number
  bathroom_count?: number
  length?: number
  width?: number
  year_built?: number
  equipment?: string
  features?: string
  sailing_area?: string
  home_port?: string
  license_number?: string
  hourly_rate: number
  daily_rate: number
  deposit_amount?: number
  insurance_required?: boolean
  cancellation_policy?: string
}

export interface SearchShipRequest {
  ship_type?: ShipType
  min_capacity?: number
  max_capacity?: number
  min_price?: number
  max_price?: number
  location?: string
  start_date?: string
  end_date?: string
  page?: number
  page_size?: number
  sort_by?: string
  sort_order?: string
}
