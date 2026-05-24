<template>
  <el-card>
    <template #header>
      <div class="header">
        <span>我的设备</span>
        <div>
          <el-button type="primary" @click="$router.push('/drone/create')">
            <el-icon><Plus /></el-icon>
            添加设备
          </el-button>
          <el-upload
            :action="'/api/drones/batch-import-excel'"
            :headers="{ Authorization: 'Bearer ' + token }"
            :show-file-list="false"
            accept=".xlsx,.xls"
            :on-success="handleBatchImport"
            :on-error="handleImportError"
          >
            <el-button type="success">
              <el-icon><Upload /></el-icon>
              批量导入
            </el-button>
          </el-upload>
        </div>
      </div>
    </template>

    <el-table :data="drones" v-loading="loading">
      <el-table-column prop="name" label="设备名称" />
      <el-table-column label="品牌/型号">
        <template #default="{ row }">{{ row.brand }} {{ row.model }}</template>
      </el-table-column>
      <el-table-column prop="price_per_day" label="日租金">
        <template #default="{ row }">¥{{ row.price_per_day }}</template>
      </el-table-column>
      <el-table-column label="可用时段">
        <template #default="{ row }">
          {{ formatAvailableDate(row.available_from) }} ~ {{ formatAvailableDate(row.available_to) }}
        </template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="statusTagType(row.status)">{{ statusText(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button v-if="row.status === 'offline'" type="success" link @click="updateStatus(row, 'online')">上架</el-button>
          <el-button v-if="row.status === 'online'" type="warning" link @click="updateStatus(row, 'offline')">下架</el-button>
          <el-button type="primary" link @click="$router.push(`/drone/${row.id}`)">详情</el-button>
          <el-button type="danger" link @click="deleteDrone(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      style="margin-top: 20px; justify-content: flex-end"
      layout="prev, pager, next, total"
      :total="total"
      :page-size="pageSize"
      v-model:current-page="currentPage"
      @current-change="fetchDrones"
    />
  </el-card>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useUserStore } from '@/stores/user'
import request from '@/utils/request'
import dayjs from 'dayjs'

const userStore = useUserStore()
const token = computed(() => userStore.token)

const loading = ref(false)
const drones = ref<Drone[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

onMounted(() => {
  fetchDrones()
})

async function fetchDrones() {
  loading.value = true
  try {
    const res: any = await request.get('/my-drones', {
      params: { page: currentPage.value, page_size: pageSize.value }
    })
    drones.value = res.data.list || []
    total.value = res.data.total || 0
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function updateStatus(row: Drone, status: string) {
  try {
    await request.put(`/drones/${row.id}/status`, { status })
    ElMessage.success('操作成功')
    fetchDrones()
  } catch (e: any) {
    ElMessage.error(e.message || '操作失败')
  }
}

async function deleteDrone(row: Drone) {
  ElMessageBox.confirm(`确定删除设备"${row.name}"吗？`, '提示', {
    type: 'warning'
  }).then(async () => {
    try {
      await request.delete(`/drones/${row.id}`)
      ElMessage.success('删除成功')
      fetchDrones()
    } catch (e: any) {
      ElMessage.error(e.message || '删除失败')
    }
  }).catch(() => {})
}

function handleBatchImport(res: any) {
  ElMessage.success(`成功导入 ${res.data.imported} 台设备`)
  fetchDrones()
}

function handleImportError(err: any) {
  ElMessage.error(err?.message || '导入失败，请检查Excel文件格式')
}

function formatAvailableDate(date: string | null | undefined) {
  if (!date) return '不限'
  return dayjs(date).format('YYYY-MM-DD')
}

function statusText(status: string) {
  const map: Record<string, string> = {
    offline: '已下架', online: '可租', rented: '已租', maintenance: '维护中'
  }
  return map[status] || status
}

function statusTagType(status: string) {
  const map: Record<string, string> = {
    offline: 'info', online: 'success', rented: 'warning', maintenance: 'danger'
  }
  return map[status] || ''
}
</script>

<style scoped>
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
