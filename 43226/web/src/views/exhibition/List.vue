<template>
  <div class="exhibition-list page-container">
    <div class="filter-bar card-shadow p-20 mb-20">
      <el-form :inline="true" :model="query">
        <el-form-item label="搜索">
          <el-input
            v-model="query.keyword"
            placeholder="搜索展览名称或描述"
            clearable
            @keyup.enter="fetchList"
          >
            <template #append>
              <el-button @click="fetchList"><el-icon><Search /></el-icon></el-button>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="query.status" clearable placeholder="全部" @change="fetchList">
            <el-option label="草稿" value="draft" />
            <el-option label="已发布" value="published" />
            <el-option label="已结束" value="closed" />
          </el-select>
        </el-form-item>
        <el-form-item label="开始日期">
          <el-date-picker
            v-model="query.start_date"
            type="date"
            placeholder="选择开始日期"
            value-format="YYYY-MM-DD"
            @change="fetchList"
          />
        </el-form-item>
        <el-form-item label="结束日期">
          <el-date-picker
            v-model="query.end_date"
            type="date"
            placeholder="选择结束日期"
            value-format="YYYY-MM-DD"
            @change="fetchList"
          />
        </el-form-item>
      </el-form>
    </div>

    <div v-loading="loading" class="exhibition-grid">
      <div
        v-for="exhibition in list"
        :key="exhibition.id"
        class="exhibition-card card-shadow"
        @click="$router.push(`/exhibitions/${exhibition.id}`)"
      >
        <div class="card-image">
          <img :src="exhibition.image_url || '/placeholder.svg'" :alt="exhibition.title" />
          <el-tag
            v-if="exhibition.status === 'published'"
            type="success"
            size="small"
            class="status-tag"
          >展览中</el-tag>
          <el-tag
            v-else-if="exhibition.status === 'closed'"
            type="info"
            size="small"
            class="status-tag"
          >已结束</el-tag>
        </div>
        <div class="card-content">
          <h3 class="card-title">{{ exhibition.title }}</h3>
          <div class="card-meta">
            <span><el-icon><Calendar /></el-icon> {{ formatDate(exhibition.start_date) }} - {{ formatDate(exhibition.end_date) }}</span>
            <span><el-icon><Location /></el-icon> {{ exhibition.hall_number || exhibition.location }}</span>
          </div>
          <p class="card-desc">{{ exhibition.description }}</p>
          <div class="card-footer">
            <span class="price" v-if="exhibition.ticket_price > 0">¥{{ exhibition.ticket_price }}</span>
            <span class="price free" v-else>免费参观</span>
            <el-button type="primary" size="small">立即预约</el-button>
          </div>
        </div>
      </div>
    </div>

    <div class="pagination mt-20 flex-center">
      <el-pagination
        v-model:current-page="query.page"
        v-model:page-size="query.page_size"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        :page-sizes="[6, 12, 24, 48]"
        @size-change="fetchList"
        @current-change="fetchList"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import * as exhibitionApi from '@/api/exhibition'
import type { Exhibition, ExhibitionQuery } from '@/types'
import dayjs from 'dayjs'

const loading = ref(false)
const list = ref<Exhibition[]>([])
const total = ref(0)

const query = reactive<ExhibitionQuery>({
  page: 1,
  page_size: 12,
  keyword: '',
  status: 'published',
  start_date: '',
  end_date: ''
})

const formatDate = (date: string) => {
  return dayjs(date).format('YYYY.MM.DD')
}

const fetchList = async () => {
  try {
    loading.value = true
    const res = await exhibitionApi.listExhibitions(query)
    list.value = res.data.list
    total.value = res.data.total
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchList()
})
</script>

<style scoped lang="scss">
.exhibition-list {
  max-width: 1400px;
  margin: 0 auto;
}

.exhibition-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 24px;
}

.exhibition-card {
  border-radius: 12px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s ease;

  &:hover {
    transform: translateY(-4px);
  }

  .card-image {
    position: relative;
    height: 200px;

    img {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }

    .status-tag {
      position: absolute;
      top: 12px;
      right: 12px;
    }
  }

  .card-content {
    padding: 20px;

    .card-title {
      font-size: 18px;
      margin-bottom: 12px;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }

    .card-meta {
      display: flex;
      flex-direction: column;
      gap: 8px;
      margin-bottom: 12px;
      font-size: 13px;
      color: #909399;

      span {
        display: flex;
        align-items: center;
        gap: 4px;
      }
    }

    .card-desc {
      color: #606266;
      font-size: 14px;
      line-height: 1.5;
      margin-bottom: 16px;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
    }

    .card-footer {
      display: flex;
      justify-content: space-between;
      align-items: center;

      .price {
        color: #f56c6c;
        font-weight: 600;
        font-size: 18px;

        &.free {
          color: #67c23a;
        }
      }
    }
  }
}
</style>
