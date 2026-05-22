<template>
  <div class="page-container">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>我的档案</span>
          <el-button type="primary" @click="editVisible = true">编辑</el-button>
        </div>
      </template>

      <el-descriptions :column="2" border v-if="customer">
        <el-descriptions-item label="姓名">{{ customer.name || '-' }}</el-descriptions-item>
        <el-descriptions-item label="手机号">{{ customer.user?.phone }}</el-descriptions-item>
        <el-descriptions-item label="性别">{{ customer.gender || '-' }}</el-descriptions-item>
        <el-descriptions-item label="年龄">{{ customer.age || '-' }}</el-descriptions-item>
        <el-descriptions-item label="皮肤类型">{{ customer.skin_type || '-' }}</el-descriptions-item>
        <el-descriptions-item label="发型偏好">{{ customer.hair_preference || '-' }}</el-descriptions-item>
        <el-descriptions-item label="过敏史">{{ customer.allergy_history || '-' }}</el-descriptions-item>
        <el-descriptions-item label="备注">{{ customer.notes || '-' }}</el-descriptions-item>
        <el-descriptions-item label="会员等级">
          <el-tag :type="getLevelType(customer.member_level)">
            {{ getLevelText(customer.member_level) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="积分">{{ customer.points }}</el-descriptions-item>
      </el-descriptions>
    </el-card>

    <el-dialog v-model="editVisible" title="编辑档案" width="500px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="姓名">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="性别">
          <el-select v-model="form.gender">
            <el-option label="男" value="male" />
            <el-option label="女" value="female" />
          </el-select>
        </el-form-item>
        <el-form-item label="年龄">
          <el-input-number v-model="form.age" :min="0" :max="120" />
        </el-form-item>
        <el-form-item label="皮肤类型">
          <el-select v-model="form.skin_type">
            <el-option label="干性" value="干性" />
            <el-option label="油性" value="油性" />
            <el-option label="混合性" value="混合性" />
            <el-option label="中性" value="中性" />
            <el-option label="敏感性" value="敏感性" />
          </el-select>
        </el-form-item>
        <el-form-item label="发型偏好">
          <el-input v-model="form.hair_preference" type="textarea" />
        </el-form-item>
        <el-form-item label="过敏史">
          <el-input v-model="form.allergy_history" type="textarea" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.notes" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { getMyCustomer, updateMyCustomer } from '@/api/auth'
import { ElMessage } from 'element-plus'
import type { Customer } from '@/types'

const customer = ref<Customer | null>(null)
const editVisible = ref(false)
const saving = ref(false)
const form = reactive({
  name: '',
  gender: '',
  age: 0,
  skin_type: '',
  hair_preference: '',
  allergy_history: '',
  notes: ''
})

const getLevelType = (level: number) => {
  const types = ['', 'info', '', 'warning', 'success', 'danger']
  return types[level] || 'info'
}

const getLevelText = (level: number) => {
  const texts = ['', '普通会员', '', '银卡会员', '金卡会员', '钻石会员']
  return texts[level] || '普通会员'
}

const fetchCustomer = async () => {
  try {
    const res = await getMyCustomer()
    customer.value = res.data
    Object.assign(form, {
      name: res.data.name || '',
      gender: res.data.gender || '',
      age: res.data.age || 0,
      skin_type: res.data.skin_type || '',
      hair_preference: res.data.hair_preference || '',
      allergy_history: res.data.allergy_history || '',
      notes: res.data.notes || ''
    })
  } catch (e) {
    console.error(e)
  }
}

const handleSave = async () => {
  saving.value = true
  try {
    await updateMyCustomer(form)
    ElMessage.success('保存成功')
    editVisible.value = false
    fetchCustomer()
  } catch (e: any) {
    ElMessage.error(e.message || '保存失败')
  } finally {
    saving.value = false
  }
}

onMounted(fetchCustomer)
</script>

<style scoped lang="scss">
.page-container {
  .card-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }
}
</style>
