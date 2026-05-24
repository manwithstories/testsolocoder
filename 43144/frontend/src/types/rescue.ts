export interface RescueStation {
  id: number
  name: string
  address?: string
  contact_person?: string
  contact_phone?: string
  contact_email?: string
  license_number?: string
  license_file?: string
  description?: string
  status: 'pending' | 'approved' | 'rejected'
  verified_by?: number
  verified_at?: string
  reject_reason?: string
  created_at: string
  updated_at: string
}

export interface ReviewRescueRequest {
  action: 'approve' | 'reject'
  reject_reason?: string
}

export interface RescueStats {
  total_pets: number
  adoptable_pets: number
  adopted_pets: number
  treatment_pets: number
  deceased_pets: number
  total_adoptions: number
  pending_applications: number
  completed_follow_ups: number
  total_follow_ups: number
  adoption_rate: number
  follow_up_rate: number
}

export interface RescueListQuery {
  page?: number
  page_size?: number
  status?: string
  search?: string
}
