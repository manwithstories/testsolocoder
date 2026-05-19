<template>
  <div v-loading="loading" class="book-detail">
    <div v-if="book" class="detail-content">
      <el-button link @click="goBack" style="margin-bottom: 16px">
        <el-icon><ArrowLeft /></el-icon>
        返回列表
      </el-button>

      <el-row :gutter="24">
        <el-col :span="6">
          <el-card class="cover-card">
            <div class="cover-container">
              <img v-if="book.cover_image" :src="book.cover_image" class="large-cover" />
              <div v-else class="large-cover placeholder">
                <el-icon :size="64"><Reading /></el-icon>
              </div>
            </div>
            <el-upload
              :show-file-list="false"
              :before-upload="beforeCoverUpload"
              :http-request="handleCoverUpload"
              accept="image/*"
            >
              <el-button type="primary" style="width: 100%; margin-top: 12px">
                <el-icon><Upload /></el-icon>
                上传封面
              </el-button>
            </el-upload>
          </el-card>
        </el-col>

        <el-col :span="18">
          <el-card class="info-card">
            <template #header>
              <div class="card-header">
                <h2>{{ book.title }}</h2>
                <div>
                  <el-tag :type="statusType[book.reading_status]">{{ statusText[book.reading_status] }}</el-tag>
                </div>
              </div>
            </template>

            <el-descriptions :column="2" border>
              <el-descriptions-item label="作者">{{ book.author || '-' }}</el-descriptions-item>
              <el-descriptions-item label="出版社">{{ book.publisher || '-' }}</el-descriptions-item>
              <el-descriptions-item label="ISBN">{{ book.isbn || '-' }}</el-descriptions-item>
              <el-descriptions-item label="总页数">{{ book.total_pages || '-' }}</el-descriptions-item>
              <el-descriptions-item label="开始阅读">
                {{ book.start_date ? formatDate(book.start_date) : '-' }}
              </el-descriptions-item>
              <el-descriptions-item label="读完日期">
                {{ book.end_date ? formatDate(book.end_date) : '-' }}
              </el-descriptions-item>
              <el-descriptions-item label="标签">
                <el-tag
                  v-for="tag in book.tags"
                  :key="tag.id"
                  style="margin-right: 8px"
                  :style="{ backgroundColor: tag.color + '20', color: tag.color, borderColor: tag.color + '40' }"
                >
                  {{ tag.name }}
                </el-tag>
              </el-descriptions-item>
              <el-descriptions-item label="分类">
                <span v-for="cat in book.categories" :key="cat.id" style="margin-right: 8px">
                  {{ cat.name }}
                </span>
              </el-descriptions-item>
              <el-descriptions-item label="简介" :span="2">
                {{ book.summary || '暂无简介' }}
              </el-descriptions-item>
            </el-descriptions>

            <div class="action-section">
              <h3>阅读进度</h3>
              <div class="progress-section">
                <div class="progress-display">
                  <span class="current-page">{{ book.current_page }}</span>
                  <span class="divider">/</span>
                  <span class="total-pages">{{ book.total_pages }}</span>
                </div>
                <el-progress
                  :percentage="Math.round(book.reading_progress)"
                  :stroke-width="12"
                />
                <div class="progress-actions">
                  <el-input-number
                    v-model="newPage"
                    :min="0"
                    :max="book.total_pages || 9999"
                    size="small"
                    placeholder="当前页码"
                  />
                  <el-button type="primary" size="small" @click="updateProgress">更新进度</el-button>
                  <el-button
                    v-if="book.reading_status !== 'completed'"
                    size="small"
                    @click="markAsCompleted"
                  >
                    标记为已读
                  </el-button>
                </div>
              </div>
            </div>

            <div class="action-section">
              <h3>阅读笔记</h3>
              <div class="note-input">
                <el-input
                  v-model="newNoteContent"
                  type="textarea"
                  :rows="2"
                  placeholder="写下你的读书笔记..."
                />
                <div class="note-actions">
                  <el-input-number
                    v-model="newNotePage"
                    :min="0"
                    :max="book.total_pages || 9999"
                    size="small"
                    placeholder="页码"
                  />
                  <el-button type="primary" size="small" @click="addNote">添加笔记</el-button>
                </div>
              </div>
              <div class="notes-list">
                <div v-if="!book.reading_notes?.length" class="empty-notes">
                  暂无笔记
                </div>
                <div v-for="note in book.reading_notes" :key="note.id" class="note-item">
                  <div class="note-header">
                    <el-tag size="small" type="info" v-if="note.page">P.{{ note.page }}</el-tag>
                    <span class="note-date">{{ formatDate(note.created_at) }}</span>
                    <el-button type="danger" link size="small" @click="deleteNote(note.id)">删除</el-button>
                  </div>
                  <div class="note-content">{{ note.content }}</div>
                </div>
              </div>
            </div>

            <div class="action-section">
              <h3>借阅管理</h3>
              <div v-if="!book.borrow_record" class="borrow-form">
                <el-input v-model="newBorrow.name" placeholder="借阅人姓名" style="width: 200px" />
                <el-input v-model="newBorrow.phone" placeholder="联系电话" style="width: 200px" />
                <el-date-picker
                  v-model="newBorrow.expectedReturn"
                  type="date"
                  placeholder="预计归还日期"
                  value-format="YYYY-MM-DD"
                />
                <el-button type="primary" @click="borrowBook">借出</el-button>
              </div>
              <div v-else class="borrow-info">
                <el-alert
                  type="warning"
                  :title="`当前借阅人：${book.borrow_record.borrower_name}`"
                  show-icon
                >
                  <template #default>
                    <div>借出日期：{{ formatDate(book.borrow_record.borrow_date) }}</div>
                    <div v-if="book.borrow_record.expected_return_date">
                      预计归还：{{ formatDate(book.borrow_record.expected_return_date) }}
                    </div>
                    <div v-if="book.borrow_record.borrower_phone">
                      联系方式：{{ book.borrow_record.borrower_phone }}
                    </div>
                    <el-button type="primary" size="small" style="margin-top: 8px" @click="returnBook">
                      归还
                    </el-button>
                  </template>
                </el-alert>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import dayjs from 'dayjs'
import { getBook, updateReadingProgress, updateBookStatus, uploadCover } from '@/api/book'
import { createNote, deleteNote as apiDeleteNote, createBorrow, returnBook as apiReturnBook } from '@/api/common'
import type { Book, ReadingStatus } from '@/types'

const route = useRoute()
const router = useRouter()

const book = ref<Book | null>(null)
const loading = ref(false)
const newPage = ref(0)
const newNoteContent = ref('')
const newNotePage = ref(0)

const newBorrow = reactive({
  name: '',
  phone: '',
  expectedReturn: ''
})

const statusText: Record<ReadingStatus, string> = {
  to_read: '想读',
  reading: '在读',
  completed: '已读',
  abandoned: '放弃'
}

const statusType: Record<ReadingStatus, string> = {
  to_read: 'info',
  reading: 'primary',
  completed: 'success',
  abandoned: 'danger'
}

const loadBook = async () => {
  const id = Number(route.params.id)
  if (!id) return
  loading.value = true
  try {
    book.value = await getBook(id)
    newPage.value = book.value.current_page
  } catch (e) {
  } finally {
    loading.value = false
  }
}

const formatDate = (dateStr: string) => {
  return dayjs(dateStr).format('YYYY-MM-DD')
}

const updateProgress = async () => {
  if (!book.value) return
  try {
    const res = await updateReadingProgress(book.value.id, newPage.value)
    book.value.current_page = res.current_page
    book.value.reading_progress = res.reading_progress
    book.value.reading_status = res.reading_status as ReadingStatus
    ElMessage.success('进度已更新')
  } catch (e) {}
}

const markAsCompleted = async () => {
  if (!book.value) return
  try {
    const res = await updateBookStatus(book.value.id, 'completed')
    book.value.reading_status = res.reading_status as ReadingStatus
    book.value.reading_progress = 100
    if (book.value.total_pages) {
      book.value.current_page = book.value.total_pages
      newPage.value = book.value.total_pages
    }
    ElMessage.success('已标记为已读')
  } catch (e) {}
}

const addNote = async () => {
  if (!book.value || !newNoteContent.value.trim()) return
  try {
    const note = await createNote({
      book_id: book.value.id,
      content: newNoteContent.value.trim(),
      page: newNotePage.value || undefined
    })
    if (!book.value.reading_notes) book.value.reading_notes = []
    book.value.reading_notes.unshift(note)
    newNoteContent.value = ''
    newNotePage.value = 0
    ElMessage.success('笔记已添加')
  } catch (e) {}
}

const deleteNote = async (noteId: number) => {
  try {
    await apiDeleteNote(noteId)
    if (book.value?.reading_notes) {
      book.value.reading_notes = book.value.reading_notes.filter(n => n.id !== noteId)
    }
    ElMessage.success('笔记已删除')
  } catch (e) {}
}

const borrowBook = async () => {
  if (!book.value || !newBorrow.name.trim()) {
    ElMessage.warning('请输入借阅人姓名')
    return
  }
  try {
    const record = await createBorrow({
      book_id: book.value.id,
      borrower_name: newBorrow.name.trim(),
      borrower_phone: newBorrow.phone.trim(),
      expected_return_date: newBorrow.expectedReturn || undefined
    })
    book.value.borrow_record = record
    newBorrow.name = ''
    newBorrow.phone = ''
    newBorrow.expectedReturn = ''
    ElMessage.success('借出成功')
  } catch (e) {}
}

const returnBook = async () => {
  if (!book.value?.borrow_record) return
  try {
    const record = await apiReturnBook(book.value.borrow_record.id)
    book.value.borrow_record = undefined
    ElMessage.success('归还成功')
  } catch (e) {}
}

const beforeCoverUpload = (file: File) => {
  const isLt5M = file.size / 1024 / 1024 < 5
  if (!isLt5M) {
    ElMessage.error('图片大小不能超过5MB')
    return false
  }
  return true
}

const handleCoverUpload = async (option: any) => {
  if (!book.value) return
  try {
    const res = await uploadCover(book.value.id, option.file)
    book.value.cover_image = res.cover_image
    ElMessage.success('封面上传成功')
  } catch (e) {}
}

const goBack = () => {
  router.push('/books')
}

onMounted(() => {
  loadBook()
})
</script>

<style scoped lang="scss">
.book-detail {
  .cover-card {
    :deep(.el-card__body) {
      display: flex;
      flex-direction: column;
      align-items: center;
    }
  }

  .cover-container {
    width: 100%;
    aspect-ratio: 3/4;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    overflow: hidden;

    .large-cover {
      width: 100%;
      height: 100%;
      object-fit: cover;

      &.placeholder {
        color: rgba(255, 255, 255, 0.6);
        display: flex;
        align-items: center;
        justify-content: center;
      }
    }
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;

    h2 {
      margin: 0;
    }
  }

  .action-section {
    margin-top: 24px;

    h3 {
      margin: 0 0 16px 0;
      font-size: 16px;
    }
  }

  .progress-section {
    .progress-display {
      text-align: center;
      margin-bottom: 12px;

      .current-page {
        font-size: 32px;
        font-weight: 600;
        color: #409eff;
      }
      .divider {
        font-size: 24px;
        color: #c0c4cc;
        margin: 0 8px;
      }
      .total-pages {
        font-size: 24px;
        color: #909399;
      }
    }

    .progress-actions {
      display: flex;
      gap: 12px;
      margin-top: 12px;
      align-items: center;
    }
  }

  .note-input {
    margin-bottom: 16px;

    .note-actions {
      display: flex;
      gap: 12px;
      margin-top: 8px;
      align-items: center;
      justify-content: flex-end;
    }
  }

  .notes-list {
    .empty-notes {
      text-align: center;
      color: #909399;
      padding: 20px;
    }

    .note-item {
      padding: 12px;
      background: #f5f7fa;
      border-radius: 6px;
      margin-bottom: 8px;

      .note-header {
        display: flex;
        gap: 12px;
        align-items: center;
        margin-bottom: 8px;

        .note-date {
          font-size: 12px;
          color: #909399;
        }
      }

      .note-content {
        color: #303133;
        line-height: 1.6;
      }
    }
  }

  .borrow-form {
    display: flex;
    gap: 12px;
    align-items: center;
    flex-wrap: wrap;
  }
}
</style>
