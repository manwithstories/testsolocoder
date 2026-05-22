<template>
  <div class="activity-list">
    <div class="page-header">
      <h2 class="page-title">活动管理</h2>
      <el-button type="primary" @click="$router.push('/activities/create')" v-if="userStore.isAdmin">
        <el-icon><Plus /></el-icon>
        创建活动
      </el-button>
    </div>

    <el-card>
      <div class="search-bar">
        <el-input v-model="search.keyword" placeholder="搜索活动名称" clearable style="width: 240px" @keyup.enter="loadData" />
        <el-select v-model="search.status" placeholder="活动状态" clearable style="width: 140px">
          <el-option label="草稿" value="draft" />
          <el-option label="已发布" value="published" />
          <el-option label="已取消" value="canceled" />
        </el-select>
        <el-button type="primary" @click="loadData">搜索</el-button>
      </div>

      <el-table :data="list" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="poster" label="海报" width="120">
          <template #default="{ row }">
            <img v-if="row.poster" :src="row.poster" class="poster-img" />
            <span v-else class="no-poster">无</span>
          </template>
        </el-table-column>
        <el-table-column prop="title" label="活动名称" min-width="180" />
        <el-table-column prop="location" label="地点" min-width="160" />
        <el-table-column label="时间" min-width="260">
          <template #default="{ row }">
            <div>{{ formatDate(row.startTime) }}</div>
            <div class="sub-text">至 {{ formatDate(row.endTime) }}</div>
          </template>
        </el-table-column>
        <el-table-column prop="capacity" label="容量" width="100" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <div class="table-actions">
              <el-button size="small" @click="$router.push(`/activities/${row.id}`)">查看</el-button>
              <el-button size="small" type="primary" @click="$router.push(`/activities/${row.id}/edit`)" v-if="userStore.isAdmin">编辑</el-button>
              <el-button size="small" type="danger" @click="handleDelete(row.id)" v-if="userStore.isAdmin">删除</el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="loadData"
        @current-change="loadData"
        style="margin-top: 20px; justify-content: flex-end"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getActivityList, deleteActivity } from '@/api/activity'
import { useUserStore } from '@/store/user'
import dayjs from 'dayjs'

const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const list = ref<any[]>([])

const search = reactive({
  keyword: '',
  status: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const loadData = async () => {
  try {
    loading.value = true
    const res = await getActivityList({
      page: pagination.page,
      pageSize: pagination.pageSize,
      keyword: search.keyword,
      status: search.status
    })
    list.value = res.list
    pagination.total = res.total
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleDelete = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定要删除这个活动吗？', '提示', { type: 'warning' })
    await deleteActivity(id)
    ElMessage.success('删除成功')
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error(error)
    }
  }
}

const getStatusType = (status: string) => {
  const map: Record<string, string> = { draft: 'info', published: 'success', canceled: 'danger' }
  return map[status] || 'info'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = { draft: '草稿', published: '已发布', canceled: '已取消' }
  return map[status] || status
}

const formatDate = (date: string) => dayjs(date).format('YYYY-MM-DD HH:mm')

onMounted(loadData)
</script>

<style scoped lang="scss">
.poster-img {
  width: 80px;
  height: 60px;
  object-fit: cover;
  border-radius: 4px;
}

.no-poster {
  color: #c0c4cc;
}

.sub-text {
  font-size: 12px;
  color: #909399;
  margin-top: 2px;
}
</style>
