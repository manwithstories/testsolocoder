<template>
  <div class="notes-page">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card shadow="never" class="books-sidebar">
          <template #header>
            <span>选择书籍</span>
          </template>
          <div class="book-list">
            <div
              v-for="book in books"
              :key="book.id"
              class="book-item"
              :class="{ active: selectedBookId === book.id }"
              @click="selectBook(book.id)"
            >
              <div class="book-title">{{ book.title }}</div>
              <div class="book-note-count">
                {{ getNoteCount(book.id) }} 条笔记
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="18">
        <el-card shadow="never">
          <template #header>
            <span v-if="selectedBook">读书笔记 - {{ selectedBook.title }}</span>
            <span v-else>请选择一本书查看笔记</span>
          </template>
          <div v-if="selectedBook" class="notes-content">
            <div class="add-note">
              <el-input
                v-model="newNote.content"
                type="textarea"
                :rows="3"
                placeholder="写下你的读书笔记..."
              />
              <div class="actions">
                <el-input-number
                  v-model="newNote.page"
                  :min="0"
                  :max="selectedBook.total_pages || 9999"
                  size="small"
                  placeholder="页码"
                />
                <el-button type="primary" size="small" @click="addNote">添加笔记</el-button>
              </div>
            </div>
            <div class="notes-list" v-loading="loading">
              <div v-if="!notes.length" class="empty">
                <el-empty description="暂无笔记，开始记录吧" :image-size="100" />
              </div>
              <div v-for="note in notes" :key="note.id" class="note-card">
                <div class="note-header">
                  <el-tag size="small" type="info" v-if="note.page">第 {{ note.page }} 页</el-tag>
                  <span class="note-date">{{ formatDate(note.created_at) }}</span>
                  <el-button type="danger" link size="small" @click="deleteNote(note.id)">删除</el-button>
                </div>
                <div class="note-content">{{ note.content }}</div>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage } from 'element-plus'
import dayjs from 'dayjs'
import { getBooks } from '@/api/book'
import { getNotesByBook, createNote, deleteNote as apiDeleteNote } from '@/api/common'
import type { Book, ReadingNote } from '@/types'

const books = ref<Book[]>([])
const notes = ref<ReadingNote[]>([])
const selectedBookId = ref<number | null>(null)
const loading = ref(false)

const newNote = reactive({
  content: '',
  page: 0
})

const selectedBook = computed(() => {
  return books.value.find(b => b.id === selectedBookId.value) || null
})

const loadBooks = async () => {
  try {
    const res = await getBooks({ page_size: 1000 })
    books.value = res.data
  } catch (e) {}
}

const loadNotes = async () => {
  if (!selectedBookId.value) return
  loading.value = true
  try {
    notes.value = await getNotesByBook(selectedBookId.value)
  } catch (e) {
  } finally {
    loading.value = false
  }
}

const selectBook = (id: number) => {
  selectedBookId.value = id
  newNote.content = ''
  newNote.page = 0
  loadNotes()
}

const addNote = async () => {
  if (!selectedBookId.value || !newNote.content.trim()) {
    ElMessage.warning('请输入笔记内容')
    return
  }
  try {
    const note = await createNote({
      book_id: selectedBookId.value,
      content: newNote.content.trim(),
      page: newNote.page || undefined
    })
    notes.value.unshift(note)
    newNote.content = ''
    newNote.page = 0
    ElMessage.success('笔记已添加')
  } catch (e) {}
}

const deleteNote = async (id: number) => {
  try {
    await apiDeleteNote(id)
    notes.value = notes.value.filter(n => n.id !== id)
    ElMessage.success('笔记已删除')
  } catch (e) {}
}

const getNoteCount = (bookId: number) => {
  return books.value.find(b => b.id === bookId)?.reading_notes?.length || 0
}

const formatDate = (dateStr: string) => {
  return dayjs(dateStr).format('YYYY-MM-DD HH:mm')
}

onMounted(() => {
  loadBooks()
})
</script>

<style scoped lang="scss">
.notes-page {
  .books-sidebar {
    :deep(.el-card__body) {
      padding: 0;
    }
  }

  .book-list {
    max-height: 600px;
    overflow-y: auto;

    .book-item {
      padding: 12px 16px;
      cursor: pointer;
      border-bottom: 1px solid #f0f2f5;
      transition: all 0.2s;

      &:hover {
        background-color: #f5f7fa;
      }

      &.active {
        background-color: #ecf5ff;
        border-left: 3px solid #409eff;
      }

      .book-title {
        font-weight: 500;
        color: #303133;
        margin-bottom: 4px;
      }

      .book-note-count {
        font-size: 12px;
        color: #909399;
      }
    }
  }

  .notes-content {
    min-height: 500px;

    .add-note {
      margin-bottom: 20px;
      padding-bottom: 20px;
      border-bottom: 1px solid #ebeef5;

      .actions {
        display: flex;
        gap: 12px;
        justify-content: flex-end;
        align-items: center;
        margin-top: 12px;
      }
    }

    .notes-list {
      .note-card {
        padding: 16px;
        background: #fafafa;
        border-radius: 8px;
        margin-bottom: 12px;
        border-left: 3px solid #409eff;

        .note-header {
          display: flex;
          align-items: center;
          gap: 12px;
          margin-bottom: 12px;

          .note-date {
            font-size: 12px;
            color: #909399;
          }
        }

        .note-content {
          line-height: 1.8;
          color: #303133;
        }
      }
    }
  }
}
</style>
