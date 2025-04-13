import { jwtDecode } from "jwt-decode"
import { apiConfig } from '@/api/config'
import { apiRequest } from '@/api/api'

interface User {
  id: string
  email: string
  companyName: string
  exp: number
}

const validateSession = async (token: string) => {
  try {
    if (!token) {
      throw new Error("No token found")
    }

    const response = await apiRequest({
      method: 'GET',
      endpoint: apiConfig.endpoints.auth.validateSession,
      isSecure: false,
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });

    if (!response.success) {
      throw new Error(response.message || 'Session validation failed');
    }

    return Promise.resolve();
  } catch (error) {
    console.error("Session validation error:", error);
    return Promise.reject(error);
  } finally {
  }
}

export async function verifyAuth(token: string): Promise<boolean> {
  try {

    if (!token) {
      return false
    }

    await validateSession(token)
    // const decoded = jwtDecode<User>(token)

    // // await validateSession(token)
    // if (decoded.exp * 1000 < Date.now()) {
    //   return false
    // }

    return true
  } catch (error) {
    console.error("Token verification failed:", error)
    return false
  }
}

export async function getUser(token: string): Promise<User | null> {
  try {
    const decoded = jwtDecode<User>(token)

    if (decoded.exp * 1000 < Date.now()) {
      return null
    }

    return decoded
  } catch (error) {
    console.error("Failed to get user from token:", error)
    return null
  }
}

