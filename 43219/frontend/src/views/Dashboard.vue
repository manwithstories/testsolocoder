<template>
  <AppLayout>
    <div class="page">
      <div class="card">
        <div class="row">
          <h2 style="margin:0">欢迎, {{ userStore.user?.real_name || userStore.user?.username }}</h2>
          <el-tag size="large" type="primary">{{ roleLabel }}</el-tag>
        </div>
        <p class="muted" style="margin-top:12px">
          根据您的角色,请从顶部菜单选择对应功能入口。
        </p>
      </div>
      <div class="row">
        <router-link to="/services">
          <el-card shadow="hover" style="width:240px">
            <div style="font-size:28px">🧹</div>
            <h3>服务项目</h3>
            <p class="muted">查看所有家政服务</p>
          </el-card>
        </router-link>
        <router-link v-if="userStore.role==='customer'" to="/booking/new">
          <el-card shadow="hover" style="width:240px">
            <div style="font-size:28px">📅</div>
            <h3>预约服务</h3>
            <p class="muted">提交新的预约订单</p>
          </el-card>
        </router-link>
        <router-link to="/orders">
          <el-card shadow="hover" style="width:240px">
            <div style="font-size:28px">📦</div>
            <h3>订单管理</h3>
            <p class="muted">查看订单与确认</p>
          </el-card>
        </router-link>
        <router-link v-if="userStore.role==='staff'" to="/staff/earnings">
          <el-card shadow="hover" style="width:240px">
            <div style="font-size:28px">💰</div>
            <h3>我的收益</h3>
            <p class="muted">查看分成与提现</p>
          </el-card>
        </router-link>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import AppLayout from '../components/AppLayout.vue'
import { useUserStore } from '../stores/user'

const userStore = useUserStore()
const roleLabel = computed(() => ({
  company: '家政公司', staff: '家政人员', customer: '服务客户', admin: '管理员',
} as Record<string, string>)[userStore.role] || '')
</script>
