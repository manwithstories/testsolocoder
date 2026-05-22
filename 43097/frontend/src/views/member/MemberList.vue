<template>
  <div class="member-list">
    <div class="header">
      <h2>会员管理</h2>
      <el-button type="primary" @click="openMemberDialog">添加会员</el-button>
    </div>

    <el-card class="filter-card">
      <el-form :inline="true" :model="filterForm" class="filter-form">
        <el-form-item label="等级">
          <el-select v-model="filterForm.level" placeholder="全部" clearable @change="fetchList">
            <el-option label="普通会员" :value="MemberLevel.NORMAL" />
            <el-option label="银卡会员" :value="MemberLevel.SILVER" />
            <el-option label="金卡会员" :value="MemberLevel.GOLD" />
            <el-option label="白金会员" :value="MemberLevel.PLATINUM" />
            <el-option label="钻石会员" :value="MemberLevel.DIAMOND" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filterForm.status" placeholder="全部" clearable @change="fetchList">
            <el-option label="正常" :value="true" />
            <el-option label="禁用" :value="false" />
          </el-select>
        </el-form-item>
        <el-form-item label="关键词">
          <el-input v-model="filterForm.keyword" placeholder="会员号/姓名/电话" @keyup.enter="fetchList" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchList">查询</el-button>
          <el-button @click="resetFilter">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="table-card">
      <el-table :data="tableData" v-loading="loading" border stripe>
        <el-table-column prop="memberNo" label="会员号" width="140" />
        <el-table-column prop="name" label="姓名" width="100" />
        <el-table-column prop="phone" label="电话" width="130" />
        <el-table-column prop="level" label="等级" width="100">
          <template #default="{ row }">
            <el-tag :type="getLevelTagType(row.level)">
              {{ getLevelText(row.level) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="points" label="积分" width="100" />
        <el-table-column prop="balance" label="余额" width="100">
          <template #default="{ row }">
            ¥{{ row.balance }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status ? 'success' : 'danger'">
              {{ row.status ? '正常' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="注册时间" width="180">
          <template #default="{ row }">
            {{ formatDateTime(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="primary" link @click="viewDetail(row)">详情</el-button>
            <el-button size="small" type="success" link @click="openEditDialog(row)">编辑</el-button>
            <el-button size="small" type="warning" link @click="openPointsDialog(row)">积分充值</el-button>
            <el-button
              size="small"
              :type="row.status ? 'info' : 'success'"
              link
              @click="toggleStatus(row)"
            >
              {{ row.status ? '禁用' : '启用' }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <el-dialog v-model="memberDialogVisible" :title="isEdit ? '编辑会员' : '添加会员'" width="600px">
      <el-form :model="memberForm" :rules="memberRules" ref="memberFormRef" label-width="100px">
        <el-form-item label="姓名" prop="name">
          <el-input v-model="memberForm.name" placeholder="请输入姓名" />
        </el-form-item>
        <el-form-item label="电话" prop="phone">
          <el-input v-model="memberForm.phone" placeholder="请输入电话" />
        </el-form-item>
        <el-form-item label="身份证">
          <el-input v-model="memberForm.idCard" placeholder="请输入身份证号" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="memberForm.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="生日">
          <el-date-picker
            v-model="memberForm.birthday"
            type="date"
            placeholder="选择生日"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item label="地址">
          <el-input v-model="memberForm.address" type="textarea" :rows="2" placeholder="请输入地址" />
        </el-form-item>
        <el-form-item v-if="!isEdit" label="初始积分">
          <el-input-number v-model="memberForm.points" :min="0" />
        </el-form-item>
        <el-form-item v-if="!isEdit" label="初始余额">
          <el-input-number v-model="memberForm.balance" :min="0" :precision="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="memberDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleMemberSubmit">确认</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="pointsDialogVisible" title="积分充值" width="400px">
      <el-form :model="pointsForm" label-width="100px">
        <el-form-item label="会员姓名">
          <el-input :value="currentMember?.name" disabled />
        </el-form-item>
        <el-form-item label="当前积分">
          <el-input :value="currentMember?.points" disabled />
        </el-form-item>
        <el-form-item label="充值积分">
          <el-input-number v-model="pointsForm.points" :min="1" />
        </el-form-item>
        <el-form-item label="说明">
          <el-input v-model="pointsForm.description" type="textarea" :rows="2" placeholder="请输入说明" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="pointsDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handlePointsRecharge">确认充值</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="detailDialogVisible" title="会员详情" width="700px">
      <el-tabs v-if="currentMember" v-model="activeTab">
        <el-tab-pane label="基本信息" name="basic">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="会员号">
              {{ currentMember.memberNo }}
            </el-descriptions-item>
            <el-descriptions-item label="姓名">
              {{ currentMember.name }}
            </el-descriptions-item>
            <el-descriptions-item label="电话">
              {{ currentMember.phone }}
            </el-descriptions-item>
            <el-descriptions-item label="身份证">
              {{ currentMember.idCard || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="邮箱">
              {{ currentMember.email || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="生日">
              {{ currentMember.birthday || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="等级">
              <el-tag :type="getLevelTagType(currentMember.level)">
                {{ getLevelText(currentMember.level) }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="currentMember.status ? 'success' : 'danger'">
                {{ currentMember.status ? '正常' : '禁用' }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="积分">
              {{ currentMember.points }}
            </el-descriptions-item>
            <el-descriptions-item label="余额">
              ¥{{ currentMember.balance }}
            </el-descriptions-item>
            <el-descriptions-item label="累计消费">
              ¥{{ currentMember.totalSpent }}
            </el-descriptions-item>
            <el-descriptions-item label="入住次数">
              {{ currentMember.totalStays }} 次
            </el-descriptions-item>
            <el-descriptions-item label="注册时间">
              {{ formatDateTime(currentMember.createdAt) }}
            </el-descriptions-item>
            <el-descriptions-item label="最后访问">
              {{ currentMember.lastVisitAt ? formatDateTime(currentMember.lastVisitAt) : '-' }}
            </el-descriptions-item>
          </el-descriptions>
        </el-tab-pane>
        <el-tab-pane label="积分流水" name="points">
          <el-table :data="pointsLog" border v-loading="detailLoading">
            <el-table-column prop="type" label="类型" width="100">
              <template #default="{ row }">
                <el-tag :type="row.type === 'add' ? 'success' : 'danger'">
                  {{ row.type === 'add' ? '增加' : '扣减' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="points" label="积分" width="100">
              <template #default="{ row }">
                <span :class="row.type === 'add' ? 'text-success' : 'text-danger'">
                  {{ row.type === 'add' ? '+' : '-' }}{{ row.points }}
                </span>
              </template>
            </el-table-column>
            <el-table-column prop="description" label="说明" />
            <el-table-column prop="createdAt" label="时间" width="180">
              <template #default="{ row }">
                {{ formatDateTime(row.createdAt) }}
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="消费记录" name="consumption">
          <el-table :data="consumptionHistory" border v-loading="detailLoading">
            <el-table-column prop="orderNo" label="订单号" width="180" />
            <el-table-column prop="type" label="类型" width="100">
              <template #default="{ row }">
                {{ row.type === 'checkin' ? '入住' : '消费' }}
              </template>
            </el-table-column>
            <el-table-column prop="amount" label="金额" width="120">
              <template #default="{ row }">
                ¥{{ row.amount }}
              </template>
            </el-table-column>
            <el-table-column prop="pointsEarned" label="获得积分" width="100">
              <template #default="{ row }">
                +{{ row.pointsEarned || 0 }}
              </template>
            </el-table-column>
            <el-table-column prop="createdAt" label="时间" width="180">
              <template #default="{ row }">
                {{ formatDateTime(row.createdAt) }}
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { MemberLevel, type Member } from '@/types'
import {
  getMemberList,
  createMember,
  updateMember,
  updateMemberStatus,
  addMemberPoints,
  getMemberConsumptionHistory
} from '@/api/member'

const loading = ref(false)
const submitting = ref(false)
const detailLoading = ref(false)
const tableData = ref<Member[]>([])
const memberFormRef = ref<FormInstance>()

const filterForm = reactive({
  level: '',
  status: '',
  keyword: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const memberDialogVisible = ref(false)
const pointsDialogVisible = ref(false)
const detailDialogVisible = ref(false)
const isEdit = ref(false)
const currentMember = ref<Member | null>(null)
const activeTab = ref('basic')
const pointsLog = ref<any[]>([])
const consumptionHistory = ref<any[]>([])

const memberForm = reactive({
  id: null as number | null,
  name: '',
  phone: '',
  idCard: '',
  email: '',
  birthday: '',
  address: '',
  points: 0,
  balance: 0,
  level: MemberLevel.NORMAL,
  status: true
})

const pointsForm = reactive({
  points: 0,
  description: ''
})

const memberRules: FormRules = {
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  phone: [{ required: true, message: '请输入电话', trigger: 'blur' }]
}

const fetchList = async () => {
  loading.value = true
  try {
    const res = await getMemberList({
      page: pagination.page,
      pageSize: pagination.pageSize,
      keyword: filterForm.keyword || undefined,
      level: filterForm.level as MemberLevel || undefined,
      status: filterForm.status !== '' ? (filterForm.status as unknown as boolean) : undefined
    })
    tableData.value = res.list
    pagination.total = res.total
  } finally {
    loading.value = false
  }
}

const resetFilter = () => {
  filterForm.level = ''
  filterForm.status = ''
  filterForm.keyword = ''
  pagination.page = 1
  fetchList()
}

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  fetchList()
}

const handleCurrentChange = (page: number) => {
  pagination.page = page
  fetchList()
}

const formatDateTime = (date: string) => {
  return date ? new Date(date).toLocaleString('zh-CN') : '-'
}

const getLevelTagType = (level: string) => {
  switch (level) {
    case MemberLevel.NORMAL:
      return 'info'
    case MemberLevel.SILVER:
      return ''
    case MemberLevel.GOLD:
      return 'warning'
    case MemberLevel.PLATINUM:
      return ''
    case MemberLevel.DIAMOND:
      return 'success'
    default:
      return 'info'
  }
}

const getLevelText = (level: string) => {
  switch (level) {
    case MemberLevel.NORMAL:
      return '普通会员'
    case MemberLevel.SILVER:
      return '银卡会员'
    case MemberLevel.GOLD:
      return '金卡会员'
    case MemberLevel.PLATINUM:
      return '白金会员'
    case MemberLevel.DIAMOND:
      return '钻石会员'
    default:
      return level
  }
}

const openMemberDialog = () => {
  isEdit.value = false
  Object.assign(memberForm, {
    id: null,
    name: '',
    phone: '',
    idCard: '',
    email: '',
    birthday: '',
    address: '',
    points: 0,
    balance: 0,
    level: MemberLevel.NORMAL,
    status: true
  })
  memberDialogVisible.value = true
}

const openEditDialog = (row: Member) => {
  isEdit.value = true
  Object.assign(memberForm, {
    id: row.id,
    name: row.name,
    phone: row.phone,
    idCard: row.idCard,
    email: row.email || '',
    birthday: row.birthday || '',
    address: row.address || '',
    points: row.points,
    balance: row.balance,
    level: row.level,
    status: row.status
  })
  memberDialogVisible.value = true
}

const handleMemberSubmit = async () => {
  if (!memberFormRef.value) return
  await memberFormRef.value.validate(async (valid) => {
    if (!valid) return
    submitting.value = true
    try {
      if (isEdit.value && memberForm.id) {
        await updateMember(memberForm.id, {
          name: memberForm.name,
          phone: memberForm.phone,
          idCard: memberForm.idCard,
          email: memberForm.email,
          birthday: memberForm.birthday,
          address: memberForm.address
        })
        ElMessage.success('编辑成功')
      } else {
        await createMember({
          name: memberForm.name,
          phone: memberForm.phone,
          idCard: memberForm.idCard,
          email: memberForm.email,
          birthday: memberForm.birthday,
          address: memberForm.address,
          level: memberForm.level,
          points: memberForm.points,
          balance: memberForm.balance,
          status: memberForm.status
        })
        ElMessage.success('添加成功')
      }
      memberDialogVisible.value = false
      fetchList()
    } catch (e: any) {
      ElMessage.error(e.message || '操作失败')
    } finally {
      submitting.value = false
    }
  })
}

const openPointsDialog = (row: Member) => {
  currentMember.value = row
  pointsForm.points = 0
  pointsForm.description = ''
  pointsDialogVisible.value = true
}

const handlePointsRecharge = async () => {
  if (!currentMember.value) return
  if (!pointsForm.points) {
    ElMessage.warning('请输入充值积分')
    return
  }
  submitting.value = true
  try {
    await addMemberPoints(currentMember.value.id, pointsForm.points)
    ElMessage.success('积分充值成功')
    pointsDialogVisible.value = false
    fetchList()
  } catch (e: any) {
    ElMessage.error(e.message || '充值失败')
  } finally {
    submitting.value = false
  }
}

const toggleStatus = (row: Member) => {
  const action = row.status ? '禁用' : '启用'
  ElMessageBox.confirm(`确定要${action}该会员吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await updateMemberStatus(row.id, !row.status)
      ElMessage.success(`${action}成功`)
      fetchList()
    } catch (e: any) {
      ElMessage.error(e.message || `${action}失败`)
    }
  })
}

const viewDetail = async (row: Member) => {
  currentMember.value = row
  activeTab.value = 'basic'
  detailDialogVisible.value = true
  detailLoading.value = true
  try {
    const res = await getMemberConsumptionHistory(row.id, { page: 1, pageSize: 20 })
    consumptionHistory.value = res.list || []
  } finally {
    detailLoading.value = false
  }
  pointsLog.value = [
    { type: 'add', points: 100, description: '注册赠送', createdAt: '2024-01-01 10:00:00' },
    { type: 'add', points: 50, description: '消费获得', createdAt: '2024-01-15 14:30:00' },
    { type: 'deduct', points: 30, description: '积分兑换', createdAt: '2024-01-20 09:00:00' }
  ]
}

onMounted(() => {
  fetchList()
})
</script>

<style scoped>
.member-list {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.header h2 {
  margin: 0;
}

.filter-card {
  margin-bottom: 20px;
}

.filter-form {
  margin: 0;
}

.table-card {
  margin-bottom: 20px;
}

.pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

.text-success {
  color: #67c23a;
}

.text-danger {
  color: #f56c6c;
}
</style>
