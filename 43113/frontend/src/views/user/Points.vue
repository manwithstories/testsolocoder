<template>
  <div class="points-page">
    <el-row :gutter="20">
      <el-col :span="8">
        <el-card>
          <div class="points-summary">
            <div class="current-points">
              <span class="label">当前积分</span>
              <span class="value">{{ userStore.userInfo?.points || 0 }}</span>
            </div>
            <div class="level-info">
              <span>Lv.{{ userStore.userInfo?.level || 1 }}</span>
              <el-progress
                :percentage="levelProgress"
                :stroke-width="8"
                style="margin-top: 8px"
              />
              <span class="next-level">距离下一等级还需 {{ nextLevelPoints }} 积分</span>
            </div>
          </div>
        </el-card>

        <el-card style="margin-top: 20px">
          <template #header>积分商城</template>
          <div class="reward-list">
            <div
              v-for="reward in rewards"
              :key="reward.id"
              class="reward-item"
            >
              <div class="reward-icon">{{ reward.image }}</div>
              <div class="reward-info">
                <div class="reward-name">{{ reward.name }}</div>
                <div class="reward-cost">{{ reward.pointsCost }} 积分</div>
              </div>
              <el-button
                type="primary"
                size="small"
                :disabled="(userStore.userInfo?.points || 0) < reward.pointsCost"
                @click="exchangeReward(reward.id)"
              >
                兑换
              </el-button>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :span="16">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>积分流水</span>
              <el-tabs v-model="activeTab" size="small">
                <el-tab-pane label="全部" name="all" />
                <el-tab-pane label="获得" name="earn" />
                <el-tab-pane label="消耗" name="spend" />
              </el-tabs>
            </div>
          </template>

          <el-table :data="filteredLogs" style="width: 100%">
            <el-table-column prop="type" label="类型" width="120">
              <template #default="{ row }">
                <el-tag :type="getPointTypeTag(row.type)" size="small">
                  {{ getPointTypeName(row.type) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="description" label="说明" />
            <el-table-column prop="points" label="积分变动" width="100">
              <template #default="{ row }">
                <span :class="{ positive: row.points > 0, negative: row.points < 0 }">
                  {{ row.points > 0 ? '+' : '' }}{{ row.points }}
                </span>
              </template>
            </el-table-column>
            <el-table-column prop="balance" label="余额" width="100" />
            <el-table-column prop="createdAt" label="时间" width="180">
              <template #default="{ row }">
                {{ formatTime(row.createdAt) }}
              </template>
            </el-table-column>
          </el-table>

          <el-pagination
            class="pagination"
            v-model:current-page="page"
            v-model:page-size="pageSize"
            :total="total"
            layout="prev, pager, next"
            @current-change="fetchLogs"
          />
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { rewardApi } from '@/api'
import type { PointLog, Reward } from '@/types'
import dayjs from 'dayjs'

const userStore = useUserStore()

const logs = ref<PointLog[]>([])
const rewards = ref<Reward[]>([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const activeTab = ref('all')

const levelThresholds = [0, 100, 500, 1000, 2000, 5000, 10000, 20000, 50000, 100000]

const levelProgress = computed(() => {
  const points = userStore.userInfo?.points || 0
  const level = userStore.userInfo?.level || 1
  const currentThreshold = levelThresholds[level - 1] || 0
  const nextThreshold = levelThresholds[level] || 100000
  return Math.min(100, Math.floor(((points - currentThreshold) / (nextThreshold - currentThreshold)) * 100))
})

const nextLevelPoints = computed(() => {
  const points = userStore.userInfo?.points || 0
  const level = userStore.userInfo?.level || 1
  const nextThreshold = levelThresholds[level] || 100000
  return Math.max(0, nextThreshold - points)
})

const filteredLogs = computed(() => {
  if (activeTab.value === 'earn') {
    return logs.value.filter(log => log.points > 0)
  } else if (activeTab.value === 'spend') {
    return logs.value.filter(log => log.points < 0)
  }
  return logs.value
})

const fetchLogs = async () => {
  try {
    const res = await rewardApi.getPointLogs({
      page: page.value,
      pageSize: pageSize.value
    })
    logs.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch (e) {
    console.error(e)
  }
}

const fetchRewards = async () => {
  try {
    const res = await rewardApi.getRewardList({ page: 1, pageSize: 50 })
    rewards.value = res.data?.list || []
  } catch (e) {
    console.error(e)
  }
}

const exchangeReward = async (id: number) => {
  try {
    await rewardApi.exchangeReward(id)
    fetchLogs()
    fetchRewards()
    if (userStore.userInfo) {
      const reward = rewards.value.find(r => r.id === id)
      if (reward) {
        userStore.userInfo.points -= reward.pointsCost
        userStore.setUserInfo(userStore.userInfo)
      }
    }
  } catch (e) {
    console.error(e)
  }
}

const getPointTypeTag = (type: string) => {
  const map: Record<string, string> = {
    register: 'success',
    answer_created: 'primary',
    answer_accepted: 'success',
    reward_received: 'warning',
    reward_question: 'danger',
    reward_exchange: 'info',
    like_received: 'success'
  }
  return map[type] || ''
}

const getPointTypeName = (type: string) => {
  const map: Record<string, string> = {
    register: '注册奖励',
    answer_created: '回答问题',
    answer_accepted: '回答被采纳',
    reward_received: '获得悬赏',
    reward_question: '设置悬赏',
    reward_exchange: '积分兑换',
    like_received: '获得点赞'
  }
  return map[type] || type
}

const formatTime = (time: string) => {
  return dayjs(time).format('YYYY-MM-DD HH:mm:ss')
}

onMounted(() => {
  fetchLogs()
  fetchRewards()
})
</script>

<style scoped lang="scss">
.points-page {
  .points-summary {
    text-align: center;
    padding: 20px 0;

    .current-points {
      margin-bottom: 20px;

      .label {
        display: block;
        font-size: 14px;
        color: #909399;
        margin-bottom: 8px;
      }

      .value {
        display: block;
        font-size: 48px;
        font-weight: bold;
        color: #409eff;
      }
    }

    .level-info {
      padding: 0 20px;

      .next-level {
        display: block;
        font-size: 12px;
        color: #909399;
        margin-top: 8px;
      }
    }
  }

  .reward-list {
    .reward-item {
      display: flex;
      align-items: center;
      gap: 12px;
      padding: 12px 0;
      border-bottom: 1px solid #e4e7ed;

      &:last-child {
        border-bottom: none;
      }

      .reward-icon {
        font-size: 32px;
      }

      .reward-info {
        flex: 1;

        .reward-name {
          font-weight: 500;
        }

        .reward-cost {
          font-size: 12px;
          color: #e6a23c;
        }
      }
    }
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .pagination {
    justify-content: center;
    margin-top: 20px;
  }

  .positive {
    color: #67c23a;
  }

  .negative {
    color: #f56c6c;
  }
}
</style>
