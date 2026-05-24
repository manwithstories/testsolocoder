declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}

interface UserInfo {
  user_id: number
  username: string
  nickname: string
  avatar: string
  role: 'client' | 'pilot' | 'owner'
  verify_status: 'pending' | 'approved' | 'rejected'
}

interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

interface PageData<T = any> {
  list: T[]
  total: number
  page: number
  page_size: number
}

interface Drone {
  id: number
  owner_id: number
  owner?: UserInfo
  name: string
  brand: string
  model: string
  serial_no: string
  weight: number
  battery_life: number
  gimbal_spec: string
  camera_spec: string
  max_speed: number
  max_altitude: number
  region: string
  description: string
  images: string
  price_per_day: number
  deposit: number
  status: 'offline' | 'online' | 'rented' | 'maintenance'
  rating: number
  rating_count: number
  created_at: string
}

interface RentalOrder {
  id: number
  order_no: string
  user_id: number
  user?: UserInfo
  drone_id: number
  drone?: Drone
  start_date: string
  end_date: string
  return_date?: string
  region: string
  address: string
  contact_name: string
  contact_phone: string
  total_days: number
  price_per_day: number
  rental_fee: number
  deposit: number
  insurance_fee: number
  late_fee: number
  total_amount: number
  paid_amount: number
  refund_amount: number
  status: 'pending' | 'paid' | 'picked' | 'returned' | 'completed' | 'cancelled'
  remark: string
  cancel_reason: string
  created_at: string
}

interface AerialService {
  id: number
  service_no: string
  user_id: number
  user?: UserInfo
  pilot_id?: number
  pilot?: UserInfo
  title: string
  description: string
  region: string
  address: string
  service_date?: string
  service_time: string
  duration: number
  budget_min: number
  budget_max: number
  final_price: number
  status: 'open' | 'assigned' | 'progress' | 'completed' | 'cancelled'
  images: string
  remark: string
  created_at: string
}

interface ServiceBid {
  id: number
  service_id: number
  pilot_id: number
  pilot?: UserInfo
  price: number
  message: string
  status: 'pending' | 'accepted' | 'rejected'
  created_at: string
}

interface FlightRecord {
  id: number
  order_id?: number
  service_id?: number
  drone_id: number
  drone?: Drone
  pilot_id: number
  pilot?: UserInfo
  start_point: string
  end_point: string
  route: string
  altitude_max: number
  altitude_avg: number
  duration: number
  distance: number
  flight_date: string
  flight_log: string
  images: string
  remark: string
  created_at: string
}

interface InsuranceClaim {
  id: number
  claim_no: string
  order_id: number
  order?: RentalOrder
  user_id: number
  user?: UserInfo
  damage_desc: string
  damage_images: string
  estimated_cost: number
  actual_cost: number
  status: 'pending' | 'reviewing' | 'approved' | 'rejected'
  reviewer_id?: number
  reviewer?: UserInfo
  review_remark: string
  reviewed_at?: string
  deducted_amount: number
  created_at: string
}

interface Review {
  id: number
  type: 'rental' | 'service'
  order_id?: number
  service_id?: number
  reviewer_id: number
  reviewer?: UserInfo
  reviewee_id: number
  reviewee?: UserInfo
  drone_id?: number
  rating: number
  content: string
  images: string
  reply: string
  created_at: string
}

interface RevenueStats {
  date: string
  amount: number
  count: number
}

interface RegionStats {
  region: string
  count: number
  amount: number
}

interface DroneStats {
  drone_id: number
  drone_name: string
  total_days: number
  utilization: number
  income: number
}
