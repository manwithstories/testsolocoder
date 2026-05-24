import { useState, useEffect } from 'react'
import { toast } from 'sonner'
import { walletApi } from '@/services/api'
import { Wallet, Transaction } from '@/types'
import { DollarSign, CreditCard, ArrowUpRight, ArrowDownLeft, Clock, CheckCircle, XCircle } from 'lucide-react'

export default function WalletPage() {
  const [wallet, setWallet] = useState<Wallet | null>(null)
  const [transactions, setTransactions] = useState<Transaction[]>([])
  const [loading, setLoading] = useState(true)
  const [showDepositModal, setShowDepositModal] = useState(false)
  const [showWithdrawModal, setShowWithdrawModal] = useState(false)
  const [depositAmount, setDepositAmount] = useState('')
  const [withdrawAmount, setWithdrawAmount] = useState('')
  const [withdrawBankInfo, setWithdrawBankInfo] = useState({
    bankName: '',
    bankAccount: '',
    accountHolder: '',
  })

  useEffect(() => {
    loadWalletData()
  }, [])

  const loadWalletData = async () => {
    try {
      setLoading(true)
      const [walletRes, transactionsRes] = await Promise.all([
        walletApi.getWallet(),
        walletApi.getTransactions(),
      ])
      setWallet(walletRes.data)
      setTransactions(transactionsRes.data)
    } catch (error) {
      console.error('Failed to load wallet data:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleDeposit = async () => {
    if (!depositAmount || Number(depositAmount) <= 0) {
      toast.error('请输入有效金额')
      return
    }

    try {
      await walletApi.deposit({
        amount: Number(depositAmount),
        currency: 'USD',
        paymentMethod: 'credit_card',
      })
      toast.success('充值成功！')
      setShowDepositModal(false)
      setDepositAmount('')
      loadWalletData()
    } catch (error: any) {
      toast.error(error.response?.data?.error || '充值失败')
    }
  }

  const handleWithdraw = async () => {
    if (!withdrawAmount || Number(withdrawAmount) <= 0) {
      toast.error('请输入有效金额')
      return
    }
    if (Number(withdrawAmount) > (wallet?.balance || 0)) {
      toast.error('余额不足')
      return
    }
    if (!withdrawBankInfo.bankName || !withdrawBankInfo.bankAccount || !withdrawBankInfo.accountHolder) {
      toast.error('请填写完整的银行信息')
      return
    }

    try {
      await walletApi.withdraw({
        amount: Number(withdrawAmount),
        currency: 'USD',
        ...withdrawBankInfo,
      })
      toast.success('提现申请已提交！')
      setShowWithdrawModal(false)
      setWithdrawAmount('')
      setWithdrawBankInfo({ bankName: '', bankAccount: '', accountHolder: '' })
      loadWalletData()
    } catch (error: any) {
      toast.error(error.response?.data?.error || '提现失败')
    }
  }

  const getTypeIcon = (type: string) => {
    switch (type) {
      case 'deposit': return <ArrowDownLeft className="h-4 w-4 text-green-500" />
      case 'withdrawal': return <ArrowUpRight className="h-4 w-4 text-red-500" />
      case 'payment': return <ArrowUpRight className="h-4 w-4 text-orange-500" />
      default: return <DollarSign className="h-4 w-4 text-gray-500" />
    }
  }

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'completed': return <CheckCircle className="h-4 w-4 text-green-500" />
      case 'pending': return <Clock className="h-4 w-4 text-yellow-500" />
      case 'failed': return <XCircle className="h-4 w-4 text-red-500" />
      default: return null
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">我的钱包</h1>
        <p className="text-gray-500">管理您的账户余额和交易记录</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="card bg-gradient-to-br from-primary-600 to-primary-800 text-white">
          <div className="flex items-center justify-between mb-4">
            <DollarSign className="h-8 w-8 opacity-80" />
            <span className="text-sm opacity-80">账户余额</span>
          </div>
          <div className="text-4xl font-bold">
            ${wallet?.balance?.toFixed(2) || '0.00'}
          </div>
          <div className="text-sm opacity-80 mt-2">
            总收入: ${wallet?.totalIncome?.toFixed(2) || '0.00'}
          </div>
        </div>

        <div className="card">
          <div className="flex items-center gap-3 mb-4">
            <div className="w-12 h-12 bg-green-100 rounded-lg flex items-center justify-center">
              <ArrowDownLeft className="h-6 w-6 text-green-600" />
            </div>
            <div>
              <div className="text-sm text-gray-500">总收入</div>
              <div className="text-2xl font-bold text-gray-900">
                ${wallet?.totalIncome?.toFixed(2) || '0.00'}
              </div>
            </div>
          </div>
          <button
            onClick={() => setShowDepositModal(true)}
            className="btn-primary w-full"
          >
            充值
          </button>
        </div>

        <div className="card">
          <div className="flex items-center gap-3 mb-4">
            <div className="w-12 h-12 bg-red-100 rounded-lg flex items-center justify-center">
              <ArrowUpRight className="h-6 w-6 text-red-600" />
            </div>
            <div>
              <div className="text-sm text-gray-500">总支出</div>
              <div className="text-2xl font-bold text-gray-900">
                ${wallet?.totalSpent?.toFixed(2) || '0.00'}
              </div>
            </div>
          </div>
          <button
            onClick={() => setShowWithdrawModal(true)}
            className="btn-secondary w-full"
          >
            提现
          </button>
        </div>
      </div>

      <div className="card">
        <h2 className="text-lg font-semibold text-gray-900 mb-4">交易记录</h2>
        {transactions.length > 0 ? (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-gray-200">
                  <th className="text-left py-3 px-4 text-sm font-medium text-gray-500">类型</th>
                  <th className="text-left py-3 px-4 text-sm font-medium text-gray-500">金额</th>
                  <th className="text-left py-3 px-4 text-sm font-medium text-gray-500">状态</th>
                  <th className="text-left py-3 px-4 text-sm font-medium text-gray-500">时间</th>
                </tr>
              </thead>
              <tbody>
                {transactions.map((tx) => (
                  <tr key={tx.id} className="border-b border-gray-100 hover:bg-gray-50">
                    <td className="py-3 px-4">
                      <div className="flex items-center gap-2">
                        {getTypeIcon(tx.type)}
                        <span className="text-sm">
                          {tx.type === 'deposit' ? '充值' :
                           tx.type === 'withdrawal' ? '提现' :
                           tx.type === 'payment' ? '支付' :
                           tx.type === 'commission' ? '佣金' : tx.type}
                        </span>
                      </div>
                    </td>
                    <td className="py-3 px-4">
                      <span className={`font-medium ${
                        tx.type === 'deposit' ? 'text-green-600' : 'text-red-600'
                      }`}>
                        {tx.type === 'deposit' ? '+' : '-'}${tx.amount.toFixed(2)}
                      </span>
                    </td>
                    <td className="py-3 px-4">
                      <div className="flex items-center gap-1">
                        {getStatusIcon(tx.status)}
                        <span className="text-sm">
                          {tx.status === 'completed' ? '已完成' :
                           tx.status === 'pending' ? '处理中' :
                           tx.status === 'failed' ? '失败' : tx.status}
                        </span>
                      </div>
                    </td>
                    <td className="py-3 px-4 text-sm text-gray-500">
                      {new Date(tx.createdAt).toLocaleDateString('zh-CN', {
                        year: 'numeric',
                        month: 'short',
                        day: 'numeric',
                        hour: '2-digit',
                        minute: '2-digit',
                      })}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        ) : (
          <div className="text-center py-12">
            <CreditCard className="h-12 w-12 text-gray-300 mx-auto mb-4" />
            <p className="text-gray-500">暂无交易记录</p>
          </div>
        )}
      </div>

      {showDepositModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-2xl p-6 w-full max-w-md mx-4">
            <h3 className="text-xl font-bold text-gray-900 mb-4">充值</h3>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">充值金额</label>
                <div className="relative">
                  <DollarSign className="absolute left-3 top-1/2 -translate-y-1/2 h-5 w-5 text-gray-400" />
                  <input
                    type="number"
                    value={depositAmount}
                    onChange={(e) => setDepositAmount(e.target.value)}
                    className="input-field pl-10"
                    placeholder="0.00"
                    min="1"
                    step="0.01"
                  />
                </div>
              </div>
            </div>
            <div className="flex gap-3 mt-6">
              <button
                onClick={() => setShowDepositModal(false)}
                className="btn-secondary flex-1"
              >
                取消
              </button>
              <button
                onClick={handleDeposit}
                className="btn-primary flex-1"
              >
                确认充值
              </button>
            </div>
          </div>
        </div>
      )}

      {showWithdrawModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-2xl p-6 w-full max-w-md mx-4 max-h-[90vh] overflow-y-auto">
            <h3 className="text-xl font-bold text-gray-900 mb-4">提现</h3>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">提现金额</label>
                <div className="relative">
                  <DollarSign className="absolute left-3 top-1/2 -translate-y-1/2 h-5 w-5 text-gray-400" />
                  <input
                    type="number"
                    value={withdrawAmount}
                    onChange={(e) => setWithdrawAmount(e.target.value)}
                    className="input-field pl-10"
                    placeholder="0.00"
                    min="100"
                    step="0.01"
                  />
                </div>
                <p className="text-sm text-gray-500 mt-1">
                  当前余额: ${wallet?.balance?.toFixed(2) || '0.00'}
                </p>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">银行名称</label>
                <input
                  type="text"
                  value={withdrawBankInfo.bankName}
                  onChange={(e) => setWithdrawBankInfo({ ...withdrawBankInfo, bankName: e.target.value })}
                  className="input-field"
                  placeholder="请输入银行名称"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">银行账号</label>
                <input
                  type="text"
                  value={withdrawBankInfo.bankAccount}
                  onChange={(e) => setWithdrawBankInfo({ ...withdrawBankInfo, bankAccount: e.target.value })}
                  className="input-field"
                  placeholder="请输入银行账号"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">账户持有人</label>
                <input
                  type="text"
                  value={withdrawBankInfo.accountHolder}
                  onChange={(e) => setWithdrawBankInfo({ ...withdrawBankInfo, accountHolder: e.target.value })}
                  className="input-field"
                  placeholder="请输入账户持有人姓名"
                />
              </div>
            </div>
            <div className="flex gap-3 mt-6">
              <button
                onClick={() => setShowWithdrawModal(false)}
                className="btn-secondary flex-1"
              >
                取消
              </button>
              <button
                onClick={handleWithdraw}
                className="btn-primary flex-1"
              >
                申请提现
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
