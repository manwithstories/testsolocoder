export type UserRole = 'admin' | 'teacher' | 'student'
export type UserStatus = 'pending' | 'active' | 'inactive' | 'banned'

export interface User {
  id: string
  email: string
  firstName: string
  lastName: string
  phone?: string
  avatarUrl?: string
  role: UserRole
  status: UserStatus
  timezone: string
  language: string
  emailVerified: boolean
  lastLoginAt?: string
  createdAt: string
  updatedAt: string
  teacherProfile?: TeacherProfile
  studentProfile?: StudentProfile
  wallet?: Wallet
}

export interface TeacherProfile {
  id: string
  userId: string
  bio?: string
  education?: string
  experience?: string
  certifications?: string
  hourlyRate: number
  currency: string
  rating: number
  reviewCount: number
  totalSessions: number
  totalHours: number
  isVerified: boolean
  approvalStatus: 'pending' | 'approved' | 'rejected'
  approvalNotes?: string
  approvedAt?: string
  approvedBy?: string
  resumeUrl?: string
  idCardUrl?: string
  subjects?: TeacherSubject[]
  availabilities?: AvailabilitySlot[]
  user?: User
}

export interface TeacherSubject {
  id: string
  teacherId: string
  subjectId: string
  level: string
  customRate: number
  isPrimary: boolean
  subject?: Subject
}

export interface Subject {
  id: string
  name: string
  category?: string
  description?: string
  iconUrl?: string
  isActive: boolean
  sortOrder: number
}

export interface AvailabilitySlot {
  id: string
  teacherId: string
  dayOfWeek: number
  startTime: string
  endTime: string
  isRecurring: boolean
  startDate?: string
  endDate?: string
  isBooked: boolean
}

export interface StudentProfile {
  id: string
  userId: string
  gradeLevel?: string
  school?: string
  learningStyle?: string
  preferredTime?: string
  notes?: string
  parentName?: string
  parentPhone?: string
  parentEmail?: string
  assessmentStatus: 'not_started' | 'in_progress' | 'completed'
  learningGoals?: LearningGoal[]
}

export interface LearningGoal {
  id: string
  studentId: string
  subjectId: string
  title: string
  description?: string
  targetScore: number
  currentScore: number
  deadline?: string
  isAchieved: boolean
  achievedAt?: string
  subject?: Subject
}

export type BookingStatus = 'pending' | 'confirmed' | 'completed' | 'cancelled' | 'rescheduled'

export interface Booking {
  id: string
  studentId: string
  teacherId: string
  subjectId: string
  startTime: string
  endTime: string
  duration: number
  hourlyRate: number
  totalAmount: number
  status: BookingStatus
  notes?: string
  studentNotes?: string
  teacherNotes?: string
  cancelledBy?: string
  cancelledAt?: string
  cancelReason?: string
  createdAt: string
  updatedAt: string
  student?: User
  teacher?: User
  subject?: Subject
  videoSession?: VideoSession
}

export interface VideoSession {
  id: string
  bookingId: string
  sessionId: string
  roomName: string
  token: string
  joinUrl: string
  status: string
  actualStartAt?: string
  actualEndAt?: string
  actualDuration: number
  qualityScore: number
  recordingUrl?: string
}

export interface Wallet {
  id: string
  userId: string
  balance: number
  currency: string
  totalIncome: number
  totalSpent: number
  transactions?: Transaction[]
}

export type TransactionType = 'deposit' | 'withdrawal' | 'payment' | 'refund' | 'commission' | 'transfer'
export type TransactionStatus = 'pending' | 'completed' | 'failed' | 'cancelled'

export interface Transaction {
  id: string
  walletId: string
  userId: string
  type: TransactionType
  amount: number
  currency: string
  balanceAfter: number
  status: TransactionStatus
  description: string
  referenceId?: string
  referenceType?: string
  bookingId?: string
  paymentMethod?: string
  transactionId?: string
  failureReason?: string
  completedAt?: string
  createdAt: string
  booking?: Booking
}

export interface LessonNote {
  id: string
  bookingId: string
  teacherId: string
  studentId: string
  subjectId: string
  title: string
  content: string
  tags?: string
  isPrivate: boolean
  createdAt: string
  updatedAt: string
  teacher?: User
  student?: User
  subject?: Subject
}

export type HomeworkStatus = 'pending' | 'submitted' | 'graded'

export interface Homework {
  id: string
  bookingId: string
  teacherId: string
  studentId: string
  subjectId: string
  title: string
  description: string
  dueDate: string
  status: HomeworkStatus
  maxScore: number
  attachments?: string
  createdAt: string
  teacher?: User
  student?: User
  subject?: Subject
  submission?: HomeworkSubmission
}

export interface HomeworkSubmission {
  id: string
  homeworkId: string
  studentId: string
  content?: string
  attachments?: string
  submittedAt: string
  score?: number
  feedback?: string
  gradedAt?: string
  gradedBy?: string
  student?: User
}

export interface Review {
  id: string
  bookingId: string
  teacherId: string
  studentId: string
  reviewerId: string
  revieweeId: string
  rating: number
  content: string
  tags?: string
  isAnonymous: boolean
  teacherReply?: string
  teacherRepliedAt?: string
  isHidden: boolean
  createdAt: string
  reviewer?: User
  reviewee?: User
  booking?: Booking
}

export type MessageType = 'text' | 'file' | 'image' | 'system' | 'booking'

export interface Message {
  id: string
  senderId: string
  receiverId: string
  bookingId?: string
  type: MessageType
  content: string
  isRead: boolean
  readAt?: string
  isDeleted: boolean
  createdAt: string
  sender?: User
  receiver?: User
  booking?: Booking
  files?: MessageFile[]
}

export interface MessageFile {
  id: string
  messageId: string
  fileName: string
  fileUrl: string
  fileSize: number
  fileType: string
}

export type NotificationType =
  | 'booking_created'
  | 'booking_updated'
  | 'booking_cancelled'
  | 'booking_reminder'
  | 'lesson_start'
  | 'lesson_end'
  | 'new_review'
  | 'payment_received'
  | 'payment_failed'
  | 'withdraw_approved'
  | 'withdraw_rejected'
  | 'new_message'
  | 'homework_assigned'
  | 'homework_graded'
  | 'milestone'
  | 'system'

export interface Notification {
  id: string
  userId: string
  type: NotificationType
  title: string
  content: string
  data?: string
  isRead: boolean
  readAt?: string
  createdAt: string
}

export interface Milestone {
  id: string
  studentId: string
  subjectId?: string
  title: string
  description?: string
  type: 'session_complete' | 'goal_achieved' | 'score_improve' | 'streak' | 'custom'
  icon?: string
  color?: string
  isAchieved: boolean
  achievedAt?: string
  subject?: Subject
}

export interface Feedback {
  id: string
  bookingId: string
  teacherId: string
  studentId: string
  subjectId: string
  content: string
  type: string
  progress?: string
  suggestions?: string
  nextSteps?: string
  createdAt: string
  teacher?: User
  student?: User
  subject?: Subject
}
