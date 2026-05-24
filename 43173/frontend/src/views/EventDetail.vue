<template>
  <div class="event-detail" v-loading="loading">
    <div v-if="event" class="event-content">
      <div class="event-header">
        <el-image :src="event.cover_url" class="cover" fit="cover">
          <template #error>
            <div class="image-placeholder">
              <el-icon :size="64"><Calendar /></el-icon>
            </div>
          </template>
        </el-image>
        
        <div class="event-info">
          <h1 class="title">{{ event.title }}</h1>
          
          <div class="info-list">
            <div class="info-item">
              <el-icon><Location /></el-icon>
              <span>{{ event.venue }}</span>
            </div>
            <div class="info-item" v-if="event.address">
              <el-icon><MapLocation /></el-icon>
              <span>{{ event.address }}</span>
            </div>
            <div class="info-item">
              <el-icon><Clock /></el-icon>
              <span>{{ formatTime(event.start_time) }} - {{ formatTime(event.end_time) }}</span>
            </div>
            <div class="info-item" v-if="event.door_time">
              <el-icon><Bell /></el-icon>
              <span>入场时间: {{ formatTime(event.door_time) }}</span>
            </div>
          </div>
          
          <div class="stats">
            <div class="stat-item">
              <div class="value">¥{{ event.ticket_price }}</div>
              <div class="label">票价</div>
            </div>
            <div class="stat-item">
              <div class="value">{{ event.total_tickets - event.sold_tickets }}</div>
              <div class="label">剩余</div>
            </div>
            <div class="stat-item">
              <div class="value">{{ event.max_per_user }}</div>
              <div class="label">每人限购</div>
            </div>
          </div>
          
          <div class="actions">
            <el-button 
              type="primary" 
              size="large" 
              :disabled="event.sold_tickets >= event.total_tickets"
              @click="showPurchaseDialog = true"
            >
              {{ event.sold_tickets >= event.total_tickets ? '已售罄' : '立即购票' }}
            </el-button>
          </div>
        </div>
      </div>
      
      <div class="event-body">
        <h3>演出详情</h3>
        <div class="description" v-if="event.description">
          {{ event.description }}
        </div>
        <el-empty v-else description="暂无详情" />
      </div>
    </div>
    
    <el-dialog v-model="showPurchaseDialog" title="购票" width="500px">
      <div class="purchase-form">
        <el-form :model="purchaseForm" label-width="100px">
          <el-form-item label="票价">
            <span>¥{{ event?.ticket_price }}</span>
          </el-form-item>
          <el-form-item label="购买数量">
            <el-input-number 
              v-model="purchaseForm.quantity" 
              :min="1" 
              :max="event?.max_per_user || 4"
            />
          </el-form-item>
          <el-form-item label="合计">
            <span class="total">¥{{ (event?.ticket_price || 0) * purchaseForm.quantity }}</span>
          </el-form-item>
        </el-form>
      </div>
      <template #footer>
        <el-button @click="showPurchaseDialog = false">取消</el-button>
        <el-button type="primary" :loading="purchasing" @click="purchase">
          确认购票
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { eventApi } from '@/api/event'
import { useUserStore } from '@/stores/user'
import dayjs from 'dayjs'
import type { Event } from '@/api/event'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const event = ref<Event | null>(null)
const showPurchaseDialog = ref(false)
const purchasing = ref(false)
const purchaseForm = reactive({
  quantity: 1
})

onMounted(() => {
  loadEvent()
})

async function loadEvent() {
  loading.value = true
  try {
    const id = parseInt(route.params.id as string)
    event.value = await eventApi.getById(id)
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function purchase() {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }
  
  if (!event.value) return
  
  purchasing.value = true
  try {
    await eventApi.purchase({
      event_id: event.value.id,
      quantity: purchaseForm.quantity
    })
    ElMessage.success('购票成功')
    showPurchaseDialog.value = false
    loadEvent()
  } catch (e) {
    console.error(e)
  } finally {
    purchasing.value = false
  }
}

function formatTime(time: string) {
  return dayjs(time).format('YYYY-MM-DD HH:mm')
}
</script>

<style scoped lang="scss">
.event-detail {
  .event-content {
    .event-header {
      display: flex;
      gap: 30px;
      margin-bottom: 30px;
      
      .cover {
        width: 400px;
        height: 225px;
        border-radius: 8px;
        overflow: hidden;
        flex-shrink: 0;
      }
      
      .event-info {
        flex: 1;
        display: flex;
        flex-direction: column;
        
        .title {
          font-size: 28px;
          font-weight: 600;
          margin: 0 0 20px 0;
        }
        
        .info-list {
          margin-bottom: 20px;
          
          .info-item {
            display: flex;
            align-items: center;
            gap: 8px;
            padding: 8px 0;
            color: var(--text-light);
          }
        }
        
        .stats {
          display: flex;
          gap: 30px;
          margin-bottom: 20px;
          
          .stat-item {
            .value {
              font-size: 24px;
              font-weight: 600;
              color: var(--primary-color);
            }
            
            .label {
              font-size: 13px;
              color: var(--text-light);
            }
          }
        }
        
        .actions {
          margin-top: auto;
        }
      }
    }
    
    .event-body {
      padding: 20px;
      background: #f5f7fa;
      border-radius: 8px;
      
      h3 {
        margin: 0 0 16px 0;
      }
      
      .description {
        line-height: 1.8;
        white-space: pre-wrap;
      }
    }
  }
  
  .purchase-form {
    .total {
      font-size: 20px;
      font-weight: 600;
      color: var(--danger-color);
    }
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
