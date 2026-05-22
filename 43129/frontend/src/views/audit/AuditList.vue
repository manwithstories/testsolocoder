<template>
  <div class="page-container">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>审计日志</span>
          <div class="header-actions">
            <el-select v-model="filterModule" placeholder="模块" clearable style="width: 120px">
              <el-option label="顾客" value="customer" />
              <el-option label="技师" value="technician" />
              <el-option label="服务" value="service" />
              <el-option label="预约" value="appointment" />
              <el-option label="支付" value="payment" />
              <el-option label="产品" value="product" />
            </el-select>
            <el-button :icon="Search" @click="fetchList">查询</el-button>
          </div>
        </div>
      </template>

      <el-table :data="logs" stripe v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="操作人">
          <template #default="{ row }">
            {{ row.user?.nickname || row.user?.phone || '系统' }}
          </template>
        </el-table-column>
        <el-table-column prop="module" label="模块" width="120">
          <template #default="{ row }">{{ getModuleText(row.module) }}</template>
        </el-table-column>
        <el-table-column prop="action" label="操作" width="120">
          <template #default="{ row }">{{ getActionText(row.action) }}</template>
        </el-table-column>
        <el-table-column prop="detail" label="详情" show-overflow-tooltip />
        <el-table-column prop="ip" label="IP地址" width="150" />
        <el-table-column prop="created_at" label="时间" width="180">
          <template #default="{ row }">{{ formatDateTime(row.created_at) }}</template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="fetchList"
          @size-change="fetchList"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getAuditLogs } from '@/api/notification'
import { Search } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import type { AuditLog } from '@/types'

const loading = ref(false)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const logs = ref<AuditLog[]>([])
const filterModule = ref('')

const getModuleText = (module: string) => {
  const modules: Record<string, string> = {
    customer: '顾客',
    technician: '技师',
    service: '服务',
    appointment: '预约',
    payment: '支付',
    product: '产品',
    auth: '认证'
  }
  return modules[module] || module
}

const getActionText = (action: string) => {
  const actions: Record<string, string> = {
    create: '创建',
    update: '更新',
    delete: '删除',
    login: '登录',
    cancel: '取消',
    complete: '完成'
  }
  return actions[action] || action
}

const formatDateTime = (date: string) => dayjs(date).format('YYYY-MM-DD HH:mm:ss')

const fetchList = async () => {
  loading.value = true
  try {
    const params: any = {
      page: page.value,
      page_size: pageSize.value
    }
    if (filterModule.value) {
      params.module = filterModule.value
    }
    const res = await getAuditLogs(params)
    logs.value = res.data.list
    total.value = res.data.total
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(fetchList)
</script>

<style scoped lang="scss">
.page-container {
  .card-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .header-actions {
    display: flex;
    gap: 10px;
  }

  .pagination {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }
}
</style>
