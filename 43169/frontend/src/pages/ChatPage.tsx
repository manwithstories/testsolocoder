import { useState, useEffect, useRef } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { List, Avatar, Input, Button, Badge, Empty, Upload, message } from 'antd'
import { UserOutlined, SendOutlined, PictureOutlined, AudioOutlined } from '@ant-design/icons'
import { useParams } from 'react-router-dom'
import { chatApi, ChatMessage, ChatSession, PageData } from '@/api/endpoints'
import { useAuthStore } from '@/store/authStore'

const { TextArea } = Input

const emptyPageData: PageData<ChatMessage> = {
  list: [],
  total: 0,
  page: 1,
  page_size: 100,
  total_pages: 0,
}

export default function ChatPage() {
  const { userId } = useParams()
  const { user } = useAuthStore()
  const [selectedSession, setSelectedSession] = useState<ChatSession | null>(null)
  const [messageText, setMessageText] = useState('')
  const [messages, setMessages] = useState<ChatMessage[]>([])
  const messagesEndRef = useRef<HTMLDivElement>(null)
  const queryClient = useQueryClient()
  const [ws, setWs] = useState<WebSocket | null>(null)

  const { data: sessionsData } = useQuery({
    queryKey: ['chatSessions'],
    queryFn: () => chatApi.getSessions({ page: 1, page_size: 50 }),
  })

  const { data: historyData, refetch: refetchHistory } = useQuery({
    queryKey: ['chatHistory', selectedSession?.id],
    queryFn: () => {
      if (!selectedSession) return Promise.resolve(emptyPageData)
      const otherId = selectedSession.user_a_id === user?.id ? selectedSession.user_b_id : selectedSession.user_a_id
      return chatApi.getHistory(otherId, { page: 1, page_size: 100 })
    },
    enabled: !!selectedSession,
  })

  useEffect(() => {
    if (historyData?.list) {
      setMessages([...historyData.list].reverse())
    }
  }, [historyData])

  useEffect(() => {
    if (sessionsData?.list && userId) {
      const session = sessionsData.list.find(
        (s) => s.user_a_id === Number(userId) || s.user_b_id === Number(userId)
      )
      if (session) setSelectedSession(session)
    }
  }, [sessionsData, userId])

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages])

  useEffect(() => {
    if (user?.id) {
      const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws'
      const wsUrl = `${protocol}://${window.location.host}/api/v1/ws?user_id=${user.id}`
      const websocket = new WebSocket(wsUrl)

      websocket.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data)
          if (data.type === 'message') {
            setMessages((prev) => [...prev, data.message])
          }
        } catch {
          // ignore
        }
      }

      setWs(websocket)
      return () => websocket.close()
    }
  }, [user?.id])

  const sendMutation = useMutation({
    mutationFn: chatApi.sendMessage,
    onSuccess: (data: any) => {
      const msg = data as ChatMessage
      setMessages((prev) => [...prev, msg])
      setMessageText('')
      queryClient.invalidateQueries({ queryKey: ['chatSessions'] })
    },
  })

  const handleSend = () => {
    if (!messageText.trim() || !selectedSession || !user) return

    const otherId = selectedSession.user_a_id === user.id ? selectedSession.user_b_id : selectedSession.user_a_id
    sendMutation.mutate({
      receiver_id: otherId,
      type: 'text',
      content: messageText.trim(),
    })

    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({
        type: 'chat',
        receiver_id: otherId,
        content: messageText.trim(),
      }))
    }
  }

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleSend()
    }
  }

  const handleFileUpload = async (file: File, type: string) => {
    try {
      const res = await chatApi.uploadFile(file, type)
      const data = res as any
      if (data?.url && selectedSession && user) {
        const otherId = selectedSession.user_a_id === user.id ? selectedSession.user_b_id : selectedSession.user_a_id
        sendMutation.mutate({
          receiver_id: otherId,
          type,
          content: data.data.url,
        })
      }
    } catch {
      // handled
    }
  }

  const getOtherUser = (session: ChatSession) => {
    return session.user_a_id === user?.id ? session.user_b_id : session.user_a_id
  }

  const getUnreadCount = (session: ChatSession) => {
    return session.user_a_id === user?.id ? session.unread_b : session.unread_a
  }

  return (
    <div className="chat-container">
      <div className="chat-sidebar">
        {sessionsData?.list.length === 0 ? (
          <Empty description="暂无会话" style={{ marginTop: 40 }} />
        ) : (
          <List
            dataSource={sessionsData?.list || []}
            renderItem={(session) => {
              const otherId = getOtherUser(session)
              const unread = getUnreadCount(session)
              return (
                <div
                  className={`session-item ${selectedSession?.id === session.id ? 'active' : ''}`}
                  onClick={() => {
                    setSelectedSession(session)
                    chatApi.markAsRead(otherId)
                    refetchHistory()
                  }}
                >
                  <Badge count={unread} offset={[-4, 4]}>
                    <Avatar size={48} icon={<UserOutlined />} />
                  </Badge>
                  <div style={{ flex: 1, overflow: 'hidden' }}>
                    <div style={{ fontWeight: 500 }}>用户#{otherId}</div>
                    <div style={{
                      fontSize: 12,
                      color: '#888',
                      overflow: 'hidden',
                      textOverflow: 'ellipsis',
                      whiteSpace: 'nowrap',
                    }}>
                      {session.last_message || '暂无消息'}
                    </div>
                  </div>
                </div>
              )
            }}
          />
        )}
      </div>

      <div className="chat-main">
        {selectedSession ? (
          <>
            <div style={{
              padding: '12px 16px',
              borderBottom: '1px solid #f0f0f0',
              display: 'flex',
              alignItems: 'center',
              gap: 12,
            }}>
              <Avatar icon={<UserOutlined />} />
              <span style={{ fontWeight: 500 }}>用户#{getOtherUser(selectedSession)}</span>
            </div>

            <div className="chat-messages">
              {messages.length === 0 ? (
                <Empty description="暂无消息" style={{ marginTop: 100 }} />
              ) : (
                messages.map((msg) => (
                  <div
                    key={msg.id}
                    className={`message-bubble ${msg.sender_id === user?.id ? 'message-sent' : 'message-received'}`}
                  >
                    {msg.type === 'image' ? (
                      <img src={msg.content} alt="" style={{ maxWidth: 200, borderRadius: 8 }} />
                    ) : msg.type === 'voice' ? (
                      <audio src={msg.content} controls style={{ maxWidth: 200 }} />
                    ) : (
                      msg.content
                    )}
                  </div>
                ))
              )}
              <div ref={messagesEndRef} />
            </div>

            <div className="chat-input">
              <div style={{ display: 'flex', gap: 8, marginBottom: 8 }}>
                <Upload
                  showUploadList={false}
                  beforeUpload={(file) => {
                    handleFileUpload(file, 'image')
                    return false
                  }}
                >
                  <Button icon={<PictureOutlined />} size="small">图片</Button>
                </Upload>
                <Upload
                  showUploadList={false}
                  beforeUpload={(file) => {
                    handleFileUpload(file, 'voice')
                    return false
                  }}
                  accept="audio/*"
                >
                  <Button icon={<AudioOutlined />} size="small">语音</Button>
                </Upload>
              </div>
              <div style={{ display: 'flex', gap: 8 }}>
                <TextArea
                  value={messageText}
                  onChange={(e) => setMessageText(e.target.value)}
                  onKeyDown={handleKeyPress}
                  placeholder="输入消息..."
                  autoSize={{ minRows: 1, maxRows: 4 }}
                  style={{ flex: 1 }}
                />
                <Button
                  type="primary"
                  icon={<SendOutlined />}
                  onClick={handleSend}
                  disabled={!messageText.trim()}
                >
                  发送
                </Button>
              </div>
            </div>
          </>
        ) : (
          <div style={{ flex: 1, display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
            <Empty description="选择一个会话开始聊天" />
          </div>
        )}
      </div>
    </div>
  )
}
