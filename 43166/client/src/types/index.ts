export type UserRole = 'admin' | 'entrepreneur' | 'agent'

export const UserRoles = {
  ADMIN: 'admin' as UserRole,
  ENTREPRENEUR: 'entrepreneur' as UserRole,
  AGENT: 'agent' as UserRole
} as const

export type UserStatus = 'active' | 'inactive' | 'banned'

export interface User {
  id: number
  username: string
  realName: string
  email: string
  phone: string
  role: UserRole
  status: UserStatus
  avatar: string
  lastLoginAt: string | null
  createdAt: string
  updatedAt: string
  agentProfile?: AgentProfile
}

export interface AgentProfile {
  id: number
  userId: number
  employeeNo: string
  specialtyTags: string
  maxApplications: number
  currentApps: number
  workStartTime: string
  workEndTime: string
  status: string
  performanceScore: number
  totalHandled: number
  createdAt: string
  updatedAt: string
}

export type ApplicationStatus =
  | 'draft'
  | 'pending_review'
  | 'reviewing'
  | 'processing'
  | 'completed'
  | 'rejected'
  | 'cancelled'
  | 'payment_pending'

export type CompanyType = 'llc' | 'joint_stock' | 'sole' | 'partnership'

export interface Application {
  id: number
  applicationNo: string
  entrepreneurId: number
  agentId: number | null
  companyName: string
  companyType: CompanyType
  registeredCapital: number
  businessScope: string
  registeredAddress: string
  shareholderInfo: string
  idCardFront: string
  idCardBack: string
  licensePreview: string
  otherMaterials: string
  status: ApplicationStatus
  reviewComments: string
  currentStep: string
  progressPercent: number
  submittedAt: string | null
  completedAt: string | null
  rejectedAt: string | null
  createdAt: string
  updatedAt: string
  entrepreneur?: User
  agent?: User
  processSteps?: ProcessStep[]
  fee?: ApplicationFee
  notifications?: Notification[]
}

export type ProcessStepType =
  | 'naming'
  | 'registration'
  | 'tax'
  | 'bank'
  | 'seal'
  | 'completion'

export type ProcessStepStatus = 'pending' | 'in_progress' | 'completed' | 'failed' | 'skipped'

export interface ProcessStep {
  id: number
  applicationId: number
  stepType: ProcessStepType
  stepName: string
  stepOrder: number
  status: ProcessStepStatus
  description: string
  remark: string
  certificateFile: string
  handlerId: number | null
  startedAt: string | null
  completedAt: string | null
  createdAt: string
  updatedAt: string
  application?: Application
  handler?: User
}

export type FeeStatus = 'pending' | 'paid' | 'refunded'

export interface FeeItem {
  id: number
  applicationFeeId: number
  itemName: string
  amount: number
  description: string
  createdAt: string
}

export interface ApplicationFee {
  id: number
  applicationId: number
  totalAmount: number
  discountAmount: number
  paidAmount: number
  status: FeeStatus
  paymentMethod: string
  paymentTime: string | null
  transactionNo: string
  createdAt: string
  updatedAt: string
  application?: Application
  feeItems?: FeeItem[]
}

export interface FeeStandard {
  id: number
  companyType: CompanyType
  namingFee: number
  registrationFee: number
  taxFee: number
  bankFee: number
  sealFee: number
  serviceFee: number
  capitalRate: number
  createdAt: string
  updatedAt: string
}

export interface DiscountPolicy {
  id: number
  name: string
  code: string
  type: string
  value: number
  minAmount: number
  maxDiscount: number
  startDate: string | null
  endDate: string | null
  isActive: boolean
  createdAt: string
  updatedAt: string
}

export type NotificationType = 'system' | 'email' | 'sms'

export type NotificationStatus = 'pending' | 'sent' | 'read' | 'failed'

export interface Notification {
  id: number
  userId: number
  applicationId: number | null
  type: NotificationType
  title: string
  content: string
  status: NotificationStatus
  isRead: boolean
  sentAt: string | null
  readAt: string | null
  createdAt: string
  user?: User
  application?: Application
}

export interface NotificationTemplate {
  id: number
  code: string
  name: string
  type: string
  title: string
  content: string
  variables: string
  isActive: boolean
  createdAt: string
  updatedAt: string
}

export interface ExportTask {
  id: number
  userId: number
  type: string
  fileName: string
  filePath: string
  status: string
  expiresAt: string | null
  downloaded: boolean
  downloadedAt: string | null
  params: string
  errorMsg: string
  createdAt: string
  updatedAt: string
}

export interface ApiResponse<T = any> {
  code: number
  message: string
  data?: T
  errors?: any
}

export interface PaginatedData<T> {
  list: T[]
  total: number
  page: number
  pageSize: number
}

export interface LoginRequest {
  username: string
  password: string
}

export interface RegisterRequest {
  username: string
  password: string
  realName: string
  email?: string
  phone?: string
  role: UserRole
}

export interface LoginResponse {
  token: string
  userInfo: User
}

export interface CreateApplicationRequest {
  companyName: string
  companyType: CompanyType
  registeredCapital: number
  businessScope: string
  registeredAddress: string
  shareholderInfo: string
}

export interface UpdateApplicationRequest {
  companyName?: string
  companyType?: string
  registeredCapital?: number
  businessScope?: string
  registeredAddress?: string
  shareholderInfo?: string
}

export interface PayFeeRequest {
  applicationId: number
  paymentMethod: string
}

export interface CalculateFeeRequest {
  applicationId: number
  companyType: CompanyType
  capital: number
  discountCode?: string
}

export interface SendNotificationRequest {
  userId: number
  applicationId?: number
  type: NotificationType
  templateCode: string
  variables?: Record<string, string>
}

export interface OverviewStats {
  totalApplications: number
  pendingApplications: number
  processingApps: number
  completedApps: number
  rejectedApps: number
  totalRevenue: number
  todayApplications: number
  todayRevenue: number
  totalAgents: number
  activeAgents: number
  totalEntrepreneurs: number
}

export interface StatusDistribution {
  status: string
  count: number
}

export interface CompanyTypeDistribution {
  companyType: string
  count: number
  totalAmount: number
}

export interface AgentPerformance {
  agentId: number
  agentName: string
  employeeNo: string
  totalHandled: number
  completedCount: number
  inProgressCount: number
  avgDuration: number
  totalRevenue: number
}

export interface TimeSeriesData {
  date: string
  count: number
  amount: number
}

export interface CreateAgentRequest {
  username: string
  password: string
  realName: string
  email?: string
  phone: string
  employeeNo: string
  specialtyTags: string
  maxApps?: number
}

export interface UpdateAgentProfileRequest {
  specialtyTags?: string
  maxApps?: number
  workStartTime?: string
  workEndTime?: string
  status?: string
}
