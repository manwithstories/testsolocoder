<template>
  <Layout>
    <div class="seat-select-page">
      <div class="seat-container">
        <div class="stage">
          <div class="stage-label">舞台</div>
        </div>

        <div class="seat-map" v-if="seats.length > 0">
          <div
            v-for="seat in seats"
            :key="seat.id"
            class="seat"
            :class="getSeatClass(seat)"
            :style="{
              position: 'absolute',
              left: seat.x + 'px',
              top: seat.y + 'px',
              width: seat.width + 'px',
              height: seat.height + 'px'
            }"
            @click="handleSeatClick(seat)"
          >
            {{ seat.col }}
          </div>
        </div>

        <el-empty v-else description="暂无座位数据" />
      </div>

      <div class="side-panel">
        <div class="legend">
          <h3>座位图例</h3>
          <div class="legend-item">
            <span class="seat-demo seat-available"></span>
            <span>可选</span>
          </div>
          <div class="legend-item">
            <span class="seat-demo seat-selected"></span>
            <span>已选</span>
          </div>
          <div class="legend-item">
            <span class="seat-demo seat-locked"></span>
            <span>已锁定</span>
          </div>
          <div class="legend-item">
            <span class="seat-demo seat-sold"></span>
            <span>已售</span>
          </div>
        </div>

        <div class="area-list" v-if="areas.length > 0">
          <h3>票价区域</h3>
          <div v-for="area in areas" :key="area.id" class="area-item">
            <span class="area-color" :style="{ backgroundColor: area.color }"></span>
            <span class="area-name">{{ area.name }}</span>
            <span class="area-price">¥{{ area.price }}</span>
          </div>
        </div>

        <div class="selected-info">
          <h3>已选座位 ({{ selectedSeats.size }})</h3>
          <div class="selected-seats" v-if="selectedSeats.size > 0">
            <el-tag
              v-for="seatId in Array.from(selectedSeats)"
              :key="seatId"
              closable
              @close="removeSeat(seatId)"
            >
              {{ getSeatLabel(seatId) }}
            </el-tag>
          </div>
          <p v-else class="empty-text">请点击座位进行选择</p>

          <div class="total-price">
            <span>总价：</span>
            <span class="price">¥{{ totalPrice.toFixed(2) }}</span>
          </div>

          <el-button
            type="primary"
            size="large"
            class="confirm-btn"
            :disabled="selectedSeats.size === 0 || locking"
            @click="handleConfirm"
            :loading="locking"
          >
            确认选座
          </el-button>
        </div>
      </div>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { seatApi } from '@/api'
import { useSeatStore } from '@/store'
import Layout from '@/components/Layout.vue'
import type { Seat, SeatArea } from '@/types'

const route = useRoute()
const router = useRouter()
const seatStore = useSeatStore()

const seats = ref<Seat[]>([])
const areas = ref<SeatArea[]>([])
const selectedSeats = ref<Set<number>>(new Set())
const locking = ref(false)

const areaMap = computed(() => {
  const map = new Map<number, SeatArea>()
  areas.value.forEach(area => map.set(area.id, area))
  return map
})

const totalPrice = computed(() => {
  let total = 0
  selectedSeats.value.forEach(seatId => {
    const seat = seats.value.find(s => s.id === seatId)
    if (seat) {
      const area = areaMap.value.get(seat.area_id)
      if (area) {
        total += area.price
      }
    }
  })
  return total
})

function getSeatClass(seat: Seat) {
  const classes: string[] = []
  if (selectedSeats.value.has(seat.id)) {
    classes.push('seat-selected')
  } else if (seat.status === 2) {
    classes.push('seat-sold')
  } else if (seat.status === 1 || seatStore.lockedSeats.has(seat.id)) {
    classes.push('seat-locked')
  } else {
    classes.push('seat-available')
  }
  return classes
}

function handleSeatClick(seat: Seat) {
  if (seat.status === 2 || seat.status === 1 || seatStore.lockedSeats.has(seat.id)) {
    return
  }
  if (selectedSeats.value.has(seat.id)) {
    selectedSeats.value.delete(seat.id)
  } else {
    selectedSeats.value.add(seat.id)
  }
}

function removeSeat(seatId: number) {
  selectedSeats.value.delete(seatId)
}

function getSeatLabel(seatId: number) {
  const seat = seats.value.find(s => s.id === seatId)
  return seat ? `${seat.row}排${seat.col}座` : ''
}

async function handleConfirm() {
  if (selectedSeats.value.size === 0) {
    ElMessage.warning('请先选择座位')
    return
  }

  const sessionId = Number(route.params.sessionId as string)
  const seatIds = Array.from(selectedSeats.value)

  locking.value = true
  try {
    await seatApi.lock({
      session_id: sessionId,
      seat_ids: seatIds
    })
    seatStore.lockSeats(seatIds)
    seatStore.selectedSeats.clear()
    seatIds.forEach(id => seatStore.selectedSeats.add(id))

    ElMessage.success('座位锁定成功，请在15分钟内完成支付')
    router.push({
      path: '/order-confirm',
      query: {
        sessionId: sessionId.toString(),
        seatIds: seatIds.join(','),
        totalPrice: totalPrice.value.toString()
      }
    })
  } catch (err: any) {
    ElMessage.error(err.message || '锁定失败，请重试')
  } finally {
    locking.value = false
  }
}

async function fetchSeats() {
  const sessionId = Number(route.params.sessionId as string)
  try {
    seats.value = await seatApi.getSeats(sessionId)
    areas.value = await seatApi.getAreas(sessionId)
  } catch (err) {
    console.error(err)
  }
}

let refreshTimer: number | null = null

onMounted(() => {
  fetchSeats()
  refreshTimer = window.setInterval(fetchSeats, 30000)
})

onUnmounted(() => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
  }
})
</script>

<style lang="scss" scoped>
.seat-select-page {
  max-width: 1400px;
  margin: 0 auto;
  padding: 20px;
  display: flex;
  gap: 24px;
}

.seat-container {
  flex: 1;
  background: white;
  padding: 30px;
  border-radius: 12px;
  min-height: 600px;
}

.stage {
  width: 60%;
  height: 50px;
  background: linear-gradient(180deg, #409eff, #66b1ff);
  margin: 0 auto 40px;
  border-radius: 0 0 50% 50%;
  display: flex;
  align-items: center;
  justify-content: center;

  .stage-label {
    color: white;
    font-weight: 600;
    font-size: 16px;
  }
}

.seat-map {
  position: relative;
  margin: 0 auto;
  min-height: 500px;
}

.seat {
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 10px;
  color: white;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
  user-select: none;

  &:hover {
    transform: scale(1.1);
    z-index: 10;
  }
}

.side-panel {
  width: 300px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.legend, .area-list, .selected-info {
  background: white;
  padding: 20px;
  border-radius: 12px;

  h3 {
    margin: 0 0 16px 0;
    font-size: 16px;
  }
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
  font-size: 14px;
  color: #666;
}

.seat-demo {
  width: 20px;
  height: 20px;
  border-radius: 4px;
}

.area-item {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 8px;
  font-size: 14px;
}

.area-color {
  width: 16px;
  height: 16px;
  border-radius: 3px;
}

.area-price {
  margin-left: auto;
  color: #f56c6c;
  font-weight: 600;
}

.selected-info {
  .empty-text {
    color: #999;
    text-align: center;
    padding: 20px 0;
  }

  .selected-seats {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    margin-bottom: 16px;
  }

  .total-price {
    padding: 16px 0;
    border-top: 1px solid #eee;
    font-size: 16px;
    display: flex;
    justify-content: space-between;
    align-items: center;

    .price {
      color: #f56c6c;
      font-size: 24px;
      font-weight: bold;
    }
  }

  .confirm-btn {
    width: 100%;
    height: 48px;
    font-size: 16px;
  }
}
</style>
