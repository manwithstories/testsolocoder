import type { CityWeather } from '../types'

const getCurrentTime = (): string => {
  const now = new Date()
  return `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')} ${String(now.getHours()).padStart(2, '0')}:${String(now.getMinutes()).padStart(2, '0')}`
}

const getNextDays = (count: number): Array<{ date: string; dayName: string }> => {
  const days = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
  const result = []
  const today = new Date()
  
  for (let i = 0; i < count; i++) {
    const date = new Date(today)
    date.setDate(today.getDate() + i)
    result.push({
      date: `${date.getMonth() + 1}/${date.getDate()}`,
      dayName: i === 0 ? '今天' : days[date.getDay()]
    })
  }
  
  return result
}

export const weatherMockData: CityWeather[] = [
  {
    city: '北京',
    current: {
      city: '北京',
      temperature: 28,
      humidity: 45,
      windSpeed: 3.2,
      condition: {
        main: 'sunny',
        description: '晴天',
        icon: '☀️'
      },
      updateTime: getCurrentTime()
    },
    forecast: getNextDays(7).map((day, index) => ({
      ...day,
      highTemp: 30 - index,
      lowTemp: 18 - index,
      condition: {
        main: index % 3 === 0 ? 'sunny' : index % 3 === 1 ? 'cloudy' : 'rainy' as const,
        description: index % 3 === 0 ? '晴天' : index % 3 === 1 ? '多云' : '小雨',
        icon: index % 3 === 0 ? '☀️' : index % 3 === 1 ? '⛅' : '🌧️'
      }
    }))
  },
  {
    city: '上海',
    current: {
      city: '上海',
      temperature: 32,
      humidity: 78,
      windSpeed: 2.5,
      condition: {
        main: 'cloudy',
        description: '多云',
        icon: '⛅'
      },
      updateTime: getCurrentTime()
    },
    forecast: getNextDays(7).map((day, index) => ({
      ...day,
      highTemp: 34 - index,
      lowTemp: 26 - index,
      condition: {
        main: index % 2 === 0 ? 'cloudy' : 'rainy' as const,
        description: index % 2 === 0 ? '多云' : '阵雨',
        icon: index % 2 === 0 ? '⛅' : '🌧️'
      }
    }))
  },
  {
    city: '广州',
    current: {
      city: '广州',
      temperature: 35,
      humidity: 85,
      windSpeed: 1.8,
      condition: {
        main: 'rainy',
        description: '雷阵雨',
        icon: '⛈️'
      },
      updateTime: getCurrentTime()
    },
    forecast: getNextDays(7).map((day, index) => ({
      ...day,
      highTemp: 36 - index,
      lowTemp: 28 - index,
      condition: {
        main: index % 4 === 0 ? 'rainy' : index % 4 === 1 ? 'sunny' : 'cloudy' as const,
        description: index % 4 === 0 ? '雷阵雨' : index % 4 === 1 ? '晴天' : '多云',
        icon: index % 4 === 0 ? '⛈️' : index % 4 === 1 ? '☀️' : '⛅'
      }
    }))
  },
  {
    city: '深圳',
    current: {
      city: '深圳',
      temperature: 33,
      humidity: 80,
      windSpeed: 2.2,
      condition: {
        main: 'sunny',
        description: '晴转多云',
        icon: '🌤️'
      },
      updateTime: getCurrentTime()
    },
    forecast: getNextDays(7).map((day, index) => ({
      ...day,
      highTemp: 35 - index,
      lowTemp: 27 - index,
      condition: {
        main: index % 2 === 0 ? 'sunny' : 'cloudy' as const,
        description: index % 2 === 0 ? '晴天' : '多云',
        icon: index % 2 === 0 ? '☀️' : '⛅'
      }
    }))
  },
  {
    city: '杭州',
    current: {
      city: '杭州',
      temperature: 26,
      humidity: 70,
      windSpeed: 2.8,
      condition: {
        main: 'rainy',
        description: '小雨',
        icon: '🌧️'
      },
      updateTime: getCurrentTime()
    },
    forecast: getNextDays(7).map((day, index) => ({
      ...day,
      highTemp: 28 - index,
      lowTemp: 20 - index,
      condition: {
        main: index % 3 === 0 ? 'rainy' : index % 3 === 1 ? 'cloudy' : 'sunny' as const,
        description: index % 3 === 0 ? '小雨' : index % 3 === 1 ? '多云' : '晴天',
        icon: index % 3 === 0 ? '🌧️' : index % 3 === 1 ? '⛅' : '☀️'
      }
    }))
  },
  {
    city: '成都',
    current: {
      city: '成都',
      temperature: 22,
      humidity: 75,
      windSpeed: 1.5,
      condition: {
        main: 'cloudy',
        description: '阴天',
        icon: '☁️'
      },
      updateTime: getCurrentTime()
    },
    forecast: getNextDays(7).map((day, index) => ({
      ...day,
      highTemp: 24 - index,
      lowTemp: 16 - index,
      condition: {
        main: index % 4 === 0 ? 'cloudy' : index % 4 === 1 ? 'rainy' : 'sunny' as const,
        description: index % 4 === 0 ? '阴天' : index % 4 === 1 ? '小雨' : '晴天',
        icon: index % 4 === 0 ? '☁️' : index % 4 === 1 ? '🌧️' : '☀️'
      }
    }))
  },
  {
    city: '武汉',
    current: {
      city: '武汉',
      temperature: 34,
      humidity: 65,
      windSpeed: 2.0,
      condition: {
        main: 'sunny',
        description: '晴天',
        icon: '☀️'
      },
      updateTime: getCurrentTime()
    },
    forecast: getNextDays(7).map((day, index) => ({
      ...day,
      highTemp: 36 - index,
      lowTemp: 28 - index,
      condition: {
        main: index % 2 === 0 ? 'sunny' : 'rainy' as const,
        description: index % 2 === 0 ? '晴天' : '雷阵雨',
        icon: index % 2 === 0 ? '☀️' : '⛈️'
      }
    }))
  },
  {
    city: '南京',
    current: {
      city: '南京',
      temperature: 30,
      humidity: 68,
      windSpeed: 2.3,
      condition: {
        main: 'cloudy',
        description: '多云',
        icon: '⛅'
      },
      updateTime: getCurrentTime()
    },
    forecast: getNextDays(7).map((day, index) => ({
      ...day,
      highTemp: 32 - index,
      lowTemp: 24 - index,
      condition: {
        main: index % 3 === 0 ? 'cloudy' : index % 3 === 1 ? 'sunny' : 'rainy' as const,
        description: index % 3 === 0 ? '多云' : index % 3 === 1 ? '晴天' : '小雨',
        icon: index % 3 === 0 ? '⛅' : index % 3 === 1 ? '☀️' : '🌧️'
      }
    }))
  },
  {
    city: '西安',
    current: {
      city: '西安',
      temperature: 25,
      humidity: 50,
      windSpeed: 3.0,
      condition: {
        main: 'sunny',
        description: '晴天',
        icon: '☀️'
      },
      updateTime: getCurrentTime()
    },
    forecast: getNextDays(7).map((day, index) => ({
      ...day,
      highTemp: 27 - index,
      lowTemp: 15 - index,
      condition: {
        main: index % 3 === 0 ? 'sunny' : 'cloudy' as const,
        description: index % 3 === 0 ? '晴天' : '多云',
        icon: index % 3 === 0 ? '☀️' : '⛅'
      }
    }))
  },
  {
    city: '重庆',
    current: {
      city: '重庆',
      temperature: 36,
      humidity: 72,
      windSpeed: 1.2,
      condition: {
        main: 'sunny',
        description: '炎热',
        icon: '🌞'
      },
      updateTime: getCurrentTime()
    },
    forecast: getNextDays(7).map((day, index) => ({
      ...day,
      highTemp: 38 - index,
      lowTemp: 30 - index,
      condition: {
        main: index % 2 === 0 ? 'sunny' : 'rainy' as const,
        description: index % 2 === 0 ? '晴天' : '阵雨',
        icon: index % 2 === 0 ? '🌞' : '🌧️'
      }
    }))
  }
]

export const getWeatherByCity = (city: string): CityWeather | undefined => {
  return weatherMockData.find(item => item.city === city)
}

export const searchCities = (keyword: string): string[] => {
  if (!keyword) return []
  const lowerKeyword = keyword.toLowerCase()
  return weatherMockData
    .map(item => item.city)
    .filter(city => city.toLowerCase().includes(lowerKeyword))
}