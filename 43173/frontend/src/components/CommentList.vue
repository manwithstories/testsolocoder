<template>
  <div class="comment-list">
    <div class="comment-input" v-if="userStore.isLoggedIn">
      <el-avatar :size="40" :src="userStore.user?.avatar">
        {{ userStore.user?.nickname?.charAt(0) || 'U' }}
      </el-avatar>
      <div class="input-area">
        <el-input
          v-model="commentText"
          type="textarea"
          :rows="2"
          placeholder="写下你的评论..."
          maxlength="500"
          show-word-limit
        />
        <div class="actions">
          <el-button type="primary" :disabled="!commentText.trim()" @click="submitComment">
            发表评论
          </el-button>
        </div>
      </div>
    </div>
    
    <div class="comment-input-login" v-else>
      <el-text type="info">请先登录后发表评论</el-text>
      <el-button type="primary" link @click="goToLogin">登录</el-button>
    </div>
    
    <div class="comments" v-loading="loading">
      <div v-for="comment in comments" :key="comment.id" class="comment-item">
        <el-avatar :size="40" :src="comment.user?.avatar">
          {{ comment.user?.nickname?.charAt(0) || 'U' }}
        </el-avatar>
        <div class="comment-content">
          <div class="comment-header">
            <span class="username">{{ comment.user?.nickname || comment.user?.username }}</span>
            <span class="time">{{ formatTime(comment.created_at) }}</span>
          </div>
          <div class="comment-text">{{ comment.content }}</div>
          <div class="comment-actions">
            <el-button text size="small" @click="likeComment(comment)">
              <el-icon><Star /></el-icon>
              {{ comment.like_count }}
            </el-button>
            <el-button text size="small" @click="replyToComment(comment)">
              <el-icon><ChatDotRound /></el-icon>
              回复
            </el-button>
            <el-button 
              v-if="comment.user_id === userStore.user?.id" 
              text 
              size="small" 
              type="danger"
              @click="deleteComment(comment)"
            >
              删除
            </el-button>
          </div>
          
          <div v-if="comment.replies?.length" class="replies">
            <div v-for="reply in comment.replies" :key="reply.id" class="reply-item">
              <el-avatar :size="28" :src="reply.user?.avatar">
                {{ reply.user?.nickname?.charAt(0) || 'U' }}
              </el-avatar>
              <div class="reply-content">
                <span class="username">{{ reply.user?.nickname || reply.user?.username }}</span>
                <span class="text">{{ reply.content }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <el-empty v-if="comments.length === 0 && !loading" description="暂无评论" />
    </div>
    
    <div class="pagination" v-if="total > 0">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        layout="prev, pager, next"
        @current-change="loadComments"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { communityApi } from '@/api/community'
import { useUserStore } from '@/stores/user'
import dayjs from 'dayjs'

const props = defineProps<{
  workId: number
}>()

const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const commentText = ref('')
const comments = ref<any[]>([])
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)

onMounted(() => {
  loadComments()
})

async function loadComments() {
  loading.value = true
  try {
    const res = await communityApi.getComments({
      work_id: props.workId,
      page: page.value,
      page_size: pageSize.value
    })
    comments.value = res.list
    total.value = res.total
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function submitComment() {
  if (!commentText.value.trim()) return
  
  try {
    await communityApi.createComment({
      work_id: props.workId,
      content: commentText.value.trim()
    })
    ElMessage.success('评论成功')
    commentText.value = ''
    page.value = 1
    loadComments()
  } catch (e) {
    console.error(e)
  }
}

function likeComment(comment: any) {
  ElMessage.info('点赞功能开发中')
}

function replyToComment(comment: any) {
  ElMessage.info('回复功能开发中')
}

async function deleteComment(comment: any) {
  try {
    await ElMessageBox.confirm('确定要删除这条评论吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await communityApi.deleteComment(comment.id)
    ElMessage.success('删除成功')
    loadComments()
  } catch (e) {
    if (e !== 'cancel') {
      console.error(e)
    }
  }
}

function formatTime(time: string) {
  return dayjs(time).format('YYYY-MM-DD HH:mm')
}

function goToLogin() {
  router.push('/login')
}
</script>

<style scoped lang="scss">
.comment-list {
  .comment-input {
    display: flex;
    gap: 12px;
    margin-bottom: 24px;
    
    .input-area {
      flex: 1;
      
      .actions {
        margin-top: 8px;
        text-align: right;
      }
    }
  }
  
  .comment-input-login {
    text-align: center;
    padding: 20px;
    margin-bottom: 24px;
    background: #f5f7fa;
    border-radius: 8px;
  }
  
  .comments {
    .comment-item {
      display: flex;
      gap: 12px;
      padding: 16px 0;
      border-bottom: 1px solid var(--border-color);
      
      .comment-content {
        flex: 1;
        
        .comment-header {
          display: flex;
          gap: 12px;
          margin-bottom: 4px;
          
          .username {
            font-weight: 500;
          }
          
          .time {
            font-size: 12px;
            color: var(--text-light);
          }
        }
        
        .comment-text {
          margin-bottom: 8px;
          line-height: 1.6;
        }
        
        .comment-actions {
          display: flex;
          gap: 16px;
          
          .el-button {
            padding: 0;
          }
        }
        
        .replies {
          margin-top: 12px;
          padding: 12px;
          background: #f5f7fa;
          border-radius: 4px;
          
          .reply-item {
            display: flex;
            gap: 8px;
            padding: 8px 0;
            
            .reply-content {
              flex: 1;
              font-size: 14px;
              
              .username {
                color: var(--primary-color);
                margin-right: 8px;
              }
            }
          }
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
