<script setup lang="ts">
import type { ForecastDay } from '../types'

defineProps<{
  forecast: ForecastDay
}>()
</script>

<template>
  <div class="weather-card">
    <div class="card-header">
      <span class="day-name">{{ forecast.dayName }}</span>
      <span class="date">{{ forecast.date }}</span>
    </div>
    
    <div class="card-body">
      <div class="weather-icon">
        {{ forecast.condition.icon }}
      </div>
      <p class="condition">{{ forecast.condition.description }}</p>
    </div>
    
    <div class="card-footer">
      <span class="high-temp">{{ forecast.highTemp }}°</span>
      <span class="temp-divider">/</span>
      <span class="low-temp">{{ forecast.lowTemp }}°</span>
    </div>
  </div>
</template>

<style scoped>
.weather-card {
  background: rgba(255, 255, 255, 0.15);
  backdrop-filter: blur(10px);
  border-radius: 20px;
  padding: 20px;
  text-align: center;
  color: white;
  min-width: 140px;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.weather-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.2);
}

.card-header {
  display: flex;
  flex-direction: column;
  gap: 4px;
  margin-bottom: 16px;
}

.day-name {
  font-size: 18px;
  font-weight: 600;
}

.date {
  font-size: 14px;
  opacity: 0.8;
}

.card-body {
  margin-bottom: 16px;
}

.weather-icon {
  font-size: 50px;
  margin-bottom: 8px;
}

.condition {
  font-size: 14px;
  opacity: 0.9;
  margin: 0;
}

.card-footer {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.high-temp {
  font-size: 22px;
  font-weight: 600;
}

.temp-divider {
  opacity: 0.6;
  font-size: 18px;
}

.low-temp {
  font-size: 22px;
  opacity: 0.8;
}

@media (max-width: 768px) {
  .weather-card {
    min-width: 120px;
    padding: 16px;
  }
  
  .weather-icon {
    font-size: 42px;
  }
  
  .high-temp,
  .low-temp {
    font-size: 18px;
  }
}
</style>