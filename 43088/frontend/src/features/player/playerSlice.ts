import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import type { Episode, PlayerState } from '@/types'

const initialState: PlayerState = {
  currentEpisode: null,
  isPlaying: false,
  currentTime: 0,
  duration: 0,
  volume: 1,
  playbackRate: 1,
  isMuted: false,
  skipSilence: false,
  silenceThreshold: 0.02,
  silenceMinDuration: 1.5,
}

const playerSlice = createSlice({
  name: 'player',
  initialState,
  reducers: {
    setCurrentEpisode: (state, action: PayloadAction<Episode>) => {
      state.currentEpisode = action.payload
      state.currentTime = 0
      state.duration = action.payload.duration
    },
    setPlaying: (state, action: PayloadAction<boolean>) => {
      state.isPlaying = action.payload
    },
    setCurrentTime: (state, action: PayloadAction<number>) => {
      state.currentTime = action.payload
    },
    setDuration: (state, action: PayloadAction<number>) => {
      state.duration = action.payload
    },
    setVolume: (state, action: PayloadAction<number>) => {
      state.volume = Math.max(0, Math.min(1, action.payload))
      state.isMuted = state.volume === 0
    },
    setPlaybackRate: (state, action: PayloadAction<number>) => {
      state.playbackRate = action.payload
    },
    toggleMute: (state) => {
      state.isMuted = !state.isMuted
    },
    playPause: (state) => {
      state.isPlaying = !state.isPlaying
    },
    resetPlayer: (state) => {
      state.currentEpisode = null
      state.isPlaying = false
      state.currentTime = 0
      state.duration = 0
    },
    toggleSkipSilence: (state) => {
      state.skipSilence = !state.skipSilence
    },
    setSilenceThreshold: (state, action: PayloadAction<number>) => {
      state.silenceThreshold = action.payload
    },
    setSilenceMinDuration: (state, action: PayloadAction<number>) => {
      state.silenceMinDuration = action.payload
    },
  },
})

export const {
  setCurrentEpisode,
  setPlaying,
  setCurrentTime,
  setDuration,
  setVolume,
  setPlaybackRate,
  toggleMute,
  playPause,
  resetPlayer,
  toggleSkipSilence,
  setSilenceThreshold,
  setSilenceMinDuration,
} = playerSlice.actions

export default playerSlice.reducer
