<template>
  <div class="page-container">
    <div class="page-header">
      <h2 class="page-title">用户管理</h2>
      <div class="page-actions">
        <el-button type="primary" :icon="Plus" @click="handleAdd">添加用户</el-button>
      </div>
    </div>

    <div class="filter-bar">
      <div class="filter-form">
        <div class="filter-item">
          <label class="filter-label">关键词</label>
          <el-input
            v-model="searchForm.keyword"
            placeholder="用户名/姓名/电话"
            clearable
            style="width: 200px"
            @keyup.enter="handleSearch"
          />
        </div>
        <div class="filter-item">
          <label class="filter-label">角色</label>
          <el-select
            v-model="searchForm.role"
            placeholder="全部角色"
            clearable
            style="width: 150px"
          >
            <el-option
              v-for="role in roleOptions"
              :key="role.value"
              :label="role.label"
              :value="role.value"
            />
          </el-select>
        </div>
        <div class="filter-item">
          <label class="filter-label">状态</label>
          <el-select
            v-model="searchForm.status"
            placeholder="全部状态"
            clearable
            style="width: 150px"
          >
            <el-option label="启用" :value="true" />
            <el-option label="禁用" :value="false" />
          </el-select>
        </div>
        <div class="filter-actions">
          <el-button type="primary" :icon="Search" @click="handleSearch">查询</el-button>
          <el-button :icon="Refresh" @click="handleReset">重置</el-button>
        </div>
      </div>
    </div>

    <div class="common-card">
      <div class="card-body">
        <el-table
          :data="tableData"
          v-loading="loading"
          border
          stripe
          class="common-table"
        >
          <el-table-column prop="id" label="ID" width="80" align="center" />
          <el-table-column prop="username" label="用户名" width="120" />
          <el-table-column prop="name" label="真实姓名" width="100" />
          <el-table-column prop="phone" label="电话" width="130" />
          <el-table-column prop="email" label="邮箱" min-width="180" show-overflow-tooltip />
          <el-table-column label="角色" width="100" align="center">
            <template #default="{ row }">
              <el-tag :type="getRoleTagType(row.role)">{{ getRoleLabel(row.role) }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="状态" width="80" align="center">
            <template #default="{ row }">
              <el-switch
                v-model="row.status"
                :active-value="true"
                :inactive-value="false"
                @change="(val: boolean) => handleStatusChange(row, val)"
              />
            </template>
          </el-table-column>
          <el-table-column prop="createdAt" label="创建时间" width="180" align="center">
            <template #default="{ row }">
              {{ formatDate(row.createdAt) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="200" fixed="right" align="center">
            <template #default="{ row }">
              <el-button type="primary" link :icon="Edit" @click="handleEdit(row)">编辑</el-button>
              <el-button type="warning" link :icon="Key" @click="handleResetPassword(row)">重置密码</el-button>
              <el-button type="danger" link :icon="Delete" @click="handleDelete(row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>

        <div class="pagination-wrapper">
          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.pageSize"
            :page-sizes="[10, 20, 50, 100]"
            :total="pagination.total"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          />
        </div>
      </div>
    </div>

    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑用户' : '添加用户'"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="100px"
      >
        <el-form-item label="用户名" prop="username">
          <el-input
            v-model="formData.username"
            placeholder="请输入用户名"
            :disabled="isEdit"
          />
        </el-form-item>
        <el-form-item label="密码" prop="password" v-if="!isEdit">
          <el-input
            v-model="formData.password"
            type="password"
            placeholder="请输入密码"
            show-password
          />
        </el-form-item>
        <el-form-item label="真实姓名" prop="name">
          <el-input v-model="formData.name" placeholder="请输入真实姓名" />
        </el-form-item>
        <el-form-item label="电话" prop="phone">
          <el-input v-model="formData.phone" placeholder="请输入电话号码" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="formData.email" placeholder="请输入邮箱地址" />
        </el-form-item>
        <el-form-item label="角色" prop="role">
          <el-select v-model="formData.role" placeholder="请选择角色" style="width: 100%">
            <el-option
              v-for="role in roleOptions"
              :key="role.value"
              :label="role.label"
              :value="role.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-switch
            v-model="formData.status"
            :active-value="true"
            :inactive-value="false"
            active-text="启用"
            inactive-text="禁用"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" :loading="submitting" @click="handleSubmit">确定</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Plus, Search, Refresh, Edit, Key, Delete } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import { getUserList, createUser, updateUser, deleteUser, updateUserStatus, resetPassword } from '@/api/user'
import { UserRole, type User } from '@/types'

const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref<FormInstance>()

const searchForm = reactive({
  keyword: '',
  role: null as UserRole | null,
  status: null as boolean | null
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const tableData = ref<User[]>([])

const roleOptions = [
  { label: '管理员', value: UserRole.ADMIN },
  { label: '经理', value: UserRole.MANAGER },
  { label: '前台', value: UserRole.RECEPTIONIST },
  { label: '员工', value: UserRole.STAFF }
]

const formData = reactive<Partial<User>>({
  username: '',
  password: '',
  name: '',
  phone: '',
  email: '',
  role: UserRole.STAFF,
  status: true
})

const formRules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度在 6 到 20 个字符', trigger: 'blur' }
  ],
  name: [
    { required: true, message: '请输入真实姓名', trigger: 'blur' }
  ],
  phone: [
    { required: true, message: '请输入电话号码', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码', trigger: 'blur' }
  ],
  email: [
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ],
  role: [
    { required: true, message: '请选择角色', trigger: 'change' }
  ]
}

const getRoleLabel = (role: UserRole) => {
  const map: Record<UserRole, string> = {
    [UserRole.ADMIN]: '管理员',
    [UserRole.MANAGER]: '经理',
    [UserRole.RECEPTIONIST]: '前台',
    [UserRole.STAFF]: '员工'
  }
  return map[role] || role
}

const getRoleTagType = (role: UserRole) => {
  const map: Record<UserRole, string> = {
    [UserRole.ADMIN]: 'danger',
    [UserRole.MANAGER]: 'warning',
    [UserRole.RECEPTIONIST]: 'primary',
    [UserRole.STAFF]: 'info'
  }
  return map[role] || 'info'
}

const formatDate = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

const generateMockData = (): User[] => {
  const names = ['张三', '李四', '王五', '赵六', '钱七', '孙八', '周九', '吴十', '郑十一', '王十二']
  const roles = [UserRole.ADMIN, UserRole.MANAGER, UserRole.RECEPTIONIST, UserRole.STAFF, UserRole.STAFF]
  
  const users: User[] = []
  for (let i = 1; i <= 25; i++) {
    const nameIndex = (i - 1) % names.length
    users.push({
      id: i,
      username: `user${i}`,
      name: names[nameIndex],
      phone: `138${String(10000000 + i).slice(-8)}`,
      email: `user${i}@hotel.com`,
      role: roles[i % roles.length],
      status: i % 5 !== 0,
      createdAt: dayjs().subtract(i, 'day').format('YYYY-MM-DD HH:mm:ss'),
      updatedAt: dayjs().subtract(i, 'day').format('YYYY-MM-DD HH:mm:ss')
    })
  }
  return users
}

const fetchData = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      pageSize: pagination.pageSize,
      keyword: searchForm.keyword || undefined,
      role: searchForm.role || undefined,
      status: searchForm.status !== null ? searchForm.status : undefined
    }

    const data = await getUserList(params) as any
    
    if (data && data.list) {
      tableData.value = data.list
      pagination.total = data.total
    } else {
      const allData = generateMockData()
      const filtered = allData.filter(item => {
        if (searchForm.keyword) {
          const keyword = searchForm.keyword.toLowerCase()
          if (!item.username.toLowerCase().includes(keyword) &&
              !item.name.toLowerCase().includes(keyword) &&
              !item.phone.includes(keyword)) {
            return false
          }
        }
        if (searchForm.role && item.role !== searchForm.role) {
          return false
        }
        if (searchForm.status !== null && item.status !== searchForm.status) {
          return false
        }
        return true
      })
      
      const start = (pagination.page - 1) * pagination.pageSize
      const end = start + pagination.pageSize
      tableData.value = filtered.slice(start, end)
      pagination.total = filtered.length
    }
  } catch (error) {
    console.error('Failed to fetch user list:', error)
    const allData = generateMockData()
    const start = (pagination.page - 1) * pagination.pageSize
    const end = start + pagination.pageSize
    tableData.value = allData.slice(start, end)
    pagination.total = allData.length
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  fetchData()
}

const handleReset = () => {
  searchForm.keyword = ''
  searchForm.role = null
  searchForm.status = null
  pagination.page = 1
  fetchData()
}

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.page = 1
  fetchData()
}

const handleCurrentChange = (page: number) => {
  pagination.page = page
  fetchData()
}

const handleAdd = () => {
  isEdit.value = false
  Object.assign(formData, {
    username: '',
    password: '',
    name: '',
    phone: '',
    email: '',
    role: UserRole.STAFF,
    status: true
  })
  dialogVisible.value = true
}

const handleEdit = (row: User) => {
  isEdit.value = true
  Object.assign(formData, {
    id: row.id,
    username: row.username,
    name: row.name,
    phone: row.phone,
    email: row.email,
    role: row.role,
    status: row.status
  })
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitting.value = true
    try {
      if (isEdit.value && formData.id) {
        await updateUser(formData.id, {
          name: formData.name,
          phone: formData.phone,
          email: formData.email,
          role: formData.role,
          status: formData.status
        })
        ElMessage.success('更新成功')
      } else {
        await createUser({
          username: formData.username!,
          password: formData.password!,
          name: formData.name!,
          phone: formData.phone!,
          email: formData.email,
          role: formData.role!,
          status: formData.status!
        })
        ElMessage.success('创建成功')
      }

      dialogVisible.value = false
      fetchData()
    } catch (error) {
      console.error('Failed to submit user:', error)
      ElMessage.error(isEdit.value ? '更新失败' : '创建失败')
    } finally {
      submitting.value = false
    }
  })
}

const handleStatusChange = async (row: User, status: boolean) => {
  try {
    await updateUserStatus(row.id, status)
    ElMessage.success(`已${status ? '启用' : '禁用'}用户`)
  } catch (error) {
    console.error('Failed to update user status:', error)
    row.status = !status
    ElMessage.error('操作失败')
  }
}

const handleResetPassword = async (row: User) => {
  ElMessageBox.confirm(
    `确定要重置用户 "${row.username}" 的密码吗？`,
    '重置密码',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      await resetPassword(row.id)
      ElMessage.success('密码已重置为默认密码 123456')
    } catch (error) {
      console.error('Failed to reset password:', error)
      ElMessage.error('重置失败')
    }
  }).catch(() => {})
}

const handleDelete = async (row: User) => {
  ElMessageBox.confirm(
    `确定要删除用户 "${row.username}" 吗？此操作不可恢复。`,
    '删除确认',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'error'
    }
  ).then(async () => {
    try {
      await deleteUser(row.id)
      ElMessage.success('删除成功')
      fetchData()
    } catch (error) {
      console.error('Failed to delete user:', error)
      ElMessage.error('删除失败')
    }
  }).catch(() => {})
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped lang="scss">
</style>
