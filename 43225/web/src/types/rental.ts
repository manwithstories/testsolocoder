export type RentalType = 'daily' | 'hourly' | 'voyage'
export type RentalStatus = 'pending' | 'confirmed' | 'active' | 'completed' | 'cancelled' | 'refunded'
export type InsuranceType = 'basic' | 'premium' | 'none'

export interface Rental {
  id: string
  tenant_id: string
  tenant?: User
  ship_id: string
  ship?: Ship
  rental_type: RentalType
  start_date: string
  end_date: string
  start_location?: string
  end_location?: string
  base_amount: number
  insurance_type: InsuranceType
  insurance_amount: number
  platform_fee: number
  total_amount: number
  deposit_amount: number
  deposit_status: string
  currency: string
  status: RentalStatus
  emergency_contact: string
  emergency_phone: string
  notes?: string
  crew_count?: number
  passenger_count?: number
  cancellation_reason?: string
  cancelled_at?: string
  completed_at?: string
  created_at: string
  updated_at: string
}

export interface CreateRentalRequest {
  ship_id: string
  rental_type: RentalType
  start_date: string
  end_date: string
  start_location?: string
  end_location?: string
  insurance_type: InsuranceType
  emergency_contact: string
  emergency_phone: string
  notes?: string
  crew_count?: number
  passenger_count?: number
  currency?: string
}

export interface UpdateRentalStatusRequest {
  status: RentalStatus
  cancellation_reason?: string
}
