<template>
  <div>
    <div class="page-header">
      <h1 class="page-title">门店管理</h1>
      <el-button type="primary" :icon="Plus" @click="showAddDialog = true">添加门店</el-button>
    </div>

    <div class="search-bar">
      <el-select v-model="filters.cityId" placeholder="城市" clearable style="width: 150px">
        <el-option
          v-for="city in cities"
          :key="city.id"
          :label="city.name"
          :value="city.id"
        />
      </el-select>
      <el-input
        v-model="filters.keyword"
        placeholder="搜索门店名称/地址"
        clearable
        style="width: 250px"
        @keyup.enter="loadStores"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
      <el-button type="primary" @click="loadStores">搜索</el-button>
    </div>

    <el-table :data="stores" v-loading="loading" style="width: 100%">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="门店名称" min-width="150" />
      <el-table-column label="城市" width="100">
        <template #default="{ row }">
          {{ row.city?.name }}
        </template>
      </el-table-column>
      <el-table-column prop="address" label="地址" min-width="200" />
      <el-table-column prop="phone" label="电话" width="130" />
      <el-table-column prop="business_hours" label="营业时间" width="150" />
      <el-table-column label="状态" width="80">
        <template #default="{ row }">
          <el-tag :type="row.is_active ? 'success' : 'danger'" size="small">
            {{ row.is_active ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="150" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link size="small" @click="editStore(row)">编辑</el-button>
          <el-button type="danger" link size="small" @click="deleteStore(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div style="text-align: center; margin-top: 20px;">
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        layout="total, prev, pager, next"
        @current-change="loadStores"
      />
    </div>

    <el-dialog v-model="showAddDialog" :title="editingStore ? '编辑门店' : '添加门店'" width="500px">
      <el-form :model="storeForm" label-width="100px">
        <el-form-item label="门店名称">
          <el-input v-model="storeForm.name" />
        </el-form-item>
        <el-form-item label="城市">
          <el-select v-model="storeForm.city_id">
            <el-option
              v-for="city in cities"
              :key="city.id"
              :label="city.name"
              :value="city.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="地址">
          <el-input v-model="storeForm.address" />
        </el-form-item>
        <el-form-item label="电话">
          <el-input v-model="storeForm.phone" />
        </el-form-item>
        <el-form-item label="营业时间">
          <el-input v-model="storeForm.business_hours" placeholder="如: 09:00-21:00" />
        </el-form-item>
        <el-form-item label="是否启用">
          <el-switch v-model="storeForm.is_active" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddDialog = false">取消</el-button>
        <el-button type="primary" @click="saveStore">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Plus } from '@element-plus/icons-vue'
import { storeApi } from '@/api'
import type { Store, City } from '@/types'

const stores = ref<Store[]>([])
const cities = ref<City[]>([])
const loading = ref(false)
const showAddDialog = ref(false)
const editingStore = ref<Store | null>(null)

const filters = reactive({
  cityId: undefined as number | undefined,
  keyword: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const storeForm = reactive({
  name: '',
  city_id: undefined as number | undefined,
  address: '',
  phone: '',
  business_hours: '',
  is_active: true
})

onMounted(() => {
  loadStores()
  loadCities()
})

const loadStores = async () => {
  loading.value = true
  try {
    const res = await storeApi.getStores({
      page: pagination.page,
      page_size: pagination.pageSize,
      city_id: filters.cityId,
      keyword: filters.keyword
    })
    stores.value = res.data.items
    pagination.total = res.data.total
  } catch {
    // ignore
  } finally {
    loading.value = false
  }
}

const loadCities = async () => {
  try {
    const res = await storeApi.getCities()
    cities.value = res.data
  } catch {
    // ignore
  }
}

const editStore = (row: Store) => {
  editingStore.value = row
  Object.assign(storeForm, {
    name: row.name,
    city_id: row.city_id,
    address: row.address,
    phone: row.phone,
    business_hours: row.business_hours,
    is_active: row.is_active
  })
  showAddDialog.value = true
}

const saveStore = async () => {
  try {
    if (editingStore.value) {
      await storeApi.updateStore(editingStore.value.id, storeForm)
      ElMessage.success('更新成功')
    } else {
      await storeApi.createStore(storeForm)
      ElMessage.success('添加成功')
    }
    showAddDialog.value = false
    editingStore.value = null
    loadStores()
  } catch (err: any) {
    ElMessage.error(err.message || '操作失败')
  }
}

const deleteStore = async (row: Store) => {
  try {
    await ElMessageBox.confirm(`确定要删除门店 ${row.name} 吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await storeApi.deleteStore(row.id)
    ElMessage.success('删除成功')
    loadStores()
  } catch (err: any) {
    if (err !== 'cancel') {
      ElMessage.error(err.message || '删除失败')
    }
  }
}
</script>
