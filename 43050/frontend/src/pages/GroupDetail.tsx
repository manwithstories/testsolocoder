import React, { useState, useEffect } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import {
  groupAPI,
  expenseAPI,
  settlementAPI,
  CreateExpenseData,
  ParticipantInput,
} from '../api/services';
import { Group, Expense, Transfer, MemberStats, GroupMember } from '../types';
import { useAuth } from '../context/AuthContext';
import dayjs from 'dayjs';

type TabType = 'expenses' | 'settlement' | 'members' | 'stats';

const GroupDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { user } = useAuth();
  const [group, setGroup] = useState<Group | null>(null);
  const [expenses, setExpenses] = useState<Expense[]>([]);
  const [members, setMembers] = useState<GroupMember[]>([]);
  const [transfers, setTransfers] = useState<Transfer[]>([]);
  const [memberStats, setMemberStats] = useState<MemberStats[]>([]);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState<TabType>('expenses');
  const [showAddExpense, setShowAddExpense] = useState(false);
  const [showInviteModal, setShowInviteModal] = useState(false);
  const [copied, setCopied] = useState(false);
  const [error, setError] = useState('');

  const [expenseForm, setExpenseForm] = useState({
    title: '',
    amount: '',
    paidBy: user?.id?.toString() || '',
    splitType: 'equal' as 'equal' | 'ratio' | 'custom',
    participants: [] as ParticipantInput[],
  });

  const groupId = parseInt(id || '0', 10);

  const fetchData = async () => {
    try {
      const [groupRes, expensesRes, membersRes, transfersRes, statsRes] = await Promise.all([
        groupAPI.getGroup(groupId),
        expenseAPI.getGroupExpenses(groupId),
        groupAPI.getGroupMembers(groupId),
        settlementAPI.getOptimalTransfers(groupId),
        settlementAPI.getGroupStats(groupId),
      ]);
      setGroup(groupRes.data);
      setExpenses(expensesRes.data);
      setMembers(membersRes.data);
      setTransfers(transfersRes.data);
      setMemberStats(statsRes.data);

      const activeMembers = membersRes.data.filter((m) => m.isActive);
      setExpenseForm((prev) => ({
        ...prev,
        paidBy: user?.id?.toString() || activeMembers[0]?.user?.id?.toString() || '',
        participants: activeMembers.map((m) => ({
          userId: m.user.id,
          amount: 0,
          ratio: 1,
        })),
      }));
    } catch (err) {
      console.error('Failed to fetch group data:', err);
      if (err.response?.status === 403) {
        navigate('/groups');
      }
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (groupId) {
      fetchData();
    }
  }, [groupId]);

  const handleAddExpense = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    const amount = parseFloat(expenseForm.amount);
    if (isNaN(amount) || amount <= 0) {
      setError('请输入有效的金额');
      return;
    }

    const paidBy = parseInt(expenseForm.paidBy, 10);
    if (!paidBy) {
      setError('请选择付款人');
      return;
    }

    const validParticipants = expenseForm.participants.filter((p) => {
      if (expenseForm.splitType === 'custom') {
        return p.amount && p.amount > 0;
      }
      return true;
    });

    if (validParticipants.length === 0) {
      setError('请至少选择一位参与成员');
      return;
    }

    const data: CreateExpenseData = {
      title: expenseForm.title,
      amount,
      paidBy,
      splitType: expenseForm.splitType,
      participants: validParticipants,
      expenseDate: new Date().toISOString(),
    };

    try {
      await expenseAPI.createExpense(groupId, data);
      setShowAddExpense(false);
      resetExpenseForm();
      fetchData();
    } catch (err: any) {
      setError(err.response?.data?.error || '添加失败，请稍后重试');
    }
  };

  const resetExpenseForm = () => {
    const activeMembers = members.filter((m) => m.isActive);
    setExpenseForm({
      title: '',
      amount: '',
      paidBy: user?.id?.toString() || activeMembers[0]?.user?.id?.toString() || '',
      splitType: 'equal',
      participants: activeMembers.map((m) => ({
        userId: m.user.id,
        amount: 0,
        ratio: 1,
      })),
    });
    setError('');
  };

  const toggleParticipant = (userId: number) => {
    setExpenseForm((prev) => {
      const existing = prev.participants.find((p) => p.userId === userId);
      if (existing) {
        return {
          ...prev,
          participants: prev.participants.filter((p) => p.userId !== userId),
        };
      } else {
        return {
          ...prev,
          participants: [
            ...prev.participants,
            { userId, amount: 0, ratio: 1 },
          ],
        };
      }
    });
  };

  const updateParticipantAmount = (userId: number, amount: number) => {
    setExpenseForm((prev) => ({
      ...prev,
      participants: prev.participants.map((p) =>
        p.userId === userId ? { ...p, amount } : p
      ),
    }));
  };

  const updateParticipantRatio = (userId: number, ratio: number) => {
    setExpenseForm((prev) => ({
      ...prev,
      participants: prev.participants.map((p) =>
        p.userId === userId ? { ...p, ratio } : p
      ),
    }));
  };

  const handleLeaveGroup = async () => {
    if (!confirm('确定要退出该群组吗？历史账单记录将保留。')) {
      return;
    }
    try {
      await groupAPI.leaveGroup(groupId);
      navigate('/groups');
    } catch (err: any) {
      alert(err.response?.data?.error || '退出失败');
    }
  };

  const copyInviteCode = () => {
    if (group?.inviteCode) {
      navigator.clipboard.writeText(group.inviteCode);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    }
  };

  const handleDeleteExpense = async (expenseId: number) => {
    if (!confirm('确定要删除这笔账单吗？')) {
      return;
    }
    try {
      await expenseAPI.deleteExpense(expenseId);
      fetchData();
    } catch (err: any) {
      alert(err.response?.data?.error || '删除失败');
    }
  };

  if (loading) {
    return <div className="text-center py-12">加载中...</div>;
  }

  if (!group) {
    return <div className="text-center py-12">群组不存在</div>;
  }

  const activeMembers = members.filter((m) => m.isActive);
  const totalExpense = expenses.reduce((sum, e) => sum + e.amount, 0);

  const getSplitPreview = () => {
    const amount = parseFloat(expenseForm.amount) || 0;
    const validCount = expenseForm.participants.length;

    if (expenseForm.splitType === 'equal' && validCount > 0) {
      const perPerson = amount / validCount;
      return expenseForm.participants.map((p) => {
        const member = activeMembers.find((m) => m.user.id === p.userId);
        return {
          name: member?.user.username || '',
          amount: perPerson,
        };
      });
    } else if (expenseForm.splitType === 'ratio') {
      const totalRatio = expenseForm.participants.reduce((sum, p) => sum + (p.ratio || 0), 0);
      if (totalRatio > 0) {
        return expenseForm.participants.map((p) => {
          const member = activeMembers.find((m) => m.user.id === p.userId);
          return {
            name: member?.user.username || '',
            amount: (amount * (p.ratio || 0)) / totalRatio,
          };
        });
      }
    } else if (expenseForm.splitType === 'custom') {
      return expenseForm.participants.map((p) => {
        const member = activeMembers.find((m) => m.user.id === p.userId);
        return {
          name: member?.user.username || '',
          amount: p.amount || 0,
        };
      });
    }
    return [];
  };

  const splitPreview = getSplitPreview();
  const totalPreview = splitPreview.reduce((sum, p) => sum + p.amount, 0);

  return (
    <div className="space-y-6">
      <div className="bg-white rounded-xl shadow-sm border border-gray-100 p-6">
        <div className="flex items-start justify-between">
          <div className="flex items-center gap-4">
            <div className="w-16 h-16 bg-primary-100 rounded-full flex items-center justify-center text-primary-600 font-bold text-2xl">
              {group.name.charAt(0).toUpperCase()}
            </div>
            <div>
              <h1 className="text-2xl font-bold text-gray-800">{group.name}</h1>
              {group.description && (
                <p className="text-gray-500 mt-1">{group.description}</p>
              )}
              <div className="flex items-center gap-4 mt-2 text-sm text-gray-500">
                <span>{activeMembers.length} 位成员</span>
                <span>{expenses.length} 笔账单</span>
                <span>总支出 ¥{totalExpense.toFixed(2)}</span>
              </div>
            </div>
          </div>
          <div className="flex gap-2">
            <button
              onClick={() => setShowInviteModal(true)}
              className="px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors"
            >
              邀请成员
            </button>
            {group.creatorId !== user?.id && (
              <button
                onClick={handleLeaveGroup}
                className="px-4 py-2 border border-red-300 text-red-600 rounded-lg hover:bg-red-50 transition-colors"
              >
                退出群组
              </button>
            )}
          </div>
        </div>
      </div>

      <div className="bg-white rounded-xl shadow-sm border border-gray-100">
        <div className="border-b border-gray-100">
          <nav className="flex">
            {[
              { key: 'expenses', label: '账单记录', icon: '📝' },
              { key: 'settlement', label: '智能结算', icon: '💱' },
              { key: 'members', label: '成员列表', icon: '👥' },
              { key: 'stats', label: '统计分析', icon: '📊' },
            ].map((tab) => (
              <button
                key={tab.key}
                onClick={() => setActiveTab(tab.key as TabType)}
                className={`flex items-center gap-2 px-6 py-4 text-sm font-medium transition-colors border-b-2 ${
                  activeTab === tab.key
                    ? 'border-primary-500 text-primary-600'
                    : 'border-transparent text-gray-500 hover:text-gray-700'
                }`}
              >
                <span>{tab.icon}</span>
                <span>{tab.label}</span>
              </button>
            ))}
          </nav>
        </div>

        <div className="p-6">
          {activeTab === 'expenses' && (
            <div className="space-y-4">
              <div className="flex justify-between items-center">
                <h2 className="text-lg font-semibold text-gray-800">账单记录</h2>
                <button
                  onClick={() => {
                    resetExpenseForm();
                    setShowAddExpense(true);
                  }}
                  className="px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition-colors"
                >
                  + 添加账单
                </button>
              </div>

              {expenses.length === 0 ? (
                <div className="text-center py-12">
                  <div className="text-5xl mb-4">📝</div>
                  <p className="text-gray-500 mb-4">暂无账单记录</p>
                  <button
                    onClick={() => setShowAddExpense(true)}
                    className="px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition-colors"
                  >
                    添加第一笔账单
                  </button>
                </div>
              ) : (
                <div className="space-y-3">
                  {expenses.map((expense) => (
                    <div
                      key={expense.id}
                      className="border border-gray-200 rounded-lg p-4 hover:border-primary-300 transition-colors"
                    >
                      <div className="flex items-start justify-between">
                        <div className="flex-1">
                          <div className="flex items-center gap-3">
                            <h3 className="font-medium text-gray-800">{expense.title}</h3>
                            <span className={`px-2 py-1 text-xs rounded-full ${
                              expense.splitType === 'equal' ? 'bg-blue-100 text-blue-700' :
                              expense.splitType === 'ratio' ? 'bg-purple-100 text-purple-700' :
                              'bg-orange-100 text-orange-700'
                            }`}>
                              {expense.splitType === 'equal' ? '均摊' :
                               expense.splitType === 'ratio' ? '按比例' : '自定义'}
                            </span>
                          </div>
                          <p className="text-sm text-gray-500 mt-1">
                            {expense.payer?.username} 支付 ¥{expense.amount.toFixed(2)}
                          </p>
                          <div className="mt-2 flex flex-wrap gap-2">
                            {expense.participants.map((p) => (
                              <span
                                key={p.id}
                                className="px-2 py-1 bg-gray-100 text-gray-600 text-xs rounded"
                              >
                                {p.user?.username}: ¥{p.amount.toFixed(2)}
                              </span>
                            ))}
                          </div>
                        </div>
                        <div className="text-right">
                          <p className="text-lg font-bold text-gray-800">
                            ¥{expense.amount.toFixed(2)}
                          </p>
                          <p className="text-xs text-gray-500 mt-1">
                            {dayjs(expense.expenseDate).format('YYYY-MM-DD')}
                          </p>
                          {expense.createdBy === user?.id && (
                            <button
                              onClick={() => handleDeleteExpense(expense.id)}
                              className="text-xs text-red-500 hover:text-red-700 mt-2"
                            >
                              删除
                            </button>
                          )}
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          )}

          {activeTab === 'settlement' && (
            <div className="space-y-6">
              <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
                <p className="text-blue-800 text-sm">
                  💡 智能结算功能会自动计算最优转账路径，将复杂的欠款关系简化为最少的转账笔数。
                </p>
              </div>

              {transfers.length === 0 ? (
                <div className="text-center py-12">
                  <div className="text-5xl mb-4">🎉</div>
                  <p className="text-gray-500">太棒了！所有账单已结清，无需转账。</p>
                </div>
              ) : (
                <div className="space-y-3">
                  <h3 className="text-lg font-semibold text-gray-800">最优转账方案</h3>
                  {transfers.map((transfer, index) => (
                    <div
                      key={index}
                      className="flex items-center justify-between p-4 bg-gray-50 rounded-lg"
                    >
                      <div className="flex items-center gap-3">
                        <div className="w-10 h-10 bg-red-100 rounded-full flex items-center justify-center text-red-600 font-bold">
                          {transfer.fromUsername.charAt(0).toUpperCase()}
                        </div>
                        <span className="font-medium text-gray-800">{transfer.fromUsername}</span>
                        <span className="text-gray-400">→</span>
                        <div className="w-10 h-10 bg-green-100 rounded-full flex items-center justify-center text-green-600 font-bold">
                          {transfer.toUsername.charAt(0).toUpperCase()}
                        </div>
                        <span className="font-medium text-gray-800">{transfer.toUsername}</span>
                      </div>
                      <span className="text-lg font-bold text-primary-600">
                        ¥{transfer.amount.toFixed(2)}
                      </span>
                    </div>
                  ))}
                  <p className="text-sm text-gray-500 text-center pt-4">
                    共需要 {transfers.length} 笔转账即可结清所有账单
                  </p>
                </div>
              )}
            </div>
          )}

          {activeTab === 'members' && (
            <div className="space-y-4">
              <h2 className="text-lg font-semibold text-gray-800">成员列表</h2>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {activeMembers.map((member) => (
                  <div
                    key={member.user.id}
                    className="flex items-center gap-4 p-4 border border-gray-200 rounded-lg"
                  >
                    <div className="w-12 h-12 bg-primary-100 rounded-full flex items-center justify-center text-primary-600 font-bold text-lg">
                      {member.user.username.charAt(0).toUpperCase()}
                    </div>
                    <div className="flex-1">
                      <p className="font-medium text-gray-800">
                        {member.user.username}
                        {member.user.id === group.creatorId && (
                          <span className="ml-2 px-2 py-0.5 bg-yellow-100 text-yellow-700 text-xs rounded-full">
                            群主
                          </span>
                        )}
                        {member.user.id === user?.id && (
                          <span className="ml-2 px-2 py-0.5 bg-blue-100 text-blue-700 text-xs rounded-full">
                            我
                          </span>
                        )}
                      </p>
                      <p className="text-sm text-gray-500">{member.user.email}</p>
                      <p className="text-xs text-gray-400 mt-1">
                        加入于 {dayjs(member.joinedAt).format('YYYY-MM-DD')}
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}

          {activeTab === 'stats' && (
            <div className="space-y-6">
              <h2 className="text-lg font-semibold text-gray-800">成员收支统计</h2>
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead>
                    <tr className="border-b border-gray-200">
                      <th className="text-left py-3 px-4 font-medium text-gray-600">成员</th>
                      <th className="text-right py-3 px-4 font-medium text-gray-600">已支付</th>
                      <th className="text-right py-3 px-4 font-medium text-gray-600">应承担</th>
                      <th className="text-right py-3 px-4 font-medium text-gray-600">余额</th>
                    </tr>
                  </thead>
                  <tbody>
                    {memberStats.map((stat) => (
                      <tr key={stat.userId} className="border-b border-gray-100">
                        <td className="py-3 px-4">
                          <div className="flex items-center gap-2">
                            <div className="w-8 h-8 bg-primary-100 rounded-full flex items-center justify-center text-primary-600 font-bold text-sm">
                              {stat.username.charAt(0).toUpperCase()}
                            </div>
                            <span className="font-medium text-gray-800">{stat.username}</span>
                          </div>
                        </td>
                        <td className="py-3 px-4 text-right text-green-600">
                          ¥{stat.totalPaid.toFixed(2)}
                        </td>
                        <td className="py-3 px-4 text-right text-red-500">
                          ¥{stat.totalOwed.toFixed(2)}
                        </td>
                        <td className={`py-3 px-4 text-right font-bold ${
                          stat.balance >= 0 ? 'text-green-600' : 'text-red-500'
                        }`}>
                          {stat.balance >= 0 ? '+' : ''}¥{stat.balance.toFixed(2)}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          )}
        </div>
      </div>

      {showAddExpense && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-xl w-full max-w-2xl max-h-[90vh] overflow-y-auto">
            <div className="p-6 border-b border-gray-100">
              <h2 className="text-xl font-bold text-gray-800">添加账单</h2>
            </div>
            <form onSubmit={handleAddExpense} className="p-6 space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">账单名称</label>
                <input
                  type="text"
                  value={expenseForm.title}
                  onChange={(e) => setExpenseForm({ ...expenseForm, title: e.target.value })}
                  className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
                  placeholder="例如：午餐、房租等"
                  required
                />
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">总金额</label>
                  <input
                    type="number"
                    step="0.01"
                    min="0"
                    value={expenseForm.amount}
                    onChange={(e) => setExpenseForm({ ...expenseForm, amount: e.target.value })}
                    className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
                    placeholder="0.00"
                    required
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">付款人</label>
                  <select
                    value={expenseForm.paidBy}
                    onChange={(e) => setExpenseForm({ ...expenseForm, paidBy: e.target.value })}
                    className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
                    required
                  >
                    {activeMembers.map((m) => (
                      <option key={m.user.id} value={m.user.id}>
                        {m.user.username}
                      </option>
                    ))}
                  </select>
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">分摊方式</label>
                <div className="flex gap-3">
                  {[
                    { value: 'equal', label: '均摊', desc: '所有参与成员平分' },
                    { value: 'ratio', label: '按比例', desc: '按设定比例分摊' },
                    { value: 'custom', label: '自定义', desc: '手动输入每人金额' },
                  ].map((type) => (
                    <button
                      key={type.value}
                      type="button"
                      onClick={() => setExpenseForm({ ...expenseForm, splitType: type.value as any })}
                      className={`flex-1 p-3 rounded-lg border-2 text-left transition-colors ${
                        expenseForm.splitType === type.value
                          ? 'border-primary-500 bg-primary-50'
                          : 'border-gray-200 hover:border-gray-300'
                      }`}
                    >
                      <p className={`font-medium ${
                        expenseForm.splitType === type.value ? 'text-primary-600' : 'text-gray-800'
                      }`}>
                        {type.label}
                      </p>
                      <p className="text-xs text-gray-500 mt-1">{type.desc}</p>
                    </button>
                  ))}
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">参与成员</label>
                <div className="space-y-2">
                  {activeMembers.map((member) => {
                    const isSelected = expenseForm.participants.some(
                      (p) => p.userId === member.user.id
                    );
                    const participant = expenseForm.participants.find(
                      (p) => p.userId === member.user.id
                    );

                    return (
                      <div
                        key={member.user.id}
                        className={`p-3 rounded-lg border-2 transition-colors ${
                          isSelected
                            ? 'border-primary-300 bg-primary-50'
                            : 'border-gray-200 hover:border-gray-300'
                        }`}
                      >
                        <div className="flex items-center gap-3">
                          <input
                            type="checkbox"
                            checked={isSelected}
                            onChange={() => toggleParticipant(member.user.id)}
                            className="w-4 h-4 text-primary-600 rounded"
                          />
                          <span className="font-medium text-gray-800">
                            {member.user.username}
                          </span>
                          {member.user.id === parseInt(expenseForm.paidBy, 10) && (
                            <span className="text-xs px-2 py-0.5 bg-green-100 text-green-700 rounded">
                              付款人
                            </span>
                          )}

                          {isSelected && expenseForm.splitType === 'ratio' && (
                            <div className="ml-auto flex items-center gap-2">
                              <span className="text-sm text-gray-500">比例:</span>
                              <input
                                type="number"
                                step="0.1"
                                min="0"
                                value={participant?.ratio || 1}
                                onChange={(e) =>
                                  updateParticipantRatio(
                                    member.user.id,
                                    parseFloat(e.target.value) || 0
                                  )
                                }
                                className="w-20 px-2 py-1 border border-gray-300 rounded text-sm"
                              />
                            </div>
                          )}

                          {isSelected && expenseForm.splitType === 'custom' && (
                            <div className="ml-auto flex items-center gap-2">
                              <span className="text-sm text-gray-500">¥</span>
                              <input
                                type="number"
                                step="0.01"
                                min="0"
                                value={participant?.amount || ''}
                                onChange={(e) =>
                                  updateParticipantAmount(
                                    member.user.id,
                                    parseFloat(e.target.value) || 0
                                  )
                                }
                                className="w-24 px-2 py-1 border border-gray-300 rounded text-sm"
                                placeholder="0.00"
                              />
                            </div>
                          )}
                        </div>
                      </div>
                    );
                  })}
                </div>
              </div>

              {expenseForm.amount && splitPreview.length > 0 && (
                <div className="bg-gray-50 rounded-lg p-4">
                  <p className="text-sm font-medium text-gray-700 mb-2">分摊预览</p>
                  <div className="grid grid-cols-2 gap-2">
                    {splitPreview.map((item, index) => (
                      <div key={index} className="flex justify-between text-sm">
                        <span className="text-gray-600">{item.name}</span>
                        <span className="font-medium">¥{item.amount.toFixed(2)}</span>
                      </div>
                    ))}
                  </div>
                  {expenseForm.splitType === 'custom' && (
                    <div className={`mt-2 pt-2 border-t border-gray-200 flex justify-between text-sm ${
                      Math.abs(totalPreview - parseFloat(expenseForm.amount)) > 0.01
                        ? 'text-red-500'
                        : 'text-green-600'
                    }`}>
                      <span>合计: ¥{totalPreview.toFixed(2)}</span>
                      <span>
                        {Math.abs(totalPreview - parseFloat(expenseForm.amount)) > 0.01
                          ? `差额: ¥${(parseFloat(expenseForm.amount) - totalPreview).toFixed(2)}`
                          : '✓ 金额匹配'}
                      </span>
                    </div>
                  )}
                </div>
              )}

              {error && (
                <div className="bg-red-50 text-red-600 px-4 py-3 rounded-lg text-sm">
                  {error}
                </div>
              )}

              <div className="flex gap-3 pt-4">
                <button
                  type="button"
                  onClick={() => setShowAddExpense(false)}
                  className="flex-1 px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors"
                >
                  取消
                </button>
                <button
                  type="submit"
                  className="flex-1 px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition-colors"
                >
                  添加
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {showInviteModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-xl w-full max-w-md p-6">
            <h2 className="text-xl font-bold text-gray-800 mb-4">邀请成员</h2>
            <p className="text-gray-600 mb-6">
              将邀请码分享给好友，他们可以通过邀请码加入群组。
            </p>
            <div className="bg-gray-100 rounded-lg p-4 text-center">
              <p className="text-sm text-gray-500 mb-2">邀请码</p>
              <p className="text-3xl font-bold text-primary-600 tracking-widest font-mono">
                {group.inviteCode}
              </p>
            </div>
            <button
              onClick={copyInviteCode}
              className="w-full mt-6 px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition-colors"
            >
              {copied ? '✓ 已复制' : '复制邀请码'}
            </button>
            <button
              onClick={() => setShowInviteModal(false)}
              className="w-full mt-2 px-4 py-2 text-gray-600 hover:text-gray-800 transition-colors"
            >
              关闭
            </button>
          </div>
        </div>
      )}
    </div>
  );
};

export default GroupDetail;
