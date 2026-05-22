<template>
  <div class="contract-detail">
    <div class="page-header">
      <h2 class="page-title">合同详情</h2>
      <el-button @click="goBack">返回</el-button>
    </div>

    <el-row :gutter="20">
      <el-col :span="16">
        <div class="card">
          <div class="detail-header">
            <h1>合同 #{{ contract?.id }}</h1>
            <el-tag :type="getStatusType(contract?.status)" size="large">
              {{ getStatusText(contract?.status) }}
            </el-tag>
          </div>

          <el-descriptions :column="2" border>
            <el-descriptions-item label="租户">
              {{ contract?.tenant?.name }} ({{ contract?.tenant?.phone }})
            </el-descriptions-item>
            <el-descriptions-item label="房源">
              {{ contract?.property?.title }}
            </el-descriptions-item>
            <el-descriptions-item label="开始日期">
              {{ formatDate(contract?.startDate) }}
            </el-descriptions-item>
            <el-descriptions-item label="结束日期">
              {{ formatDate(contract?.endDate) }}
            </el-descriptions-item>
            <el-descriptions-item label="月租金">
              ¥{{ contract?.rent }}
            </el-descriptions-item>
            <el-descriptions-item label="押金">
              ¥{{ contract?.deposit }}
            </el-descriptions-item>
            <el-descriptions-item label="付款方式">
              {{ contract?.paymentType }}
            </el-descriptions-item>
            <el-descriptions-item label="创建时间">
              {{ formatDateTime(contract?.createdAt) }}
            </el-descriptions-item>
          </el-descriptions>

          <div class="file-section" v-if="contract?.fileUrl">
            <h3>合同文件</h3>
            <el-link :href="contract.fileUrl" type="primary" target="_blank">
              <el-icon><Document /></el-icon>
              查看合同文件
            </el-link>
          </div>
        </div>
      </el-col>

      <el-col :span="8">
        <div class="card">
          <h3>租金记录</h3>
          <el-table :data="rentRecords" style="width: 100%">
            <el-table-column prop="month" label="月份" width="100" />
            <el-table-column prop="amount" label="金额" width="100">
              <template #default="{ row }">¥{{ row.amount }}</template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="80">
              <template #default="{ row }">
                <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
                  {{ row.status === 1 ? '已缴' : '未缴' }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import type { Contract, RentRecord } from '@/types'
import { getContract } from '@/api/tenant'
import { getRentBills } from '@/api/business'
import dayjs from 'dayjs'

const route = useRoute()
const router = useRouter()

const contract = ref<Contract | null>(null)
const rentRecords = ref<RentRecord[]>([])

onMounted(async () => {
  await loadContract()
  await loadRentRecords()
})

async function loadContract() {
  try {
    const res = await getContract(Number(route.params.id))
    contract.value = res.data
  } catch (error) {
    console.error('Failed to load contract:', error)
  }
}

async function loadRentRecords() {
  try {
    const res = await getRentBills({
      pageSize: 100,
      contractId: Number(route.params.id)
    })
    rentRecords.value = res.data.list
  } catch (error) {
    console.error('Failed to load rent records:', error)
  }
}

function formatDate(date?: string) {
  return date ? dayjs(date).format('YYYY-MM-DD') : '-'
}

function formatDateTime(date?: string) {
  return date ? dayjs(date).format('YYYY-MM-DD HH:mm') : '-'
}

function getStatusType(status?: number) {
  switch (status) {
    case 1: return 'success'
    case 2: return 'warning'
    default: return 'info'
  }
}

function getStatusText(status?: number) {
  switch (status) {
    case 1: return '生效中'
    case 2: return '已到期'
    default: return '已终止'
  }
}

function goBack() {
  router.back()
}
</script>

<style scoped>
.contract-detail {
  padding: 0;
}

.detail-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.detail-header h1 {
  font-size: 20px;
  margin: 0;
}

.file-section {
  margin-top: 20px;
}

.file-section h3 {
  font-size: 16px;
  margin-bottom: 10px;
}
</style>
