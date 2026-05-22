<template>
  <div>
    <div class="page-header">
      <h1 class="page-title">车辆管理</h1>
      <el-button type="primary" :icon="Plus" @click="showAddDialog = true">添加车辆</el-button>
    </div>

    <div class="search-bar">
      <el-input
        v-model="filters.keyword"
        placeholder="搜索品牌/型号/车牌"
        clearable
        style="width: 250px"
        @keyup.enter="loadCars"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
      <el-select v-model="filters.status" placeholder="状态" clearable style="width: 120px">
        <el-option label="可用" value="available" />
        <el-option label="出租中" value="rented" />
        <el-option label="维护中" value="maintenance" />
        <el-option label="已停用" value="disabled" />
      </el-select>
      <el-button type="primary" @click="loadCars">搜索</el-button>
    </div>

    <el-table :data="cars" v-loading="loading" style="width: 100%">
      <el-table-column label="图片" width="100">
        <template #default="{ row }">
          <el-image
            :src="row.images?.[0]?.url || 'https://picsum.photos/100/100?random=' + row.id"
            style="width: 60px; height: 60px; border-radius: 4px;"
            fit="cover"
          />
        </template>
      </el-table-column>
      <el-table-column label="车型" min-width="150">
        <template #default="{ row }">
          {{ row.brand }} {{ row.model }}
        </template>
      </el-table-column>
      <el-table-column prop="year" label="年份" width="80" />
      <el-table-column prop="seats" label="座位" width="70" />
      <el-table-column prop="transmission" label="变速箱" width="100" />
      <el-table-column prop="daily_rent" label="日租金" width="100">
        <template #default="{ row }">
          ¥{{ row.daily_rent.toFixed(2) }}
        </template>
      </el-table-column>
      <el-table-column prop="license_plate" label="车牌号" width="120" />
      <el-table-column label="评分" width="150">
        <template #default="{ row }">
          <el-rate :model-value="row.rating" disabled size="small" />
          <span style="margin-left: 5px;">({{ row.review_count }})</span>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link size="small" @click="editCar(row)">编辑</el-button>
          <el-button type="primary" link size="small" @click="uploadImages(row)">上传图片</el-button>
          <el-button type="danger" link size="small" @click="deleteCar(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div style="text-align: center; margin-top: 20px;">
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next, jumper"
        :page-sizes="[10, 20, 50]"
        @size-change="loadCars"
        @current-change="loadCars"
      />
    </div>

    <el-dialog v-model="showAddDialog" :title="editingCar ? '编辑车辆' : '添加车辆'" width="600px">
      <el-form :model="carForm" label-width="100px">
        <el-form-item label="品牌">
          <el-input v-model="carForm.brand" />
        </el-form-item>
        <el-form-item label="型号">
          <el-input v-model="carForm.model" />
        </el-form-item>
        <el-form-item label="年份">
          <el-input-number v-model="carForm.year" :min="2010" :max="2025" />
        </el-form-item>
        <el-form-item label="座位数">
          <el-input-number v-model="carForm.seats" :min="2" :max="20" />
        </el-form-item>
        <el-form-item label="变速箱">
          <el-select v-model="carForm.transmission">
            <el-option label="自动" value="自动" />
            <el-option label="手动" value="手动" />
          </el-select>
        </el-form-item>
        <el-form-item label="燃料类型">
          <el-select v-model="carForm.fuel_type">
            <el-option label="汽油" value="汽油" />
            <el-option label="柴油" value="柴油" />
            <el-option label="电动" value="电动" />
            <el-option label="混动" value="混动" />
          </el-select>
        </el-form-item>
        <el-form-item label="日租金">
          <el-input-number v-model="carForm.daily_rent" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="押金">
          <el-input-number v-model="carForm.deposit" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="车牌号">
          <el-input v-model="carForm.license_plate" />
        </el-form-item>
        <el-form-item label="颜色">
          <el-input v-model="carForm.color" />
        </el-form-item>
        <el-form-item label="门店">
          <el-select v-model="carForm.store_id">
            <el-option
              v-for="store in stores"
              :key="store.id"
              :label="store.name"
              :value="store.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="carForm.description" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddDialog = false">取消</el-button>
        <el-button type="primary" @click="saveCar">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showUploadDialog" title="上传车辆图片" width="500px">
      <el-upload
        :action="uploadUrl"
        :headers="uploadHeaders"
        :data="{}"
        multiple
        :show-file-list="true"
        :on-success="handleUploadSuccess"
        accept="image/*"
      >
        <el-button type="primary">选择图片</el-button>
      </el-upload>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Plus } from '@element-plus/icons-vue'
import { carApi, storeApi } from '@/api'
import { useUserStore } from '@/stores/user'
import type { Car, Store } from '@/types'

const userStore = useUserStore()
const cars = ref<Car[]>([])
const stores = ref<Store[]>([])
const loading = ref(false)
const showAddDialog = ref(false)
const showUploadDialog = ref(false)
const editingCar = ref<Car | null>(null)
const uploadingCar = ref<Car | null>(null)

const uploadUrl = ref('')
const uploadHeaders = ref({})

const filters = reactive({
  keyword: '',
  status: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const carForm = reactive({
  brand: '',
  model: '',
  year: 2024,
  seats: 5,
  transmission: '自动',
  fuel_type: '汽油',
  daily_rent: 0,
  deposit: 0,
  license_plate: '',
  color: '',
  store_id: undefined as number | undefined,
  description: ''
})

onMounted(() => {
  loadCars()
  loadStores()
})

const loadCars = async () => {
  loading.value = true
  try {
    const res = await carApi.getCars({
      page: pagination.page,
      page_size: pagination.pageSize,
      keyword: filters.keyword,
      status: filters.status
    })
    cars.value = res.data.items
    pagination.total = res.data.total
  } catch {
    // ignore
  } finally {
    loading.value = false
  }
}

const loadStores = async () => {
  try {
    const res = await storeApi.getStores({ page: 1, page_size: 100 })
    stores.value = res.data.items
  } catch {
    // ignore
  }
}

const editCar = (row: Car) => {
  editingCar.value = row
  Object.assign(carForm, {
    brand: row.brand,
    model: row.model,
    year: row.year,
    seats: row.seats,
    transmission: row.transmission,
    fuel_type: row.fuel_type,
    daily_rent: row.daily_rent,
    deposit: row.deposit,
    license_plate: row.license_plate,
    color: row.color,
    store_id: row.store_id,
    description: row.description
  })
  showAddDialog.value = true
}

const saveCar = async () => {
  try {
    if (editingCar.value) {
      await carApi.updateCar(editingCar.value.id, carForm)
      ElMessage.success('更新成功')
    } else {
      await carApi.createCar(carForm)
      ElMessage.success('添加成功')
    }
    showAddDialog.value = false
    editingCar.value = null
    loadCars()
  } catch (err: any) {
    ElMessage.error(err.message || '操作失败')
  }
}

const uploadImages = (row: Car) => {
  uploadingCar.value = row
  uploadUrl.value = `/api/cars/${row.id}/batch-upload`
  uploadHeaders.value = {
    Authorization: `Bearer ${userStore.accessToken}`
  }
  showUploadDialog.value = true
}

const handleUploadSuccess = () => {
  ElMessage.success('上传成功')
  showUploadDialog.value = false
  loadCars()
}

const deleteCar = async (row: Car) => {
  try {
    await ElMessageBox.confirm(`确定要删除车辆 ${row.brand} ${row.model} 吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await carApi.deleteCar(row.id)
    ElMessage.success('删除成功')
    loadCars()
  } catch (err: any) {
    if (err !== 'cancel') {
      ElMessage.error(err.message || '删除失败')
    }
  }
}

const getStatusType = (status: string) => {
  const map: Record<string, string> = {
    available: 'success',
    rented: 'primary',
    maintenance: 'warning',
    disabled: 'danger'
  }
  return map[status] || 'info'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    available: '可用',
    rented: '出租中',
    maintenance: '维护中',
    disabled: '已停用'
  }
  return map[status] || status
}
</script>
