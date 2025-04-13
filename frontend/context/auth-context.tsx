"use client"

import { apiConfig } from '@/api/config';
import { jwtDecode } from "jwt-decode"
import { apiRequest } from '@/api/api';
import { createContext, useContext, useState, useEffect, type ReactNode } from "react"

interface User {
  id: string
  email: string
  companyName: string
}

interface UserResponse {
  id: string;
  email: string;
  name: string;
}

interface LoginResponse {
  success: boolean;
  message: string;
  user: UserResponse;
  token?: string;
}

interface RefreshSessionResponse {
  success: boolean;
  message: string;
}

interface RegisterData {
  companyName: string
  email: string
  phone: string
  document: string
  password: string
}

interface AuthContextType {
  userSet: User | null
  isLoading: boolean
  isAuthenticated: boolean
  login: (email: string, password: string) => Promise<void>
  logout: () => void
  register: (data: RegisterData) => Promise<void>
  refreshSession: (id: string, email: string) => Promise<void>
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [userSet, setUser] = useState<User | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [isAuthenticated, setIsAuthenticated] = useState(false)

  useEffect(() => {
    let isMounted = true

    const checkAuth = async () => {
      const token = getCookie("auth-token")

      if (token) {
        try {
          const userData = getCookie("user-data")
          const user = JSON.parse(userData)

          const decodedToken = jwtDecode<{ exp: number }>(token)
          if (decodedToken.exp * 1000 < Date.now()) {
            deleteCookie('user-data')
            deleteCookie('auth-token')
            throw new Error("Token expired")
          }

          if (isMounted) {
            setUser(user)
            setIsAuthenticated(true)
            saveUserToCookie(user)
          }
        } catch (error) {
          console.error("Auth check error:", error)
          if (isMounted) {
            setUser(null)
            setIsAuthenticated(false)
            deleteCookie('user-data')
            deleteCookie('auth-token')
          }
        }
      }

      if (isMounted) {
        setIsLoading(false)
      }
    }

    checkAuth()

    return () => {
      isMounted = false
    }
  }, [])

  const refreshSession = async (id: string, email: string) => {
    setIsLoading(true)

    try {
      const response = await apiRequest<RefreshSessionResponse>({
        method: 'POST',
        endpoint: apiConfig.endpoints.auth.refreshToken,
        body: { id, email },
        isSecure: true,
      });

      if (!response.success) {
        throw new Error(response.message || 'Refresh session failed');
      }

      setIsAuthenticated(true);

      return Promise.resolve();
    } catch (error) {
      console.error("Refresh error:", error);
      setIsAuthenticated(false);
      return Promise.reject(error);
    } finally {
      setIsLoading(false);
    }
  }

  const login = async (email: string, password: string) => {
    setIsLoading(true)

    try {
      const response = await apiRequest<LoginResponse>({
        method: 'POST',
        endpoint: apiConfig.endpoints.auth.customerLogin,
        body: { email, password },
        isSecure: true,
      });

      const responseData = typeof response.data === 'string'
        ? JSON.parse(response.data)
        : response.data;

      if (!response.success || !responseData?.user) {
        throw new Error(response.message || 'Login failed');
      }

      const user = {
        id: responseData.user.id,
        email: responseData.user.email,
        companyName: responseData.user.name,
      };

      setUser(user);
      setIsAuthenticated(true);
      saveUserToCookie(user);

      return Promise.resolve();
    } catch (error) {
      console.error("Login error:", error);
      setIsAuthenticated(false);
      return Promise.reject(error);
    } finally {
      setIsLoading(false);
    }
  }

  const register = async (data: RegisterData) => {
    setIsLoading(true)
    try {
      const response = await apiRequest({
        method: 'POST',
        endpoint: apiConfig.endpoints.auth.customer,
        body: {
          name: data.companyName,
          email: data.email,
          phone: data.phone,
          cpf_cnpj: data.document,
          password: data.password
        },
        isSecure: true
      })

      if (!response.success) {
        throw new Error(response.message || 'Registration error. Please try again.')
      }

      return
    } catch (error) {
      console.error("Registration error:", error)
      throw error
    } finally {
      setIsLoading(false)
    }
  }

  const logout = async () => {
    try {
      const token = getCookie("auth-token");

      await apiRequest({
        method: 'POST',
        endpoint: apiConfig.endpoints.auth.logout,
        isSecure: true,
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });
    } catch (error) {
      console.error("Logout error:", error);
    } finally {
      setUser(null);
      setIsAuthenticated(false);
      deleteCookie('user-data');
      deleteCookie('auth-token');

      window.location.href = "/";
    }
  }

  return (
    <AuthContext.Provider
      value={{
        userSet,
        isAuthenticated,
        isLoading,
        login,
        register,
        logout,
        refreshSession
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider")
  }
  return context
}

export function getCookie(name: string): string {
  if (typeof window === 'undefined') return "";
  const value = `; ${document.cookie}`;
  const parts = value.split(`; ${name}=`);
  if (parts.length === 2) return parts.pop()?.split(';').shift() || "";
  return "";
}

export function deleteCookie(name: string) {
  if (typeof window === 'undefined') return;
  document.cookie = `${name}=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;`;
}

function saveUserToCookie(user: User) {
  if (typeof window === 'undefined') return;
  document.cookie = `user-data=${JSON.stringify(user)}; path=/; max-age=31536000`;
}
