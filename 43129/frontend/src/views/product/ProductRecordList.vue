<template>
  <div class="page-container">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>库存记录</span>
          <el-button @click="$router.back()">返回</el-button>
        </div>
      </template>

      <el-table :data="records" stripe v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="产品">
          <template #default="{ row }">
            {{ row.product?.name }}
          </template>
        </el-table-column>
        <el-table-column prop="change_type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="getTypeTag(row.change_type)">{{ getTypeText(row.change_type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="quantity" label="数量" width="100">
          <template #default="{ row }">
            <span :class="{ 'text-green': row.quantity > 0, 'text-red': row.quantity < 0 }">
              {{ row.quantity > 0 ? '+' : '' }}{{ row.quantity }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="before_stock" label="变动前" width="100" />
        <el-table-column prop="after_stock" label="变动后" width="100" />
        <el-table-column prop="remark" label="备注" show-overflow-tooltip />
        <el-table-column prop="created_at" label="时间" width="180">
          <template #default="{ row }">{{ formatDateTime(row.created_at) }}</template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="fetchList"
          @size-change="fetchList"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getProductRecords } from '@/api/product'
import dayjs from 'dayjs'
import type { ProductRecord } from '@/types'

const loading = ref(false)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const records = ref<ProductRecord[]>([])

const getTypeTag = (type: string) => {
  const tags: Record<string, string> = {
    in: 'success',
    out: 'danger',
    sale: 'warning',
    adjust: 'info'
  }
  return tags[type] || 'info'
}

const getTypeText = (type: string) => {
  const texts: Record<string, string> = {
    in: '入库',
    out: '出库',
    sale: '销售',
    adjust: '调整'
  }
  return texts[type] || type
}

const formatDateTime = (date: string) => dayjs(date).format('YYYY-MM-DD HH:mm:ss')

const fetchList = async () => {
  loading.value = true
  try {
    const res = await getProductRecords({ page: page.value, page_size: pageSize.value })
    records.value = res.data.list
    total.value = res.data.total
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(fetchList)
</script>

<style scoped lang="scss">
.page-container {
  .card-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .pagination {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }

  .text-green {
    color: #67c23a;
  }

  .text-red {
    color: #f56c6c;
  }
}
</style>
