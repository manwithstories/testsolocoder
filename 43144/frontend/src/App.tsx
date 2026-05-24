import React from 'react'
import { Routes, Route, Navigate } from 'react-router-dom'
import MainLayout from './components/Layout'
import ProtectedRoute from './components/ProtectedRoute'
import Home from './pages/Home'
import Login from './pages/Login'
import Register from './pages/Register'
import PetList from './pages/PetList'
import PetDetail from './pages/PetDetail'
import MyPets from './pages/MyPets'
import MyAdoptedPets from './pages/MyAdoptedPets'
import AdoptionApplications from './pages/AdoptionApplications'
import HealthRecords from './pages/HealthRecords'
import Appointments from './pages/Appointments'
import RescueStats from './pages/RescueStats'
import Profile from './pages/Profile'
import AdminRescues from './pages/AdminRescues'
import AdminUsers from './pages/AdminUsers'

const App: React.FC = () => {
  return (
    <Routes>
      <Route path="/login" element={<Login />} />
      <Route path="/register" element={<Register />} />
      <Route
        path="/"
        element={
          <ProtectedRoute>
            <MainLayout>
              <Home />
            </MainLayout>
          </ProtectedRoute>
        }
      />
      <Route
        path="/pets"
        element={
          <ProtectedRoute>
            <MainLayout>
              <PetList />
            </MainLayout>
          </ProtectedRoute>
        }
      />
      <Route
        path="/pets/:id"
        element={
          <ProtectedRoute>
            <MainLayout>
              <PetDetail />
            </MainLayout>
          </ProtectedRoute>
        }
      />
      <Route
        path="/my-pets"
        element={
          <ProtectedRoute roles={['rescue', 'admin']}>
            <MainLayout>
              <MyPets />
            </MainLayout>
          </ProtectedRoute>
        }
      />
      <Route
        path="/adopted"
        element={
          <ProtectedRoute roles={['adopter']}>
            <MainLayout>
              <MyAdoptedPets />
            </MainLayout>
          </ProtectedRoute>
        }
      />
      <Route
        path="/adoption-applications"
        element={
          <ProtectedRoute roles={['rescue', 'admin']}>
            <MainLayout>
              <AdoptionApplications />
            </MainLayout>
          </ProtectedRoute>
        }
      />
      <Route
        path="/health-records"
        element={
          <ProtectedRoute>
            <MainLayout>
              <HealthRecords />
            </MainLayout>
          </ProtectedRoute>
        }
      />
      <Route
        path="/appointments"
        element={
          <ProtectedRoute>
            <MainLayout>
              <Appointments />
            </MainLayout>
          </ProtectedRoute>
        }
      />
      <Route
        path="/rescue-stats"
        element={
          <ProtectedRoute roles={['rescue', 'admin']}>
            <MainLayout>
              <RescueStats />
            </MainLayout>
          </ProtectedRoute>
        }
      />
      <Route
        path="/profile"
        element={
          <ProtectedRoute>
            <MainLayout>
              <Profile />
            </MainLayout>
          </ProtectedRoute>
        }
      />
      <Route
        path="/admin/rescues"
        element={
          <ProtectedRoute roles={['admin']}>
            <MainLayout>
              <AdminRescues />
            </MainLayout>
          </ProtectedRoute>
        }
      />
      <Route
        path="/admin/users"
        element={
          <ProtectedRoute roles={['admin']}>
            <MainLayout>
              <AdminUsers />
            </MainLayout>
          </ProtectedRoute>
        }
      />
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}

export default App
