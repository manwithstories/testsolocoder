<template>
  <div class="not-found-page">
    <div class="not-found-container">
      <div class="error-code">
        <span class="digit">4</span>
        <span class="digit digit-0">0</span>
        <span class="digit">4</span>
      </div>
      <h1 class="error-title">页面未找到</h1>
      <p class="error-message">抱歉，您访问的页面不存在或已被移除</p>
      <div class="action-buttons">
        <el-button type="primary" size="large" @click="goHome">
          <el-icon><House /></el-icon>
          返回首页
        </el-button>
        <el-button size="large" @click="goBack">
          <el-icon><ArrowLeft /></el-icon>
          返回上一页
        </el-button>
      </div>
      <div class="suggestions">
        <h3>您可能想要：</h3>
        <div class="suggestion-links">
          <div class="suggestion-item" @click="$router.push('/exhibitions')">
            <el-icon size="24"><Picture /></el-icon>
            <span>浏览展览</span>
          </div>
          <div class="suggestion-item" @click="$router.push('/collections')">
            <el-icon size="24"><Collection /></el-icon>
            <span>探索藏品</span>
          </div>
          <div class="suggestion-item" @click="$router.push('/login')">
            <el-icon size="24"><User /></el-icon>
            <span>登录账户</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

const goHome = () => {
  router.push('/')
}

const goBack = () => {
  if (window.history.length > 1) {
    router.back()
  } else {
    router.push('/')
  }
}
</script>

<style scoped lang="scss">
.not-found-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.not-found-container {
  text-align: center;
  color: #fff;
  max-width: 600px;
  width: 100%;
}

.error-code {
  display: flex;
  justify-content: center;
  gap: 8px;
  margin-bottom: 32px;

  .digit {
    font-size: 120px;
    font-weight: 900;
    line-height: 1;
    color: #fff;
    text-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
    animation: bounce 2s ease-in-out infinite;

    &:nth-child(1) {
      animation-delay: 0s;
    }

    &:nth-child(2) {
      animation-delay: 0.2s;
    }

    &:nth-child(3) {
      animation-delay: 0.4s;
    }

    &.digit-0 {
      position: relative;

      &::before {
        content: '';
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        width: 60px;
        height: 60px;
        background: rgba(255, 255, 255, 0.2);
        border-radius: 50%;
        animation: pulse 2s ease-in-out infinite;
      }
    }
  }
}

.error-title {
  font-size: 36px;
  font-weight: 700;
  margin-bottom: 16px;
}

.error-message {
  font-size: 18px;
  opacity: 0.9;
  margin-bottom: 40px;
}

.action-buttons {
  display: flex;
  gap: 16px;
  justify-content: center;
  margin-bottom: 60px;
  flex-wrap: wrap;

  .el-button {
    min-width: 160px;
    height: 48px;
    font-size: 16px;
    border-radius: 24px;

    &.el-button--primary {
      background: #fff;
      color: #667eea;
      border-color: #fff;

      &:hover {
        background: rgba(255, 255, 255, 0.9);
      }
    }

    &:not(.el-button--primary) {
      background: rgba(255, 255, 255, 0.2);
      border-color: rgba(255, 255, 255, 0.3);
      color: #fff;

      &:hover {
        background: rgba(255, 255, 255, 0.3);
      }
    }
  }
}

.suggestions {
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  padding: 32px;
  border: 1px solid rgba(255, 255, 255, 0.2);

  h3 {
    font-size: 18px;
    margin-bottom: 24px;
    opacity: 0.9;
  }

  .suggestion-links {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 16px;

    .suggestion-item {
      background: rgba(255, 255, 255, 0.15);
      border-radius: 12px;
      padding: 20px 16px;
      cursor: pointer;
      transition: all 0.3s ease;
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 12px;

      &:hover {
        background: rgba(255, 255, 255, 0.25);
        transform: translateY(-4px);
      }

      .el-icon {
        opacity: 0.9;
      }

      span {
        font-size: 14px;
        opacity: 0.9;
      }
    }
  }
}

@keyframes bounce {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-16px);
  }
}

@keyframes pulse {
  0%, 100% {
    transform: translate(-50%, -50%) scale(1);
    opacity: 0.5;
  }
  50% {
    transform: translate(-50%, -50%) scale(1.5);
    opacity: 0;
  }
}

@media (max-width: 768px) {
  .error-code {
    .digit {
      font-size: 80px;
    }

    .digit-0::before {
      width: 40px;
      height: 40px;
    }
  }

  .error-title {
    font-size: 28px;
  }

  .error-message {
    font-size: 16px;
  }

  .action-buttons {
    .el-button {
      min-width: 140px;
      font-size: 14px;
    }
  }

  .suggestions {
    padding: 24px 16px;

    .suggestion-links {
      grid-template-columns: 1fr;
      gap: 12px;

      .suggestion-item {
        flex-direction: row;
        justify-content: flex-start;
        padding: 16px;
      }
    }
  }
}
</style>
