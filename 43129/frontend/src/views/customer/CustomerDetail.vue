<template>
  <div class="page-container" v-loading="loading">
    <el-card shadow="never" v-if="customer">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <el-button :icon="ArrowLeft" @click="$router.back()">返回</el-button>
            <span style="margin-left: 10px; font-weight: 600">顾客详情</span>
          </div>
        </div>
      </template>

      <el-descriptions :column="2" border>
        <el-descriptions-item label="姓名">{{ customer.name || '-' }}</el-descriptions-item>
        <el-descriptions-item label="手机号">{{ customer.user?.phone }}</el-descriptions-item>
        <el-descriptions-item label="性别">{{ customer.gender || '-' }}</el-descriptions-item>
        <el-descriptions-item label="年龄">{{ customer.age || '-' }}</el-descriptions-item>
        <el-descriptions-item label="皮肤类型">{{ customer.skin_type || '-' }}</el-descriptions-item>
        <el-descriptions-item label="发型偏好">{{ customer.hair_preference || '-' }}</el-descriptions-item>
        <el-descriptions-item label="过敏史">{{ customer.allergy_history || '-' }}</el-descriptions-item>
        <el-descriptions-item label="备注">{{ customer.notes || '-' }}</el-descriptions-item>
        <el-descriptions-item label="会员等级">
          <el-tag :type="getLevelType(customer.member_level)">
            {{ getLevelText(customer.member_level) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="积分">{{ customer.points }}</el-descriptions-item>
        <el-descriptions-item label="累计消费">¥{{ customer.total_spent?.toFixed(2) }}</el-descriptions-item>
        <el-descriptions-item label="到店次数">{{ customer.visit_count }}</el-descriptions-item>
      </el-descriptions>

      <el-divider />

      <el-tabs v-model="activeTab">
        <el-tab-pane label="预约记录" name="appointments">
          <appointment-list :customer-id="customer.id" />
        </el-tab-pane>

        <el-tab-pane label="会员卡" name="cards">
          <member-card-list :customer-id="customer.id" />
        </el-tab-pane>

        <el-tab-pane label="套餐" name="packages">
          <package-list :customer-id="customer.id" />
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getCustomer } from '@/api/auth'
import { ArrowLeft } from '@element-plus/icons-vue'
import type { Customer } from '@/types'
import AppointmentList from './components/CustomerAppointments.vue'
import MemberCardList from './components/CustomerCards.vue'
import PackageList from './components/CustomerPackages.vue'

const route = useRoute()
const loading = ref(false)
const customer = ref<Customer | null>(null)
const activeTab = ref('appointments')

const getLevelType = (level: number) => {
  const types = ['', 'info', '', 'warning', 'success', 'danger']
  return types[level] || 'info'
}

const getLevelText = (level: number) => {
  const texts = ['', '普通会员', '', '银卡会员', '金卡会员', '钻石会员']
  return texts[level] || '普通会员'
}

const fetchCustomer = async () => {
  loading.value = true
  try {
    const id = Number(route.params.id)
    const res = await getCustomer(id)
    customer.value = res.data
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(fetchCustomer)
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
