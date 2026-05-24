<template>
  <div class="login-container">
    <div class="login-card" style="text-align: center;">
      <el-icon :size="80" color="#409eff"><WarningFilled /></el-icon>
      <h2 class="login-title" style="margin-top: 16px;">页面未找到</h2>
      <p style="color: #909399; margin-bottom: 20px;">
        您访问的页面不存在或已被删除
      </p>
      <el-button type="primary" @click="handleBack">返回首页</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { WarningFilled } from '@element-plus/icons-vue'

const router = useRouter()

const handleBack = () => {
  router.push('/')
}
</script>
