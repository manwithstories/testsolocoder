<template>
  <div class="user-management">
    <div class="page-header">
      <h2>用户管理</h2>
    </div>
    
    <div class="filter-bar">
      <el-input 
        v-model="searchKeyword" 
        placeholder="搜索用户名" 
        clearable
        style="width: 200px;"
        @change="loadUsers"
      />
      <el-select 
        v-model="filterRole" 
        placeholder="筛选角色" 
        clearable
        style="width: 120px;"
        @change="loadUsers"
      >
        <el-option label="管理员" value="admin" />
        <el-option label="音乐人" value="artist" />
        <el-option label="厂牌" value="label" />
        <el-option label="乐迷" value="fan" />
      </el-select>
    </div>
    
    <el-table :data="users" v-loading="loading" style="width: 100%">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="username" label="用户名" />
      <el-table-column prop="nickname" label="昵称" />
      <el-table-column prop="email" label="邮箱" />
      <el-table-column prop="role" label="角色" width="100">
        <template #default="{ row }">
          <el-tag :type="getRoleTag(row.role)">{{ getRoleText(row.role) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'danger'">
            {{ row.status === 1 ? '正常' : '已禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="注册时间" width="180" />
      <el-table-column label="操作" width="150">
        <template #default="{ row }">
          <el-button 
            type="primary" 
            text 
            size="small"
            @click="editRole(row)"
          >
            修改角色
          </el-button>
          <el-button 
            :type="row.status === 1 ? 'danger' : 'success'" 
            text 
            size="small"
            @click="toggleStatus(row)"
          >
            {{ row.status === 1 ? '禁用' : '启用' }}
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
        @current-change="loadUsers"
      />
    </div>
    
    <el-dialog v-model="roleDialog.visible" title="修改角色" width="400px">
      <el-form :model="roleDialog.form" label-width="80px">
        <el-form-item label="角色">
          <el-select v-model="roleDialog.form.role">
            <el-option label="管理员" value="admin" />
            <el-option label="音乐人" value="artist" />
            <el-option label="厂牌" value="label" />
            <el-option label="乐迷" value="fan" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="roleDialog.visible = false">取消</el-button>
        <el-button type="primary" @click="saveRole">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { userApi } from '@/api/auth'

const loading = ref(false)
const users = ref<any[]>([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const searchKeyword = ref('')
const filterRole = ref('')

const roleDialog = reactive({
  visible: false,
  userId: 0,
  form: {
    role: ''
  }
})

onMounted(() => {
  loadUsers()
})

async function loadUsers() {
  loading.value = true
  try {
    const params: any = {
      page: page.value,
      page_size: pageSize.value
    }
    
    if (searchKeyword.value) {
      params.keyword = searchKeyword.value
    }
    if (filterRole.value) {
      params.role = filterRole.value
    }
    
    const res = await userApi.list(params)
    users.value = res.list
    total.value = res.total
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function editRole(user: any) {
  roleDialog.userId = user.id
  roleDialog.form.role = user.role
  roleDialog.visible = true
}

async function saveRole() {
  try {
    await userApi.updateRole(roleDialog.userId, roleDialog.form)
    ElMessage.success('角色修改成功')
    roleDialog.visible = false
    loadUsers()
  } catch (e) {
    console.error(e)
  }
}

async function toggleStatus(user: any) {
  try {
    const action = user.status === 1 ? '禁用' : '启用'
    await ElMessageBox.confirm(`确定要${action}该用户吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await userApi.updateStatus(user.id, { status: user.status === 1 ? 0 : 1 })
    ElMessage.success(`${action}成功`)
    loadUsers()
  } catch (e) {
    if (e !== 'cancel') {
      console.error(e)
    }
  }
}

function getRoleTag(role: string): 'primary' | 'success' | 'warning' | 'info' | 'danger' {
  const tags: Record<string, 'primary' | 'success' | 'warning' | 'info' | 'danger'> = {
    admin: 'danger',
    artist: 'success',
    label: 'warning',
    fan: 'info'
  }
  return tags[role] || 'info'
}

function getRoleText(role: string) {
  const texts: Record<string, string> = {
    admin: '管理员',
    artist: '音乐人',
    label: '厂牌',
    fan: '乐迷'
  }
  return texts[role] || role
}
</script>

<style scoped lang="scss">
.user-management {
  .filter-bar {
    display: flex;
    gap: 12px;
    margin-bottom: 16px;
  }
  
  .pagination {
    display: flex;
    justify-content: center;
    margin-top: 20px;
  }
}
</style>
