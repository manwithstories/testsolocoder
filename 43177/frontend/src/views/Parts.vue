<template>
  <div class="parts-page">
    <div class="container">
      <div class="page-header">
        <h2 class="page-title">配件管理</h2>
        <el-button
          v-if="userStore.isTechnician"
          type="primary"
          @click="router.push('/part-requests')"
        >
          申请配件
        </el-button>
      </div>

      <el-card v-if="lowStockCount > 0" class="warning-card" type="warning">
        <el-icon><Warning /></el-icon>
        <span>当前有 {{ lowStockCount }} 种配件库存不足</span>
      </el-card>

      <div class="filter-bar">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索配件"
          clearable
          style="width: 200px;"
          @input="handleSearch"
        />
        <el-select v-model="categoryFilter" placeholder="选择分类" clearable style="width: 150px;">
          <el-option label="家电配件" value="家电配件" />
          <el-option label="数码配件" value="数码配件" />
          <el-option label="汽车配件" value="汽车配件" />
        </el-select>
      </div>

      <div class="parts-grid">
        <el-card v-for="part in filteredParts" :key="part.id" class="part-card" shadow="hover">
          <div class="part-image">
            <img v-if="part.image" :src="part.image" alt="part" />
            <el-icon v-else :size="60"><Tools /></el-icon>
          </div>
          <div class="part-info">
            <h3 class="part-name">{{ part.name }}</h3>
            <div class="part-category">{{ part.category }}</div>
            <div class="part-price">¥{{ part.price }}</div>
            <div class="part-stock">
              库存：<span :class="{ 'low-stock': part.stock <= part.min_stock }">{{ part.stock }}</span>
            </div>
            <el-button
              v-if="userStore.isTechnician"
              type="primary"
              size="small"
              :disabled="part.stock <= 0"
              @click="usePart(part)"
            >
              使用
            </el-button>
          </div>
        </el-card>
      </div>

      <div v-if="filteredParts.length === 0" class="empty-state">
        <el-empty description="暂无配件" />
      </div>

      <div v-if="total > pageSize" class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          :page-size="pageSize"
          :total="total"
          layout="prev, pager, next"
          @current-change="loadParts"
        />
      </div>
    </div>

    <el-dialog v-model="showUseDialog" title="使用配件" width="400px">
      <el-form :model="useForm" label-width="80px">
        <el-form-item label="工单号">
          <el-input v-model="useForm.order_id" placeholder="请输入工单号" />
        </el-form-item>
        <el-form-item label="数量">
          <el-input-number v-model="useForm.quantity" :min="1" :max="selectedPart?.stock || 1" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showUseDialog = false">取消</el-button>
        <el-button type="primary" @click="submitUsePart" :loading="submitting">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Tools, Warning } from '@element-plus/icons-vue'
import { useUserStore } from '@/store/user'
import { partApi } from '@/api/part'
import type { Part } from '@/types'

const router = useRouter()
const userStore = useUserStore()

const parts = ref<Part[]>([])
const searchKeyword = ref('')
const categoryFilter = ref('')
const currentPage = ref(1)
const pageSize = ref(12)
const total = ref(0)
const lowStockCount = ref(0)
const showUseDialog = ref(false)
const submitting = ref(false)
const selectedPart = ref<Part | null>(null)

const useForm = reactive({
  order_id: '',
  quantity: 1
})

const filteredParts = computed(() => {
  let result = parts.value

  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    result = result.filter(part =>
      part.name.toLowerCase().includes(keyword) ||
      part.code.toLowerCase().includes(keyword)
    )
  }

  if (categoryFilter.value) {
    result = result.filter(part => part.category === categoryFilter.value)
  }

  return result
})

onMounted(() => {
  loadParts()
})

async function loadParts() {
  try {
    const res = await partApi.getParts({
      page: currentPage.value,
      page_size: pageSize.value
    })
    parts.value = res.data?.list || []
    total.value = res.data?.total || 0
    lowStockCount.value = (res.data as any)?.low_stock_count || 0
  } catch (error) {
    console.error('Failed to load parts:', error)
  }
}

function handleSearch() {
  // Filtering handled by computed
}

function usePart(part: Part) {
  selectedPart.value = part
  useForm.order_id = ''
  useForm.quantity = 1
  showUseDialog.value = true
}

async function submitUsePart() {
  if (!useForm.order_id) {
    ElMessage.error('请输入工单号')
    return
  }

  submitting.value = true
  try {
    await partApi.usePart({
      order_id: Number(useForm.order_id),
      part_id: selectedPart.value!.id,
      quantity: useForm.quantity
    })
    ElMessage.success('使用记录已保存')
    showUseDialog.value = false
    loadParts()
  } catch (error) {
    console.error('Failed to use part:', error)
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.parts-page {
  min-height: 100vh;
  background-color: #f5f7fa;
}

.warning-card {
  margin-bottom: 20px;
  display: flex;
  align-items: center;
  gap: 10px;
}

.filter-bar {
  display: flex;
  gap: 15px;
  margin-bottom: 20px;
}

.parts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 20px;
}

.part-card {
  text-align: center;
}

.part-image {
  height: 100px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 15px;
  color: #909399;
}

.part-image img {
  max-height: 100px;
  max-width: 100%;
}

.part-info {
  text-align: left;
}

.part-name {
  margin: 0 0 8px 0;
  font-size: 16px;
  color: #303133;
}

.part-category {
  color: #909399;
  font-size: 12px;
  margin-bottom: 10px;
}

.part-price {
  color: #f56c6c;
  font-weight: 600;
  margin-bottom: 10px;
}

.part-stock {
  font-size: 14px;
  color: #606266;
  margin-bottom: 15px;
}

.low-stock {
  color: #f56c6c;
  font-weight: 600;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 30px;
}
</style>
