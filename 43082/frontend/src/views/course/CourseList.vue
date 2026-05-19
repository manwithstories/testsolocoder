<template>
  <div class="course-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span style="font-weight: 600">课程列表</span>
          <div class="header-actions">
            <el-input
              v-model="searchKeyword"
              placeholder="搜索课程名称"
              style="width: 200px; margin-right: 12px"
              clearable
              @clear="loadCourses"
              @keyup.enter="loadCourses"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
            <el-select v-model="filterType" placeholder="课程类型" style="width: 120px; margin-right: 12px" clearable @change="loadCourses">
              <el-option label="单次课" value="single" />
              <el-option label="周课" value="weekly" />
              <el-option label="月课" value="monthly" />
            </el-select>
            <el-button type="primary" @click="openAddDialog">
              <el-icon><Plus /></el-icon>
              添加课程
            </el-button>
          </div>
        </div>
      </template>

      <el-table :data="courses" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="课程名称" />
        <el-table-column prop="coach.name" label="教练" />
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="getCourseTagType(row.type)">
              {{ getCourseTypeName(row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="capacity" label="容量" width="80" />
        <el-table-column prop="duration" label="时长(分钟)" width="100" />
        <el-table-column prop="start_time" label="开始时间" width="120" />
        <el-table-column prop="location" label="地点" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : row.status === 2 ? 'danger' : 'info'">
              {{ getStatusName(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="viewDetail(row.id)">详情</el-button>
            <el-button type="warning" link @click="openEditDialog(row)">编辑</el-button>
            <el-button type="danger" link @click="deleteCourse(row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :page-sizes="[10, 20, 50]"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next, jumper"
        style="margin-top: 20px; justify-content: flex-end"
        @size-change="loadCourses"
        @current-change="loadCourses"
      />
    </el-card>

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑课程' : '添加课程'" width="600px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="课程名称" prop="name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="课程描述">
          <el-input v-model="form.description" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item label="教练" prop="coach_id">
          <el-select v-model="form.coach_id" style="width: 100%" placeholder="请选择教练">
            <el-option
              v-for="coach in coaches"
              :key="coach.id"
              :label="coach.name"
              :value="coach.id"
            />
          </el-select>
        </el-form-item>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="课程容量" prop="capacity">
              <el-input-number v-model="form.capacity" :min="1" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="课程时长" prop="duration">
              <el-input-number v-model="form.duration" :min="1" style="width: 100%" />
              <span style="color: #909399; font-size: 12px">分钟</span>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="课程类型" prop="type">
          <el-radio-group v-model="form.type">
            <el-radio value="single">单次课</el-radio>
            <el-radio value="weekly">周课</el-radio>
            <el-radio value="monthly">月课</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="form.type !== 'single'" label="上课周几">
          <el-checkbox-group v-model="selectedWeekdays">
            <el-checkbox label="1">周一</el-checkbox>
            <el-checkbox label="2">周二</el-checkbox>
            <el-checkbox label="3">周三</el-checkbox>
            <el-checkbox label="4">周四</el-checkbox>
            <el-checkbox label="5">周五</el-checkbox>
            <el-checkbox label="6">周六</el-checkbox>
            <el-checkbox label="7">周日</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="开始日期" prop="start_date">
              <el-date-picker
                v-model="form.start_date"
                type="date"
                placeholder="选择日期"
                style="width: 100%"
                value-format="YYYY-MM-DD"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item v-if="form.type !== 'single'" label="结束日期">
              <el-date-picker
                v-model="form.end_date"
                type="date"
                placeholder="选择日期"
                style="width: 100%"
                value-format="YYYY-MM-DD"
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="开始时间" prop="start_time">
          <el-time-picker
            v-model="form.start_time"
            format="HH:mm"
            value-format="HH:mm"
            placeholder="选择时间"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="上课地点">
          <el-input v-model="form.location" placeholder="例如：1号操房" />
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
import { reactive, ref, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { courseApi, coachApi } from '@/api/course'
import type { Course, Coach } from '@/types'

const router = useRouter()
const loading = ref(false)
const submitting = ref(false)
const courses = ref<Course[]>([])
const coaches = ref<Coach[]>([])
const searchKeyword = ref('')
const filterType = ref('')
const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref<FormInstance>()
const selectedWeekdays = ref<string[]>([])

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const form = reactive({
  id: 0,
  name: '',
  description: '',
  coach_id: 0,
  capacity: 20,
  duration: 60,
  type: 'single' as any,
  weekdays: '',
  start_date: '',
  end_date: '',
  start_time: '',
  location: ''
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入课程名称', trigger: 'blur' }],
  coach_id: [{ required: true, message: '请选择教练', trigger: 'change' }],
  capacity: [{ required: true, message: '请输入课程容量', trigger: 'blur' }],
  duration: [{ required: true, message: '请输入课程时长', trigger: 'blur' }],
  type: [{ required: true, message: '请选择课程类型', trigger: 'change' }],
  start_date: [{ required: true, message: '请选择开始日期', trigger: 'change' }],
  start_time: [{ required: true, message: '请选择开始时间', trigger: 'change' }]
}

const getCourseTypeName = (type: string) => {
  const map: Record<string, string> = { single: '单次课', weekly: '周课', monthly: '月课' }
  return map[type] || type
}

const getCourseTagType = (type: string) => {
  const map: Record<string, string> = { single: '', weekly: 'warning', monthly: 'success' }
  return map[type] || ''
}

const getStatusName = (status: number) => {
  const map: Record<number, string> = { 1: '正常', 2: '取消', 3: '结束' }
  return map[status] || '未知'
}

const loadCourses = async () => {
  try {
    loading.value = true
    const res = await courseApi.getList({
      page: pagination.page,
      page_size: pagination.pageSize,
      keyword: searchKeyword.value,
      type: filterType.value
    })
    courses.value = res.data
    pagination.total = res.pagination.total
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const loadCoaches = async () => {
  try {
    const res = await coachApi.getList({ page_size: 100 })
    coaches.value = res.data
  } catch (error) {
    console.error(error)
  }
}

const viewDetail = (id: number) => {
  router.push(`/courses/${id}`)
}

const openAddDialog = () => {
  isEdit.value = false
  Object.assign(form, {
    id: 0,
    name: '',
    description: '',
    coach_id: coaches.value[0]?.id || 0,
    capacity: 20,
    duration: 60,
    type: 'single',
    weekdays: '',
    start_date: '',
    end_date: '',
    start_time: '',
    location: ''
  })
  selectedWeekdays.value = []
  dialogVisible.value = true
}

const openEditDialog = (row: Course) => {
  isEdit.value = true
  Object.assign(form, { ...row })
  selectedWeekdays.value = row.weekdays ? row.weekdays.split(',') : []
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (form.type !== 'single' && selectedWeekdays.value.length === 0) {
    ElMessage.warning('请选择上课周几')
    return
  }
  
  form.weekdays = selectedWeekdays.value.join(',')
  
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    
    try {
      submitting.value = true
      if (isEdit.value) {
        const { id, ...data } = form
        await courseApi.update(id, data)
        ElMessage.success('编辑成功')
      } else {
        await courseApi.create(form as any)
        ElMessage.success('添加成功')
      }
      dialogVisible.value = false
      loadCourses()
    } catch (error) {
      console.error(error)
    } finally {
      submitting.value = false
    }
  })
}

const deleteCourse = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定要删除该课程吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await courseApi.delete(id)
    ElMessage.success('删除成功')
    loadCourses()
  } catch (error) {
    if (error !== 'cancel') {
      console.error(error)
    }
  }
}

watch(selectedWeekdays, (val) => {
  form.weekdays = val.join(',')
})

onMounted(() => {
  loadCourses()
  loadCoaches()
})
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
