import { Routes, Route, Navigate } from 'react-router-dom'
import { useEffect } from 'react'
import { useAppStore } from './store'
import MainLayout from './components/Layout/MainLayout'
import SellerLayout from './components/Layout/SellerLayout'
import AdminLayout from './components/Layout/AdminLayout'
import Login from './pages/Auth/Login'
import Register from './pages/Auth/Register'
import Home from './pages/Home'
import ProductList from './pages/Product/List'
import ProductDetail from './pages/Product/Detail'
import ShopDetail from './pages/Shop/Detail'
import Cart from './pages/Cart'
import Orders from './pages/Order/List'
import OrderDetail from './pages/Order/Detail'
import Favorites from './pages/Favorite'
import Notifications from './pages/Notification'
import SellerShop from './pages/Seller/Shop'
import SellerProducts from './pages/Seller/Products'
import SellerProductCreate from './pages/Seller/ProductCreate'
import SellerOrders from './pages/Seller/Orders'
import SellerStatistics from './pages/Seller/Statistics'
import AdminShops from './pages/Admin/Shops'
import AdminCategories from './pages/Admin/Categories'
import AdminDisputes from './pages/Admin/Disputes'
import AdminStatistics from './pages/Admin/Statistics'
import AdminUsers from './pages/Admin/Users'

const App = () => {
  const { token, loadUser, loadCart } = useAppStore()

  useEffect(() => {
    if (token) {
      loadUser()
      loadCart()
    }
  }, [token])

  return (
    <Routes>
      <Route path="/login" element={<Login />} />
      <Route path="/register" element={<Register />} />

      <Route path="/" element={<MainLayout />}>
        <Route index element={<Home />} />
        <Route path="products" element={<ProductList />} />
        <Route path="products/:id" element={<ProductDetail />} />
        <Route path="shops/:id" element={<ShopDetail />} />
        <Route path="cart" element={token ? <Cart /> : <Navigate to="/login" />} />
        <Route path="orders" element={token ? <Orders /> : <Navigate to="/login" />} />
        <Route path="orders/:id" element={token ? <OrderDetail /> : <Navigate to="/login" />} />
        <Route path="favorites" element={token ? <Favorites /> : <Navigate to="/login" />} />
        <Route path="notifications" element={token ? <Notifications /> : <Navigate to="/login" />} />
      </Route>

      <Route path="/seller" element={token ? <SellerLayout /> : <Navigate to="/login" />}>
        <Route index element={<Navigate to="/seller/shop" replace />} />
        <Route path="shop" element={<SellerShop />} />
        <Route path="products" element={<SellerProducts />} />
        <Route path="products/create" element={<SellerProductCreate />} />
        <Route path="orders" element={<SellerOrders />} />
        <Route path="statistics" element={<SellerStatistics />} />
      </Route>

      <Route path="/admin" element={token ? <AdminLayout /> : <Navigate to="/login" />}>
        <Route index element={<Navigate to="/admin/statistics" replace />} />
        <Route path="statistics" element={<AdminStatistics />} />
        <Route path="shops" element={<AdminShops />} />
        <Route path="categories" element={<AdminCategories />} />
        <Route path="disputes" element={<AdminDisputes />} />
        <Route path="users" element={<AdminUsers />} />
      </Route>

      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}

export default App
