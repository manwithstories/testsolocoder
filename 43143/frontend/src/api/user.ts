import { apiRequest } from './index';
import { User, TokenPair, PaginatedData, SkillTag } from '@/types';

export const authApi = {
  register: (data: { email?: string; phone?: string; password: string; nickname: string; auth_type: string }) =>
    apiRequest.post<{ user: User; token: TokenPair }>('/auth/register', data),

  login: (data: { email?: string; phone?: string; password: string }) =>
    apiRequest.post<{ user: User; token: TokenPair }>('/auth/login', data),

  refreshToken: (refresh_token: string) =>
    apiRequest.post<{ token: TokenPair }>('/auth/refresh', { refresh_token }),
};

export const userApi = {
  getProfile: () => apiRequest.get<User>('/users/me'),

  updateProfile: (data: any) => apiRequest.put<User>('/users/me', data),

  getUser: (id: string) => apiRequest.get<User>(`/users/${id}`),

  listUsers: (params: { page?: number; page_size?: number; keyword?: string }) =>
    apiRequest.get<PaginatedData<User>>('/users', { params }),

  addSkillTags: (tag_ids: string[]) =>
    apiRequest.post('/users/me/tags', { tag_ids }),

  removeSkillTag: (tag_id: string) =>
    apiRequest.delete(`/users/me/tags/${tag_id}`),

  getSkillTags: () => apiRequest.get<SkillTag[]>('/users/me/tags'),
};
