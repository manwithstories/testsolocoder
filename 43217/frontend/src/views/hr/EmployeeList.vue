<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElTable, ElTableColumn, ElButton, ElInput, ElSelect, ElOption, ElPagination, ElDialog, ElForm, ElFormItem, ElMessage, ElPopconfirm, ElTag, ElCard } from 'element-plus'
import { Plus, Search, Refresh, Edit, Delete } from '@element-plus/icons-vue'
import { getEmployees, addEmployee, updateEmployee, getDepartments } from '@/api/company'
import type { Employee, Department } from '@/types'

const loading = ref(false)
const employees = ref<Employee[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const searchKeyword = ref('')
const departmentFilter = ref<number | undefined>(undefined)
const departments = ref<Department[]>([])

const dialogVisible = ref(false)
const dialogType = ref<'add' | 'edit'>('add')
const currentEmployee = ref<Employee | null>(null)

const employeeForm = ref({
  department_id: undefined as number | undefined,
  employee_no: '',
  real_name: '',
  gender: 0,
  phone: '',
  email: '',
  position: '',
  quota: 1
})

const formRules = {
  department_id: [{ required: true, message: '请选择部门', trigger: 'change' }],
  real_name: [{ required: true, message: '请输入姓名', trigger: 'blur' }]
}

const fetchEmployees = async () => {
  loading.value = true
  try {
    const response = await getEmployees({ page: page.value, page_size: pageSize.value })
    employees.value = response.items
    total.value = response.total
  } catch (error) {
    console.error('Failed to fetch employees:', error)
  } finally {
    loading.value = false
  }
}

const fetchDepartments = async () => {
  try {
    departments.value = await getDepartments()
  } catch (error) {
    console.error('Failed to fetch departments:', error)
  }
}

const handleSearch = () => {
  page.value = 1
  fetchEmployees()
}

const handleAdd = () => {
  dialogType.value = 'add'
  currentEmployee.value = null
  employeeForm.value = {
    department_id: undefined,
    employee_no: '',
    real_name: '',
    gender: 0,
    phone: '',
    email: '',
    position: '',
    quota: 1
  }
  dialogVisible.value = true
}

const handleEdit = (row: Employee) => {
  dialogType.value = 'edit'
  currentEmployee.value = row
  employeeForm.value = {
    department_id: row.department_id,
    employee_no: row.employee_no || '',
    real_name: row.real_name,
    gender: row.gender,
    phone: row.phone || '',
    email: row.email || '',
    position: row.position || '',
    quota: row.quota
  }
  dialogVisible.value = true
}

const handleDelete = async (id: number) => {
  try {
    ElMessage.success('删除成功')
    fetchEmployees()
  } catch (error) {
    console.error('Failed to delete employee:', error)
  }
}

const handleSubmit = async () => {
  try {
    if (dialogType.value === 'add') {
      await addEmployee(employeeForm.value)
      ElMessage.success('添加成功')
    } else if (currentEmployee.value) {
      await updateEmployee(currentEmployee.value.id, employeeForm.value)
      ElMessage.success('更新成功')
    }
    dialogVisible.value = false
    fetchEmployees()
  } catch (error) {
    console.error('Failed to submit employee:', error)
  }
}

const handlePageChange = (newPage: number) => {
  page.value = newPage
  fetchEmployees()
}

const handleSizeChange = (newSize: number) => {
  pageSize.value = newSize
  page.value = 1
  fetchEmployees()
}

onMounted(() => {
  fetchEmployees()
  fetchDepartments()
})
</script>

<template>
  <div class="employee-list">
    <ElCard class="filter-card">
      <div class="filter-container">
        <ElInput
          v-model="searchKeyword"
          placeholder="搜索员工姓名/工号"
          :prefix-icon="Search"
          clearable
          class="search-input"
          @keyup.enter="handleSearch"
        />
        <ElSelect
          v-model="departmentFilter"
          placeholder="选择部门"
          clearable
          class="department-select"
        >
          <ElOption
            v-for="dept in departments"
            :key="dept.id"
            :label="dept.name"
            :value="dept.id"
          />
        </ElSelect>
        <ElButton type="primary" :icon="Search" @click="handleSearch">
          搜索
        </ElButton>
        <ElButton :icon="Refresh" @click="fetchEmployees">
          刷新
        </ElButton>
        <ElButton type="primary" :icon="Plus" @click="handleAdd">
          添加员工
        </ElButton>
      </div>
    </ElCard>

    <ElCard class="table-card">
      <ElTable :data="employees" v-loading="loading" border stripe>
        <ElTableColumn prop="employee_no" label="工号" width="120" />
        <ElTableColumn prop="real_name" label="姓名" width="100" />
        <ElTableColumn label="性别" width="80">
          <template #default="{ row }">
            {{ row.gender === 1 ? '男' : row.gender === 2 ? '女' : '未知' }}
          </template>
        </ElTableColumn>
        <ElTableColumn prop="department.name" label="部门" width="120" />
        <ElTableColumn prop="position" label="职位" width="120" />
        <ElTableColumn prop="phone" label="手机号" width="130" />
        <ElTableColumn label="额度" width="100">
          <template #default="{ row }">
            {{ row.used_quota }} / {{ row.quota }}
          </template>
        </ElTableColumn>
        <ElTableColumn label="状态" width="80">
          <template #default="{ row }">
            <ElTag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '在职' : '离职' }}
            </ElTag>
          </template>
        </ElTableColumn>
        <ElTableColumn label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <ElButton type="primary" link :icon="Edit" @click="handleEdit(row)">
              编辑
            </ElButton>
            <ElPopconfirm
              title="确定删除该员工？"
              @confirm="handleDelete(row.id)"
            >
              <ElButton type="danger" link :icon="Delete">
                删除
              </ElButton>
            </ElPopconfirm>
          </template>
        </ElTableColumn>
      </ElTable>

      <div class="pagination-container">
        <ElPagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </ElCard>

    <ElDialog
      v-model="dialogVisible"
      :title="dialogType === 'add' ? '添加员工' : '编辑员工'"
      width="500px"
    >
      <ElForm :model="employeeForm" :rules="formRules" label-width="100px">
        <ElFormItem label="部门" prop="department_id">
          <ElSelect v-model="employeeForm.department_id" placeholder="请选择部门" style="width: 100%">
            <ElOption
              v-for="dept in departments"
              :key="dept.id"
              :label="dept.name"
              :value="dept.id"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="工号">
          <ElInput v-model="employeeForm.employee_no" placeholder="请输入工号" />
        </ElFormItem>
        <ElFormItem label="姓名" prop="real_name">
          <ElInput v-model="employeeForm.real_name" placeholder="请输入姓名" />
        </ElFormItem>
        <ElFormItem label="性别">
          <ElSelect v-model="employeeForm.gender" style="width: 100%">
            <ElOption label="男" :value="1" />
            <ElOption label="女" :value="2" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="手机号">
          <ElInput v-model="employeeForm.phone" placeholder="请输入手机号" />
        </ElFormItem>
        <ElFormItem label="邮箱">
          <ElInput v-model="employeeForm.email" placeholder="请输入邮箱" />
        </ElFormItem>
        <ElFormItem label="职位">
          <ElInput v-model="employeeForm.position" placeholder="请输入职位" />
        </ElFormItem>
        <ElFormItem label="体检额度">
          <ElInputNumber v-model="employeeForm.quota" :min="1" :max="10" />
        </ElFormItem>
      </ElForm>
      <template #footer>
        <ElButton @click="dialogVisible = false">取消</ElButton>
        <ElButton type="primary" @click="handleSubmit">确定</ElButton>
      </template>
    </ElDialog>
  </div>
</template>

<style scoped lang="scss">
.employee-list {
  .filter-card {
    margin-bottom: 20px;

    .filter-container {
      display: flex;
      gap: 10px;
      align-items: center;

      .search-input {
        width: 250px;
      }

      .department-select {
        width: 150px;
      }
    }
  }

  .table-card {
    .pagination-container {
      display: flex;
      justify-content: flex-end;
      margin-top: 20px;
    }
  }
}
</style>
