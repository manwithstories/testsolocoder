<template>
  <div class="book-list">
    <el-card shadow="never" class="filter-card">
      <el-form :inline="true" :model="filters">
        <el-form-item label="搜索">
          <el-input v-model="filters.search" placeholder="书名、作者、ISBN" clearable style="width: 240px" @keyup.enter="loadBooks" />
        </el-form-item>
        <el-form-item label="阅读状态">
          <el-select v-model="filters.status" clearable placeholder="全部" style="width: 140px">
            <el-option label="想读" value="to_read" />
            <el-option label="在读" value="reading" />
            <el-option label="已读" value="completed" />
            <el-option label="放弃" value="abandoned" />
          </el-select>
        </el-form-item>
        <el-form-item label="标签">
          <el-select v-model="filters.tag_ids" multiple filterable placeholder="选择标签" style="width: 200px">
            <el-option v-for="tag in tags" :key="tag.id" :label="tag.name" :value="tag.id">
              <span :style="{ color: tag.color }">●</span> {{ tag.name }}
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadBooks">搜索</el-button>
          <el-button @click="resetFilters">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card shadow="never" class="list-card" style="margin-top: 16px">
      <div v-loading="loading" class="book-grid">
        <div v-for="book in books" :key="book.id" class="book-card card-hover" @click="goToDetail(book.id)">
          <div class="cover-wrapper">
            <img v-if="book.cover_image" :src="book.cover_image" class="book-cover" />
            <div v-else class="book-cover placeholder">
              <el-icon :size="40"><Reading /></el-icon>
            </div>
            <div class="status-badge" :class="book.reading_status">
              {{ statusText[book.reading_status] }}
            </div>
          </div>
          <div class="book-info">
            <div class="book-title">{{ book.title }}</div>
            <div class="book-author">{{ book.author || '未知作者' }}</div>
            <div class="book-tags">
              <el-tag
                v-for="tag in book.tags?.slice(0, 3)"
                :key="tag.id"
                size="small"
                :style="{ backgroundColor: tag.color + '20', color: tag.color, borderColor: tag.color + '40' }"
              >
                {{ tag.name }}
              </el-tag>
            </div>
            <div class="book-progress">
              <el-progress
                :percentage="Math.round(book.reading_progress)"
                :stroke-width="6"
                :show-text="false"
              />
              <span class="progress-text">{{ book.current_page }}/{{ book.total_pages || '-' }}</span>
            </div>
          </div>
          <div class="book-actions" @click.stop>
            <el-button type="primary" link @click="goToDetail(book.id)">
              <el-icon><View /></el-icon>
            </el-button>
            <el-button type="danger" link @click="handleDelete(book)">
              <el-icon><Delete /></el-icon>
            </el-button>
          </div>
        </div>
        <div v-if="!loading && !books.length" class="empty-state">
          <el-empty description="暂无图书，点击右上角添加" />
        </div>
      </div>
      <div class="pagination">
        <el-pagination
          v-model:current-page="filters.page"
          v-model:page-size="filters.page_size"
          :total="total"
          :page-sizes="[12, 24, 48]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="loadBooks"
          @current-change="loadBooks"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { deleteBook, type BookQueryParams } from '@/api/book'
import { useBookStore } from '@/store/book'
import type { Book, ReadingStatus } from '@/types'

const router = useRouter()
const store = useBookStore()

const books = ref<Book[]>([])
const tags = ref(store.tags)
const loading = ref(false)
const total = ref(0)

const statusText: Record<ReadingStatus, string> = {
  to_read: '想读',
  reading: '在读',
  completed: '已读',
  abandoned: '放弃'
}

const filters = reactive<BookQueryParams>({
  page: 1,
  page_size: 12,
  search: '',
  status: '',
  tag_ids: [],
  sort_by: 'created_at',
  sort_order: 'desc'
})

const loadBooks = async () => {
  loading.value = true
  try {
    const res = await store.fetchBooks(filters)
    books.value = store.books
    total.value = store.total
  } catch (e) {
  } finally {
    loading.value = false
  }
}

const resetFilters = () => {
  filters.page = 1
  filters.search = ''
  filters.status = ''
  filters.tag_ids = []
  loadBooks()
}

const goToDetail = (id: number) => {
  router.push(`/books/${id}`)
}

const handleDelete = async (book: Book) => {
  try {
    await ElMessageBox.confirm(`确定要删除《${book.title}》吗？`, '删除确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await deleteBook(book.id)
    ElMessage.success('删除成功')
    loadBooks()
  } catch (e) {}
}

onMounted(() => {
  store.fetchTags()
  loadBooks()
})
</script>

<style scoped lang="scss">
.book-list {
  .filter-card {
    :deep(.el-card__body) {
      padding: 16px 20px;
    }
  }

  .book-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
    gap: 16px;
    min-height: 400px;
  }

  .book-card {
    border: 1px solid #ebeef5;
    border-radius: 8px;
    overflow: hidden;
    cursor: pointer;
    background: #fff;
    position: relative;

    .cover-wrapper {
      position: relative;
      height: 180px;
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      display: flex;
      align-items: center;
      justify-content: center;

      .book-cover {
        height: 100%;
        width: auto;
        max-width: 100%;
        object-fit: cover;

        &.placeholder {
          color: rgba(255, 255, 255, 0.6);
          display: flex;
          align-items: center;
          justify-content: center;
        }
      }

      .status-badge {
        position: absolute;
        top: 8px;
        right: 8px;
        padding: 2px 8px;
        border-radius: 4px;
        font-size: 12px;
        color: #fff;

        &.to_read { background-color: #909399; }
        &.reading { background-color: #409eff; }
        &.completed { background-color: #67c23a; }
        &.abandoned { background-color: #f56c6c; }
      }
    }

    .book-info {
      padding: 12px;

      .book-title {
        font-weight: 600;
        color: #303133;
        margin-bottom: 4px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }

      .book-author {
        font-size: 13px;
        color: #909399;
        margin-bottom: 8px;
      }

      .book-tags {
        margin-bottom: 8px;
        min-height: 20px;
      }

      .book-progress {
        display: flex;
        align-items: center;
        gap: 8px;

        :deep(.el-progress) {
          flex: 1;
        }

        .progress-text {
          font-size: 12px;
          color: #909399;
          flex-shrink: 0;
        }
      }
    }

    .book-actions {
      position: absolute;
      top: 8px;
      left: 8px;
      display: flex;
      gap: 4px;

      .el-button {
        padding: 4px;
        color: #fff;

        &:hover {
          color: #fff;
          opacity: 0.9;
        }
      }
    }
  }

  .pagination {
    display: flex;
    justify-content: center;
    margin-top: 20px;
  }
}
</style>
