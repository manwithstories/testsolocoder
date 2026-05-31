<template>
  <div class="dashboard-timeslots">
    <div class="card-shadow p-20">
      <div class="flex-between mb-20">
        <h2 class="page-title">时段管理 - {{ exhibition?.title }}</h2>
        <div class="actions">
          <el-button @click="$router.push('/dashboard/exhibitions')">返回列表</el-button>
          <el-button type="primary" @click="showBatchDialog = true">
            <el-icon><Plus /></el-icon> 批量生成时段
          </el-button>
        </div>
      </div>

      <div class="filter-bar mb-20">
        <el-date-picker
          v-model="selectedDate"
          type="date"
          placeholder="选择日期"
          value-format="YYYY-MM-DD"
          @change="fetchSlots"
        />
      </div>

      <el-table :data="slots" v-loading="loading" border>
        <el-table-column prop="date" label="日期" width="120">
          <template #default="{ row }">
            {{ formatDate(row.date) }}
          </template>
        </el-table-column>
        <el-table-column prop="start_time" label="开始时间" width="100" />
        <el-table-column prop="end_time" label="结束时间" width="100" />
        <el-table-column prop="max_capacity" label="最大容量" width="100" />
        <el-table-column label="已预约/总量" width="150">
          <template #default="{ row }">
            <el-progress
              :percentage="Math.round((row.booked_count / row.max_capacity) * 100)"
              :status="row.booked_count >= row.max_capacity ? 'exception' : undefined"
            />
            <span style="margin-left: 10px;">{{ row.booked_count }}/{{ row.max_capacity }}</span>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.booked_count >= row.max_capacity ? 'danger' : 'success'" size="small">
              {{ row.booked_count >= row.max_capacity ? '已约满' : '可预约' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-button type="danger" size="small" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-dialog v-model="showBatchDialog" title="批量生成时段" width="500px">
      <el-form :model="batchForm" :rules="batchRules" ref="batchFormRef" label-width="120px">
        <el-form-item label="开始日期" prop="start_date">
          <el-date-picker v-model="batchForm.start_date" type="date" value-format="YYYY-MM-DD" style="width: 100%" />
        </el-form-item>
        <el-form-item label="结束日期" prop="end_date">
          <el-date-picker v-model="batchForm.end_date" type="date" value-format="YYYY-MM-DD" style="width: 100%" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="开始时间" prop="start_time">
              <el-time-picker v-model="batchForm.start_time" format="HH:mm" value-format="HH:mm" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="结束时间" prop="end_time">
              <el-time-picker v-model="batchForm.end_time" format="HH:mm" value-format="HH:mm" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="间隔(分钟)" prop="interval">
              <el-input-number v-model="batchForm.interval" :min="15" :step="15" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="最大人数" prop="max_capacity">
              <el-input-number v-model="batchForm.max_capacity" :min="1" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="showBatchDialog = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleBatchCreate">生成</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import * as exhibitionApi from '@/api/exhibition'
import type { Exhibition, TimeSlot } from '@/types'
import dayjs from 'dayjs'

const route = useRoute()
const exhibitionId = Number(route.params.id)

const loading = ref(false)
const submitting = ref(false)
const exhibition = ref<Exhibition | null>(null)
const slots = ref<TimeSlot[]>([])
const showBatchDialog = ref(false)
const selectedDate = ref(dayjs().format('YYYY-MM-DD'))
const batchFormRef = ref<FormInstance>()

const batchForm = reactive({
  start_date: '',
  end_date: '',
  start_time: '09:00',
  end_time: '17:00',
  interval: 30,
  max_capacity: 50
})

const batchRules: FormRules = {
  start_date: [{ required: true, message: '请选择开始日期', trigger: 'change' }],
  end_date: [{ required: true, message: '请选择结束日期', trigger: 'change' }],
  start_time: [{ required: true, message: '请选择开始时间', trigger: 'change' }],
  end_time: [{ required: true, message: '请选择结束时间', trigger: 'change' }],
  interval: [{ required: true, message: '请输入间隔时间', trigger: 'blur' }],
  max_capacity: [{ required: true, message: '请输入最大人数', trigger: 'blur' }]
}

const formatDate = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD')
}

const fetchExhibition = async () => {
  try {
    const res = await exhibitionApi.getExhibition(exhibitionId)
    exhibition.value = res.data
  } catch (e) {
    console.error(e)
  }
}

const fetchSlots = async () => {
  try {
    loading.value = true
    const res = await exhibitionApi.listTimeSlots(exhibitionId, selectedDate.value)
    slots.value = res.data
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const handleDelete = (row: TimeSlot) => {
  ElMessageBox.confirm('确定要删除该时段吗？', '提示', {
    type: 'warning'
  }).then(() => {
    ElMessage.success('删除成功')
    fetchSlots()
  }).catch(() => {})
}

const handleBatchCreate = async () => {
  if (!batchFormRef.value) return
  await batchFormRef.value.validate()
  try {
    submitting.value = true
    await exhibitionApi.batchCreateTimeSlots({
      exhibition_id: exhibitionId,
      ...batchForm
    })
    ElMessage.success('生成成功')
    showBatchDialog.value = false
    fetchSlots()
  } catch (e) {
    console.error(e)
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  fetchExhibition()
  fetchSlots()
})
</script>
