import { RescueStation } from './rescue'
import { User } from './user'

export interface Pet {
  id: number
  archive_number: string
  name: string
  species: string
  breed: string
  age: string
  gender: string
  weight: number
  color: string
  description: string
  status: 'adoptable' | 'adopted' | 'treatment' | 'deceased'
  photos: string
  videos: string
  health_status: string
  vaccinated: boolean
  neutered: boolean
  rescue_id: number
  rescue?: RescueStation
  adopter_id?: number
  adopter?: User
  found_location?: string
  found_date?: string
  adopted_date?: string
  medical_history: string
  personality: string
  special_needs: string
  created_at: string
  updated_at: string
}

export interface CreatePetRequest {
  name: string
  species: string
  breed?: string
  age?: string
  gender: string
  weight?: number
  color?: string
  description?: string
  health_status?: string
  vaccinated?: boolean
  neutered?: boolean
  found_location?: string
  found_date?: string
  medical_history?: string
  personality?: string
  special_needs?: string
}

export interface PetListQuery {
  page?: number
  page_size?: number
  status?: string
  species?: string
  gender?: string
  search?: string
  rescue_id?: number
  breed?: string
}

export interface PetStatus {
  ADOPTABLE: 'adoptable'
  ADOPTED: 'adopted'
  TREATMENT: 'treatment'
  DECEASED: 'deceased'
}
