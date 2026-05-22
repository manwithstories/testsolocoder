<template>
  <div>
    <div class="page-header">
      <h1 class="page-title">优惠码管理</h1>
      <el-button type="primary" :icon="Plus" @click="showAddDialog = true">添加优惠码</el-button>
    </div>

    <el-table :data="promos" v-loading="loading" style="width: 100%">
      <el-table-column prop="code" label="优惠码" width="150" />
      <el-table-column prop="name" label="名称" min-width="150" />
      <el-table-column prop="type" label="类型" width="100" />
      <el-table-column prop="value" label="折扣/金额" width="120">
        <template #default="{ row }">
          {{ row.type === 'percent' ? row.value + '%' : '¥' + row.value }}
        </template>
      </el-table-column>
      <el-table-column prop="min_amount" label="最低消费" width="120" />
      <el-table-column label="使用次数" width="120">
        <template #default="{ row }">
          {{ row.used_count }} / {{ row.usage_limit }}
        </template>
      </el-table-column>
      <el-table-column label="有效期" width="240">
        <template #default="{ row }">
          {{ formatDate(row.start_date) }} - {{ formatDate(row.end_date) }}
        </template>
      </el-table-column>
      <el-table-column label="状态" width="80">
        <template #default="{ row }">
          <el-tag :type="row.is_active ? 'success' : 'info'" size="small">
            {{ row.is_active ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="120" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link size="small" @click="editPromo(row)">编辑</el-button>
          <el-button type="danger" link size="small" @click="deletePromo(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div style="text-align: center; margin-top: 20px;">
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        layout="total, prev, pager, next"
        @current-change="loadPromos"
      />
    </div>

    <el-dialog v-model="showAddDialog" :title="editingPromo ? '编辑优惠码' : '添加优惠码'" width="500px">
      <el-form :model="promoForm" label-width="100px">
        <el-form-item label="名称">
          <el-input v-model="promoForm.name" />
        </el-form-item>
        <el-form-item label="优惠码">
          <el-input v-model="promoForm.code" />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="promoForm.type">
            <el-option label="百分比" value="percent" />
            <el-option label="固定金额" value="fixed" />
          </el-select>
        </el-form-item>
        <el-form-item label="折扣/金额">
          <el-input-number v-model="promoForm.value" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="最低消费">
          <el-input-number v-model="promoForm.min_amount" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="最大优惠">
          <el-input-number v-model="promoForm.max_discount" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="使用次数">
          <el-input-number v-model="promoForm.usage_limit" :min="1" />
        </el-form-item>
        <el-form-item label="开始时间">
          <el-date-picker v-model="promoForm.start_date" type="datetime" style="width: 100%" />
        </el-form-item>
        <el-form-item label="结束时间">
          <el-date-picker v-model="promoForm.end_date" type="datetime" style="width: 100%" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddDialog = false">取消</el-button>
        <el-button type="primary" @click="savePromo">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import { promoApi } from '@/api'
import type { PromoCode } from '@/types'

const promos = ref<PromoCode[]>([])
const loading = ref(false)
const showAddDialog = ref(false)
const editingPromo = ref<PromoCode | null>(null)

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const promoForm = reactive({
  name: '',
  code: '',
  type: 'percent',
  value: 10,
  min_amount: 0,
  max_discount: 0,
  usage_limit: 1,
  start_date: undefined as Date | undefined,
  end_date: undefined as Date | undefined
})

onMounted(() => {
  loadPromos()
})

const loadPromos = async () => {
  loading.value = true
  try {
    const res = await promoApi.getPromos({
      page: pagination.page,
      page_size: pagination.pageSize
    })
    promos.value = res.data.items
    pagination.total = res.data.total
  } catch {
    // ignore
  } finally {
    loading.value = false
  }
}

const editPromo = (row: PromoCode) => {
  editingPromo.value = row
  Object.assign(promoForm, {
    name: row.name,
    code: row.code,
    type: row.type,
    value: row.value,
    min_amount: row.min_amount,
    max_discount: row.max_discount,
    usage_limit: row.usage_limit,
    start_date: new Date(row.start_date),
    end_date: new Date(row.end_date)
  })
  showAddDialog.value = true
}

const savePromo = async () => {
  try {
    const data = {
      ...promoForm,
      start_date: promoForm.start_date?.toISOString(),
      end_date: promoForm.end_date?.toISOString()
    }
    if (editingPromo.value) {
      await promoApi.updatePromo(editingPromo.value.id, data)
      ElMessage.success('更新成功')
    } else {
      await promoApi.createPromo(data)
      ElMessage.success('添加成功')
    }
    showAddDialog.value = false
    editingPromo.value = null
    loadPromos()
  } catch (err: any) {
    ElMessage.error(err.message || '操作失败')
  }
}

const deletePromo = async (row: PromoCode) => {
  try {
    await ElMessageBox.confirm(`确定要删除优惠码 ${row.code} 吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await promoApi.deletePromo(row.id)
    ElMessage.success('删除成功')
    loadPromos()
  } catch (err: any) {
    if (err !== 'cancel') {
      ElMessage.error(err.message || '删除失败')
    }
  }
}

const formatDate = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}
</script>
