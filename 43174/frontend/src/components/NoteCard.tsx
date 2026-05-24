import React from 'react';
import { Card, Tag, Rate, Button, Avatar } from 'antd';
import { DownloadOutlined, EyeOutlined, StarOutlined, UserOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { Note } from '../types';

interface NoteCardProps {
  note: Note;
  onDownload?: (note: Note) => void;
}

export const NoteCard: React.FC<NoteCardProps> = ({ note, onDownload }) => {
  const navigate = useNavigate();

  const getFileTypeIcon = (fileType: string) => {
    if (fileType?.includes('pdf')) return '📄';
    if (fileType?.includes('image')) return '🖼️';
    return '📁';
  };

  const formatFileSize = (bytes: number) => {
    if (bytes < 1024) return bytes + ' B';
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB';
    return (bytes / (1024 * 1024)).toFixed(1) + ' MB';
  };

  return (
    <Card
      hoverable
      className="card-hover h-full"
      cover={
        <div className="h-36 bg-gradient-to-br from-blue-100 to-purple-100 flex items-center justify-center overflow-hidden">
          {note.cover_image ? (
            <img
              src={note.cover_image}
              alt={note.title}
              className="h-full w-full object-cover"
            />
          ) : (
            <div className="text-5xl">{getFileTypeIcon(note.file_type || '')}</div>
          )}
        </div>
      }
      actions={[
        <Button
          type="text"
          icon={<EyeOutlined />}
          onClick={() => navigate(`/notes/${note.id}`)}
        >
          查看
        </Button>,
        <Button
          type="primary"
          icon={<DownloadOutlined />}
          onClick={() => onDownload?.(note)}
        >
          下载
        </Button>,
      ]}
    >
      <div className="space-y-2">
        <div className="flex justify-between items-start">
          <h3 className="font-semibold text-base truncate flex-1 pr-2">{note.title}</h3>
          {note.is_featured && (
            <Tag color="gold" icon={<StarOutlined />}>精选</Tag>
          )}
        </div>
        {note.subject && (
          <Tag color="blue">{note.subject}</Tag>
        )}
        {note.course_name && (
          <p className="text-gray-500 text-sm">课程: {note.course_name}</p>
        )}
        {note.description && (
          <p className="text-gray-600 text-sm line-clamp-2">{note.description}</p>
        )}
        <div className="flex items-center justify-between pt-2">
          <div className="flex items-center gap-2">
            <Rate disabled allowHalf value={note.rating} className="text-xs" />
            <span className="text-gray-500 text-sm">({note.rating_count})</span>
          </div>
          <span className="text-gray-400 text-sm">{formatFileSize(note.file_size)}</span>
        </div>
        {note.uploader && (
          <div className="flex items-center gap-2 text-gray-500 text-sm pt-1">
            <Avatar size="small" src={note.uploader.avatar} icon={<UserOutlined />} />
            <span>{note.uploader.username}</span>
            <span className="text-gray-400">|</span>
            <span>👁️ {note.view_count}</span>
            <span>⬇️ {note.download_count}</span>
          </div>
        )}
      </div>
    </Card>
  );
};
