import { useState, useEffect } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Loader2, Building2, UserCog } from 'lucide-react'
import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { register as registerUser, selectAuth, clearError } from '@/store/slices/authSlice'

const registerSchema = z.object({
  email: z.string().email('Invalid email address'),
  password: z.string().min(6, 'Password must be at least 6 characters'),
  name: z.string().min(2, 'Name must be at least 2 characters'),
  phone: z.string().optional(),
  role: z.enum(['company', 'jobseeker'], { required_error: 'Please select a role' }),
  company_name: z.string().optional(),
  industry: z.string().optional(),
})

type RegisterForm = z.infer<typeof registerSchema>

export default function Register() {
  const navigate = useNavigate()
  const dispatch = useAppDispatch()
  const { loading, error, user } = useAppSelector(selectAuth)
  const [serverError, setServerError] = useState('')
  const [selectedRole, setSelectedRole] = useState<'company' | 'jobseeker'>('jobseeker')

  const {
    register,
    handleSubmit,
    setValue,
    watch,
    formState: { errors },
  } = useForm<RegisterForm>({
    resolver: zodResolver(registerSchema),
    defaultValues: {
      role: 'jobseeker',
    },
  })

  useEffect(() => {
    if (error) {
      setServerError(error)
      dispatch(clearError())
    }
  }, [error, dispatch])

  useEffect(() => {
    if (user) {
      if (user.role === 'company') navigate('/company/dashboard')
      else if (user.role === 'jobseeker') navigate('/jobseeker/dashboard')
    }
  }, [user, navigate])

  const onSubmit = async (data: RegisterForm) => {
    setServerError('')
    const result = await dispatch(registerUser(data))
    if (registerUser.fulfilled.match(result)) {
      const newUser = result.payload.user
      if (newUser.role === 'company') navigate('/company/dashboard')
      else if (newUser.role === 'jobseeker') navigate('/jobseeker/dashboard')
    }
  }

  const handleRoleSelect = (role: 'company' | 'jobseeker') => {
    setSelectedRole(role)
    setValue('role', role)
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-primary-50 to-primary-100 p-4">
      <div className="w-full max-w-lg">
        <div className="bg-white rounded-2xl shadow-xl p-8">
          <div className="text-center mb-8">
            <h1 className="text-3xl font-bold text-gray-900">Create Account</h1>
            <p className="text-gray-500 mt-2">Join our recruitment platform</p>
          </div>

          {serverError && (
            <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-4">
              {serverError}
            </div>
          )}

          <div className="grid grid-cols-2 gap-4 mb-6">
            <button
              type="button"
              onClick={() => handleRoleSelect('jobseeker')}
              className={`p-4 rounded-lg border-2 transition-all ${
                selectedRole === 'jobseeker'
                  ? 'border-primary-500 bg-primary-50'
                  : 'border-gray-200 hover:border-gray-300'
              }`}
            >
              <UserCog className={`w-8 h-8 mx-auto mb-2 ${selectedRole === 'jobseeker' ? 'text-primary-600' : 'text-gray-400'}`} />
              <div className={`font-medium ${selectedRole === 'jobseeker' ? 'text-primary-700' : 'text-gray-600'}`}>
                Job Seeker
              </div>
            </button>
            <button
              type="button"
              onClick={() => handleRoleSelect('company')}
              className={`p-4 rounded-lg border-2 transition-all ${
                selectedRole === 'company'
                  ? 'border-primary-500 bg-primary-50'
                  : 'border-gray-200 hover:border-gray-300'
              }`}
            >
              <Building2 className={`w-8 h-8 mx-auto mb-2 ${selectedRole === 'company' ? 'text-primary-600' : 'text-gray-400'}`} />
              <div className={`font-medium ${selectedRole === 'company' ? 'text-primary-700' : 'text-gray-600'}`}>
                Company
              </div>
            </button>
          </div>

          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Name</label>
              <input
                {...register('name')}
                type="text"
                className="input-field"
                placeholder="Your name"
              />
              {errors.name && (
                <p className="text-red-500 text-sm mt-1">{errors.name.message}</p>
              )}
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Email</label>
              <input
                {...register('email')}
                type="email"
                className="input-field"
                placeholder="you@example.com"
              />
              {errors.email && (
                <p className="text-red-500 text-sm mt-1">{errors.email.message}</p>
              )}
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Phone</label>
              <input
                {...register('phone')}
                type="tel"
                className="input-field"
                placeholder="Phone number"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Password</label>
              <input
                {...register('password')}
                type="password"
                className="input-field"
                placeholder="••••••••"
              />
              {errors.password && (
                <p className="text-red-500 text-sm mt-1">{errors.password.message}</p>
              )}
            </div>

            {selectedRole === 'company' && (
              <>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Company Name</label>
                  <input
                    {...register('company_name')}
                    type="text"
                    className="input-field"
                    placeholder="Company name"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Industry</label>
                  <input
                    {...register('industry')}
                    type="text"
                    className="input-field"
                    placeholder="e.g. Technology, Healthcare, Finance"
                  />
                </div>
              </>
            )}

            <button
              type="submit"
              disabled={loading}
              className="w-full btn-primary disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
            >
              {loading && <Loader2 className="w-5 h-5 animate-spin" />}
              {loading ? 'Creating account...' : 'Create Account'}
            </button>
          </form>

          <div className="mt-6 text-center">
            <p className="text-gray-500">
              Already have an account?{' '}
              <Link to="/login" className="text-primary-600 hover:text-primary-700 font-medium">
                Sign in
              </Link>
            </p>
          </div>
        </div>
      </div>
    </div>
  )
}
