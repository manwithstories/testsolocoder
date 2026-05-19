export type Theme = 'light' | 'dark' | 'system';
export type ViewMode = 'grid' | 'list';

export interface Settings {
  notificationEnabled: boolean;
  notificationTime: string;
  theme: Theme;
  defaultView: ViewMode;
}

export interface ExportOptions {
  format: 'json' | 'csv';
  includeBooks: boolean;
  includeReadingSessions: boolean;
  includeNotes: boolean;
  includeStats: boolean;
}
