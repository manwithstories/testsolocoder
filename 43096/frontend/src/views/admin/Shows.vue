<template>
  <div class="admin-shows">
    <el-card class="filter-card">
      <el-form :inline="true" :model="filterForm">
        <el-form-item label="演出名称">
          <el-input v-model="filterForm.keyword" placeholder="请输入演出名称" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filterForm.status" placeholder="全部" clearable style="width: 150px;">
            <el-option label="未发布" :value="0" />
            <el-option label="已发布" :value="1" />
            <el-option label="已结束" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchShows">查询</el-button>
          <el-button @click="resetFilter">重置</el-button>
          <el-button type="success" @click="showCreateDialog = true">+ 新增演出</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="table-card">
      <el-table :data="shows" style="width: 100%" v-loading="loading">
        <el-table-column prop="name" label="演出名称" min-width="200" />
        <el-table-column prop="artist" label="艺人" width="150" />
        <el-table-column prop="venue" label="场馆" width="200" />
        <el-table-column label="场次">
          <template #default="{ row }">
            {{ row.sessions?.length || 0 }} 场
          </template>
        </el-table-column>
        <el-table-column prop="duration" label="时长(分钟)" width="120" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusTagType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="240" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" type="primary" @click="handleManageSessions(row)">场次管理</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        class="pagination"
        @current-change="fetchShows"
      />
    </el-card>

    <el-dialog v-model="showCreateDialog" :title="isEdit ? '编辑演出' : '新增演出'" width="600px" @close="resetForm">
      <el-form :model="showForm" :rules="showRules" ref="showFormRef" label-width="100px">
        <el-form-item label="演出名称" prop="name">
          <el-input v-model="showForm.name" />
        </el-form-item>
        <el-form-item label="演出描述" prop="description">
          <el-input v-model="showForm.description" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item label="艺人" prop="artist">
          <el-input v-model="showForm.artist" />
        </el-form-item>
        <el-form-item label="时长(分钟)" prop="duration">
          <el-input-number v-model="showForm.duration" :min="30" :max="600" />
        </el-form-item>
        <el-form-item label="场馆" prop="venue">
          <el-input v-model="showForm.venue" />
        </el-form-item>
        <el-form-item label="地址" prop="address">
          <el-input v-model="showForm.address" />
        </el-form-item>
        <el-form-item label="主办方" prop="organizer">
          <el-input v-model="showForm.organizer" />
        </el-form-item>
        <el-form-item label="海报" prop="poster">
          <el-input v-model="showForm.poster" placeholder="海报URL" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="showForm.status">
            <el-radio :value="0">未发布</el-radio>
            <el-radio :value="1">已发布</el-radio>
            <el-radio :value="2">已结束</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确认</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showSessionDialog" title="场次管理" width="800px">
      <template v-if="currentShow">
        <div class="session-header">
          <h3>{{ currentShow.name }}</h3>
          <el-button type="primary" size="small" @click="showAddSession = true">+ 新增场次</el-button>
        </div>

        <el-table :data="sessions" style="width: 100%; margin-top: 16px;">
          <el-table-column prop="start_time" label="开始时间" width="180">
            <template #default="{ row }">
              {{ formatTime(row.start_time) }}
            </template>
          </el-table-column>
          <el-table-column prop="end_time" label="结束时间" width="180">
            <template #default="{ row }">
              {{ formatTime(row.end_time) }}
            </template>
          </el-table-column>
          <el-table-column prop="total_seats" label="总座位数" width="120" />
          <el-table-column prop="sold_seats" label="已售" width="100" />
          <el-table-column label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="row.status === 1 ? 'success' : 'info'">
                {{ row.status === 1 ? '销售中' : '已结束' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="200">
            <template #default="{ row }">
              <el-button size="small" @click="handleManageSeats(row)">座位管理</el-button>
              <el-button size="small" type="danger" @click="handleDeleteSession(row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </template>
    </el-dialog>

    <el-dialog v-model="showAddSession" title="新增场次" width="500px">
      <el-form :model="sessionForm" :rules="sessionRules" ref="sessionFormRef" label-width="100px">
        <el-form-item label="开始时间" prop="start_time">
          <el-date-picker
            v-model="sessionForm.start_time"
            type="datetime"
            placeholder="选择开始时间"
            style="width: 100%;"
          />
        </el-form-item>
        <el-form-item label="结束时间" prop="end_time">
          <el-date-picker
            v-model="sessionForm.end_time"
            type="datetime"
            placeholder="选择结束时间"
            style="width: 100%;"
          />
        </el-form-item>
        <el-form-item label="总座位数" prop="total_seats">
          <el-input-number v-model="sessionForm.total_seats" :min="1" :max="10000" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddSession = false">取消</el-button>
        <el-button type="primary" @click="handleAddSession">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { ElMessage, ElMessageBox, ElForm } from 'element-plus'
import dayjs from 'dayjs'
import { showApi } from '@/api'
import type { Show, Session } from '@/types'

const loading = ref(false)
const shows = ref<Show[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

const filterForm = reactive({
  keyword: '',
  status: undefined as number | undefined
})

const showCreateDialog = ref(false)
const isEdit = ref(false)
const editId = ref(0)
const showForm = reactive({
  name: '',
  description: '',
  artist: '',
  duration: 120,
  venue: '',
  address: '',
  organizer: '',
  poster: '',
  status: 0
})

const showRules = {
  name: [{ required: true, message: '请输入演出名称', trigger: 'blur' }],
  venue: [{ required: true, message: '请输入场馆', trigger: 'blur' }]
}

const showSessionDialog = ref(false)
const currentShow = ref<Show | null>(null)
const sessions = ref<Session[]>([])
const showAddSession = ref(false)
const sessionForm = reactive({
  start_time: '',
  end_time: '',
  total_seats: 500
})
const sessionRules = {
  start_time: [{ required: true, message: '请选择开始时间', trigger: 'change' }],
  end_time: [{ required: true, message: '请选择结束时间', trigger: 'change' }],
  total_seats: [{ required: true, message: '请输入总座位数', trigger: 'blur' }]
}

const showFormRef = ref<InstanceType<typeof ElForm>>()
const sessionFormRef = ref<InstanceType<typeof ElForm>>()

async function fetchShows() {
  loading.value = true
  try {
    const res = await showApi.list({
      page: page.value,
      page_size: pageSize.value,
      keyword: filterForm.keyword,
      status: filterForm.status
    })
    shows.value = res.list
    total.value = res.pagination?.total || 0
  } catch (err) {
    console.error(err)
  } finally {
    loading.value = false
  }
}

function resetFilter() {
  filterForm.keyword = ''
  filterForm.status = undefined
  page.value = 1
  fetchShows()
}

function resetForm() {
  showForm.name = ''
  showForm.description = ''
  showForm.artist = ''
  showForm.duration = 120
  showForm.venue = ''
  showForm.address = ''
  showForm.organizer = ''
  showForm.poster = ''
  showForm.status = 0
  isEdit.value = false
  editId.value = 0
}

async function handleSubmit() {
  try {
    await showFormRef?.validate()
    if (isEdit.value) {
      await showApi.update(editId.value, showForm)
      ElMessage.success('更新成功')
    } else {
      await showApi.create(showForm)
      ElMessage.success('创建成功')
    }
    showCreateDialog.value = false
    fetchShows()
  } catch (err) {
    console.error(err)
  }
}

function handleEdit(row: Show) {
  isEdit.value = true
  editId.value = row.id
  showForm.name = row.name
  showForm.description = row.description
  showForm.artist = row.artist
  showForm.duration = row.duration
  showForm.venue = row.venue
  showForm.address = row.address
  showForm.organizer = row.organizer
  showForm.poster = row.poster
  showForm.status = row.status
  showCreateDialog.value = true
}

async function handleDelete(row: Show) {
  try {
    await ElMessageBox.confirm('确定要删除该演出吗？', '提示', {
      type: 'warning'
    })
    await showApi.delete(row.id)
    ElMessage.success('删除成功')
    fetchShows()
  } catch (err) {
    if (err !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

async function handleManageSessions(row: Show) {
  currentShow.value = row
  try {
    const res = await showApi.getSessions(row.id)
    sessions.value = res
    showSessionDialog.value = true
  } catch (err) {
    console.error(err)
  }
}

async function handleAddSession() {
  try {
    await sessionFormRef?.validate()
    await showApi.createSession({
      show_id: currentShow.value!.id,
      start_time: dayjs(sessionForm.start_time).format('YYYY-MM-DD HH:mm:ss'),
      end_time: dayjs(sessionForm.end_time).format('YYYY-MM-DD HH:mm:ss'),
      total_seats: sessionForm.total_seats
    })
    ElMessage.success('场次创建成功')
    showAddSession.value = false
    handleManageSessions(currentShow.value!)
  } catch (err) {
    console.error(err)
  }
}

function handleManageSeats(row: Session) {
  ElMessage.info(`座位管理功能开发中，场次ID: ${row.id}`)
}

async function handleDeleteSession(row: Session) {
  try {
    await ElMessageBox.confirm('确定要删除该场次吗？', '提示', {
      type: 'warning'
    })
    ElMessage.success('删除成功')
    handleManageSessions(currentShow.value!)
  } catch (err) {
    if (err !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

function getStatusText(status: number) {
  const texts: Record<number, string> = {
    0: '未发布',
    1: '已发布',
    2: '已结束'
  }
  return texts[status] || '未知'
}

function getStatusTagType(status: number) {
  const types: Record<number, string> = {
    0: 'info',
    1: 'success',
    2: 'danger'
  }
  return types[status] || 'info'
}

function formatTime(time: string) {
  return dayjs(time).format('YYYY-MM-DD HH:mm')
}

fetchShows()
</script>

<style lang="scss" scoped>
.admin-shows {
  .filter-card {
    margin-bottom: 20px;
  }

  .table-card {
    .pagination {
      margin-top: 20px;
      justify-content: flex-end;
    }
  }

  .session-header {
    display: flex;
    justify-content: space-between;
    align-items: center;

    h3 {
      margin: 0;
    }
  }
}
</style>
