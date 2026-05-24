<template>
  <el-card>
    <template #header>
      <span>评价中心</span>
    </template>

    <el-table :data="reviews" v-loading="loading">
      <el-table-column label="类型" width="80">
        <template #default="{ row }">
          <el-tag size="small">{{ row.type === 'rental' ? '租赁' : '服务' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="评分" width="150">
        <template #default="{ row }">
          <el-rate :model-value="row.rating" disabled size="small" />
        </template>
      </el-table-column>
      <el-table-column prop="content" label="评价内容" show-overflow-tooltip />
      <el-table-column label="评价人" width="120">
        <template #default="{ row }">{{ row.reviewer?.nickname || row.reviewer?.username }}</template>
      </el-table-column>
      <el-table-column label="被评价人" width="120">
        <template #default="{ row }">{{ row.reviewee?.nickname || row.reviewee?.username }}</template>
      </el-table-column>
      <el-table-column prop="reply" label="回复" show-overflow-tooltip />
      <el-table-column prop="created_at" label="时间" width="180" />
      <el-table-column label="操作" width="120" fixed="right">
        <template #default="{ row }">
          <el-button v-if="!row.reply && role !== row.reviewer_id" type="primary" link @click="reply(row)">回复</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      style="margin-top: 20px; justify-content: flex-end"
      layout="prev, pager, next, total"
      :total="total"
      :page-size="pageSize"
      v-model:current-page="currentPage"
      @current-change="fetchReviews"
    />

    <el-dialog v-model="showReply" title="回复评价" width="400px">
      <el-input v-model="replyContent" type="textarea" :rows="3" placeholder="请输入回复内容" />
      <template #footer>
        <el-button @click="showReply = false">取消</el-button>
        <el-button type="primary" @click="submitReply" :loading="submitting">提交</el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'

const userStore = useUserStore()
const role = computed(() => userStore.role)

const loading = ref(false)
const reviews = ref<Review[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const showReply = ref(false)
const submitting = ref(false)
const replyContent = ref('')
const currentReview = ref<Review | null>(null)

onMounted(() => {
  fetchReviews()
})

async function fetchReviews() {
  loading.value = true
  try {
    const res: any = await request.get('/reviews', {
      params: { page: currentPage.value, page_size: pageSize.value }
    })
    reviews.value = res.data.list || []
    total.value = res.data.total || 0
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function reply(row: Review) {
  currentReview.value = row
  replyContent.value = ''
  showReply.value = true
}

async function submitReply() {
  if (!replyContent.value.trim()) {
    ElMessage.warning('请输入回复内容')
    return
  }
  submitting.value = true
  try {
    await request.post('/reviews/reply', {
      review_id: currentReview.value?.id,
      reply: replyContent.value
    })
    ElMessage.success('回复成功')
    showReply.value = false
    fetchReviews()
  } catch (e: any) {
    ElMessage.error(e.message || '回复失败')
  } finally {
    submitting.value = false
  }
}
</script>
