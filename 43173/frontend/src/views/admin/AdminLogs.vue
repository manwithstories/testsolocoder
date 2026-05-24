<template>
  <div class="admin-logs">
    <div class="page-header">
      <h2>操作日志</h2>
    </div>
    
    <div class="filter-bar">
      <el-input 
        v-model="searchKeyword" 
        placeholder="搜索用户ID或操作" 
        clearable
        style="width: 200px;"
        @change="loadLogs"
      />
      <el-date-picker
        v-model="dateRange"
        type="daterange"
        range-separator="至"
        start-placeholder="开始日期"
        end-placeholder="结束日期"
        @change="loadLogs"
      />
      <el-button type="primary" @click="loadLogs">
        <el-icon><Search /></el-icon>
        搜索
      </el-button>
      <el-button type="success" @click="exportExcel">
        <el-icon><Download /></el-icon>
        导出
      </el-button>
    </div>
    
    <el-table :data="logs" v-loading="loading" style="width: 100%">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="user_id" label="用户ID" width="100" />
      <el-table-column prop="username" label="用户名" width="120" />
      <el-table-column prop="action" label="操作" width="150" />
      <el-table-column prop="method" label="请求方法" width="100">
        <template #default="{ row }">
          <el-tag :type="getMethodTag(row.method)" size="small">
            {{ row.method }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="path" label="请求路径" />
      <el-table-column prop="ip" label="IP地址" width="140" />
      <el-table-column prop="status" label="状态" width="80">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
            {{ row.status === 1 ? '成功' : '失败' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="时间" width="180" />
      <el-table-column label="操作" width="80">
        <template #default="{ row }">
          <el-button 
            type="primary" 
            text 
            size="small"
            @click="viewDetail(row)"
          >
            详情
          </el-button>
        </template>
      </el-table-column>
    </el-table>
    
    <div class="pagination">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="loadLogs"
      />
    </div>
    
    <el-dialog v-model="detailDialog.visible" title="日志详情" width="600px">
      <el-descriptions :column="1" border>
        <el-descriptions-item label="用户ID">
          {{ detailDialog.data.user_id }}
        </el-descriptions-item>
        <el-descriptions-item label="用户名">
          {{ detailDialog.data.username }}
        </el-descriptions-item>
        <el-descriptions-item label="操作">
          {{ detailDialog.data.action }}
        </el-descriptions-item>
        <el-descriptions-item label="请求方法">
          {{ detailDialog.data.method }}
        </el-descriptions-item>
        <el-descriptions-item label="请求路径">
          {{ detailDialog.data.path }}
        </el-descriptions-item>
        <el-descriptions-item label="请求参数">
          <pre class="params-pre">{{ detailDialog.data.params }}</pre>
        </el-descriptions-item>
        <el-descriptions-item label="IP地址">
          {{ detailDialog.data.ip }}
        </el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="detailDialog.data.status === 1 ? 'success' : 'danger'">
            {{ detailDialog.data.status === 1 ? '成功' : '失败' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item v-if="detailDialog.data.error_msg" label="错误信息">
          {{ detailDialog.data.error_msg }}
        </el-descriptions-item>
        <el-descriptions-item label="操作时间">
          {{ detailDialog.data.created_at }}
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { userApi } from '@/api/auth'

const loading = ref(false)
const logs = ref<any[]>([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const searchKeyword = ref('')
const dateRange = ref<string[]>([])

const detailDialog = reactive({
  visible: false,
  data: {} as any
})

onMounted(() => {
  loadLogs()
})

async function loadLogs() {
  loading.value = true
  try {
    const params: any = {
      page: page.value,
      page_size: pageSize.value
    }
    
    if (searchKeyword.value) {
      params.keyword = searchKeyword.value
    }
    if (dateRange.value?.length === 2) {
      params.start_date = dateRange.value[0]
      params.end_date = dateRange.value[1]
    }
    
    const res = await userApi.getOperationLogs(params)
    logs.value = res.list
    total.value = res.total
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function viewDetail(log: any) {
  detailDialog.data = { ...log }
  detailDialog.visible = true
}

function getMethodTag(method: string): 'primary' | 'success' | 'warning' | 'info' | 'danger' {
  const tags: Record<string, 'primary' | 'success' | 'warning' | 'info' | 'danger'> = {
    GET: 'success',
    POST: 'primary',
    PUT: 'warning',
    DELETE: 'danger',
    PATCH: 'info'
  }
  return tags[method] || 'info'
}

function exportExcel() {
  const params = new URLSearchParams()
  if (searchKeyword.value) {
    params.append('keyword', searchKeyword.value)
  }
  if (dateRange.value?.length === 2) {
    params.append('start_date', dateRange.value[0])
    params.append('end_date', dateRange.value[1])
  }
  
  window.open(`/api/v1/export/logs?${params.toString()}`, '_blank')
}
</script>

<style scoped lang="scss">
.admin-logs {
  .filter-bar {
    display: flex;
    gap: 12px;
    margin-bottom: 16px;
  }
  
  .pagination {
    display: flex;
    justify-content: center;
    margin-top: 20px;
  }
  
  .params-pre {
    background: #f5f7fa;
    padding: 8px;
    border-radius: 4px;
    font-size: 12px;
    max-height: 200px;
    overflow: auto;
    margin: 0;
    white-space: pre-wrap;
    word-break: break-all;
  }
}
</style>
