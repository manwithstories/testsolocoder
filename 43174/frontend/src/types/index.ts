export interface User {
  id: string;
  username: string;
  email: string;
  phone?: string;
  role: 'student' | 'merchant' | 'admin';
  status: 'pending' | 'active' | 'rejected' | 'banned';
  real_name?: string;
  school_name?: string;
  student_id?: string;
  student_card_url?: string;
  business_license?: string;
  avatar?: string;
  rating: number;
  rating_count: number;
  created_at: string;
  updated_at: string;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user: User;
}

export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
  phone?: string;
  role: 'student' | 'merchant';
  real_name?: string;
  school_name?: string;
  student_id?: string;
  student_card_url?: string;
  business_license?: string;
}

export interface Category {
  id: string;
  name: string;
  parent_id?: string;
  sort_order: number;
  created_at: string;
  updated_at: string;
  children?: Category[];
}

export interface Textbook {
  id: string;
  isbn: string;
  title: string;
  author?: string;
  course_name?: string;
  edition?: string;
  publisher?: string;
  original_price: number;
  price: number;
  condition: 'new' | 'like_new' | 'good' | 'fair';
  description?: string;
  cover_image?: string;
  status: 'available' | 'reserved' | 'sold';
  seller_id: string;
  category_id?: string;
  view_count: number;
  created_at: string;
  updated_at: string;
  seller?: User;
  category?: Category;
}

export interface CreateTextbookRequest {
  isbn: string;
  title: string;
  author?: string;
  course_name?: string;
  edition?: string;
  publisher?: string;
  original_price?: number;
  price: number;
  condition: string;
  description?: string;
  cover_image?: string;
  category_id?: string;
}

export interface Note {
  id: string;
  title: string;
  subject?: string;
  course_name?: string;
  description?: string;
  file_url: string;
  file_type?: string;
  file_size: number;
  cover_image?: string;
  uploader_id: string;
  category_id?: string;
  download_count: number;
  view_count: number;
  rating: number;
  rating_count: number;
  is_featured: boolean;
  created_at: string;
  updated_at: string;
  uploader?: User;
  category?: Category;
}

export interface CreateNoteRequest {
  title: string;
  subject?: string;
  course_name?: string;
  description?: string;
  file_url: string;
  file_type?: string;
  file_size?: number;
  cover_image?: string;
  category_id?: string;
}

export interface Transaction {
  id: string;
  textbook_id: string;
  seller_id: string;
  buyer_id: string;
  type: 'sell' | 'exchange';
  agreed_price: number;
  status: 'pending' | 'negotiating' | 'confirmed' | 'completed' | 'cancelled';
  exchange_item?: string;
  negotiation_history?: string;
  created_at: string;
  updated_at: string;
  textbook?: Textbook;
  seller?: User;
  buyer?: User;
}

export interface Order {
  id: string;
  order_no: string;
  buyer_id: string;
  seller_id: string;
  total_amount: number;
  status: 'pending' | 'paid' | 'shipped' | 'delivered' | 'completed' | 'cancelled' | 'refunded';
  payment_method?: string;
  payment_status?: string;
  transaction_id?: string;
  shipping_address?: string;
  tracking_number?: string;
  remark?: string;
  created_at: string;
  updated_at: string;
  buyer?: User;
  seller?: User;
  transaction?: Transaction;
  items?: OrderItem[];
  status_history?: OrderStatusHistory[];
}

export interface OrderItem {
  id: string;
  order_id: string;
  textbook_id: string;
  quantity: number;
  price: number;
  subtotal: number;
  created_at: string;
  textbook?: Textbook;
}

export interface OrderStatusHistory {
  id: string;
  order_id: string;
  status: string;
  remark?: string;
  created_at: string;
}

export interface Message {
  id: string;
  sender_id: string;
  receiver_id: string;
  content: string;
  is_read: boolean;
  related_order_id?: string;
  is_dispute: boolean;
  created_at: string;
  updated_at: string;
  sender?: User;
  receiver?: User;
}

export interface Review {
  id: string;
  user_id: string;
  target_type: 'textbook' | 'note';
  textbook_id?: string;
  note_id?: string;
  rating: number;
  content?: string;
  is_hidden: boolean;
  is_malicious: boolean;
  created_at: string;
  updated_at: string;
  user?: User;
}

export interface Notification {
  id: string;
  user_id: string;
  type: string;
  title: string;
  content: string;
  is_read: boolean;
  created_at: string;
}

export interface PaginatedResponse<T> {
  code: number;
  message: string;
  data: T[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

export interface Statistics {
  textbook_stats: {
    total_count: number;
    available_count: number;
    sold_count: number;
    total_value: number;
  };
  user_stats: {
    total_count: number;
    student_count: number;
    merchant_count: number;
    pending_count: number;
  };
  order_stats: {
    total_count: number;
    total_revenue: number;
    pending_count: number;
    completed_count: number;
  };
}
