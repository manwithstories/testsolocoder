import request, { PageData, PageParams } from '@/utils/request'

export interface Event {
  id: number
  user_id: number
  artist_id: number
  title: string
  description?: string
  cover_url: string
  venue: string
  address?: string
  city?: string
  longitude?: number
  latitude?: number
  start_time: string
  end_time: string
  door_time?: string
  ticket_price: number
  total_tickets: number
  sold_tickets: number
  max_per_user: number
  status: number
  published_at?: string
  view_count: number
  like_count: number
  share_count: number
  seat_map?: string
  user?: {
    id: number
    username: string
    nickname: string
    avatar: string
  }
  created_at: string
  updated_at: string
}

export interface Order {
  id: number
  order_no: string
  user_id: number
  event_id: number
  artist_id: number
  total_amount: number
  quantity: number
  status: number
  payment_method?: string
  paid_at?: string
  remark?: string
  user?: any
  event?: Event
  tickets?: Ticket[]
  created_at: string
  updated_at: string
}

export interface Ticket {
  id: number
  order_id: number
  user_id: number
  event_id: number
  seat_row: number
  seat_col: number
  seat_label: string
  qr_code: string
  status: number
  used_at?: string
  event?: Event
  created_at: string
  updated_at: string
}

export interface CreateEventParams {
  title: string
  description?: string
  cover_url?: string
  venue: string
  address?: string
  city?: string
  longitude?: number
  latitude?: number
  start_time: string
  end_time: string
  door_time?: string
  ticket_price: number
  total_tickets: number
  max_per_user?: number
  seat_map?: string
}

export interface PurchaseTicketParams {
  event_id: number
  quantity: number
  seats?: Array<{ row: number; col: number }>
}

export const eventApi = {
  list: (params: PageParams & { keyword?: string; city?: string; status?: number; artist_id?: number }) => 
    request.get<PageData<Event>>('/events', params),
  getAll: (params: PageParams) => request.get<PageData<Event>>('/events', params),
  getById: (id: number) => request.get<Event>(`/events/${id}`),
  getMyEvents: (params: PageParams) => request.get<PageData<Event>>('/events/my', params),
  getMyTickets: (params: PageParams) => request.get<PageData<any>>('/tickets/my', params),
  create: (data: CreateEventParams) => request.post<Event>('/events', data),
  update: (id: number, data: Partial<CreateEventParams>) => request.put(`/events/${id}`, data),
  delete: (id: number) => request.delete(`/events/${id}`),
  publish: (id: number) => request.post(`/events/${id}/publish`),
  getStats: (id: number) => request.get(`/events/${id}/stats`),
  getSeats: (id: number) => request.get(`/events/${id}/seats`),
  purchase: (data: PurchaseTicketParams) => request.post<Order>('/events/purchase', data),
  getOrders: (params: PageParams) => request.get<PageData<Order>>('/orders', params),
  getOrderById: (id: number) => request.get<Order>(`/orders/${id}`),
  getOrdersByArtist: (artistId: number, params: PageParams) => 
    request.get<PageData<Order>>(`/orders/artist/${artistId}`, params),
  getTickets: (params: PageParams) => request.get<PageData<Ticket>>('/tickets', params),
  useTicket: (id: number) => request.post(`/tickets/${id}/use`)
}
