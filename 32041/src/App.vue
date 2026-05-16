<script setup lang="ts">
import { ref } from 'vue'
import Timer from './components/Timer.vue'
import TaskList from './components/TaskList.vue'
import Statistics from './components/Statistics.vue'
import Settings from './components/Settings.vue'

type Tab = 'tasks' | 'stats' | 'settings'

const activeTab = ref<Tab>('tasks')

const tabs: { id: Tab; label: string; icon: string }[] = [
  { id: 'tasks', label: '任务', icon: '📋' },
  { id: 'stats', label: '统计', icon: '📊' },
  { id: 'settings', label: '设置', icon: '⚙️' }
]
</script>

<template>
  <div class="app-container">
    <header class="app-header">
      <div class="logo">
        <span class="logo-icon">🍅</span>
        <h1 class="logo-text">专注计时器</h1>
      </div>
    </header>

    <main class="app-main">
      <div class="timer-section">
        <Timer />
      </div>

      <div class="sidebar">
        <nav class="tab-nav">
          <button
            v-for="tab in tabs"
            :key="tab.id"
            class="tab-btn"
            :class="{ active: activeTab === tab.id }"
            @click="activeTab = tab.id"
          >
            <span class="tab-icon">{{ tab.icon }}</span>
            <span class="tab-label">{{ tab.label }}</span>
          </button>
        </nav>

        <div class="tab-content">
          <TaskList v-show="activeTab === 'tasks'" />
          <Statistics v-show="activeTab === 'stats'" />
          <Settings v-show="activeTab === 'settings'" />
        </div>
      </div>
    </main>

    <footer class="app-footer">
      <p>专注当下，高效工作 🍅</p>
    </footer>
  </div>
</template>

<style scoped>
.app-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: linear-gradient(135deg, #ff9a9e 0%, #fecfef 100%);
}

.app-header {
  padding: 20px 32px;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo-icon {
  font-size: 32px;
}

.logo-text {
  margin: 0;
  font-size: 24px;
  font-weight: 700;
  color: #333;
}

.app-main {
  flex: 1;
  display: grid;
  grid-template-columns: 1fr 400px;
  gap: 32px;
  padding: 32px;
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
  box-sizing: border-box;
}

.timer-section {
  display: flex;
  align-items: center;
  justify-content: center;
}

.sidebar {
  display: flex;
  flex-direction: column;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  min-height: 600px;
}

.tab-nav {
  display: flex;
  border-bottom: 1px solid #f0f0f0;
}

.tab-btn {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 12px 16px;
  border: none;
  background: transparent;
  cursor: pointer;
  transition: all 0.2s;
  position: relative;
}

.tab-btn::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 0;
  height: 3px;
  background: #e74c3c;
  border-radius: 2px;
  transition: width 0.2s;
}

.tab-btn:hover {
  background: #fafafa;
}

.tab-btn.active::after {
  width: 60%;
}

.tab-btn.active .tab-label {
  color: #e74c3c;
  font-weight: 600;
}

.tab-icon {
  font-size: 20px;
}

.tab-label {
  font-size: 12px;
  color: #666;
  transition: color 0.2s;
}

.tab-content {
  flex: 1;
  padding: 20px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.app-footer {
  padding: 16px;
  text-align: center;
  color: rgba(255, 255, 255, 0.8);
  font-size: 14px;
}

.app-footer p {
  margin: 0;
}

@media (max-width: 900px) {
  .app-main {
    grid-template-columns: 1fr;
    padding: 16px;
  }

  .sidebar {
    min-height: auto;
    max-height: 500px;
  }

  .app-header {
    padding: 16px 20px;
  }

  .logo-text {
    font-size: 20px;
  }
}
</style>
