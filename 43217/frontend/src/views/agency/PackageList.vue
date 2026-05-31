<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElTable, ElTableColumn, ElButton, ElDialog, ElForm, ElFormItem, ElInput, ElInputNumber, ElSelect, ElOption, ElMessage, ElPopconfirm, ElCard, ElTag, ElSwitch, ElIcon } from 'element-plus'
import { Plus, Edit, Delete, SwitchButton } from '@element-plus/icons-vue'
import { getAgencyPackages, createPackage, updatePackage, updatePackagePrice, updatePackageStatus } from '@/api/agency'
import type { Package } from '@/types'

const loading = ref(false)
const packages = ref<Package[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

const dialogVisible = ref(false)
const dialogType = ref<'add' | 'edit' | 'price'>('add')
const currentPackage = ref<Package | null>(null)

const packageForm = ref({
  name: '',
  description: '',
  original_price: 0,
  price: 0,
  suitable_for: '',
  gender_limit: 0,
  min_age: 0,
  max_age: 150,
  notes: ''
})

const priceForm = ref({
  price: 0
})

const formRules = {
  name: [{ required: true, message: '请输入套餐名称', trigger: 'blur' }],
  price: [{ required: true, message: '请输入套餐价格', trigger: 'blur' }]
}

const fetchPackages = async () => {
  loading.value = true
  try {
    const response = await getAgencyPackages({ page: page.value, page_size: pageSize.value })
    packages.value = response.items
    total.value = response.total
  } catch (error) {
    console.error('Failed to fetch packages:', error)
  } finally {
    loading.value = false
  }
}

const handleAdd = () => {
  dialogType.value = 'add'
  currentPackage.value = null
  packageForm.value = {
    name: '',
    description: '',
    original_price: 0,
    price: 0,
    suitable_for: '',
    gender_limit: 0,
    min_age: 0,
    max_age: 150,
    notes: ''
  }
  dialogVisible.value = true
}

const handleEdit = (row: Package) => {
  dialogType.value = 'edit'
  currentPackage.value = row
  packageForm.value = {
    name: row.name,
    description: row.description || '',
    original_price: row.original_price,
    price: row.price,
    suitable_for: row.suitable_for || '',
    gender_limit: row.gender_limit,
    min_age: row.min_age,
    max_age: row.max_age,
    notes: row.notes || ''
  }
  dialogVisible.value = true
}

const handleEditPrice = (row: Package) => {
  dialogType.value = 'price'
  currentPackage.value = row
  priceForm.value = { price: row.price }
  dialogVisible.value = true
}

const handleStatusChange = async (row: Package, status: number) => {
  try {
    await updatePackageStatus(row.id, status)
    ElMessage.success(status === 1 ? '套餐已上架' : '套餐已下架')
    fetchPackages()
  } catch (error) {
    console.error('Failed to update status:', error)
  }
}

const handleSubmit = async () => {
  try {
    if (dialogType.value === 'add') {
      await createPackage(packageForm.value)
      ElMessage.success('添加成功')
    } else if (dialogType.value === 'edit' && currentPackage.value) {
      await updatePackage(currentPackage.value.id, packageForm.value)
      ElMessage.success('更新成功')
    } else if (dialogType.value === 'price' && currentPackage.value) {
      await updatePackagePrice(currentPackage.value.id, priceForm.value.price)
      ElMessage.success('调价成功')
    }
    dialogVisible.value = false
    fetchPackages()
  } catch (error) {
    console.error('Failed to submit:', error)
  }
}

const handlePageChange = (newPage: number) => {
  page.value = newPage
  fetchPackages()
}

const handleSizeChange = (newSize: number) => {
  pageSize.value = newSize
  page.value = 1
  fetchPackages()
}

onMounted(() => {
  fetchPackages()
})
</script>

<template>
  <div class="package-list">
    <ElCard>
      <template #header>
        <div class="card-header">
          <span>套餐管理</span>
          <ElButton type="primary" :icon="Plus" @click="handleAdd">
            添加套餐
          </ElButton>
        </div>
      </template>

      <ElTable :data="packages" v-loading="loading" border stripe>
        <ElTableColumn prop="name" label="套餐名称" width="200" />
        <ElTableColumn prop="description" label="描述" show-overflow-tooltip />
        <ElTableColumn label="价格" width="150">
          <template #default="{ row }">
            <div class="price-info">
              <span class="current-price">¥{{ row.price }}</span>
              <span v-if="row.original_price > row.price" class="original-price">
                ¥{{ row.original_price }}
              </span>
            </div>
          </template>
        </ElTableColumn>
        <ElTableColumn prop="suitable_for" label="适用人群" width="120" />
        <ElTableColumn label="浏览/销量" width="120">
          <template #default="{ row }">
            <div>{{ row.view_count }} / {{ row.sale_count }}</div>
          </template>
        </ElTableColumn>
        <ElTableColumn label="状态" width="100">
          <template #default="{ row }">
            <ElTag :type="row.status === 1 ? 'success' : 'info'">
              {{ row.status === 1 ? '上架' : '下架' }}
            </ElTag>
          </template>
        </ElTableColumn>
        <ElTableColumn label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <ElButton type="primary" link :icon="Edit" @click="handleEdit(row)">
              编辑
            </ElButton>
            <ElButton type="success" link @click="handleEditPrice(row)">
              调价
            </ElButton>
            <ElSwitch
              :model-value="row.status === 1"
              size="small"
              @change="(val: boolean) => handleStatusChange(row, val ? 1 : 0)"
            />
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
      :title="dialogType === 'add' ? '添加套餐' : dialogType === 'edit' ? '编辑套餐' : '调整价格'"
      width="600px"
    >
      <template v-if="dialogType !== 'price'">
        <ElForm :model="packageForm" :rules="formRules" label-width="100px">
          <ElFormItem label="套餐名称" prop="name">
            <ElInput v-model="packageForm.name" placeholder="请输入套餐名称" />
          </ElFormItem>
          <ElFormItem label="套餐描述">
            <ElInput v-model="packageForm.description" type="textarea" :rows="3" placeholder="请输入套餐描述" />
          </ElFormItem>
          <ElFormItem label="原价">
            <ElInputNumber v-model="packageForm.original_price" :min="0" :precision="2" />
          </ElFormItem>
          <ElFormItem label="现价" prop="price">
            <ElInputNumber v-model="packageForm.price" :min="0" :precision="2" />
          </ElFormItem>
          <ElFormItem label="适用人群">
            <ElInput v-model="packageForm.suitable_for" placeholder="如：全人群/男性/女性" />
          </ElFormItem>
          <ElFormItem label="性别限制">
            <ElSelect v-model="packageForm.gender_limit" style="width: 100%">
              <ElOption label="不限" :value="0" />
              <ElOption label="男" :value="1" />
              <ElOption label="女" :value="2" />
            </ElSelect>
          </ElFormItem>
          <ElFormItem label="年龄范围">
            <div class="age-range">
              <ElInputNumber v-model="packageForm.min_age" :min="0" placeholder="最小年龄" />
              <span>-</span>
              <ElInputNumber v-model="packageForm.max_age" :min="0" placeholder="最大年龄" />
            </div>
          </ElFormItem>
          <ElFormItem label="注意事项">
            <ElInput v-model="packageForm.notes" type="textarea" :rows="2" placeholder="请输入注意事项" />
          </ElFormItem>
        </ElForm>
      </template>
      <template v-else>
        <ElForm :model="priceForm" label-width="100px">
          <ElFormItem label="当前价格">
            <span class="price-display">¥{{ currentPackage?.price }}</span>
          </ElFormItem>
          <ElFormItem label="新价格" prop="price">
            <ElInputNumber v-model="priceForm.price" :min="0" :precision="2" />
          </ElFormItem>
        </ElForm>
      </template>
      <template #footer>
        <ElButton @click="dialogVisible = false">取消</ElButton>
        <ElButton type="primary" @click="handleSubmit">确定</ElButton>
      </template>
    </ElDialog>
  </div>
</template>

<style scoped lang="scss">
.package-list {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .price-info {
    .current-price {
      color: #f56c6c;
      font-size: 16px;
      font-weight: bold;
    }

    .original-price {
      color: #909399;
      font-size: 12px;
      text-decoration: line-through;
      margin-left: 5px;
    }
  }

  .age-range {
    display: flex;
    align-items: center;
    gap: 10px;
  }

  .price-display {
    color: #f56c6c;
    font-size: 18px;
    font-weight: bold;
  }

  .pagination-container {
    display: flex;
    justify-content: flex-end;
    margin-top: 20px;
  }
}
</style>
