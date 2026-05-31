import { Routes, Route, Navigate } from 'react-router-dom'
import { useAuthStore } from '@/store/auth'
import { useEffect } from 'react'

import Layout from '@/components/Layout'
import ProtectedRoute from '@/components/ProtectedRoute'

import Login from '@/pages/Login'
import Register from '@/pages/Register'
import Home from '@/pages/Home'
import Products from '@/pages/Products'
import ProductDetail from '@/pages/ProductDetail'
import Cart from '@/pages/Cart'
import Orders from '@/pages/Orders'
import OrderDetail from '@/pages/OrderDetail'
import RoastingRecords from '@/pages/RoastingRecords'
import RoastingDetail from '@/pages/RoastingDetail'
import Cupping from '@/pages/Cupping'
import Roasters from '@/pages/Roasters'
import RoasterProfile from '@/pages/RoasterProfile'
import MyProfile from '@/pages/MyProfile'
import MyCertification from '@/pages/MyCertification'
import Search from '@/pages/Search'

import AdminDashboard from '@/pages/admin/AdminDashboard'
import AdminUsers from '@/pages/admin/AdminUsers'
import AdminProducts from '@/pages/admin/AdminProducts'
import AdminCertifications from '@/pages/admin/AdminCertifications'
import AdminStats from '@/pages/admin/AdminStats'

import RoasterDashboard from '@/pages/roaster/RoasterDashboard'
import RoasterProducts from '@/pages/roaster/RoasterProducts'
import RoasterRoastingRecords from '@/pages/roaster/RoasterRoastingRecords'

function App() {
  const { isAuthenticated, loadUser } = useAuthStore()

  useEffect(() => {
    if (isAuthenticated) {
      loadUser()
    }
  }, [isAuthenticated, loadUser])

  return (
    <Routes>
      <Route path="/login" element={<Login />} />
      <Route path="/register" element={<Register />} />

      <Route path="/" element={<Layout />}>
        <Route index element={<Home />} />
        <Route path="products" element={<Products />} />
        <Route path="products/:id" element={<ProductDetail />} />
        <Route path="roasters" element={<Roasters />} />
        <Route path="roasters/:id" element={<RoasterProfile />} />
        <Route path="search" element={<Search />} />

        <Route element={<ProtectedRoute />}>
          <Route path="cart" element={<Cart />} />
          <Route path="orders" element={<Orders />} />
          <Route path="orders/:id" element={<OrderDetail />} />
          <Route path="cupping" element={<Cupping />} />
          <Route path="profile" element={<MyProfile />} />
          <Route path="certification" element={<MyCertification />} />
          <Route path="roasting" element={<RoastingRecords />} />
          <Route path="roasting/:id" element={<RoastingDetail />} />
        </Route>

        <Route element={<ProtectedRoute roles={['admin']} />}>
          <Route path="admin" element={<AdminDashboard />} />
          <Route path="admin/users" element={<AdminUsers />} />
          <Route path="admin/products" element={<AdminProducts />} />
          <Route path="admin/certifications" element={<AdminCertifications />} />
          <Route path="admin/stats" element={<AdminStats />} />
        </Route>

        <Route element={<ProtectedRoute roles={['admin', 'roaster']} />}>
          <Route path="roaster/dashboard" element={<RoasterDashboard />} />
          <Route path="roaster/products" element={<RoasterProducts />} />
          <Route path="roaster/roasting" element={<RoasterRoastingRecords />} />
        </Route>

        <Route path="*" element={<Navigate to="/" replace />} />
      </Route>
    </Routes>
  )
}

export default App
