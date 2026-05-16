<script setup lang="ts">
import type { Task } from '../types'

interface Props {
  task: Task
  isSelected: boolean
}

defineProps<Props>()

const emit = defineEmits<{
  toggle: [id: string]
  select: [id: string]
  delete: [id: string]
}>()
</script>

<template>
  <div 
    class="task-item"
    :class="{ completed: task.completed, selected: isSelected }"
    @click="emit('select', task.id)"
  >
    <button
      class="checkbox"
      :class="{ checked: task.completed }"
      @click.stop="emit('toggle', task.id)"
    >
      <svg v-if="task.completed" width="14" height="14" viewBox="0 0 14 14" fill="none">
        <path d="M3 7L6 10L11 4" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>
    </button>
    
    <div class="task-content">
      <span class="task-title">{{ task.title }}</span>
      <div class="task-meta">
        <span class="pomodoro-count">
          🍅 {{ task.completedPomodoros }}
        </span>
      </div>
    </div>
    
    <button 
      class="delete-btn"
      @click.stop="emit('delete', task.id)"
      title="删除任务"
    >
      <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
        <path d="M4 4L12 12M12 4L4 12" stroke="#999" stroke-width="2" stroke-linecap="round"/>
      </svg>
    </button>
  </div>
</template>

<style scoped>
.task-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: white;
  border-radius: 8px;
  border: 2px solid transparent;
  cursor: pointer;
  transition: all 0.2s;
}

.task-item:hover {
  background: #fafafa;
}

.task-item.selected {
  border-color: #e74c3c;
  background: #fef5f5;
}

.task-item.completed .task-title {
  text-decoration: line-through;
  color: #999;
}

.checkbox {
  width: 22px;
  height: 22px;
  border-radius: 50%;
  border: 2px solid #ddd;
  background: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: all 0.2s;
}

.checkbox:hover {
  border-color: #e74c3c;
}

.checkbox.checked {
  background: #27ae60;
  border-color: #27ae60;
}

.task-content {
  flex: 1;
  min-width: 0;
}

.task-title {
  font-size: 14px;
  color: #333;
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.task-meta {
  margin-top: 4px;
}

.pomodoro-count {
  font-size: 12px;
  color: #e74c3c;
}

.delete-btn {
  padding: 4px;
  border: none;
  background: transparent;
  cursor: pointer;
  opacity: 0;
  transition: opacity 0.2s;
  border-radius: 4px;
}

.task-item:hover .delete-btn {
  opacity: 1;
}

.delete-btn:hover {
  background: #ffe5e5;
}
</style>
