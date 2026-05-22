import { configureStore } from '@reduxjs/toolkit'
import { api } from './services/api'
import playerReducer from '@/features/player/playerSlice'

export const store = configureStore({
  reducer: {
    [api.reducerPath]: api.reducer,
    player: playerReducer,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware().concat(api.middleware),
})

export type RootState = ReturnType<typeof store.getState>
export type AppDispatch = typeof store.dispatch
