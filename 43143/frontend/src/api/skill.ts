import { apiRequest } from './index';
import { SkillCategory, SkillTag, Skill, SkillPosting, PaginatedData } from '@/types';

export const categoryApi = {
  create: (data: { name: string; icon?: string; sort_order?: number }) =>
    apiRequest.post<SkillCategory>('/categories', data),

  list: () => apiRequest.get<SkillCategory[]>('/categories'),

  update: (id: string, data: { name: string; icon?: string; sort_order?: number }) =>
    apiRequest.put<SkillCategory>(`/categories/${id}`, data),

  delete: (id: string) => apiRequest.delete(`/categories/${id}`),
};

export const tagApi = {
  create: (data: { name: string; category_id?: string }) =>
    apiRequest.post<SkillTag>('/tags', data),

  list: (params?: { category_id?: string }) =>
    apiRequest.get<SkillTag[]>('/tags', { params }),

  update: (id: string, data: { name: string; category_id?: string }) =>
    apiRequest.put<SkillTag>(`/tags/${id}`, data),

  delete: (id: string) => apiRequest.delete(`/tags/${id}`),
};

export const skillApi = {
  create: (data: any) => apiRequest.post<Skill>('/skills', data),

  list: (params?: { page?: number; page_size?: number; category_id?: string; keyword?: string }) =>
    apiRequest.get<PaginatedData<Skill>>('/skills', { params }),

  get: (id: string) => apiRequest.get<Skill>(`/skills/${id}`),

  update: (id: string, data: any) => apiRequest.put<Skill>(`/skills/${id}`, data),

  delete: (id: string) => apiRequest.delete(`/skills/${id}`),

  getPopular: (params?: { limit?: number }) =>
    apiRequest.get<Skill[]>('/skills/popular', { params }),
};

export const postingApi = {
  create: (data: any) => apiRequest.post<SkillPosting>('/postings', data),

  list: (params?: { page?: number; page_size?: number; skill_id?: string; teacher_id?: string }) =>
    apiRequest.get<PaginatedData<SkillPosting>>('/postings', { params }),

  get: (id: string) => apiRequest.get<SkillPosting>(`/postings/${id}`),

  update: (id: string, data: any) => apiRequest.put<SkillPosting>(`/postings/${id}`, data),

  delete: (id: string) => apiRequest.delete(`/postings/${id}`),

  match: (params?: { page?: number; page_size?: number; skill_id?: string; min_rating?: number; method?: string; max_price?: number }) =>
    apiRequest.get<PaginatedData<SkillPosting>>('/postings/match', { params }),
};
