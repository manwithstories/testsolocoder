import React, { useEffect } from 'react';
import { BrowserRouter, Routes, Route, useNavigate, useLocation } from 'react-router-dom';
import { ConfigProvider, message } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import { Provider, useDispatch } from 'react-redux';
import { store } from './store';
import { setUser } from './store/authSlice';
import MainLayout from './components/MainLayout';
import ProtectedRoute from './components/ProtectedRoute';
import AuthPage from './pages/AuthPage';
import HomePage from './pages/HomePage';
import ModelMarketPage from './pages/ModelMarketPage';
import ModelDetailPage from './pages/ModelDetailPage';
import ModelUploadPage from './pages/ModelUploadPage';
import MyModelsPage from './pages/MyModelsPage';
import MyOrdersPage from './pages/MyOrdersPage';
import PrinterOrdersPage from './pages/PrinterOrdersPage';
import PrinterDevicesPage from './pages/PrinterDevicesPage';
import StatsPage from './pages/StatsPage';
import ProfilePage from './pages/ProfilePage';
import WalletPage from './pages/WalletPage';
import FavoritesPage from './pages/FavoritesPage';
import NotFoundPage from './pages/NotFoundPage';
import ForbiddenPage from './pages/ForbiddenPage';
import { User } from './types';

const AuthInitializer: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const dispatch = useDispatch();

  useEffect(() => {
    const userStr = localStorage.getItem('user');
    if (userStr) {
      try {
        const user = JSON.parse(userStr) as User;
        dispatch(setUser(user));
      } catch (e) {
        localStorage.removeItem('user');
      }
    }
  }, [dispatch]);

  return <>{children}</>;
};

const ScrollToTop: React.FC = () => {
  const { pathname } = useLocation();

  useEffect(() => {
    window.scrollTo(0, 0);
  }, [pathname]);

  return null;
};

const AppContent: React.FC = () => {
  return (
    <>
      <ScrollToTop />
      <Routes>
        <Route path="/login" element={<AuthPage mode="login" />} />
        <Route path="/register" element={<AuthPage mode="register" />} />
        <Route path="/403" element={<ForbiddenPage />} />

        <Route
          path="/"
          element={
            <MainLayout>
              <HomePage />
            </MainLayout>
          }
        />
        <Route
          path="/models"
          element={
            <MainLayout>
              <ModelMarketPage />
            </MainLayout>
          }
        />
        <Route
          path="/models/:id"
          element={
            <MainLayout>
              <ModelDetailPage />
            </MainLayout>
          }
        />

        <Route
          path="/model-upload"
          element={
            <ProtectedRoute roles={['designer', 'admin']}>
              <MainLayout>
                <ModelUploadPage />
              </MainLayout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/my-models"
          element={
            <ProtectedRoute roles={['designer', 'admin']}>
              <MainLayout>
                <MyModelsPage />
              </MainLayout>
            </ProtectedRoute>
          }
        />

        <Route
          path="/my-orders"
          element={
            <ProtectedRoute roles={['customer', 'admin']}>
              <MainLayout>
                <MyOrdersPage />
              </MainLayout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/my-orders/:id"
          element={
            <ProtectedRoute roles={['customer', 'admin']}>
              <MainLayout>
                <MyOrdersPage />
              </MainLayout>
            </ProtectedRoute>
          }
        />

        <Route
          path="/printer-orders"
          element={
            <ProtectedRoute roles={['printer', 'admin']}>
              <MainLayout>
                <PrinterOrdersPage />
              </MainLayout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/printer-devices"
          element={
            <ProtectedRoute roles={['printer', 'admin']}>
              <MainLayout>
                <PrinterDevicesPage />
              </MainLayout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/printer-inventory"
          element={
            <ProtectedRoute roles={['printer', 'admin']}>
              <MainLayout>
                <PrinterDevicesPage />
              </MainLayout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/printer-schedules"
          element={
            <ProtectedRoute roles={['printer', 'admin']}>
              <MainLayout>
                <PrinterDevicesPage />
              </MainLayout>
            </ProtectedRoute>
          }
        />

        <Route
          path="/stats"
          element={
            <ProtectedRoute roles={['designer', 'printer', 'admin']}>
              <MainLayout>
                <StatsPage />
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
          path="/settings"
          element={
            <ProtectedRoute>
              <MainLayout>
                <ProfilePage />
              </MainLayout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/wallet"
          element={
            <ProtectedRoute>
              <MainLayout>
                <WalletPage />
              </MainLayout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/favorites"
          element={
            <ProtectedRoute>
              <MainLayout>
                <FavoritesPage />
              </MainLayout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/purchases"
          element={
            <ProtectedRoute>
              <MainLayout>
                <MyModelsPage />
              </MainLayout>
            </ProtectedRoute>
          }
        />

        <Route path="*" element={<NotFoundPage />} />
      </Routes>
    </>
  );
};

const App: React.FC = () => {
  return (
    <Provider store={store}>
      <ConfigProvider
        locale={zhCN}
        theme={{
          token: {
            colorPrimary: '#3b82f6',
            colorInfo: '#3b82f6',
            borderRadius: 8,
          },
        }}
      >
        <BrowserRouter>
          <AuthInitializer>
            <AppContent />
          </AuthInitializer>
        </BrowserRouter>
      </ConfigProvider>
    </Provider>
  );
};

export default App;
