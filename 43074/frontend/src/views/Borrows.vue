<template>
  <div class="borrows-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <el-tabs v-model="activeTab" type="border-card">
            <el-tab-pane label="借出中" name="borrowed" />
            <el-tab-pane label="已归还" name="returned" />
            <el-tab-pane label="全部记录" name="all" />
          </el-tabs>
        </div>
      </template>

      <div v-loading="loading" class="borrows-list">
        <el-table :data="borrows" stripe>
          <el-table-column prop="id" label="ID" width="60" />
          <el-table-column label="图书">
            <template #default="{ row }">
              <div class="book-info">
                <img v-if="row.book?.cover_image" :src="row.book.cover_image" class="mini-cover" />
                <span>{{ row.book?.title || '未知图书' }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="borrower_name" label="借阅人" />
          <el-table-column prop="borrower_phone" label="联系电话" />
          <el-table-column label="借出日期">
            <template #default="{ row }">
              {{ formatDate(row.borrow_date) }}
            </template>
          </el-table-column>
          <el-table-column label="预计归还">
            <template #default="{ row }">
              <span :class="{ overdue: isOverdue(row) }">
                {{ row.expected_return_date ? formatDate(row.expected_return_date) : '-' }}
              </span>
            </template>
          </el-table-column>
          <el-table-column label="归还日期">
            <template #default="{ row }">
              {{ row.return_date ? formatDate(row.return_date) : '-' }}
            </template>
          </el-table-column>
          <el-table-column label="状态">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row)">
                {{ getStatusText(row) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="140" fixed="right">
            <template #default="{ row }">
              <el-button
                v-if="!row.return_date"
                type="primary"
                size="small"
                @click="returnBook(row)"
              >
                归还
              </el-button>
              <el-button
                type="danger"
                size="small"
                @click="deleteBorrow(row)"
              >
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>
        <div v-if="!borrows.length && !loading" class="empty">
          <el-empty description="暂无借阅记录" />
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import dayjs from 'dayjs'
import { getBorrows, returnBook as apiReturnBook, deleteBorrow as apiDeleteBorrow } from '@/api/common'
import type { BorrowRecord } from '@/types'

const activeTab = ref('borrowed')
const borrows = ref<BorrowRecord[]>([])
const loading = ref(false)

const loadBorrows = async () => {
  loading.value = true
  try {
    let status: string | undefined
    if (activeTab.value === 'borrowed') status = 'borrowed'
    else if (activeTab.value === 'returned') status = 'returned'
    borrows.value = await getBorrows(status)
  } catch (e) {
  } finally {
    loading.value = false
  }
}

const formatDate = (dateStr: string) => {
  return dayjs(dateStr).format('YYYY-MM-DD')
}

const isOverdue = (row: BorrowRecord) => {
  if (row.return_date || !row.expected_return_date) return false
  return dayjs(row.expected_return_date).isBefore(dayjs(), 'day')
}

const getStatusText = (row: BorrowRecord) => {
  if (row.return_date) return '已归还'
  if (isOverdue(row)) return '已逾期'
  return '借出中'
}

const getStatusType = (row: BorrowRecord) => {
  if (row.return_date) return 'success'
  if (isOverdue(row)) return 'danger'
  return 'warning'
}

const returnBook = async (row: BorrowRecord) => {
  try {
    await ElMessageBox.confirm(`确定「${row.borrower_name}」已归还图书吗？`, '确认归还', {
      type: 'success'
    })
    await apiReturnBook(row.id)
    ElMessage.success('归还成功')
    loadBorrows()
  } catch (e) {}
}

const deleteBorrow = async (row: BorrowRecord) => {
  try {
    await ElMessageBox.confirm('确定删除该借阅记录吗？', '确认', {
      type: 'warning'
    })
    await apiDeleteBorrow(row.id)
    ElMessage.success('删除成功')
    loadBorrows()
  } catch (e) {}
}

watch(activeTab, loadBorrows)
onMounted(() => {
  loadBorrows()
})
</script>

<style scoped lang="scss">
.borrows-page {
  .card-header {
    padding: 0;
    margin: -12px -20px -20px -20px;
  }

  .book-info {
    display: flex;
    align-items: center;
    gap: 8px;

    .mini-cover {
      width: 36px;
      height: 48px;
      object-fit: cover;
      border-radius: 2px;
    }
  }

  .overdue {
    color: #f56c6c;
    font-weight: 500;
  }
}
</style>
