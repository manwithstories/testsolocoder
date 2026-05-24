<template>
  <div class="admin-works">
    <div class="page-header">
      <h2>作品管理</h2>
    </div>
    
    <div class="filter-bar">
      <el-input 
        v-model="searchKeyword" 
        placeholder="搜索作品标题" 
        clearable
        style="width: 200px;"
        @change="loadWorks"
      />
      <el-select 
        v-model="filterStatus" 
        placeholder="筛选状态" 
        clearable
        style="width: 120px;"
        @change="loadWorks"
      >
        <el-option label="待审核" :value="1" />
        <el-option label="已发布" :value="2" />
        <el-option label="已拒绝" :value="3" />
      </el-select>
      <el-button type="primary" @click="loadWorks">
        <el-icon><Search /></el-icon>
        搜索
      </el-button>
    </div>
    
    <el-table :data="works" v-loading="loading" style="width: 100%">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="title" label="标题" />
      <el-table-column prop="artist_name" label="艺术家" width="120" />
      <el-table-column prop="genre" label="风格" width="100" />
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getStatusTag(row.status)">{{ getStatusText(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="play_count" label="播放量" width="100" />
      <el-table-column prop="created_at" label="上传时间" width="180" />
      <el-table-column label="操作" width="200">
        <template #default="{ row }">
          <template v-if="row.status === 1">
            <el-button 
              type="success" 
              text 
              size="small"
              @click="approveWork(row)"
            >
              通过
            </el-button>
            <el-button 
              type="danger" 
              text 
              size="small"
              @click="rejectWork(row)"
            >
              拒绝
            </el-button>
          </template>
          <template v-else>
            <el-button 
              type="primary" 
              text 
              size="small"
              @click="viewDetail(row)"
            >
              详情
            </el-button>
            <el-button 
              :type="row.status === 2 ? 'warning' : 'success'" 
              text 
              size="small"
              @click="toggleStatus(row)"
            >
              {{ row.status === 2 ? '下架' : '上架' }}
            </el-button>
          </template>
        </template>
      </el-table-column>
    </el-table>
    
    <div class="pagination">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="loadWorks"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { workApi } from '@/api/work'

const loading = ref(false)
const works = ref<any[]>([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const searchKeyword = ref('')
const filterStatus = ref<number | ''>('')

onMounted(() => {
  loadWorks()
})

async function loadWorks() {
  loading.value = true
  try {
    const params: any = {
      page: page.value,
      page_size: pageSize.value
    }
    
    if (searchKeyword.value) {
      params.keyword = searchKeyword.value
    }
    if (filterStatus.value !== '') {
      params.status = filterStatus.value
    }
    
    const res = await workApi.search(params)
    works.value = res.list
    total.value = res.total
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function approveWork(work: any) {
  try {
    await ElMessageBox.confirm('确定要通过该作品吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'success'
    })
    
    await workApi.approve(work.id)
    ElMessage.success('审核通过')
    loadWorks()
  } catch (e) {
    if (e !== 'cancel') {
      console.error(e)
    }
  }
}

async function rejectWork(work: any) {
  try {
    const { value } = await ElMessageBox.prompt('请输入拒绝原因', '拒绝作品', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      inputPattern: /.+/,
      inputErrorMessage: '请输入拒绝原因'
    })
    
    await workApi.reject(work.id, { reason: value })
    ElMessage.success('已拒绝')
    loadWorks()
  } catch (e) {
    if (e !== 'cancel') {
      console.error(e)
    }
  }
}

async function toggleStatus(work: any) {
  try {
    const action = work.status === 2 ? '下架' : '上架'
    await ElMessageBox.confirm(`确定要${action}该作品吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await workApi.updateStatus(work.id, { 
      status: work.status === 2 ? 4 : 2 
    })
    ElMessage.success(`${action}成功`)
    loadWorks()
  } catch (e) {
    if (e !== 'cancel') {
      console.error(e)
    }
  }
}

function viewDetail(work: any) {
  window.open(`/works/${work.id}`, '_blank')
}

function getStatusTag(status: number): 'primary' | 'success' | 'warning' | 'info' | 'danger' {
  const tags: Record<number, 'primary' | 'success' | 'warning' | 'info' | 'danger'> = {
    0: 'info',
    1: 'warning',
    2: 'success',
    3: 'danger',
    4: 'info'
  }
  return tags[status] || 'info'
}

function getStatusText(status: number) {
  const texts: Record<number, string> = {
    0: '草稿',
    1: '待审核',
    2: '已发布',
    3: '已拒绝',
    4: '已下架'
  }
  return texts[status] || '未知'
}
</script>

<style scoped lang="scss">
.admin-works {
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
}
</style>
