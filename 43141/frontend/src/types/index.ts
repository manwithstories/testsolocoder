export interface User {
  id: number
  email: string
  full_name: string
  role: 'admin' | 'captain' | 'player' | 'referee'
  phone?: string
  team_id?: number
  is_active: boolean
}

export interface League {
  id: number
  name: string
  description: string
  sport: string
  status: string
  created_by: number
  seasons?: Season[]
  created_at: string
  updated_at: string
}

export interface Season {
  id: number
  league_id: number
  name: string
  start_date: string
  end_date: string
  format: string
  group_count: number
  points_for_win: number
  points_for_draw: number
  points_for_loss: number
  max_teams: number
  registration_fee: number
  venue_fee: number
  status: string
  custom_rules?: string
  awards?: string
  matches?: Match[]
  created_at: string
  updated_at: string
}

export interface Team {
  id: number
  name: string
  logo?: string
  captain_id: number
  captain?: User
  description?: string
  contact_email?: string
  contact_phone?: string
  players?: Player[]
  registrations?: Registration[]
  created_at: string
  updated_at: string
}

export interface Player {
  id: number
  team_id: number
  name: string
  number: number
  position: string
  birth_date: string
  user_id?: number
  user?: User
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface Registration {
  id: number
  season_id: number
  season?: Season
  team_id: number
  team?: Team
  group_name: string
  status: string
  paid: boolean
  paid_at?: string
  note?: string
  created_at: string
  updated_at: string
}

export interface Match {
  id: number
  season_id: number
  season?: Season
  round: number
  group_name: string
  home_team_id: number
  home_team?: Team
  away_team_id: number
  away_team?: Team
  venue_id?: number
  venue?: Venue
  referee_id?: number
  referee?: User
  match_time?: string
  home_score?: number
  away_score?: number
  has_ot: boolean
  ot_home_score?: number
  ot_away_score?: number
  has_penalty: boolean
  pen_home_score?: number
  pen_away_score?: number
  status: string
  knockout_stage?: string
  winner_team_id?: number
  created_at: string
  updated_at: string
}

export interface Venue {
  id: number
  name: string
  address?: string
  capacity: number
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface RefereeAssignment {
  id: number
  match_id: number
  match?: Match
  referee_id: number
  referee?: User
  status: string
  assigned_at: string
  responded_at?: string
  created_at: string
  updated_at: string
}

export interface Fee {
  id: number
  season_id: number
  season?: Season
  team_id: number
  team?: Team
  type: string
  amount: number
  status: string
  paid_at?: string
  invoice_no: string
  note?: string
  created_at: string
  updated_at: string
}

export interface Notification {
  id: number
  user_id: number
  user?: User
  title: string
  content: string
  type: string
  is_read: boolean
  created_at: string
}

export interface PlayerStat {
  id: number
  player_id: number
  player?: Player
  match_id: number
  match?: Match
  season_id: number
  team_id: number
  goals: number
  assists: number
  fouls: number
  yellow_card: number
  red_card: number
  minutes: number
  created_at: string
  updated_at: string
}

export interface Standings {
  season_id: number
  team_id: number
  team_name: string
  group_name: string
  played: number
  wins: number
  draws: number
  losses: number
  goals_for: number
  goals_against: number
  goal_diff: number
  points: number
}

export interface PlayerRanking {
  player_id: number
  player_name: string
  team_name: string
  matches: number
  goals: number
  assists: number
  fouls: number
  yellow_card: number
  red_card: number
  minutes: number
}
