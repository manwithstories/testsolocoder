import { useEffect, useState } from 'react';
import { jwtDecode } from 'jwt-decode';
import { JWTPayload, UserRole } from '@/types';

export function useAuth() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [userRole, setUserRole] = useState<UserRole | null>(null);
  const [userId, setUserId] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = localStorage.getItem('access_token');
    if (token) {
      try {
        const payload = jwtDecode<JWTPayload>(token);
        const isExpired = payload.exp * 1000 < Date.now();
        
        if (!isExpired) {
          setIsAuthenticated(true);
          setUserRole(payload.role);
          setUserId(payload.user_id);
        } else {
          localStorage.removeItem('access_token');
          localStorage.removeItem('refresh_token');
          localStorage.removeItem('user');
        }
      } catch (error) {
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        localStorage.removeItem('user');
      }
    }
    setLoading(false);
  }, []);

  const logout = () => {
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    localStorage.removeItem('user');
    setIsAuthenticated(false);
    setUserRole(null);
    setUserId(null);
  };

  return {
    isAuthenticated,
    userRole,
    userId,
    loading,
    logout,
  };
}

export function useAuthGuard(requiredRoles?: UserRole[]) {
  const { isAuthenticated, userRole, loading } = useAuth();
  const [hasAccess, setHasAccess] = useState(false);

  useEffect(() => {
    if (!loading) {
      if (!isAuthenticated) {
        setHasAccess(false);
      } else if (requiredRoles && userRole) {
        setHasAccess(requiredRoles.includes(userRole));
      } else {
        setHasAccess(isAuthenticated);
      }
    }
  }, [isAuthenticated, userRole, loading, requiredRoles]);

  return { hasAccess, loading, isAuthenticated, userRole };
}
