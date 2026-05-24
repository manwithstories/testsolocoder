import React from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import { useAuth } from './context/AuthContext';
import Layout from './components/Layout';
import Login from './pages/Login';
import Dashboard from './pages/Dashboard';
import Templates from './pages/Templates';
import Orders from './pages/Orders';
import OrderCustomize from './pages/OrderCustomize';
import Pricing from './pages/Pricing';
import Production from './pages/Production';
import Customers from './pages/Customers';
import Invoices from './pages/Invoices';

function PrivateRoute({ children }: { children: React.ReactNode }) {
  const { token, loading } = useAuth();
  if (loading) return <div>加载中...</div>;
  return token ? <>{children}</> : <Navigate to="/login" replace />;
}

export default function App() {
  return (
    <Routes>
      <Route path="/login" element={<Login />} />
      <Route path="/" element={<PrivateRoute><Layout /></PrivateRoute>}>
        <Route index element={<Dashboard />} />
        <Route path="templates" element={<Templates />} />
        <Route path="orders" element={<Orders />} />
        <Route path="orders/new" element={<OrderCustomize />} />
        <Route path="pricing" element={<Pricing />} />
        <Route path="production" element={<Production />} />
        <Route path="customers" element={<Customers />} />
        <Route path="invoices" element={<Invoices />} />
      </Route>
    </Routes>
  );
}
