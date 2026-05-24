import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { ConfigProvider, App as AntApp } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import { MainLayout, AdminLayout } from './components/Layouts';
import { LoginPage, RegisterPage } from './pages/AuthPages';
import { HomePage } from './pages/HomePage';
import { TextbookListPage, TextbookDetailPage } from './pages/TextbookPages';
import { NoteListPage, NoteDetailPage } from './pages/NotePages';
import { OrderListPage } from './pages/OrderPages';
import { MessagePage } from './pages/MessagePage';
import { ProfilePage } from './pages/ProfilePage';
import {
  AdminDashboard,
  AdminUsersPage,
  AdminTextbooksPage,
  AdminNotesPage,
  AdminOrdersPage,
  AdminReviewsPage,
} from './pages/AdminPages';
import { useAuthStore } from './context/authStore';

const ProtectedRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { isAuthenticated } = useAuthStore();
  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }
  return <>{children}</>;
};

const AdminRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { isAuthenticated, user } = useAuthStore();
  if (!isAuthenticated || user?.role !== 'admin') {
    return <Navigate to="/" replace />;
  }
  return <>{children}</>;
};

const App: React.FC = () => {
  return (
    <ConfigProvider locale={zhCN}>
      <AntApp>
        <BrowserRouter>
          <Routes>
            <Route path="/login" element={<LoginPage />} />
            <Route path="/register" element={<RegisterPage />} />

            <Route
              path="/"
              element={
                <MainLayout>
                  <HomePage />
                </MainLayout>
              }
            />

            <Route
              path="/textbooks"
              element={
                <MainLayout>
                  <TextbookListPage />
                </MainLayout>
              }
            />
            <Route
              path="/textbooks/:id"
              element={
                <MainLayout>
                  <TextbookDetailPage />
                </MainLayout>
              }
            />

            <Route
              path="/notes"
              element={
                <MainLayout>
                  <NoteListPage />
                </MainLayout>
              }
            />
            <Route
              path="/notes/:id"
              element={
                <MainLayout>
                  <NoteDetailPage />
                </MainLayout>
              }
            />

            <Route
              path="/orders"
              element={
                <ProtectedRoute>
                  <MainLayout>
                    <OrderListPage />
                  </MainLayout>
                </ProtectedRoute>
              }
            />

            <Route
              path="/messages"
              element={
                <ProtectedRoute>
                  <MainLayout>
                    <MessagePage />
                  </MainLayout>
                </ProtectedRoute>
              }
            />

            <Route
              path="/profile"
              element={
                <ProtectedRoute>
                  <MainLayout>
                    <ProfilePage />
                  </MainLayout>
                </ProtectedRoute>
              }
            />

            <Route
              path="/admin/dashboard"
              element={
                <AdminRoute>
                  <AdminLayout>
                    <AdminDashboard />
                  </AdminLayout>
                </AdminRoute>
              }
            />
            <Route
              path="/admin/users"
              element={
                <AdminRoute>
                  <AdminLayout>
                    <AdminUsersPage />
                  </AdminLayout>
                </AdminRoute>
              }
            />
            <Route
              path="/admin/textbooks"
              element={
                <AdminRoute>
                  <AdminLayout>
                    <AdminTextbooksPage />
                  </AdminLayout>
                </AdminRoute>
              }
            />
            <Route
              path="/admin/notes"
              element={
                <AdminRoute>
                  <AdminLayout>
                    <AdminNotesPage />
                  </AdminLayout>
                </AdminRoute>
              }
            />
            <Route
              path="/admin/orders"
              element={
                <AdminRoute>
                  <AdminLayout>
                    <AdminOrdersPage />
                  </AdminLayout>
                </AdminRoute>
              }
            />
            <Route
              path="/admin/reviews"
              element={
                <AdminRoute>
                  <AdminLayout>
                    <AdminReviewsPage />
                  </AdminLayout>
                </AdminRoute>
              }
            />

            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </BrowserRouter>
      </AntApp>
    </ConfigProvider>
  );
};

export default App;
