<script setup lang="ts">
import type { ForecastDay } from '../types'
import WeatherCard from './WeatherCard.vue'

defineProps<{
  forecast: ForecastDay[]
}>()
</script>

<template>
  <div class="forecast-section">
    <h2 class="section-title">未来7天预报</h2>
    <div class="forecast-list">
      <WeatherCard
        v-for="(day, index) in forecast"
        :key="index"
        :forecast="day"
      />
    </div>
  </div>
</template>

<style scoped>
.forecast-section {
  margin-top: 40px;
}

.section-title {
  color: white;
  font-size: 24px;
  font-weight: 600;
  margin-bottom: 24px;
  text-align: center;
}

.forecast-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
  gap: 16px;
  justify-items: center;
}

@media (max-width: 768px) {
  .forecast-list {
    grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
    gap: 12px;
  }
  
  .section-title {
    font-size: 20px;
  }
}

@media (max-width: 480px) {
  .forecast-list {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>