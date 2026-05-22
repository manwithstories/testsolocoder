<template>
  <div class="coupon-list">
    <div class="page-header">
      <h2 class="page-title">优惠券管理</h2>
      <el-button type="primary" @click="dialogVisible = true" v-if="userStore.isAdmin">
        <el-icon><Plus /></el-icon>
        创建优惠券
      </el-button>
    </div>

    <el-card>
      <div class="search-bar">
        <el-input v-model="search.code" placeholder="优惠券码" clearable style="width: 200px" />
        <el-select v-model="search.status" placeholder="状态" clearable style="width: 140px">
          <el-option label="有效" value="active" />
          <el-option label="已用完" value="used" />
          <el-option label="已过期" value="expired" />
        </el-select>
        <el-select v-model="search.type" placeholder="类型" clearable style="width: 140px">
          <el-option label="固定金额" value="fixed" />
          <el-option label="折扣" value="discount" />
        </el-select>
        <el-button type="primary" @click="loadData">搜索</el-button>
      </div>

      <el-table :data="list" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="code" label="优惠券码" width="140" />
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="row.type === 'fixed' ? '' : 'warning'">
              {{ row.type === 'fixed' ? '固定金额' : '折扣' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="面值" width="120">
          <template #default="{ row }">
            {{ row.type === 'fixed' ? '¥' + row.value : row.value + '折' }}
          </template>
        </el-table-column>
        <el-table-column prop="minAmount" label="最低消费(元)" width="140" />
        <el-table-column label="使用情况" width="140">
          <template #default="{ row }">{{ row.usedCount }} / {{ row.totalCount }}</template>
        </el-table-column>
        <el-table-column label="有效期" min-width="240">
          <template #default="{ row }">
            <div>{{ formatDate(row.startTime) }}</div>
            <div class="sub-text">至 {{ formatDate(row.endTime) }}</div>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100" v-if="userStore.isAdmin">
          <template #default="{ row }">
            <el-button size="small" type="danger" @click="handleDelete(row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        layout="total, prev, pager, next, jumper"
        @current-change="loadData"
        style="margin-top: 20px; justify-content: flex-end"
      />
    </el-card>

    <el-dialog v-model="dialogVisible" title="创建优惠券" width="500px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="类型" prop="type">
          <el-select v-model="form.type" placeholder="选择类型" style="width: 100%">
            <el-option label="固定金额" value="fixed" />
            <el-option label="折扣" value="discount" />
          </el-select>
        </el-form-item>
        <el-form-item label="面值" prop="value">
          <el-input-number v-model="form.value" :min="0" :max="form.type === 'discount' ? 10 : undefined" :precision="2" style="width: 100%" />
          <span class="tip">{{ form.type === 'fixed' ? '元' : '折' }}</span>
        </el-form-item>
        <el-form-item label="最低消费" prop="minAmount">
          <el-input-number v-model="form.minAmount" :min="0" :precision="2" style="width: 100%" />
        </el-form-item>
        <el-form-item label="发放数量" prop="totalCount">
          <el-input-number v-model="form.totalCount" :min="1" style="width: 100%" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="开始时间" prop="startTime">
              <el-date-picker
                v-model="form.startTime"
                type="datetime"
                placeholder="选择开始时间"
                value-format="YYYY-MM-DD HH:mm:ss"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="结束时间" prop="endTime">
              <el-date-picker
                v-model="form.endTime"
                type="datetime"
                placeholder="选择结束时间"
                value-format="YYYY-MM-DD HH:mm:ss"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { getCouponList, createCoupon, deleteCoupon } from '@/api/coupon'
import { useUserStore } from '@/store/user'
import dayjs from 'dayjs'

const userStore = useUserStore()

const loading = ref(false)
const submitting = ref(false)
const list = ref<any[]>([])
const dialogVisible = ref(false)
const formRef = ref<FormInstance>()

const search = reactive({
  code: '',
  status: '',
  type: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const form = reactive({
  type: 'fixed',
  value: 0,
  minAmount: 0,
  totalCount: 100,
  startTime: '',
  endTime: ''
})

const rules: FormRules = {
  type: [{ required: true, message: '请选择类型', trigger: 'change' }],
  value: [{ required: true, message: '请输入面值', trigger: 'blur' }],
  totalCount: [{ required: true, message: '请输入发放数量', trigger: 'blur' }]
}

const loadData = async () => {
  try {
    loading.value = true
    const res = await getCouponList({
      page: pagination.page,
      pageSize: pagination.pageSize,
      code: search.code,
      status: search.status,
      type: search.type
    })
    list.value = res.list
    pagination.total = res.total
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleDelete = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定要删除这个优惠券吗？', '提示', { type: 'warning' })
    await deleteCoupon(id)
    ElMessage.success('删除成功')
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error(error)
    }
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    try {
      submitting.value = true
      await createCoupon(form as any)
      ElMessage.success('创建成功')
      dialogVisible.value = false
      loadData()
    } catch (error) {
      console.error(error)
    } finally {
      submitting.value = false
    }
  })
}

const getStatusType = (status: string) => {
  const map: Record<string, string> = { active: 'success', used: 'info', expired: 'danger' }
  return map[status] || 'info'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = { active: '有效', used: '已用完', expired: '已过期' }
  return map[status] || status
}

const formatDate = (date: string) => date ? dayjs(date).format('YYYY-MM-DD HH:mm') : '长期有效'

onMounted(loadData)
</script>

<style scoped lang="scss">
.sub-text {
  font-size: 12px;
  color: #909399;
  margin-top: 2px;
}

.tip {
  margin-left: 8px;
  color: #909399;
}
</style>
