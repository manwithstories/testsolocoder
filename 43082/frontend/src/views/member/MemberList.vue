<template>
  <div class="member-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span style="font-weight: 600">会员列表</span>
          <div class="header-actions">
            <el-input
              v-model="searchKeyword"
              placeholder="搜索姓名/手机号"
              style="width: 240px; margin-right: 12px"
              clearable
              @clear="loadMembers"
              @keyup.enter="loadMembers"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
            <el-button type="primary" @click="openAddDialog">
              <el-icon><Plus /></el-icon>
              添加会员
            </el-button>
          </div>
        </div>
      </template>

      <el-table :data="members" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="姓名" />
        <el-table-column prop="phone" label="手机号" />
        <el-table-column prop="email" label="邮箱" />
        <el-table-column prop="gender" label="性别" width="80" />
        <el-table-column label="会员卡">
          <template #default="{ row }">
            <el-tag v-if="row.membership" :type="getMembershipTagType(row.membership.type)">
              {{ getMembershipTypeName(row.membership.type) }}
            </el-tag>
            <span v-else style="color: #909399">无</span>
          </template>
        </el-table-column>
        <el-table-column label="到期时间">
          <template #default="{ row }">
            <span v-if="row.membership">{{ row.membership.end_date }}</span>
            <span v-else style="color: #909399">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">
              {{ row.status === 1 ? '正常' : '冻结' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="240" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="viewDetail(row.id)">详情</el-button>
            <el-button type="warning" link @click="openEditDialog(row)">编辑</el-button>
            <el-button :type="row.status === 1 ? 'info' : 'success'" link @click="toggleStatus(row)">
              {{ row.status === 1 ? '冻结' : '解冻' }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next, jumper"
        style="margin-top: 20px; justify-content: flex-end"
        @size-change="loadMembers"
        @current-change="loadMembers"
      />
    </el-card>

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑会员' : '添加会员'" width="500px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="80px">
        <el-form-item label="姓名" prop="name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="form.phone" maxlength="11" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="form.email" />
        </el-form-item>
        <el-form-item label="性别">
          <el-radio-group v-model="form.gender">
            <el-radio value="男">男</el-radio>
            <el-radio value="女">女</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="!isEdit" label="密码" prop="password">
          <el-input v-model="form.password" type="password" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { memberApi } from '@/api/member'
import type { Member } from '@/types'

const router = useRouter()
const loading = ref(false)
const submitting = ref(false)
const members = ref<Member[]>([])
const searchKeyword = ref('')
const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref<FormInstance>()

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const form = reactive({
  id: 0,
  name: '',
  phone: '',
  email: '',
  gender: '男',
  password: ''
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  phone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { len: 11, message: '手机号格式不正确', trigger: 'blur' }
  ],
  email: [{ type: 'email', message: '邮箱格式不正确', trigger: 'blur' }],
  password: [{ min: 6, message: '密码长度不能少于6位', trigger: 'blur' }]
}

const getMembershipTypeName = (type: string) => {
  const map: Record<string, string> = { monthly: '月卡', quarter: '季卡', yearly: '年卡' }
  return map[type] || type
}

const getMembershipTagType = (type: string) => {
  const map: Record<string, string> = { monthly: '', quarter: 'warning', yearly: 'success' }
  return map[type] || ''
}

const loadMembers = async () => {
  try {
    loading.value = true
    const res = await memberApi.getList({
      page: pagination.page,
      page_size: pagination.pageSize,
      keyword: searchKeyword.value
    })
    members.value = res.data
    pagination.total = res.pagination.total
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const viewDetail = (id: number) => {
  router.push(`/members/${id}`)
}

const openAddDialog = () => {
  isEdit.value = false
  Object.assign(form, { id: 0, name: '', phone: '', email: '', gender: '男', password: '' })
  dialogVisible.value = true
}

const openEditDialog = (row: Member) => {
  isEdit.value = true
  Object.assign(form, { ...row, password: '' })
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    
    try {
      submitting.value = true
      if (isEdit.value) {
        const { id, ...data } = form
        await memberApi.update(id, data)
        ElMessage.success('编辑成功')
      } else {
        await memberApi.register(form)
        ElMessage.success('添加成功')
      }
      dialogVisible.value = false
      loadMembers()
    } catch (error) {
      console.error(error)
    } finally {
      submitting.value = false
    }
  })
}

const toggleStatus = async (row: Member) => {
  const newStatus = row.status === 1 ? 2 : 1
  const action = newStatus === 1 ? '解冻' : '冻结'
  
  try {
    await ElMessageBox.confirm(`确定要${action}该会员吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await memberApi.updateStatus(row.id, newStatus)
    ElMessage.success(`${action}成功`)
    loadMembers()
  } catch (error) {
    if (error !== 'cancel') {
      console.error(error)
    }
  }
}

onMounted(loadMembers)
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  align-items: center;
}
</style>
