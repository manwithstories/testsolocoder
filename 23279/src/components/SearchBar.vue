<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { searchCities } from '../mock/weatherData'

const props = defineProps<{
  currentCity: string
}>()

const emit = defineEmits<{
  (e: 'cityChange', city: string): void
}>()

const searchKeyword = ref('')
const showDropdown = ref(false)
const inputRef = ref<HTMLInputElement | null>(null)

const filteredCities = computed(() => {
  if (!searchKeyword.value) return []
  return searchCities(searchKeyword.value)
})

const handleInput = () => {
  showDropdown.value = filteredCities.value.length > 0
}

const selectCity = (city: string) => {
  emit('cityChange', city)
  searchKeyword.value = ''
  showDropdown.value = false
}

const handleFocus = () => {
  if (filteredCities.value.length > 0) {
    showDropdown.value = true
  }
}

const handleBlur = () => {
  setTimeout(() => {
    showDropdown.value = false
  }, 200)
}

watch(() => props.currentCity, (newCity) => {
  if (inputRef.value) {
    inputRef.value.placeholder = `当前城市: ${newCity}`
  }
}, { immediate: true })
</script>

<template>
  <div class="search-container">
    <div class="search-box">
      <input
        ref="inputRef"
        v-model="searchKeyword"
        type="text"
        class="search-input"
        placeholder="搜索城市..."
        @input="handleInput"
        @focus="handleFocus"
        @blur="handleBlur"
      />
      <span class="search-icon">🔍</span>
    </div>
    
    <div v-if="showDropdown && filteredCities.length > 0" class="dropdown">
      <div
        v-for="city in filteredCities"
        :key="city"
        class="dropdown-item"
        @mousedown="selectCity(city)"
      >
        📍 {{ city }}
      </div>
    </div>
  </div>
</template>

<style scoped>
.search-container {
  position: relative;
  width: 100%;
  max-width: 400px;
}

.search-box {
  position: relative;
  display: flex;
  align-items: center;
}

.search-input {
  width: 100%;
  padding: 12px 40px 12px 16px;
  border: none;
  border-radius: 12px;
  font-size: 16px;
  background: rgba(255, 255, 255, 0.9);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  outline: none;
  transition: all 0.3s ease;
}

.search-input:focus {
  background: white;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
}

.search-icon {
  position: absolute;
  right: 16px;
  font-size: 18px;
  pointer-events: none;
}

.dropdown {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  margin-top: 8px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
  overflow: hidden;
  z-index: 100;
}

.dropdown-item {
  padding: 12px 16px;
  cursor: pointer;
  transition: background 0.2s ease;
  font-size: 16px;
}

.dropdown-item:hover {
  background: #f0f4ff;
}
</style>