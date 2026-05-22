import { request } from '@/utils/request'
import { ApiResponse, BudgetItem, Payment } from '@/types'

export const budgetApi = {
  getSummary: (weddingId: number) => 
    request.get<ApiResponse<any>>(`/weddings/${weddingId}/budget/summary`),
  
  getItems: (weddingId: number) => 
    request.get<ApiResponse<BudgetItem[]>>(`/weddings/${weddingId}/budget/items`),
  
  createItem: (weddingId: number, data: Partial<BudgetItem>) => 
    request.post<ApiResponse<BudgetItem>>(`/weddings/${weddingId}/budget/items`, data),
  
  updateItem: (weddingId: number, id: number, data: Partial<BudgetItem>) => 
    request.put<ApiResponse<BudgetItem>>(`/weddings/${weddingId}/budget/items/${id}`, data),
  
  deleteItem: (weddingId: number, id: number) => 
    request.delete<ApiResponse>(`/weddings/${weddingId}/budget/items/${id}`),
  
  getPayments: (weddingId: number) => 
    request.get<ApiResponse<Payment[]>>(`/weddings/${weddingId}/budget/payments`),
  
  recordPayment: (weddingId: number, data: Partial<Payment> & { budget_item_id: number; amount: number }) => 
    request.post<ApiResponse>(`/weddings/${weddingId}/budget/payments`, data),
  
  getCategories: (weddingId: number) => 
    request.get<ApiResponse<string[]>>(`/weddings/${weddingId}/budget/categories`),
  
  checkAlerts: (weddingId: number) => 
    request.get<ApiResponse<any[]>>(`/weddings/${weddingId}/budget/alerts`)
}
