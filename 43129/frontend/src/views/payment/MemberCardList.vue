<template>
  <div class="page-container">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>会员卡管理</span>
        </div>
      </template>

      <el-table :data="cards" stripe v-loading="loading">
        <el-table-column prop="card_no" label="卡号" />
        <el-table-column label="顾客">
          <template #default="{ row }">
            {{ row.customer?.name || row.customer?.user?.phone }}
          </template>
        </el-table-column>
        <el-table-column prop="card_type" label="卡类型" />
        <el-table-column prop="balance" label="余额">
          <template #default="{ row }">¥{{ row.balance?.toFixed(2) }}</template>
        </el-table-column>
        <el-table-column prop="discount" label="折扣">
          <template #default="{ row }">{{ (row.discount * 10).toFixed(1) }}折</template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '正常' : '已停用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">{{ formatDateTime(row.created_at) }}</template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getMemberCards } from '@/api/payment'
import dayjs from 'dayjs'
import type { MemberCard } from '@/types'

const loading = ref(false)
const cards = ref<MemberCard[]>([])

const formatDateTime = (date: string) => dayjs(date).format('YYYY-MM-DD HH:mm:ss')

const fetchList = async () => {
  loading.value = true
  try {
    const res: any = await getMemberCards(0)
    cards.value = res.data || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(fetchList)
</script>
