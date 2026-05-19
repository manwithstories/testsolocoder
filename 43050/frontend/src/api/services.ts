import api from './client';
import {
  User,
  Group,
  Expense,
  Settlement,
  Balance,
  Transfer,
  StatisticsSummary,
  MonthlyStats,
  MemberStats,
  GroupMember,
} from '../types';

export interface RegisterData {
  username: string;
  email: string;
  password: string;
  avatar?: string;
}

export interface LoginData {
  email: string;
  password: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export const authAPI = {
  register: (data: RegisterData) => api.post<AuthResponse>('/auth/register', data),
  login: (data: LoginData) => api.post<AuthResponse>('/auth/login', data),
  getCurrentUser: () => api.get<User>('/auth/me'),
  searchUsers: (query: string) => api.get<User[]>(`/users/search?q=${query}`),
};

export interface CreateGroupData {
  name: string;
  description?: string;
}

export const groupAPI = {
  createGroup: (data: CreateGroupData) => api.post<Group>('/groups', data),
  getUserGroups: () => api.get<Group[]>('/groups'),
  getGroup: (id: number) => api.get<Group>(`/groups/${id}`),
  getGroupMembers: (id: number) => api.get<GroupMember[]>(`/groups/${id}/members`),
  joinGroup: (inviteCode: string) => api.post<Group>('/groups/join', { inviteCode }),
  leaveGroup: (id: number) => api.post(`/groups/${id}/leave`),
};

export interface ParticipantInput {
  userId: number;
  amount?: number;
  ratio?: number;
}

export interface CreateExpenseData {
  title: string;
  amount: number;
  paidBy: number;
  splitType: string;
  expenseDate?: string;
  participants: ParticipantInput[];
}

export const expenseAPI = {
  createExpense: (groupId: number, data: CreateExpenseData) =>
    api.post<Expense>(`/groups/${groupId}/expenses`, data),
  getGroupExpenses: (groupId: number) => api.get<Expense[]>(`/groups/${groupId}/expenses`),
  getExpense: (id: number) => api.get<Expense>(`/expenses/${id}`),
  updateExpense: (id: number, data: CreateExpenseData) =>
    api.put<Expense>(`/expenses/${id}`, data),
  deleteExpense: (id: number) => api.delete(`/expenses/${id}`),
};

export const settlementAPI = {
  getGroupBalances: (groupId: number) => api.get<Balance[]>(`/groups/${groupId}/balances`),
  getOptimalTransfers: (groupId: number) => api.get<Transfer[]>(`/groups/${groupId}/transfers`),
  getGroupStats: (groupId: number) => api.get<MemberStats[]>(`/groups/${groupId}/stats`),
  createSettlement: (groupId: number, data: { fromUserId: number; toUserId: number; amount: number }) =>
    api.post<Settlement>(`/groups/${groupId}/settlements`, data),
  getGroupSettlements: (groupId: number) => api.get<Settlement[]>(`/groups/${groupId}/settlements`),
  markSettlementPaid: (id: number) => api.patch<Settlement>(`/settlements/${id}/paid`),
};

export const statsAPI = {
  getSummary: (params?: { startDate?: string; endDate?: string; groupId?: string; memberId?: string }) =>
    api.get<StatisticsSummary>('/stats/summary', { params }),
  getMonthlyStats: () => api.get<MonthlyStats[]>('/stats/monthly'),
  getHistory: (params?: { startDate?: string; endDate?: string; groupId?: string; memberId?: string }) =>
    api.get<Expense[]>('/stats/history', { params }),
};
