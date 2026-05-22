<template>
  <el-table :data="cards" stripe v-loading="loading">
    <el-table-column prop="card_no" label="卡号" />
    <el-table-column prop="card_type" label="卡类型" />
    <el-table-column prop="balance" label="余额">
      <template #default="{ row }">
        ¥{{ row.balance?.toFixed(2) }}
      </template>
    </el-table-column>
    <el-table-column prop="discount" label="折扣">
      <template #default="{ row }">
        {{ (row.discount * 10).toFixed(1) }}折
      </template>
    </el-table-column>
    <el-table-column label="状态">
      <template #default="{ row }">
        <el-tag :type="row.status === 1 ? 'success' : 'danger'">
          {{ row.status === 1 ? '正常' : '已停用' }}
        </el-tag>
      </template>
    </el-table-column>
  </el-table>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { getMemberCards } from '@/api/payment'
import type { MemberCard } from '@/types'

const props = defineProps<{
  customerId: number
}>()

const loading = ref(false)
const cards = ref<MemberCard[]>([])

const fetchList = async () => {
  loading.value = true
  try {
    const res = await getMemberCards(props.customerId)
    cards.value = res.data
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

watch(() => props.customerId, fetchList)

onMounted(fetchList)
</script>
