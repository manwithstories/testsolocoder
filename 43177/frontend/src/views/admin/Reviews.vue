<template>
  <div class="admin-reviews-page">
    <div class="page-header">
      <h2 class="page-title">差评处理</h2>
    </div>

    <el-card>
      <el-table :data="reviews" style="width: 100%">
        <el-table-column label="工单号" width="200">
          <template #default="{ row }">{{ row.order?.order_no }}</template>
        </el-table-column>
        <el-table-column label="客户" width="150">
          <template #default="{ row }">{{ row.customer?.username }}</template>
        </el-table-column>
        <el-table-column label="技师" width="150">
          <template #default="{ row }">{{ row.technician?.username }}</template>
        </el-table-column>
        <el-table-column prop="rating" label="评分" width="100">
          <template #default="{ row }">
            <el-rate :model-value="row.rating" disabled size="small" />
          </template>
        </el-table-column>
        <el-table-column prop="content" label="评价内容" min-width="200" />
        <el-table-column label="是否介入" width="100">
          <template #default="{ row }">
            <el-tag :type="row.is_intervened ? 'success' : 'warning'">
              {{ row.is_intervened ? '已介入' : '待处理' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <el-button
              v-if="!row.is_intervened"
              type="primary"
              size="small"
              @click="intervene(row)"
            >
              介入处理
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div v-if="reviews.length === 0" class="empty-state">
        <el-empty description="暂无差评" />
      </div>
    </el-card>

    <el-dialog v-model="showInterveneDialog" title="平台介入" width="500px">
      <el-form label-width="80px">
        <el-form-item label="处理备注">
          <el-input v-model="interveneNote" type="textarea" :rows="3" placeholder="请输入处理意见" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showInterveneDialog = false">取消</el-button>
        <el-button type="primary" @click="submitIntervene" :loading="submitting">
          提交
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { adminApi } from '@/api/admin'
import type { Review } from '@/types'

const reviews = ref<Review[]>([])
const showInterveneDialog = ref(false)
const selectedReview = ref<Review | null>(null)
const interveneNote = ref('')
const submitting = ref(false)

onMounted(() => {
  loadReviews()
})

async function loadReviews() {
  try {
    const res = await adminApi.getLowRatingReviews({ page: 1, page_size: 50 })
    reviews.value = res.data?.list || []
  } catch (error) {
    console.error('Failed to load reviews:', error)
  }
}

function intervene(review: Review) {
  selectedReview.value = review
  interveneNote.value = ''
  showInterveneDialog.value = true
}

async function submitIntervene() {
  if (!interveneNote.value) {
    ElMessage.error('请输入处理意见')
    return
  }

  submitting.value = true
  try {
    await adminApi.interveneReview(selectedReview.value!.id, { note: interveneNote.value })
    ElMessage.success('已介入处理')
    showInterveneDialog.value = false
    loadReviews()
  } catch (error) {
    console.error('Failed to intervene:', error)
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.admin-reviews-page {
  padding: 0;
}

.empty-state {
  padding: 40px;
}
</style>
