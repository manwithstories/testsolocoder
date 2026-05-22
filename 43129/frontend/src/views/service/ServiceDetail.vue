<template>
  <div class="page-container" v-loading="loading">
    <el-card shadow="never" v-if="service">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <el-button :icon="ArrowLeft" @click="$router.back()">返回</el-button>
            <span style="margin-left: 10px; font-weight: 600">服务详情</span>
          </div>
        </div>
      </template>

      <el-descriptions :column="2" border>
        <el-descriptions-item label="服务名称">{{ service.name }}</el-descriptions-item>
        <el-descriptions-item label="分类">{{ service.category }}</el-descriptions-item>
        <el-descriptions-item label="描述" :span="2">{{ service.description || '-' }}</el-descriptions-item>
        <el-descriptions-item label="价格">¥{{ service.price?.toFixed(2) }}</el-descriptions-item>
        <el-descriptions-item label="时长">{{ service.duration }}分钟</el-descriptions-item>
        <el-descriptions-item label="所需技能">{{ service.required_skill || '-' }}</el-descriptions-item>
        <el-descriptions-item label="是否套餐">
          <el-tag v-if="service.is_package" type="warning">是 ({{ service.package_count }}次)</el-tag>
          <span v-else>否</span>
        </el-descriptions-item>
        <el-descriptions-item label="动态定价">
          <el-tag v-if="service.dynamic_pricing" type="success">已开启</el-tag>
          <span v-else>未开启</span>
        </el-descriptions-item>
        <el-descriptions-item v-if="service.dynamic_pricing" label="周末价格">
          ¥{{ service.weekend_price?.toFixed(2) }}
        </el-descriptions-item>
      </el-descriptions>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getService } from '@/api/service'
import { ArrowLeft } from '@element-plus/icons-vue'
import type { Service } from '@/types'

const route = useRoute()
const loading = ref(false)
const service = ref<Service | null>(null)

const fetchService = async () => {
  loading.value = true
  try {
    const id = Number(route.params.id)
    const res = await getService(id)
    service.value = res.data
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(fetchService)
</script>

<style scoped lang="scss">
.page-container {
  .card-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .header-left {
    display: flex;
    align-items: center;
  }
}
</style>
