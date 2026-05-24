export interface User {
  id: string;
  email: string;
  phone: string;
  nickname: string;
  avatar: string;
  bio: string;
  gender: string;
  birthday: string;
  location: string;
  latitude: number;
  longitude: number;
  role: UserRole;
  status: UserStatus;
  auth_type: AuthType;
  email_verified: boolean;
  phone_verified: boolean;
  skill_tags: SkillTag[];
  learn_tags: SkillTag[];
  teach_tags: SkillTag[];
  teaching_hours: number;
  student_count: number;
  rating: number;
  review_count: number;
  balance: number;
  last_login_at: string;
  created_at: string;
  updated_at: string;
}

export type UserRole = 'learner' | 'teacher' | 'both' | 'admin';
export type UserStatus = 'active' | 'suspended' | 'banned';
export type AuthType = 'email' | 'phone' | 'wechat' | 'google';

export interface SkillCategory {
  id: string;
  name: string;
  icon: string;
  sort_order: number;
  created_at: string;
  updated_at: string;
}

export interface SkillTag {
  id: string;
  name: string;
  category_id: string;
  category?: SkillCategory;
  usage_count: number;
  created_at: string;
  updated_at: string;
}

export type DifficultyLevel = 'beginner' | 'intermediate' | 'advanced' | 'expert';

export interface Skill {
  id: string;
  title: string;
  description: string;
  category_id: string;
  category?: SkillCategory;
  tags: SkillTag[];
  difficulty: DifficultyLevel;
  cover_image: string;
  video_url: string;
  prerequisites: string;
  outcomes: string;
  is_active: boolean;
  posting_count: number;
  rating: number;
  review_count: number;
  created_at: string;
  updated_at: string;
}

export type TeachingMethod = 'online' | 'offline' | 'both';
export type TeachingMode = 'one_to_one' | 'small_class' | 'group';

export interface SkillPosting {
  id: string;
  teacher_id: string;
  teacher?: User;
  skill_id: string;
  skill?: Skill;
  title: string;
  description: string;
  teaching_method: TeachingMethod;
  teaching_mode: TeachingMode;
  max_students: number;
  price_per_hour: number;
  currency: string;
  session_duration: number;
  location: string;
  latitude: number;
  longitude: number;
  availability: string;
  is_active: boolean;
  rating: number;
  review_count: number;
  booking_count: number;
  total_hours: number;
  created_at: string;
  updated_at: string;
}

export type BookingStatus = 'pending' | 'confirmed' | 'rejected' | 'cancelled' | 'completed' | 'no_show';
export type PaymentStatus = 'pending' | 'paid' | 'held' | 'released' | 'refunded' | 'failed';

export interface Booking {
  id: string;
  posting_id: string;
  posting?: SkillPosting;
  student_id: string;
  student?: User;
  teacher_id: string;
  teacher?: User;
  scheduled_start: string;
  scheduled_end: string;
  actual_start: string;
  actual_end: string;
  status: BookingStatus;
  price: number;
  platform_fee: number;
  teacher_earnings: number;
  note: string;
  reject_reason: string;
  cancel_reason: string;
  cancelled_by: string;
  reviewed_by_student: boolean;
  reviewed_by_teacher: boolean;
  payment_status: PaymentStatus;
  created_at: string;
  updated_at: string;
}

export interface Review {
  id: string;
  booking_id: string;
  booking?: Booking;
  reviewer_id: string;
  reviewer?: User;
  reviewee_id: string;
  reviewee?: User;
  posting_id: string;
  rating: number;
  content: string;
  is_public: boolean;
  is_featured: boolean;
  helpful_count: number;
  status: string;
  created_at: string;
  updated_at: string;
}

export interface Complaint {
  id: string;
  reporter_id: string;
  reporter?: User;
  respondent_id: string;
  respondent?: User;
  type: ComplaintType;
  target_id: string;
  title: string;
  description: string;
  evidence: string;
  status: ComplaintStatus;
  handler_id: string;
  handler?: User;
  result: string;
  handled_at: string;
  created_at: string;
  updated_at: string;
}

export type ComplaintType = 'review' | 'booking' | 'payment' | 'user';
export type ComplaintStatus = 'pending' | 'processing' | 'resolved' | 'rejected';

export type MessageType = 'text' | 'image' | 'file' | 'system';

export interface Message {
  id: string;
  sender_id: string;
  sender?: User;
  receiver_id: string;
  receiver?: User;
  type: MessageType;
  content: string;
  encrypted: boolean;
  file_url: string;
  file_name: string;
  file_size: number;
  is_read: boolean;
  read_at: string;
  created_at: string;
  updated_at: string;
}

export type PaymentMethod = 'alipay' | 'wechat' | 'bank_transfer' | 'cash';
export type TransactionType = 'payment' | 'refund' | 'withdraw' | 'deposit' | 'fee';
export type TransactionStatus = 'pending' | 'processing' | 'completed' | 'failed' | 'cancelled';

export interface Payment {
  id: string;
  booking_id: string;
  booking?: Booking;
  payer_id: string;
  payer?: User;
  payee_id: string;
  payee?: User;
  amount: number;
  currency: string;
  method: PaymentMethod;
  type: TransactionType;
  status: TransactionStatus;
  platform_fee: number;
  net_amount: number;
  escrow_release_at: string;
  released: boolean;
  released_at: string;
  transaction_id: string;
  third_party_id: string;
  failure_reason: string;
  remark: string;
  created_at: string;
  updated_at: string;
}

export interface Wallet {
  id: string;
  user_id: string;
  user?: User;
  balance: number;
  frozen: number;
  currency: string;
  created_at: string;
  updated_at: string;
}

export type ScheduleType = 'availability' | 'busy' | 'booking';
export type DayOfWeek = 'monday' | 'tuesday' | 'wednesday' | 'thursday' | 'friday' | 'saturday' | 'sunday';

export interface Schedule {
  id: string;
  user_id: string;
  user?: User;
  type: ScheduleType;
  day_of_week: DayOfWeek;
  specific_date: string;
  start_time: string;
  end_time: string;
  is_recurring: boolean;
  booking_id: string;
  booking?: Booking;
  title: string;
  description: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface ApiResponse<T = any> {
  code: number;
  message: string;
  data: T;
  errors?: any;
}

export interface PaginatedData<T = any> {
  items: T[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface TokenPair {
  access_token: string;
  refresh_token: string;
  expires_in: number;
}

export interface TeacherStats {
  teaching_hours: number;
  student_count: number;
  total_income: number;
  avg_rating: number;
}

export interface MonthlyReport {
  month: string;
  teaching_hours: number;
  student_count: number;
  income: number;
  bookings: number;
}
