<template>
  <div class="page-container">
    <div class="page-header">
      <h1 class="page-title">我的设备</h1>
      <el-button type="primary" @click="goToCreate">
        <el-icon><Plus /></el-icon>
        添加设备
      </el-button>
    </div>

    <div v-loading="loading" class="equipment-list">
      <el-table :data="equipments" stripe style="width: 100%">
        <el-table-column label="设备名称" prop="name" min-width="150">
          <template #default="{ row }">
            <div class="equipment-cell" @click="goToDetail(row.id)">
              <el-avatar
                v-if="row.images && row.images.length > 0"
                :size="48"
                :src="`/uploads/${row.images[0].imageUrl}`"
                shape="square"
              />
              <el-avatar v-else :size="48" shape="square">
                <el-icon><Camera /></el-icon>
              </el-avatar>
              <span class="equipment-name">{{ row.name }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="品牌型号" min-width="150">
          <template #default="{ row }">
            {{ row.brand }} {{ row.model }}
          </template>
        </el-table-column>
        <el-table-column label="分类" prop="category" width="100" />
        <el-table-column label="日租金" width="100">
          <template #default="{ row }">
            <span class="price">¥{{ row.dailyRent }}</span>
          </template>
        </el-table-column>
        <el-table-column label="押金" width="100">
          <template #default="{ row }">
            ¥{{ row.deposit }}
          </template>
        </el-table-column>
        <el-table-column label="评分" width="120">
          <template #default="{ row }">
            <el-rate :model-value="row.rating" disabled size="small" />
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="goToEdit(row.id)">编辑</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div v-if="!loading && equipments.length === 0" class="empty-state">
        <el-icon><Box /></el-icon>
        <p>暂无设备</p>
        <el-button type="primary" @click="goToCreate">立即添加</el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { equipmentApi } from '@/api/equipment'
import type { Equipment } from '@/types'

const router = useRouter()
const loading = ref(false)
const equipments = ref<Equipment[]>([])

onMounted(() => {
  loadEquipments()
})

async function loadEquipments() {
  loading.value = true
  try {
    const response = await equipmentApi.getMyEquipments()
    equipments.value = response.data
  } catch (error) {
    console.error('Failed to load equipments:', error)
    ElMessage.error('加载设备列表失败')
  } finally {
    loading.value = false
  }
}

function goToCreate() {
  router.push('/equipments/create')
}

function goToDetail(id: number) {
  router.push(`/equipments/${id}`)
}

function goToEdit(id: number) {
  router.push(`/equipments/${id}/edit`)
}

async function handleDelete(row: Equipment) {
  try {
    await ElMessageBox.confirm(
      `确定要删除设备"${row.name}"吗？此操作不可恢复。`,
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    await equipmentApi.deleteEquipment(row.id)
    ElMessage.success('删除成功')
    loadEquipments()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete equipment:', error)
    }
  }
}

function getStatusType(status: string) {
  const typeMap: Record<string, string> = {
    available: 'success',
    rented: 'warning',
    maintenance: 'info'
  }
  return typeMap[status] || 'info'
}

function getStatusText(status: string) {
  const textMap: Record<string, string> = {
    available: '可出租',
    rented: '已出租',
    maintenance: '维护中'
  }
  return textMap[status] || status
}
</script>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.equipment-cell {
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
}

.equipment-name {
  font-weight: 500;
}

.price {
  color: #f56c6c;
  font-weight: 500;
}
</style>
