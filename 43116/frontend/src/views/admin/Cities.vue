<template>
  <div>
    <div class="page-header">
      <h1 class="page-title">城市管理</h1>
      <el-button type="primary" :icon="Plus" @click="showAddDialog = true">添加城市</el-button>
    </div>

    <el-table :data="cities" v-loading="loading" style="width: 100%">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="城市名称" min-width="150" />
      <el-table-column prop="code" label="城市代码" width="120" />
      <el-table-column prop="province" label="省份" width="120" />
      <el-table-column label="门店数量" width="100">
        <template #default="{ row }">
          {{ row.stores?.length || 0 }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="150" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link size="small" @click="editCity(row)">编辑</el-button>
          <el-button type="danger" link size="small" @click="deleteCity(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="showAddDialog" :title="editingCity ? '编辑城市' : '添加城市'" width="400px">
      <el-form :model="cityForm" label-width="80px">
        <el-form-item label="城市">
          <el-input v-model="cityForm.name" />
        </el-form-item>
        <el-form-item label="代码">
          <el-input v-model="cityForm.code" />
        </el-form-item>
        <el-form-item label="省份">
          <el-input v-model="cityForm.province" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddDialog = false">取消</el-button>
        <el-button type="primary" @click="saveCity">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { storeApi } from '@/api'
import type { City } from '@/types'

const cities = ref<City[]>([])
const loading = ref(false)
const showAddDialog = ref(false)
const editingCity = ref<City | null>(null)

const cityForm = reactive({
  name: '',
  code: '',
  province: ''
})

onMounted(() => {
  loadCities()
})

const loadCities = async () => {
  loading.value = true
  try {
    const res = await storeApi.getCities()
    cities.value = res.data
  } catch {
    // ignore
  } finally {
    loading.value = false
  }
}

const editCity = (row: City) => {
  editingCity.value = row
  Object.assign(cityForm, {
    name: row.name,
    code: row.code,
    province: row.province
  })
  showAddDialog.value = true
}

const saveCity = async () => {
  try {
    if (editingCity.value) {
      await storeApi.updateCity(editingCity.value.id, cityForm)
      ElMessage.success('更新成功')
    } else {
      await storeApi.createCity(cityForm)
      ElMessage.success('添加成功')
    }
    showAddDialog.value = false
    editingCity.value = null
    loadCities()
  } catch (err: any) {
    ElMessage.error(err.message || '操作失败')
  }
}

const deleteCity = async (row: City) => {
  try {
    await ElMessageBox.confirm(`确定要删除城市 ${row.name} 吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await storeApi.deleteCity(row.id)
    ElMessage.success('删除成功')
    loadCities()
  } catch (err: any) {
    if (err !== 'cancel') {
      ElMessage.error(err.message || '删除失败')
    }
  }
}
</script>
