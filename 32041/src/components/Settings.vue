<script setup lang="ts">
import { ref, watch } from 'vue'
import { useStore } from '../stores/useStore'
import { playBeep } from '../utils/audio'

const store = useStore()

const localSettings = ref({ ...store.settings.value })

watch(
  () => store.settings.value,
  (newSettings) => {
    localSettings.value = { ...newSettings }
  },
  { deep: true }
)

function updateWorkDuration(value: number) {
  localSettings.value.workDuration = Math.max(1, Math.min(60, value))
  store.updateSettings({ workDuration: localSettings.value.workDuration })
}

function updateShortBreakDuration(value: number) {
  localSettings.value.shortBreakDuration = Math.max(1, Math.min(30, value))
  store.updateSettings({ shortBreakDuration: localSettings.value.shortBreakDuration })
}

function updateLongBreakDuration(value: number) {
  localSettings.value.longBreakDuration = Math.max(1, Math.min(60, value))
  store.updateSettings({ longBreakDuration: localSettings.value.longBreakDuration })
}

function updatePomodorosBeforeLongBreak(value: number) {
  localSettings.value.pomodorosBeforeLongBreak = Math.max(2, Math.min(10, value))
  store.updateSettings({ pomodorosBeforeLongBreak: localSettings.value.pomodorosBeforeLongBreak })
}

function toggleSound() {
  localSettings.value.soundEnabled = !localSettings.value.soundEnabled
  store.updateSettings({ soundEnabled: localSettings.value.soundEnabled })
  
  if (localSettings.value.soundEnabled) {
    playBeep()
  }
}

function toggleAutoStartBreak() {
  localSettings.value.autoStartBreak = !localSettings.value.autoStartBreak
  store.updateSettings({ autoStartBreak: localSettings.value.autoStartBreak })
}

function toggleAutoStartWork() {
  localSettings.value.autoStartWork = !localSettings.value.autoStartWork
  store.updateSettings({ autoStartWork: localSettings.value.autoStartWork })
}

function resetToDefaults() {
  store.updateSettings({
    workDuration: 25,
    shortBreakDuration: 5,
    longBreakDuration: 15,
    pomodorosBeforeLongBreak: 4,
    soundEnabled: true,
    autoStartBreak: false,
    autoStartWork: false
  })
}
</script>

<template>
  <div class="settings">
    <h3>设置</h3>
    
    <div class="settings-section">
      <h4>时长设置</h4>
      
      <div class="setting-item">
        <div class="setting-info">
          <span class="setting-label">工作时长</span>
          <span class="setting-desc">每个番茄钟的专注时间</span>
        </div>
        <div class="setting-control">
          <button 
            class="control-btn" 
            @click="updateWorkDuration(localSettings.workDuration - 1)"
          >
            -
          </button>
          <span class="control-value">{{ localSettings.workDuration }} 分钟</span>
          <button 
            class="control-btn" 
            @click="updateWorkDuration(localSettings.workDuration + 1)"
          >
            +
          </button>
        </div>
      </div>

      <div class="setting-item">
        <div class="setting-info">
          <span class="setting-label">短休息时长</span>
          <span class="setting-desc">工作后的短暂休息</span>
        </div>
        <div class="setting-control">
          <button 
            class="control-btn" 
            @click="updateShortBreakDuration(localSettings.shortBreakDuration - 1)"
          >
            -
          </button>
          <span class="control-value">{{ localSettings.shortBreakDuration }} 分钟</span>
          <button 
            class="control-btn" 
            @click="updateShortBreakDuration(localSettings.shortBreakDuration + 1)"
          >
            +
          </button>
        </div>
      </div>

      <div class="setting-item">
        <div class="setting-info">
          <span class="setting-label">长休息时长</span>
          <span class="setting-desc">完成一轮后的长时间休息</span>
        </div>
        <div class="setting-control">
          <button 
            class="control-btn" 
            @click="updateLongBreakDuration(localSettings.longBreakDuration - 1)"
          >
            -
          </button>
          <span class="control-value">{{ localSettings.longBreakDuration }} 分钟</span>
          <button 
            class="control-btn" 
            @click="updateLongBreakDuration(localSettings.longBreakDuration + 1)"
          >
            +
          </button>
        </div>
      </div>

      <div class="setting-item">
        <div class="setting-info">
          <span class="setting-label">长休息间隔</span>
          <span class="setting-desc">多少个番茄钟后进入长休息</span>
        </div>
        <div class="setting-control">
          <button 
            class="control-btn" 
            @click="updatePomodorosBeforeLongBreak(localSettings.pomodorosBeforeLongBreak - 1)"
          >
            -
          </button>
          <span class="control-value">{{ localSettings.pomodorosBeforeLongBreak }} 个</span>
          <button 
            class="control-btn" 
            @click="updatePomodorosBeforeLongBreak(localSettings.pomodorosBeforeLongBreak + 1)"
          >
            +
          </button>
        </div>
      </div>
    </div>

    <div class="settings-section">
      <h4>通知设置</h4>
      
      <div class="setting-item">
        <div class="setting-info">
          <span class="setting-label">提示音</span>
          <span class="setting-desc">计时结束时播放提示音</span>
        </div>
        <button 
          class="toggle-btn"
          :class="{ active: localSettings.soundEnabled }"
          @click="toggleSound"
        >
          <span class="toggle-track">
            <span class="toggle-thumb"></span>
          </span>
        </button>
      </div>
    </div>

    <div class="settings-section">
      <h4>自动化</h4>
      
      <div class="setting-item">
        <div class="setting-info">
          <span class="setting-label">自动开始休息</span>
          <span class="setting-desc">工作结束后自动开始休息</span>
        </div>
        <button 
          class="toggle-btn"
          :class="{ active: localSettings.autoStartBreak }"
          @click="toggleAutoStartBreak"
        >
          <span class="toggle-track">
            <span class="toggle-thumb"></span>
          </span>
        </button>
      </div>

      <div class="setting-item">
        <div class="setting-info">
          <span class="setting-label">自动开始工作</span>
          <span class="setting-desc">休息结束后自动开始工作</span>
        </div>
        <button 
          class="toggle-btn"
          :class="{ active: localSettings.autoStartWork }"
          @click="toggleAutoStartWork"
        >
          <span class="toggle-track">
            <span class="toggle-thumb"></span>
          </span>
        </button>
      </div>
    </div>

    <div class="settings-actions">
      <button class="reset-btn" @click="resetToDefaults">
        恢复默认设置
      </button>
    </div>
  </div>
</template>

<style scoped>
.settings {
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow-y: auto;
}

.settings h3 {
  margin: 0 0 20px 0;
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.settings-section {
  margin-bottom: 24px;
}

.settings-section h4 {
  margin: 0 0 12px 0;
  font-size: 14px;
  font-weight: 600;
  color: #666;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.setting-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 0;
  border-bottom: 1px solid #f0f0f0;
}

.setting-item:last-child {
  border-bottom: none;
}

.setting-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.setting-label {
  font-size: 14px;
  font-weight: 500;
  color: #333;
}

.setting-desc {
  font-size: 12px;
  color: #999;
}

.setting-control {
  display: flex;
  align-items: center;
  gap: 12px;
}

.control-btn {
  width: 32px;
  height: 32px;
  border: 2px solid #e0e0e0;
  background: white;
  border-radius: 8px;
  cursor: pointer;
  font-size: 18px;
  font-weight: 500;
  color: #666;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.control-btn:hover {
  border-color: #e74c3c;
  color: #e74c3c;
}

.control-value {
  min-width: 70px;
  text-align: center;
  font-size: 14px;
  font-weight: 500;
  color: #333;
  font-variant-numeric: tabular-nums;
}

.toggle-btn {
  background: none;
  border: none;
  padding: 0;
  cursor: pointer;
}

.toggle-track {
  display: block;
  width: 48px;
  height: 28px;
  background: #e0e0e0;
  border-radius: 14px;
  position: relative;
  transition: background 0.2s;
}

.toggle-btn.active .toggle-track {
  background: #e74c3c;
}

.toggle-thumb {
  position: absolute;
  top: 2px;
  left: 2px;
  width: 24px;
  height: 24px;
  background: white;
  border-radius: 50%;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  transition: transform 0.2s;
}

.toggle-btn.active .toggle-thumb {
  transform: translateX(20px);
}

.settings-actions {
  margin-top: auto;
  padding-top: 20px;
  border-top: 1px solid #f0f0f0;
}

.reset-btn {
  width: 100%;
  padding: 12px;
  border: 2px solid #e0e0e0;
  background: white;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  color: #666;
  transition: all 0.2s;
}

.reset-btn:hover {
  border-color: #e74c3c;
  color: #e74c3c;
}
</style>
