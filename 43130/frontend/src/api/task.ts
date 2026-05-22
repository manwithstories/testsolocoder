import { request } from '@/utils/request'
import { ApiResponse, Task, TaskTemplate } from '@/types'

export const taskApi = {
  create: (weddingId: number, data: Partial<Task>) => 
    request.post<ApiResponse<Task>>(`/weddings/${weddingId}/tasks`, data),
  
  getList: (weddingId: number, params?: { status?: string; category?: string }) => 
    request.get<ApiResponse<Task[]>>(`/weddings/${weddingId}/tasks`, { params }),
  
  getById: (weddingId: number, id: number) => 
    request.get<ApiResponse<{ task: Task; subtasks: Task[] }>>(`/weddings/${weddingId}/tasks/${id}`),
  
  update: (weddingId: number, id: number, data: Partial<Task>) => 
    request.put<ApiResponse<Task>>(`/weddings/${weddingId}/tasks/${id}`, data),
  
  delete: (weddingId: number, id: number) => 
    request.delete<ApiResponse>(`/weddings/${weddingId}/tasks/${id}`),
  
  updateStatus: (weddingId: number, id: number, status: string) => 
    request.put<ApiResponse<Task>>(`/weddings/${weddingId}/tasks/${id}/status`, { status }),
  
  getCategories: (weddingId: number) => 
    request.get<ApiResponse<string[]>>(`/weddings/${weddingId}/tasks/categories`),
  
  getTemplates: () => 
    request.get<ApiResponse<TaskTemplate[]>>('/task-templates'),
  
  createTemplate: (data: Partial<TaskTemplate>) => 
    request.post<ApiResponse<TaskTemplate>>('/task-templates', data),
  
  deleteTemplate: (id: number) => 
    request.delete<ApiResponse>(`/task-templates/${id}`),
  
  applyTemplate: (weddingId: number, templateId: number) => 
    request.post<ApiResponse>(`/task-templates/${templateId}/apply?wedding_id=${weddingId}`)
}
