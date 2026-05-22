<template>
  <div class="notice-list">
    <div class="page-header">
      <h2 class="page-title">公告通知</h2>
      <el-button type="primary" @click="goToCreate">
        <el-icon><Plus /></el-icon>发布公告
      </el-button>
    </div>

    <div class="search-bar card">
      <el-select v-model="typeFilter" placeholder="类型" clearable style="width: 150px">
        <el-option label="普通公告" :value="1" />
        <el-option label="紧急通知" :value="2" />
      </el-select>
      <el-select v-model="buildingFilter" placeholder="楼栋" clearable style="width: 150px">
        <el-option label="A栋" value="A栋" />
        <el-option label="B栋" value="B栋" />
        <el-option label="C栋" value="C栋" />
        <el-option label="D栋" value="D栋" />
      </el-select>
      <el-button type="primary" @click="loadData">搜索</el-button>
    </div>

    <div class="card">
      <el-table :data="notices" v-loading="loading">
        <el-table-column label="标题" min-width="200">
          <template #default="{ row }">
            <div class="title-cell">
              <el-tag v-if="row.isTop" type="danger" size="small" style="margin-right: 8px;">置顶</el-tag>
              <el-tag v-if="row.type === 2" type="warning" size="small" style="margin-right: 8px;">紧急</el-tag>
              {{ row.title }}
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="building" label="推送范围" width="120">
          <template #default="{ row }">
            {{ row.building || '全部' }}
          </template>
        </el-table-column>
        <el-table-column label="发布人" width="120">
          <template #default="{ row }">
            {{ row.publisher?.realName }}
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="发布时间" width="180">
          <template #default="{ row }">{{ formatDate(row.createdAt) }}</template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">
              {{ row.status === 1 ? '已发布' : '已下架' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="viewDetail(row.id)">查看</el-button>
            <el-button
              link
              type="warning"
              @click="toggleTop(row)"
            >{{ row.isTop ? '取消置顶' : '置顶' }}</el-button>
            <el-button
              link
              :type="row.status === 1 ? 'info' : 'success'"
              @click="toggleStatus(row)"
            >{{ row.status === 1 ? '下架' : '发布' }}</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        class="pagination"
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next, jumper"
        @current-change="loadData"
        @size-change="loadData"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { Notice } from '@/types'
import { getNotices, updateNoticeStatus, updateNoticeFields } from '@/api/business'
import { updateNotice } from '@/api/business'
import dayjs from 'dayjs'

const router = useRouter()
const loading = ref(false)
const notices = ref<Notice[]>([])
const typeFilter = ref<number | ''>('')
const buildingFilter = ref('')

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

onMounted(() => {
  loadData()
})

async function loadData() {
  loading.value = true
  try {
    const res = await getNotices({
      page: pagination.page,
      pageSize: pagination.pageSize,
      type: typeFilter.value || undefined,
      building: buildingFilter.value || undefined
    })
    notices.value = res.data.list
    pagination.total = res.data.total
  } catch (error) {
    console.error('Failed to load notices:', error)
  } finally {
    loading.value = false
  }
}

function formatDate(date: string) {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

function goToCreate() {
  router.push('/notices/create')
}

function viewDetail(id: number) {
  router.push(`/notices/${id}`)
}

async function toggleTop(row: Notice) {
  try {
    await updateNoticeFields(row.id, { isTop: row.isTop ? 0 : 1 })
    ElMessage.success('操作成功')
    loadData()
  } catch (error) {
    console.error(error)
  }
}

async function toggleStatus(row: Notice) {
  try {
    await updateNoticeStatus(row.id, row.status === 1 ? 0 : 1)
    ElMessage.success('操作成功')
    loadData()
  } catch (error) {
    console.error(error)
  }
}
</script>

<style scoped>
.notice-list {
  padding: 0;
}

.title-cell {
  display: flex;
  align-items: center;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
