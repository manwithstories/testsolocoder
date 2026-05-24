<template>
  <div class="page-container">
    <div class="page-header">
      <h2>通知消息</h2>
      <div>
        <el-button :icon="Check" @click="markAllRead">全部已读</el-button>
      </div>
    </div>

    <div class="search-bar">
      <el-select v-model="filterType" placeholder="消息类型" style="width: 160px" clearable @change="loadNotifications">
        <el-option label="告警" value="alert" />
        <el-option label="信息" value="info" />
        <el-option label="邀请" value="invitation" />
      </el-select>
      <el-select v-model="filterRead" placeholder="阅读状态" style="width: 140px" clearable @change="loadNotifications">
        <el-option label="未读" value="false" />
        <el-option label="已读" value="true" />
      </el-select>
    </div>

    <div class="table-container">
      <el-table :data="notifications" v-loading="loading" stripe>
        <el-table-column width="80">
          <template #default="{ row }">
            <el-icon :size="24" :color="getTypeColor(row.type)">
              <component :is="getTypeIcon(row.type)" />
            </el-icon>
          </template>
        </el-table-column>
        <el-table-column prop="title" label="标题" min-width="200">
          <template #default="{ row }">
            <span :style="{ fontWeight: row.isRead ? 'normal' : 'bold' }">{{ row.title }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="content" label="内容" show-overflow-tooltip />
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="getTagType(row.type)" size="small">{{ getTypeLabel(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="时间" width="180" />
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button v-if="!row.isRead" size="small" type="primary" link @click="markRead(row)">
              标为已读
            </el-button>
            <el-button size="small" type="danger" link @click="deleteItem(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-empty v-if="notifications.length === 0 && !loading" description="暂无通知" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Check, Warning, InfoFilled, Bell } from '@element-plus/icons-vue'
import {
  listNotifications, markNotificationRead, markAllRead as apiMarkAll,
  deleteNotification, notificationTypeMap
} from '@/api/notification'

const loading = ref(false)
const notifications = ref<any[]>([])
const filterType = ref('')
const filterRead = ref('')

onMounted(() => loadNotifications())

async function loadNotifications() {
  loading.value = true
  try {
    const res = await listNotifications({
      type: filterType.value || undefined,
      isRead: filterRead.value || undefined
    })
    notifications.value = res.notifications || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function markRead(row: any) {
  try {
    await markNotificationRead(row.id)
    row.isRead = true
    ElMessage.success('已标为已读')
  } catch (e) {
    console.error(e)
  }
}

async function markAllRead() {
  try {
    await apiMarkAll()
    notifications.value.forEach(n => n.isRead = true)
    ElMessage.success('已全部标为已读')
  } catch (e) {
    console.error(e)
  }
}

async function deleteItem(row: any) {
  try {
    await ElMessageBox.confirm('确定要删除该通知吗？', '提示', { type: 'warning' })
    await deleteNotification(row.id)
    ElMessage.success('删除成功')
    loadNotifications()
  } catch (e: any) {
    if (e !== 'cancel') console.error(e)
  }
}

function getTypeIcon(type: string) {
  const map: Record<string, any> = { alert: Warning, info: InfoFilled, invitation: Bell }
  return map[type] || InfoFilled
}

function getTypeColor(type: string) {
  const map: Record<string, string> = { alert: '#f56c6c', info: '#409eff', invitation: '#e6a23c' }
  return map[type] || '#909399'
}

function getTagType(type: string) {
  return notificationTypeMap[type]?.color || 'info'
}

function getTypeLabel(type: string) {
  return notificationTypeMap[type]?.label || type
}
</script>
