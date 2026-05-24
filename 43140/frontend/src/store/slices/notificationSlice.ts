import { createSlice, PayloadAction, createAsyncThunk } from '@reduxjs/toolkit'
import { statisticsApi, Notification } from '@/api/jobs'

interface NotificationState {
  notifications: Notification[]
  loading: boolean
  error: string | null
}

const initialState: NotificationState = {
  notifications: [],
  loading: false,
  error: null,
}

export const fetchNotifications = createAsyncThunk(
  'notifications/fetch',
  async () => {
    const response = await statisticsApi.getNotifications()
    return response.data
  }
)

export const markAsRead = createAsyncThunk(
  'notifications/markAsRead',
  async (id: number) => {
    await statisticsApi.markNotificationRead(id)
    return id
  }
)

export const markAllAsRead = createAsyncThunk(
  'notifications/markAllAsRead',
  async () => {
    await statisticsApi.markAllNotificationsRead()
  }
)

const notificationSlice = createSlice({
  name: 'notifications',
  initialState,
  reducers: {
    addNotification: (state, action: PayloadAction<Notification>) => {
      state.notifications.unshift(action.payload)
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchNotifications.fulfilled, (state, action) => {
        state.notifications = action.payload
      })
      .addCase(markAsRead.fulfilled, (state, action) => {
        const id = action.payload
        state.notifications = state.notifications.map((n) =>
          n.id === id ? { ...n, read: true } : n
        )
      })
      .addCase(markAllAsRead.fulfilled, (state) => {
        state.notifications = state.notifications.map((n) => ({ ...n, read: true }))
      })
  },
})

export const { addNotification } = notificationSlice.actions

export const selectNotifications = (state: { notifications: NotificationState }) => state.notifications
export const selectUnreadCount = (state: { notifications: NotificationState }) =>
  state.notifications.notifications.filter((n) => !n.read).length

export default notificationSlice.reducer
