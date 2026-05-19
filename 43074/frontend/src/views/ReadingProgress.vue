<template>
  <div class="reading-progress-page">
    <el-card shadow="never">
      <template #header>
        <span>正在阅读的书籍</span>
      </template>
      <div v-loading="loading" class="reading-list">
        <div v-if="!books.length" class="empty">
          <el-empty description="暂无在读书籍" />
        </div>
        <el-row :gutter="20">
          <el-col :span="12" v-for="book in books" :key="book.id">
            <div class="reading-card card-hover" @click="goToBook(book.id)">
              <div class="card-header">
                <img v-if="book.cover_image" :src="book.cover_image" class="cover" />
                <div v-else class="cover placeholder">
                  <el-icon :size="32"><Reading /></el-icon>
                </div>
                <div class="book-meta">
                  <h3 class="title">{{ book.title }}</h3>
                  <p class="author">{{ book.author || '未知作者' }}</p>
                </div>
              </div>
              <div class="card-body">
                <div class="progress-info">
                  <span class="current">{{ book.current_page }}</span>
                  <span class="separator">/</span>
                  <span class="total">{{ book.total_pages }}</span>
                  <span class="percentage">{{ Math.round(book.reading_progress) }}%</span>
                </div>
                <el-progress
                  :percentage="Math.round(book.reading_progress)"
                  :stroke-width="10"
                  :color="progressColor(book.reading_progress)"
                />
                <div class="quick-actions" @click.stop>
                  <el-input-number
                    v-model="progressMap[book.id]"
                    :min="0"
                    :max="book.total_pages || 9999"
                    size="small"
                    :controls="false"
                  />
                  <el-button type="primary" size="small" @click="updateProgress(book)">
                    更新
                  </el-button>
                  <el-button size="small" @click="markCompleted(book)">
                    已读完
                  </el-button>
                </div>
              </div>
            </div>
          </el-col>
        </el-row>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getCurrentlyReading, updateReadingProgress, updateBookStatus } from '@/api/book'
import type { Book } from '@/types'

const router = useRouter()
const books = ref<Book[]>([])
const loading = ref(false)
const progressMap = reactive<Record<number, number>>({})

const loadBooks = async () => {
  loading.value = true
  try {
    books.value = await getCurrentlyReading()
    books.value.forEach(b => {
      progressMap[b.id] = b.current_page
    })
  } catch (e) {
  } finally {
    loading.value = false
  }
}

const progressColor = (percentage: number) => {
  if (percentage < 30) return '#e6a23c'
  if (percentage < 70) return '#409eff'
  return '#67c23a'
}

const updateProgress = async (book: Book) => {
  try {
    const newPage = progressMap[book.id] ?? book.current_page
    const res = await updateReadingProgress(book.id, newPage)
    book.current_page = res.current_page
    book.reading_progress = res.reading_progress
    if (res.reading_status !== 'reading') {
      books.value = books.value.filter(b => b.id !== book.id)
    }
    ElMessage.success('进度已更新')
  } catch (e) {}
}

const markCompleted = async (book: Book) => {
  try {
    await updateBookStatus(book.id, 'completed')
    books.value = books.value.filter(b => b.id !== book.id)
    ElMessage.success('已标记为已读')
  } catch (e) {}
}

const goToBook = (id: number) => {
  router.push(`/books/${id}`)
}

onMounted(() => {
  loadBooks()
})
</script>

<style scoped lang="scss">
.reading-progress-page {
  .reading-list {
    min-height: 400px;
  }

  .reading-card {
    border: 1px solid #ebeef5;
    border-radius: 8px;
    padding: 16px;
    margin-bottom: 16px;
    cursor: pointer;
    background: #fff;

    .card-header {
      display: flex;
      gap: 12px;
      margin-bottom: 16px;

      .cover {
        width: 60px;
        height: 80px;
        border-radius: 4px;
        object-fit: cover;
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);

        &.placeholder {
          display: flex;
          align-items: center;
          justify-content: center;
          color: rgba(255, 255, 255, 0.6);
        }
      }

      .book-meta {
        flex: 1;

        .title {
          margin: 0 0 4px 0;
          font-size: 16px;
          font-weight: 600;
        }

        .author {
          margin: 0;
          font-size: 13px;
          color: #909399;
        }
      }
    }

    .progress-info {
      display: flex;
      align-items: baseline;
      gap: 6px;
      margin-bottom: 8px;

      .current {
        font-size: 24px;
        font-weight: 600;
        color: #409eff;
      }

      .separator {
        color: #c0c4cc;
      }

      .total {
        color: #909399;
      }

      .percentage {
        margin-left: auto;
        font-weight: 600;
        color: #67c23a;
      }
    }

    .quick-actions {
      display: flex;
      gap: 8px;
      margin-top: 12px;
      align-items: center;
    }
  }
}
</style>
