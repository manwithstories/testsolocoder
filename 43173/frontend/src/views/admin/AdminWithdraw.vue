<template>
  <div class="admin-withdraw">
    <div class="page-header">
      <h2>提现审核</h2>
    </div>
    
    <div class="filter-bar">
      <el-tabs v-model="activeTab" @tab-change="loadWithdrawList">
        <el-tab-pane label="待审核" name="0" />
        <el-tab-pane label="已通过" name="1" />
        <el-tab-pane label="已拒绝" name="2" />
        <el-tab-pane label="已打款" name="3" />
      </el-tabs>
    </div>
    
    <el-table :data="withdrawList" v-loading="loading" style="width: 100%">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="user_id" label="用户ID" width="100" />
      <el-table-column prop="amount" label="申请金额" width="120">
        <template #default="{ row }">
          ¥{{ row.amount.toFixed(2) }}
        </template>
      </el-table-column>
      <el-table-column prop="actual_amount" label="实际到账" width="120">
        <template #default="{ row }">
          ¥{{ row.actual_amount.toFixed(2) }}
        </template>
      </el-table-column>
      <el-table-column prop="method" label="方式" width="100">
        <template #default="{ row }">
          {{ getMethodText(row.method) }}
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getStatusTag(row.status)">{{ getStatusText(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="申请时间" width="180" />
      <el-table-column label="操作" width="200">
        <template #default="{ row }">
          <template v-if="row.status === 0">
            <el-button 
              type="success" 
              text 
              size="small"
              @click="approveWithdraw(row)"
            >
              通过
            </el-button>
            <el-button 
              type="danger" 
              text 
              size="small"
              @click="rejectWithdraw(row)"
            >
              拒绝
            </el-button>
          </template>
          <template v-else-if="row.status === 1">
            <el-button 
              type="primary" 
              text 
              size="small"
              @click="markPaid(row)"
            >
              标记已打款
            </el-button>
          </template>
          <el-button 
            type="primary" 
            text 
            size="small"
            @click="viewDetail(row)"
          >
            详情
          </el-button>
        </template>
      </el-table-column>
    </el-table>
    
    <div class="pagination">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="loadWithdrawList"
      />
    </div>
    
    <el-dialog v-model="detailDialog.visible" title="提现详情" width="500px">
      <el-descriptions :column="1" border>
        <el-descriptions-item label="用户ID">
          {{ detailDialog.data.user_id }}
        </el-descriptions-item>
        <el-descriptions-item label="申请金额">
          ¥{{ detailDialog.data.amount?.toFixed(2) }}
        </el-descriptions-item>
        <el-descriptions-item label="实际到账">
          ¥{{ detailDialog.data.actual_amount?.toFixed(2) }}
        </el-descriptions-item>
        <el-descriptions-item label="提现方式">
          {{ getMethodText(detailDialog.data.method) }}
        </el-descriptions-item>
        <el-descriptions-item label="收款账号">
          {{ detailDialog.data.account }}
        </el-descriptions-item>
        <el-descriptions-item label="账户名称">
          {{ detailDialog.data.account_name }}
        </el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusTag(detailDialog.data.status)">
            {{ getStatusText(detailDialog.data.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="申请时间">
          {{ detailDialog.data.created_at }}
        </el-descriptions-item>
        <el-descriptions-item v-if="detailDialog.data.reason" label="拒绝原因">
          {{ detailDialog.data.reason }}
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { revenueApi } from '@/api/revenue'

const loading = ref(false)
const withdrawList = ref<any[]>([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const activeTab = ref('0')

const detailDialog = reactive({
  visible: false,
  data: {} as any
})

onMounted(() => {
  loadWithdrawList()
})

async function loadWithdrawList() {
  loading.value = true
  try {
    const res = await revenueApi.getAdminWithdrawList({
      page: page.value,
      page_size: pageSize.value,
      status: parseInt(activeTab.value)
    })
    withdrawList.value = res.list
    total.value = res.total
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function approveWithdraw(item: any) {
  try {
    await ElMessageBox.confirm('确定要通过该提现申请吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'success'
    })
    
    await revenueApi.approveWithdraw(item.id)
    ElMessage.success('审核通过')
    loadWithdrawList()
  } catch (e) {
    if (e !== 'cancel') {
      console.error(e)
    }
  }
}

async function rejectWithdraw(item: any) {
  try {
    const { value } = await ElMessageBox.prompt('请输入拒绝原因', '拒绝提现', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      inputPattern: /.+/,
      inputErrorMessage: '请输入拒绝原因'
    })
    
    await revenueApi.rejectWithdraw(item.id, { reason: value })
    ElMessage.success('已拒绝')
    loadWithdrawList()
  } catch (e) {
    if (e !== 'cancel') {
      console.error(e)
    }
  }
}

async function markPaid(item: any) {
  try {
    await ElMessageBox.confirm('确定要标记为已打款吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'success'
    })
    
    await revenueApi.markWithdrawPaid(item.id)
    ElMessage.success('操作成功')
    loadWithdrawList()
  } catch (e) {
    if (e !== 'cancel') {
      console.error(e)
    }
  }
}

function viewDetail(item: any) {
  detailDialog.data = { ...item }
  detailDialog.visible = true
}

function getMethodText(method: string) {
  const texts: Record<string, string> = {
    alipay: '支付宝',
    wechat: '微信',
    bank: '银行卡'
  }
  return texts[method] || method
}

function getStatusTag(status: number): 'primary' | 'success' | 'warning' | 'info' | 'danger' {
  const tags: Record<number, 'primary' | 'success' | 'warning' | 'info' | 'danger'> = {
    0: 'warning',
    1: 'success',
    2: 'danger',
    3: 'success',
    4: 'info'
  }
  return tags[status] || 'info'
}

function getStatusText(status: number) {
  const texts: Record<number, string> = {
    0: '待审核',
    1: '已通过',
    2: '已拒绝',
    3: '已打款',
    4: '失败'
  }
  return texts[status] || '未知'
}
</script>

<style scoped lang="scss">
.admin-withdraw {
  .filter-bar {
    margin-bottom: 16px;
  }
  
  .pagination {
    display: flex;
    justify-content: center;
    margin-top: 20px;
  }
}
</style>
