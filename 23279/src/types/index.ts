export interface WeatherCondition {
  main: 'sunny' | 'cloudy' | 'rainy' | 'snowy' | 'windy'
  description: string
  icon: string
}

export interface CurrentWeather {
  city: string
  temperature: number
  humidity: number
  windSpeed: number
  condition: WeatherCondition
  updateTime: string
}

export interface ForecastDay {
  date: string
  dayName: string
  highTemp: number
  lowTemp: number
  condition: WeatherCondition
}

export interface CityWeather {
  city: string
  current: CurrentWeather
  forecast: ForecastDay[]
}