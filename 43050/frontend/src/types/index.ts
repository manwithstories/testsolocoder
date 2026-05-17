export interface User {
  id: number;
  username: string;
  email: string;
  avatar?: string;
  createdAt: string;
  updatedAt: string;
}

export interface Group {
  id: number;
  name: string;
  description: string;
  creatorId: number;
  creator?: User;
  members?: User[];
  inviteCode: string;
  createdAt: string;
  updatedAt: string;
}

export interface GroupMember {
  groupId: number;
  userId: number;
  joinedAt: string;
  isActive: boolean;
  user: User;
}

export type SplitType = 'equal' | 'ratio' | 'custom';

export interface ExpenseParticipant {
  id: number;
  expenseId: number;
  userId: number;
  user?: User;
  amount: number;
  ratio?: number;
  isSettled: boolean;
}

export interface Expense {
  id: number;
  groupId: number;
  title: string;
  amount: number;
  paidBy: number;
  payer?: User;
  splitType: SplitType;
  createdBy: number;
  version: number;
  expenseDate: string;
  createdAt: string;
  updatedAt: string;
  participants: ExpenseParticipant[];
}

export interface Settlement {
  id: number;
  groupId: number;
  fromUserId: number;
  fromUser?: User;
  toUserId: number;
  toUser?: User;
  amount: number;
  isPaid: boolean;
  paidAt?: string;
  createdAt: string;
  updatedAt: string;
}

export interface Balance {
  userId: number;
  balance: number;
}

export interface Transfer {
  fromUserId: number;
  fromUsername: string;
  toUserId: number;
  toUsername: string;
  amount: number;
}

export interface StatisticsSummary {
  totalPaid: number;
  totalOwed: number;
  netBalance: number;
  expenseCount: number;
}

export interface MonthlyStats {
  month: string;
  totalPaid: number;
  totalOwed: number;
}

export interface MemberStats {
  userId: number;
  username: string;
  totalPaid: number;
  totalOwed: number;
  balance: number;
}
