<template>
  <div class="page-container">
    <div class="page-header">
      <h2>婚礼管理</h2>
      <el-button type="primary" :icon="Plus" @click="showCreateDialog = true">
        新建婚礼
      </el-button>
    </div>

    <el-card shadow="never">
      <div class="filter-bar">
        <el-input
          v-model="searchQuery"
          placeholder="搜索婚礼名称..."
          :prefix-icon="Search"
          clearable
          style="width: 240px"
          @clear="fetchWeddings"
          @keyup.enter="fetchWeddings"
        />
        <el-select v-model="statusFilter" placeholder="状态" clearable style="width: 140px" @change="fetchWeddings">
          <el-option label="筹备中" value="planning" />
          <el-option label="已确认" value="confirmed" />
          <el-option label="已完成" value="completed" />
          <el-option label="已取消" value="cancelled" />
        </el-select>
      </div>

      <el-table :data="weddings" v-loading="loading" stripe>
        <el-table-column prop="title" label="婚礼名称" min-width="180" />
        <el-table-column label="新人" min-width="160">
          <template #default="{ row }">
            {{ row.groom_name }} & {{ row.bride_name }}
          </template>
        </el-table-column>
        <el-table-column label="婚礼日期" width="120">
          <template #default="{ row }">
            {{ formatDate(row.wedding_date) }}
          </template>
        </el-table-column>
        <el-table-column prop="venue" label="场地" min-width="140" show-overflow-tooltip />
        <el-table-column label="预算" width="120">
          <template #default="{ row }">
            ¥{{ formatNumber(row.budget) }}
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="statusType(row.status)">{{ statusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="viewDetail(row)">详情</el-button>
            <el-button type="primary" link @click="editWedding(row)">编辑</el-button>
            <el-button type="danger" link @click="deleteWedding(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="fetchWeddings"
        @current-change="fetchWeddings"
        class="pagination"
      />
    </el-card>

    <el-dialog
      v-model="showCreateDialog"
      :title="editingWedding ? '编辑婚礼' : '新建婚礼'"
      width="600px"
      @close="resetForm"
    >
      <el-form ref="weddingForm" :model="weddingForm" :rules="weddingRules" label-width="100px">
        <el-form-item label="婚礼名称" prop="title">
          <el-input v-model="weddingForm.title" placeholder="请输入婚礼名称" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="新郎" prop="groom_name">
              <el-input v-model="weddingForm.groom_name" placeholder="新郎姓名" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="新娘" prop="bride_name">
              <el-input v-model="weddingForm.bride_name" placeholder="新娘姓名" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="婚礼日期" prop="wedding_date">
          <el-date-picker
            v-model="weddingForm.wedding_date"
            type="date"
            placeholder="选择婚礼日期"
            style="width: 100%"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item label="预算" prop="budget">
          <el-input-number v-model="weddingForm.budget" :min="0" :precision="2" style="width: 100%" />
        </el-form-item>
        <el-form-item label="婚礼风格">
          <el-input v-model="weddingForm.style" placeholder="如：浪漫、简约、奢华" />
        </el-form-item>
        <el-form-item label="场地">
          <el-input v-model="weddingForm.venue" placeholder="场地名称" />
        </el-form-item>
        <el-form-item label="场地地址">
          <el-input v-model="weddingForm.venue_address" placeholder="场地详细地址" />
        </el-form-item>
        <el-form-item label="预计人数">
          <el-input-number v-model="weddingForm.guest_count" :min="0" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="weddingForm.description" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="saveWedding">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useWeddingStore } from '@/store/wedding'
import { weddingApi } from '@/api/wedding'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Plus, Search } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import type { Wedding } from '@/types'

const router = useRouter()
const weddingStore = useWeddingStore()

const loading = ref(false)
const saving = ref(false)
const weddings = ref<Wedding[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const searchQuery = ref('')
const statusFilter = ref('')
const showCreateDialog = ref(false)
const editingWedding = ref<Wedding | null>(null)

const weddingForm = reactive({
  title: '',
  groom_name: '',
  bride_name: '',
  wedding_date: '',
  budget: 0,
  style: '',
  venue: '',
  venue_address: '',
  guest_count: 0,
  description: ''
})

const weddingRules: FormRules = {
  title: [{ required: true, message: '请输入婚礼名称', trigger: 'blur' }],
  groom_name: [{ required: true, message: '请输入新郎姓名', trigger: 'blur' }],
  bride_name: [{ required: true, message: '请输入新娘姓名', trigger: 'blur' }],
  wedding_date: [{ required: true, message: '请选择婚礼日期', trigger: 'change' }]
}

const weddingFormRef = ref<FormInstance>()

function formatDate(date: string) {
  return dayjs(date).format('YYYY-MM-DD')
}

function formatNumber(num: number) {
  return num?.toLocaleString() || '0'
}

function statusType(status: string) {
  const types: Record<string, string> = {
    planning: 'warning',
    confirmed: 'success',
    completed: 'primary',
    cancelled: 'danger'
  }
  return types[status] || 'info'
}

function statusText(status: string) {
  const texts: Record<string, string> = {
    planning: '筹备中',
    confirmed: '已确认',
    completed: '已完成',
    cancelled: '已取消'
  }
  return texts[status] || status
}

async function fetchWeddings() {
  loading.value = true
  try {
    const res = await weddingApi.getList({
      search: searchQuery.value,
      status: statusFilter.value,
      page: page.value,
      page_size: pageSize.value
    })
    weddings.value = res.data.list
    total.value = res.data.total
  } catch (error) {
    console.error('Failed to fetch weddings:', error)
  } finally {
    loading.value = false
  }
}

function viewDetail(wedding: Wedding) {
  weddingStore.setCurrentWedding(wedding)
  router.push(`/weddings/${wedding.id}`)
}

function editWedding(wedding: Wedding) {
  editingWedding.value = wedding
  Object.assign(weddingForm, {
    title: wedding.title,
    groom_name: wedding.groom_name,
    bride_name: wedding.bride_name,
    wedding_date: dayjs(wedding.wedding_date).format('YYYY-MM-DD'),
    budget: wedding.budget,
    style: wedding.style || '',
    venue: wedding.venue || '',
    venue_address: wedding.venue_address || '',
    guest_count: wedding.guest_count,
    description: wedding.description || ''
  })
  showCreateDialog.value = true
}

async function saveWedding() {
  if (!weddingFormRef.value) return
  
  await weddingFormRef.value.validate(async (valid) => {
    if (valid) {
      saving.value = true
      try {
        if (editingWedding.value) {
          await weddingApi.update(editingWedding.value.id, weddingForm)
          ElMessage.success('更新成功')
        } else {
          await weddingApi.create(weddingForm)
          ElMessage.success('创建成功')
        }
        showCreateDialog.value = false
        fetchWeddings()
      } catch (error: any) {
        ElMessage.error(error.message || '保存失败')
      } finally {
        saving.value = false
      }
    }
  })
}

async function deleteWedding(wedding: Wedding) {
  try {
    await ElMessageBox.confirm(`确定要删除婚礼"${wedding.title}"吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await weddingApi.delete(wedding.id)
    ElMessage.success('删除成功')
    fetchWeddings()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete wedding:', error)
    }
  }
}

function resetForm() {
  editingWedding.value = null
  Object.assign(weddingForm, {
    title: '',
    groom_name: '',
    bride_name: '',
    wedding_date: '',
    budget: 0,
    style: '',
    venue: '',
    venue_address: '',
    guest_count: 0,
    description: ''
  })
}

onMounted(fetchWeddings)
</script>

<style scoped>
.page-container {
  padding: 0;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-header h2 {
  margin: 0;
}

.filter-bar {
  display: flex;
  gap: 16px;
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  justify-content: flex-end;
}
</style>
