<template>
  <div class="pagination-wrapper">
    <el-pagination
      v-model:current-page="currentPage"
      v-model:page-size="pageSize"
      :page-sizes="pageSizes"
      :total="total"
      :layout="layout"
      :background="background"
      :small="small"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
    />
  </div>
</template>

<script setup lang="ts">
import { watch } from 'vue'

const props = withDefaults(
  defineProps<{
    currentPage: number
    pageSize: number
    total: number
    pageSizes?: number[]
    layout?: string
    background?: boolean
    small?: boolean
  }>(),
  {
    pageSizes: () => [10, 20, 50, 100],
    layout: 'total, sizes, prev, pager, next, jumper',
    background: true,
    small: false
  }
)

const emit = defineEmits<{
  (e: 'update:currentPage', value: number): void
  (e: 'update:pageSize', value: number): void
  (e: 'change', page: number, pageSize: number): void
}>

function handleSizeChange(size: number) {
  emit('update:pageSize', size)
  emit('update:currentPage', 1)
  emit('change', 1, size)
}

function handleCurrentChange(page: number) {
  emit('update:currentPage', page)
  emit('change', page, props.pageSize)
}

watch(
  () => props.total,
  (newTotal) => {
    const maxPage = Math.max(1, Math.ceil(newTotal / props.pageSize))
    if (props.currentPage > maxPage) {
      emit('update:currentPage', maxPage)
    }
  }
)
</script>

<style lang="scss" scoped>
.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  padding: 16px 0;
}
</style>