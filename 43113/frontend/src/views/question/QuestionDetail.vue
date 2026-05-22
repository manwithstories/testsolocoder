<template>
  <div class="question-detail-page">
    <el-skeleton v-if="loading" :rows="8" animated />

    <div v-else-if="questionData">
      <div class="question-header">
        <div class="question-title">
          <el-tag v-if="questionData.question.rewardPoints > 0" type="warning" size="large">
            悬赏 {{ questionData.question.rewardPoints }} 积分
          </el-tag>
          <el-tag v-if="questionData.question.isSolved" type="success" size="large">已解决</el-tag>
          <h1>{{ questionData.question.title }}</h1>
        </div>

        <div class="question-meta">
          <el-avatar :size="32" :src="questionData.user.avatar">
            {{ questionData.user.nickname?.charAt(0) || 'U' }}
          </el-avatar>
          <span class="username">{{ questionData.user.nickname || questionData.user.username }}</span>
          <el-tag v-if="questionData.user.isExpert" type="primary" size="small" effect="dark">
            <el-icon><Medal /></el-icon>
            专家
          </el-tag>
          <span>Lv.{{ questionData.user.level }}</span>
          <span>提问于 {{ formatTime(questionData.question.createdAt) }}</span>
        </div>

        <div class="question-stats">
          <span><el-icon><View /></el-icon> {{ questionData.question.views }} 浏览</span>
          <span><el-icon><ChatDotRound /></el-icon> {{ questionData.answerCount }} 回答</span>
          <span><el-icon><Star /></el-icon> {{ questionData.question.collectCount }} 收藏</span>
        </div>

        <div class="question-tags">
          <el-tag
            v-for="tag in questionData.tags"
            :key="tag.id"
            size="large"
            effect="plain"
          >
            {{ tag.name }}
          </el-tag>
        </div>
      </div>

      <el-divider />

      <div class="question-content" v-highlight>
        <div v-html="renderMarkdown(questionData.question.content)"></div>
      </div>

      <div class="question-actions">
        <el-button
          :type="questionData.isFavorited ? 'warning' : 'default'"
          @click="handleFavorite"
        >
          <el-icon><Star /></el-icon>
          {{ questionData.isFavorited ? '已收藏' : '收藏' }}
        </el-button>
        <el-button @click="handleLike">
          <el-icon><ThumbUp /></el-icon>
          点赞 ({{ questionData.question.likeCount }})
        </el-button>
        <el-button @click="showReportDialog = true">
          <el-icon><Warning /></el-icon>
          举报
        </el-button>
      </div>

      <el-divider />

      <div class="answers-section">
        <h2>{{ questionData.answerCount }} 个回答</h2>

        <div class="answer-form">
          <h3>撰写回答</h3>
          <el-input
            v-model="answerContent"
            type="textarea"
            :rows="6"
            placeholder="请输入您的回答，支持 Markdown 格式..."
            v-if="userStore.isLoggedIn"
          />
          <template v-else>
            <el-empty description="请登录后回答问题">
              <router-link to="/login">
                <el-button type="primary">登录</el-button>
              </router-link>
            </el-empty>
          </template>
          <div class="form-actions" v-if="userStore.isLoggedIn">
            <el-button type="primary" :loading="submitting" @click="submitAnswer">
              提交回答
            </el-button>
          </div>
        </div>

        <div class="answers-list">
          <div
            v-for="answer in questionData.answers"
            :key="answer.answer.id"
            class="answer-item"
            :class="{ accepted: answer.answer.isAccepted }"
          >
            <div class="answer-header">
              <el-tag v-if="answer.answer.isAccepted" type="success" size="small">
                <el-icon><Select /></el-icon>
                最佳答案
              </el-tag>
              <el-avatar :size="28" :src="answer.user.avatar">
                {{ answer.user.nickname?.charAt(0) || 'U' }}
              </el-avatar>
              <span class="username">{{ answer.user.nickname || answer.user.username }}</span>
              <el-tag v-if="answer.user.isExpert" type="primary" size="small">专家</el-tag>
              <span>Lv.{{ answer.user.level }}</span>
              <span>{{ formatTime(answer.answer.createdAt) }}</span>
            </div>

            <div class="answer-content" v-highlight>
              <div v-html="renderMarkdown(answer.answer.content)"></div>
            </div>

            <div class="answer-actions">
              <el-button size="small" @click="likeAnswer(answer.answer.id)">
                <el-icon><ThumbUp /></el-icon>
                {{ answer.answer.likeCount }}
              </el-button>
              <el-button size="small" @click="dislikeAnswer(answer.answer.id)">
                <el-icon><ThumbDown /></el-icon>
              </el-button>
              <el-button
                v-if="canAcceptAnswer && !answer.answer.isAccepted"
                size="small"
                type="success"
                @click="acceptAnswer(answer.answer.id)"
              >
                <el-icon><Select /></el-icon>
                采纳
              </el-button>
              <el-dropdown v-if="answer.answer.userId === userStore.userInfo?.id || userStore.isAdmin">
                <el-button size="small">
                  更多 <el-icon><ArrowDown /></el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item @click="deleteAnswer(answer.answer.id)">
                      <el-icon><Delete /></el-icon>
                      删除
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>

            <div class="comments-section">
              <div
                v-for="comment in answer.comments"
                :key="comment.comment.id"
                class="comment-item"
              >
                <el-avatar :size="20" :src="comment.user.avatar">
                  {{ comment.user.nickname?.charAt(0) || 'U' }}
                </el-avatar>
                <span class="comment-user">{{ comment.user.nickname || comment.user.username }}</span>
                <span class="comment-content">{{ comment.comment.content }}</span>
                <span class="comment-time">{{ formatTime(comment.comment.createdAt) }}</span>
              </div>

              <div class="comment-input" v-if="userStore.isLoggedIn">
                <el-input
                  v-model="newComment"
                  size="small"
                  placeholder="发表评论..."
                  @keyup.enter="submitComment(answer.answer.id)"
                >
                  <template #append>
                    <el-button @click="submitComment(answer.answer.id)">发送</el-button>
                  </template>
                </el-input>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <el-dialog v-model="showReportDialog" title="举报内容" width="400px">
      <el-form label-position="top">
        <el-form-item label="举报原因">
          <el-select v-model="reportReason" placeholder="请选择原因">
            <el-option label="垃圾广告" value="spam" />
            <el-option label="违规内容" value="violation" />
            <el-option label="抄袭" value="plagiarism" />
            <el-option label="其他" value="other" />
          </el-select>
        </el-form-item>
        <el-form-item label="详细描述">
          <el-input
            v-model="reportDescription"
            type="textarea"
            :rows="3"
            placeholder="请详细描述举报原因..."
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showReportDialog = false">取消</el-button>
        <el-button type="primary" @click="submitReport">提交举报</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { questionApi, answerApi, commentApi, favoriteApi, auditApi } from '@/api'
import type { QuestionDetail } from '@/types'
import MarkdownIt from 'markdown-it'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'

dayjs.extend(relativeTime)

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true,
  highlight: (str, lang) => {
    return `<pre><code class="hljs language-${lang}">${str}</code></pre>`
  }
})

const loading = ref(true)
const questionData = ref<QuestionDetail | null>(null)
const answerContent = ref('')
const newComment = ref('')
const submitting = ref(false)
const showReportDialog = ref(false)
const reportReason = ref('')
const reportDescription = ref('')

const canAcceptAnswer = computed(() => {
  return userStore.isLoggedIn &&
    questionData.value &&
    questionData.value.question.userId === userStore.userInfo?.id &&
    !questionData.value.question.isSolved
})

const fetchQuestion = async () => {
  loading.value = true
  try {
    const id = Number(route.params.id)
    const res = await questionApi.getQuestion(id)
    questionData.value = res.data || null
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const renderMarkdown = (content: string) => {
  return md.render(content)
}

const formatTime = (time: string) => {
  return dayjs(time).fromNow()
}

const submitAnswer = async () => {
  if (!answerContent.value.trim()) return
  submitting.value = true
  try {
    await answerApi.createAnswer({
      questionId: questionData.value!.question.id,
      content: answerContent.value
    })
    answerContent.value = ''
    fetchQuestion()
  } catch (e) {
    console.error(e)
  } finally {
    submitting.value = false
  }
}

const likeAnswer = async (id: number) => {
  try {
    await answerApi.likeAnswer(id)
    fetchQuestion()
  } catch (e) {
    console.error(e)
  }
}

const dislikeAnswer = async (id: number) => {
  try {
    await answerApi.dislikeAnswer(id)
  } catch (e) {
    console.error(e)
  }
}

const acceptAnswer = async (answerId: number) => {
  try {
    await questionApi.acceptAnswer(questionData.value!.question.id, answerId)
    fetchQuestion()
  } catch (e) {
    console.error(e)
  }
}

const deleteAnswer = async (id: number) => {
  try {
    await answerApi.deleteAnswer(id)
    fetchQuestion()
  } catch (e) {
    console.error(e)
  }
}

const submitComment = async (answerId: number) => {
  if (!newComment.value.trim()) return
  try {
    await commentApi.createComment({
      answerId,
      content: newComment.value
    })
    newComment.value = ''
    fetchQuestion()
  } catch (e) {
    console.error(e)
  }
}

const handleFavorite = async () => {
  if (!questionData.value) return
  try {
    if (questionData.value.isFavorited) {
      await favoriteApi.removeFavorite({
        targetType: 'question',
        targetId: questionData.value.question.id
      })
    } else {
      await favoriteApi.addFavorite({
        targetType: 'question',
        targetId: questionData.value.question.id
      })
    }
    fetchQuestion()
  } catch (e) {
    console.error(e)
  }
}

const handleLike = async () => {
  if (!questionData.value) return
  try {
    await questionApi.likeQuestion(questionData.value.question.id)
    fetchQuestion()
  } catch (e) {
    console.error(e)
  }
}

const submitReport = async () => {
  if (!reportReason.value) return
  try {
    await auditApi.createReport({
      targetType: 'question',
      targetId: questionData.value!.question.id,
      reason: reportReason.value,
      description: reportDescription.value
    })
    showReportDialog.value = false
    reportReason.value = ''
    reportDescription.value = ''
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  fetchQuestion()
})
</script>

<style scoped lang="scss">
.question-detail-page {
  .question-header {
    background: white;
    padding: 24px;
    border-radius: 8px;

    .question-title {
      display: flex;
      align-items: center;
      gap: 8px;
      margin-bottom: 16px;
      flex-wrap: wrap;

      h1 {
        margin: 0;
        font-size: 22px;
      }
    }

    .question-meta {
      display: flex;
      align-items: center;
      gap: 12px;
      margin-bottom: 16px;
      color: #606266;
      font-size: 14px;

      .username {
        font-weight: 500;
      }
    }

    .question-stats {
      display: flex;
      gap: 20px;
      margin-bottom: 16px;
      color: #909399;
      font-size: 14px;
    }

    .question-tags {
      display: flex;
      gap: 8px;
      flex-wrap: wrap;
    }
  }

  .question-content {
    background: white;
    padding: 24px;
    border-radius: 8px;
    margin-bottom: 20px;
    line-height: 1.8;

    :deep(img) {
      max-width: 100%;
    }

    :deep(pre) {
      background: #f5f7fa;
      padding: 12px;
      border-radius: 4px;
      overflow-x: auto;
    }

    :deep(code) {
      font-family: 'Monaco', 'Menlo', monospace;
    }
  }

  .question-actions {
    display: flex;
    gap: 12px;
    margin-bottom: 20px;
  }

  .answers-section {
    background: white;
    padding: 24px;
    border-radius: 8px;

    h2 {
      margin: 0 0 20px 0;
    }

    .answer-form {
      margin-bottom: 24px;

      h3 {
        margin: 0 0 12px 0;
        font-size: 16px;
      }

      .form-actions {
        margin-top: 12px;
        text-align: right;
      }
    }

    .answers-list {
      .answer-item {
        border-bottom: 1px solid #e4e7ed;
        padding: 20px 0;

        &:last-child {
          border-bottom: none;
        }

        &.accepted {
          background: #f0f9eb;
          padding: 20px;
          border-radius: 8px;
          border: 2px solid #e1f3d8;
        }

        .answer-header {
          display: flex;
          align-items: center;
          gap: 8px;
          margin-bottom: 12px;
          font-size: 14px;
          color: #606266;

          .username {
            font-weight: 500;
          }
        }

        .answer-content {
          line-height: 1.8;
          margin-bottom: 12px;

          :deep(pre) {
            background: #f5f7fa;
            padding: 12px;
            border-radius: 4px;
            overflow-x: auto;
          }
        }

        .answer-actions {
          display: flex;
          gap: 8px;
          margin-bottom: 12px;
        }

        .comments-section {
          background: #fafafa;
          padding: 12px;
          border-radius: 4px;

          .comment-item {
            display: flex;
            align-items: center;
            gap: 8px;
            padding: 8px 0;
            font-size: 13px;

            .comment-user {
              font-weight: 500;
            }

            .comment-content {
              flex: 1;
            }

            .comment-time {
              color: #909399;
            }
          }

          .comment-input {
            margin-top: 12px;
          }
        }
      }
    }
  }
}
</style>
