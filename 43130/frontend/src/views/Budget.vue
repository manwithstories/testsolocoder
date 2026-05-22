<template>
  <div class="budget-page">
    <div class="summary-section">
      <el-row :gutter="20">
        <el-col :span="6">
          <el-card shadow="hover">
            <div class="stat-item">
              <p class="stat-label">总预算</p>
              <p class="stat-value primary">¥{{ formatNumber(summary.total_budget) }}</p>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover">
            <div class="stat-item">
              <p class="stat-label">预计支出</p>
              <p class="stat-value">¥{{ formatNumber(summary.total_estimated) }}</p>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover">
            <div class="stat-item">
              <p class="stat-label">实际支出</p>
              <p class="stat-value" :class="{ danger: summary.over_budget }">
                ¥{{ formatNumber(summary.total_actual) }}
              </p>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover">
            <div class="stat-item">
              <p class="stat-label">已支付</p>
              <p class="stat-value success">¥{{ formatNumber(summary.total_paid) }}</p>
            </div>
          </el-card>
        </el-col>
      </el-row>
      
      <el-alert
        v-if="summary.over_budget"
        title="预算超支预警！"
        type="error"
        :closable="false"
        style="margin-top: 20px"
      >
        当前实际支出已超过预算 ¥{{ formatNumber(summary.total_actual - summary.total_budget) }}
      </el-alert>
    </div>

    <el-card shadow="never" style="margin-top: 20px">
      <div class="page-header">
        <h3>预算明细</h3>
        <el-button type="primary" :icon="Plus" @click="showCreateDialog = true">
          添加预算项
        </el-button>
      </div>

      <el-table :data="budgetItems" v-loading="loading" stripe>
        <el-table-column prop="category" label="类别" width="120" />
        <el-table-column prop="description" label="描述" min-width="160" show-overflow-tooltip />
        <el-table-column label="预计金额" width="120">
          <template #default="{ row }">
            ¥{{ formatNumber(row.estimated_cost) }}
          </template>
        </el-table-column>
        <el-table-column label="实际金额" width="120">
          <template #default="{ row }">
            ¥{{ formatNumber(row.actual_cost) }}
          </template>
        </el-table-column>
        <el-table-column label="已支付" width="120">
          <template #default="{ row }">
            ¥{{ formatNumber(row.paid_amount) }}
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="statusType(row.status)" size="small">{{ statusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="editItem(row)">编辑</el-button>
            <el-button type="success" link @click="recordPayment(row)">付款</el-button>
            <el-button type="danger" link @click="deleteItem(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="showCreateDialog" :title="editingItem ? '编辑预算项' : '添加预算项'" width="500px">
      <el-form ref="itemForm" :model="itemForm" :rules="itemRules" label-width="100px">
        <el-form-item label="类别" prop="category">
          <el-select v-model="itemForm.category" style="width: 100%">
            <el-option v-for="cat in categories" :key="cat" :label="cat" :value="cat" />
          </el-select>
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="itemForm.description" />
        </el-form-item>
        <el-form-item label="预计金额">
          <el-input-number v-model="itemForm.estimated_cost" :min="0" :precision="2" style="width: 100%" />
        </el-form-item>
        <el-form-item label="实际金额">
          <el-input-number v-model="itemForm.actual_cost" :min="0" :precision="2" style="width: 100%" />
        </el-form-item>
        <el-form-item label="到期日期">
          <el-date-picker
            v-model="itemForm.due_date"
            type="date"
            style="width: 100%"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="itemForm.notes" type="textarea" :rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="saveItem">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showPaymentDialog" title="记录付款" width="400px">
      <el-form :model="paymentForm" label-width="80px">
        <el-form-item label="金额">
          <el-input-number v-model="paymentForm.amount" :min="0" :precision="2" style="width: 100%" />
        </el-form-item>
        <el-form-item label="支付方式">
          <el-select v-model="paymentForm.method" style="width: 100%">
            <el-option label="现金" value="cash" />
            <el-option label="银行转账" value="transfer" />
            <el-option label="支付宝" value="alipay" />
            <el-option label="微信" value="wechat" />
            <el-option label="信用卡" value="credit_card" />
          </el-select>
        </el-form-item>
        <el-form-item label="支付日期">
          <el-date-picker
            v-model="paymentForm.paid_at"
            type="date"
            style="width: 100%"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="paymentForm.notes" type="textarea" :rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showPaymentDialog = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="savePayment">确认付款</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { budgetApi } from '@/api/budget'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import type { BudgetItem } from '@/types'

const props = defineProps<{
  weddingId: number
}>()

const loading = ref(false)
const saving = ref(false)
const budgetItems = ref<BudgetItem[]>([])
const categories = ref<string[]>([])
const summary = ref<any>({
  total_budget: 0,
  total_estimated: 0,
  total_actual: 0,
  total_paid: 0,
  over_budget: false
})
const showCreateDialog = ref(false)
const showPaymentDialog = ref(false)
const editingItem = ref<BudgetItem | null>(null)
const currentItem = ref<BudgetItem | null>(null)

const itemForm = reactive({
  category: '',
  description: '',
  estimated_cost: 0,
  actual_cost: 0,
  due_date: '',
  notes: ''
})

const paymentForm = reactive({
  amount: 0,
  method: '',
  paid_at: '',
  notes: ''
})

const itemRules: FormRules = {
  category: [{ required: true, message: '请选择类别', trigger: 'change' }]
}

const itemFormRef = ref<FormInstance>()

function formatNumber(num?: number) {
  return num?.toLocaleString() || '0'
}

function statusType(status: string) {
  const types: Record<string, string> = {
    pending: 'warning',
    partial: 'primary',
    paid: 'success',
    cancelled: 'info'
  }
  return types[status] || 'info'
}

function statusText(status: string) {
  const texts: Record<string, string> = {
    pending: '待支付',
    partial: '部分支付',
    paid: '已付清',
    cancelled: '已取消'
  }
  return texts[status] || status
}

async function fetchCategories() {
  try {
    const res = await budgetApi.getCategories(props.weddingId)
    categories.value = res.data
  } catch (error) {
    console.error('Failed to fetch categories:', error)
  }
}

async function fetchSummary() {
  try {
    const res = await budgetApi.getSummary(props.weddingId)
    summary.value = res.data
  } catch (error) {
    console.error('Failed to fetch summary:', error)
  }
}

async function fetchItems() {
  loading.value = true
  try {
    const res = await budgetApi.getItems(props.weddingId)
    budgetItems.value = res.data
  } catch (error) {
    console.error('Failed to fetch budget items:', error)
  } finally {
    loading.value = false
  }
}

function editItem(item: BudgetItem) {
  editingItem.value = item
  Object.assign(itemForm, {
    category: item.category,
    description: item.description || '',
    estimated_cost: item.estimated_cost,
    actual_cost: item.actual_cost,
    due_date: item.due_date ? item.due_date.split('T')[0] : '',
    notes: item.notes || ''
  })
  showCreateDialog.value = true
}

function recordPayment(item: BudgetItem) {
  currentItem.value = item
  paymentForm.amount = item.actual_cost - item.paid_amount
  paymentForm.method = ''
  paymentForm.paid_at = ''
  paymentForm.notes = ''
  showPaymentDialog.value = true
}

async function saveItem() {
  if (!itemFormRef.value) return
  
  await itemFormRef.value.validate(async (valid) => {
    if (valid) {
      saving.value = true
      try {
        if (editingItem.value) {
          await budgetApi.updateItem(props.weddingId, editingItem.value.id, itemForm)
          ElMessage.success('更新成功')
        } else {
          await budgetApi.createItem(props.weddingId, itemForm)
          ElMessage.success('创建成功')
        }
        showCreateDialog.value = false
        fetchItems()
        fetchSummary()
      } catch (error: any) {
        ElMessage.error(error.message || '保存失败')
      } finally {
        saving.value = false
      }
    }
  })
}

async function savePayment() {
  if (!currentItem.value || paymentForm.amount <= 0) {
    ElMessage.warning('请输入有效金额')
    return
  }
  
  saving.value = true
  try {
    await budgetApi.recordPayment(props.weddingId, {
      budget_item_id: currentItem.value.id,
      amount: paymentForm.amount,
      method: paymentForm.method,
      paid_at: paymentForm.paid_at,
      notes: paymentForm.notes
    })
    ElMessage.success('付款记录成功')
    showPaymentDialog.value = false
    fetchItems()
    fetchSummary()
  } catch (error: any) {
    ElMessage.error(error.message || '记录失败')
  } finally {
    saving.value = false
  }
}

async function deleteItem(item: BudgetItem) {
  try {
    await ElMessageBox.confirm(`确定要删除此预算项吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await budgetApi.deleteItem(props.weddingId, item.id)
    ElMessage.success('删除成功')
    fetchItems()
    fetchSummary()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete item:', error)
    }
  }
}

onMounted(() => {
  fetchCategories()
  fetchSummary()
  fetchItems()
})
</script>

<style scoped>
.budget-page {
  padding: 0;
}

.summary-section {
  margin-bottom: 20px;
}

.stat-item {
  text-align: center;
}

.stat-label {
  color: #909399;
  font-size: 14px;
  margin: 0 0 8px;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  margin: 0;
}

.stat-value.primary {
  color: #409EFF;
}

.stat-value.success {
  color: #67C23A;
}

.stat-value.danger {
  color: #F56C6C;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-header h3 {
  margin: 0;
}
</style>
