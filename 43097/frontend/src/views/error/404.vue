<template>
  <div class="not-found-container">
    <div class="not-found-content">
      <h1 class="error-code">404</h1>
      <h2 class="error-title">页面不存在</h2>
      <p class="error-desc">抱歉，您访问的页面不存在或已被移除</p>
      <el-button type="primary" @click="goHome">
        <el-icon><House /></el-icon>
        返回首页
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { House } from '@element-plus/icons-vue'

const router = useRouter()

const goHome = () => {
  router.push('/')
}
</script>

<style scoped lang="scss">
.not-found-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.not-found-content {
  text-align: center;
  color: #fff;
  padding: 40px;

  .error-code {
    font-size: 120px;
    font-weight: bold;
    margin: 0;
    text-shadow: 4px 4px 8px rgba(0, 0, 0, 0.3);
  }

  .error-title {
    font-size: 32px;
    margin: 20px 0;
  }

  .error-desc {
    font-size: 16px;
    margin-bottom: 30px;
    opacity: 0.9;
  }
}
</style>
