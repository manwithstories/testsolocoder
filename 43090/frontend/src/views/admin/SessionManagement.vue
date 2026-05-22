<template>
  <div class="session-management">
    <div class="header">
      <h3>拍卖会管理</h3>
      <el-button type="primary" @click="showCreateDialog = true">创建拍卖会</el-button>
    </div>
    <el-table :data="sessions" v-loading="loading">
      <el-table-column prop="name" label="名称" min-width="200" />
      <el-table-column prop="min_increment" label="最小加价" width="120">
        <template #default="{ row }">¥{{ row.min_increment.toFixed(2) }}</template>
      </el-table-column>
      <el-table-column prop="extend_time" label="延时(秒)" width="100" />
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="statusType(row.status)" size="small">{{ statusText(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="start_time" label="开始时间" width="160">
        <template #default="{ row }">{{ formatTime(row.start_time) }}</template>
      </el-table-column>
      <el-table-column prop="end_time" label="结束时间" width="160">
        <template #default="{ row }">{{ formatTime(row.end_time) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" size="small" @click="editSession(row)">编辑</el-button>
          <el-button
            v-if="row.status === 0"
            link
            type="success"
            size="small"
            @click="startSession(row)"
          >开始</el-button>
          <el-button
            v-if="row.status === 1"
            link
            type="warning"
            size="small"
            @click="endSession(row)"
          >结束</el-button>
          <el-button
            v-if="row.status !== 2"
            link
            type="danger"
            size="small"
            @click="cancelSession(row)"
          >取消</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination
      class="pagination"
      v-model:current-page="page"
      v-model:page-size="pageSize"
      :total="total"
      layout="total, prev, pager, next"
      @current-change="fetchSessions"
    />

    <el-dialog v-model="showCreateDialog" :title="editingSession ? '编辑拍卖会' : '创建拍卖会'" width="500px">
      <el-form :model="sessionForm" label-width="100px">
        <el-form-item label="名称">
          <el-input v-model="sessionForm.name" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="sessionForm.description" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item label="开始时间">
          <el-date-picker
            v-model="sessionForm.start_time"
            type="datetime"
            placeholder="选择开始时间"
            value-format="YYYY-MM-DD HH:mm:ss"
          />
        </el-form-item>
        <el-form-item label="结束时间">
          <el-date-picker
            v-model="sessionForm.end_time"
            type="datetime"
            placeholder="选择结束时间"
            value-format="YYYY-MM-DD HH:mm:ss"
          />
        </el-form-item>
        <el-form-item label="最小加价">
          <el-input-number v-model="sessionForm.min_increment" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="延时时间(秒)">
          <el-input-number v-model="sessionForm.extend_time" :min="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" @click="saveSession">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import dayjs from 'dayjs'
import type { AuctionSession } from '@/types'
import { sessionApi } from '@/api'

const sessions = ref<AuctionSession[]>([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const showCreateDialog = ref(false)
const editingSession = ref<AuctionSession | null>(null)

const sessionForm = reactive({
  name: '',
  description: '',
  start_time: '',
  end_time: '',
  min_increment: 10,
  extend_time: 300,
})

const fetchSessions = async () => {
  loading.value = true
  try {
    const res = await sessionApi.getList({ page: page.value, page_size: pageSize.value })
    sessions.value = res.list
    total.value = res.total
  } finally {
    loading.value = false
  }
}

const statusType = (status: number) => {
  const map: Record<number, string> = { 0: 'info', 1: 'success', 2: 'info', 3: 'danger' }
  return map[status] || 'info'
}

const statusText = (status: number) => {
  const map: Record<number, string> = { 0: '未开始', 1: '进行中', 2: '已结束', 3: '已取消' }
  return map[status] || ''
}

const formatTime = (time: string) => dayjs(time).format('YYYY-MM-DD HH:mm')

const editSession = (session: AuctionSession) => {
  editingSession.value = session
  Object.assign(sessionForm, {
    name: session.name,
    description: session.description,
    start_time: session.start_time,
    end_time: session.end_time,
    min_increment: session.min_increment,
    extend_time: session.extend_time,
  })
  showCreateDialog.value = true
}

const saveSession = async () => {
  try {
    if (editingSession.value) {
      await sessionApi.update(editingSession.value.id, sessionForm)
      ElMessage.success('更新成功')
    } else {
      await sessionApi.create(sessionForm)
      ElMessage.success('创建成功')
    }
    showCreateDialog.value = false
    fetchSessions()
  } catch (e) {}
}

const startSession = async (session: AuctionSession) => {
  try {
    await sessionApi.start(session.id)
    ElMessage.success('已开始')
    fetchSessions()
  } catch (e) {}
}

const endSession = async (session: AuctionSession) => {
  try {
    await sessionApi.end(session.id)
    ElMessage.success('已结束')
    fetchSessions()
  } catch (e) {}
}

const cancelSession = async (session: AuctionSession) => {
  try {
    await ElMessageBox.confirm('确定要取消这个拍卖会吗？', '提示', { type: 'warning' })
    await sessionApi.cancel(session.id)
    ElMessage.success('已取消')
    fetchSessions()
  } catch (e) {}
}

onMounted(() => {
  fetchSessions()
})
</script>

<style scoped>
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.header h3 {
  margin: 0;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}
</style>
