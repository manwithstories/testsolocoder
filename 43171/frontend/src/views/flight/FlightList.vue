<template>
  <el-card>
    <template #header>
      <div class="header">
        <span>飞行记录</span>
        <el-button type="primary" @click="$router.push('/flight/create')" v-if="role === 'pilot'">
          <el-icon><Plus /></el-icon>
          添加记录
        </el-button>
      </div>
    </template>

    <el-table :data="flights" v-loading="loading">
      <el-table-column label="设备" min-width="120">
        <template #default="{ row }">{{ row.drone?.name || '-' }}</template>
      </el-table-column>
      <el-table-column label="飞手" width="120">
        <template #default="{ row }">{{ row.pilot?.nickname || row.pilot?.username }}</template>
      </el-table-column>
      <el-table-column prop="start_point" label="起飞点" />
      <el-table-column prop="end_point" label="降落点" />
      <el-table-column prop="duration" label="时长(分钟)" width="100" />
      <el-table-column prop="distance" label="距离(km)" width="100" />
      <el-table-column prop="flight_date" label="飞行日期" width="120" />
      <el-table-column label="操作" width="100" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link @click="viewDetail(row)">详情</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      style="margin-top: 20px; justify-content: flex-end"
      layout="prev, pager, next, total"
      :total="total"
      :page-size="pageSize"
      v-model:current-page="currentPage"
      @current-change="fetchFlights"
    />
  </el-card>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'

const userStore = useUserStore()
const role = computed(() => userStore.role)

const loading = ref(false)
const flights = ref<FlightRecord[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

onMounted(() => {
  fetchFlights()
})

async function fetchFlights() {
  loading.value = true
  try {
    const res: any = await request.get('/flights', {
      params: { page: currentPage.value, page_size: pageSize.value }
    })
    flights.value = res.data.list || []
    total.value = res.data.total || 0
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function viewDetail(row: FlightRecord) {
  ElMessage.info(`飞行记录详情: ${row.drone?.name}`)
}
</script>

<style scoped>
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
