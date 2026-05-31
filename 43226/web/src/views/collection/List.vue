<template>
  <div class="collection-list page-container">
    <div class="filter-bar card-shadow p-20 mb-20">
      <el-form :inline="true" :model="query">
        <el-form-item label="搜索">
          <el-input
            v-model="query.keyword"
            placeholder="搜索藏品名称、编号、描述"
            clearable
            @keyup.enter="fetchList"
          >
            <template #append>
              <el-button @click="fetchList"><el-icon><Search /></el-icon></el-button>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="分类">
          <el-select v-model="query.category_id" clearable placeholder="全部分类" @change="fetchList">
            <el-option
              v-for="cat in categories"
              :key="cat.id"
              :label="cat.name"
              :value="cat.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="年代">
          <el-input v-model="query.era" placeholder="年代" clearable @keyup.enter="fetchList" />
        </el-form-item>
        <el-form-item label="材质">
          <el-input v-model="query.material" placeholder="材质" clearable @keyup.enter="fetchList" />
        </el-form-item>
        <el-form-item label="排序">
          <el-select v-model="query.sort_by" @change="fetchList">
            <el-option label="最新发布" value="created_at" />
            <el-option label="浏览量" value="view_count" />
            <el-option label="年代" value="era" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-radio-group v-model="query.sort_order" @change="fetchList">
            <el-radio-button value="desc">降序</el-radio-button>
            <el-radio-button value="asc">升序</el-radio-button>
          </el-radio-group>
        </el-form-item>
      </el-form>
    </div>

    <div v-loading="loading" class="collection-grid">
      <div
        v-for="item in list"
        :key="item.id"
        class="collection-card card-shadow"
        @click="$router.push(`/collections/${item.id}`)"
      >
        <div class="card-image">
          <img :src="item.image_url || '/placeholder.svg'" :alt="item.name" />
        </div>
        <div class="card-content">
          <div class="flex-between">
            <h3 class="card-title">{{ item.name }}</h3>
            <el-tag size="small" type="info">{{ item.category?.name }}</el-tag>
          </div>
          <div class="card-meta">
            <span v-if="item.era"><el-icon><Clock /></el-icon> {{ item.era }}</span>
            <span v-if="item.material"><el-icon><Collection /></el-icon> {{ item.material }}</span>
          </div>
          <p class="card-desc">{{ item.description }}</p>
          <div class="card-footer">
            <span v-if="item.size" class="meta-item">尺寸: {{ item.size }}</span>
            <span class="views"><el-icon><View /></el-icon> {{ item.view_count }}</span>
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
        :page-sizes="[12, 24, 48, 96]"
        @size-change="fetchList"
        @current-change="fetchList"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import * as collectionApi from '@/api/collection'
import type { Collection, CollectionQuery, CollectionCategory } from '@/types'

const loading = ref(false)
const list = ref<Collection[]>([])
const total = ref(0)
const categories = ref<CollectionCategory[]>([])

const query = reactive<CollectionQuery>({
  page: 1,
  page_size: 12,
  keyword: '',
  category_id: undefined,
  era: '',
  material: '',
  sort_by: 'created_at',
  sort_order: 'desc'
})

const fetchCategories = async () => {
  try {
    const res = await collectionApi.listCategories()
    categories.value = res.data
  } catch (e) {
    console.error(e)
  }
}

const fetchList = async () => {
  try {
    loading.value = true
    const res = await collectionApi.listCollections(query)
    list.value = res.data.list
    total.value = res.data.total
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchCategories()
  fetchList()
})
</script>

<style scoped lang="scss">
.collection-list {
  max-width: 1400px;
  margin: 0 auto;
}

.collection-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 24px;
}

.collection-card {
  border-radius: 12px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s ease;

  &:hover {
    transform: translateY(-4px);
  }

  .card-image {
    height: 200px;
    overflow: hidden;

    img {
      width: 100%;
      height: 100%;
      object-fit: cover;
      transition: transform 0.3s;
    }
  }

  &:hover .card-image img {
    transform: scale(1.05);
  }

  .card-content {
    padding: 16px;

    .card-title {
      font-size: 16px;
      margin: 0 0 10px;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }

    .card-meta {
      display: flex;
      gap: 16px;
      margin-bottom: 10px;
      font-size: 12px;
      color: #909399;

      span {
        display: flex;
        align-items: center;
        gap: 4px;
      }
    }

    .card-desc {
      color: #606266;
      font-size: 13px;
      line-height: 1.5;
      margin-bottom: 12px;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
    }

    .card-footer {
      display: flex;
      justify-content: space-between;
      align-items: center;
      font-size: 12px;
      color: #909399;

      .views {
        display: flex;
        align-items: center;
        gap: 4px;
      }
    }
  }
}
</style>
