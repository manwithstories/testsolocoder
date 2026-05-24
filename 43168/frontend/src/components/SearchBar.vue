<template>
  <div class="search-bar">
    <el-form :inline="true" :model="searchForm" @submit.prevent>
      <el-form-item v-for="item in fields" :key="item.prop" :label="item.label">
        <el-input
          v-if="item.type === 'input'"
          v-model="searchForm[item.prop]"
          :placeholder="item.placeholder || `请输入${item.label}`"
          :clearable="item.clearable !== false"
          style="width: 200px"
        />
        <el-select
          v-else-if="item.type === 'select'"
          v-model="searchForm[item.prop]"
          :placeholder="item.placeholder || `请选择${item.label}`"
          :clearable="item.clearable !== false"
          style="width: 200px"
        >
          <el-option
            v-for="opt in item.options"
            :key="opt.value"
            :label="opt.label"
            :value="opt.value"
          />
        </el-select>
        <el-date-picker
          v-else-if="item.type === 'date'"
          v-model="searchForm[item.prop]"
          :placeholder="item.placeholder || `请选择${item.label}`"
          type="date"
          value-format="YYYY-MM-DD"
          style="width: 200px"
        />
        <el-date-picker
          v-else-if="item.type === 'daterange'"
          v-model="searchForm[item.prop]"
          :placeholder="item.placeholder || '选择日期范围'"
          type="daterange"
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          value-format="YYYY-MM-DD"
          style="width: 240px"
        />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" :icon="Search" @click="handleSearch">搜索</el-button>
        <el-button :icon="Refresh" @click="handleReset">重置</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { reactive } from 'vue'
import { Search, Refresh } from '@element-plus/icons-vue'

interface SearchField {
  prop: string
  label: string
  type?: 'input' | 'select' | 'date' | 'daterange'
  placeholder?: string
  options?: { label: string; value: string | number }[]
  clearable?: boolean
}

const props = withDefaults(
  defineProps<{
    fields: SearchField[]
    modelValue: Record<string, unknown>
  }>(),
  {
    modelValue: () => ({})
  }
)

const emit = defineEmits<{
  (e: 'update:modelValue', value: Record<string, unknown>): void
  (e: 'search', value: Record<string, unknown>): void
  (e: 'reset'): void
}>

const searchForm = reactive<Record<string, unknown>>({ ...props.modelValue })

function handleSearch() {
  const result: Record<string, unknown> = {}
  Object.keys(searchForm).forEach((key) => {
    if (searchForm[key] !== '' && searchForm[key] !== null && searchForm[key] !== undefined) {
      result[key] = searchForm[key]
    }
  })
  emit('update:modelValue', result)
  emit('search', result)
}

function handleReset() {
  Object.keys(searchForm).forEach((key) => {
    searchForm[key] = props.modelValue[key] ?? ''
  })
  emit('update:modelValue', {})
  emit('reset')
}
</script>

<style lang="scss" scoped>
.search-bar {
  background-color: #fff;
  padding: 16px 16px 0 16px;
  border-radius: 4px;
  margin-bottom: 16px;
}
</style>