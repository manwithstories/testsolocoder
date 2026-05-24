<template>
  <div class="login-container">
    <div class="login-card" style="text-align: center;">
      <el-icon :size="64" color="#67c23a"><CircleCheckFilled /></el-icon>
      <h2 class="login-title" style="margin-top: 16px;">感谢您的参与！</h2>
      <p style="color: #909399; margin-bottom: 20px;">
        您的问卷已成功提交，感谢您抽出宝贵的时间参与我们的调查。
      </p>
      <el-button type="primary" @click="handleBack">返回首页</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { CircleCheckFilled } from '@element-plus/icons-vue'

const router = useRouter()

const handleBack = () => {
  router.push('/login')
}
</script>
