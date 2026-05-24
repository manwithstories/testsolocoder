import React, { useEffect, useState, useRef } from 'react';
import { Card, List, Input, Button, Avatar, Empty, Badge, Upload, message, Spin } from 'antd';
import { SendOutlined, PictureOutlined, FileOutlined, UserOutlined } from '@ant-design/icons';
import { useSelector } from 'react-redux';
import { messageApi } from '@/api/message';
import { Message } from '@/types';
import { RootState } from '@/store';

const { TextArea } = Input;

const Messages: React.FC = () => {
  const user = useSelector((state: RootState) => state.auth.user);
  const [conversations, setConversations] = useState<Message[]>([]);
  const [selectedUser, setSelectedUser] = useState<string | null>(null);
  const [messages, setMessages] = useState<Message[]>([]);
  const [newMessage, setNewMessage] = useState('');
  const [loading, setLoading] = useState(false);
  const [sending, setSending] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    fetchConversations();
  }, []);

  useEffect(() => {
    if (selectedUser) {
      fetchMessages(selectedUser);
    }
  }, [selectedUser]);

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const fetchConversations = async () => {
    try {
      const data = await messageApi.getConversations();
      setConversations(data);
    } catch (error) {
      console.error('获取会话列表失败:', error);
    }
  };

  const fetchMessages = async (userId: string) => {
    setLoading(true);
    try {
      const data = await messageApi.getConversation(userId, { page_size: 50 });
      setMessages(data.items || []);
    } catch (error) {
      console.error('获取消息失败:', error);
    } finally {
      setLoading(false);
    }
  };

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  const handleSendMessage = async () => {
    if (!newMessage.trim() || !selectedUser) return;

    setSending(true);
    try {
      await messageApi.send({
        receiver_id: selectedUser,
        type: 'text',
        content: newMessage.trim(),
      });

      setNewMessage('');
      fetchMessages(selectedUser);
      fetchConversations();
    } catch (error: any) {
      message.error(error.message || '发送消息失败');
    } finally {
      setSending(false);
    }
  };

  const handleFileUpload = async (file: File, type: 'image' | 'file') => {
    if (!selectedUser) return;

    setSending(true);
    try {
      await messageApi.send({
        receiver_id: selectedUser,
        type,
        file_name: file.name,
        file_size: file.size,
        file_url: URL.createObjectURL(file),
      });

      fetchMessages(selectedUser);
      fetchConversations();
    } catch (error: any) {
      message.error(error.message || '发送文件失败');
    } finally {
      setSending(false);
    }
  };

  const getOtherUser = (msg: Message) => {
    if (msg.sender_id === user?.id) {
      return { id: msg.receiver_id, name: msg.receiver?.nickname, avatar: msg.receiver?.avatar };
    }
    return { id: msg.sender_id, name: msg.sender?.nickname, avatar: msg.sender?.avatar };
  };

  return (
    <div style={{ display: 'flex', height: 'calc(100vh - 180px)' }}>
      <Card
        style={{ width: 300, marginRight: 16, overflow: 'auto' }}
        title="消息列表"
        bodyStyle={{ padding: 0 }}
      >
        {conversations.length === 0 ? (
          <Empty description="暂无消息" style={{ padding: 40 }} />
        ) : (
          <List
            dataSource={conversations}
            renderItem={(item) => {
              const otherUser = getOtherUser(item);
              return (
                <List.Item
                  onClick={() => setSelectedUser(otherUser.id)}
                  style={{
                    cursor: 'pointer',
                    background: selectedUser === otherUser.id ? '#e6f7ff' : 'transparent',
                    padding: 12,
                  }}
                >
                  <List.Item.Meta
                    avatar={
                      <Badge dot={!item.is_read && item.receiver_id === user?.id}>
                        <Avatar src={otherUser.avatar} icon={<UserOutlined />} />
                      </Badge>
                    }
                    title={otherUser.name}
                    description={
                      <div style={{ overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                        {item.type === 'text' ? item.content : item.type === 'image' ? '[图片]' : '[文件]'}
                      </div>
                    }
                  />
                </List.Item>
              );
            }}
          />
        )}
      </Card>

      <Card
        style={{ flex: 1, display: 'flex', flexDirection: 'column' }}
        title={selectedUser ? '聊天' : '请选择一个会话'}
        bodyStyle={{ flex: 1, display: 'flex', flexDirection: 'column', padding: 0 }}
      >
        {selectedUser ? (
          <>
            <div style={{ flex: 1, overflow: 'auto', padding: 16, background: '#f5f5f5' }}>
              {loading ? (
                <div style={{ textAlign: 'center', padding: 40 }}>
                  <Spin />
                </div>
              ) : messages.length === 0 ? (
                <Empty description="暂无消息" style={{ padding: 40 }} />
              ) : (
                messages.map((msg) => (
                  <div
                    key={msg.id}
                    style={{
                      display: 'flex',
                      justifyContent: msg.sender_id === user?.id ? 'flex-end' : 'flex-start',
                      marginBottom: 16,
                    }}
                  >
                    {msg.sender_id !== user?.id && (
                      <Avatar
                        src={msg.sender?.avatar}
                        icon={<UserOutlined />}
                        style={{ marginRight: 8 }}
                      />
                    )}
                    <div
                      style={{
                        maxWidth: '60%',
                        padding: '8px 12px',
                        borderRadius: 8,
                        background: msg.sender_id === user?.id ? '#1890ff' : '#fff',
                        color: msg.sender_id === user?.id ? '#fff' : '#333',
                        wordBreak: 'break-word',
                      }}
                    >
                      {msg.type === 'text' ? (
                        msg.content
                      ) : msg.type === 'image' ? (
                        <img
                          src={msg.file_url}
                          alt={msg.file_name}
                          style={{ maxWidth: 200, borderRadius: 4 }}
                        />
                      ) : (
                        <a href={msg.file_url} target="_blank" rel="noopener noreferrer">
                          <FileOutlined /> {msg.file_name}
                        </a>
                      )}
                      <div
                        style={{
                          fontSize: 10,
                          marginTop: 4,
                          color: msg.sender_id === user?.id ? 'rgba(255,255,255,0.7)' : '#999',
                        }}
                      >
                        {new Date(msg.created_at).toLocaleTimeString()}
                      </div>
                    </div>
                    {msg.sender_id === user?.id && (
                      <Avatar
                        src={user?.avatar}
                        icon={<UserOutlined />}
                        style={{ marginLeft: 8 }}
                      />
                    )}
                  </div>
                ))
              )}
              <div ref={messagesEndRef} />
            </div>

            <div style={{ padding: 16, borderTop: '1px solid #e8e8e8' }}>
              <div style={{ display: 'flex', gap: 8, marginBottom: 8 }}>
                <Upload
                  showUploadList={false}
                  beforeUpload={(file) => {
                    handleFileUpload(file, file.type.startsWith('image/') ? 'image' : 'file');
                    return false;
                  }}
                >
                  <Button icon={<PictureOutlined />} size="small">
                    图片
                  </Button>
                </Upload>
                <Upload
                  showUploadList={false}
                  beforeUpload={(file) => {
                    handleFileUpload(file, 'file');
                    return false;
                  }}
                >
                  <Button icon={<FileOutlined />} size="small">
                    文件
                  </Button>
                </Upload>
              </div>
              <div style={{ display: 'flex', gap: 8 }}>
                <TextArea
                  value={newMessage}
                  onChange={(e) => setNewMessage(e.target.value)}
                  onPressEnter={handleSendMessage}
                  placeholder="输入消息..."
                  autoSize={{ minRows: 1, maxRows: 4 }}
                  disabled={sending}
                />
                <Button
                  type="primary"
                  icon={<SendOutlined />}
                  onClick={handleSendMessage}
                  loading={sending}
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
      </Card>
    </div>
  );
};

export default Messages;
