<template>
  <el-table :data="packages" stripe v-loading="loading">
    <el-table-column label="服务">
      <template #default="{ row }">
        {{ row.service?.name }}
      </template>
    </el-table-column>
    <el-table-column prop="total_count" label="总次数" width="100" />
    <el-table-column prop="used_count" label="已使用" width="100" />
    <el-table-column label="剩余次数" width="100">
      <template #default="{ row }">
        {{ row.total_count - row.used_count }}
      </template>
    </el-table-column>
    <el-table-column prop="purchase_date" label="购买日期">
      <template #default="{ row }">
        {{ formatDate(row.purchase_date) }}
      </template>
    </el-table-column>
    <el-table-column prop="expire_date" label="有效期至">
      <template #default="{ row }">
        {{ formatDate(row.expire_date) }}
      </template>
    </el-table-column>
    <el-table-column label="状态">
      <template #default="{ row }">
        <el-tag :type="isExpired(row.expire_date) ? 'danger' : 'success'">
          {{ isExpired(row.expire_date) ? '已过期' : '有效' }}
        </el-tag>
      </template>
    </el-table-column>
  </el-table>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { getCustomerPackages } from '@/api/payment'
import dayjs from 'dayjs'
import type { CustomerPackage } from '@/types'

const props = defineProps<{
  customerId: number
}>()

const loading = ref(false)
const packages = ref<CustomerPackage[]>([])

const formatDate = (date: string) => dayjs(date).format('YYYY-MM-DD')
const isExpired = (date: string) => dayjs(date).isBefore(dayjs())

const fetchList = async () => {
  loading.value = true
  try {
    const res = await getCustomerPackages(props.customerId)
    packages.value = res.data
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

watch(() => props.customerId, fetchList)

onMounted(fetchList)
</script>
