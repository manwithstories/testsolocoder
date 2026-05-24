import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { toast } from 'sonner'
import { useForm } from 'react-hook-form'
import { authApi } from '@/services/api'
import { GraduationCap, Mail, Lock, User, Eye, EyeOff } from 'lucide-react'

interface RegisterForm {
  email: string
  password: string
  firstName: string
  lastName: string
  role: 'teacher' | 'student'
}

export default function Register() {
  const navigate = useNavigate()
  const [showPassword, setShowPassword] = useState(false)
  const [loading, setLoading] = useState(false)

  const { register, handleSubmit, formState: { errors } } = useForm<RegisterForm>({
    defaultValues: {
      role: 'student',
    },
  })

  const onSubmit = async (data: RegisterForm) => {
    try {
      setLoading(true)
      await authApi.register(data)
      toast.success('注册成功！请登录')
      navigate('/login')
    } catch (error: any) {
      toast.error(error.response?.data?.error || '注册失败')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-primary-50 to-white flex items-center justify-center px-4 py-8">
      <div className="max-w-md w-full">
        <div className="text-center mb-8">
          <div className="inline-flex items-center justify-center w-16 h-16 bg-primary-600 rounded-full mb-4">
            <GraduationCap className="h-8 w-8 text-white" />
          </div>
          <h1 className="text-3xl font-bold text-gray-900">创建账户</h1>
          <p className="text-gray-500 mt-2">加入我们的在线家教平台</p>
        </div>

        <div className="bg-white rounded-2xl shadow-xl p-8">
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">名</label>
                <div className="relative">
                  <User className="absolute left-3 top-1/2 -translate-y-1/2 h-5 w-5 text-gray-400" />
                  <input
                    {...register('firstName', { required: '请输入名' })}
                    className="input-field pl-10"
                    placeholder="张"
                  />
                </div>
                {errors.firstName && <p className="text-red-500 text-sm mt-1">{errors.firstName.message}</p>}
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">姓</label>
                <div className="relative">
                  <User className="absolute left-3 top-1/2 -translate-y-1/2 h-5 w-5 text-gray-400" />
                  <input
                    {...register('lastName', { required: '请输入姓' })}
                    className="input-field pl-10"
                    placeholder="三"
                  />
                </div>
                {errors.lastName && <p className="text-red-500 text-sm mt-1">{errors.lastName.message}</p>}
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">邮箱</label>
              <div className="relative">
                <Mail className="absolute left-3 top-1/2 -translate-y-1/2 h-5 w-5 text-gray-400" />
                <input
                  type="email"
                  {...register('email', { required: '请输入邮箱' })}
                  className="input-field pl-10"
                  placeholder="your@email.com"
                />
              </div>
              {errors.email && <p className="text-red-500 text-sm mt-1">{errors.email.message}</p>}
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">密码</label>
              <div className="relative">
                <Lock className="absolute left-3 top-1/2 -translate-y-1/2 h-5 w-5 text-gray-400" />
                <input
                  type={showPassword ? 'text' : 'password'}
                  {...register('password', { required: '请输入密码', minLength: { value: 6, message: '密码至少6位' } })}
                  className="input-field pl-10 pr-10"
                  placeholder="••••••••"
                />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
                >
                  {showPassword ? <EyeOff className="h-5 w-5" /> : <Eye className="h-5 w-5" />}
                </button>
              </div>
              {errors.password && <p className="text-red-500 text-sm mt-1">{errors.password.message}</p>}
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">身份</label>
              <div className="grid grid-cols-2 gap-4">
                <label className={`flex items-center justify-center p-4 border-2 rounded-lg cursor-pointer transition-colors ${
                  // @ts-ignore
                  errors.role ? 'border-red-300' : 'border-gray-200'
                }`}>
                  <input
                    type="radio"
                    value="student"
                    {...register('role', { required: '请选择身份' })}
                    className="hidden"
                  />
                  <div className="text-center">
                    <div className="text-2xl mb-1">🎓</div>
                    <div className="font-medium">学生</div>
                  </div>
                </label>

                <label className={`flex items-center justify-center p-4 border-2 rounded-lg cursor-pointer transition-colors ${
                  // @ts-ignore
                  errors.role ? 'border-red-300' : 'border-gray-200'
                }`}>
                  <input
                    type="radio"
                    value="teacher"
                    {...register('role', { required: '请选择身份' })}
                    className="hidden"
                  />
                  <div className="text-center">
                    <div className="text-2xl mb-1">👨‍🏫</div>
                    <div className="font-medium">老师</div>
                  </div>
                </label>
              </div>
              {errors.role && <p className="text-red-500 text-sm mt-1">{errors.role.message}</p>}
            </div>

            <button
              type="submit"
              disabled={loading}
              className="btn-primary w-full py-3 disabled:opacity-50"
            >
              {loading ? '注册中...' : '创建账户'}
            </button>
          </form>

          <div className="mt-6 text-center">
            <p className="text-gray-500">
              已有账户？{' '}
              <button
                onClick={() => navigate('/login')}
                className="text-primary-600 font-medium hover:text-primary-700"
              >
                立即登录
              </button>
            </p>
          </div>
        </div>
      </div>
    </div>
  )
}
