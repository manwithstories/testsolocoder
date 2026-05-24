<template>
  <div class="not-found-container">
    <div class="not-found-content">
      <h1>404</h1>
      <p>抱歉，您访问的页面不存在</p>
      <el-button type="primary" @click="$router.push('/')">
        返回首页
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
</script>

<style lang="scss" scoped>
.not-found-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #f5f7fa 0%, #e4e7ed 100%);
}

.not-found-content {
  text-align: center;
  
  h1 {
    font-size: 120px;
    font-weight: 700;
    color: var(--primary-color);
    margin-bottom: 16px;
  }
  
  p {
    font-size: 18px;
    color: var(--text-secondary);
    margin-bottom: 32px;
  }
}
</style>
