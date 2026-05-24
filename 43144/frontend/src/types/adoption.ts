import { Pet } from './pet'
import { User } from './user'

export interface AdoptionApplication {
  id: number
  pet_id: number
  pet?: Pet
  adopter_id: number
  adopter?: User
  rescue_id: number
  status: 'pending' | 'approved' | 'rejected' | 'signed' | 'completed' | 'cancelled'
  reason: string
  living_situation: string
  pet_experience: string
  family_members: number
  has_children: boolean
  has_other_pets: boolean
  other_pets_desc: string
  housing_type: string
  income_level: string
  can_afford_vet: boolean
  agree_to_visit: boolean
  reviewed_by?: number
  reviewed_at?: string
  reject_reason?: string
  signed_at?: string
  completed_at?: string
  created_at: string
  updated_at: string
}

export interface AdoptionAgreement {
  id: number
  application_id: number
  application?: AdoptionApplication
  adopter_sign: boolean
  adopter_signed_at?: string
  rescue_sign: boolean
  rescue_signed_at?: string
  agreement_terms: string
  agreement_file?: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface FollowUpRecord {
  id: number
  application_id: number
  application?: AdoptionApplication
  pet_id: number
  pet?: Pet
  adopter_id: number
  rescue_id: number
  follow_up_date?: string
  health_status: string
  living_condition: string
  notes: string
  photo_evidence?: string
  recorded_by: number
  created_at: string
  updated_at: string
}

export interface CreateAdoptionRequest {
  pet_id: number
  reason: string
  living_situation: string
  pet_experience?: string
  family_members?: number
  has_children?: boolean
  has_other_pets?: boolean
  other_pets_desc?: string
  housing_type: string
  income_level?: string
  can_afford_vet?: boolean
  agree_to_visit: boolean
}

export interface ReviewAdoptionRequest {
  action: 'approve' | 'reject'
  reject_reason?: string
}

export interface CreateFollowUpRequest {
  application_id: number
  follow_up_date: string
  health_status?: string
  living_condition?: string
  notes?: string
}

export interface AdoptionListQuery {
  page?: number
  page_size?: number
  status?: string
  pet_id?: number
  rescue_id?: number
}
