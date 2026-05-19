import { create } from 'zustand';
import type { Settings, Theme, ViewMode } from '../types/settings';
import { STORAGE_KEYS, DEFAULT_SETTINGS } from '../constants/config';
import { safeStorage } from '../utils/storage';
import { logger } from '../utils/logger';
import { validateData, settingsSchema } from '../utils/validation';
import { scheduleDailyReminder, requestNotificationPermission } from '../services/notificationService';

interface SettingsState extends Settings {}

interface SettingsActions {
  setNotificationEnabled: (enabled: boolean) => Promise<void>;
  setNotificationTime: (time: string) => void;
  setTheme: (theme: Theme) => void;
  setDefaultView: (view: ViewMode) => void;
  updateSettings: (settings: Partial<Settings>) => void;
  resetSettings: () => void;
  applyTheme: () => void;
  loadSettings: () => void;
  saveSettings: () => void;
}

const initialState: SettingsState = {
  ...DEFAULT_SETTINGS,
};

export const useSettingsStore = create<SettingsState & SettingsActions>((set, get) => ({
  ...initialState,

  setNotificationEnabled: async (enabled) => {
    if (enabled) {
      const granted = await requestNotificationPermission();
      if (!granted) {
        logger.warn('Notification permission not granted');
        return;
      }
    }
    set({ notificationEnabled: enabled });
    scheduleDailyReminder(get().notificationTime, enabled);
    get().saveSettings();
    logger.info(`Notification ${enabled ? 'enabled' : 'disabled'}`);
  },

  setNotificationTime: (time) => {
    set({ notificationTime: time });
    if (get().notificationEnabled) {
      scheduleDailyReminder(time, true);
    }
    get().saveSettings();
    logger.info(`Notification time set to ${time}`);
  },

  setTheme: (theme) => {
    set({ theme });
    get().applyTheme();
    get().saveSettings();
    logger.info(`Theme set to ${theme}`);
  },

  setDefaultView: (view) => {
    set({ defaultView: view });
    get().saveSettings();
    logger.info(`Default view set to ${view}`);
  },

  updateSettings: (settings) => {
    set((state) => ({ ...state, ...settings }));
    get().applyTheme();
    get().saveSettings();
  },

  resetSettings: () => {
    set(DEFAULT_SETTINGS);
    get().applyTheme();
    get().saveSettings();
    logger.info('Settings reset to defaults');
  },

  applyTheme: () => {
    const { theme } = get();
    const root = document.documentElement;

    if (theme === 'dark') {
      root.classList.add('dark');
    } else if (theme === 'light') {
      root.classList.remove('dark');
    } else {
      if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
        root.classList.add('dark');
      } else {
        root.classList.remove('dark');
      }
    }
  },

  loadSettings: () => {
    try {
      const storedSettings = safeStorage.get<Settings>(STORAGE_KEYS.SETTINGS, DEFAULT_SETTINGS);
      const validated = validateData(settingsSchema, storedSettings);
      if (validated) {
        set(validated);
        get().applyTheme();
        if (validated.notificationEnabled) {
          scheduleDailyReminder(validated.notificationTime, true);
        }
        logger.info('Settings loaded', validated);
      } else {
        set(DEFAULT_SETTINGS);
      }
    } catch (error) {
      logger.error('Failed to load settings', error);
      set(DEFAULT_SETTINGS);
    }
  },

  saveSettings: () => {
    const { notificationEnabled, notificationTime, theme, defaultView } = get();
    const settings: Settings = {
      notificationEnabled,
      notificationTime,
      theme,
      defaultView,
    };
    const success = safeStorage.set(STORAGE_KEYS.SETTINGS, settings);
    if (!success) {
      logger.warn('Failed to save settings to storage');
    }
  },
}));
