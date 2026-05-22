<template>
  <div class="page-container">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>库存管理</span>
          <div class="header-actions">
            <el-button type="warning" @click="$router.push('/product-records')">库存记录</el-button>
            <el-button type="primary" :icon="Plus" @click="addVisible = true">添加产品</el-button>
          </div>
        </div>
      </template>

      <el-table :data="products" stripe v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="产品名称" />
        <el-table-column prop="category" label="分类" width="100" />
        <el-table-column prop="unit" label="单位" width="80" />
        <el-table-column prop="stock" label="库存" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.stock <= row.threshold" type="danger">{{ row.stock }}</el-tag>
            <span v-else>{{ row.stock }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="threshold" label="预警阈值" width="100" />
        <el-table-column prop="price" label="成本价" width="100">
          <template #default="{ row }">¥{{ row.price?.toFixed(2) }}</template>
        </el-table-column>
        <el-table-column prop="retail_price" label="零售价" width="100">
          <template #default="{ row }">¥{{ row.retail_price?.toFixed(2) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button type="success" link @click="handleAddStock(row)">入库</el-button>
            <el-button type="warning" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="fetchList"
          @size-change="fetchList"
        />
      </div>
    </el-card>

    <el-dialog v-model="addVisible" :title="editing ? '编辑产品' : '添加产品'" width="500px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="产品名称" required>
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="分类">
          <el-input v-model="form.category" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" />
        </el-form-item>
        <el-form-item label="单位">
          <el-input v-model="form.unit" />
        </el-form-item>
        <el-form-item label="初始库存" v-if="!editing">
          <el-input-number v-model="form.stock" :min="0" />
        </el-form-item>
        <el-form-item label="预警阈值">
          <el-input-number v-model="form.threshold" :min="1" />
        </el-form-item>
        <el-form-item label="成本价">
          <el-input-number v-model="form.price" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="零售价">
          <el-input-number v-model="form.retail_price" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="供应商">
          <el-input v-model="form.supplier" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="addVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="stockVisible" title="入库" width="400px">
      <el-form :model="stockForm" label-width="80px">
        <el-form-item label="数量" required>
          <el-input-number v-model="stockForm.quantity" :min="1" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="stockForm.remark" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="stockVisible = false">取消</el-button>
        <el-button type="primary" :loading="stocking" @click="handleStockSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { getProducts, createProduct, updateProduct, deleteProduct, addStock } from '@/api/product'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import type { Product } from '@/types'

const loading = ref(false)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const products = ref<Product[]>([])
const addVisible = ref(false)
const stockVisible = ref(false)
const editing = ref(false)
const editingId = ref<number | null>(null)
const saving = ref(false)
const stocking = ref(false)
const stockingId = ref<number | null>(null)

const form = reactive({
  name: '',
  category: '',
  description: '',
  unit: '件',
  stock: 0,
  threshold: 10,
  price: 0,
  retail_price: 0,
  supplier: ''
})

const stockForm = reactive({
  quantity: 1,
  remark: ''
})

const fetchList = async () => {
  loading.value = true
  try {
    const res = await getProducts({ page: page.value, page_size: pageSize.value })
    products.value = res.data.list
    total.value = res.data.total
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const handleEdit = (row: Product) => {
  editing.value = true
  editingId.value = row.id
  Object.assign(form, {
    name: row.name,
    category: row.category,
    description: row.description,
    unit: row.unit,
    threshold: row.threshold,
    price: row.price,
    retail_price: row.retail_price,
    supplier: row.supplier
  })
  addVisible.value = true
}

const handleDelete = async (row: Product) => {
  try {
    await ElMessageBox.confirm(`确定删除产品"${row.name}"吗？`, '提示', { type: 'warning' })
    await deleteProduct(row.id)
    ElMessage.success('删除成功')
    fetchList()
  } catch (e) {
    if (e !== 'cancel') console.error(e)
  }
}

const handleAddStock = (row: Product) => {
  stockingId.value = row.id
  stockForm.quantity = 1
  stockForm.remark = ''
  stockVisible.value = true
}

const handleStockSubmit = async () => {
  if (!stockForm.quantity) {
    ElMessage.warning('请输入数量')
    return
  }

  stocking.value = true
  try {
    await addStock({
      product_id: stockingId.value!,
      quantity: stockForm.quantity,
      remark: stockForm.remark
    })
    ElMessage.success('入库成功')
    stockVisible.value = false
    fetchList()
  } catch (e: any) {
    ElMessage.error(e.message || '入库失败')
  } finally {
    stocking.value = false
  }
}

const handleSave = async () => {
  if (!form.name) {
    ElMessage.warning('请填写产品名称')
    return
  }

  saving.value = true
  try {
    if (editing.value && editingId.value) {
      await updateProduct(editingId.value, form)
      ElMessage.success('更新成功')
    } else {
      await createProduct(form)
      ElMessage.success('添加成功')
    }
    addVisible.value = false
    fetchList()
  } catch (e: any) {
    ElMessage.error(e.message || '保存失败')
  } finally {
    saving.value = false
  }
}

onMounted(fetchList)
</script>

<style scoped lang="scss">
.page-container {
  .card-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .header-actions {
    display: flex;
    gap: 10px;
  }

  .pagination {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }
}
</style>
