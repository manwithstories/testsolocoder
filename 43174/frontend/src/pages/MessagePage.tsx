import React, { useState, useEffect, useRef } from 'react';
import { List, Button, Avatar, Input as AntInput, Empty, Badge } from 'antd';
import { SendOutlined, UserOutlined } from '@ant-design/icons';
import { messageApi } from '../services/api';
import { Message } from '../types';
import { Loading } from '../components/Loading';
import { useAuthStore } from '../context/authStore';

export const MessagePage: React.FC = () => {
  const [loading, setLoading] = useState(true);
  const [messages, setMessages] = useState<Message[]>([]);
  const [newMessage, setNewMessage] = useState('');
  const [contactUserId] = useState<string | null>(null);
  const [unreadCount, setUnreadCount] = useState(0);
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const { user } = useAuthStore();

  useEffect(() => {
    loadUnreadCount();
    if (contactUserId) {
      loadMessages();
    } else {
      setLoading(false);
    }
  }, [contactUserId]);

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const loadUnreadCount = async () => {
    try {
      const response: any = await messageApi.getUnreadCount();
      setUnreadCount(response.data?.unread_count || 0);
    } catch (error) {
      console.error('Failed to load unread count:', error);
    }
  };

  const loadMessages = async () => {
    if (!user?.id || !contactUserId) return;
    setLoading(true);
    try {
      const response: any = await messageApi.getConversation(user.id, contactUserId, {
        page: 1, page_size: 50,
      });
      setMessages(response.data || []);
      await messageApi.markAsRead();
      setUnreadCount(0);
    } catch (error) {
      console.error('Failed to load messages:', error);
    } finally {
      setLoading(false);
    }
  };

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  const handleSendMessage = async () => {
    if (!newMessage.trim() || !user?.id || !contactUserId) return;

    try {
      await messageApi.create({
        sender_id: user.id,
        receiver_id: contactUserId,
        content: newMessage.trim(),
      });
      setNewMessage('');
      loadMessages();
    } catch (error: any) {
      console.error('Failed to send message:', error);
    }
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSendMessage();
    }
  };

  if (loading && contactUserId) return <Loading />;

  return (
    <div className="max-w-4xl mx-auto px-4 py-8 h-[calc(100vh-64px)]">
      <h1 className="text-2xl font-bold mb-6">消息中心</h1>

      <div className="flex h-full bg-white rounded-lg shadow overflow-hidden">
        <div className="w-64 border-r">
          <div className="p-4 border-b">
            <h3 className="font-semibold">联系人</h3>
            {unreadCount > 0 && (
              <Badge count={unreadCount} className="ml-2" />
            )}
          </div>
          <div className="overflow-y-auto h-[calc(100%-60px)]">
            <List
              dataSource={[]}
              renderItem={() => null}
              locale={{ emptyText: <Empty description="暂无联系人" /> }}
            />
          </div>
        </div>

        <div className="flex-1 flex flex-col">
          {contactUserId ? (
            <>
              <div className="p-4 border-b flex items-center">
                <Avatar icon={<UserOutlined />} className="mr-3" />
                <span className="font-medium">聊天</span>
              </div>

              <div className="flex-1 overflow-y-auto p-4 bg-gray-50">
                {messages.length === 0 ? (
                  <Empty description="暂无消息" className="py-8" />
                ) : (
                  <div className="space-y-4">
                    {messages.map((msg) => (
                    <div
                      key={msg.id}
                      className={`flex ${msg.sender_id === user?.id ? 'justify-end' : 'justify-start'}`}
                    >
                      <div
                        className={`max-w-[70%] px-4 py-2 rounded-lg ${
                          msg.sender_id === user?.id
                            ? 'bg-blue-500 text-white'
                            : 'bg-white'
                        }`}
                      >
                        <p className="text-sm">{msg.content}</p>
                        <p className={`text-xs mt-1 ${msg.sender_id === user?.id ? 'text-blue-100' : 'text-gray-400'}`}>
                          {new Date(msg.created_at).toLocaleTimeString()}
                        </p>
                      </div>
                    </div>
                  ))}
                    <div ref={messagesEndRef} />
                  </div>
                )}
              </div>

              <div className="p-4 border-t">
                <div className="flex gap-2">
                  <AntInput
                    value={newMessage}
                    onChange={(e) => setNewMessage(e.target.value)}
                    onKeyPress={handleKeyPress}
                    placeholder="输入消息..."
                  />
                  <Button
                    type="primary"
                    icon={<SendOutlined />}
                    onClick={handleSendMessage}
                  >
                    发送
                  </Button>
                </div>
              </div>
            </>
          ) : (
            <div className="flex-1 flex items-center justify-center">
            <Empty description="选择一个联系人开始聊天" />
          </div>
          )}
        </div>
      </div>
    </div>
  );
};
