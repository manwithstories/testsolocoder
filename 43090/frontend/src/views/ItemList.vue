<template>
  <div class="item-list">
    <el-card class="filter-card">
      <el-form :inline="true" :model="filters" @submit.prevent>
        <el-form-item label="搜索">
          <el-input v-model="filters.keyword" placeholder="搜索拍卖品" clearable />
        </el-form-item>
        <el-form-item label="分类">
          <el-select v-model="filters.category_id" placeholder="全部分类" clearable>
            <el-option v-for="cat in categories" :key="cat.id" :label="cat.name" :value="cat.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filters.status" placeholder="全部状态" clearable>
            <el-option label="拍卖中" :value="1" />
            <el-option label="已售出" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchItems">搜索</el-button>
          <el-button @click="resetFilters">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-row :gutter="20" class="items-grid">
      <el-col v-for="item in items" :key="item.id" :xs="24" :sm="12" :md="8" :lg="6">
        <ItemCard :item="item" />
      </el-col>
    </el-row>

    <el-empty v-if="items.length === 0" description="暂无拍卖品" />

    <el-pagination
      v-if="total > 0"
      class="pagination"
      v-model:current-page="filters.page"
      v-model:page-size="filters.page_size"
      :total="total"
      :page-sizes="[12, 24, 48]"
      layout="total, sizes, prev, pager, next, jumper"
      @size-change="fetchItems"
      @current-change="fetchItems"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import type { AuctionItem, Category } from '@/types'
import { itemApi } from '@/api'
import ItemCard from '@/components/ItemCard.vue'

const items = ref<AuctionItem[]>([])
const total = ref(0)
const categories = ref<Category[]>([])

const filters = reactive({
  page: 1,
  page_size: 12,
  keyword: '',
  category_id: undefined as number | undefined,
  status: undefined as number | undefined,
  sort_by: 'created_at',
  sort_order: 'desc',
})

const fetchItems = async () => {
  try {
    const res = await itemApi.getList(filters)
    items.value = res.list
    total.value = res.total
  } catch (e) {}
}

const resetFilters = () => {
  filters.page = 1
  filters.keyword = ''
  filters.category_id = undefined
  filters.status = undefined
  fetchItems()
}

onMounted(() => {
  fetchItems()
})
</script>

<style scoped>
.filter-card {
  margin-bottom: 20px;
}

.items-grid {
  margin-bottom: 30px;
}

.pagination {
  display: flex;
  justify-content: center;
}
</style>
