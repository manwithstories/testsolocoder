import { useState, useEffect, useRef } from 'react'
import { messageApi, userApi } from '@/services/api'
import { Message, User } from '@/types'
import { Send, Paperclip, Search, Phone, Video, MoreVertical } from 'lucide-react'

export default function Messages() {
  const [conversations, setConversations] = useState<any[]>([])
  const [messages, setMessages] = useState<Message[]>([])
  const [selectedUser, setSelectedUser] = useState<User | null>(null)
  const [newMessage, setNewMessage] = useState('')
  const messagesEndRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    loadConversations()
  }, [])

  useEffect(() => {
    if (selectedUser) {
      loadMessages()
    }
  }, [selectedUser])

  useEffect(() => {
    scrollToBottom()
  }, [messages])

  const loadConversations = async () => {
    try {
      const res = await messageApi.getConversations()
      setConversations(res.data)
    } catch (error) {
      console.error('Failed to load conversations:', error)
    }
  }

  const loadMessages = async () => {
    try {
      const res = await messageApi.getMessages({ userId: selectedUser?.id })
      setMessages(res.data)
      if (selectedUser) {
        await messageApi.markRead({ senderId: selectedUser.id })
      }
    } catch (error) {
      console.error('Failed to load messages:', error)
    }
  }

  const handleSendMessage = async () => {
    if (!newMessage.trim() || !selectedUser) return

    try {
      await messageApi.send({
        receiverId: selectedUser.id,
        content: newMessage,
        type: 'text',
      })
      setNewMessage('')
      loadMessages()
    } catch (error) {
      console.error('Failed to send message:', error)
    }
  }

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleSendMessage()
    }
  }

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }

  const handleSelectUser = async (conversation: any) => {
    try {
      const res = await userApi.getUserById(conversation.userId)
      setSelectedUser(res.data)
    } catch (error) {
      console.error('Failed to load user:', error)
    }
  }

  return (
    <div className="h-[calc(100vh-8rem)]">
      <div className="card h-full flex overflow-hidden">
        <div className="w-80 border-r border-gray-200 flex flex-col">
          <div className="p-4 border-b border-gray-200">
            <h2 className="text-lg font-semibold text-gray-900">消息</h2>
          </div>

          <div className="p-3 border-b border-gray-200">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
              <input
                type="text"
                placeholder="搜索..."
                className="input-field pl-10 py-2 text-sm"
              />
            </div>
          </div>

          <div className="flex-1 overflow-y-auto scrollbar-thin">
            {conversations.map((conv) => (
              <button
                key={conv.userId}
                onClick={() => handleSelectUser(conv)}
                className={`w-full p-3 flex items-center gap-3 hover:bg-gray-50 transition-colors ${
                  selectedUser?.id === conv.userId ? 'bg-primary-50' : ''
                }`}
              >
                <div className="w-10 h-10 rounded-full bg-primary-600 flex items-center justify-center text-white text-sm font-medium">
                  {conv.firstName?.[0]}{conv.lastName?.[0]}
                </div>
                <div className="flex-1 min-w-0">
                  <div className="flex items-center justify-between">
                    <span className="font-medium text-gray-900 truncate">
                      {conv.firstName} {conv.lastName}
                    </span>
                    {conv.unreadCount > 0 && (
                      <span className="bg-primary-600 text-white text-xs rounded-full w-5 h-5 flex items-center justify-center">
                        {conv.unreadCount}
                      </span>
                    )}
                  </div>
                  <p className="text-sm text-gray-500 truncate">{conv.lastMessage}</p>
                </div>
              </button>
            ))}
          </div>
        </div>

        <div className="flex-1 flex flex-col">
          {selectedUser ? (
            <>
              <div className="p-4 border-b border-gray-200 flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 rounded-full bg-primary-600 flex items-center justify-center text-white text-sm font-medium">
                    {selectedUser.firstName?.[0]}{selectedUser.lastName?.[0]}
                  </div>
                  <div>
                    <div className="font-medium text-gray-900">
                      {selectedUser.firstName} {selectedUser.lastName}
                    </div>
                    <div className="text-sm text-gray-500">
                      {selectedUser.role === 'teacher' ? '老师' : '学生'}
                    </div>
                  </div>
                </div>
                <div className="flex items-center gap-2">
                  <button className="p-2 hover:bg-gray-100 rounded-lg">
                    <Phone className="h-5 w-5 text-gray-500" />
                  </button>
                  <button className="p-2 hover:bg-gray-100 rounded-lg">
                    <Video className="h-5 w-5 text-gray-500" />
                  </button>
                  <button className="p-2 hover:bg-gray-100 rounded-lg">
                    <MoreVertical className="h-5 w-5 text-gray-500" />
                  </button>
                </div>
              </div>

              <div className="flex-1 overflow-y-auto p-4 space-y-4 scrollbar-thin bg-gray-50">
                {messages.map((msg) => (
                  <div
                    key={msg.id}
                    className={`flex ${msg.senderId === selectedUser.id ? 'justify-start' : 'justify-end'}`}
                  >
                    <div
                      className={`max-w-[70%] rounded-2xl px-4 py-2 ${
                        msg.senderId === selectedUser.id
                          ? 'bg-white text-gray-900'
                          : 'bg-primary-600 text-white'
                      }`}
                    >
                      <p className="text-sm whitespace-pre-wrap">{msg.content}</p>
                      <span className={`text-xs mt-1 block ${
                        msg.senderId === selectedUser.id ? 'text-gray-400' : 'text-primary-200'
                      }`}>
                        {new Date(msg.createdAt).toLocaleTimeString('zh-CN', {
                          hour: '2-digit',
                          minute: '2-digit',
                        })}
                      </span>
                    </div>
                  </div>
                ))}
                <div ref={messagesEndRef} />
              </div>

              <div className="p-4 border-t border-gray-200">
                <div className="flex items-end gap-2">
                  <button className="p-2 hover:bg-gray-100 rounded-lg">
                    <Paperclip className="h-5 w-5 text-gray-500" />
                  </button>
                  <textarea
                    value={newMessage}
                    onChange={(e) => setNewMessage(e.target.value)}
                    onKeyPress={handleKeyPress}
                    placeholder="输入消息..."
                    className="input-field resize-none h-10 py-2"
                    rows={1}
                  />
                  <button
                    onClick={handleSendMessage}
                    disabled={!newMessage.trim()}
                    className="p-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 disabled:opacity-50"
                  >
                    <Send className="h-5 w-5" />
                  </button>
                </div>
              </div>
            </>
          ) : (
            <div className="flex-1 flex items-center justify-center">
              <div className="text-center">
                <div className="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
                  <Send className="h-8 w-8 text-gray-400" />
                </div>
                <p className="text-gray-500">选择一个对话开始聊天</p>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
