export type WeatherCondition = 'sunny' | 'cloudy' | 'rainy' | 'stormy' | 'foggy'

export interface VoyageLog {
  id: string
  rental_id: string
  rental?: Rental
  ship_id: string
  ship?: Ship
  captain_id: string
  captain?: User
  log_date: string
  departure_port?: string
  arrival_port?: string
  departure_time?: string
  arrival_time?: string
  distance?: number
  engine_hours?: number
  fuel_consumed?: number
  fuel_unit?: string
  avg_speed?: number
  max_speed?: number
  weather: WeatherCondition
  wind_speed?: number
  wind_direction?: string
  wave_height?: number
  water_temperature?: number
  air_temperature?: number
  passenger_count?: number
  crew_count?: number
  notes?: string
  incidents?: string
  created_at: string
  updated_at: string
}

export interface CreateVoyageLogRequest {
  rental_id: string
  log_date: string
  departure_port?: string
  arrival_port?: string
  departure_time?: string
  arrival_time?: string
  distance?: number
  engine_hours?: number
  fuel_consumed?: number
  fuel_unit?: string
  avg_speed?: number
  max_speed?: number
  weather: WeatherCondition
  wind_speed?: number
  wind_direction?: string
  wave_height?: number
  water_temperature?: number
  air_temperature?: number
  passenger_count?: number
  crew_count?: number
  notes?: string
  incidents?: string
}
