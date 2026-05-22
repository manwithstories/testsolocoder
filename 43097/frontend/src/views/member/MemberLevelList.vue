<template>
  <div class="member-level-list">
    <div class="header">
      <h2>会员等级管理</h2>
      <el-button type="primary" @click="openLevelDialog">添加等级</el-button>
    </div>

    <el-card class="table-card">
      <el-table :data="tableData" v-loading="loading" border stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="等级名称" width="150" />
        <el-table-column prop="discountRate" label="折扣率" width="120">
          <template #default="{ row }">
            {{ (row.discountRate * 100).toFixed(0) }}%
          </template>
        </el-table-column>
        <el-table-column prop="pointsRate" label="积分倍率" width="120">
          <template #default="{ row }">
            {{ row.pointsRate }}x
          </template>
        </el-table-column>
        <el-table-column label="积分范围" width="200">
          <template #default="{ row }">
            {{ row.minPoints }} - {{ row.maxPoints > 0 ? row.maxPoints : '无上限' }}
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDateTime(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="primary" link @click="openEditDialog(row)">编辑</el-button>
            <el-button size="small" type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
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
    </el-card>

    <el-dialog v-model="levelDialogVisible" :title="isEdit ? '编辑等级' : '添加等级'" width="500px">
      <el-form :model="levelForm" :rules="levelRules" ref="levelFormRef" label-width="100px">
        <el-form-item label="等级名称" prop="name">
          <el-input v-model="levelForm.name" placeholder="请输入等级名称" />
        </el-form-item>
        <el-form-item label="折扣率" prop="discountRate">
          <el-slider
            v-model="levelForm.discountRate"
            :min="0.1"
            :max="1"
            :step="0.05"
            :marks="marks"
            :format-tooltip="formatDiscount"
          />
          <div class="slider-value">当前折扣：{{ (levelForm.discountRate * 100).toFixed(0) }}%</div>
        </el-form-item>
        <el-form-item label="积分倍率" prop="pointsRate">
          <el-input-number v-model="levelForm.pointsRate" :min="0.5" :max="10" :step="0.5" />
          <span class="form-tip">消费1元获得的积分倍数</span>
        </el-form-item>
        <el-form-item label="最低积分" prop="minPoints">
          <el-input-number v-model="levelForm.minPoints" :min="0" />
        </el-form-item>
        <el-form-item label="最高积分" prop="maxPoints">
          <el-input-number v-model="levelForm.maxPoints" :min="0" />
          <span class="form-tip">0表示无上限</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="levelDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleLevelSubmit">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import {
  getMemberLevelList,
  createMemberLevel,
  updateMemberLevel,
  deleteMemberLevel,
  type MemberLevelData
} from '@/api/member'

const loading = ref(false)
const submitting = ref(false)
const tableData = ref<MemberLevelData[]>([])
const levelFormRef = ref<FormInstance>()

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const levelDialogVisible = ref(false)
const isEdit = ref(false)

const levelForm = reactive<MemberLevelData>({
  id: undefined,
  name: '',
  discountRate: 1,
  pointsRate: 1,
  minPoints: 0,
  maxPoints: 0
})

const levelRules: FormRules = {
  name: [{ required: true, message: '请输入等级名称', trigger: 'blur' }],
  discountRate: [{ required: true, message: '请输入折扣率', trigger: 'change' }],
  pointsRate: [{ required: true, message: '请输入积分倍率', trigger: 'change' }],
  minPoints: [{ required: true, message: '请输入最低积分', trigger: 'change' }],
  maxPoints: [{ required: true, message: '请输入最高积分', trigger: 'change' }]
}

const marks = {
  0.1: '10%',
  0.5: '50%',
  0.7: '70%',
  0.9: '90%',
  1: '100%'
}

const fetchList = async () => {
  loading.value = true
  try {
    const res = await getMemberLevelList({
      page: pagination.page,
      pageSize: pagination.pageSize
    })
    tableData.value = res.list
    pagination.total = res.total
  } finally {
    loading.value = false
  }
}

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  fetchList()
}

const handleCurrentChange = (page: number) => {
  pagination.page = page
  fetchList()
}

const formatDateTime = (date?: string) => {
  return date ? new Date(date).toLocaleString('zh-CN') : '-'
}

const formatDiscount = (value: number) => {
  return `${(value * 100).toFixed(0)}%`
}

const openLevelDialog = () => {
  isEdit.value = false
  Object.assign(levelForm, {
    id: undefined,
    name: '',
    discountRate: 1,
    pointsRate: 1,
    minPoints: 0,
    maxPoints: 0
  })
  levelDialogVisible.value = true
}

const openEditDialog = (row: MemberLevelData) => {
  isEdit.value = true
  Object.assign(levelForm, {
    id: row.id,
    name: row.name,
    discountRate: row.discountRate,
    pointsRate: row.pointsRate,
    minPoints: row.minPoints,
    maxPoints: row.maxPoints
  })
  levelDialogVisible.value = true
}

const handleLevelSubmit = async () => {
  if (!levelFormRef.value) return
  await levelFormRef.value.validate(async (valid) => {
    if (!valid) return
    submitting.value = true
    try {
      if (isEdit.value && levelForm.id) {
        await updateMemberLevel(levelForm.id, {
          name: levelForm.name,
          discountRate: levelForm.discountRate,
          pointsRate: levelForm.pointsRate,
          minPoints: levelForm.minPoints,
          maxPoints: levelForm.maxPoints
        })
        ElMessage.success('编辑成功')
      } else {
        await createMemberLevel({
          name: levelForm.name,
          discountRate: levelForm.discountRate,
          pointsRate: levelForm.pointsRate,
          minPoints: levelForm.minPoints,
          maxPoints: levelForm.maxPoints
        })
        ElMessage.success('添加成功')
      }
      levelDialogVisible.value = false
      fetchList()
    } catch (e: any) {
      ElMessage.error(e.message || '操作失败')
    } finally {
      submitting.value = false
    }
  })
}

const handleDelete = (row: MemberLevelData) => {
  ElMessageBox.confirm('确定要删除该会员等级吗？删除后关联的会员将受到影响。', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await deleteMemberLevel(row.id!)
      ElMessage.success('删除成功')
      fetchList()
    } catch (e: any) {
      ElMessage.error(e.message || '删除失败')
    }
  })
}

onMounted(() => {
  fetchList()
})
</script>

<style scoped>
.member-level-list {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.header h2 {
  margin: 0;
}

.table-card {
  margin-bottom: 20px;
}

.pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

.slider-value {
  margin-top: 8px;
  font-size: 14px;
  color: #606266;
}

.form-tip {
  margin-left: 10px;
  font-size: 12px;
  color: #909399;
}
</style>
