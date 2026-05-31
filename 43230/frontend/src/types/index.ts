export type UserRole = 'designer' | 'printer' | 'customer' | 'admin';
export type UserStatus = 'active' | 'inactive' | 'banned';
export type ModelStatus = 'draft' | 'published' | 'rejected' | 'banned';
export type LicenseType = 'per_purchase' | 'subscription' | 'commercial';
export type OrderStatus = 'pending' | 'paid' | 'printing' | 'quality_check' | 'shipped' | 'delivered' | 'completed' | 'cancelled' | 'refunded';
export type PrintQuality = 'draft' | 'standard' | 'high' | 'ultra';
export type MaterialType = 'pla' | 'abs' | 'petg' | 'tpu' | 'resin' | 'nylon' | 'pc';
export type PrinterStatus = 'idle' | 'printing' | 'maintenance' | 'offline';

export interface User {
  id: string;
  email: string;
  username: string;
  role: UserRole;
  status: UserStatus;
  phone: string;
  avatar: string;
  credit_score: number;
  real_name: string;
  balance: number;
  email_verified: boolean;
  designer_profile?: DesignerProfile;
  printer_profile?: PrinterProfile;
  created_at: string;
  updated_at: string;
}

export interface DesignerProfile {
  id: string;
  user_id: string;
  nickname: string;
  bio: string;
  portfolio_url: string;
  total_models: number;
  total_sales: number;
  rating: number;
  rating_count: number;
  specialties: string[];
  experience_years: number;
  created_at: string;
  updated_at: string;
}

export interface PrinterProfile {
  id: string;
  user_id: string;
  company_name: string;
  address: string;
  max_print_size: string;
  supported_materials: string[];
  rating: number;
  rating_count: number;
  total_orders: number;
  total_revenue: number;
  created_at: string;
  updated_at: string;
}

export interface Model3D {
  id: string;
  designer_id: string;
  designer?: User;
  title: string;
  description: string;
  category: string;
  tags: string[];
  price: number;
  license_type: LicenseType;
  subscription_price: number;
  status: ModelStatus;
  thumbnail_url: string;
  file_url: string;
  file_size: number;
  file_type: string;
  version: string;
  volume: number;
  bounding_box: string;
  print_time_hours: number;
  recommended_materials: string[];
  polygon_count: number;
  view_count: number;
  download_count: number;
  purchase_count: number;
  favorite_count: number;
  rating: number;
  rating_count: number;
  is_featured: boolean;
  previous_versions?: ModelVersion[];
  created_at: string;
  updated_at: string;
}

export interface ModelVersion {
  id: string;
  model_id: string;
  version_number: string;
  change_log: string;
  created_at: string;
}

export interface ModelPurchase {
  id: string;
  model_id: string;
  user_id: string;
  model?: Model3D;
  user?: User;
  purchase_type: LicenseType;
  amount: number;
  transaction_id: string;
  expires_at?: string;
  created_at: string;
}

export interface Material {
  id: string;
  type: MaterialType;
  name: string;
  description: string;
  price_per_gram: number;
  print_speed: number;
  density: number;
  strength: string;
  temperature_resistance: string;
  is_available: boolean;
  image_url: string;
  created_at: string;
}

export interface PrintOrder {
  id: string;
  order_no: string;
  customer_id: string;
  printer_id?: string;
  model_id: string;
  customer?: User;
  printer?: User;
  model?: Model3D;
  quantity: number;
  material_id: string;
  material?: Material;
  color: string;
  quality: PrintQuality;
  layer_height: number;
  infill_percent: number;
  supports: boolean;
  estimated_volume: number;
  estimated_weight: number;
  estimated_print_time: number;
  base_price: number;
  material_cost: number;
  service_fee: number;
  shipping_fee: number;
  total_amount: number;
  status: OrderStatus;
  shipping_address: string;
  tracking_number: string;
  notes: string;
  cancelled_reason?: string;
  print_started_at?: string;
  print_finished_at?: string;
  shipped_at?: string;
  delivered_at?: string;
  completed_at?: string;
  cancelled_at?: string;
  created_at: string;
  updated_at: string;
}

export interface OrderHistory {
  id: string;
  order_id: string;
  status: OrderStatus;
  description: string;
  created_at: string;
}

export interface Settlement {
  id: string;
  order_id: string;
  order?: PrintOrder;
  total_amount: number;
  platform_fee: number;
  designer_share: number;
  printer_share: number;
  status: string;
  settled_at?: string;
  created_at: string;
}

export interface PrinterDevice {
  id: string;
  printer_id: string;
  name: string;
  model: string;
  manufacturer: string;
  max_print_size: string;
  max_print_volume: number;
  supported_materials: string[];
  supported_qualities: PrintQuality[];
  status: PrinterStatus;
  total_print_hours: number;
  total_print_jobs: number;
  ip_address: string;
  firmware_version: string;
  created_at: string;
  updated_at: string;
}

export interface MaterialInventory {
  id: string;
  printer_id: string;
  material_id: string;
  material?: Material;
  color: string;
  quantity_grams: number;
  reorder_level: number;
  last_updated: string;
  created_at: string;
}

export interface PrintSchedule {
  id: string;
  printer_id: string;
  device_id: string;
  order_id: string;
  order?: PrintOrder;
  scheduled_start?: string;
  scheduled_end?: string;
  actual_start?: string;
  actual_end?: string;
  status: string;
  priority: number;
  created_at: string;
}

export interface Review {
  id: string;
  order_id: string;
  order?: PrintOrder;
  customer_id: string;
  customer?: User;
  model_id: string;
  model?: Model3D;
  printer_id: string;
  printer?: User;
  model_rating: number;
  print_rating: number;
  model_comment: string;
  print_comment: string;
  images: string[];
  is_anonymous: boolean;
  designer_reply?: string;
  printer_reply?: string;
  helpful_count: number;
  created_at: string;
  updated_at: string;
}

export interface Notification {
  id: string;
  user_id: string;
  title: string;
  content: string;
  type: string;
  related_id: string;
  is_read: boolean;
  created_at: string;
  read_at?: string;
}

export interface Transaction {
  id: string;
  user_id: string;
  order_id?: string;
  type: 'income' | 'expense' | 'refund';
  amount: number;
  balance_after: number;
  description: string;
  payment_method: string;
  transaction_no: string;
  status: string;
  created_at: string;
}

export interface PriceEstimateRequest {
  model_id: string;
  quantity: number;
  material_id: string;
  color?: string;
  quality: PrintQuality;
  layer_height: number;
  infill_percent: number;
  supports: boolean;
}

export interface PriceEstimateResponse {
  base_price: number;
  material_cost: number;
  service_fee: number;
  shipping_fee: number;
  total_amount: number;
  estimated_volume: number;
  estimated_weight: number;
  estimated_print_time: number;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  username: string;
  password: string;
  role: UserRole;
  phone?: string;
  real_name?: string;
  company_name?: string;
}

export interface LoginResponse {
  access_token: string;
  refresh_token: string;
  expires_at: number;
  user: User;
}

export interface ApiResponse<T = any> {
  data?: T;
  total?: number;
  page?: number;
  size?: number;
  message?: string;
  error?: string;
}

export interface JWTPayload {
  user_id: string;
  email: string;
  username: string;
  role: UserRole;
  exp: number;
  iat: number;
  iss: string;
  sub: string;
}
