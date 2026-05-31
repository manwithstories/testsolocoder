<template>
  <AppLayout>
    <div class="page">
      <div class="row" style="margin-bottom:16px">
      <h2 style="margin:0">服务项目</h2>
      <el-input v-model="keyword" placeholder="搜索服务名/描述" clearable style="width:240px" @change="load" />
      <el-select v-model="category" placeholder="全部分类" clearable style="width:160px" @change="load">
        <el-option v-for="c in categories" :key="c" :value="c" :label="c" />
      </el-select>
      </div>
      <el-row :gutter="16">
        <el-col v-for="s in services" :key="s.id" :xs="24" :sm="12" :md="8">
          <el-card class="service-card">
          <template #header>
            <div class="row">
              <strong>{{ s.name }}</strong>
              <el-tag size="small">{{ s.category || '通用' }}</el-tag>
            </div>
          </template>
          <p class="muted">{{ s.desc || '暂无描述' }}</p>
          <p>价格区间: <b>¥{{ s.min_price }} - ¥{{ s.max_price }}</p>
          <p>时长: {{ s.duration }} 分钟</p>
          <p v-if="s.skills" class="muted">技能: {{ s.skills }}</p>
          <template #footer>
            <div class="row">
              <el-button
                v-if="userStore.role==='customer'" type="primary" @click="goBook(s)">立即预约</el-button>
              <el-button @click="$router.push(`/services/${s.id}`)">查看详情</el-button>
            </div>
          </template>
        </el-card>
      </el-col>
    </el-row>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import AppLayout from '../components/AppLayout.vue'
import { listServices, type ServiceItem } from '../api/service'
import { useUserStore } from '../stores/user'

const userStore = useUserStore()
const router = useRouter()
const services = ref<ServiceItem[]>([])
const keyword = ref('')
const category = ref('')
const categories = ['保洁', '月嫂', '护工', '维修', '育婴', '老人陪护']

async function load() {
  const params: Record<string, string> = {}
  if (keyword.value) params.keyword = keyword.value
  if (category.value) params.category = category.value
  const res = await listServices(params)
  services.value = (res.data as any).data || []
}

function goBook(s: ServiceItem) {
  router.push({ path: '/booking/new', query: { service_id: String(s.id) })
}

onMounted(load)
</script>

<style scoped>
.service-card { margin-bottom: 16px; }
</style>
