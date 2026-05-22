import React from 'react'
import { Outlet } from 'react-router-dom'
import { Layout } from 'antd'
import AppHeader from './AppHeader'
import AppFooter from './AppFooter'

const { Content } = Layout

const MainLayout: React.FC = () => {
  return (
    <Layout>
      <AppHeader />
      <Content style={{ minHeight: 'calc(100vh - 134px)', padding: '24px', background: '#f5f5f5' }}>
        <div style={{ maxWidth: 1400, margin: '0 auto' }}>
          <Outlet />
        </div>
      </Content>
      <AppFooter />
    </Layout>
  )
}

export default MainLayout
