// 角色类型
export type Role = 'manufacturer' | 'designer' | 'owner'

// 通用分页参数
export interface PaginationParams {
  page: number
  pageSize: number
}

// 通用分页响应
export interface PaginationResult<T> {
  list: T[]
  total: number
  page: number
  pageSize: number
}

// 用户信息
export interface User {
  id: number
  username: string
  nickname: string
  role: Role
  email?: string
  phone?: string
  avatar?: string
  status: number
  createdAt: string
  updatedAt: string
}

// 登录请求
export interface LoginRequest {
  username: string
  password: string
}

// 登录响应
export interface LoginResponse {
  token: string
  userInfo: User
}

// 产品分类
export interface ProductCategory {
  id: number
  name: string
  parentId?: number
  sort?: number
}

// 产品图片
export interface ProductImage {
  id: number
  productId: number
  url: string
  sort?: number
  createdAt?: string
}

// 产品自定义选项类型
export type OptionType = 'size' | 'material' | 'color'

// 产品自定义选项值
export interface ProductOptionValue {
  id?: number
  optionId?: number
  value: string
  priceAdjustment?: number
  sort?: number
}

// 产品自定义选项
export interface ProductOption {
  id?: number
  productId?: number
  type: OptionType
  name: string
  required?: boolean
  values: ProductOptionValue[]
}

// 产品
export interface Product {
  id: number
  name: string
  sku: string
  description?: string
  imageUrl?: string
  images?: ProductImage[]
  manufacturerId: number
  manufacturerName?: string
  category?: string
  categoryId?: number
  price: number
  stock: number
  isHot?: boolean
  status: number
  options?: ProductOption[]
  createdAt: string
  updatedAt: string
}

// 产品查询参数
export interface ListProductsParams extends PaginationParams {
  keyword?: string
  category?: string
  categoryId?: number
  status?: number
}

// 产品表单提交
export interface ProductFormData {
  id?: number
  name: string
  sku?: string
  description?: string
  categoryId?: number
  category?: string
  price: number
  stock: number
  isHot?: boolean
  status?: number
  images?: ProductImage[]
  options?: ProductOption[]
}

// 订单状态枚举
export const OrderStatus = {
  PENDING: 0,
  QUOTED: 1,
  CONFIRMED: 2,
  PRODUCING: 3,
  SHIPPED: 4,
  COMPLETED: 5,
  CANCELLED: 9
} as const

export type OrderStatusValue = (typeof OrderStatus)[keyof typeof OrderStatus]

export const OrderStatusMap: Record<number, { label: string; color: string }> = {
  0: { label: '待报价', color: 'warning' },
  1: { label: '待确认', color: 'info' },
  2: { label: '已确认', color: 'primary' },
  3: { label: '生产中', color: '' },
  4: { label: '已发货', color: '' },
  5: { label: '已完成', color: 'success' },
  9: { label: '已取消', color: 'danger' }
}

// 订单自定义选项选择
export interface OrderOptionSelection {
  optionId?: number
  type: OptionType
  name: string
  value: string
  priceAdjustment?: number
}

// 订单项
export interface OrderItem {
  id: number
  orderId: number
  productId: number
  productName: string
  productImage?: string
  quantity: number
  unitPrice: number
  totalPrice: number
  options?: OrderOptionSelection[]
}

// 订单状态流转
export interface OrderStatusLog {
  id: number
  orderId: number
  status: number
  operatorId?: number
  operatorName?: string
  remark?: string
  createdAt: string
}

// 订单
export interface Order {
  id: number
  orderNo: string
  productId: number
  productName?: string
  productImage?: string
  quantity: number
  unitPrice: number
  totalPrice: number
  status: number
  manufacturerId: number
  manufacturerName?: string
  designerId?: number
  designerName?: string
  ownerId?: number
  ownerName?: string
  designId?: number
  remark?: string
  items?: OrderItem[]
  statusLogs?: OrderStatusLog[]
  createdAt: string
  updatedAt: string
}

// 订单查询参数
export interface ListOrdersParams extends PaginationParams {
  keyword?: string
  orderNo?: string
  status?: number
  startTime?: string
  endTime?: string
}

// 报价请求
export interface QuoteOrderParams {
  unitPrice: number
  remark?: string
}

// 订单通用操作
export interface OrderActionParams {
  remark?: string
}

// 设计项目
export interface DesignProject {
  id: number
  name: string
  description?: string
  productId?: number
  productName?: string
  designerId: number
  designerName?: string
  ownerId?: number
  ownerName?: string
  status: number
  deadline?: string
  files?: DesignFile[]
  createdAt: string
  updatedAt: string
}

// 设计文件
export interface DesignFile {
  id: number
  projectId: number
  fileName: string
  fileUrl: string
  fileType: string
  fileSize: number
  uploadedAt: string
}

// 交付物
export interface Delivery {
  id: number
  projectId: number
  projectName?: string
  title: string
  description?: string
  fileUrl: string
  fileName: string
  designerId: number
  designerName?: string
  status: number
  reviewerId?: number
  reviewerName?: string
  reviewComment?: string
  reviewedAt?: string
  createdAt: string
  updatedAt: string
}

// 评审
export interface Review {
  id: number
  deliveryId: number
  deliveryTitle?: string
  reviewerId: number
  reviewerName?: string
  status: number
  comment?: string
  score?: number
  createdAt: string
  updatedAt: string
}

// 工单
export interface Ticket {
  id: number
  title: string
  description: string
  type: number
  status: number
  priority: number
  creatorId: number
  creatorName?: string
  handlerId?: number
  handlerName?: string
  relatedOrderId?: number
  relatedProjectId?: number
  createdAt: string
  updatedAt: string
}

// 统计数据
export interface Statistics {
  totalUsers: number
  totalProducts: number
  totalOrders: number
  totalDesigns: number
  totalDeliveries: number
  totalTickets: number
  recentOrders: Order[]
  recentDeliveries: Delivery[]
  orderStatusCounts: { status: number; count: number }[]
  monthlyData: { month: string; orders: number; designs: number }[]
}

// 通用 API 响应
export interface ApiResponse<T = unknown> {
  code: number
  message: string
  data: T
}