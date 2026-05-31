<template>
  <div class="ship-detail" v-if="ship">
    <el-page-header @back="$router.back()" :content="ship.name" class="page-header" />

    <el-row :gutter="24">
      <el-col :xs="24" :md="16">
        <div class="card-container">
          <el-carousel :interval="4000" height="400px">
            <el-carousel-item v-for="image in ship.images" :key="image.id">
              <img :src="image.url" :alt="ship.name" class="ship-image" />
            </el-carousel-item>
            <el-carousel-item v-if="!ship.images || ship.images.length === 0">
              <div class="no-image">暂无图片</div>
            </el-carousel-item>
          </el-carousel>
        </div>

        <div class="card-container" style="margin-top: 16px">
          <h3>船只详情</h3>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="船型">{{ getShipTypeText(ship.ship_type) }}</el-descriptions-item>
            <el-descriptions-item label="载客量">{{ ship.capacity }}人</el-descriptions-item>
            <el-descriptions-item label="客舱数">{{ ship.cabin_count }}间</el-descriptions-item>
            <el-descriptions-item label="卫生间">{{ ship.bathroom_count }}间</el-descriptions-item>
            <el-descriptions-item label="船长">{{ ship.length }}米</el-descriptions-item>
            <el-descriptions-item label="船宽">{{ ship.width }}米</el-descriptions-item>
            <el-descriptions-item label="建造年份">{{ ship.year_built }}</el-descriptions-item>
            <el-descriptions-item label="母港">{{ ship.home_port }}</el-descriptions-item>
          </el-descriptions>

          <h4 style="margin-top: 20px">设备配置</h4>
          <p>{{ ship.equipment || '暂无信息' }}</p>

          <h4 style="margin-top: 20px">特色功能</h4>
          <p>{{ ship.features || '暂无信息' }}</p>

          <h4 style="margin-top: 20px">航行区域</h4>
          <p>{{ ship.sailing_area || '暂无信息' }}</p>

          <h4 style="margin-top: 20px">取消政策</h4>
          <p>{{ ship.cancellation_policy || '暂无信息' }}</p>
        </div>

        <div class="card-container" style="margin-top: 16px">
          <h3>用户评价</h3>
          <div class="rating-summary">
            <el-rate :model-value="ship.average_rating" disabled />
            <span class="rating-text">{{ ship.average_rating }}分 / {{ ship.review_count }}条评价</span>
          </div>
        </div>
      </el-col>

      <el-col :xs="24" :md="8">
        <div class="card-container sticky">
          <div class="price-box">
            <div class="price">
              <span class="currency">¥</span>
              <span class="amount">{{ ship.daily_rate }}</span>
              <span class="unit">/天</span>
            </div>
            <div class="hourly-price">或 ¥{{ ship.hourly_rate }}/小时</div>
          </div>

          <el-divider />

          <div class="booking-form">
            <el-form :model="bookingForm" label-position="top">
              <el-form-item label="租赁类型">
                <el-select v-model="bookingForm.rental_type" style="width: 100%">
                  <el-option label="按天租赁" value="daily" />
                  <el-option label="按小时租赁" value="hourly" />
                  <el-option label="航程租赁" value="voyage" />
                </el-select>
              </el-form-item>
              <el-form-item label="开始日期">
                <el-date-picker
                  v-model="bookingForm.start_date"
                  type="datetime"
                  placeholder="选择开始时间"
                  style="width: 100%"
                />
              </el-form-item>
              <el-form-item label="结束日期">
                <el-date-picker
                  v-model="bookingForm.end_date"
                  type="datetime"
                  placeholder="选择结束时间"
                  style="width: 100%"
                />
              </el-form-item>
              <el-form-item label="保险">
                <el-select v-model="bookingForm.insurance_type" style="width: 100%">
                  <el-option label="无保险" value="none" />
                  <el-option label="基础保险" value="basic" />
                  <el-option label="高级保险" value="premium" />
                </el-select>
              </el-form-item>
            </el-form>

            <el-button
              type="primary"
              size="large"
              style="width: 100%"
              :disabled="!userStore.isLoggedIn"
              @click="handleBook"
            >
              {{ userStore.isLoggedIn ? '立即预订' : '请先登录' }}
            </el-button>
          </div>

          <el-divider />

          <div class="owner-info">
            <div class="owner-label">船东信息</div>
            <el-avatar :size="48" :src="ship.owner?.avatar_url">
              {{ ship.owner?.full_name?.[0] }}
            </el-avatar>
            <div class="owner-name">{{ ship.owner?.full_name }}</div>
          </div>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getShipApi } from '@/api/ship'
import { createRentalApi } from '@/api/rental'
import { useUserStore } from '@/stores/user'
import type { Ship } from '@/types/ship'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const ship = ref<Ship | null>(null)

const bookingForm = reactive({
  rental_type: 'daily',
  start_date: '',
  end_date: '',
  insurance_type: 'none'
})

const getShipTypeText = (type: string) => {
  const map: Record<string, string> = {
    sailboat: '帆船',
    motorboat: '摩托艇',
    yacht: '游艇',
    fishing: '渔船',
    cargo: '货船'
  }
  return map[type] || type
}

const handleBook = () => {
  if (!bookingForm.start_date || !bookingForm.end_date) {
    ElMessage.warning('请选择租赁时间')
    return
  }
  router.push({
    path: '/rental-create',
    query: {
      ship_id: ship.value?.id,
      start_date: bookingForm.start_date,
      end_date: bookingForm.end_date,
      rental_type: bookingForm.rental_type,
      insurance_type: bookingForm.insurance_type
    }
  })
}

onMounted(async () => {
  const id = route.params.id as string
  try {
    const res: any = await getShipApi(id)
    ship.value = res.data
  } catch (error) {
    console.error('Failed to fetch ship:', error)
  }
})
</script>

<style lang="scss" scoped>
.ship-detail {
  .page-header {
    margin-bottom: 16px;
    padding: 12px 20px;
    background: #fff;
    border-radius: 8px;
  }

  .card-container {
    background: #fff;
    border-radius: 8px;
    padding: 20px;

    .ship-image {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }

    .no-image {
      width: 100%;
      height: 100%;
      display: flex;
      align-items: center;
      justify-content: center;
      background: #f5f5f5;
      color: #999;
    }
  }

  .rating-summary {
    display: flex;
    align-items: center;
    gap: 12px;

    .rating-text {
      color: rgba(0, 0, 0, 0.45);
    }
  }

  .sticky {
    position: sticky;
    top: 24px;
  }

  .price-box {
    text-align: center;

    .price {
      color: #1890ff;
      font-weight: 600;

      .currency {
        font-size: 16px;
      }

      .amount {
        font-size: 32px;
      }

      .unit {
        font-size: 14px;
        color: rgba(0, 0, 0, 0.45);
      }
    }

    .hourly-price {
      color: rgba(0, 0, 0, 0.45);
      font-size: 14px;
    }
  }

  .owner-info {
    text-align: center;

    .owner-label {
      color: rgba(0, 0, 0, 0.45);
      margin-bottom: 8px;
    }

    .owner-name {
      margin-top: 8px;
      font-weight: 500;
    }
  }
}
</style>
