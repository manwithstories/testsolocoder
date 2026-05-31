<template>
  <div class="my-ships">
    <div class="page-header">
      <h2 class="page-title">我的船只</h2>
      <el-button type="primary" @click="$router.push('/ship-create')">
        <el-icon><Plus /></el-icon>
        发布新船
      </el-button>
    </div>

    <div class="card-container">
      <el-table :data="ships" v-loading="loading" style="width: 100%">
        <el-table-column label="船只" min-width="200">
          <template #default="{ row }">
            <div class="ship-cell">
              <el-avatar :size="48" :src="row.images?.[0]?.url">
                {{ row.name?.[0] }}
              </el-avatar>
              <div class="ship-info">
                <div class="ship-name">{{ row.name }}</div>
                <div class="ship-type">
                  <el-tag size="small">{{ getShipTypeText(row.ship_type) }}</el-tag>
                </div>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="价格">
          <template #default="{ row }">
            ¥{{ row.daily_rate }}/天
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'available' ? 'success' : 'warning'" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="评分" width="120">
          <template #default="{ row }">
            <el-rate :model-value="row.average_rating" disabled size="small" />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="$router.push(`/ships/${row.id}`)">查看</el-button>
            <el-button type="primary" link @click="$router.push(`/ships/${row.id}/edit`)">编辑</el-button>
            <el-popconfirm title="确定删除这艘船吗？" @confirm="handleDelete(row.id)">
              <template #reference>
                <el-button type="danger" link>删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <el-empty v-if="ships.length === 0 && !loading" description="您还没有发布任何船只" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getMyShipsApi, deleteShipApi } from '@/api/ship'
import type { Ship } from '@/types/ship'

const loading = ref(false)
const ships = ref<Ship[]>([])

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

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    available: '可用',
    rented: '已出租',
    maintenance: '维护中',
    inactive: '已下架'
  }
  return map[status] || status
}

const fetchShips = async () => {
  loading.value = true
  try {
    const res: any = await getMyShipsApi()
    ships.value = res.data || []
  } catch (error) {
    console.error('Failed to fetch ships:', error)
  } finally {
    loading.value = false
  }
}

const handleDelete = async (id: string) => {
  try {
    await deleteShipApi(id)
    ElMessage.success('删除成功')
    fetchShips()
  } catch (error) {
    ElMessage.error('删除失败')
  }
}

onMounted(fetchShips)
</script>

<style lang="scss" scoped>
.my-ships {
  .card-container {
    background: #fff;
    border-radius: 8px;
    padding: 20px;
  }

  .ship-cell {
    display: flex;
    align-items: center;
    gap: 12px;

    .ship-info {
      .ship-name {
        font-weight: 500;
      }

      .ship-type {
        margin-top: 4px;
      }
    }
  }
}
</style>
