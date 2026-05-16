<script setup lang="ts">
import { computed, onMounted, onUnmounted } from 'vue'
import { useStore } from '../stores/useStore'
import { formatTime } from '../utils/storage'

const store = useStore()

const phaseLabels: Record<string, string> = {
  work: '专注时间',
  shortBreak: '短休息',
  longBreak: '长休息'
}

const totalTime = computed(() => {
  switch (store.timerPhase.value) {
    case 'work':
      return store.settings.value.workDuration * 60
    case 'shortBreak':
      return store.settings.value.shortBreakDuration * 60
    case 'longBreak':
      return store.settings.value.longBreakDuration * 60
  }
})

const progress = computed(() => {
  return ((totalTime.value - store.timeRemaining.value) / totalTime.value) * 100
})

const circumference = 2 * Math.PI * 140
const strokeDashoffset = computed(() => {
  return circumference - (progress.value / 100) * circumference
})

let intervalId: number | null = null

onMounted(() => {
  intervalId = window.setInterval(() => {
    store.tick()
  }, 1000)
})

onUnmounted(() => {
  if (intervalId) {
    clearInterval(intervalId)
  }
})

const phaseColors: Record<string, string> = {
  work: '#e74c3c',
  shortBreak: '#27ae60',
  longBreak: '#3498db'
}
</script>

<template>
  <div class="timer-container">
    <div class="phase-tabs">
      <button
        v-for="phase in ['work', 'shortBreak', 'longBreak'] as const"
        :key="phase"
        class="phase-tab"
        :class="{ active: store.timerPhase.value === phase }"
        @click="store.setPhase(phase)"
      >
        {{ phaseLabels[phase] }}
      </button>
    </div>

    <div class="timer-display">
      <svg class="progress-ring" width="320" height="320">
        <circle
          class="progress-ring-bg"
          cx="160"
          cy="160"
          r="140"
          fill="none"
          stroke="#e0e0e0"
          stroke-width="8"
        />
        <circle
          class="progress-ring-fg"
          cx="160"
          cy="160"
          r="140"
          fill="none"
          :stroke="phaseColors[store.timerPhase.value]"
          stroke-width="8"
          stroke-linecap="round"
          :stroke-dasharray="circumference"
          :stroke-dashoffset="strokeDashoffset"
          transform="rotate(-90 160 160)"
        />
      </svg>
      <div class="time-text">
        <span class="time">{{ formatTime(store.timeRemaining.value) }}</span>
        <span class="task-info" v-if="store.currentTask.value">
          当前任务: {{ store.currentTask.value.title }}
        </span>
      </div>
    </div>

    <div class="timer-controls">
      <template v-if="store.timerStatus.value === 'idle'">
        <button class="btn btn-primary" @click="store.startTimer()">
          开始
        </button>
      </template>
      <template v-else-if="store.timerStatus.value === 'running'">
        <button class="btn btn-secondary" @click="store.pauseTimer()">
          暂停
        </button>
        <button class="btn btn-outline" @click="store.skipPhase()">
          跳过
        </button>
      </template>
      <template v-else-if="store.timerStatus.value === 'paused'">
        <button class="btn btn-primary" @click="store.resumeTimer()">
          继续
        </button>
        <button class="btn btn-outline" @click="store.resetTimer()">
          重置
        </button>
      </template>
    </div>
  </div>
</template>

<style scoped>
.timer-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 24px;
}

.phase-tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 32px;
  background: #f5f5f5;
  padding: 4px;
  border-radius: 12px;
}

.phase-tab {
  padding: 8px 20px;
  border: none;
  background: transparent;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  color: #666;
  transition: all 0.2s;
}

.phase-tab:hover {
  color: #333;
}

.phase-tab.active {
  background: white;
  color: #333;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.timer-display {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 32px;
}

.progress-ring {
  transform: scale(1);
}

.progress-ring-fg {
  transition: stroke-dashoffset 0.5s ease;
}

.time-text {
  position: absolute;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.time {
  font-size: 64px;
  font-weight: 300;
  color: #333;
  font-variant-numeric: tabular-nums;
}

.task-info {
  font-size: 14px;
  color: #666;
  margin-top: 8px;
  max-width: 240px;
  text-align: center;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.timer-controls {
  display: flex;
  gap: 12px;
}

.btn {
  padding: 12px 32px;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary {
  background: #e74c3c;
  color: white;
}

.btn-primary:hover {
  background: #c0392b;
}

.btn-secondary {
  background: #f39c12;
  color: white;
}

.btn-secondary:hover {
  background: #d68910;
}

.btn-outline {
  background: transparent;
  border: 2px solid #ddd;
  color: #666;
}

.btn-outline:hover {
  border-color: #999;
  color: #333;
}
</style>
