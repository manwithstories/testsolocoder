import axios from 'axios';
import { EventWithStatus, Registration, User } from '../types';

const api = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json',
  },
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export const authAPI = {
  register: (data: { email: string; username: string; password: string }) =>
    api.post('/auth/register', data),
  login: (data: { email: string; password: string }) =>
    api.post('/auth/login', data),
  verify: (data: { email: string; code: string }) =>
    api.post('/auth/verify', data),
  getProfile: () => api.get<User>('/user/profile'),
  getMyRegistrations: () => api.get<Registration[]>('/user/registrations'),
};

export const eventAPI = {
  getEvents: () => api.get<EventWithStatus[]>('/events'),
  getEvent: (id: number) => api.get(`/events/${id}`),
  createEvent: (data: {
    title: string;
    description: string;
    location: string;
    start_time: string;
    end_time: string;
    capacity: number;
    deadline: string;
  }) => api.post('/events', data),
  updateEvent: (id: number, data: any) => api.put(`/events/${id}`, data),
  deleteEvent: (id: number) => api.delete(`/events/${id}`),
  registerEvent: (id: number) => api.post(`/events/${id}/register`),
  cancelRegistration: (id: number) => api.post(`/events/${id}/cancel`),
  getRegistrations: (id: number) => api.get<Registration[]>(`/events/${id}/registrations`),
  exportRegistrations: (id: number) => {
    const token = localStorage.getItem('token');
    window.open(`/api/events/${id}/export?token=${token}`, '_blank');
  },
};

export default api;
