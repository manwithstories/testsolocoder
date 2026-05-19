<template>
  <div class="goals-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>阅读目标</span>
          <el-button type="primary" @click="showAddGoal">
            <el-icon><Plus /></el-icon>
            新建目标
          </el-button>
        </div>
      </template>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-card shadow="hover" class="goal-card">
            <template #header>
              <div class="goal-header">
                <span class="year">{{ selectedYear }} 年度目标</span>
                <el-button type="primary" link size="small" @click="editYearlyGoal">
                  <el-icon><Edit /></el-icon>
                  编辑
                </el-button>
              </div>
            </template>
            <div v-if="yearlyGoal" class="goal-content">
              <div class="goal-item">
                <div class="goal-label">读书目标（本）</div>
                <div class="goal-progress">
                  <div class="progress-text">
                    <span class="completed">{{ yearlyGoal.completed_books }}</span>
                    <span class="target">/ {{ yearlyGoal.target_books }}</span>
                  </div>
                  <el-progress
                    :percentage="Math.round(yearlyGoal.book_progress)"
                    :stroke-width="12"
                    :color="progressColor(yearlyGoal.book_progress)"
                  />
                </div>
              </div>
              <div class="goal-item">
                <div class="goal-label">阅读页数目标</div>
                <div class="goal-progress">
                  <div class="progress-text">
                    <span class="completed">{{ yearlyGoal.completed_pages }}</span>
                    <span class="target">/ {{ yearlyGoal.target_pages }}</span>
                  </div>
                  <el-progress
                    :percentage="Math.round(yearlyGoal.page_progress)"
                    :stroke-width="12"
                    :color="progressColor(yearlyGoal.page_progress)"
                  />
                </div>
              </div>
            </div>
            <div v-else class="empty-goal">
              <el-empty description="还未设置年度目标" :image-size="80">
                <el-button type="primary" @click="showAddGoal">设置年度目标</el-button>
              </el-empty>
            </div>
          </el-card>
        </el-col>

        <el-col :span="12">
          <el-card shadow="hover">
            <template #header>
              <span>各月完成情况</span>
            </template>
            <div class="monthly-goals">
              <div v-loading="loading" class="monthly-list">
                <div v-for="goal in monthlyGoals" :key="goal.id" class="monthly-item">
                  <div class="month-name">{{ getMonthName(goal.month) }}</div>
                  <div class="month-progress">
                    <el-progress
                      :percentage="Math.round(goal.book_progress)"
                      :stroke-width="10"
                      :color="progressColor(goal.book_progress)"
                    />
                    <span class="progress-label">{{ goal.completed_books }}/{{ goal.target_books }} 本</span>
                  </div>
                </div>
                <div v-if="!monthlyGoals.length && !loading" class="empty">
                  暂无月度目标
                </div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </el-card>

    <el-dialog v-model="goalDialogVisible" :title="editingGoal ? '编辑目标' : '新建目标'" width="450px">
      <el-form :model="goalForm" label-width="100px">
        <el-form-item label="类型">
          <el-radio-group v-model="goalForm.type">
            <el-radio label="yearly">年度目标</el-radio>
            <el-radio label="monthly">月度目标</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="goalForm.type === 'monthly'" label="月份">
          <el-select v-model="goalForm.month" placeholder="选择月份">
            <el-option v-for="m in 12" :key="m" :label="getMonthName(m)" :value="m" />
          </el-select>
        </el-form-item>
        <el-form-item label="读书目标">
          <el-input-number v-model="goalForm.target_books" :min="0" />
          <span style="margin-left: 8px; color: #909399">本</span>
        </el-form-item>
        <el-form-item label="页数目标">
          <el-input-number v-model="goalForm.target_pages" :min="0" />
          <span style="margin-left: 8px; color: #909399">页</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="goalDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveGoal">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getYearlyGoalProgress, createGoal, updateGoal } from '@/api/common'
import type { GoalProgress } from '@/types'

const selectedYear = ref(new Date().getFullYear())
const yearlyGoal = ref<GoalProgress | null>(null)
const monthlyGoals = ref<GoalProgress[]>([])
const loading = ref(false)

const goalDialogVisible = ref(false)
const editingGoal = ref<GoalProgress | null>(null)
const goalForm = reactive({
  type: 'yearly' as 'yearly' | 'monthly',
  month: 1,
  target_books: 12,
  target_pages: 3000
})

const loadGoals = async () => {
  loading.value = true
  try {
    const data = await getYearlyGoalProgress(selectedYear.value)
    yearlyGoal.value = data.yearly
    monthlyGoals.value = data.monthly
  } catch (e) {
  } finally {
    loading.value = false
  }
}

const progressColor = (percentage: number) => {
  if (percentage < 30) return '#e6a23c'
  if (percentage < 70) return '#409eff'
  return '#67c23a'
}

const getMonthName = (month?: number) => {
  if (!month) return ''
  const months = ['一月', '二月', '三月', '四月', '五月', '六月', '七月', '八月', '九月', '十月', '十一月', '十二月']
  return months[month - 1]
}

const showAddGoal = () => {
  editingGoal.value = null
  goalForm.type = 'yearly'
  goalForm.month = 1
  goalForm.target_books = yearlyGoal.value?.target_books || 12
  goalForm.target_pages = yearlyGoal.value?.target_pages || 3000
  goalDialogVisible.value = true
}

const editYearlyGoal = () => {
  if (yearlyGoal.value) {
    editingGoal.value = yearlyGoal.value
    goalForm.type = 'yearly'
    goalForm.target_books = yearlyGoal.value.target_books
    goalForm.target_pages = yearlyGoal.value.target_pages
  } else {
    editingGoal.value = null
    goalForm.type = 'yearly'
    goalForm.target_books = 12
    goalForm.target_pages = 3000
  }
  goalDialogVisible.value = true
}

const saveGoal = async () => {
  try {
    if (editingGoal.value) {
      await updateGoal(editingGoal.value.id, {
        target_books: goalForm.target_books,
        target_pages: goalForm.target_pages
      })
      ElMessage.success('更新成功')
    } else {
      await createGoal({
        year: selectedYear.value,
        month: goalForm.type === 'monthly' ? goalForm.month : undefined,
        target_books: goalForm.target_books,
        target_pages: goalForm.target_pages
      })
      ElMessage.success('创建成功')
    }
    goalDialogVisible.value = false
    loadGoals()
  } catch (e) {}
}

onMounted(() => {
  loadGoals()
})
</script>

<style scoped lang="scss">
.goals-page {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .goal-card {
    .goal-header {
      display: flex;
      justify-content: space-between;
      align-items: center;

      .year {
        font-size: 16px;
        font-weight: 600;
      }
    }

    .goal-item {
      margin-bottom: 24px;

      &:last-child {
        margin-bottom: 0;
      }

      .goal-label {
        font-size: 14px;
        color: #909399;
        margin-bottom: 8px;
      }

      .progress-text {
        margin-bottom: 8px;

        .completed {
          font-size: 24px;
          font-weight: 600;
          color: #409eff;
        }

        .target {
          font-size: 16px;
          color: #909399;
          margin-left: 4px;
        }
      }
    }

    .empty-goal {
      padding: 20px 0;
    }
  }

  .monthly-list {
    max-height: 400px;
    overflow-y: auto;
  }

  .monthly-item {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 12px 0;
    border-bottom: 1px solid #f0f2f5;

    &:last-child {
      border-bottom: none;
    }

    .month-name {
      width: 60px;
      font-weight: 500;
    }

    .month-progress {
      flex: 1;
      display: flex;
      align-items: center;
      gap: 12px;

      :deep(.el-progress) {
        flex: 1;
      }

      .progress-label {
        width: 80px;
        font-size: 13px;
        color: #909399;
        flex-shrink: 0;
      }
    }
  }

  .empty {
    text-align: center;
    color: #909399;
    padding: 40px 0;
  }
}
</style>
