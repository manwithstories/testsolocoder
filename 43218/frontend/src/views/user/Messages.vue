<template>
  <div class="messages-page">
    <div class="page-header">
      <h2 class="page-title">消息中心</h2>
    </div>

    <div class="messages-content">
      <el-card class="contacts-card">
        <template #header>
          <span>联系人</span>
        </template>
        <div class="contact-list" v-loading="contactsLoading">
          <div
            v-for="contact in contacts"
            :key="contact.userId"
            class="contact-item"
            :class="{ active: currentUserId === contact.userId }"
            @click="selectContact(contact.userId)"
          >
            <el-avatar :size="40">
              {{ contact.userId }}
            </el-avatar>
            <div class="contact-info">
              <div class="contact-msg text-ellipsis">{{ contact.lastMsg }}</div>
              <div class="contact-time">{{ formatTime(contact.lastTime) }}</div>
            </div>
            <el-badge v-if="contact.unreadCount > 0" :value="contact.unreadCount" class="unread-badge" />
          </div>
          <div class="empty-state" v-if="!contactsLoading && contacts.length === 0">
            <el-empty description="暂无联系人" :image-size="60" />
          </div>
        </div>
      </el-card>

      <el-card class="chat-card">
        <template #header>
          <div class="chat-header">
            <span v-if="currentUserId">与用户 {{ currentUserId }} 的对话</span>
            <span v-else>请选择联系人</span>
          </div>
        </template>
        <div class="chat-content" v-loading="messagesLoading">
          <div class="message-list" ref="messageListRef">
            <div
              v-for="msg in messages"
              :key="msg.id"
              class="message-item"
              :class="{ 'is-self': msg.senderId === userStore.userInfo?.id }"
            >
              <el-avatar :size="32">
                {{ msg.senderId }}
              </el-avatar>
              <div class="message-content">
                <div class="message-text">{{ msg.content }}</div>
                <div class="message-time">{{ formatTime(msg.createdAt) }}</div>
              </div>
            </div>
          </div>
        </div>
        <div class="chat-input" v-if="currentUserId">
          <el-input
            v-model="messageInput"
            placeholder="输入消息..."
            @keyup.enter="sendMessage"
          >
            <template #append>
              <el-button :disabled="!messageInput.trim()" @click="sendMessage">
                发送
              </el-button>
            </template>
          </el-input>
        </div>
        <div class="empty-state" v-if="!currentUserId">
          <el-empty description="请选择联系人开始对话" :image-size="80" />
        </div>
      </el-card>

      <el-card class="notifications-card">
        <template #header>
          <div class="notifications-header">
            <span>通知消息</span>
            <el-button type="primary" link @click="markAllAsRead">
              全部已读
            </el-button>
          </div>
        </template>
        <div class="notification-list" v-loading="notificationsLoading">
          <div
            v-for="notification in notifications"
            :key="notification.id"
            class="notification-item"
            :class="{ unread: !notification.isRead }"
            @click="readNotification(notification)"
          >
            <el-icon :size="20" class="notification-icon">
              <Bell />
            </el-icon>
            <div class="notification-content">
              <div class="notification-title">{{ notification.title }}</div>
              <div class="notification-desc text-ellipsis">{{ notification.content }}</div>
              <div class="notification-time">{{ formatTime(notification.createdAt) }}</div>
            </div>
          </div>
          <div class="empty-state" v-if="!notificationsLoading && notifications.length === 0">
            <el-empty description="暂无通知" :image-size="60" />
          </div>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/store/user'
import { messageApi, notificationApi } from '@/api'
import type { Message, Notification } from '@/types'

const userStore = useUserStore()

const contactsLoading = ref(false)
const messagesLoading = ref(false)
const notificationsLoading = ref(false)

const contacts = ref<any[]>([])
const messages = ref<Message[]>([])
const notifications = ref<Notification[]>([])
const currentUserId = ref<number | null>(null)
const messageInput = ref('')
const messageListRef = ref<HTMLElement | null>(null)

function formatTime(time: string): string {
  return new Date(time).toLocaleString()
}

async function fetchContacts() {
  contactsLoading.value = true
  try {
    const res = await messageApi.getContacts()
    contacts.value = res.data
  } catch (error) {
    console.error('Failed to fetch contacts:', error)
  } finally {
    contactsLoading.value = false
  }
}

async function selectContact(userId: number) {
  currentUserId.value = userId
  await fetchMessages(userId)
  await messageApi.markAsRead(userId)
  await fetchContacts()
}

async function fetchMessages(userId: number) {
  messagesLoading.value = true
  try {
    const res = await messageApi.getList(userId, { page: 1, pageSize: 50 })
    messages.value = res.data.reverse()
    await nextTick()
    if (messageListRef.value) {
      messageListRef.value.scrollTop = messageListRef.value.scrollHeight
    }
  } catch (error) {
    console.error('Failed to fetch messages:', error)
  } finally {
    messagesLoading.value = false
  }
}

async function sendMessage() {
  if (!currentUserId.value || !messageInput.value.trim()) return

  try {
    await messageApi.send({
      receiverId: currentUserId.value,
      content: messageInput.value.trim()
    })
    messageInput.value = ''
    await fetchMessages(currentUserId.value)
    await fetchContacts()
  } catch (error) {
    console.error('Failed to send message:', error)
    ElMessage.error('发送失败')
  }
}

async function fetchNotifications() {
  notificationsLoading.value = true
  try {
    const res = await notificationApi.getList({ page: 1, pageSize: 20 })
    notifications.value = res.data
  } catch (error) {
    console.error('Failed to fetch notifications:', error)
  } finally {
    notificationsLoading.value = false
  }
}

async function readNotification(notification: Notification) {
  if (notification.isRead) return
  try {
    await notificationApi.markAsRead(notification.id)
    notification.isRead = true
  } catch (error) {
    console.error('Failed to mark as read:', error)
  }
}

async function markAllAsRead() {
  try {
    await notificationApi.markAllAsRead()
    notifications.value.forEach(n => n.isRead = true)
    ElMessage.success('已全部标记为已读')
  } catch (error) {
    console.error('Failed to mark all as read:', error)
  }
}

onMounted(() => {
  fetchContacts()
  fetchNotifications()
})
</script>

<style lang="scss" scoped>
.messages-page {
  .messages-content {
    display: grid;
    grid-template-columns: 280px 1fr 320px;
    gap: 20px;
    height: calc(100vh - 120px);
  }

  .contacts-card,
  .chat-card,
  .notifications-card {
    height: 100%;
    display: flex;
    flex-direction: column;

    :deep(.el-card__body) {
      flex: 1;
      overflow: hidden;
      display: flex;
      flex-direction: column;
    }
  }

  .contact-list {
    flex: 1;
    overflow-y: auto;

    .contact-item {
      display: flex;
      align-items: center;
      gap: 12px;
      padding: 12px;
      border-radius: 8px;
      cursor: pointer;
      margin-bottom: 8px;
      transition: background 0.2s;

      &:hover {
        background: #f5f7fa;
      }

      &.active {
        background: #e8f3ff;
      }

      .contact-info {
        flex: 1;
        min-width: 0;

        .contact-msg {
          font-size: 13px;
          margin-bottom: 4px;
        }

        .contact-time {
          font-size: 11px;
          color: var(--text-lighter-color);
        }
      }

      .unread-badge {
        margin-left: 8px;
      }
    }
  }

  .chat-card {
    .chat-header {
      font-weight: 500;
    }

    .chat-content {
      flex: 1;
      overflow-y: auto;
      padding: 20px 0;

      .message-list {
        .message-item {
          display: flex;
          gap: 12px;
          margin-bottom: 20px;

          &.is-self {
            flex-direction: row-reverse;

            .message-content {
              text-align: right;
            }

            .message-text {
              background: var(--primary-color);
              color: #fff;
            }
          }

          .message-content {
            max-width: 70%;

            .message-text {
              background: #f5f7fa;
              padding: 12px 16px;
              border-radius: 12px;
              margin-bottom: 4px;
              word-break: break-word;
            }

            .message-time {
              font-size: 11px;
              color: var(--text-lighter-color);
            }
          }
        }
      }
    }

    .chat-input {
      padding-top: 16px;
      border-top: 1px solid #f0f0f0;
    }
  }

  .notifications-card {
    .notifications-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
    }

    .notification-list {
      flex: 1;
      overflow-y: auto;

      .notification-item {
        display: flex;
        gap: 12px;
        padding: 12px;
        border-radius: 8px;
        margin-bottom: 8px;
        cursor: pointer;
        transition: background 0.2s;

        &:hover {
          background: #f5f7fa;
        }

        &.unread {
          background: #f0f9ff;
        }

        .notification-icon {
          color: var(--primary-color);
          flex-shrink: 0;
        }

        .notification-content {
          flex: 1;
          min-width: 0;

          .notification-title {
            font-size: 13px;
            font-weight: 500;
            margin-bottom: 4px;
          }

          .notification-desc {
            font-size: 12px;
            color: var(--text-lighter-color);
            margin-bottom: 4px;
          }

          .notification-time {
            font-size: 11px;
            color: var(--text-lighter-color);
          }
        }
      }
    }
  }

  .empty-state {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
  }
}
</style>
