export interface User {
  id: number
  username: string
  realName: string
  phone: string
  email: string
  role: string
  avatar: string
  status: number
  createdAt: string
  updatedAt: string
}

export interface Property {
  id: number
  title: string
  community: string
  address: string
  area: number
  layout: string
  floor: string
  rent: number
  deposit: number
  paymentType: string
  description: string
  status: number
  ownerId: number
  owner?: User
  images?: PropertyImage[]
  facilities?: Facility[]
  region: string
  building: string
  roomNo: string
  createdAt: string
  updatedAt: string
}

export interface PropertyImage {
  id: number
  propertyId: number
  url: string
  sort: number
}

export interface Facility {
  id: number
  name: string
  icon: string
}

export interface Tenant {
  id: number
  name: string
  phone: string
  idCard: string
  email: string
  avatar: string
  status: number
  createdAt: string
  updatedAt: string
}

export interface Appointment {
  id: number
  tenantId: number
  tenant?: Tenant
  propertyId: number
  property?: Property
  visitTime: string
  status: number
  remark: string
  createdAt: string
  updatedAt: string
}

export interface Contract {
  id: number
  tenantId: number
  tenant?: Tenant
  propertyId: number
  property?: Property
  startDate: string
  endDate: string
  rent: number
  deposit: number
  paymentType: string
  status: number
  fileUrl: string
  isReminded: boolean
  createdAt: string
  updatedAt: string
}

export interface RentRecord {
  id: number
  contractId: number
  contract?: Contract
  tenantId: number
  propertyId: number
  month: string
  amount: number
  status: number
  dueDate: string
  paidAt?: string
  lateFee: number
  remark: string
  createdAt: string
  updatedAt: string
}

export interface RepairOrder {
  id: number
  tenantId: number
  tenant?: Tenant
  propertyId: number
  property?: Property
  title: string
  description: string
  category: string
  images: string
  priority: number
  status: number
  handlerId?: number
  handler?: User
  processNote: string
  completedAt?: string
  createdAt: string
  updatedAt: string
}

export interface UtilityFee {
  id: number
  propertyId: number
  property?: Property
  type: string
  month: string
  totalAmount: number
  units: number
  unitPrice: number
  status: number
  dueDate: string
  paidAt?: string
  createdAt: string
  updatedAt: string
}

export interface Notice {
  id: number
  title: string
  content: string
  type: number
  building: string
  publisherId: number
  publisher?: User
  isTop: number
  status: number
  createdAt: string
  updatedAt: string
}

export interface PaginationParams {
  page?: number
  pageSize?: number
  keyword?: string
  status?: number
  [key: string]: any
}

export interface PaginationResult<T> {
  list: T[]
  total: number
  page: number
  pageSize: number
}

export interface ApiResponse<T> {
  code: number
  message: string
  data: T
}
