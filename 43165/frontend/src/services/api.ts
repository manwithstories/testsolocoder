import axios, { AxiosInstance, InternalAxiosRequestConfig } from 'axios';

const api: AxiosInstance = axios.create({
  baseURL: '/api/v1',
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json',
  },
});

api.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

api.interceptors.response.use(
  (response) => {
    if (response.data?.code === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return response;
  },
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export const authApi = {
  register: (data: any) => api.post('/auth/register', data),
  login: (data: any) => api.post('/auth/login', data),
  getCurrentUser: () => api.get('/auth/me'),
  updateProfile: (data: any) => api.put('/auth/profile', data),
  changePassword: (data: any) => api.put('/auth/password', data),
};

export const jobApi = {
  getJobs: (params?: any) => api.get('/jobs', { params }),
  getMyJobs: (params?: any) => api.get('/jobs/mine', { params }),
  getJob: (id: string) => api.get(`/jobs/${id}`),
  createJob: (data: any) => api.post('/jobs', data),
  updateJob: (id: string, data: any) => api.put(`/jobs/${id}`, data),
  deleteJob: (id: string) => api.delete(`/jobs/${id}`),
  applyJob: (id: string, data?: any) => api.post(`/jobs/${id}/apply`, data),
  getApplications: (id: string, params?: any) => api.get(`/jobs/${id}/applications`, { params }),
  getMyApplications: (params?: any) => api.get('/applications/mine', { params }),
  reviewApplication: (id: string, data: any) => api.put(`/applications/${id}/review`, data),
};

export const scheduleApi = {
  getSchedules: (params?: any) => api.get('/schedules', { params }),
  getMySchedules: (params?: any) => api.get('/schedules/mine', { params }),
  getSchedule: (id: string) => api.get(`/schedules/${id}`),
  createSchedule: (data: any) => api.post('/schedules', data),
  batchCreateSchedules: (data: any) => api.post('/schedules/batch', data),
  updateSchedule: (id: string, data: any) => api.put(`/schedules/${id}`, data),
  deleteSchedule: (id: string) => api.delete(`/schedules/${id}`),
  checkConflict: (data: any) => api.post('/schedules/check-conflict', data),
  exportSchedules: (params?: any) => api.get('/schedules/export', { params, responseType: 'blob' }),
  generateQRCode: (id: string) => api.get(`/schedules/${id}/qrcode`),
};

export const checkInApi = {
  getCheckIns: (params?: any) => api.get('/checkins', { params }),
  getStats: () => api.get('/checkins/stats'),
  checkIn: (data: any) => api.post('/checkins', data),
  checkOut: (data: any) => api.post('/checkins/checkout', data),
  verifyFace: (data: any) => api.post('/checkins/verify-face', data),
  registerFace: (data: any) => api.post('/checkins/register-face', data),
};

export const salaryApi = {
  getSalaries: (params?: any) => api.get('/salaries', { params }),
  getSalary: (id: string) => api.get(`/salaries/${id}`),
  calculateSalary: (data: any) => api.post('/salaries', data),
  paySalary: (id: string) => api.post(`/salaries/${id}/pay`),
  batchPaySalary: (data: any) => api.post('/salaries/batch-pay', data),
  exportSalary: (id: string) => api.get(`/salaries/${id}/export`, { responseType: 'blob' }),
};

export const evaluationApi = {
  getEvaluations: (params?: any) => api.get('/evaluations', { params }),
  getMyEvaluations: (params?: any) => api.get('/evaluations/mine', { params }),
  getEvaluationStats: (id: string) => api.get(`/evaluations/${id}/stats`),
  createEvaluation: (data: any) => api.post('/evaluations', data),
  updateEvaluation: (id: string, data: any) => api.put(`/evaluations/${id}`, data),
};

export const templateApi = {
  getTemplates: (params?: any) => api.get('/job-templates', { params }),
  getTemplate: (id: string) => api.get(`/job-templates/${id}`),
  createTemplate: (data: any) => api.post('/job-templates', data),
  updateTemplate: (id: string, data: any) => api.put(`/job-templates/${id}`, data),
  deleteTemplate: (id: string) => api.delete(`/job-templates/${id}`),
  batchImportTemplates: (data: any) => api.post('/job-templates/batch-import', data),
  applyTemplate: (data: any) => api.post('/job-templates/apply', data),
};

export const matchApi = {
  matchTemporaries: (data: any) => api.post('/match/temporaries', data),
  quickAssign: (data: any) => api.post('/match/quick-assign', data),
  getMatchHistory: (params?: any) => api.get('/match/history', { params }),
};

export const statsApi = {
  getOverview: () => api.get('/stats/overview'),
  getActivityStats: (params?: any) => api.get('/stats/activities', { params }),
  getPersonnelStats: () => api.get('/stats/personnel'),
  getSalaryStats: (params?: any) => api.get('/stats/salary', { params }),
};

export default api;
