<template>
  <div class="vendor-list">
    <div class="page-header">
      <el-input
        v-model="searchQuery"
        placeholder="搜索供应商..."
        :prefix-icon="Search"
        clearable
        style="width: 200px"
        @clear="fetchVendors"
        @keyup.enter="fetchVendors"
      />
      <el-select v-model="categoryFilter" placeholder="分类" clearable style="width: 140px" @change="fetchVendors">
        <el-option v-for="cat in categories" :key="cat" :label="cat" :value="cat" />
      </el-select>
      <el-button type="primary" :icon="Plus" @click="showCreateDialog = true">
        添加供应商
      </el-button>
    </div>

    <el-table :data="vendors" v-loading="loading" stripe>
      <el-table-column prop="name" label="供应商名称" min-width="160" />
      <el-table-column prop="category" label="分类" width="120" />
      <el-table-column prop="contact_person" label="联系人" width="120" />
      <el-table-column prop="phone" label="电话" width="130" />
      <el-table-column label="评分" width="140">
        <template #default="{ row }">
          <div class="rating-cell">
            <el-rate :model-value="row.rating" disabled size="small" />
            <span style="margin-left: 6px; color: #909399;">{{ row.rating?.toFixed(1) || 0 }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 'active' ? 'success' : 'info'">
            {{ row.status === 'active' ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="220" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link @click="viewVendor(row)">详情</el-button>
          <el-button type="primary" link @click="editVendor(row)">编辑</el-button>
          <el-button type="success" link @click="openReviewDialog(row)">评价</el-button>
          <el-button type="danger" link @click="deleteVendor(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="showCreateDialog" :title="editingVendor ? '编辑供应商' : '添加供应商'" width="500px">
      <el-form ref="vendorForm" :model="vendorForm" :rules="vendorRules" label-width="100px">
        <el-form-item label="供应商名称" prop="name">
          <el-input v-model="vendorForm.name" />
        </el-form-item>
        <el-form-item label="分类" prop="category">
          <el-select v-model="vendorForm.category" placeholder="请选择分类" style="width: 100%">
            <el-option v-for="cat in categories" :key="cat" :label="cat" :value="cat" />
          </el-select>
        </el-form-item>
        <el-form-item label="联系人">
          <el-input v-model="vendorForm.contact_person" />
        </el-form-item>
        <el-form-item label="电话">
          <el-input v-model="vendorForm.phone" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="vendorForm.email" />
        </el-form-item>
        <el-form-item label="地址">
          <el-input v-model="vendorForm.address" />
        </el-form-item>
        <el-form-item label="服务范围">
          <el-input v-model="vendorForm.service_area" />
        </el-form-item>
        <el-form-item label="价格区间">
          <el-input v-model="vendorForm.price_range" placeholder="如：¥1000-5000" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="vendorForm.notes" type="textarea" :rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="saveVendor">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showDetailDialog" :title="currentVendor?.name || '供应商详情'" width="600px">
      <el-descriptions :column="2" border v-if="currentVendor">
        <el-descriptions-item label="分类">{{ currentVendor.category }}</el-descriptions-item>
        <el-descriptions-item label="联系人">{{ currentVendor.contact_person || '-' }}</el-descriptions-item>
        <el-descriptions-item label="电话">{{ currentVendor.phone || '-' }}</el-descriptions-item>
        <el-descriptions-item label="邮箱">{{ currentVendor.email || '-' }}</el-descriptions-item>
        <el-descriptions-item label="地址" :span="2">{{ currentVendor.address || '-' }}</el-descriptions-item>
        <el-descriptions-item label="服务范围" :span="2">{{ currentVendor.service_area || '-' }}</el-descriptions-item>
        <el-descriptions-item label="价格区间">{{ currentVendor.price_range || '-' }}</el-descriptions-item>
        <el-descriptions-item label="评分">
          <el-rate :model-value="currentVendor.rating" disabled size="small" />
          ({{ currentVendor.rating?.toFixed(1) || 0 }})
        </el-descriptions-item>
        <el-descriptions-item label="备注" :span="2">{{ currentVendor.notes || '-' }}</el-descriptions-item>
      </el-descriptions>
      
      <div style="margin-top: 24px;">
        <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px;">
          <h4 style="margin: 0;">用户评价 ({{ currentVendor?.review_count || 0 }})</h4>
        </div>
        <div v-if="vendorReviews.length === 0" style="color: #909399; text-align: center; padding: 20px;">
          暂无评价
        </div>
        <div v-else class="review-list">
          <div v-for="review in vendorReviews" :key="review.id" class="review-item">
            <div class="review-header">
              <el-rate :model-value="review.rating" disabled size="small" />
              <span class="review-date">{{ formatDate(review.created_at) }}</span>
            </div>
            <div v-if="review.content" class="review-content">{{ review.content }}</div>
          </div>
        </div>
      </div>
    </el-dialog>

    <el-dialog v-model="showReviewDialog" title="评价供应商" width="450px" :close-on-click-modal="false">
      <div v-if="currentVendor" style="margin-bottom: 16px;">
        供应商：<strong>{{ currentVendor.name }}</strong>
      </div>
      <el-form ref="reviewFormRef" :model="reviewForm" :rules="reviewRules" label-width="80px">
        <el-form-item label="评分" prop="rating">
          <el-rate v-model="reviewForm.rating" :max="5" />
        </el-form-item>
        <el-form-item label="评价">
          <el-input
            v-model="reviewForm.content"
            type="textarea"
            :rows="4"
            placeholder="请分享您的使用体验..."
            maxlength="500"
            show-word-limit
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showReviewDialog = false">取消</el-button>
        <el-button type="primary" :loading="submittingReview" @click="submitReview">提交评价</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { vendorApi } from '@/api/vendor'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Plus, Search } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import type { Vendor } from '@/types'

const props = defineProps<{
  weddingId?: number
}>()

const loading = ref(false)
const saving = ref(false)
const submittingReview = ref(false)
const vendors = ref<Vendor[]>([])
const categories = ref<string[]>([])
const searchQuery = ref('')
const categoryFilter = ref('')
const showCreateDialog = ref(false)
const showDetailDialog = ref(false)
const showReviewDialog = ref(false)
const editingVendor = ref<Vendor | null>(null)
const currentVendor = ref<Vendor | null>(null)
const vendorReviews = ref<any[]>([])

const vendorForm = reactive({
  name: '',
  category: '',
  contact_person: '',
  phone: '',
  email: '',
  address: '',
  service_area: '',
  price_range: '',
  notes: '',
  wedding_id: undefined as number | undefined
})

const reviewForm = reactive({
  rating: 5,
  content: ''
})

const vendorRules: FormRules = {
  name: [{ required: true, message: '请输入供应商名称', trigger: 'blur' }],
  category: [{ required: true, message: '请选择分类', trigger: 'change' }]
}

const reviewRules: FormRules = {
  rating: [{ required: true, message: '请选择评分', trigger: 'change' }]
}

const vendorFormRef = ref<FormInstance>()
const reviewFormRef = ref<FormInstance>()

function formatDate(date: string) {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

async function fetchCategories() {
  try {
    const res = await vendorApi.getCategories()
    categories.value = res.data
  } catch (error) {
    console.error('Failed to fetch categories:', error)
  }
}

async function fetchVendors() {
  loading.value = true
  try {
    const params: any = {
      search: searchQuery.value,
      category: categoryFilter.value
    }
    if (props.weddingId) {
      params.wedding_id = props.weddingId
    }
    const res = await vendorApi.getList(params)
    vendors.value = res.data.list
  } catch (error) {
    console.error('Failed to fetch vendors:', error)
  } finally {
    loading.value = false
  }
}

async function viewVendor(vendor: Vendor) {
  currentVendor.value = vendor
  try {
    const res = await vendorApi.getById(vendor.id)
    vendorReviews.value = res.data.reviews || []
  } catch (error) {
    console.error('Failed to fetch vendor details:', error)
  }
  showDetailDialog.value = true
}

function editVendor(vendor: Vendor) {
  editingVendor.value = vendor
  Object.assign(vendorForm, {
    name: vendor.name,
    category: vendor.category,
    contact_person: vendor.contact_person || '',
    phone: vendor.phone || '',
    email: vendor.email || '',
    address: vendor.address || '',
    service_area: vendor.service_area || '',
    price_range: vendor.price_range || '',
    notes: vendor.notes || '',
    wedding_id: vendor.wedding_id
  })
  showCreateDialog.value = true
}

async function saveVendor() {
  if (!vendorFormRef.value) return
  
  await vendorFormRef.value.validate(async (valid) => {
    if (valid) {
      saving.value = true
      try {
        if (props.weddingId) {
          vendorForm.wedding_id = props.weddingId
        }
        if (editingVendor.value) {
          await vendorApi.update(editingVendor.value.id, vendorForm)
          ElMessage.success('更新成功')
        } else {
          await vendorApi.create(vendorForm)
          ElMessage.success('创建成功')
        }
        showCreateDialog.value = false
        fetchVendors()
      } catch (error: any) {
        ElMessage.error(error.message || '保存失败')
      } finally {
        saving.value = false
      }
    }
  })
}

async function deleteVendor(vendor: Vendor) {
  try {
    await ElMessageBox.confirm(`确定要删除供应商"${vendor.name}"吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await vendorApi.delete(vendor.id)
    ElMessage.success('删除成功')
    fetchVendors()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete vendor:', error)
    }
  }
}

function openReviewDialog(vendor: Vendor) {
  currentVendor.value = vendor
  reviewForm.rating = 5
  reviewForm.content = ''
  showReviewDialog.value = true
}

async function submitReview() {
  if (!reviewFormRef.value || !currentVendor.value) return
  
  await reviewFormRef.value.validate(async (valid) => {
    if (valid) {
      submittingReview.value = true
      try {
        await vendorApi.addReview(currentVendor.value.id, reviewForm)
        ElMessage.success('评价提交成功')
        showReviewDialog.value = false
        fetchVendors()
        if (showDetailDialog.value) {
          viewVendor(currentVendor.value)
        }
      } catch (error: any) {
        ElMessage.error(error.message || '提交失败')
      } finally {
        submittingReview.value = false
      }
    }
  })
}

onMounted(() => {
  fetchCategories()
  fetchVendors()
})
</script>

<style scoped>
.vendor-list {
  padding: 0;
}

.page-header {
  display: flex;
  gap: 16px;
  margin-bottom: 20px;
}

.rating-cell {
  display: flex;
  align-items: center;
}

.review-list {
  max-height: 300px;
  overflow-y: auto;
}

.review-item {
  padding: 12px;
  border: 1px solid #ebeef5;
  border-radius: 6px;
  margin-bottom: 12px;
}

.review-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.review-date {
  color: #909399;
  font-size: 12px;
}

.review-content {
  color: #606266;
  font-size: 14px;
  line-height: 1.6;
}
</style>
