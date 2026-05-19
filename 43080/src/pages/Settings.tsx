import React, { useState } from 'react';
import { useSettingsStore } from '../store/useSettingsStore';
import { useBookStore } from '../store/useBookStore';
import { useReadingStore } from '../store/useReadingStore';
import { useNoteStore } from '../store/useNoteStore';
import { useToast } from '../components/ui/Toast';
import { useConfirm } from '../hooks/useConfirm';
import { ConfirmDialog } from '../components/ui/ConfirmDialog';
import { Button } from '../components/ui/Button';
import { Input } from '../components/ui/Input';
import { Select } from '../components/ui/Select';
import { Checkbox } from '../components/ui/Checkbox';
import {
  Bell,
  Moon,
  Download,
  Trash2,
  RefreshCw,
  FileJson,
  FileText,
  AlertTriangle,
  CheckCircle2,
} from 'lucide-react';
import { exportToJSON, exportToCSV } from '../utils/export';
import { logger } from '../utils/logger';
import type { Theme, ViewMode } from '../types/settings';

const Settings: React.FC = () => {
  const {
    notificationEnabled,
    notificationTime,
    theme,
    defaultView,
    setNotificationEnabled,
    setNotificationTime,
    setTheme,
    setDefaultView,
    resetSettings,
  } = useSettingsStore();
  const { books } = useBookStore();
  const { sessions, plans } = useReadingStore();
  const { notes, bookmarks } = useNoteStore();
  const { showToast } = useToast();
  const { confirmState, confirm, cancel } = useConfirm();
  const [exportFormat, setExportFormat] = useState<'json' | 'csv'>('json');

  const handleExport = () => {
    try {
      const data = {
        books,
        readingPlans: plans,
        readingSessions: sessions,
        notes,
        bookmarks,
        exportedAt: new Date().toISOString(),
        version: '1.0.0',
      };

      const fileName = `reading-tracker-export-${new Date().toISOString().split('T')[0]}`;

      if (exportFormat === 'json') {
        exportToJSON(data, fileName);
      } else {
        exportToCSV(data, fileName);
      }

      showToast('success', '数据导出成功');
      logger.info('Data exported', { format: exportFormat });
    } catch (error) {
      showToast('error', '数据导出失败');
      logger.error('Export failed', error);
    }
  };

  const handleResetSettings = () => {
    confirm(
      '重置设置',
      '确定要重置所有设置为默认值吗？这不会影响你的书籍数据和阅读记录。',
      () => {
        resetSettings();
        showToast('success', '设置已重置');
      }
    );
  };

  const handleClearAllData = () => {
    confirm(
      '清除所有数据',
      '这将删除所有书籍、阅读记录、笔记和书签。此操作不可撤销！',
      () => {
        localStorage.clear();
        window.location.reload();
      }
    );
  };

  const dataStats = [
    { label: '书籍数量', value: books.length },
    { label: '阅读计划', value: plans.length },
    { label: '阅读记录', value: sessions.length },
    { label: '笔记数量', value: notes.length },
    { label: '书签数量', value: bookmarks.length },
  ];

  const themeOptions: { value: Theme; label: string }[] = [
    { value: 'light', label: '浅色' },
    { value: 'dark', label: '深色' },
    { value: 'system', label: '跟随系统' },
  ];

  const viewOptions: { value: ViewMode; label: string }[] = [
    { value: 'grid', label: '网格视图' },
    { value: 'list', label: '列表视图' },
  ];

  return (
    <div className="max-w-3xl mx-auto space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-gray-900">设置</h1>
        <p className="text-gray-500 mt-1">管理应用设置和数据</p>
      </div>

      <div className="bg-white rounded-xl shadow-sm border p-6 space-y-6">
        <div>
          <h2 className="text-lg font-semibold text-gray-900 flex items-center gap-2">
            <Bell className="w-5 h-5 text-blue-500" />
            阅读提醒
          </h2>
          <p className="text-sm text-gray-500 mt-1">设置每日阅读提醒时间</p>

          <div className="mt-4 space-y-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="font-medium text-gray-900">启用每日提醒</p>
                <p className="text-sm text-gray-500">在指定时间发送阅读提醒通知</p>
              </div>
              <Checkbox
                checked={notificationEnabled}
                onChange={(e) => setNotificationEnabled(e.target.checked)}
              />
            </div>

            {notificationEnabled && (
              <div className="flex items-center gap-4">
                <label className="text-sm text-gray-600 whitespace-nowrap">提醒时间:</label>
                <Input
                  type="time"
                  value={notificationTime}
                  onChange={(e) => setNotificationTime(e.target.value)}
                  className="w-32"
                />
              </div>
            )}
          </div>
        </div>
      </div>

      <div className="bg-white rounded-xl shadow-sm border p-6 space-y-6">
        <div>
          <h2 className="text-lg font-semibold text-gray-900 flex items-center gap-2">
            <Moon className="w-5 h-5 text-purple-500" />
            外观设置
          </h2>

          <div className="mt-4 space-y-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="font-medium text-gray-900">主题</p>
                <p className="text-sm text-gray-500">选择应用的颜色主题</p>
              </div>
              <Select
                value={theme}
                onChange={(e) => setTheme(e.target.value as Theme)}
                className="w-32"
                options={themeOptions}
              />
            </div>

            <div className="flex items-center justify-between">
              <div>
                <p className="font-medium text-gray-900">默认视图</p>
                <p className="text-sm text-gray-500">书架的默认展示方式</p>
              </div>
              <Select
                value={defaultView}
                onChange={(e) => setDefaultView(e.target.value as ViewMode)}
                className="w-32"
                options={viewOptions}
              />
            </div>
          </div>
        </div>
      </div>

      <div className="bg-white rounded-xl shadow-sm border p-6 space-y-6">
        <div>
          <h2 className="text-lg font-semibold text-gray-900 flex items-center gap-2">
            <Download className="w-5 h-5 text-green-500" />
            数据导出
          </h2>
          <p className="text-sm text-gray-500 mt-1">导出你的阅读数据进行备份</p>

          <div className="mt-4 grid grid-cols-2 sm:grid-cols-5 gap-4">
            {dataStats.map((stat) => (
              <div key={stat.label} className="text-center p-3 bg-gray-50 rounded-lg">
                <p className="text-2xl font-bold text-gray-900">{stat.value}</p>
                <p className="text-xs text-gray-500 mt-1">{stat.label}</p>
              </div>
            ))}
          </div>

          <div className="mt-4 flex flex-wrap gap-4 items-center">
            <div className="flex items-center gap-2">
              <span className="text-sm text-gray-600">导出格式:</span>
              <div className="flex bg-gray-100 rounded-lg p-1">
                <button
                  onClick={() => setExportFormat('json')}
                  className={`flex items-center gap-1 px-3 py-1.5 rounded-md text-sm transition-colors ${
                    exportFormat === 'json'
                      ? 'bg-white text-gray-900 shadow-sm'
                      : 'text-gray-600 hover:text-gray-900'
                  }`}
                >
                  <FileJson className="w-4 h-4" />
                  JSON
                </button>
                <button
                  onClick={() => setExportFormat('csv')}
                  className={`flex items-center gap-1 px-3 py-1.5 rounded-md text-sm transition-colors ${
                    exportFormat === 'csv'
                      ? 'bg-white text-gray-900 shadow-sm'
                      : 'text-gray-600 hover:text-gray-900'
                  }`}
                >
                  <FileText className="w-4 h-4" />
                  CSV
                </button>
              </div>
            </div>

            <Button onClick={handleExport}>
              <Download className="w-4 h-4 mr-2" />
              导出数据
            </Button>
          </div>
        </div>
      </div>

      <div className="bg-white rounded-xl shadow-sm border p-6 space-y-6">
        <div>
          <h2 className="text-lg font-semibold text-gray-900 flex items-center gap-2">
            <RefreshCw className="w-5 h-5 text-orange-500" />
            重置与清除
          </h2>

          <div className="mt-4 space-y-4">
            <div className="flex items-center justify-between p-4 bg-gray-50 rounded-lg">
              <div>
                <p className="font-medium text-gray-900">重置设置</p>
                <p className="text-sm text-gray-500">将所有设置恢复为默认值</p>
              </div>
              <Button variant="outline" onClick={handleResetSettings}>
                <RefreshCw className="w-4 h-4 mr-2" />
                重置
              </Button>
            </div>

            <div className="flex items-center justify-between p-4 bg-red-50 rounded-lg border border-red-200">
              <div className="flex items-start gap-3">
                <AlertTriangle className="w-5 h-5 text-red-500 mt-0.5" />
                <div>
                  <p className="font-medium text-red-700">清除所有数据</p>
                  <p className="text-sm text-red-600">
                    删除所有书籍、阅读记录、笔记和书签。此操作不可撤销！
                  </p>
                </div>
              </div>
              <Button variant="ghost" onClick={handleClearAllData} className="text-red-600 hover:text-red-700 hover:bg-red-100">
                <Trash2 className="w-4 h-4 mr-2" />
                清除
              </Button>
            </div>
          </div>
        </div>
      </div>

      <div className="bg-white rounded-xl shadow-sm border p-6">
        <div className="flex items-center gap-3">
          <CheckCircle2 className="w-5 h-5 text-green-500" />
          <div>
            <p className="font-medium text-gray-900">数据存储</p>
            <p className="text-sm text-gray-500">
              所有数据都安全地存储在你的浏览器本地存储中，不会上传到任何服务器。
            </p>
          </div>
        </div>
      </div>

      <ConfirmDialog
        isOpen={confirmState.isOpen}
        title={confirmState.title}
        message={confirmState.message}
        onConfirm={confirmState.onConfirm}
        onClose={cancel}
      />
    </div>
  );
};

export default Settings;
