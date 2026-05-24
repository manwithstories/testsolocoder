export interface User {
  id: number
  email: string
  nickname: string
  avatar: string
  role: string
  status: number
  last_login?: string
  created_at: string
}

export interface Role {
  id: number
  name: string
  description: string
}

export interface Survey {
  id: number
  title: string
  description: string
  cover_image: string
  user_id: number
  user?: User
  status: number
  start_time?: string
  end_time?: string
  anonymous: boolean
  has_password: boolean
  max_responses: number
  max_per_user: number
  requires_login: boolean
  allow_resume: boolean
  response_count: number
  category: string
  tags: string
  questions?: Question[]
  created_at: string
  updated_at: string
}

export type QuestionType = 'single_choice' | 'multi_choice' | 'fill_in' | 'rating' | 'ranking' | 'matrix'

export interface Question {
  id: number
  survey_id: number
  title: string
  type: QuestionType
  is_required: boolean
  order_index: number
  description?: string
  placeholder?: string
  min_value?: number
  max_value?: number
  default_value?: string
  validation_rule?: string
  display_logic?: string
  options?: Option[]
  logic_jumps?: LogicJump[]
  status: number
}

export interface Option {
  id: number
  question_id: number
  text: string
  order_index: number
  is_other: boolean
  jump_target?: string
  score: number
}

export interface LogicJump {
  id: number
  question_id: number
  condition: string
  value: string
  jump_to: string
}

export interface SurveyResponse {
  id: number
  survey_id: number
  user_id?: number
  session_id: string
  status: number
  start_time?: string
  complete_time?: string
  duration: number
  answers?: Answer[]
  created_at: string
}

export interface Answer {
  id: number
  response_id: number
  question_id: number
  option_id?: number
  text_value?: string
  numeric_value?: number
  matrix_values?: string
  ranking_order?: string
}

export interface DistributionLink {
  id: number
  survey_id: number
  link_token: string
  channel: string
  max_uses: number
  use_count: number
  expires_at?: string
  is_active: boolean
  full_url?: string
  qrcode_url?: string
  created_at: string
}

export interface Invitation {
  id: number
  survey_id: number
  email: string
  link_token: string
  status: number
  sent_at?: string
  opened_at?: string
  responded_at?: string
  retry_count: number
  error_message?: string
  created_at: string
}

export interface Statistics {
  survey_id: number
  total_responses: number
  completed_count: number
  in_progress_count: number
  abandoned_count: number
  completion_rate: number
  avg_duration: number
  question_stats: QuestionStat[]
  time_distribution: TimeDataPoint[]
  channel_stats: ChannelStat[]
}

export interface QuestionStat {
  question_id: number
  question_title: string
  question_type: string
  response_count: number
  option_stats?: OptionStat[]
  text_stats?: TextStat
  rating_stats?: RatingStat
  word_cloud?: WordCloudItem[]
}

export interface OptionStat {
  option_id: number
  text: string
  count: number
  percent: number
}

export interface TextStat {
  avg_length: number
  max_length: number
  min_length: number
}

export interface RatingStat {
  average: number
  median: number
  std_dev: number
  min: number
  max: number
  distribution: Record<number, number>
}

export interface WordCloudItem {
  word: string
  count: number
}

export interface TimeDataPoint {
  date: string
  count: number
}

export interface ChannelStat {
  channel: string
  count: number
  percent: number
}

export interface ApiResponse<T = any> {
  code: number
  message: string
  data?: T
}

export interface PaginatedData<T> {
  items: T[]
  total: number
  page: number
  page_size: number
}
