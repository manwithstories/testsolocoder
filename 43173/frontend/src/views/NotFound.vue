<template>
  <div class="not-found">
    <div class="content">
      <div class="error-code">404</div>
      <div class="title">页面未找到</div>
      <div class="description">抱歉，您访问的页面不存在或已被移除</div>
      <div class="actions">
        <el-button type="primary" @click="goHome">返回首页</el-button>
        <el-button @click="goBack">返回上一页</el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

function goHome() {
  router.push('/')
}

function goBack() {
  router.back()
}
</script>

<style scoped lang="scss">
.not-found {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  
  .content {
    text-align: center;
    color: #fff;
    
    .error-code {
      font-size: 120px;
      font-weight: bold;
      line-height: 1;
      margin-bottom: 20px;
      text-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
    }
    
    .title {
      font-size: 24px;
      font-weight: 500;
      margin-bottom: 16px;
    }
    
    .description {
      font-size: 16px;
      opacity: 0.9;
      margin-bottom: 32px;
    }
    
    .actions {
      display: flex;
      gap: 12px;
      justify-content: center;
    }
  }
}
</style>
