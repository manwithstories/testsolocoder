<template>
  <div class="product-list">
    <SearchBar
      :fields="searchFields"
      v-model="searchForm"
      @search="handleSearch"
      @reset="handleReset"
    />

    <el-card class="list-card" shadow="never">
      <div class="card-header">
        <el-button
          v-if="isManufacturer"
          type="primary"
          :icon="Plus"
          @click="goCreate"
        >
          新增产品
        </el-button>
        <el-button :icon="Refresh" @click="fetchList">刷新</el-button>
      </div>

      <el-table
        v-loading="loading"
        :data="tableData"
        stripe
        border
        style="width: 100%"
      >
        <el-table-column label="图片" width="90" align="center">
          <template #default="{ row }">
            <el-image
              :src="row.imageUrl"
              :preview-src-list="row.images?.map(i => i.url) || [row.imageUrl].filter(Boolean)"
              fit="cover"
              style="width: 56px; height: 56px; border-radius: 4px"
            >
              <template #error>
                <div class="image-placeholder">
                  <el-icon :size="20"><Picture /></el-icon>
                </div>
              </template>
            </el-image>
          </template>
        </el-table-column>
        <el-table-column prop="name" label="名称" min-width="160" show-overflow-tooltip />
        <el-table-column prop="category" label="分类" width="120" />
        <el-table-column label="价格" width="120" align="right">
          <template #default="{ row }">¥{{ formatPrice(row.price) }}</template>
        </el-table-column>
        <el-table-column prop="stock" label="库存" width="100" align="center" />
        <el-table-column label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">
              {{ row.status === 1 ? '上架' : '下架' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" align="center" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="goDetail(row)">查看</el-button>
            <el-button
              v-if="isManufacturer"
              link
              type="primary"
              @click="goEdit(row)"
            >
              编辑
            </el-button>
            <el-button
              v-if="isManufacturer"
              link
              type="danger"
              @click="handleDelete(row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <Pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        @change="handlePageChange"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Refresh, Picture } from '@element-plus/icons-vue'
import SearchBar from '@/components/SearchBar.vue'
import Pagination from '@/components/Pagination.vue'
import { listProducts, deleteProduct } from '@/api/product'
import { useUserStore } from '@/stores/user'
import type { Product } from '@/types'

const router = useRouter()
const userStore = useUserStore()
const isManufacturer = computed(() => userStore.userRole === 'manufacturer')

const loading = ref(false)
const tableData = ref<Product[]>([])
const searchForm = reactive<Record<string, unknown>>({})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const searchFields = [
  { prop: 'keyword', label: '关键字', type: 'input' as const, placeholder: '名称/SKU' },
  {
    prop: 'category',
    label: '分类',
    type: 'select' as const,
    options: [
      { label: '沙发', value: '沙发' },
      { label: '椅子', value: '椅子' },
      { label: '桌子', value: '桌子' },
      { label: '床', value: '床' },
      { label: '柜', value: '柜' }
    ]
  },
  {
    prop: 'status',
    label: '状态',
    type: 'select' as const,
    options: [
      { label: '上架', value: 1 },
      { label: '下架', value: 0 }
    ]
  }
]

async function fetchList() {
  loading.value = true
  try {
    const data = await listProducts({
      page: pagination.page,
      pageSize: pagination.pageSize,
      ...(searchForm as Record<string, unknown>)
    })
    tableData.value = data.list
    pagination.total = data.total
  } catch (err) {
    console.error(err)
  } finally {
    loading.value = false
  }
}

function handleSearch(form: Record<string, unknown>) {
  Object.assign(searchForm, form)
  pagination.page = 1
  fetchList()
}

function handleReset() {
  Object.keys(searchForm).forEach((k) => delete searchForm[k])
  pagination.page = 1
  fetchList()
}

function handlePageChange(page: number, pageSize: number) {
  pagination.page = page
  pagination.pageSize = pageSize
  fetchList()
}

function goDetail(row: Product) {
  router.push(`/products/${row.id}`)
}

function goEdit(row: Product) {
  router.push(`/products/${row.id}/edit`)
}

function goCreate() {
  router.push('/products/create')
}

async function handleDelete(row: Product) {
  try {
    await ElMessageBox.confirm(`确定删除产品「${row.name}」吗？`, '提示', {
      type: 'warning'
    })
    await deleteProduct(row.id)
    ElMessage.success('删除成功')
    fetchList()
  } catch {
    // 取消
  }
}

function formatPrice(price: number) {
  return (price || 0).toFixed(2)
}

onMounted(fetchList)
</script>

<style lang="scss" scoped>
.product-list {
  .card-header {
    display: flex;
    gap: 8px;
    margin-bottom: 16px;
  }

  .image-placeholder {
    width: 56px;
    height: 56px;
    background-color: #f5f7fa;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #c0c4cc;
    border-radius: 4px;
  }
}
</style>
