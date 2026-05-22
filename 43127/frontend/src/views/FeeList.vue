<template>
  <div class="fee-list">
    <div class="page-header">
      <h2 class="page-title">公共费用</h2>
      <el-button type="primary" @click="showCreateDialog">
        <el-icon><Plus /></el-icon>添加费用
      </el-button>
      <el-button type="success" @click="showBatchDialog">
        <el-icon><DocumentCopy /></el-icon>批量生成
      </el-button>
    </div>

    <div class="search-bar card">
      <el-input
        v-model="monthFilter"
        placeholder="选择月份"
        style="width: 150px"
      />
      <el-select v-model="typeFilter" placeholder="类型" clearable style="width: 150px">
        <el-option label="水费" value="水费" />
        <el-option label="电费" value="电费" />
        <el-option label="燃气费" value="燃气费" />
        <el-option label="物业费" value="物业费" />
      </el-select>
      <el-select v-model="statusFilter" placeholder="状态" clearable style="width: 150px">
        <el-option label="未缴" :value="0" />
        <el-option label="已缴" :value="1" />
      </el-select>
      <el-button type="primary" @click="loadData">搜索</el-button>
    </div>

    <div class="card">
      <div class="summary">
        <el-descriptions :column="3" border>
          <el-descriptions-item label="总金额">¥{{ summary.totalAmount }}</el-descriptions-item>
          <el-descriptions-item label="已缴金额">
            <span style="color: #67c23a;">¥{{ summary.paidAmount }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="未缴金额">
            <span style="color: #f56c6c;">¥{{ summary.totalAmount - summary.paidAmount }}</span>
          </el-descriptions-item>
        </el-descriptions>
      </div>

      <el-table :data="fees" v-loading="loading">
        <el-table-column prop="month" label="月份" width="120" />
        <el-table-column prop="type" label="类型" width="100" />
        <el-table-column label="房源" min-width="150">
          <template #default="{ row }">
            {{ row.property?.title }}
          </template>
        </el-table-column>
        <el-table-column prop="totalAmount" label="金额" width="120">
          <template #default="{ row }">¥{{ row.totalAmount }}</template>
        </el-table-column>
        <el-table-column prop="units" label="用量" width="100">
          <template #default="{ row }">{{ row.units || '-' }}</template>
        </el-table-column>
        <el-table-column prop="unitPrice" label="单价" width="100">
          <template #default="{ row }">{{ row.unitPrice || '-' }}</template>
        </el-table-column>
        <el-table-column prop="dueDate" label="应缴日期" width="120">
          <template #default="{ row }">{{ formatDate(row.dueDate) }}</template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '已缴' : '未缴' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button
              link
              type="primary"
              @click="handlePay(row)"
              v-if="row.status === 0"
            >缴费</el-button>
            <el-button link type="danger" @click="deleteFee(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        class="pagination"
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next, jumper"
        @current-change="loadData"
        @size-change="loadData"
      />
    </div>

    <el-dialog v-model="createDialogVisible" title="添加费用" width="500px">
      <el-form :model="feeForm" :rules="feeRules" ref="feeFormRef" label-width="80px">
        <el-form-item label="房源" prop="propertyId">
          <el-select v-model="feeForm.propertyId" placeholder="请选择房源" filterable style="width: 100%">
            <el-option
              v-for="p in properties"
              :key="p.id"
              :label="p.title"
              :value="p.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-select v-model="feeForm.type" placeholder="请选择类型" style="width: 100%">
            <el-option label="水费" value="水费" />
            <el-option label="电费" value="电费" />
            <el-option label="燃气费" value="燃气费" />
            <el-option label="物业费" value="物业费" />
          </el-select>
        </el-form-item>
        <el-form-item label="月份" prop="month">
          <el-input v-model="feeForm.month" placeholder="如: 2024-01" />
        </el-form-item>
        <el-form-item label="金额" prop="totalAmount">
          <el-input-number v-model="feeForm.totalAmount" :min="0" :precision="2" style="width: 100%" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitFee">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="batchDialogVisible" title="批量生成费用" width="500px">
      <el-form :model="batchForm" :rules="batchRules" ref="batchFormRef" label-width="80px">
        <el-form-item label="月份" prop="month">
          <el-input v-model="batchForm.month" placeholder="如: 2024-01" />
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-select v-model="batchForm.type" placeholder="请选择类型" style="width: 100%">
            <el-option label="水费" value="水费" />
            <el-option label="电费" value="电费" />
            <el-option label="燃气费" value="燃气费" />
            <el-option label="物业费" value="物业费" />
          </el-select>
        </el-form-item>
        <el-form-item label="金额" prop="amount">
          <el-input-number v-model="batchForm.amount" :min="0" :precision="2" style="width: 100%" />
        </el-form-item>
        <el-form-item label="选择房源" prop="propertyIds">
          <el-select
            v-model="batchForm.propertyIds"
            multiple
            placeholder="请选择房源"
            filterable
            style="width: 100%"
          >
            <el-option
              v-for="p in properties"
              :key="p.id"
              :label="p.title"
              :value="p.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="batchDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitBatch">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import type { UtilityFee, Property } from '@/types'
import { getFees, createFee, deleteFee as deleteFeeApi, payFee, batchGenerateFees } from '@/api/business'
import { getProperties } from '@/api/property'
import dayjs from 'dayjs'

const loading = ref(false)
const fees = ref<UtilityFee[]>([])
const properties = ref<Property[]>([])
const monthFilter = ref('')
const typeFilter = ref('')
const statusFilter = ref<number | ''>('')
const createDialogVisible = ref(false)
const batchDialogVisible = ref(false)
const feeFormRef = ref<FormInstance>()
const batchFormRef = ref<FormInstance>()

const summary = reactive({
  totalAmount: 0,
  paidAmount: 0
})

const feeForm = reactive({
  propertyId: 0,
  type: '',
  month: '',
  totalAmount: 0
})

const batchForm = reactive({
  month: '',
  type: '',
  amount: 0,
  propertyIds: [] as number[]
})

const feeRules: FormRules = {
  propertyId: [{ required: true, message: '请选择房源', trigger: 'change' }],
  type: [{ required: true, message: '请选择类型', trigger: 'change' }],
  month: [{ required: true, message: '请输入月份', trigger: 'blur' }],
  totalAmount: [{ required: true, message: '请输入金额', trigger: 'blur' }]
}

const batchRules: FormRules = {
  month: [{ required: true, message: '请输入月份', trigger: 'blur' }],
  type: [{ required: true, message: '请选择类型', trigger: 'change' }],
  amount: [{ required: true, message: '请输入金额', trigger: 'blur' }],
  propertyIds: [{ required: true, message: '请选择房源', trigger: 'change' }]
}

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

onMounted(() => {
  loadData()
  loadProperties()
})

async function loadData() {
  loading.value = true
  try {
    const res = await getFees({
      page: pagination.page,
      pageSize: pagination.pageSize,
      month: monthFilter.value || undefined,
      type: typeFilter.value || undefined,
      status: statusFilter.value || undefined
    })
    fees.value = res.data.list
    pagination.total = res.data.total
    summary.totalAmount = res.data.totalAmount || 0
    summary.paidAmount = res.data.paidAmount || 0
  } catch (error) {
    console.error('Failed to load fees:', error)
  } finally {
    loading.value = false
  }
}

async function loadProperties() {
  try {
    const res = await getProperties({ pageSize: 100 })
    properties.value = res.data.list
  } catch (error) {
    console.error('Failed to load properties:', error)
  }
}

function formatDate(date: string) {
  return dayjs(date).format('YYYY-MM-DD')
}

function showCreateDialog() {
  feeForm.propertyId = 0
  feeForm.type = ''
  feeForm.month = ''
  feeForm.totalAmount = 0
  createDialogVisible.value = true
}

async function submitFee() {
  if (!feeFormRef.value) return
  
  await feeFormRef.value.validate(async (valid) => {
    if (valid) {
      try {
        await createFee(feeForm)
        ElMessage.success('创建成功')
        createDialogVisible.value = false
        loadData()
      } catch (error) {
        console.error(error)
      }
    }
  })
}

function showBatchDialog() {
  batchForm.month = ''
  batchForm.type = ''
  batchForm.amount = 0
  batchForm.propertyIds = []
  batchDialogVisible.value = true
}

async function submitBatch() {
  if (!batchFormRef.value) return
  
  await batchFormRef.value.validate(async (valid) => {
    if (valid) {
      try {
        await batchGenerateFees(batchForm)
        ElMessage.success('批量生成成功')
        batchDialogVisible.value = false
        loadData()
      } catch (error) {
        console.error(error)
      }
    }
  })
}

async function handlePay(row: UtilityFee) {
  try {
    await payFee(row.id)
    ElMessage.success('缴费成功')
    loadData()
  } catch (error) {
    console.error(error)
  }
}

async function deleteFee(row: UtilityFee) {
  try {
    await ElMessageBox.confirm('确定要删除该费用记录吗？', '提示', { type: 'warning' })
    await deleteFeeApi(row.id)
    ElMessage.success('删除成功')
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error(error)
    }
  }
}
</script>

<style scoped>
.fee-list {
  padding: 0;
}

.summary {
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
