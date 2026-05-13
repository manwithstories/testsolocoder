<script setup lang="ts">
import type { CurrentWeather } from '../types'

defineProps<{
  weather: CurrentWeather
}>()
</script>

<template>
  <div class="current-weather">
    <div class="weather-main">
      <div class="weather-icon">
        {{ weather.condition.icon }}
      </div>
      <div class="temperature">
        <span class="temp-value">{{ weather.temperature }}</span>
        <span class="temp-unit">°C</span>
      </div>
    </div>
    
    <div class="weather-info">
      <p class="condition">{{ weather.condition.description }}</p>
      <div class="details">
        <div class="detail-item">
          <span class="detail-icon">💧</span>
          <span>湿度 {{ weather.humidity }}%</span>
        </div>
        <div class="detail-item">
          <span class="detail-icon">🌬️</span>
          <span>风速 {{ weather.windSpeed }} m/s</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.current-weather {
  background: rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(10px);
  border-radius: 24px;
  padding: 40px;
  text-align: center;
  color: white;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}

.weather-main {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
  margin-bottom: 30px;
}

.weather-icon {
  font-size: 120px;
  line-height: 1;
}

.temperature {
  display: flex;
  align-items: flex-start;
}

.temp-value {
  font-size: 100px;
  font-weight: 300;
  line-height: 1;
}

.temp-unit {
  font-size: 32px;
  margin-top: 10px;
}

.weather-info {
  text-align: center;
}

.condition {
  font-size: 28px;
  margin-bottom: 24px;
  font-weight: 500;
}

.details {
  display: flex;
  justify-content: center;
  gap: 40px;
  flex-wrap: wrap;
}

.detail-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 18px;
}

.detail-icon {
  font-size: 24px;
}

@media (max-width: 640px) {
  .current-weather {
    padding: 30px 20px;
  }
  
  .weather-icon {
    font-size: 80px;
  }
  
  .temp-value {
    font-size: 72px;
  }
  
  .temp-unit {
    font-size: 24px;
    margin-top: 8px;
  }
  
  .condition {
    font-size: 22px;
  }
  
  .details {
    gap: 24px;
  }
  
  .detail-item {
    font-size: 16px;
  }
}
</style>