<template>
  <div class="fee-standards">
    <div class="page-header flex-between">
      <h2 class="page-title">费用标准管理</h2>
    </div>

    <div class="card">
      <el-table :data="standards" style="width: 100%" v-loading="loading">
        <el-table-column prop="companyType" label="公司类型" width="160">
          <template #default="{ row }">
            {{ getCompanyTypeText(row.companyType) }}
          </template>
        </el-table-column>
        <el-table-column prop="namingFee" label="核名费" width="120">
          <template #default="{ row }">
            {{ formatMoney(row.namingFee) }}
          </template>
        </el-table-column>
        <el-table-column prop="registrationFee" label="登记费" width="120">
          <template #default="{ row }">
            {{ formatMoney(row.registrationFee) }}
          </template>
        </el-table-column>
        <el-table-column prop="taxFee" label="税务登记费" width="120">
          <template #default="{ row }">
            {{ formatMoney(row.taxFee) }}
          </template>
        </el-table-column>
        <el-table-column prop="bankFee" label="银行开户费" width="120">
          <template #default="{ row }">
            {{ formatMoney(row.bankFee) }}
          </template>
        </el-table-column>
        <el-table-column prop="sealFee" label="刻章费" width="120">
          <template #default="{ row }">
            {{ formatMoney(row.sealFee) }}
          </template>
        </el-table-column>
        <el-table-column prop="serviceFee" label="服务费" width="120">
          <template #default="{ row }">
            {{ formatMoney(row.serviceFee) }}
          </template>
        </el-table-column>
        <el-table-column prop="capitalRate" label="资本费率(‰)" width="120" />
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="editStandard(row)">编辑</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-dialog v-model="editDialogVisible" title="编辑费用标准" width="500px">
      <el-form :model="editForm" label-width="120px">
        <el-form-item label="公司类型">
          <el-input :value="getCompanyTypeText(editForm.companyType)" disabled />
        </el-form-item>
        <el-form-item label="核名费">
          <el-input-number v-model="editForm.namingFee" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="登记费">
          <el-input-number v-model="editForm.registrationFee" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="税务登记费">
          <el-input-number v-model="editForm.taxFee" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="银行开户费">
          <el-input-number v-model="editForm.bankFee" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="刻章费">
          <el-input-number v-model="editForm.sealFee" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="服务费">
          <el-input-number v-model="editForm.serviceFee" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="资本费率(‰)">
          <el-input-number v-model="editForm.capitalRate" :min="0" :precision="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSaveStandard">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { feeApi } from '@/api/fee'
import { FeeStandard, CompanyType } from '@/types'

const loading = ref(false)
const standards = ref<FeeStandard[]>([])
const editDialogVisible = ref(false)
const saving = ref(false)

const editForm = reactive<FeeStandard>({
  id: 0,
  companyType: 'llc',
  namingFee: 0,
  registrationFee: 0,
  taxFee: 0,
  bankFee: 0,
  sealFee: 0,
  serviceFee: 0,
  capitalRate: 0,
  createdAt: '',
  updatedAt: ''
})

const fetchData = async () => {
  loading.value = true
  try {
    const res = await feeApi.getStandards()
    standards.value = res || []
  } catch (error) {
    console.error('获取费用标准失败:', error)
  } finally {
    loading.value = false
  }
}

const getCompanyTypeText = (type: CompanyType) => {
  const map: Record<CompanyType, string> = {
    llc: '有限责任公司',
    joint_stock: '股份有限公司',
    sole: '个人独资',
    partnership: '合伙企业'
  }
  return map[type] || type
}

const formatMoney = (amount: number) => {
  return `¥${amount.toLocaleString()}`
}

const editStandard = (row: FeeStandard) => {
  Object.assign(editForm, row)
  editDialogVisible.value = true
}

const handleSaveStandard = async () => {
  saving.value = true
  try {
    await feeApi.updateStandard(editForm.id, {
      namingFee: editForm.namingFee,
      registrationFee: editForm.registrationFee,
      taxFee: editForm.taxFee,
      bankFee: editForm.bankFee,
      sealFee: editForm.sealFee,
      serviceFee: editForm.serviceFee,
      capitalRate: editForm.capitalRate
    })
    ElMessage.success('保存成功')
    editDialogVisible.value = false
    fetchData()
  } catch (error: any) {
    ElMessage.error(error.message || '保存失败')
  } finally {
    saving.value = false
  }
}

onMounted(fetchData)
</script>
