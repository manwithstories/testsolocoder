<template>
  <div class="dock-list">
    <div class="page-header">
      <h2 class="page-title">码头列表</h2>
      <el-button v-if="userStore.hasRole(['admin'])" type="primary" @click="showCreateDialog = true">
        <el-icon><Plus /></el-icon>
        添加码头
      </el-button>
    </div>

    <el-row :gutter="16">
      <el-col :xs="24" :sm="12" :md="8" :lg="6" v-for="dock in docks" :key="dock.id">
        <el-card class="dock-card" shadow="hover" @click="$router.push(`/docks/${dock.id}`)">
          <div class="dock-image">
            <img :src="dock.image_url || '/placeholder-dock.jpg'" :alt="dock.name" />
          </div>
          <div class="dock-info">
            <div class="dock-name">{{ dock.name }}</div>
            <div class="dock-address">
              <el-icon><Location /></el-icon>
              {{ dock.address }}
            </div>
            <div class="dock-rating">
              <el-rate :model-value="dock.average_rating" disabled size="small" />
              <span>({{ dock.review_count }})</span>
            </div>
            <div class="dock-amenities" v-if="dock.amenities">
              <el-tag v-for="amenity in dock.amenities.split(',').slice(0, 3)" :key="amenity" size="small">
                {{ amenity.trim() }}
              </el-tag>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-empty v-if="docks.length === 0 && !loading" description="暂无码头信息" />
  </div>

  <el-dialog v-model="showCreateDialog" title="添加码头" width="500px">
    <el-form :model="dockForm" label-width="80px">
      <el-form-item label="名称">
        <el-input v-model="dockForm.name" />
      </el-form-item>
      <el-form-item label="地址">
        <el-input v-model="dockForm.address" />
      </el-form-item>
      <el-form-item label="城市">
        <el-input v-model="dockForm.city" />
      </el-form-item>
      <el-form-item label="国家">
        <el-input v-model="dockForm.country" />
      </el-form-item>
      <el-form-item label="描述">
        <el-input v-model="dockForm.description" type="textarea" :rows="3" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="showCreateDialog = false">取消</el-button>
      <el-button type="primary" @click="handleCreateDock">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getDocksApi, createDockApi } from '@/api/berth'
import { useUserStore } from '@/stores/user'
import type { Dock } from '@/types/berth'

const userStore = useUserStore()
const loading = ref(false)
const docks = ref<Dock[]>([])
const showCreateDialog = ref(false)

const dockForm = reactive({
  name: '',
  address: '',
  city: '',
  country: '',
  description: ''
})

const fetchDocks = async () => {
  loading.value = true
  try {
    const res: any = await getDocksApi({ page_size: 100 })
    docks.value = res.data || []
  } catch (error) {
    console.error('Failed to fetch docks:', error)
  } finally {
    loading.value = false
  }
}

const handleCreateDock = async () => {
  try {
    await createDockApi(dockForm)
    ElMessage.success('码头创建成功')
    showCreateDialog.value = false
    fetchDocks()
  } catch (error) {
    ElMessage.error('创建失败')
  }
}

onMounted(fetchDocks)
</script>

<style lang="scss" scoped>
.dock-list {
  .dock-card {
    margin-bottom: 16px;
    cursor: pointer;
    transition: transform 0.3s;

    &:hover {
      transform: translateY(-4px);
    }

    .dock-image {
      height: 150px;
      overflow: hidden;

      img {
        width: 100%;
        height: 100%;
        object-fit: cover;
      }
    }

    .dock-info {
      padding: 16px;

      .dock-name {
        font-size: 16px;
        font-weight: 600;
        margin-bottom: 8px;
      }

      .dock-address {
        color: rgba(0, 0, 0, 0.45);
        font-size: 13px;
        margin-bottom: 8px;
        display: flex;
        align-items: center;
        gap: 4px;
      }

      .dock-rating {
        margin-bottom: 8px;

        span {
          font-size: 12px;
          color: rgba(0, 0, 0, 0.45);
        }
      }

      .dock-amenities {
        display: flex;
        gap: 4px;
        flex-wrap: wrap;
      }
    }
  }
}
</style>
