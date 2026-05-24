<template>
  <div class="technicians-page">
    <div class="container">
      <div class="page-header">
        <h2 class="page-title">技师列表</h2>
      </div>

      <div class="filter-bar">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索技师"
          clearable
          style="width: 200px;"
          @input="filterTechnicians"
        />
        <el-select v-model="specialtyFilter" placeholder="选择专业" clearable style="width: 150px;" @change="filterTechnicians">
          <el-option label="家电维修" value="家电" />
          <el-option label="数码维修" value="数码" />
          <el-option label="汽车维修" value="汽车" />
        </el-select>
      </div>

      <div class="technician-grid">
        <el-card
          v-for="tech in filteredTechnicians"
          :key="tech.id"
          class="technician-card"
          shadow="hover"
          @click="router.push(`/technicians/${tech.id}`)"
        >
          <div class="tech-info">
            <el-avatar :size="80" :src="tech.avatar">
              {{ tech.username?.charAt(0) }}
            </el-avatar>
            <div class="tech-detail">
              <h3 class="tech-name">{{ tech.real_name || tech.username }}</h3>
              <div class="tech-specialty">{{ tech.specialty }}</div>
              <div class="tech-rating">
                <el-rate :model-value="tech.rating" disabled size="small" />
                <span class="rating-text">{{ tech.rating }}</span>
              </div>
              <div class="tech-stats">
                <span>完成{{ tech.completed_orders }}单</span>
                <span>{{ tech.experience_years }}年经验</span>
              </div>
            </div>
          </div>
        </el-card>
      </div>

      <div v-if="filteredTechnicians.length === 0" class="empty-state">
        <el-empty description="暂无技师" />
      </div>

      <div v-if="total > pageSize" class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          :page-size="pageSize"
          :total="total"
          layout="prev, pager, next"
          @current-change="loadTechnicians"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { userApi } from '@/api/user'

const router = useRouter()

const technicians = ref<any[]>([])
const searchKeyword = ref('')
const specialtyFilter = ref('')
const currentPage = ref(1)
const pageSize = ref(12)
const total = ref(0)

const filteredTechnicians = computed(() => {
  let result = technicians.value

  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    result = result.filter(tech =>
      tech.username?.toLowerCase().includes(keyword) ||
      tech.real_name?.toLowerCase().includes(keyword) ||
      tech.specialty?.toLowerCase().includes(keyword)
    )
  }

  if (specialtyFilter.value) {
    result = result.filter(tech => tech.specialty?.includes(specialtyFilter.value))
  }

  return result
})

onMounted(() => {
  loadTechnicians()
})

async function loadTechnicians() {
  try {
    const res = await userApi.getTechnicians({
      page: currentPage.value,
      page_size: pageSize.value
    })
    technicians.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch (error) {
    console.error('Failed to load technicians:', error)
  }
}

function filterTechnicians() {
  // Filtering is handled by computed property
}
</script>

<style scoped>
.technicians-page {
  min-height: 100vh;
  background-color: #f5f7fa;
}

.filter-bar {
  display: flex;
  gap: 15px;
  margin-bottom: 20px;
}

.technician-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

.technician-card {
  cursor: pointer;
  transition: all 0.3s;
}

.technician-card:hover {
  transform: translateY(-5px);
}

.tech-info {
  text-align: center;
}

.tech-detail {
  margin-top: 15px;
}

.tech-name {
  margin: 0 0 8px 0;
  font-size: 18px;
  color: #303133;
}

.tech-specialty {
  color: #909399;
  font-size: 14px;
  margin-bottom: 10px;
}

.tech-rating {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 8px;
  margin-bottom: 10px;
}

.rating-text {
  color: #fba934;
  font-weight: 600;
}

.tech-stats {
  display: flex;
  justify-content: center;
  gap: 15px;
  font-size: 12px;
  color: #606266;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 30px;
}
</style>
