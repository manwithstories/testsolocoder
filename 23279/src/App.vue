<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import type { CityWeather } from './types'
import { getWeatherByCity, weatherMockData } from './mock/weatherData'
import SearchBar from './components/SearchBar.vue'
import CurrentWeather from './components/CurrentWeather.vue'
import ForecastList from './components/ForecastList.vue'

const currentCity = ref('北京')
const weatherData = ref<CityWeather | null>(null)

const loadWeatherData = (city: string) => {
  const data = getWeatherByCity(city)
  if (data) {
    weatherData.value = data
  }
}

const handleCityChange = (city: string) => {
  currentCity.value = city
  loadWeatherData(city)
}

const backgroundClass = computed(() => {
  if (!weatherData.value) return 'bg-sunny'
  
  const condition = weatherData.value.current.condition.main
  switch (condition) {
    case 'sunny':
      return 'bg-sunny'
    case 'cloudy':
      return 'bg-cloudy'
    case 'rainy':
      return 'bg-rainy'
    case 'snowy':
      return 'bg-snowy'
    default:
      return 'bg-sunny'
  }
})

onMounted(() => {
  loadWeatherData(currentCity.value)
})
</script>

<template>
  <div class="app-container" :class="backgroundClass">
    <div class="app-wrapper">
      <header class="app-header">
        <div class="header-left">
          <h1 class="city-name">
            <span class="location-icon">📍</span>
            {{ currentCity }}
          </h1>
          <p v-if="weatherData" class="update-time">
            更新时间: {{ weatherData.current.updateTime }}
          </p>
        </div>
        <SearchBar
          :current-city="currentCity"
          @city-change="handleCityChange"
        />
      </header>
      
      <main class="app-main">
        <CurrentWeather
          v-if="weatherData"
          :weather="weatherData.current"
        />
        
        <ForecastList
          v-if="weatherData"
          :forecast="weatherData.forecast"
        />
      </main>
      
      <footer class="app-footer">
        <p>数据来源: Mock 数据</p>
      </footer>
    </div>
  </div>
</template>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body, #app {
  height: 100%;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
}

.app-container {
  min-height: 100vh;
  transition: background 0.5s ease;
}

.bg-sunny {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.bg-cloudy {
  background: linear-gradient(135deg, #606c88 0%, #3f4c6b 100%);
}

.bg-rainy {
  background: linear-gradient(135deg, #2c3e50 0%, #4ca1af 100%);
}

.bg-snowy {
  background: linear-gradient(135deg, #e0eafc 0%, #cfdef3 100%);
}

.app-wrapper {
  max-width: 1200px;
  margin: 0 auto;
  padding: 30px 20px;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.app-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 40px;
  gap: 20px;
}

.header-left {
  color: white;
}

.city-name {
  font-size: 32px;
  font-weight: 700;
  display: flex;
  align-items: center;
  gap: 10px;
}

.location-icon {
  font-size: 28px;
}

.update-time {
  font-size: 14px;
  opacity: 0.8;
  margin-top: 8px;
}

.app-main {
  flex: 1;
}

.app-footer {
  text-align: center;
  color: rgba(255, 255, 255, 0.7);
  font-size: 14px;
  margin-top: 40px;
  padding-top: 20px;
  border-top: 1px solid rgba(255, 255, 255, 0.2);
}

@media (max-width: 768px) {
  .app-wrapper {
    padding: 20px 16px;
  }
  
  .app-header {
    flex-direction: column;
    align-items: stretch;
  }
  
  .city-name {
    font-size: 26px;
  }
  
  .location-icon {
    font-size: 24px;
  }
}

@media (max-width: 480px) {
  .city-name {
    font-size: 22px;
  }
  
  .location-icon {
    font-size: 20px;
  }
}
</style>