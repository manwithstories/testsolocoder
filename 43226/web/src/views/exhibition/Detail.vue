<template>
  <div class="exhibition-detail page-container" v-loading="loading">
    <div v-if="exhibition" class="detail-content">
      <div class="header-section card-shadow">
        <div class="header-image">
          <img :src="exhibition.image_url || '/placeholder.svg'" :alt="exhibition.title" />
        </div>
        <div class="header-info">
          <div class="flex-between mb-20">
            <h1>{{ exhibition.title }}</h1>
            <el-tag :type="statusType">{{ statusText }}</el-tag>
          </div>
          <p class="description">{{ exhibition.description }}</p>
          <div class="info-grid">
            <div class="info-item">
              <el-icon size="20" color="#409EFF"><Calendar /></el-icon>
              <div>
                <span class="label">展览时间</span>
                <span class="value">{{ formatDate(exhibition.start_date) }} - {{ formatDate(exhibition.end_date) }}</span>
              </div>
            </div>
            <div class="info-item">
              <el-icon size="20" color="#67c23a"><Location /></el-icon>
              <div>
                <span class="label">展览地点</span>
                <span class="value">{{ exhibition.location }} {{ exhibition.hall_number }}</span>
              </div>
            </div>
            <div class="info-item">
              <el-icon size="20" color="#e6a23c"><Tickets /></el-icon>
              <div>
                <span class="label">门票价格</span>
                <span class="value price" v-if="exhibition.ticket_price > 0">¥{{ exhibition.ticket_price }}</span>
                <span class="value free" v-else>免费</span>
              </div>
            </div>
            <div class="info-item">
              <el-icon size="20" color="#f56c6c"><User /></el-icon>
              <div>
                <span class="label">容纳人数</span>
                <span class="value">每场 {{ exhibition.max_visitors }} 人</span>
              </div>
            </div>
          </div>
          <div class="action-buttons">
            <el-button type="primary" size="large" @click="showReservation = true" :disabled="exhibition.status !== 'published'">
              立即预约
            </el-button>
            <el-button size="large" v-if="exhibition.is_virtual">
              <el-icon><VideoCamera /></el-icon> 虚拟展厅
            </el-button>
            <el-button size="large" @click="scrollToCollections">
              查看藏品
            </el-button>
          </div>
          <div class="stats-row">
            <span><el-icon><View /></el-icon> {{ exhibition.view_count }} 浏览</span>
          </div>
        </div>
      </div>

      <div id="collections-section" class="collections-section card-shadow p-20">
        <h2 class="section-title">展览藏品 ({{ collections.length }})</h2>
        <div class="collection-grid">
          <div
            v-for="item in collections"
            :key="item.id"
            class="collection-item"
            @click="$router.push(`/collections/${item.id}`)"
          >
            <img :src="item.image_url || '/placeholder.svg'" :alt="item.name" />
            <div class="collection-info">
              <h4>{{ item.name }}</h4>
              <p>{{ item.era }}</p>
            </div>
          </div>
        </div>
      </div>

      <div class="time-slots-section card-shadow p-20">
        <h2 class="section-title">预约时段</h2>
        <div class="date-selector">
          <el-date-picker
            v-model="selectedDate"
            type="date"
            placeholder="选择日期"
            value-format="YYYY-MM-DD"
            @change="fetchTimeSlots"
          />
        </div>
        <div v-loading="slotLoading" class="time-slot-grid">
          <div
            v-for="slot in timeSlots"
            :key="slot.id"
            class="time-slot-item"
            :class="{ full: slot.booked_count >= slot.max_capacity }"
            @click="selectSlot(slot)"
          >
            <div class="time-range">{{ slot.start_time }} - {{ slot.end_time }}</div>
            <div class="slot-count">{{ slot.booked_count }}/{{ slot.max_capacity }} 人</div>
            <div class="slot-status" v-if="slot.booked_count >= slot.max_capacity">已约满</div>
            <div class="slot-status" v-else>可预约</div>
          </div>
        </div>
      </div>
    </div>

    <el-dialog v-model="showReservation" title="预约参观" width="500px">
      <el-form v-if="exhibition" :model="reservationForm" label-width="100px">
        <el-form-item label="展览名称">
          <span>{{ exhibition.title }}</span>
        </el-form-item>
        <el-form-item label="预约日期">
          <el-date-picker
            v-model="reservationForm.date"
            type="date"
            placeholder="选择日期"
            value-format="YYYY-MM-DD"
            style="width: 100%"
            @change="fetchTimeSlots"
          />
        </el-form-item>
        <el-form-item label="选择时段">
          <el-select v-model="reservationForm.time_slot_id" placeholder="请选择时段" style="width: 100%">
            <el-option
              v-for="slot in timeSlots"
              :key="slot.id"
              :label="`${slot.start_time} - ${slot.end_time} (${slot.booked_count}/${slot.max_capacity})`"
              :value="slot.id"
              :disabled="slot.booked_count >= slot.max_capacity"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="参观人数">
          <el-input-number
            v-model="reservationForm.visitor_count"
            :min="1"
            :max="10"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="导览服务">
          <el-radio-group v-model="reservationForm.guide_type">
            <el-radio value="standard">标准参观</el-radio>
            <el-radio value="audio">语音导览</el-radio>
            <el-radio value="human">人工讲解</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="预计费用">
          <span class="total-price">¥{{ totalPrice }}</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showReservation = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitReservation">确认预约</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import * as exhibitionApi from '@/api/exhibition'
import * as reservationApi from '@/api/reservation'
import type { Exhibition, Collection, TimeSlot } from '@/types'
import { useUserStore } from '@/stores/user'
import dayjs from 'dayjs'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const id = Number(route.params.id)

const loading = ref(false)
const slotLoading = ref(false)
const submitting = ref(false)
const exhibition = ref<Exhibition | null>(null)
const collections = ref<Collection[]>([])
const timeSlots = ref<TimeSlot[]>([])
const showReservation = ref(false)
const selectedDate = ref(dayjs().format('YYYY-MM-DD'))

const reservationForm = reactive({
  date: dayjs().format('YYYY-MM-DD'),
  time_slot_id: 0,
  visitor_count: 1,
  guide_type: 'standard'
})

const statusType = computed(() => {
  if (exhibition.value?.status === 'published') return 'success'
  if (exhibition.value?.status === 'closed') return 'info'
  return 'warning'
})

const statusText = computed(() => {
  if (exhibition.value?.status === 'published') return '展览中'
  if (exhibition.value?.status === 'closed') return '已结束'
  return '草稿'
})

const totalPrice = computed(() => {
  const price = exhibition.value?.ticket_price || 0
  const guideExtra = reservationForm.guide_type === 'human' ? 50 : reservationForm.guide_type === 'audio' ? 10 : 0
  return (price + guideExtra) * reservationForm.visitor_count
})

const formatDate = (date: string) => {
  return dayjs(date).format('YYYY年MM月DD日')
}

const scrollToCollections = () => {
  document.getElementById('collections-section')?.scrollIntoView({ behavior: 'smooth' })
}

const selectSlot = (slot: TimeSlot) => {
  if (slot.booked_count >= slot.max_capacity) {
    ElMessage.warning('该时段已约满')
    return
  }
  reservationForm.time_slot_id = slot.id
}

const fetchDetail = async () => {
  try {
    loading.value = true
    const res = await exhibitionApi.getExhibition(id)
    exhibition.value = res.data
    if (res.data.collections) {
      collections.value = res.data.collections
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const fetchCollections = async () => {
  try {
    const res = await exhibitionApi.getExhibitionCollections(id)
    collections.value = res.data
  } catch (e) {
    console.error(e)
  }
}

const fetchTimeSlots = async () => {
  if (!selectedDate.value) return
  try {
    slotLoading.value = true
    const res = await exhibitionApi.listTimeSlots(id, selectedDate.value)
    timeSlots.value = res.data
  } catch (e) {
    console.error(e)
  } finally {
    slotLoading.value = false
  }
}

const submitReservation = async () => {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }
  if (!reservationForm.time_slot_id) {
    ElMessage.warning('请选择预约时段')
    return
  }
  try {
    submitting.value = true
    await reservationApi.createReservation({
      exhibition_id: id,
      time_slot_id: reservationForm.time_slot_id,
      visitor_count: reservationForm.visitor_count,
      guide_type: reservationForm.guide_type
    })
    ElMessage.success('预约成功')
    showReservation.value = false
    router.push('/dashboard/my-reservations')
  } catch (e) {
    console.error(e)
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  fetchDetail()
  fetchCollections()
  fetchTimeSlots()
})
</script>

<style scoped lang="scss">
.exhibition-detail {
  max-width: 1400px;
  margin: 0 auto;

  .header-section {
    display: flex;
    gap: 40px;
    padding: 30px;
    margin-bottom: 24px;
    border-radius: 12px;

    .header-image {
      width: 45%;
      flex-shrink: 0;

      img {
        width: 100%;
        height: 350px;
        object-fit: cover;
        border-radius: 8px;
      }
    }

    .header-info {
      flex: 1;
      display: flex;
      flex-direction: column;

      h1 {
        font-size: 32px;
        margin: 0;
      }

      .description {
        color: #606266;
        line-height: 1.8;
        margin-bottom: 24px;
      }

      .info-grid {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 20px;
        margin-bottom: 24px;

        .info-item {
          display: flex;
          align-items: flex-start;
          gap: 12px;

          .label {
            display: block;
            font-size: 13px;
            color: #909399;
            margin-bottom: 4px;
          }

          .value {
            font-size: 15px;
            font-weight: 500;

            &.price {
              color: #f56c6c;
            }

            &.free {
              color: #67c23a;
            }
          }
        }
      }

      .action-buttons {
        display: flex;
        gap: 12px;
        margin-bottom: 16px;
      }

      .stats-row {
        color: #909399;
        font-size: 13px;

        span {
          display: flex;
          align-items: center;
          gap: 4px;
        }
      }
    }
  }

  .section-title {
    font-size: 22px;
    margin-bottom: 20px;
    padding-bottom: 12px;
    border-bottom: 1px solid #ebeef5;
  }

  .collection-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
    gap: 16px;

    .collection-item {
      cursor: pointer;
      border-radius: 8px;
      overflow: hidden;
      transition: transform 0.3s;

      &:hover {
        transform: translateY(-2px);
      }

      img {
        width: 100%;
        height: 140px;
        object-fit: cover;
      }

      .collection-info {
        padding: 10px;
        background: #f5f7fa;

        h4 {
          font-size: 14px;
          margin-bottom: 4px;
          white-space: nowrap;
          overflow: hidden;
          text-overflow: ellipsis;
        }

        p {
          font-size: 12px;
          color: #909399;
          margin: 0;
        }
      }
    }
  }

  .time-slots-section {
    margin-top: 24px;
    border-radius: 12px;

    .date-selector {
      margin-bottom: 20px;
    }

    .time-slot-grid {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
      gap: 16px;

      .time-slot-item {
        border: 2px solid #ebeef5;
        border-radius: 8px;
        padding: 16px;
        text-align: center;
        cursor: pointer;
        transition: all 0.3s;

        &:hover {
          border-color: #409eff;
        }

        &.full {
          opacity: 0.6;
          cursor: not-allowed;

          &:hover {
            border-color: #ebeef5;
          }
        }

        .time-range {
          font-size: 16px;
          font-weight: 500;
          margin-bottom: 8px;
        }

        .slot-count {
          font-size: 13px;
          color: #909399;
          margin-bottom: 8px;
        }

        .slot-status {
          font-size: 12px;
          color: #67c23a;
        }

        &.full .slot-status {
          color: #f56c6c;
        }
      }
    }
  }

  .total-price {
    font-size: 24px;
    font-weight: 600;
    color: #f56c6c;
  }
}

@media (max-width: 768px) {
  .header-section {
    flex-direction: column;

    .header-image {
      width: 100%;
    }
  }

  .info-grid {
    grid-template-columns: 1fr !important;
  }
}
</style>
