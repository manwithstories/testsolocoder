<template>
  <div class="online-guide" v-loading="loading">
    <div v-if="collection" class="guide-container">
      <div class="collection-header">
        <div class="collection-image">
          <img :src="collection.image_url || '/placeholder.svg'" :alt="getLocalizedName()" />
        </div>
        <div class="collection-info">
          <h1 class="collection-name">{{ getLocalizedName() }}</h1>
          <div class="collection-meta">
            <span class="meta-item">
              <el-icon><Collection /></el-icon>
              {{ collection.code }}
            </span>
            <span class="meta-item" v-if="collection.era">
              <el-icon><Clock /></el-icon>
              {{ collection.era }}
            </span>
            <span class="meta-item" v-if="collection.material">
              <el-icon><Box /></el-icon>
              {{ collection.material }}
            </span>
          </div>
        </div>
      </div>

      <div class="control-panel">
        <div class="control-group">
          <span class="control-label">{{ t('language') }}:</span>
          <el-radio-group v-model="currentLang" @change="onLangChange" size="default">
            <el-radio-button value="zh">中文</el-radio-button>
            <el-radio-button value="en">English</el-radio-button>
            <el-radio-button value="ja">日本語</el-radio-button>
          </el-radio-group>
        </div>
        <div class="control-group">
          <span class="control-label">{{ t('guideType') }}:</span>
          <el-radio-group v-model="guideMode" @change="onGuideModeChange" size="default">
            <el-radio-button value="text">
              <el-icon><Document /></el-icon>
              {{ t('text') }}
            </el-radio-button>
            <el-radio-button value="audio">
              <el-icon><Microphone /></el-icon>
              {{ t('audio') }}
            </el-radio-button>
          </el-radio-group>
        </div>
      </div>

      <div class="guide-content-section">
        <div v-if="guideMode === 'text'" class="text-guide">
          <div class="text-content">
            <div class="content-header">
              <el-icon size="24" color="#409eff"><Reading /></el-icon>
              <h2>{{ t('guideContent') }}</h2>
            </div>
            <div class="content-body" v-if="currentGuide">
              <p class="guide-text">{{ currentGuide.content }}</p>
            </div>
            <div class="content-body empty" v-else>
              <el-empty :description="t('noContent')" />
            </div>
          </div>
        </div>

        <div v-else class="audio-guide">
          <div class="audio-player-card">
            <div class="player-header">
              <div class="player-icon">
                <el-icon size="48" color="#67c23a"><Headset /></el-icon>
              </div>
              <div class="player-info">
                <h3>{{ t('audioGuide') }}</h3>
                <p>{{ getLocalizedName() }}</p>
              </div>
            </div>

            <div class="player-controls" v-if="currentGuide?.audio_url">
              <audio
                ref="audioPlayer"
                :src="currentGuide.audio_url"
                @timeupdate="onTimeUpdate"
                @loadedmetadata="onLoadedMetadata"
                @ended="onAudioEnded"
              ></audio>

              <div class="progress-section">
                <span class="time-text">{{ formatTime(currentTime) }}</span>
                <el-slider
                  v-model="playProgress"
                  :max="duration"
                  @change="onProgressChange"
                  :show-tooltip="false"
                  class="progress-slider"
                />
                <span class="time-text">{{ formatTime(duration) }}</span>
              </div>

              <div class="control-buttons">
                <el-button circle size="large" @click="rewind">
                  <el-icon size="20"><RefreshLeft /></el-icon>
                </el-button>
                <el-button type="primary" circle size="large" @click="togglePlay">
                  <el-icon v-if="isPlaying" size="24"><VideoPause /></el-icon>
                  <el-icon v-else size="24"><VideoPlay /></el-icon>
                </el-button>
                <el-button circle size="large" @click="forward">
                  <el-icon size="20"><RefreshRight /></el-icon>
                </el-button>
                <div class="volume-control">
                  <el-icon size="20"><Volume /></el-icon>
                  <el-slider
                    v-model="volume"
                    :min="0"
                    :max="100"
                    @change="onVolumeChange"
                    class="volume-slider"
                  />
                </div>
                <div class="speed-control">
                  <el-dropdown @command="onSpeedChange">
                    <el-button size="small">{{ playbackRate }}x</el-button>
                    <template #dropdown>
                      <el-dropdown-menu>
                        <el-dropdown-item command="0.5">0.5x</el-dropdown-item>
                        <el-dropdown-item command="0.75">0.75x</el-dropdown-item>
                        <el-dropdown-item command="1">1x</el-dropdown-item>
                        <el-dropdown-item command="1.25">1.25x</el-dropdown-item>
                        <el-dropdown-item command="1.5">1.5x</el-dropdown-item>
                        <el-dropdown-item command="2">2x</el-dropdown-item>
                      </el-dropdown-menu>
                    </template>
                  </el-dropdown>
                </div>
              </div>
            </div>

            <div class="player-empty" v-else>
              <el-empty :description="t('noAudio')" />
            </div>
          </div>
        </div>
      </div>

      <div class="guide-list-section">
        <h3 class="section-title">{{ t('otherCollections') }}</h3>
        <div class="guide-list">
          <div
            v-for="content in allGuideContents"
            :key="content.id"
            class="guide-item"
            :class="{ active: currentGuide?.id === content.id }"
            @click="selectGuide(content)"
          >
            <div class="item-thumb">
              <img :src="content.collection?.image_url || '/placeholder.svg'" :alt="content.collection?.name" />
            </div>
            <div class="item-info">
              <h4>{{ content.collection?.name }}</h4>
              <p>{{ getLanguageName(content.language) }}</p>
            </div>
            <div class="item-status">
              <el-icon v-if="content.audio_url" color="#67c23a"><Headset /></el-icon>
              <el-icon v-else color="#909399"><Document /></el-icon>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import * as collectionApi from '@/api/collection'
import * as guideApi from '@/api/guide'
import type { Collection, GuideContent } from '@/types'

const route = useRoute()
const collectionId = Number(route.params.id)

const loading = ref(false)
const collection = ref<Collection | null>(null)
const allGuideContents = ref<GuideContent[]>([])
const currentLang = ref('zh')
const guideMode = ref<'text' | 'audio'>('text')
const audioPlayer = ref<HTMLAudioElement | null>(null)
const isPlaying = ref(false)
const currentTime = ref(0)
const duration = ref(0)
const volume = ref(80)
const playbackRate = ref(1)
const playProgress = ref(0)

const translations = reactive({
  zh: {
    language: '语言',
    guideType: '导览方式',
    text: '文字',
    audio: '语音',
    guideContent: '藏品讲解',
    noContent: '暂无该语言的导览内容',
    audioGuide: '语音导览',
    noAudio: '暂无该语言的语音导览',
    otherCollections: '相关导览列表',
    play: '播放',
    pause: '暂停'
  },
  en: {
    language: 'Language',
    guideType: 'Guide Type',
    text: 'Text',
    audio: 'Audio',
    guideContent: 'Collection Guide',
    noContent: 'No guide content available in this language',
    audioGuide: 'Audio Guide',
    noAudio: 'No audio guide available in this language',
    otherCollections: 'Guide List',
    play: 'Play',
    pause: 'Pause'
  },
  ja: {
    language: '言語',
    guideType: 'ガイドタイプ',
    text: 'テキスト',
    audio: '音声',
    guideContent: '収蔵品ガイド',
    noContent: 'この言語のガイドコンテンツはありません',
    audioGuide: '音声ガイド',
    noAudio: 'この言語の音声ガイドはありません',
    otherCollections: 'ガイドリスト',
    play: '再生',
    pause: '一時停止'
  }
})

const t = (key: string) => {
  return (translations as any)[currentLang.value]?.[key] || (translations.zh as any)[key]
}

const currentGuide = computed(() => {
  return allGuideContents.value.find(
    g => g.language === currentLang.value && g.collection_id === collectionId
  )
})

const getLocalizedName = () => {
  return collection.value?.name || ''
}

const getLanguageName = (lang: string) => {
  const names: Record<string, string> = { zh: '中文', en: 'English', ja: '日本語' }
  return names[lang] || lang
}

const formatTime = (seconds: number) => {
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
}

const fetchCollection = async () => {
  try {
    loading.value = true
    const res = await collectionApi.getCollection(collectionId)
    collection.value = res.data
  } catch (e) {
    console.error(e)
    ElMessage.error('获取藏品信息失败')
  } finally {
    loading.value = false
  }
}

const fetchGuideContents = async () => {
  try {
    const res = await guideApi.listGuideContents({
      collection_id: collectionId
    })
    allGuideContents.value = res.data
  } catch (e) {
    console.error(e)
  }
}

const onLangChange = () => {
  if (isPlaying.value) {
    stopAudio()
  }
}

const onGuideModeChange = () => {
  if (isPlaying.value) {
    stopAudio()
  }
}

const togglePlay = () => {
  if (!audioPlayer.value || !currentGuide.value?.audio_url) return

  if (isPlaying.value) {
    audioPlayer.value.pause()
  } else {
    audioPlayer.value.play()
  }
  isPlaying.value = !isPlaying.value
}

const stopAudio = () => {
  if (audioPlayer.value) {
    audioPlayer.value.pause()
    audioPlayer.value.currentTime = 0
    isPlaying.value = false
    currentTime.value = 0
    playProgress.value = 0
  }
}

const rewind = () => {
  if (audioPlayer.value) {
    audioPlayer.value.currentTime = Math.max(0, audioPlayer.value.currentTime - 10)
  }
}

const forward = () => {
  if (audioPlayer.value) {
    audioPlayer.value.currentTime = Math.min(
      audioPlayer.value.duration || 0,
      audioPlayer.value.currentTime + 10
    )
  }
}

const onTimeUpdate = () => {
  if (audioPlayer.value) {
    currentTime.value = audioPlayer.value.currentTime
    playProgress.value = audioPlayer.value.currentTime
  }
}

const onLoadedMetadata = () => {
  if (audioPlayer.value) {
    duration.value = audioPlayer.value.duration || 0
  }
}

const onAudioEnded = () => {
  isPlaying.value = false
  currentTime.value = 0
  playProgress.value = 0
}

const onProgressChange = (value: number) => {
  if (audioPlayer.value) {
    audioPlayer.value.currentTime = value
  }
}

const onVolumeChange = (value: number) => {
  if (audioPlayer.value) {
    audioPlayer.value.volume = value / 100
  }
}

const onSpeedChange = (value: string | number) => {
  playbackRate.value = Number(value)
  if (audioPlayer.value) {
    audioPlayer.value.playbackRate = playbackRate.value
  }
}

const selectGuide = (content: GuideContent) => {
  if (isPlaying.value) {
    stopAudio()
  }
  currentLang.value = content.language
}

watch(volume, (newVal) => {
  if (audioPlayer.value) {
    audioPlayer.value.volume = newVal / 100
  }
})

onMounted(() => {
  fetchCollection()
  fetchGuideContents()
})

onUnmounted(() => {
  stopAudio()
})
</script>

<style scoped lang="scss">
.online-guide {
  min-height: 100vh;
  background: linear-gradient(180deg, #f0f2f5 0%, #fff 100%);
  padding-bottom: 40px;
}

.guide-container {
  max-width: 900px;
  margin: 0 auto;
  padding: 20px;
}

.collection-header {
  background: #fff;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  margin-bottom: 24px;

  .collection-image {
    width: 100%;
    height: 300px;
    overflow: hidden;

    img {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }
  }

  .collection-info {
    padding: 24px;

    .collection-name {
      font-size: 28px;
      font-weight: 700;
      margin-bottom: 16px;
      color: #303133;
    }

    .collection-meta {
      display: flex;
      flex-wrap: wrap;
      gap: 24px;

      .meta-item {
        display: flex;
        align-items: center;
        gap: 6px;
        color: #606266;
        font-size: 14px;

        .el-icon {
          color: #409eff;
        }
      }
    }
  }
}

.control-panel {
  background: #fff;
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 24px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
  display: flex;
  flex-wrap: wrap;
  gap: 32px;

  .control-group {
    display: flex;
    align-items: center;
    gap: 12px;

    .control-label {
      font-weight: 600;
      color: #303133;
      font-size: 15px;
    }
  }
}

.guide-content-section {
  margin-bottom: 32px;
}

.text-guide {
  .text-content {
    background: #fff;
    border-radius: 12px;
    padding: 32px;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);

    .content-header {
      display: flex;
      align-items: center;
      gap: 12px;
      margin-bottom: 24px;
      padding-bottom: 16px;
      border-bottom: 2px solid #f0f2f5;

      h2 {
        font-size: 22px;
        font-weight: 600;
        margin: 0;
      }
    }

    .content-body {
      .guide-text {
        font-size: 16px;
        line-height: 2;
        color: #303133;
        text-indent: 2em;
      }

      &.empty {
        padding: 40px 0;
      }
    }
  }
}

.audio-guide {
  .audio-player-card {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border-radius: 16px;
    padding: 32px;
    color: #fff;
    box-shadow: 0 8px 32px rgba(102, 126, 234, 0.3);

    .player-header {
      display: flex;
      align-items: center;
      gap: 20px;
      margin-bottom: 32px;

      .player-icon {
        width: 80px;
        height: 80px;
        background: rgba(255, 255, 255, 0.2);
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
      }

      .player-info {
        h3 {
          font-size: 24px;
          margin: 0 0 8px;
        }

        p {
          opacity: 0.9;
          margin: 0;
        }
      }
    }

    .player-controls {
      .progress-section {
        display: flex;
        align-items: center;
        gap: 16px;
        margin-bottom: 24px;

        .time-text {
          font-size: 14px;
          font-family: monospace;
          min-width: 50px;
        }

        .progress-slider {
          flex: 1;

          :deep(.el-slider__runway) {
            background: rgba(255, 255, 255, 0.3);
          }

          :deep(.el-slider__bar) {
            background: #fff;
          }

          :deep(.el-slider__button) {
            border-color: #fff;
          }
        }
      }

      .control-buttons {
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 20px;

        .el-button {
          background: rgba(255, 255, 255, 0.2);
          border-color: rgba(255, 255, 255, 0.3);
          color: #fff;

          &:hover {
            background: rgba(255, 255, 255, 0.3);
          }

          &.el-button--primary {
            background: #fff;
            color: #667eea;
            border-color: #fff;

            &:hover {
              background: rgba(255, 255, 255, 0.9);
            }
          }
        }

        .volume-control {
          display: flex;
          align-items: center;
          gap: 8px;
          margin-left: 20px;

          .el-icon {
            color: #fff;
          }

          .volume-slider {
            width: 80px;

            :deep(.el-slider__runway) {
              background: rgba(255, 255, 255, 0.3);
            }

            :deep(.el-slider__bar) {
              background: #fff;
            }

            :deep(.el-slider__button) {
              width: 12px;
              height: 12px;
              border-color: #fff;
            }
          }
        }

        .speed-control {
          margin-left: auto;

          .el-button {
            background: rgba(255, 255, 255, 0.2);
            border-color: rgba(255, 255, 255, 0.3);
            color: #fff;

            &:hover {
              background: rgba(255, 255, 255, 0.3);
            }
          }

          :deep(.el-dropdown-menu__item) {
            text-align: center;
          }
        }
      }
    }

    .player-empty {
      padding: 40px 0;

      :deep(.el-empty__description) {
      }
    }
  }
}

.guide-list-section {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);

  .section-title {
    font-size: 18px;
    font-weight: 600;
    margin: 0 0 20px;
    padding-bottom: 12px;
    border-bottom: 1px solid #f0f2f5;
  }

  .guide-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .guide-item {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 12px;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.3s ease;
    border: 2px solid transparent;

    &:hover {
      background: #f5f7fa;
    }

    &.active {
      background: #ecf5ff;
      border-color: #409eff;
    }

    .item-thumb {
      width: 60px;
      height: 60px;
      border-radius: 8px;
      overflow: hidden;
      flex-shrink: 0;

      img {
        width: 100%;
        height: 100%;
        object-fit: cover;
      }
    }

    .item-info {
      flex: 1;
      min-width: 0;

      h4 {
        font-size: 15px;
        font-weight: 500;
        margin: 0 0 4px;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
      }

      p {
        font-size: 13px;
        color: #909399;
        margin: 0;
      }
    }

    .item-status {
      flex-shrink: 0;
    }
  }
}

@media (max-width: 768px) {
  .collection-header {
    .collection-image {
      height: 200px;
    }

    .collection-info {
      padding: 16px;

      .collection-name {
        font-size: 22px;
      }

      .collection-meta {
        gap: 16px;
      }
    }
  }

  .control-panel {
    flex-direction: column;
    gap: 16px;
  }

  .audio-player-card {
    padding: 20px;

    .player-header {
      flex-direction: column;
      text-align: center;

      .player-icon {
        width: 60px;
        height: 60px;
      }
    }

    .player-controls {
      .control-buttons {
        flex-wrap: wrap;

        .volume-control {
          margin-left: 0;
          width: 100%;
          justify-content: center;
        }

        .speed-control {
          margin-left: 0;
          width: 100%;
          display: flex;
          justify-content: center;
        }
      }
    }
  }

  .text-guide {
    .text-content {
      padding: 20px;
    }
  }
}
</style>
