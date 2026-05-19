export interface User {
  id: number;
  email: string;
  username: string;
  verified: boolean;
}

export interface Event {
  id: number;
  title: string;
  description: string;
  location: string;
  start_time: string;
  end_time: string;
  capacity: number;
  deadline: string;
  organizer_id: number;
  organizer: User;
  created_at: string;
}

export interface EventWithStatus extends Event {
  registered_count: number;
  is_registered: boolean;
  is_full: boolean;
  is_deadline_passed: boolean;
  can_register: boolean;
  is_organizer: boolean;
  cancel_count: number;
}

export interface Registration {
  id: number;
  user_id: number;
  event_id: number;
  user: User;
  event: Event;
  status: string;
  cancel_count: number;
  created_at: string;
  updated_at: string;
}

export interface ApiResponse<T> {
  data?: T;
  message?: string;
  error?: string;
}
