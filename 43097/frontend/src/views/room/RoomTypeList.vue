<template>
  <div class="room-type-container">
    <el-card shadow="hover">
      <div class="header-bar">
        <div class="search-bar">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索房型名称"
            clearable
            :prefix-icon="Search"
            style="width: 280px"
            @keyup.enter="handleSearch"
          />
          <el-button type="primary" :icon="Search" @click="handleSearch">搜索</el-button>
          <el-button :icon="Refresh" @click="resetSearch">重置</el-button>
        </div>
        <el-button type="primary" :icon="Plus" @click="handleAdd">添加房型</el-button>
      </div>

      <el-table :data="tableData" v-loading="loading" border stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="房型名称" min-width="120" />
        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
        <el-table-column prop="basePrice" label="基础价格" width="120">
          <template #default="{ row }">
            <span class="price-text">¥{{ row.basePrice?.toFixed(2) || '0.00' }}/晚</span>
          </template>
        </el-table-column>
        <el-table-column prop="bedCount" label="床数" width="100">
          <template #default="{ row }">
            {{ row.bedCount }}张
          </template>
        </el-table-column>
        <el-table-column prop="maxGuests" label="最多入住" width="100">
          <template #default="{ row }">
            {{ row.maxGuests }}人
          </template>
        </el-table-column>
        <el-table-column prop="facilities" label="设施" min-width="200">
          <template #default="{ row }">
            <el-tag
              v-for="(facility, index) in row.facilities"
              :key="index"
              size="small"
              style="margin-right: 4px; margin-bottom: 4px"
            >
              {{ facility }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link :icon="Edit" @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" link :icon="Delete" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-bar">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="fetchList"
          @current-change="fetchList"
        />
      </div>
    </el-card>

    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑房型' : '添加房型'"
      width="600px"
      destroy-on-close
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="100px"
        label-position="right"
      >
        <el-form-item label="房型名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入房型名称" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="formData.description"
            type="textarea"
            :rows="3"
            placeholder="请输入房型描述"
          />
        </el-form-item>
        <el-form-item label="基础价格" prop="basePrice">
          <el-input-number
            v-model="formData.basePrice"
            :min="0"
            :precision="2"
            :step="10"
            style="width: 100%"
            placeholder="请输入基础价格"
          />
        </el-form-item>
        <el-form-item label="床数" prop="bedCount">
          <el-input-number
            v-model="formData.bedCount"
            :min="1"
            :max="10"
            style="width: 100%"
            placeholder="请输入床数"
          />
        </el-form-item>
        <el-form-item label="最多入住" prop="maxGuests">
          <el-input-number
            v-model="formData.maxGuests"
            :min="1"
            :max="10"
            style="width: 100%"
            placeholder="请输入最多入住人数"
          />
        </el-form-item>
        <el-form-item label="设施" prop="facilities">
          <el-select
            v-model="formData.facilities"
            multiple
            filterable
            allow-create
            default-first-option
            placeholder="选择或输入设施"
            style="width: 100%"
          >
            <el-option
              v-for="facility in facilityOptions"
              :key="facility"
              :label="facility"
              :value="facility"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="detailVisible" title="房型详情" width="500px">
      <el-descriptions :column="2" border v-if="currentRow">
        <el-descriptions-item label="房型名称">{{ currentRow.name }}</el-descriptions-item>
        <el-descriptions-item label="基础价格">¥{{ currentRow.basePrice?.toFixed(2) || '0.00' }}/晚</el-descriptions-item>
        <el-descriptions-item label="床数">{{ currentRow.bedCount }}张</el-descriptions-item>
        <el-descriptions-item label="最多入住">{{ currentRow.maxGuests }}人</el-descriptions-item>
        <el-descriptions-item label="描述" :span="2">
          {{ currentRow.description || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="设施" :span="2">
          <el-tag
            v-for="(facility, index) in currentRow.facilities"
            :key="index"
            size="small"
            style="margin-right: 4px; margin-bottom: 4px"
          >
            {{ facility }}
          </el-tag>
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Search, Refresh, Plus, Edit, Delete } from '@element-plus/icons-vue'
import {
  getRoomTypeList,
  createRoomType,
  updateRoomType,
  deleteRoomType
} from '@/api/room'
import type { RoomType, PageParams } from '@/types'

const loading = ref(false)
const submitLoading = ref(false)
const searchKeyword = ref('')
const tableData = ref<RoomType[]>([])
const dialogVisible = ref(false)
const detailVisible = ref(false)
const isEdit = ref(false)
const currentRow = ref<RoomType | null>(null)
const formRef = ref<FormInstance>()

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const facilityOptions = [
  '免费WiFi',
  '空调',
  '电视',
  '独立卫浴',
  '24小时热水',
  '吹风机',
  '保险箱',
  '迷你吧',
  '冰箱',
  '电热水壶',
  '咖啡机',
  '浴袍',
  '拖鞋',
  '免费洗漱用品',
  '书桌',
  '衣柜',
  '窗户',
  '阳台',
  '景观房',
  '禁烟房'
]

const formData = reactive<Partial<RoomType>>({
  name: '',
  description: '',
  basePrice: 0,
  bedCount: 1,
  maxGuests: 1,
  facilities: []
})

const formRules: FormRules = {
  name: [{ required: true, message: '请输入房型名称', trigger: 'blur' }],
  basePrice: [{ required: true, message: '请输入基础价格', trigger: 'blur' }],
  bedCount: [{ required: true, message: '请输入床数', trigger: 'blur' }],
  maxGuests: [{ required: true, message: '请输入最多入住人数', trigger: 'blur' }]
}

const fetchList = async () => {
  loading.value = true
  try {
    const params: PageParams = {
      page: pagination.page,
      pageSize: pagination.pageSize
    }
    if (searchKeyword.value) {
      params.keyword = searchKeyword.value
    }
    const res = await getRoomTypeList(params)
    tableData.value = res.list
    pagination.total = res.total
  } catch (error) {
    console.error('Failed to fetch room type list:', error)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  fetchList()
}

const resetSearch = () => {
  searchKeyword.value = ''
  pagination.page = 1
  fetchList()
}

const handleAdd = () => {
  isEdit.value = false
  Object.assign(formData, {
    name: '',
    description: '',
    basePrice: 0,
    bedCount: 1,
    maxGuests: 1,
    facilities: []
  })
  dialogVisible.value = true
}

const handleEdit = (row: RoomType) => {
  isEdit.value = true
  currentRow.value = row
  Object.assign(formData, {
    id: row.id,
    name: row.name,
    description: row.description,
    basePrice: row.basePrice,
    bedCount: row.bedCount,
    maxGuests: row.maxGuests,
    facilities: [...row.facilities]
  })
  dialogVisible.value = true
}

const handleDelete = (row: RoomType) => {
  ElMessageBox.confirm(`确定要删除房型"${row.name}"吗？`, '删除确认', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })
    .then(async () => {
      try {
        await deleteRoomType(row.id)
        ElMessage.success('删除成功')
        fetchList()
      } catch (error) {
        console.error('Failed to delete room type:', error)
      }
    })
    .catch(() => {})
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitLoading.value = true
      try {
        if (isEdit.value && formData.id) {
          await updateRoomType(formData.id, formData)
          ElMessage.success('编辑成功')
        } else {
          await createRoomType(formData as Omit<RoomType, 'id' | 'createdAt' | 'updatedAt'>)
          ElMessage.success('添加成功')
        }
        dialogVisible.value = false
        fetchList()
      } catch (error) {
        console.error('Failed to submit room type:', error)
      } finally {
        submitLoading.value = false
      }
    }
  })
}

onMounted(() => {
  fetchList()
})
</script>

<style scoped lang="scss">
.room-type-container {
  .header-bar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;

    .search-bar {
      display: flex;
      gap: 12px;
    }
  }

  .pagination-bar {
    display: flex;
    justify-content: flex-end;
    margin-top: 20px;
  }

  .price-text {
    color: #f56c6c;
    font-weight: 600;
  }
}
</style>
