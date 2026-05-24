import { useState, useEffect, useRef } from 'react'
import { useParams } from 'react-router-dom'
import { toast } from 'sonner'
import { videoApi, homeworkApi, bookingApi } from '@/services/api'
import { useAuthStore } from '@/store/auth'
import { Video, Phone, Mic, VideoOff, MicOff, Monitor, Settings, MoreVertical, BookOpen, CheckCircle } from 'lucide-react'

export default function VideoSessionPage() {
  const { sessionId } = useParams<{ sessionId: string }>()
  const { user } = useAuthStore()
  const [session, setSession] = useState<any>(null)
  const [loading, setLoading] = useState(true)
  const [isMuted, setIsMuted] = useState(false)
  const [isVideoOn, setIsVideoOn] = useState(true)
  const [isScreenSharing, setIsScreenSharing] = useState(false)
  const [quality, setQuality] = useState<any>(null)
  const [showSettings, setShowSettings] = useState(false)
  const [sessionEnded, setSessionEnded] = useState(false)
  const [showHomeworkModal, setShowHomeworkModal] = useState(false)
  const [homeworkForm, setHomeworkForm] = useState({
    title: '',
    description: '',
    dueDate: '',
    maxScore: 100,
  })
  const videoRef = useRef<HTMLVideoElement>(null)

  useEffect(() => {
    if (sessionId) {
      loadSession()
    }
  }, [sessionId])

  const loadSession = async () => {
    try {
      setLoading(true)
      const res = await videoApi.getSession(sessionId!)
      setSession(res.data)
    } catch (error) {
      console.error('Failed to load session:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleEndSession = async () => {
    if (confirm('确定要结束课程吗？')) {
      try {
        await videoApi.endSession({ sessionId: sessionId! })
        toast.success('课程已结束')
        setSessionEnded(true)
      } catch (error: any) {
        toast.error(error.response?.data?.error || '结束课程失败')
      }
    }
  }

  const handleCompleteBooking = async () => {
    if (!session?.bookingId) return
    if (confirm('确定要完成这个课程吗？完成后将自动结算费用。')) {
      try {
        await bookingApi.complete(session.bookingId)
        toast.success('课程已完成，费用已结算')
      } catch (error: any) {
        toast.error(error.response?.data?.error || '完成课程失败')
      }
    }
  }

  const handleCreateHomework = async () => {
    if (!session?.bookingId || !homeworkForm.title || !homeworkForm.description || !homeworkForm.dueDate) {
      toast.error('请填写完整的作业信息')
      return
    }

    try {
      await homeworkApi.create({
        bookingId: session.bookingId,
        subjectId: session.subjectId || '',
        title: homeworkForm.title,
        description: homeworkForm.description,
        dueDate: new Date(homeworkForm.dueDate).toISOString(),
        maxScore: homeworkForm.maxScore,
      })
      toast.success('作业已布置成功')
      setShowHomeworkModal(false)
    } catch (error: any) {
      toast.error(error.response?.data?.error || '布置作业失败')
    }
  }

  const handleQualityCheck = async () => {
    try {
      const res = await videoApi.getQuality(sessionId!)
      setQuality(res.data)
      toast.success('质量检查完成')
    } catch (error) {
      toast.error('质量检查失败')
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-screen bg-gray-900">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-white mx-auto mb-4"></div>
          <p className="text-white">加载视频会话...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="h-screen bg-gray-900 flex flex-col">
      <div className="flex-1 relative">
        <div className="absolute inset-0 flex items-center justify-center">
          <div className="w-full max-w-4xl aspect-video bg-gray-800 rounded-lg overflow-hidden">
            <video
              ref={videoRef}
              autoPlay
              playsInline
              muted={isMuted}
              className="w-full h-full object-cover"
            />
            {!isVideoOn && (
              <div className="absolute inset-0 flex items-center justify-center bg-gray-800">
                <div className="text-center">
                  <div className="w-24 h-24 rounded-full bg-gray-700 flex items-center justify-center mx-auto mb-4">
                    <VideoOff className="h-12 w-12 text-gray-500" />
                  </div>
                  <p className="text-gray-400">视频已关闭</p>
                </div>
              </div>
            )}
          </div>
        </div>

        <div className="absolute top-4 right-4 w-48 h-32 bg-gray-800 rounded-lg overflow-hidden border-2 border-gray-700">
          <div className="w-full h-full flex items-center justify-center">
            <div className="text-white text-sm">您的画面</div>
          </div>
        </div>

        <div className="absolute top-4 left-4 bg-black bg-opacity-50 rounded-lg px-4 py-2">
          <div className="flex items-center gap-2 text-white">
            <Video className="h-4 w-4" />
            <span className="text-sm">会话ID: {session?.sessionId}</span>
          </div>
          {session?.actualStartAt && (
            <div className="text-white text-sm mt-1">
              时长: {Math.floor((Date.now() - new Date(session.actualStartAt).getTime()) / 60000)}分钟
            </div>
          )}
        </div>
      </div>

      <div className="bg-gray-800 py-4 px-6">
        {!sessionEnded ? (
          <>
            <div className="flex items-center justify-center gap-4">
              <button
                onClick={() => setIsMuted(!isMuted)}
                className={`w-14 h-14 rounded-full flex items-center justify-center transition-colors ${
                  isMuted ? 'bg-red-600' : 'bg-gray-700 hover:bg-gray-600'
                }`}
              >
                {isMuted ? (
                  <MicOff className="h-6 w-6 text-white" />
                ) : (
                  <Mic className="h-6 w-6 text-white" />
                )}
              </button>

              <button
                onClick={() => setIsVideoOn(!isVideoOn)}
                className={`w-14 h-14 rounded-full flex items-center justify-center transition-colors ${
                  !isVideoOn ? 'bg-red-600' : 'bg-gray-700 hover:bg-gray-600'
                }`}
              >
                {!isVideoOn ? (
                  <VideoOff className="h-6 w-6 text-white" />
                ) : (
                  <Video className="h-6 w-6 text-white" />
                )}
              </button>

              <button
                onClick={() => setIsScreenSharing(!isScreenSharing)}
                className={`w-14 h-14 rounded-full flex items-center justify-center transition-colors ${
                  isScreenSharing ? 'bg-primary-600' : 'bg-gray-700 hover:bg-gray-600'
                }`}
              >
                <Monitor className="h-6 w-6 text-white" />
              </button>

              <button
                onClick={() => setShowSettings(!showSettings)}
                className="w-14 h-14 rounded-full bg-gray-700 hover:bg-gray-600 flex items-center justify-center"
              >
                <Settings className="h-6 w-6 text-white" />
              </button>

              <button
                onClick={handleEndSession}
                className="w-14 h-14 rounded-full bg-red-600 hover:bg-red-700 flex items-center justify-center"
              >
                <Phone className="h-6 w-6 text-white rotate-[135deg]" />
              </button>
            </div>

            <div className="flex justify-center mt-4">
              <button
                onClick={handleQualityCheck}
                className="text-gray-400 hover:text-white text-sm flex items-center gap-1"
              >
                <MoreVertical className="h-4 w-4" />
                检查连接质量
              </button>
            </div>
          </>
        ) : (
          <div className="text-center py-4">
            <h3 className="text-white text-xl font-semibold mb-4">课程已结束</h3>
            <div className="flex items-center justify-center gap-4">
              <button
                onClick={handleCompleteBooking}
                className="btn-primary flex items-center gap-2"
              >
                <CheckCircle className="h-5 w-5" />
                确认完成并结算
              </button>
              {user?.role === 'teacher' && (
                <button
                  onClick={() => setShowHomeworkModal(true)}
                  className="btn-secondary flex items-center gap-2 text-white bg-blue-600 hover:bg-blue-700 border-blue-600"
                >
                  <BookOpen className="h-5 w-5" />
                  布置作业
                </button>
              )}
            </div>
          </div>
        )}
      </div>

      {showSettings && (
        <div className="absolute bottom-32 left-1/2 -translate-x-1/2 bg-gray-800 rounded-lg p-4 w-64 shadow-xl">
          <h3 className="text-white font-medium mb-3">设置</h3>
          <div className="space-y-3">
            <div>
              <label className="block text-gray-400 text-sm mb-1">摄像头</label>
              <select className="w-full bg-gray-700 text-white rounded px-3 py-2 text-sm">
                <option>默认摄像头</option>
              </select>
            </div>
            <div>
              <label className="block text-gray-400 text-sm mb-1">麦克风</label>
              <select className="w-full bg-gray-700 text-white rounded px-3 py-2 text-sm">
                <option>默认麦克风</option>
              </select>
            </div>
            <div>
              <label className="block text-gray-400 text-sm mb-1">扬声器</label>
              <select className="w-full bg-gray-700 text-white rounded px-3 py-2 text-sm">
                <option>默认扬声器</option>
              </select>
            </div>
          </div>
        </div>
      )}

      {quality && (
        <div className="absolute top-20 right-4 bg-gray-800 rounded-lg p-4 w-64">
          <h3 className="text-white font-medium mb-2">连接质量</h3>
          <div className="space-y-2">
            <div className="flex justify-between text-sm">
              <span className="text-gray-400">状态</span>
              <span className="text-green-400">{quality.status}</span>
            </div>
            <div className="flex justify-between text-sm">
              <span className="text-gray-400">时长</span>
              <span className="text-white">{quality.actualDuration}分钟</span>
            </div>
            <div className="flex justify-between text-sm">
              <span className="text-gray-400">质量评分</span>
              <span className="text-yellow-400">{quality.qualityScore}/100</span>
            </div>
            <div className="flex justify-between text-sm">
              <span className="text-gray-400">重连次数</span>
              <span className="text-white">{quality.reconnects}</span>
            </div>
          </div>
        </div>
      )}

      {showHomeworkModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 w-full max-w-md">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">布置作业</h3>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">作业标题</label>
                <input
                  type="text"
                  value={homeworkForm.title}
                  onChange={(e) => setHomeworkForm({ ...homeworkForm, title: e.target.value })}
                  className="input-field"
                  placeholder="请输入作业标题"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">作业描述</label>
                <textarea
                  value={homeworkForm.description}
                  onChange={(e) => setHomeworkForm({ ...homeworkForm, description: e.target.value })}
                  className="input-field min-h-[100px]"
                  placeholder="请输入作业描述"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">截止日期</label>
                <input
                  type="datetime-local"
                  value={homeworkForm.dueDate}
                  onChange={(e) => setHomeworkForm({ ...homeworkForm, dueDate: e.target.value })}
                  className="input-field"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">满分</label>
                <input
                  type="number"
                  value={homeworkForm.maxScore}
                  onChange={(e) => setHomeworkForm({ ...homeworkForm, maxScore: Number(e.target.value) })}
                  className="input-field"
                  min="1"
                />
              </div>
            </div>
            <div className="flex gap-3 mt-6">
              <button
                onClick={() => setShowHomeworkModal(false)}
                className="btn-secondary flex-1"
              >
                取消
              </button>
              <button
                onClick={handleCreateHomework}
                className="btn-primary flex-1"
              >
                布置作业
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
