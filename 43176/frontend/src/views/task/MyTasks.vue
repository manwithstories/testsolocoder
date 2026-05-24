<template>
  <div class="my-tasks-container">
    <el-tabs v-model="activeTab" class="tabs">
      <el-tab-pane label="我发布的" name="published">
        <TaskListCard
          :tasks="publishedTasks"
          :loading="loading"
          :show-accept="false"
        />
      </el-tab-pane>
      <el-tab-pane label="我接取的" name="accepted">
        <TaskListCard
          :tasks="acceptedTasks"
          :loading="loading"
          :show-accept="false"
        />
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useUserStore } from '@/stores/user'
import { taskApi, orderApi } from '@/api'
import TaskListCard from '@/components/TaskListCard.vue'
import type { Task } from '@/types'

const userStore = useUserStore()
const loading = ref(false)
const activeTab = ref('published')
const publishedTasks = ref<Task[]>([])
const acceptedTasks = ref<Task[]>([])

const fetchPublishedTasks = async () => {
  loading.value = true
  try {
    const res = await taskApi.list({ page_size: 50 })
    if (res.code === 200) {
      publishedTasks.value = res.data.items.filter(
        (task: Task) => task.publisher_id === userStore.userInfo?.id
      )
    }
  } catch (error) {
    console.error('Failed to fetch published tasks:', error)
  } finally {
    loading.value = false
  }
}

const fetchAcceptedTasks = async () => {
  loading.value = true
  try {
    const res = await orderApi.list({ page_size: 50 })
    if (res.code === 200) {
      acceptedTasks.value = res.data.items
        .filter((order: any) => order.courier_id === userStore.userInfo?.id)
        .map((order: any) => order.task)
        .filter(Boolean)
    }
  } catch (error) {
    console.error('Failed to fetch accepted tasks:', error)
  } finally {
    loading.value = false
  }
}

watch(activeTab, (tab) => {
  if (tab === 'published') {
    fetchPublishedTasks()
  } else {
    fetchAcceptedTasks()
  }
})

onMounted(() => {
  fetchPublishedTasks()
})
</script>

<style lang="scss" scoped>
.my-tasks-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;

  .tabs {
    background: #fff;
    border-radius: 8px;
    padding: 16px;
  }
}
</style>
