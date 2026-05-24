import { apiRequest } from './index';
import { Booking, Review, Complaint, PaginatedData } from '@/types';

export const bookingApi = {
  create: (data: { posting_id: string; scheduled_start: string; scheduled_end: string; note?: string }) =>
    apiRequest.post<Booking>('/bookings', data),

  list: (params: { role: string; page?: number; page_size?: number; status?: string }) =>
    apiRequest.get<PaginatedData<Booking>>('/bookings', { params }),

  get: (id: string) => apiRequest.get<Booking>(`/bookings/${id}`),

  confirm: (id: string) => apiRequest.put(`/bookings/${id}/confirm`),

  reject: (id: string, reason?: string) =>
    apiRequest.put(`/bookings/${id}/reject`, { reason }),

  cancel: (id: string, reason?: string) =>
    apiRequest.put(`/bookings/${id}/cancel`, { reason }),

  complete: (id: string) => apiRequest.put(`/bookings/${id}/complete`),
};

export const reviewApi = {
  create: (data: { booking_id: string; rating: number; content: string; is_public: boolean }) =>
    apiRequest.post<Review>('/reviews', data),

  getByPosting: (postingId: string, params?: { page?: number; page_size?: number }) =>
    apiRequest.get<PaginatedData<Review>>(`/reviews/posting/${postingId}`, { params }),

  getByUser: (userId: string, params?: { page?: number; page_size?: number }) =>
    apiRequest.get<PaginatedData<Review>>(`/reviews/user/${userId}`, { params }),
};

export const complaintApi = {
  create: (data: { type: string; target_id: string; title: string; description: string; evidence?: string }) =>
    apiRequest.post<Complaint>('/complaints', data),

  list: (params?: { page?: number; page_size?: number; status?: string }) =>
    apiRequest.get<PaginatedData<Complaint>>('/complaints', { params }),

  handle: (id: string, result: string) =>
    apiRequest.put(`/complaints/${id}/handle`, { result }),
};
