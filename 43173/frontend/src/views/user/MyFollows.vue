<template>
  <div class="my-follows">
    <div class="page-header">
      <h2>我的关注</h2>
    </div>
    
    <el-tabs v-model="activeTab">
      <el-tab-pane label="关注的音乐人" name="artists" />
      <el-tab-pane label="关注的用户" name="users" />
      <el-tab-pane label="粉丝" name="followers" />
    </el-tabs>
    
    <div class="follows-list" v-loading="loading">
      <div 
        v-for="item in list" 
        :key="item.id"
        class="follow-item"
      >
        <el-avatar :size="50" :src="item.avatar">
          {{ item.nickname?.charAt(0) || 'U' }}
        </el-avatar>
        <div class="info">
          <div class="name">{{ item.nickname }}</div>
          <div class="desc text-ellipsis">{{ item.bio || item.username }}</div>
        </div>
        <el-button 
          :type="isFollowing(item) ? '' : 'primary'"
          @click="toggleFollow(item)"
        >
          {{ isFollowing(item) ? '已关注' : '关注' }}
        </el-button>
      </div>
      
      <el-empty v-if="list.length === 0 && !loading" description="暂无数据" />
    </div>
    
    <div class="pagination" v-if="total > 0">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="loadList"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { communityApi } from '@/api/community'

const loading = ref(false)
const activeTab = ref('artists')
const list = ref<any[]>([])
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const followingIds = ref<number[]>([])

onMounted(() => {
  loadList()
})

watch(activeTab, () => {
  page.value = 1
  loadList()
})

async function loadList() {
  loading.value = true
  try {
    let res: any
    
    if (activeTab.value === 'artists') {
      res = await communityApi.getFollowingArtists({
        page: page.value,
        page_size: pageSize.value
      })
    } else if (activeTab.value === 'users') {
      res = await communityApi.getFollowingUsers({
        page: page.value,
        page_size: pageSize.value
      })
    } else {
      res = await communityApi.getMyFollowers({
        page: page.value,
        page_size: pageSize.value
      })
    }
    
    list.value = res.list || []
    total.value = res.total || 0
    followingIds.value = list.value.map((item: any) => item.id)
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function isFollowing(item: any) {
  return followingIds.value.includes(item.id) && activeTab.value !== 'followers'
}

async function toggleFollow(item: any) {
  try {
    if (isFollowing(item)) {
      await communityApi.unfollow(item.id)
      ElMessage.success('已取消关注')
    } else {
      await communityApi.follow(item.id)
      ElMessage.success('关注成功')
    }
    loadList()
  } catch (e) {
    console.error(e)
  }
}
</script>

<style scoped lang="scss">
.my-follows {
  .follows-list {
    .follow-item {
      display: flex;
      align-items: center;
      gap: 16px;
      padding: 16px;
      border-bottom: 1px solid var(--border-color);
      
      .info {
        flex: 1;
        min-width: 0;
        
        .name {
          font-weight: 500;
          margin-bottom: 4px;
        }
        
        .desc {
          font-size: 13px;
          color: var(--text-light);
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
