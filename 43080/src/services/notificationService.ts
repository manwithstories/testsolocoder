import { logger } from '../utils/logger';

type NotificationPermission = 'default' | 'granted' | 'denied';

export const requestNotificationPermission = async (): Promise<boolean> => {
  if (!('Notification' in window)) {
    logger.warn('Notifications are not supported in this browser');
    return false;
  }

  if (Notification.permission === 'granted') {
    return true;
  }

  if (Notification.permission === 'denied') {
    logger.warn('Notification permission has been denied');
    return false;
  }

  try {
    const permission = await Notification.requestPermission();
    return permission === 'granted';
  } catch (error) {
    logger.error('Failed to request notification permission', error);
    return false;
  }
};

export const getNotificationPermission = (): NotificationPermission => {
  if (!('Notification' in window)) {
    return 'denied';
  }
  return Notification.permission as NotificationPermission;
};

export const showNotification = (title: string, options?: NotificationOptions): void => {
  if (!('Notification' in window)) {
    logger.warn('Notifications are not supported');
    return;
  }

  if (Notification.permission !== 'granted') {
    logger.warn('Notification permission not granted');
    return;
  }

  try {
    new Notification(title, {
      icon: '/favicon.svg',
      badge: '/favicon.svg',
      ...options,
    });
    logger.info(`Notification shown: ${title}`);
  } catch (error) {
    logger.error('Failed to show notification', error);
  }
};

export const showReadingReminder = (): void => {
  showNotification('阅读时间到了！', {
    body: '今天的阅读还没完成哦，拿起书开始阅读吧 📚',
    tag: 'reading-reminder',
    requireInteraction: true,
  });
};

let reminderIntervalId: number | null = null;

export const scheduleDailyReminder = (time: string, enabled: boolean): void => {
  if (reminderIntervalId !== null) {
    clearInterval(reminderIntervalId);
    reminderIntervalId = null;
  }

  if (!enabled) {
    logger.info('Daily reminder disabled');
    return;
  }

  const [hours, minutes] = time.split(':').map(Number);
  if (isNaN(hours) || isNaN(minutes)) {
    logger.error(`Invalid time format: ${time}`);
    return;
  }

  const checkAndTrigger = () => {
    const now = new Date();
    if (now.getHours() === hours && now.getMinutes() === minutes) {
      showReadingReminder();
    }
  };

  checkAndTrigger();
  reminderIntervalId = window.setInterval(checkAndTrigger, 60000);
  logger.info(`Daily reminder scheduled at ${time}`);
};

export const cleanupNotificationService = (): void => {
  if (reminderIntervalId !== null) {
    clearInterval(reminderIntervalId);
    reminderIntervalId = null;
  }
};
