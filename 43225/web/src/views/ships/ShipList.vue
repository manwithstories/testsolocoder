<template>
  <div class="ship-list">
    <div class="page-header">
      <h2 class="page-title">船舶列表</h2>
      <el-button type="primary" @click="$router.push('/ship-create')" v-if="userStore.hasRole(['owner', 'admin'])">
        <el-icon><Plus /></el-icon>
        发布船只
      </el-button>
    </div>

    <div class="filter-bar">
      <el-select v-model="filters.ship_type" placeholder="船型" clearable style="width: 140px">
        <el-option label="帆船" value="sailboat" />
        <el-option label="摩托艇" value="motorboat" />
        <el-option label="游艇" value="yacht" />
        <el-option label="渔船" value="fishing" />
        <el-option label="货船" value="cargo" />
      </el-select>
      <el-input v-model="filters.location" placeholder="地点" clearable style="width: 160px" :prefix-icon="Location" />
      <el-input-number v-model="filters.min_price" :min="0" placeholder="最低价格" style="width: 140px" />
      <el-input-number v-model="filters.max_price" :min="0" placeholder="最高价格" style="width: 140px" />
      <el-select v-model="filters.sort_by" placeholder="排序" clearable style="width: 140px">
        <el-option label="评分最高" value="rating" />
        <el-option label="价格从低到高" value="price_asc" />
        <el-option label="价格从高到低" value="price_desc" />
        <el-option label="最新发布" value="created_at" />
      </el-select>
      <el-button type="primary" @click="fetchShips">
        <el-icon><Search /></el-icon>
        搜索
      </el-button>
    </div>

    <el-row :gutter="16">
      <el-col :xs="24" :sm="12" :md="8" :lg="6" v-for="ship in ships" :key="ship.id">
        <el-card class="ship-card" shadow="hover" @click="$router.push(`/ships/${ship.id}`)">
          <template #cover>
            <div class="ship-image">
              <img
                :src="ship.images?.find(i => i.is_primary)?.url || ship.images?.[0]?.url || '/placeholder-ship.jpg'"
                :alt="ship.name"
              />
            </div>
          </template>
          <div class="ship-info">
            <div class="ship-name">{{ ship.name }}</div>
            <div class="ship-type">
              <el-tag size="small">{{ getShipTypeText(ship.ship_type) }}</el-tag>
              <el-rate :model-value="ship.average_rating" disabled size="small" />
            </div>
            <div class="ship-specs">
              <span><el-icon><User /></el-icon> {{ ship.capacity }}人</span>
              <span><el-icon><HomeFilled /></el-icon> {{ ship.home_port }}</span>
            </div>
            <div class="ship-price">
              <span class="price">¥{{ ship.daily_rate }}</span>
              <span class="unit">/天</span>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-empty v-if="ships.length === 0 && !loading" description="暂无船只" />

    <div class="pagination">
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.page_size"
        :total="pagination.total"
        layout="total, prev, pager, next"
        @current-change="fetchShips"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { getShipsApi } from '@/api/ship'
import { useUserStore } from '@/stores/user'
import type { Ship } from '@/types/ship'

const userStore = useUserStore()
const loading = ref(false)
const ships = ref<Ship[]>([])

const filters = reactive({
  ship_type: '',
  location: '',
  min_price: undefined as number | undefined,
  max_price: undefined as number | undefined,
  sort_by: 'rating'
})

const pagination = reactive({
  page: 1,
  page_size: 12,
  total: 0
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

const fetchShips = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.page_size,
      sort_by: filters.sort_by
    }
    if (filters.ship_type) params.ship_type = filters.ship_type
    if (filters.location) params.location = filters.location
    if (filters.min_price) params.min_price = filters.min_price
    if (filters.max_price) params.max_price = filters.max_price

    const res: any = await getShipsApi(params)
    ships.value = res.data || []
    pagination.total = res.total || 0
  } catch (error) {
    console.error('Failed to fetch ships:', error)
  } finally {
    loading.value = false
  }
}

onMounted(fetchShips)
</script>

<style lang="scss" scoped>
.ship-list {
  .ship-card {
    margin-bottom: 16px;
    cursor: pointer;
    transition: transform 0.3s;

    &:hover {
      transform: translateY(-4px);
    }

    .ship-image {
      height: 180px;
      overflow: hidden;

      img {
        width: 100%;
        height: 100%;
        object-fit: cover;
      }
    }

    .ship-info {
      .ship-name {
        font-size: 16px;
        font-weight: 600;
        margin-bottom: 8px;
      }

      .ship-type {
        display: flex;
        align-items: center;
        gap: 8px;
        margin-bottom: 8px;
      }

      .ship-specs {
        color: rgba(0, 0, 0, 0.45);
        font-size: 13px;
        margin-bottom: 8px;
        display: flex;
        gap: 12px;

        span {
          display: flex;
          align-items: center;
          gap: 4px;
        }
      }

      .ship-price {
        color: #1890ff;
        font-weight: 600;

        .price {
          font-size: 18px;
        }

        .unit {
          font-size: 12px;
          color: rgba(0, 0, 0, 0.45);
        }
      }
    }
  }

  .pagination {
    display: flex;
    justify-content: center;
    margin-top: 24px;
  }
}
</style>
