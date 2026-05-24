<template>
  <div class="page-container">
    <div class="card">
      <h2 class="section-title">鉴定任务</h2>
      <div class="filter-bar">
        <el-radio-group v-model="status" size="default" @change="loadAuthentications">
          <el-radio-button label="">全部</el-radio-button>
          <el-radio-button label="pending">待接单</el-radio-button>
          <el-radio-button label="accepted">鉴定中</el-radio-button>
          <el-radio-button label="completed">已完成</el-radio-button>
        </el-radio-group>
      </div>
    </div>

    <div v-if="loading" class="loading-wrapper">
      <el-icon class="is-loading" :size="32"><Loading /></el-icon>
    </div>

    <el-table
      v-else-if="authentications.length > 0"
      :data="authentications"
      style="width: 100%"
      stripe
    >
      <el-table-column label="鉴定ID" prop="id" width="80" />
      <el-table-column label="商品" min-width="200">
        <template #default="{ row }">
          <div class="product-cell">
            <img v-if="row.product?.images?.[0]" :src="row.product.images[0].image_url" />
            <div class="product-info">
              <span class="product-title text-ellipsis">{{ row.product?.title }}</span>
              <span class="product-brand">{{ row.product?.brand_name || '其他品牌' }}</span>
            </div>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="买家" width="120">
        <template #default="{ row }">
          {{ row.buyer?.username }}
        </template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getAuthStatusType(row.status)">
            {{ getAuthStatusLabel(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="结果" width="100">
        <template #default="{ row }">
          <el-tag v-if="row.result" :type="getAuthResultType(row.result)">
            {{ getAuthResultLabel(row.result) }}
          </el-tag>
          <span v-else>待鉴定</span>
        </template>
      </el-table-column>
      <el-table-column label="申请时间" width="180">
        <template #default="{ row }">
          {{ formatDate(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button
            v-if="row.status === 'pending' && userStore.userRole === 'authenticator'"
            type="primary"
            size="small"
            @click="handleAccept(row)"
          >
            接单
          </el-button>
          <el-button
            v-if="row.status === 'accepted' && isMyTask(row)"
            type="success"
            size="small"
            @click="handleComplete(row)"
          >
            完成鉴定
          </el-button>
          <el-button
            size="small"
            @click="$router.push(`/authentications/${row.id}`)"
          >
            详情
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <div v-else class="empty-state">
      <el-icon :size="64"><Box /></el-icon>
      <p>暂无鉴定任务</p>
    </div>

    <div v-if="!loading && total > 0" class="pagination-wrapper">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="loadAuthentications"
      />
    </div>

    <el-dialog v-model="completeDialogVisible" title="完成鉴定" width="600px">
      <el-form :model="completeForm" label-width="100px">
        <el-form-item label="鉴定结果" required>
          <el-radio-group v-model="completeForm.result">
            <el-radio-button label="genuine">正品</el-radio-button>
            <el-radio-button label="counterfeit">赝品</el-radio-button>
            <el-radio-button label="inconclusive">无法鉴定</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="鉴定报告">
          <el-input
            v-model="completeForm.report_content"
            type="textarea"
            :rows="6"
            placeholder="请输入详细的鉴定报告"
          />
        </el-form-item>
        <el-form-item label="备注">
          <el-input
            v-model="completeForm.authenticator_notes"
            type="textarea"
            :rows="3"
            placeholder="可选，输入备注信息"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="completeDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="completing" @click="confirmComplete">
          提交鉴定结果
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { authServiceApi } from '@/api/authentication'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import type { Authentication } from '@/types'
import { AUTH_STATUS_OPTIONS, AUTH_RESULT_OPTIONS } from '@/types'
import dayjs from 'dayjs'
import { Loading, Box } from '@element-plus/icons-vue'

const userStore = useUserStore()

const authentications = ref<Authentication[]>([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const status = ref('')

const completeDialogVisible = ref(false)
const currentAuth = ref<Authentication | null>(null)
const completing = ref(false)

const completeForm = ref({
  result: 'genuine',
  report_content: '',
  authenticator_notes: ''
})

const loadAuthentications = async () => {
  loading.value = true
  try {
    const res = await authServiceApi.listAuthentications({
      page: page.value,
      page_size: pageSize.value,
      status: status.value || undefined
    })
    if (res.code === 200) {
      const data = res.data as any
      authentications.value = data?.list || []
      total.value = data?.total || 0
    }
  } catch (error) {
    console.error('Load authentications error:', error)
  } finally {
    loading.value = false
  }
}

const getAuthStatusType = (status: string) => {
  const opt = AUTH_STATUS_OPTIONS.find(o => o.value === status)
  return opt?.type || 'info'
}

const getAuthStatusLabel = (status: string) => {
  const opt = AUTH_STATUS_OPTIONS.find(o => o.value === status)
  return opt?.label || status
}

const getAuthResultType = (result: string) => {
  const opt = AUTH_RESULT_OPTIONS.find(o => o.value === result)
  return opt?.type || 'info'
}

const getAuthResultLabel = (result: string) => {
  const opt = AUTH_RESULT_OPTIONS.find(o => o.value === result)
  return opt?.label || result
}

const formatDate = (date: string) => dayjs(date).format('YYYY-MM-DD HH:mm:ss')

const isMyTask = (auth: Authentication) => {
  return auth.authenticator_id === userStore.userID
}

const handleAccept = async (auth: Authentication) => {
  try {
    const res = await authServiceApi.acceptAuthentication(auth.id)
    if (res.code === 200) {
      ElMessage.success('接单成功')
      loadAuthentications()
    }
  } catch (error) {
    console.error('Accept error:', error)
  }
}

const handleComplete = (auth: Authentication) => {
  currentAuth.value = auth
  completeForm.value = {
    result: 'genuine',
    report_content: '',
    authenticator_notes: ''
  }
  completeDialogVisible.value = true
}

const confirmComplete = async () => {
  if (!currentAuth.value) return

  completing.value = true
  try {
    const res = await authServiceApi.completeAuthentication(
      currentAuth.value.id,
      completeForm.value
    )
    if (res.code === 200) {
      ElMessage.success('鉴定完成')
      completeDialogVisible.value = false
      loadAuthentications()
    }
  } catch (error) {
    console.error('Complete error:', error)
  } finally {
    completing.value = false
  }
}

onMounted(() => {
  loadAuthentications()
})
</script>

<style lang="scss" scoped>
.filter-bar {
  display: flex;
  gap: 12px;
  margin-top: 16px;
}

.product-cell {
  display: flex;
  gap: 12px;
  align-items: center;
  
  img {
    width: 60px;
    height: 60px;
    object-fit: cover;
    border-radius: 4px;
  }
  
  .product-info {
    display: flex;
    flex-direction: column;
    gap: 4px;
    
    .product-title {
      font-size: 14px;
      max-width: 200px;
    }
    
    .product-brand {
      font-size: 12px;
      color: var(--text-light);
    }
  }
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 24px;
}
</style>
