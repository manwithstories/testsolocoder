<template>
  <div class="my-works">
    <div class="page-header">
      <h2>我的作品</h2>
      <el-button type="primary" @click="goToUpload">
        <el-icon><Upload /></el-icon>
        上传作品
      </el-button>
    </div>
    
    <el-tabs v-model="activeTab">
      <el-tab-pane label="全部" name="all" />
      <el-tab-pane label="草稿" name="draft" />
      <el-tab-pane label="已发布" name="published" />
      <el-tab-pane label="已下架" name="offline" />
    </el-tabs>
    
    <div class="works-list" v-loading="loading">
      <div 
        v-for="work in works" 
        :key="work.id"
        class="work-item"
      >
        <el-image :src="work.cover_url" class="cover" fit="cover">
          <template #error>
            <div class="image-placeholder">
              <el-icon :size="20"><Headset /></el-icon>
            </div>
          </template>
        </el-image>
        <div class="info">
          <div class="title text-ellipsis">{{ work.title }}</div>
          <div class="meta">
            <el-tag :type="getStatusType(work.status)" size="small">
              {{ getStatusText(work.status) }}
            </el-tag>
            <span>{{ formatCount(work.play_count) }} 播放</span>
          </div>
        </div>
        <div class="actions">
          <el-button text @click="editWork(work)">
            <el-icon><Edit /></el-icon>
          </el-button>
          <el-button text type="danger" @click="deleteWork(work)">
            <el-icon><Delete /></el-icon>
          </el-button>
        </div>
      </div>
      
      <el-empty v-if="works.length === 0 && !loading" description="暂无作品" />
    </div>
    
    <div class="pagination" v-if="total > 0">
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
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { workApi } from '@/api/work'
import { useUserStore } from '@/stores/user'
import type { Work } from '@/api/work'

const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const activeTab = ref('all')
const works = ref<Work[]>([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

onMounted(() => {
  loadWorks()
})

async function loadWorks() {
  if (!userStore.user?.artist_info) return
  
  loading.value = true
  try {
    const status = activeTab.value === 'all' ? -1 : 
      activeTab.value === 'draft' ? 0 :
      activeTab.value === 'published' ? 2 : 4
    
    const res = await workApi.getByArtist(userStore.user.artist_info.id, {
      page: page.value,
      page_size: pageSize.value,
      status: status
    })
    works.value = res.list
    total.value = res.total
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function goToUpload() {
  router.push('/user/upload')
}

function editWork(work: Work) {
  ElMessage.info('编辑功能开发中')
}

async function deleteWork(work: Work) {
  try {
    await ElMessageBox.confirm('确定要删除这个作品吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await workApi.delete(work.id)
    ElMessage.success('删除成功')
    loadWorks()
  } catch (e) {
    if (e !== 'cancel') {
      console.error(e)
    }
  }
}

function getStatusType(status: number): 'primary' | 'success' | 'warning' | 'info' | 'danger' {
  const types: Record<number, 'primary' | 'success' | 'warning' | 'info' | 'danger'> = {
    0: 'info',
    1: 'warning',
    2: 'success',
    3: 'danger',
    4: 'info'
  }
  return types[status] || 'info'
}

function getStatusText(status: number) {
  const texts: Record<number, string> = {
    0: '草稿',
    1: '审核中',
    2: '已发布',
    3: '已拒绝',
    4: '已下架'
  }
  return texts[status] || '未知'
}

function formatCount(count: number) {
  if (count >= 10000) {
    return (count / 10000).toFixed(1) + 'w'
  }
  return count
}
</script>

<style scoped lang="scss">
.my-works {
  .works-list {
    .work-item {
      display: flex;
      align-items: center;
      gap: 16px;
      padding: 16px;
      border-radius: 8px;
      border-bottom: 1px solid var(--border-color);
      
      &:hover {
        background: rgba(64, 158, 255, 0.05);
      }
      
      .cover {
        width: 80px;
        height: 80px;
        border-radius: 4px;
        overflow: hidden;
        flex-shrink: 0;
      }
      
      .info {
        flex: 1;
        min-width: 0;
        
        .title {
          font-weight: 500;
          margin-bottom: 8px;
        }
        
        .meta {
          display: flex;
          align-items: center;
          gap: 12px;
          font-size: 13px;
          color: var(--text-light);
        }
      }
      
      .actions {
        display: flex;
        gap: 8px;
      }
    }
  }
  
  .pagination {
    display: flex;
    justify-content: center;
    margin-top: 20px;
  }
  
  .image-placeholder {
    width: 100%;
    height: 100%;
    background: #f5f7fa;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #c0c4cc;
  }
}
</style>
