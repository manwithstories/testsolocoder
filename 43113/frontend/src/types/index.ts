export interface ApiResponse<T = any> {
  code: number
  message: string
  data?: T
}

export interface PageResult<T = any> {
  list: T[]
  total: number
  page: number
  pageSize: number
}

export interface User {
  id: number
  username: string
  email: string
  nickname: string
  avatar: string
  phone: string
  level: number
  points: number
  isExpert: boolean
  expertStatus: string
  role: string
  status: string
  bio: string
  lastLoginAt: string
  createdAt: string
  updatedAt: string
}

export interface LoginResponse {
  token: string
  user: User
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
  nickname?: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface Category {
  id: number
  name: string
  description: string
  icon: string
  sortOrder: number
  status: string
  createdAt: string
  updatedAt: string
}

export interface Tag {
  id: number
  name: string
  description: string
  usageCount: number
  createdAt: string
  updatedAt: string
}

export interface Question {
  id: number
  userId: number
  user: User
  title: string
  content: string
  categoryId: number
  category: Category
  tags: Tag[]
  views: number
  answerCount: number
  likeCount: number
  collectCount: number
  hotScore: number
  rewardPoints: number
  isSolved: boolean
  acceptedAnswerId: number | null
  status: string
  auditStatus: string
  createdAt: string
  updatedAt: string
}

export interface QuestionDetail {
  question: Question
  tags: Tag[]
  user: User
  category: Category
  answers: AnswerWithUser[]
  answerCount: number
  isFavorited: boolean
}

export interface Answer {
  id: number
  questionId: number
  userId: number
  content: string
  likeCount: number
  dislikeCount: number
  collectCount: number
  isAccepted: boolean
  status: string
  auditStatus: string
  createdAt: string
  updatedAt: string
}

export interface AnswerWithUser {
  answer: Answer
  user: User
  comments: CommentWithUser[]
  isLiked: boolean
  isFavorited: boolean
}

export interface Comment {
  id: number
  userId: number
  answerId: number
  content: string
  likeCount: number
  status: string
  auditStatus: string
  createdAt: string
  updatedAt: string
}

export interface CommentWithUser {
  comment: Comment
  user: User
  isLiked: boolean
}

export interface CreateQuestionRequest {
  title: string
  content: string
  categoryId: number
  tagIds?: number[]
  rewardPoints?: number
}

export interface CreateAnswerRequest {
  questionId: number
  content: string
}

export interface CreateCommentRequest {
  answerId: number
  content: string
}

export interface AuditRecord {
  id: number
  adminId: number
  admin: User
  targetType: string
  targetId: number
  action: string
  reason: string
  status: string
  operatorIp: string
  createdAt: string
}

export interface Report {
  id: number
  reporterId: number
  reporter: User
  targetType: string
  targetId: number
  reason: string
  description: string
  status: string
  handlerId: number | null
  handler: User | null
  handleResult: string
  handledAt: string | null
  createdAt: string
}

export interface SensitiveWord {
  id: number
  word: string
  category: string
  level: number
  replaceTo: string
  createdAt: string
}

export interface PointLog {
  id: number
  userId: number
  user: User
  type: string
  points: number
  balance: number
  description: string
  refType: string
  refId: number
  createdAt: string
}

export interface Reward {
  id: number
  name: string
  description: string
  image: string
  pointsCost: number
  stock: number
  status: string
  createdAt: string
  updatedAt: string
}

export interface RewardExchange {
  id: number
  userId: number
  user: User
  rewardId: number
  reward: Reward
  pointsCost: number
  status: string
  createdAt: string
}

export interface Favorite {
  id: number
  userId: number
  targetType: string
  targetId: number
  createdAt: string
}

export interface Follow {
  id: number
  followerId: number
  follower: User
  followingType: string
  followingId: number
  createdAt: string
}

export interface Notification {
  id: number
  userId: number
  type: string
  title: string
  content: string
  refType: string
  refId: number
  isRead: boolean
  createdAt: string
}

export interface ExpertApplication {
  id: number
  userId: number
  user: User
  field: string
  description: string
  status: string
  reviewerId: number | null
  reviewer: User | null
  reviewRemark: string
  reviewedAt: string | null
  createdAt: string
}

export interface DashboardStats {
  totalQuestions: number
  totalAnswers: number
  totalUsers: number
  totalComments: number
  pendingAuditCount: number
  todayNewQuestions: number
  todayNewAnswers: number
  todayNewUsers: number
}

export interface ActivityReport {
  date: string
  newQuestions: number
  newAnswers: number
  newUsers: number
  newComments: number
}

export interface AuditReport {
  date: string
  reviewedCount: number
  approvedCount: number
  rejectedCount: number
}

export interface QuestionListQuery {
  page?: number
  pageSize?: number
  categoryId?: number
  tagId?: number
  userId?: number
  keyword?: string
  sort?: string
}

export interface FavoriteRequest {
  targetType: string
  targetId: number
}

export interface FollowRequest {
  followingType: string
  followingId: number
}
