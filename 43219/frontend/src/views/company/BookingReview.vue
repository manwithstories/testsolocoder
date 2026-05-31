<template>
  <AppLayout>
    <div class="page">
      <h2>预约审核</h2>
      <el-alert v-if="needReviewList.length" :title="`有 ${needReviewList.length} 个预约改期超3次需人工审核`" type="warning" show-icon :closable="false" style="margin-bottom:16px" />
      <el-table :data="list" border stripe>
        <el-table-column prop="id" label="编号" width="80" />
        <el-table-column label="服务" width="200">
          <template #default="{ row }">{{ row.service?.name }}</template>
        </el-table-column>
        <el-table-column label="家政人员" width="140">
          <template #default="{ row }">{{ row.staff?.real_name || row.staff?.username }}</template>
        </el-table-column>
        <el-table-column prop="start_at" label="开始时间" width="170" />
        <el-table-column prop="end_at" label="结束时间" width="170" />
        <el-table-column prop="price" label="金额" width="100" />
        <el-table-column prop="reschedule_count" label="改期次数" width="100" />
        <el-table-column label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="row.need_review ? 'warning' : ''">
              {{ row.need_review ? '待人工审核' : statusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="240">
          <template #default="{ row }">
            <el-button size="small" type="primary" v-if="row.need_review" @click="review(row)">审核改期</el-button>
            <el-button size="small" v-if="row.status==='pending'" @click="confirm(row)">确认</el-button>
            <el-button size="small" type="danger" v-if="row.status==='pending' || row.status==='confirmed'" @click="cancel(row)">取消</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-dialog v-model="reviewVisible" title="改期审核">
      <el-form label-width="80px">
        <el-form-item label="新开始">
          <el-date-picker v-model="reviewForm.start_at" type="datetime" value-format="YYYY-MM-DDTHH:mm:ss" />
        </el-form-item>
        <el-form-item label="新结束">
          <el-date-picker v-model="reviewForm.end_at" type="datetime" value-format="YYYY-MM-DDTHH:mm:ss" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="submitReview(false)">拒绝</el-button>
        <el-button type="primary" @click="submitReview(true)">同意</el-button>
      </template>
    </el-dialog>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, computed, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import AppLayout from '../../components/AppLayout.vue'
import { listBookings, confirmBooking, cancelBooking, reviewReschedule, type Booking } from '../../api/booking'

const list = ref<Booking[]>([])
const reviewVisible = ref(false)
const currentId = ref<number | null>(null)
const reviewForm = reactive({ start_at: '', end_at: '' })
const needReviewList = computed(() => list.value.filter(b => b.need_review))

async function load() {
  const res = await listBookings({})
  list.value = (res.data as any).data || []
}

function statusLabel(s: string) {
  return { pending: '待确认', confirmed: '已确认', canceled: '已取消', rejected: '已拒绝', completed: '已完成' }[s] || s
}

async function confirm(row: Booking) { await confirmBooking(row.id); ElMessage.success('已确认'); load() }
async function cancel(row: Booking) { await cancelBooking(row.id); ElMessage.success('已取消'); load() }
function review(row: Booking) { currentId.value = row.id; reviewForm.start_at = ''; reviewForm.end_at = ''; reviewVisible.value = true }
async function submitReview(approve: boolean) {
  await reviewReschedule(currentId.value!, { approve, ...(approve ? reviewForm : {}) })
  ElMessage.success('已处理')
  reviewVisible.value = false
  load()
}

onMounted(load)
</script>
