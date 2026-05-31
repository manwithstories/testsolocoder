import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { useAuthStore } from './store/authStore'
import MainLayout from './components/layouts/MainLayout'
import LoginPage from './pages/auth/LoginPage'
import RegisterPage from './pages/auth/RegisterPage'
import DashboardPage from './pages/DashboardPage'
import BeehivesPage from './pages/beehive/BeehivesPage'
import BeehiveDetailPage from './pages/beehive/BeehiveDetailPage'
import HealthRecordsPage from './pages/health/HealthRecordsPage'
import HarvestsPage from './pages/harvest/HarvestsPage'
import InventoryPage from './pages/inventory/InventoryPage'
import ProductsPage from './pages/product/ProductsPage'
import ProductDetailPage from './pages/product/ProductDetailPage'
import MyProductsPage from './pages/product/MyProductsPage'
import OrdersPage from './pages/order/OrdersPage'
import OrderDetailPage from './pages/order/OrderDetailPage'
import InspectionsPage from './pages/inspection/InspectionsPage'
import CommunityPage from './pages/community/CommunityPage'
import PostDetailPage from './pages/community/PostDetailPage'
import AnalyticsPage from './pages/analytics/AnalyticsPage'

function PrivateRoute({ children }: { children: React.ReactNode }) {
  const { isAuthenticated } = useAuthStore()
  return isAuthenticated ? children : <Navigate to="/login" />
}

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />
        
        <Route path="/" element={<PrivateRoute><MainLayout /></PrivateRoute>}>
          <Route index element={<DashboardPage />} />
          
          <Route path="beehives" element={<BeehivesPage />} />
          <Route path="beehives/:id" element={<BeehiveDetailPage />} />
          
          <Route path="health-records" element={<HealthRecordsPage />} />
          
          <Route path="harvests" element={<HarvestsPage />} />
          
          <Route path="inventory" element={<InventoryPage />} />
          
          <Route path="products" element={<ProductsPage />} />
          <Route path="products/:id" element={<ProductDetailPage />} />
          <Route path="my-products" element={<MyProductsPage />} />
          
          <Route path="orders" element={<OrdersPage />} />
          <Route path="orders/:id" element={<OrderDetailPage />} />
          
          <Route path="inspections" element={<InspectionsPage />} />
          
          <Route path="community" element={<CommunityPage />} />
          <Route path="community/:id" element={<PostDetailPage />} />
          
          <Route path="analytics" element={<AnalyticsPage />} />
        </Route>
      </Routes>
    </BrowserRouter>
  )
}

export default App
