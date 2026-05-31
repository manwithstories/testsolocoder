export type ReviewType = 'ship' | 'dock'

export interface Review {
  id: string
  rental_id: string
  rental?: Rental
  reviewer_id: string
  reviewer?: User
  target_type: ReviewType
  target_id: string
  rating: number
  content?: string
  response?: string
  response_by?: string
  responder?: User
  response_at?: string
  is_recommended: boolean
  helpful_count: number
  is_deleted: boolean
  deleted_at?: string
  created_at: string
  updated_at: string
}

export interface CreateReviewRequest {
  rental_id: string
  target_type: ReviewType
  target_id: string
  rating: number
  content?: string
  is_recommended?: boolean
}

export interface RespondToReviewRequest {
  response: string
}
