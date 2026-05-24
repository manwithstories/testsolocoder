<template>
  <div class="agent-management">
    <div class="page-header flex-between">
      <h2 class="page-title">代办专员管理</h2>
      <el-button type="primary" @click="showCreateDialog = true">
        <el-icon><Plus /></el-icon>
        添加专员
      </el-button>
    </div>

    <div class="filter-bar">
      <el-form :inline="true" :model="filterForm" @submit.prevent>
        <el-form-item label="状态">
          <el-select v-model="filterForm.status" placeholder="全部" clearable style="width: 160px">
            <el-option label="全部" value="" />
            <el-option label="可用" value="available" />
            <el-option label="忙碌" value="busy" />
            <el-option label="休息" value="rest" />
          </el-select>
        </el-form-item>
        <el-form-item label="关键词">
          <el-input
            v-model="filterForm.keyword"
            placeholder="姓名/工号/手机号"
            clearable
            style="width: 200px"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchData">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button @click="resetFilter">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="table-wrapper">
      <el-table :data="agents" style="width: 100%" v-loading="loading">
        <el-table-column prop="agentProfile.employeeNo" label="工号" width="140">
          <template #default="{ row }">
            {{ row.agentProfile?.employeeNo || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="realName" label="姓名" width="120" />
        <el-table-column prop="phone" label="手机号" width="140" />
        <el-table-column prop="email" label="邮箱" width="180" />
        <el-table-column prop="agentProfile.specialtyTags" label="专业领域" width="150">
          <template #default="{ row }">
            <el-tag v-if="row.agentProfile?.specialtyTags">{{ row.agentProfile.specialtyTags }}</el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column label="工作负荷" width="150">
          <template #default="{ row }">
            {{ row.agentProfile?.currentApps || 0 }} / {{ row.agentProfile?.maxApplications || 0 }}
          </template>
        </el-table-column>
        <el-table-column prop="agentProfile.status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.agentProfile?.status)">
              {{ getStatusText(row.agentProfile?.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="editAgent(row)">编辑</el-button>
            <el-button type="danger" link @click="deleteAgent(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper mt-16">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :total="pagination.total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="fetchData"
          @current-change="fetchData"
        />
      </div>
    </div>

    <el-dialog v-model="showCreateDialog" title="添加代办专员" width="600px">
      <el-form ref="createFormRef" :model="createForm" :rules="createRules" label-width="100px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="createForm.username" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="createForm.password" type="password" show-password />
        </el-form-item>
        <el-form-item label="真实姓名" prop="realName">
          <el-input v-model="createForm.realName" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="createForm.phone" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="createForm.email" />
        </el-form-item>
        <el-form-item label="工号" prop="employeeNo">
          <el-input v-model="createForm.employeeNo" />
        </el-form-item>
        <el-form-item label="专业领域" prop="specialtyTags">
          <el-input v-model="createForm.specialtyTags" placeholder="如：工商注册、税务登记" />
        </el-form-item>
        <el-form-item label="最大处理数" prop="maxApps">
          <el-input-number v-model="createForm.maxApps" :min="1" :max="20" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" :loading="creating" @click="handleCreateAgent">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Plus, Search } from '@element-plus/icons-vue'
import { agentApi } from '@/api/agent'
import { User } from '@/types'

const loading = ref(false)
const agents = ref<User[]>([])
const showCreateDialog = ref(false)
const creating = ref(false)
const createFormRef = ref<FormInstance>()

const filterForm = reactive({
  status: '',
  keyword: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const createForm = reactive({
  username: '',
  password: '',
  realName: '',
  phone: '',
  email: '',
  employeeNo: '',
  specialtyTags: '',
  maxApps: 5
})

const createRules: FormRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
  realName: [{ required: true, message: '请输入真实姓名', trigger: 'blur' }],
  phone: [{ required: true, message: '请输入手机号', trigger: 'blur' }],
  employeeNo: [{ required: true, message: '请输入工号', trigger: 'blur' }],
  specialtyTags: [{ required: true, message: '请输入专业领域', trigger: 'blur' }]
}

const fetchData = async () => {
  loading.value = true
  try {
    const res = await agentApi.getList({
      page: pagination.page,
      pageSize: pagination.pageSize,
      status: filterForm.status || undefined,
      keyword: filterForm.keyword || undefined
    })
    agents.value = res?.list || []
    pagination.total = res?.total || 0
  } catch (error) {
    console.error('获取专员列表失败:', error)
  } finally {
    loading.value = false
  }
}

const resetFilter = () => {
  filterForm.status = ''
  filterForm.keyword = ''
  pagination.page = 1
  fetchData()
}

const getStatusType = (status?: string) => {
  const map: Record<string, string> = {
    available: 'success',
    busy: 'warning',
    rest: 'info'
  }
  return map[status || ''] || 'info'
}

const getStatusText = (status?: string) => {
  const map: Record<string, string> = {
    available: '可用',
    busy: '忙碌',
    rest: '休息'
  }
  return map[status || ''] || status || '-'
}

const editAgent = (row: User) => {
  ElMessage.info('编辑功能开发中')
}

const deleteAgent = async (row: User) => {
  try {
    await ElMessageBox.confirm(`确认删除专员 ${row.realName}？`, '提示', {
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await agentApi.delete(row.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  }
}

const handleCreateAgent = async () => {
  if (!createFormRef.value) return

  await createFormRef.value.validate(async (valid) => {
    if (valid) {
      creating.value = true
      try {
        await agentApi.create({
          username: createForm.username,
          password: createForm.password,
          realName: createForm.realName,
          phone: createForm.phone,
          email: createForm.email,
          employeeNo: createForm.employeeNo,
          specialtyTags: createForm.specialtyTags,
          maxApps: createForm.maxApps
        })
        ElMessage.success('添加成功')
        showCreateDialog.value = false
        fetchData()
      } catch (error: any) {
        ElMessage.error(error.message || '添加失败')
      } finally {
        creating.value = false
      }
    }
  })
}

onMounted(fetchData)
</script>

<style scoped>
.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
}
</style>
