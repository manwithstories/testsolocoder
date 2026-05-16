<script setup lang="ts">
import { ref } from 'vue'
import { useStore } from '../stores/useStore'
import TaskItem from './TaskItem.vue'

const store = useStore()
const newTaskTitle = ref('')
const showCompleted = ref(false)

function handleAddTask() {
  if (newTaskTitle.value.trim()) {
    store.addTask(newTaskTitle.value)
    newTaskTitle.value = ''
  }
}

function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Enter') {
    handleAddTask()
  }
}
</script>

<template>
  <div class="task-list">
    <div class="task-header">
      <h3>今日任务</h3>
      <span class="task-count">{{ store.tasks.value.length }} 个任务</span>
    </div>
    
    <div class="add-task">
      <input
        v-model="newTaskTitle"
        type="text"
        placeholder="添加新任务..."
        class="task-input"
        @keydown="handleKeydown"
      />
      <button class="add-btn" @click="handleAddTask" :disabled="!newTaskTitle.trim()">
        <svg width="20" height="20" viewBox="0 0 20 20" fill="none">
          <path d="M10 4V16M4 10H16" stroke="white" stroke-width="2" stroke-linecap="round"/>
        </svg>
      </button>
    </div>

    <div class="tasks-container">
      <div class="tasks-list">
        <TaskItem
          v-for="task in store.tasks.value"
          :key="task.id"
          :task="task"
          :is-selected="store.currentTaskId.value === task.id"
          @toggle="store.toggleTaskComplete"
          @select="store.selectTask"
          @delete="store.deleteTask"
        />
        
        <div v-if="store.tasks.value.length === 0" class="empty-state">
          <p>还没有任务，添加一个开始专注吧！</p>
        </div>
      </div>

      <div v-if="store.completedTasks.value.length > 0" class="completed-section">
        <button 
          class="completed-toggle"
          @click="showCompleted = !showCompleted"
        >
          <span>已完成 ({{ store.completedTasks.value.length }})</span>
          <svg 
            width="16" 
            height="16" 
            viewBox="0 0 16 16" 
            fill="none"
            :style="{ transform: showCompleted ? 'rotate(180deg)' : 'rotate(0)', transition: 'transform 0.2s' }"
          >
            <path d="M4 6L8 10L12 6" stroke="#666" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </button>
        
        <div v-show="showCompleted" class="completed-list">
          <TaskItem
            v-for="task in store.completedTasks.value"
            :key="task.id"
            :task="task"
            :is-selected="store.currentTaskId.value === task.id"
            @toggle="store.toggleTaskComplete"
            @select="store.selectTask"
            @delete="store.deleteTask"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.task-list {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.task-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.task-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.task-count {
  font-size: 13px;
  color: #999;
}

.add-task {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
}

.task-input {
  flex: 1;
  padding: 10px 14px;
  border: 2px solid #e0e0e0;
  border-radius: 8px;
  font-size: 14px;
  transition: border-color 0.2s;
  outline: none;
}

.task-input:focus {
  border-color: #e74c3c;
}

.add-btn {
  width: 40px;
  height: 40px;
  border: none;
  background: #e74c3c;
  color: white;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.2s;
}

.add-btn:hover:not(:disabled) {
  background: #c0392b;
}

.add-btn:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.tasks-container {
  flex: 1;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.tasks-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.empty-state {
  text-align: center;
  padding: 32px 16px;
  color: #999;
}

.empty-state p {
  margin: 0;
  font-size: 14px;
}

.completed-section {
  margin-top: 8px;
}

.completed-toggle {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 12px;
  background: #f5f5f5;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 13px;
  color: #666;
  transition: background 0.2s;
}

.completed-toggle:hover {
  background: #eee;
}

.completed-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-top: 8px;
}
</style>
