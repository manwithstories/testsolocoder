<template>
  <div class="activity-detail" v-if="activity">
    <div class="page-header">
      <h2 class="page-title">{{ activity.title }}</h2>
      <div>
        <el-button @click="$router.back()">返回</el-button>
        <el-button type="primary" @click="$router.push(`/activities/${activity.id}/edit`)" v-if="userStore.isAdmin">编辑</el-button>
      </div>
    </div>

    <el-row :gutter="20">
      <el-col :span="16">
        <el-card>
          <img v-if="activity.poster" :src="activity.poster" class="poster" />
          <div class="content" v-html="activity.description"></div>
          
          <el-divider>票型信息</el-divider>
          
          <el-row :gutter="20">
            <el-col :span="12" v-for="ticket in activity.ticketTypes" :key="ticket.id">
              <el-card class="ticket-card" :class="{ soldout: ticket.status === 'sold_out' }">
                <div class="ticket-header">
                  <span class="ticket-name">{{ ticket.name }}</span>
                  <el-tag :type="getTicketTypeColor(ticket.type)">{{ getTicketTypeText(ticket.type) }}</el-tag>
                </div>
                <div class="ticket-price">
                  <span class="currency">¥</span>
                  <span class="amount">{{ ticket.price }}</span>
                </div>
                <div class="ticket-info">
                  <span>库存: {{ ticket.stock }}</span>
                  <span>已售: {{ ticket.soldCount }}</span>
                </div>
                <el-tag v-if="ticket.status === 'sold_out'" type="info" class="soldout-tag">已售罄</el-tag>
                <el-button 
                  v-else 
                  type="primary" 
                  size="small" 
                  style="width: 100%; margin-top: 10px"
                  @click="handleBuy(ticket)"
                >
                  立即购买
                </el-button>
              </el-card>
            </el-col>
          </el-row>
        </el-card>
      </el-col>
      
      <el-col :span="8">
        <el-card>
          <template #header>活动信息</template>
          <div class="info-item">
            <span class="label">状态</span>
            <el-tag :type="getStatusType(activity.status)">{{ getStatusText(activity.status) }}</el-tag>
          </div>
          <div class="info-item">
            <span class="label">开始时间</span>
            <span>{{ formatDate(activity.startTime) }}</span>
          </div>
          <div class="info-item">
            <span class="label">结束时间</span>
            <span>{{ formatDate(activity.endTime) }}</span>
          </div>
          <div class="info-item">
            <span class="label">活动地点</span>
            <span>{{ activity.location }}</span>
          </div>
          <div class="info-item">
            <span class="label">活动容量</span>
            <span>{{ activity.capacity }} 人</span>
          </div>
        </el-card>
        
        <el-card style="margin-top: 20px" v-if="userStore.isAdmin">
          <template #header>快捷操作</template>
          <el-select v-model="newStatus" placeholder="修改状态" style="width: 100%; margin-bottom: 10px">
            <el-option label="草稿" value="draft" />
            <el-option label="发布" value="published" />
            <el-option label="取消" value="canceled" />
          </el-select>
          <el-button type="primary" style="width: 100%" @click="handleUpdateStatus">更新状态</el-button>
        </el-card>
      </el-col>
    </el-row>

    <el-dialog v-model="buyDialogVisible" title="购买门票" width="500px">
      <el-form :model="buyForm" label-width="80px">
        <el-form-item label="票型">
          <span>{{ selectedTicket?.name }}</span>
        </el-form-item>
        <el-form-item label="单价">
          <span>¥{{ selectedTicket?.price }}</span>
        </el-form-item>
        <el-form-item label="数量">
          <el-input-number v-model="buyForm.quantity" :min="1" :max="selectedTicket?.stock || 1" />
        </el-form-item>
        <el-form-item label="优惠券">
          <el-input v-model="buyForm.couponCode" placeholder="选填，输入优惠券码" />
        </el-form-item>
        <el-form-item label="总价">
          <span class="total-price">¥{{ (selectedTicket?.price || 0) * buyForm.quantity }}</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="buyDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="buyLoading" @click="submitOrder">确认购买</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getActivity, updateActivityStatus } from '@/api/activity'
import { createOrder } from '@/api/order'
import { useUserStore } from '@/store/user'
import dayjs from 'dayjs'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const activity = ref<any>(null)
const newStatus = ref('')
const buyDialogVisible = ref(false)
const selectedTicket = ref<any>(null)
const buyLoading = ref(false)
const buyForm = reactive({
  quantity: 1,
  couponCode: ''
})

const loadData = async () => {
  try {
    const res = await getActivity(Number(route.params.id))
    activity.value = res
    newStatus.value = res.status
  } catch (error) {
    console.error(error)
  }
}

const handleBuy = (ticket: any) => {
  selectedTicket.value = ticket
  buyForm.quantity = 1
  buyForm.couponCode = ''
  buyDialogVisible.value = true
}

const submitOrder = async () => {
  try {
    buyLoading.value = true
    await createOrder({
      activityId: activity.value.id,
      tickets: [{
        ticketTypeId: selectedTicket.value.id,
        quantity: buyForm.quantity
      }],
      couponCode: buyForm.couponCode
    })
    ElMessage.success('下单成功')
    buyDialogVisible.value = false
    router.push('/orders')
  } catch (error) {
    console.error(error)
  } finally {
    buyLoading.value = false
  }
}

const handleUpdateStatus = async () => {
  if (!newStatus.value) return
  try {
    await updateActivityStatus(activity.value.id, newStatus.value)
    ElMessage.success('状态更新成功')
    loadData()
  } catch (error) {
    console.error(error)
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

const getTicketTypeColor = (type: string) => {
  const map: Record<string, string> = { normal: '', vip: 'warning', early_bird: 'success' }
  return map[type] || ''
}

const getTicketTypeText = (type: string) => {
  const map: Record<string, string> = { normal: '普通票', vip: 'VIP票', early_bird: '早鸟票' }
  return map[type] || type
}

const formatDate = (date: string) => dayjs(date).format('YYYY-MM-DD HH:mm')

onMounted(loadData)
</script>

<style scoped lang="scss">
.poster {
  width: 100%;
  max-height: 400px;
  object-fit: cover;
  border-radius: 8px;
  margin-bottom: 20px;
}

.content {
  line-height: 1.8;
  color: #606266;
}

.ticket-card {
  margin-bottom: 20px;
  
  &.soldout {
    opacity: 0.7;
  }
}

.ticket-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.ticket-name {
  font-size: 16px;
  font-weight: 600;
}

.ticket-price {
  color: #f56c6c;
  margin-bottom: 10px;

  .currency {
    font-size: 14px;
  }

  .amount {
    font-size: 28px;
    font-weight: 600;
  }
}

.ticket-info {
  display: flex;
  justify-content: space-between;
  font-size: 13px;
  color: #909399;
}

.soldout-tag {
  width: 100%;
  margin-top: 10px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid #f0f0f0;

  &:last-child {
    border-bottom: none;
  }

  .label {
    color: #909399;
  }
}

.total-price {
  font-size: 20px;
  font-weight: 600;
  color: #f56c6c;
}
</style>
