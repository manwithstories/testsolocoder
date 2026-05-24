<template>
  <div class="discount-policies">
    <div class="page-header flex-between">
      <h2 class="page-title">优惠策略管理</h2>
      <el-button type="primary" @click="showCreateDialog = true">
        <el-icon><Plus /></el-icon>
        添加策略
      </el-button>
    </div>

    <div class="card">
      <el-table :data="policies" style="width: 100%" v-loading="loading">
        <el-table-column prop="name" label="策略名称" />
        <el-table-column prop="code" label="优惠码" width="150" />
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            {{ row.type === 'fixed' ? '固定金额' : '折扣比例' }}
          </template>
        </el-table-column>
        <el-table-column prop="value" label="优惠值" width="120">
          <template #default="{ row }">
            {{ row.type === 'fixed' ? `¥${row.value}` : `${row.value}%` }}
          </template>
        </el-table-column>
        <el-table-column prop="minAmount" label="最低消费" width="120">
          <template #default="{ row }">
            {{ formatMoney(row.minAmount) }}
          </template>
        </el-table-column>
        <el-table-column prop="maxDiscount" label="最大优惠" width="120">
          <template #default="{ row }">
            {{ formatMoney(row.maxDiscount) }}
          </template>
        </el-table-column>
        <el-table-column label="有效期" width="200">
          <template #default="{ row }">
            {{ formatDate(row.startDate) }} - {{ formatDate(row.endDate) }}
          </template>
        </el-table-column>
        <el-table-column prop="isActive" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.isActive ? 'success' : 'info'">
              {{ row.isActive ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="editPolicy(row)">编辑</el-button>
            <el-button type="danger" link @click="deletePolicy(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-dialog v-model="showCreateDialog" :title="isEdit ? '编辑优惠策略' : '添加优惠策略'" width="500px">
      <el-form ref="policyFormRef" :model="policyForm" :rules="policyRules" label-width="100px">
        <el-form-item label="策略名称" prop="name">
          <el-input v-model="policyForm.name" />
        </el-form-item>
        <el-form-item label="优惠码" prop="code">
          <el-input v-model="policyForm.code" />
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-radio-group v-model="policyForm.type">
            <el-radio label="fixed">固定金额</el-radio>
            <el-radio label="percent">折扣比例</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="优惠值" prop="value">
          <el-input-number v-model="policyForm.value" :min="0" :precision="2" />
          <span v-if="policyForm.type === 'percent'" style="margin-left: 8px">%</span>
        </el-form-item>
        <el-form-item label="最低消费" prop="minAmount">
          <el-input-number v-model="policyForm.minAmount" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="最大优惠" prop="maxDiscount">
          <el-input-number v-model="policyForm.maxDiscount" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="开始日期" prop="startDate">
          <el-date-picker v-model="policyForm.startDate" type="date" value-format="YYYY-MM-DD" />
        </el-form-item>
        <el-form-item label="结束日期" prop="endDate">
          <el-date-picker v-model="policyForm.endDate" type="date" value-format="YYYY-MM-DD" />
        </el-form-item>
        <el-form-item label="启用状态">
          <el-switch v-model="policyForm.isActive" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSavePolicy">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { feeApi } from '@/api/fee'
import { DiscountPolicy } from '@/types'
import dayjs from 'dayjs'

const loading = ref(false)
const policies = ref<DiscountPolicy[]>([])
const showCreateDialog = ref(false)
const isEdit = ref(false)
const saving = ref(false)
const policyFormRef = ref<FormInstance>()

const policyForm = reactive({
  id: 0,
  name: '',
  code: '',
  type: 'fixed',
  value: 0,
  minAmount: 0,
  maxDiscount: 0,
  startDate: '',
  endDate: '',
  isActive: true
})

const policyRules: FormRules = {
  name: [{ required: true, message: '请输入策略名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入优惠码', trigger: 'blur' }],
  type: [{ required: true, message: '请选择类型', trigger: 'change' }],
  value: [{ required: true, message: '请输入优惠值', trigger: 'blur' }]
}

const fetchData = async () => {
  loading.value = true
  try {
    const res = await feeApi.getDiscounts()
    policies.value = res || []
  } catch (error) {
    console.error('获取优惠策略失败:', error)
  } finally {
    loading.value = false
  }
}

const formatMoney = (amount: number) => {
  return `¥${amount.toLocaleString()}`
}

const formatDate = (date?: string | null) => {
  if (!date) return '长期有效'
  return dayjs(date).format('YYYY-MM-DD')
}

const editPolicy = (row: DiscountPolicy) => {
  isEdit.value = true
  Object.assign(policyForm, row)
  showCreateDialog.value = true
}

const deletePolicy = async (row: DiscountPolicy) => {
  try {
    await ElMessageBox.confirm(`确认删除策略 ${row.name}？`, '提示', {
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await feeApi.deleteDiscount(row.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  }
}

const handleSavePolicy = async () => {
  if (!policyFormRef.value) return

  await policyFormRef.value.validate(async (valid) => {
    if (valid) {
      saving.value = true
      try {
        if (isEdit.value) {
          await feeApi.updateDiscount(policyForm.id, policyForm)
        } else {
          await feeApi.createDiscount(policyForm as any)
        }
        ElMessage.success('保存成功')
        showCreateDialog.value = false
        fetchData()
      } catch (error: any) {
        ElMessage.error(error.message || '保存失败')
      } finally {
        saving.value = false
      }
    }
  })
}

onMounted(fetchData)
</script>
