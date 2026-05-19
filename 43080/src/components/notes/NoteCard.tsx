import React, { useState } from 'react';
import { Edit2, Trash2, BookOpen, Highlighter, Bookmark, Calendar } from 'lucide-react';
import type { Note } from '../../types/note';
import { Card } from '../ui/Card';
import { Badge } from '../ui/Badge';
import { Button } from '../ui/Button';
import { useBookStore } from '../../store/useBookStore';
import { useNoteStore } from '../../store/useNoteStore';
import { useToast } from '../ui/Toast';
import { formatDateTime } from '../../utils/date';
import { Textarea } from '../ui/Textarea';

interface NoteCardProps {
  note: Note;
  showBookTitle?: boolean;
}

export const NoteCard: React.FC<NoteCardProps> = ({ note, showBookTitle = false }) => {
  const { getBookById } = useBookStore();
  const { updateNote, deleteNote } = useNoteStore();
  const { showToast } = useToast();
  const book = getBookById(note.bookId);
  const [isEditing, setIsEditing] = useState(false);
  const [editContent, setEditContent] = useState(note.content);

  const typeIcons = {
    note: BookOpen,
    highlight: Highlighter,
    bookmark: Bookmark,
  };

  const typeLabels = {
    note: '笔记',
    highlight: '高亮',
    bookmark: '书签',
  };

  const Icon = typeIcons[note.type];

  const handleSave = () => {
    if (editContent.trim()) {
      updateNote(note.id, { content: editContent });
      setIsEditing(false);
      showToast('success', '笔记已更新');
    }
  };

  const handleDelete = () => {
    deleteNote(note.id);
    showToast('success', '笔记已删除');
  };

  return (
    <Card className="p-4">
      <div className="flex items-start justify-between mb-3">
        <div className="flex items-center gap-2">
          <div className={`p-2 rounded-lg ${
            note.type === 'highlight' ? 'bg-yellow-100' :
            note.type === 'bookmark' ? 'bg-blue-100' : 'bg-green-100'
          }`}>
            <Icon className={`w-4 h-4 ${
              note.type === 'highlight' ? 'text-yellow-700' :
              note.type === 'bookmark' ? 'text-blue-700' : 'text-green-700'
            }`} />
          </div>
          <div>
            <div className="flex items-center gap-2">
              <Badge variant={
                note.type === 'highlight' ? 'warning' :
                note.type === 'bookmark' ? 'info' : 'success'
              } size="sm">
                {typeLabels[note.type]}
              </Badge>
              {note.page > 0 && (
                <span className="text-xs text-gray-500">第 {note.page} 页</span>
              )}
            </div>
            {showBookTitle && book && (
              <p className="text-xs text-gray-500 mt-0.5">{book.title}</p>
            )}
          </div>
        </div>
        <div className="flex items-center gap-1">
          <Button variant="ghost" size="sm" onClick={() => setIsEditing(!isEditing)}>
            <Edit2 className="w-4 h-4" />
          </Button>
          <Button variant="ghost" size="sm" onClick={handleDelete}>
            <Trash2 className="w-4 h-4 text-red-500" />
          </Button>
        </div>
      </div>

      {note.chapter && (
        <p className="text-xs text-gray-500 mb-2">{note.chapter}</p>
      )}

      {note.highlightColor && (
        <div
          className="mb-3 px-3 py-2 rounded text-sm"
          style={{ backgroundColor: note.highlightColor }}
        >
          <p className="text-gray-800">{note.content}</p>
        </div>
      )}

      {!note.highlightColor && !isEditing && (
        <p className="text-sm text-gray-700 whitespace-pre-wrap">{note.content}</p>
      )}

      {isEditing && (
        <div className="space-y-3">
          <Textarea
            value={editContent}
            onChange={(e) => setEditContent(e.target.value)}
            rows={4}
          />
          <div className="flex justify-end gap-2">
            <Button variant="outline" size="sm" onClick={() => setIsEditing(false)}>
              取消
            </Button>
            <Button size="sm" onClick={handleSave}>
              保存
            </Button>
          </div>
        </div>
      )}

      <div className="flex items-center gap-2 mt-3 text-xs text-gray-400">
        <Calendar className="w-3 h-3" />
        <span>{formatDateTime(note.updatedAt)}</span>
      </div>
    </Card>
  );
};
