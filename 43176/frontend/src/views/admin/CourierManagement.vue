<template>
  <div class="courier-management-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <el-icon><Postcard /></el-icon>
          <span>跑腿员审核</span>
        </div>
      </template>

      <el-tabs v-model="activeTab">
        <el-tab-pane label="待审核" name="pending">
          <el-table :data="pendingCouriers" v-loading="loading" style="width: 100%">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column label="用户" width="200">
              <template #default="{ row }">
                <div class="user-info">
                  <el-avatar :size="40" :src="row.user?.avatar">
                    {{ row.user?.nickname?.charAt(0) }}
                  </el-avatar>
                  <div class="user-detail">
                    <span class="nickname">{{ row.user?.nickname }}</span>
                    <span class="phone">{{ row.user?.phone }}</span>
                  </div>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="level" label="等级" width="80" />
            <el-table-column label="评分" width="120">
              <template #default="{ row }">
                <el-rate :model-value="row.rating || 5" disabled size="small" />
              </template>
            </el-table-column>
            <el-table-column label="审核信息" min-width="200">
              <template #default="{ row }">
                <div class="verify-info">
                  <span>姓名: {{ row.verification?.id_card_name || '-' }}</span>
                  <span>身份证: {{ maskIDCard(row.verification?.id_card_no) }}</span>
                </div>
              </template>
            </el-table-column>
            <el-table-column label="证件照" width="200">
              <template #default="{ row }">
                <div class="id-images">
                  <el-image
                    v-if="row.verification?.id_card_front"
                    :src="row.verification?.id_card_front"
                    :preview-src-list="[row.verification?.id_card_front, row.verification?.id_card_back]"
                    fit="cover"
                    class="id-image"
                  />
                  <el-image
                    v-if="row.verification?.id_card_back"
                    :src="row.verification?.id_card_back"
                    :preview-src-list="[row.verification?.id_card_front, row.verification?.id_card_back]"
                    fit="cover"
                    class="id-image"
                  />
                </div>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="200">
              <template #default="{ row }">
                <el-button type="success" size="small" @click="handleApprove(row)">
                  通过
                </el-button>
                <el-button type="danger" size="small" @click="handleReject(row)">
                  拒绝
                </el-button>
              </template>
            </el-table-column>
          </el-table>
          <el-empty v-if="pendingCouriers.length === 0" description="暂无待审核的跑腿员" />
        </el-tab-pane>
        <el-tab-pane label="已通过" name="approved">
          <el-table :data="approvedCouriers" v-loading="loading" style="width: 100%">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column label="用户" width="200">
              <template #default="{ row }">
                <div class="user-info">
                  <el-avatar :size="40" :src="row.user?.avatar">
                    {{ row.user?.nickname?.charAt(0) }}
                  </el-avatar>
                  <div class="user-detail">
                    <span class="nickname">{{ row.user?.nickname }}</span>
                    <span class="phone">{{ row.user?.phone }}</span>
                  </div>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="level" label="等级" width="80" />
            <el-table-column label="评分" width="120">
              <template #default="{ row }">
                <el-rate :model-value="row.rating || 5" disabled size="small" />
              </template>
            </el-table-column>
            <el-table-column prop="total_orders" label="总订单" width="100" />
            <el-table-column prop="completed_orders" label="已完成" width="100" />
          </el-table>
          <el-empty v-if="approvedCouriers.length === 0" description="暂无已通过的跑腿员" />
        </el-tab-pane>
        <el-tab-pane label="已拒绝" name="rejected">
          <el-table :data="rejectedCouriers" v-loading="loading" style="width: 100%">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column label="用户" width="200">
              <template #default="{ row }">
                <div class="user-info">
                  <el-avatar :size="40" :src="row.user?.avatar">
                    {{ row.user?.nickname?.charAt(0) }}
                  </el-avatar>
                  <div class="user-detail">
                    <span class="nickname">{{ row.user?.nickname }}</span>
                    <span class="phone">{{ row.user?.phone }}</span>
                  </div>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="level" label="等级" width="80" />
            <el-table-column label="拒绝原因" min-width="200">
              <template #default="{ row }">
                {{ row.verification?.reason || '-' }}
              </template>
            </el-table-column>
          </el-table>
          <el-empty v-if="rejectedCouriers.length === 0" description="暂无已拒绝的跑腿员" />
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Postcard } from '@element-plus/icons-vue'
import { adminApi } from '@/api'

const loading = ref(false)
const activeTab = ref('pending')
const couriers = ref<any[]>([])

const pendingCouriers = computed(() => couriers.value.filter(c => c.status === 'pending'))
const approvedCouriers = computed(() => couriers.value.filter(c => c.status === 'approved'))
const rejectedCouriers = computed(() => couriers.value.filter(c => c.status === 'rejected'))

const maskIDCard = (idCard?: string) => {
  if (!idCard) return '-'
  return idCard.replace(/^(\d{4})\d+(\d{4})$/, '$1**********$2')
}

const fetchCouriers = async () => {
  loading.value = true
  try {
    const res = await adminApi.listUsers({ role: 'courier' })
    if (res.code === 200) {
      couriers.value = res.data.items
    }
  } catch (error) {
    console.error('Failed to fetch couriers:', error)
  } finally {
    loading.value = false
  }
}

const handleApprove = async (courier: any) => {
  try {
    await ElMessageBox.confirm(`确定要通过 ${courier.user?.nickname} 的审核吗？`, '审核确认', {
      type: 'success'
    })
    const res = await adminApi.approveCourier(courier.user_id)
    if (res.code === 200) {
      ElMessage.success('审核通过')
      fetchCouriers()
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to approve courier:', error)
    }
  }
}

const handleReject = async (courier: any) => {
  try {
    await ElMessageBox.prompt(`请输入拒绝 ${courier.user?.nickname} 的原因`, '拒绝确认', {
      confirmButtonText: '确认拒绝',
      cancelButtonText: '取消',
      inputPattern: /.+/,
      inputErrorMessage: '请输入拒绝原因'
    }).then(async ({ value }) => {
      const res = await adminApi.rejectCourier(courier.user_id, { reason: value })
      if (res.code === 200) {
        ElMessage.success('已拒绝')
        fetchCouriers()
      }
    })
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to reject courier:', error)
    }
  }
}

watch(activeTab, fetchCouriers)

onMounted(() => {
  fetchCouriers()
})
</script>

<style lang="scss" scoped>
.courier-management-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;

  .card-header {
    display: flex;
    align-items: center;
    gap: 8px;
    font-weight: 600;
  }

  .user-info {
    display: flex;
    align-items: center;
    gap: 12px;

    .user-detail {
      display: flex;
      flex-direction: column;

      .nickname {
        font-weight: 500;
      }

      .phone {
        color: #909399;
        font-size: 12px;
      }
    }
  }

  .verify-info {
    display: flex;
    flex-direction: column;
    gap: 4px;

    span {
      font-size: 13px;
    }
  }

  .id-images {
    display: flex;
    gap: 8px;

    .id-image {
      width: 80px;
      height: 60px;
      border-radius: 4px;
    }
  }
}
</style>
