<template>
  <Layout>
    <div class="profile-page">
      <div class="profile-header">
        <el-avatar :size="80" src="https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f31cpng.png" />
        <div class="user-info">
          <h2>{{ userStore.user?.username }}</h2>
          <div class="member-badge" v-if="currentLevel">
            <el-tag :type="getLevelTagType(currentLevel.level)" size="large">
              {{ currentLevel.name }}
            </el-tag>
            <span class="discount-text">享 {{ (currentLevel.discount * 10).toFixed(1) }} 折优惠</span>
          </div>
        </div>
      </div>

      <el-row :gutter="24">
        <el-col :span="8">
          <el-card class="stats-card">
            <div class="stat-item">
              <div class="stat-value">{{ userStore.user?.points || 0 }}</div>
              <div class="stat-label">积分余额</div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="8">
          <el-card class="stats-card">
            <div class="stat-item">
              <div class="stat-value">{{ couponCount }}</div>
              <div class="stat-label">可用优惠券</div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="8">
          <el-card class="stats-card">
            <div class="stat-item">
              <div class="stat-value">{{ orderCount }}</div>
              <div class="stat-label">历史订单</div>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <el-tabs v-model="activeTab" class="profile-tabs">
        <el-tab-pane label="基本信息" name="info">
          <el-card>
            <el-form :model="form" label-width="100px" :rules="rules" ref="formRef">
              <el-form-item label="用户名" prop="username">
                <el-input v-model="form.username" disabled />
              </el-form-item>
              <el-form-item label="真实姓名" prop="real_name">
                <el-input v-model="form.real_name" placeholder="请输入真实姓名" />
              </el-form-item>
              <el-form-item label="身份证号" prop="id_card">
                <el-input v-model="form.id_card" placeholder="请输入身份证号" maxlength="18" />
              </el-form-item>
              <el-form-item label="手机号" prop="phone">
                <el-input v-model="form.phone" placeholder="请输入手机号" maxlength="11" />
              </el-form-item>
              <el-form-item label="邮箱" prop="email">
                <el-input v-model="form.email" placeholder="请输入邮箱" />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" @click="handleUpdate">保存修改</el-button>
              </el-form-item>
            </el-form>
          </el-card>
        </el-tab-pane>

        <el-tab-pane label="会员等级" name="level">
          <el-card>
            <div class="level-progress" v-if="currentLevel && nextLevel">
              <div class="level-info">
                <span>当前等级：{{ currentLevel.name }}</span>
                <span>距离下一等级还需 {{ nextLevel.min_points - (userStore.user?.points || 0) }} 积分</span>
              </div>
              <el-progress
                :percentage="Math.min(100, Math.round(((userStore.user?.points || 0) - currentLevel.min_points) / (nextLevel.min_points - currentLevel.min_points) * 100))"
                :stroke-width="20"
              />
            </div>

            <el-table :data="memberLevels" style="width: 100%; margin-top: 20px;">
              <el-table-column prop="level" label="等级" width="100">
                <template #default="{ row }">
                  <el-tag :type="getLevelTagType(row.level)">{{ row.name }}</el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="min_points" label="所需积分" width="120" />
              <el-table-column label="折扣">
                <template #default="{ row }">
                  {{ (row.discount * 10).toFixed(1) }} 折
                </template>
              </el-table-column>
              <el-table-column prop="priority" label="优先购票权" width="120">
                <template #default="{ row }">
                  {{ row.priority > 0 ? `提前${row.priority}天` : '无' }}
                </template>
              </el-table-column>
              <el-table-column prop="description" label="说明" />
              <el-table-column label="状态" width="100">
                <template #default="{ row }">
                  <el-tag v-if="currentLevel && row.level === currentLevel.level" type="success">当前</el-tag>
                  <el-tag v-else-if="userStore.user && userStore.user.points >= row.min_points" type="info">已达成</el-tag>
                  <el-tag v-else type="info" effect="plain">未达成</el-tag>
                </template>
              </el-table-column>
            </el-table>
          </el-card>
        </el-tab-pane>

        <el-tab-pane label="我的优惠券" name="coupons">
          <el-card>
            <div class="coupon-actions">
              <el-button type="primary" @click="showExchangeDialog = true">积分兑换</el-button>
            </div>

            <div class="coupons-list" v-if="coupons.length > 0">
              <div v-for="coupon in coupons" :key="coupon.id" class="coupon-card" :class="{ disabled: coupon.status !== 1 }">
                <div class="coupon-left">
                  <div class="coupon-value">¥{{ coupon.value }}</div>
                  <div class="coupon-condition" v-if="coupon.min_amount > 0">满{{ coupon.min_amount }}可用</div>
                </div>
                <div class="coupon-right">
                  <div class="coupon-name">{{ coupon.name }}</div>
                  <div class="coupon-code">券码：{{ coupon.code }}</div>
                  <div class="coupon-expire">有效期至：{{ formatTime(coupon.expire_at) }}</div>
                </div>
                <div class="coupon-status">
                  <el-tag v-if="coupon.status === 1" type="success">未使用</el-tag>
                  <el-tag v-else type="info">已使用</el-tag>
                </div>
              </div>
            </div>

            <el-empty v-else description="暂无优惠券" />
          </el-card>
        </el-tab-pane>

        <el-tab-pane label="修改密码" name="password">
          <el-card>
            <el-form :model="passwordForm" label-width="100px" :rules="passwordRules" ref="passwordFormRef">
              <el-form-item label="原密码" prop="old_password">
                <el-input v-model="passwordForm.old_password" type="password" show-password />
              </el-form-item>
              <el-form-item label="新密码" prop="new_password">
                <el-input v-model="passwordForm.new_password" type="password" show-password />
              </el-form-item>
              <el-form-item label="确认密码" prop="confirm_password">
                <el-input v-model="passwordForm.confirm_password" type="password" show-password />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" @click="handleUpdatePassword">修改密码</el-button>
              </el-form-item>
            </el-form>
          </el-card>
        </el-tab-pane>
      </el-tabs>

      <el-dialog v-model="showExchangeDialog" title="积分兑换优惠券" width="400px">
        <el-form :model="exchangeForm" label-width="80px">
          <el-form-item label="兑换积分">
            <el-input-number v-model="exchangeForm.points" :min="100" :step="100" />
          </el-form-item>
          <el-form-item>
            <span class="exchange-tip">100积分可兑换10元优惠券</span>
          </el-form-item>
        </el-form>
        <template #footer>
          <el-button @click="showExchangeDialog = false">取消</el-button>
          <el-button type="primary" @click="handleExchange">确认兑换</el-button>
        </template>
      </el-dialog>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElForm } from 'element-plus'
import dayjs from 'dayjs'
import { useUserStore } from '@/store'
import { userApi, orderApi } from '@/api'
import Layout from '@/components/Layout.vue'
import type { MemberLevel, Coupon } from '@/types'

const userStore = useUserStore()

const activeTab = ref('info')
const memberLevels = ref<MemberLevel[]>([])
const coupons = ref<Coupon[]>([])
const couponCount = ref(0)
const orderCount = ref(0)
const showExchangeDialog = ref(false)

const form = reactive({
  username: userStore.user?.username || '',
  real_name: userStore.user?.real_name || '',
  id_card: userStore.user?.id_card || '',
  phone: userStore.user?.phone || '',
  email: userStore.user?.email || ''
})

const passwordForm = reactive({
  old_password: '',
  new_password: '',
  confirm_password: ''
})

const exchangeForm = reactive({
  points: 100
})

const formRef = ref<InstanceType<typeof ElForm>>()
const passwordFormRef = ref<InstanceType<typeof ElForm>>()

const rules = {
  real_name: [{ required: true, message: '请输入真实姓名', trigger: 'blur' }],
  id_card: [
    { required: true, message: '请输入身份证号', trigger: 'blur' },
    { pattern: /(^\d{15}$)|(^\d{18}$)|(^\d{17}(\d|X|x)$)/, message: '身份证号格式不正确', trigger: 'blur' }
  ],
  phone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '手机号格式不正确', trigger: 'blur' }
  ],
  email: [{ type: 'email', message: '邮箱格式不正确', trigger: 'blur' }]
}

const passwordRules = {
  old_password: [{ required: true, message: '请输入原密码', trigger: 'blur' }],
  new_password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ],
  confirm_password: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    {
      validator: (_rule: any, value: string, callback: any) => {
        if (value !== passwordForm.new_password) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

const currentLevel = computed(() => {
  const userPoints = userStore.user?.points || 0
  let level: MemberLevel | null = null
  for (const l of memberLevels.value) {
    if (userPoints >= l.min_points) {
      level = l
    }
  }
  return level
})

const nextLevel = computed(() => {
  const userPoints = userStore.user?.points || 0
  for (const l of memberLevels.value) {
    if (userPoints < l.min_points) {
      return l
    }
  }
  return null
})

async function fetchMemberLevels() {
  try {
    const res = await userApi.getMemberLevels()
    memberLevels.value = res.sort((a: MemberLevel, b: MemberLevel) => a.level - b.level)
  } catch (err) {
    console.error(err)
  }
}

async function fetchCoupons() {
  try {
    const res = await userApi.getCoupons()
    coupons.value = res
    couponCount.value = res.filter((c: Coupon) => c.status === 1).length
  } catch (err) {
    console.error(err)
  }
}

async function fetchOrderCount() {
  try {
    const res = await orderApi.list({ page: 1, page_size: 1 })
    orderCount.value = res.pagination?.total || 0
  } catch (err) {
    console.error(err)
  }
}

async function handleUpdate() {
  try {
    await formRef?.validate()
    await userApi.updateInfo(form)
    ElMessage.success('修改成功')
    userStore.fetchUserInfo()
  } catch (err) {
    console.error(err)
  }
}

async function handleUpdatePassword() {
  try {
    await passwordFormRef?.validate()
    await userApi.changePassword({
      old_password: passwordForm.old_password,
      new_password: passwordForm.new_password
    })
    ElMessage.success('密码修改成功')
    passwordForm.old_password = ''
    passwordForm.new_password = ''
    passwordForm.confirm_password = ''
  } catch (err: any) {
    ElMessage.error(err.message || '密码修改失败')
  }
}

async function handleExchange() {
  try {
    await userApi.exchangeCoupon(exchangeForm)
    ElMessage.success('兑换成功')
    showExchangeDialog.value = false
    fetchCoupons()
    userStore.fetchUserInfo()
  } catch (err) {
    console.error(err)
  }
}

function getLevelTagType(level: number) {
  const types: Record<number, string> = {
    1: 'info',
    2: 'success',
    3: 'warning',
    4: 'danger'
  }
  return types[level] || 'info'
}

function formatTime(time: string) {
  return dayjs(time).format('YYYY-MM-DD')
}

onMounted(() => {
  fetchMemberLevels()
  fetchCoupons()
  fetchOrderCount()
})
</script>

<style lang="scss" scoped>
.profile-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;

  .profile-header {
    display: flex;
    align-items: center;
    gap: 24px;
    padding: 24px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border-radius: 12px;
    margin-bottom: 24px;
    color: white;

    .user-info {
      h2 {
        margin: 0 0 12px 0;
        font-size: 24px;
      }

      .member-badge {
        display: flex;
        align-items: center;
        gap: 12px;

        .discount-text {
          font-size: 14px;
          opacity: 0.9;
        }
      }
    }
  }

  .stats-card {
    text-align: center;

    .stat-item {
      .stat-value {
        font-size: 28px;
        font-weight: bold;
        color: #409eff;
        margin-bottom: 8px;
      }

      .stat-label {
        color: #999;
        font-size: 14px;
      }
    }
  }

  .profile-tabs {
    margin-top: 24px;
  }

  .level-progress {
    margin-bottom: 20px;

    .level-info {
      display: flex;
      justify-content: space-between;
      margin-bottom: 12px;
    }
  }

  .coupon-actions {
    margin-bottom: 20px;
  }

  .coupons-list {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 16px;

    .coupon-card {
      display: flex;
      align-items: center;
      padding: 16px;
      background: linear-gradient(135deg, #ff6b6b 0%, #ee5a24 100%);
      border-radius: 8px;
      color: white;
      position: relative;

      &.disabled {
        background: linear-gradient(135deg, #ccc 0%, #999 100%);
      }

      .coupon-left {
        border-right: 2px dashed rgba(255, 255, 255, 0.3);
        padding-right: 16px;
        margin-right: 16px;
        text-align: center;

        .coupon-value {
          font-size: 32px;
          font-weight: bold;

          &::before {
            content: '¥';
            font-size: 16px;
          }
        }

        .coupon-condition {
          font-size: 12px;
          opacity: 0.9;
        }
      }

      .coupon-right {
        flex: 1;

        .coupon-name {
          font-size: 16px;
          font-weight: 600;
          margin-bottom: 8px;
        }

        .coupon-code,
        .coupon-expire {
          font-size: 12px;
          opacity: 0.9;
          margin-bottom: 4px;
        }
      }

      .coupon-status {
        position: absolute;
        top: 12px;
        right: 12px;
      }
    }
  }

  .exchange-tip {
    color: #999;
    font-size: 12px;
  }
}
</style>
