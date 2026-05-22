import React from 'react'
import { Layout } from 'antd'

const { Footer } = Layout

const AppFooter: React.FC = () => {
  return (
    <Footer
      style={{
        textAlign: 'center',
        background: '#fff',
        borderTop: '1px solid #f0f0f0',
        padding: '16px',
      }}
    >
      在线学习与课程管理平台 ©{new Date().getFullYear()} All Rights Reserved
    </Footer>
  )
}

export default AppFooter
