import { Navigate } from 'react-router-dom'
import Login from './Login'

export default function Register() {
  return <Navigate to="/login" replace />
}
