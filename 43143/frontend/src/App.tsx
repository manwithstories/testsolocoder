import React, { useEffect } from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { ConfigProvider } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import { useSelector } from 'react-redux';
import { RootState } from '@/store';

import MainLayout from '@/components/layout/MainLayout';
import Login from '@/pages/Login';
import Register from '@/pages/Register';
import Home from '@/pages/Home';
import SkillDetail from '@/pages/SkillDetail';
import SkillList from '@/pages/SkillList';
import PostingDetail from '@/pages/PostingDetail';
import PostingCreate from '@/pages/PostingCreate';
import BookingList from '@/pages/BookingList';
import BookingDetail from '@/pages/BookingDetail';
import Messages from '@/pages/Messages';
import Profile from '@/pages/Profile';
import ProfileEdit from '@/pages/ProfileEdit';
import Wallet from '@/pages/Wallet';
import SchedulePage from '@/pages/Schedule';
import Statistics from '@/pages/Statistics';

const PrivateRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const isAuthenticated = useSelector((state: RootState) => state.auth.isAuthenticated);

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  return <>{children}</>;
};

const App: React.FC = () => {
  useEffect(() => {
    document.title = '技能共享平台';
  }, []);

  return (
    <ConfigProvider locale={zhCN}>
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />

          <Route
            path="/"
            element={
              <PrivateRoute>
                <MainLayout />
              </PrivateRoute>
            }
          >
            <Route index element={<Home />} />
            <Route path="skills" element={<SkillList />} />
            <Route path="skills/:id" element={<SkillDetail />} />
            <Route path="postings/create" element={<PostingCreate />} />
            <Route path="postings/:id" element={<PostingDetail />} />
            <Route path="bookings" element={<BookingList />} />
            <Route path="bookings/:id" element={<BookingDetail />} />
            <Route path="messages" element={<Messages />} />
            <Route path="profile" element={<Profile />} />
            <Route path="profile/edit" element={<ProfileEdit />} />
            <Route path="wallet" element={<Wallet />} />
            <Route path="schedule" element={<SchedulePage />} />
            <Route path="statistics" element={<Statistics />} />
          </Route>

          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </BrowserRouter>
    </ConfigProvider>
  );
};

export default App;
