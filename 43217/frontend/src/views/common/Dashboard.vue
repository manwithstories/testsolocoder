<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElCard, ElRow, ElCol, ElStatistic, ElIcon } from 'element-plus'
import { User, Calendar, Document, DataAnalysis, Bell } from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()

const stats = computed(() => {
  const role = userStore.userRole
  if (role === 'hr') {
    return [
      { title: '员工总数', value: 0, icon: User, color: '#409eff', path: '/hr/employees' },
      { title: '本月预约', value: 0, icon: Calendar, color: '#67c23a', path: '/hr/appointments' },
      { title: '体检报告', value: 0, icon: Document, color: '#e6a23c', path: '/hr/reports' },
      { title: '数据分析', value: '查看', icon: DataAnalysis, color: '#f56c6c', path: '/hr/statistics' }
    ]
  } else if (role === 'agency') {
    return [
      { title: '套餐数量', value: 0, icon: Document, color: '#409eff', path: '/agency/packages' },
      { title: '今日预约', value: 0, icon: Calendar, color: '#67c23a', path: '/agency/appointments' },
      { title: '待上传报告', value: 0, icon: Bell, color: '#e6a23c', path: '/agency/reports' },
      { title: '账单管理', value: '查看', icon: DataAnalysis, color: '#f56c6c', path: '/agency/billings' }
    ]
  } else {
    return [
      { title: '我的预约', value: 0, icon: Calendar, color: '#409eff', path: '/employee/appointments' },
      { title: '体检报告', value: 0, icon: Document, color: '#67c23a', path: '/employee/reports' },
      { title: '健康档案', value: '查看', icon: User, color: '#e6a23c', path: '/employee/health-records' },
      { title: '异常提醒', value: 0, icon: Bell, color: '#f56c6c', path: '/employee/reminders' }
    ]
  }
})

const navigateTo = (path: string) => {
  router.push(path)
}
</script>

<template>
  <div class="dashboard">
    <ElCard class="welcome-card">
      <div class="welcome-content">
        <h2>欢迎回来，{{ userStore.userInfo?.real_name || userStore.userInfo?.username }}！</h2>
        <p>今天是 {{ new Date().toLocaleDateString('zh-CN') }}，祝您工作愉快！</p>
      </div>
    </ElCard>

    <ElRow :gutter="20" class="stats-row">
      <ElCol :xs="12" :sm="12" :md="6" :lg="6" v-for="(stat, index) in stats" :key="index">
        <ElCard class="stat-card" @click="navigateTo(stat.path)" shadow="hover">
          <div class="stat-content">
            <ElIcon :size="40" :color="stat.color">
              <component :is="stat.icon" />
            </ElIcon>
            <ElStatistic :title="stat.title" :value="stat.value" />
          </div>
        </ElCard>
      </ElCol>
    </ElRow>

    <ElRow :gutter="20">
      <ElCol :span="12">
        <ElCard class="info-card">
          <template #header>
            <span>快捷功能</span>
          </template>
          <div class="quick-actions">
            <button 
              v-if="userStore.userRole === 'hr'" 
              class="action-btn"
              @click="router.push('/hr/employees')"
            >
              员工管理
            </button>
            <button 
              v-if="userStore.userRole === 'hr'" 
              class="action-btn"
              @click="router.push('/hr/appointments')"
            >
              预约管理
            </button>
            <button 
              v-if="userStore.userRole === 'agency'" 
              class="action-btn"
              @click="router.push('/agency/packages')"
            >
              套餐管理
            </button>
            <button 
              v-if="userStore.userRole === 'agency'" 
              class="action-btn"
              @click="router.push('/agency/reports')"
            >
              报告上传
            </button>
            <button 
              v-if="userStore.userRole === 'employee'" 
              class="action-btn"
              @click="router.push('/employee/new-appointment')"
            >
              新建预约
            </button>
            <button 
              v-if="userStore.userRole === 'employee'" 
              class="action-btn"
              @click="router.push('/employee/reports')"
            >
              查看报告
            </button>
          </div>
        </ElCard>
      </ElCol>
      <ElCol :span="12">
        <ElCard class="info-card">
          <template #header>
            <span>系统通知</span>
          </template>
          <div class="notifications">
            <p class="no-notification">暂无新通知</p>
          </div>
        </ElCard>
      </ElCol>
    </ElRow>
  </div>
</template>

<style scoped lang="scss">
.dashboard {
  .welcome-card {
    margin-bottom: 20px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border: none;
    color: #fff;

    :deep(.el-card__body) {
      padding: 30px;
    }

    .welcome-content {
      h2 {
        margin: 0 0 10px 0;
        font-size: 24px;
        color: #fff;
      }

      p {
        margin: 0;
        color: rgba(255, 255, 255, 0.8);
        font-size: 14px;
      }
    }
  }

  .stats-row {
    margin-bottom: 20px;

    .stat-card {
      cursor: pointer;
      transition: all 0.3s;

      &:hover {
        transform: translateY(-5px);
      }

      .stat-content {
        display: flex;
        align-items: center;
        gap: 15px;

        :deep(.el-statistic__head) {
          font-size: 14px;
          color: #909399;
        }

        :deep(.el-statistic__content) {
          font-size: 28px;
        }
      }
    }
  }

  .info-card {
    .quick-actions {
      display: flex;
      flex-wrap: wrap;
      gap: 10px;

      .action-btn {
        padding: 10px 20px;
        border: 1px solid #dcdfe6;
        border-radius: 4px;
        background-color: #fff;
        cursor: pointer;
        transition: all 0.3s;

        &:hover {
          background-color: #409eff;
          color: #fff;
          border-color: #409eff;
        }
      }
    }

    .notifications {
      .no-notification {
        color: #909399;
        text-align: center;
        padding: 20px 0;
      }
    }
  }
}
</style>
