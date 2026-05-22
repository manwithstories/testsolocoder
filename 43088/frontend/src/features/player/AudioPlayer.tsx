import { useRef, useEffect, useCallback } from 'react'
import { useSelector, useDispatch } from 'react-redux'
import {
  Play,
  Pause,
  Volume2,
  VolumeX,
  SkipBack,
  SkipForward,
  Settings,
  Zap,
  ZapOff,
} from 'lucide-react'
import { RootState } from '@/app/store'
import {
  setPlaying,
  setCurrentTime,
  setDuration,
  setVolume,
  setPlaybackRate,
  toggleMute,
  toggleSkipSilence,
} from './playerSlice'
import { formatTime } from '@/utils/format'
import { useUpdatePlaybackProgressMutation } from '@/app/services/api'

export default function AudioPlayer() {
  const audioRef = useRef<HTMLAudioElement>(null)
  const analyserRef = useRef<AnalyserNode | null>(null)
  const sourceRef = useRef<MediaElementAudioSourceNode | null>(null)
  const audioContextRef = useRef<AudioContext | null>(null)
  const silenceStartRef = useRef<number | null>(null)
  const audioContextInitializedRef = useRef(false)
  const dispatch = useDispatch()

  const {
    currentEpisode,
    isPlaying,
    currentTime,
    duration,
    volume,
    playbackRate,
    isMuted,
    skipSilence,
    silenceThreshold,
    silenceMinDuration,
  } = useSelector((state: RootState) => state.player)

  const [updateProgress] = useUpdatePlaybackProgressMutation()

  const initAudioContext = useCallback(() => {
    if (audioContextInitializedRef.current || !audioRef.current) return

    try {
      audioContextRef.current = new (window.AudioContext || (window as any).webkitAudioContext)()
      sourceRef.current = audioContextRef.current.createMediaElementSource(audioRef.current)
      analyserRef.current = audioContextRef.current.createAnalyser()
      analyserRef.current.fftSize = 256
      sourceRef.current.connect(analyserRef.current)
      analyserRef.current.connect(audioContextRef.current.destination)
      audioContextInitializedRef.current = true
    } catch (err) {
      console.warn('Failed to initialize audio context for silence detection:', err)
    }
  }, [])

  useEffect(() => {
    if (audioRef.current && currentEpisode) {
      audioRef.current.load()
      audioRef.current.play().catch(() => dispatch(setPlaying(false)))
    }
  }, [currentEpisode, dispatch])

  useEffect(() => {
    if (audioRef.current) {
      if (isPlaying) {
        audioRef.current.play().catch(() => dispatch(setPlaying(false)))
      } else {
        audioRef.current.pause()
      }
    }
  }, [isPlaying, dispatch])

  useEffect(() => {
    if (audioRef.current) {
      audioRef.current.volume = isMuted ? 0 : volume
    }
  }, [volume, isMuted])

  useEffect(() => {
    if (audioRef.current) {
      audioRef.current.playbackRate = playbackRate
    }
  }, [playbackRate])

  useEffect(() => {
    if (skipSilence && !audioContextInitializedRef.current) {
      initAudioContext()
    }
  }, [skipSilence, initAudioContext])

  useEffect(() => {
    return () => {
      if (audioContextRef.current) {
        try {
          audioContextRef.current.close()
        } catch (err) {
          console.warn('Error closing audio context:', err)
        }
        audioContextRef.current = null
      }
      sourceRef.current = null
      analyserRef.current = null
      audioContextInitializedRef.current = false
    }
  }, [])

  const detectSilence = useCallback(() => {
    if (!skipSilence || !analyserRef.current || !audioRef.current) return

    const bufferLength = analyserRef.current.frequencyBinCount
    const dataArray = new Uint8Array(bufferLength)
    analyserRef.current.getByteFrequencyData(dataArray)

    let sum = 0
    for (let i = 0; i < bufferLength; i++) {
      sum += dataArray[i]
    }
    const average = sum / bufferLength
    const normalizedAverage = average / 255

    if (normalizedAverage < silenceThreshold) {
      if (silenceStartRef.current === null) {
        silenceStartRef.current = audioRef.current.currentTime
      } else {
        const silenceDuration = audioRef.current.currentTime - silenceStartRef.current
        if (silenceDuration >= silenceMinDuration) {
          audioRef.current.currentTime += 0.1
        }
      }
    } else {
      silenceStartRef.current = null
    }
  }, [skipSilence, silenceThreshold, silenceMinDuration])

  useEffect(() => {
    let animationId: number

    const animate = () => {
      detectSilence()
      animationId = requestAnimationFrame(animate)
    }

    if (skipSilence && isPlaying && audioContextInitializedRef.current) {
      animationId = requestAnimationFrame(animate)
    }

    return () => {
      if (animationId) {
        cancelAnimationFrame(animationId)
      }
    }
  }, [skipSilence, isPlaying, detectSilence])

  const handleTimeUpdate = useCallback(() => {
    if (audioRef.current) {
      dispatch(setCurrentTime(audioRef.current.currentTime))
    }
  }, [dispatch])

  const handleLoadedMetadata = useCallback(() => {
    if (audioRef.current) {
      dispatch(setDuration(audioRef.current.duration))
    }
  }, [dispatch])

  const handleSeek = (e: React.ChangeEvent<HTMLInputElement>) => {
    const time = parseFloat(e.target.value)
    if (audioRef.current) {
      audioRef.current.currentTime = time
      dispatch(setCurrentTime(time))
    }
  }

  const handleSkip = (seconds: number) => {
    if (audioRef.current) {
      audioRef.current.currentTime += seconds
    }
  }

  const handleEnded = useCallback(() => {
    dispatch(setPlaying(false))
    silenceStartRef.current = null
    if (currentEpisode) {
      updateProgress({
        id: currentEpisode.id,
        current_time: duration,
        duration: duration,
      })
    }
  }, [currentEpisode, duration, dispatch, updateProgress])

  useEffect(() => {
    if (currentEpisode && currentTime > 0) {
      const timer = setInterval(() => {
        updateProgress({
          id: currentEpisode.id,
          current_time: currentTime,
          duration: duration,
        })
      }, 10000)

      return () => clearInterval(timer)
    }
  }, [currentEpisode, currentTime, duration, updateProgress])

  if (!currentEpisode) {
    return null
  }

  const progress = duration > 0 ? (currentTime / duration) * 100 : 0

  return (
    <div className="audio-player">
      <audio
        ref={audioRef}
        src={currentEpisode.audio_url}
        onTimeUpdate={handleTimeUpdate}
        onLoadedMetadata={handleLoadedMetadata}
        onEnded={handleEnded}
        crossOrigin="anonymous"
      />

      <div className="max-w-7xl mx-auto px-4 py-3">
        <div className="flex items-center gap-4">
          <img
            src={currentEpisode.podcast?.cover_image || '/placeholder.png'}
            alt={currentEpisode.title}
            className="w-14 h-14 rounded-lg object-cover"
          />

          <div className="flex-1 min-w-0">
            <h4 className="font-medium text-gray-900 truncate">{currentEpisode.title}</h4>
            <p className="text-sm text-gray-500 truncate">{currentEpisode.podcast?.title}</p>
          </div>

          <div className="flex items-center gap-2">
            <button
              onClick={() => handleSkip(-10)}
              className="p-2 rounded-full hover:bg-gray-100 transition-colors"
              title="后退10秒"
            >
              <SkipBack className="w-5 h-5" />
            </button>

            <button
              onClick={() => dispatch(setPlaying(!isPlaying))}
              className="p-3 bg-indigo-600 text-white rounded-full hover:bg-indigo-700 transition-colors"
            >
              {isPlaying ? <Pause className="w-6 h-6" /> : <Play className="w-6 h-6 ml-0.5" />}
            </button>

            <button
              onClick={() => handleSkip(30)}
              className="p-2 rounded-full hover:bg-gray-100 transition-colors"
              title="前进30秒"
            >
              <SkipForward className="w-5 h-5" />
            </button>
          </div>

          <div className="flex-1 max-w-xl">
            <div className="flex items-center gap-3">
              <span className="text-sm text-gray-500 w-12 text-right">{formatTime(currentTime)}</span>
              <input
                type="range"
                min={0}
                max={duration || 0}
                step={1}
                value={currentTime}
                onChange={handleSeek}
                className="flex-1 h-2 bg-gray-200 rounded-full appearance-none cursor-pointer accent-indigo-600"
                style={{
                  background: `linear-gradient(to right, #6366f1 ${progress}%, #e5e7eb ${progress}%)`,
                }}
              />
              <span className="text-sm text-gray-500 w-12">{formatTime(duration)}</span>
            </div>
          </div>

          <div className="flex items-center gap-3">
            <button
              onClick={() => dispatch(toggleSkipSilence())}
              className={`p-2 rounded-full transition-colors ${
                skipSilence
                  ? 'bg-indigo-100 text-indigo-600'
                  : 'hover:bg-gray-100 text-gray-600'
              }`}
              title={skipSilence ? '已开启跳过静音' : '跳过静音片段'}
            >
              {skipSilence ? <Zap className="w-5 h-5" /> : <ZapOff className="w-5 h-5" />}
            </button>

            <div className="flex items-center gap-2">
              <button
                onClick={() => dispatch(toggleMute())}
                className="p-2 rounded-full hover:bg-gray-100 transition-colors"
              >
                {isMuted || volume === 0 ? <VolumeX className="w-5 h-5" /> : <Volume2 className="w-5 h-5" />}
              </button>
              <input
                type="range"
                min={0}
                max={1}
                step={0.01}
                value={isMuted ? 0 : volume}
                onChange={(e) => dispatch(setVolume(parseFloat(e.target.value)))}
                className="w-20 h-2 bg-gray-200 rounded-full appearance-none cursor-pointer accent-indigo-600"
              />
            </div>

            <div className="relative group">
              <button className="p-2 rounded-full hover:bg-gray-100 transition-colors">
                <Settings className="w-5 h-5" />
              </button>
              <div className="absolute bottom-full right-0 mb-2 hidden group-hover:block bg-white rounded-lg shadow-lg border border-gray-200 p-2 w-32">
                <p className="text-xs text-gray-500 mb-2 px-2">播放速度</p>
                {[0.5, 0.75, 1, 1.25, 1.5, 2].map((rate) => (
                  <button
                    key={rate}
                    onClick={() => dispatch(setPlaybackRate(rate))}
                    className={`w-full text-left px-2 py-1 rounded text-sm transition-colors ${
                      playbackRate === rate ? 'bg-indigo-50 text-indigo-600' : 'hover:bg-gray-50'
                    }`}
                  >
                    {rate}x
                  </button>
                ))}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
