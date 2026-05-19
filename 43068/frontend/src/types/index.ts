export interface User {
  id: number;
  email: string;
  first_name: string;
  last_name: string;
}

export interface Client {
  id: number;
  user_id: number;
  name: string;
  email: string;
  phone: string;
  address: string;
  company: string;
  contract_url: string;
  default_rate: number;
  created_at: string;
  projects?: Project[];
}

export type ProjectStatus = 'draft' | 'active' | 'completed' | 'archived';

export interface Project {
  id: number;
  user_id: number;
  client_id: number;
  name: string;
  description: string;
  status: ProjectStatus;
  hourly_rate: number;
  deadline?: string;
  budget: number;
  created_at: string;
  client?: Client;
  milestones?: Milestone[];
  time_entries?: TimeEntry[];
  invoices?: Invoice[];
}

export interface Milestone {
  id: number;
  project_id: number;
  title: string;
  description: string;
  due_date?: string;
  completed: boolean;
  created_at: string;
}

export interface TimeEntry {
  id: number;
  user_id: number;
  project_id: number;
  date: string;
  hours: number;
  description: string;
  start_time?: string;
  end_time?: string;
  is_timer: boolean;
  billable: boolean;
  created_at: string;
  project?: Project;
}

export type InvoiceStatus = 'draft' | 'sent' | 'paid' | 'overdue' | 'cancelled';

export interface Invoice {
  id: number;
  user_id: number;
  client_id: number;
  project_id?: number;
  invoice_number: string;
  status: InvoiceStatus;
  issue_date: string;
  due_date: string;
  subtotal: number;
  tax_rate: number;
  tax_amount: number;
  total: number;
  notes: string;
  created_at: string;
  client?: Client;
  project?: Project;
  items?: InvoiceItem[];
}

export interface InvoiceItem {
  id: number;
  invoice_id: number;
  description: string;
  quantity: number;
  unit_price: number;
  amount: number;
  created_at: string;
}

export interface TokenPair {
  access_token: string;
  refresh_token: string;
  expires_in: number;
  token_type: string;
}

export interface AuthResponse {
  user: User;
  tokens: TokenPair;
}

export interface DashboardStats {
  monthly_revenue: number;
  total_clients: number;
  active_projects: number;
  total_hours: number;
  overdue_invoices: number;
  pending_invoices: number;
  project_progress: ProjectProgressStats[];
  monthly_earnings: MonthlyEarningsStats[];
  overdue_reminders: OverdueReminder[];
}

export interface ProjectProgressStats {
  id: number;
  name: string;
  status: string;
  total_hours: number;
  budget: number;
  earned_amount: number;
  progress: number;
  days_left: number;
}

export interface MonthlyEarningsStats {
  month: string;
  revenue: number;
  hours: number;
}

export interface OverdueReminder {
  id: number;
  type: string;
  title: string;
  due_date: string;
  days_overdue: number;
  amount?: number;
}

export interface ApiResponse<T> {
  success: boolean;
  message?: string;
  data?: T;
  error?: string;
}

export interface PaginatedResponse<T> extends ApiResponse<T> {
  meta: {
    total: number;
    page: number;
    per_page: number;
    total_pages: number;
  };
}
