export type TransactionType = 'income' | 'expense'
export type TransactionStatus = 'pending' | 'completed' | 'failed' | 'refunded'

export interface Transaction {
  id: string
  rental_id?: string
  rental?: Rental
  payer_id: string
  payer?: User
  payee_id: string
  payee?: User
  amount: number
  currency: string
  transaction_type: TransactionType
  description?: string
  status: TransactionStatus
  payment_method?: string
  transaction_ref?: string
  exchange_rate?: number
  original_amount?: number
  original_currency?: string
  platform_fee: number
  dock_fee: number
  net_amount: number
  paid_at?: string
  failed_reason?: string
  created_at: string
  updated_at: string
}

export interface Settlement {
  id: string
  user_id: string
  user?: User
  period_start: string
  period_end: string
  total_income: number
  total_expense: number
  currency: string
  status: string
  paid_at?: string
  created_at: string
  updated_at: string
}

export interface FinancialSummary {
  total_income: number
  total_expense: number
  net_balance: number
  platform_fee: number
  dock_fee: number
}
