import React, { createContext, useContext, useState, useEffect } from 'react';
import type { ReactNode } from 'react';

interface User {
  id: string;
  email: string;
  display_name?: string;
}

interface AuthContextType {
  user: User | null;
  isLoading: boolean;
  login: (email: string, password: string) => Promise<boolean>;
  signup: (email: string, password: string, displayName?: string) => Promise<{ success: boolean; error?: string }>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

interface AuthProviderProps {
  children: ReactNode;
}

export function AuthProvider({ children }: AuthProviderProps) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  // Check for existing session on mount
  useEffect(() => {
    checkAuthStatus();
  }, []);

  const checkAuthStatus = async () => {
    try {
      const response = await fetch('/auth/user', {
        credentials: 'include',
      });
      
      if (response.ok) {
        const userData = await response.json();
        // Transform the auth service response to our User interface
        setUser({
          id: userData.id,
          email: userData.name, // The auth service uses 'name' for email
          display_name: userData.display_name || userData.name,
        });
      }
    } catch (error) {
      console.error('Auth check failed:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const login = async (email: string, password: string): Promise<boolean> => {
    try {
      const formData = new URLSearchParams();
      formData.append('user', email);
      formData.append('passwd', password);

      const response = await fetch('/auth/local/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
        credentials: 'include',
        body: formData,
      });

      if (response.ok) {
        await checkAuthStatus();
        return true;
      }
      
      // Log response details for debugging
      console.error('Login failed:', response.status, response.statusText);
      const responseText = await response.text();
      console.error('Response body:', responseText);
      
      return false;
    } catch (error) {
      console.error('Login failed:', error);
      return false;
    }
  };

  const signup = async (email: string, password: string, displayName?: string): Promise<{ success: boolean; error?: string }> => {
    try {
      // First, create the account
      const signupResponse = await fetch('/signup', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          email,
          password,
          display_name: displayName,
        }),
      });

      if (!signupResponse.ok) {
        const errorData = await signupResponse.json();
        return { 
          success: false, 
          error: errorData.message || `HTTP ${signupResponse.status}: ${signupResponse.statusText}` 
        };
      }

      // After successful signup, automatically log in
      const loginSuccess = await login(email, password);
      if (!loginSuccess) {
        return { 
          success: false, 
          error: 'Something went wrong. Please try again.' 
        };
      }

      return { success: true };
    } catch (error) {
      console.error('Signup failed:', error);
      return { success: false, error: 'Network error occurred' };
    }
  };

  const logout = async () => {
    try {
      await fetch('/auth/logout', {
        method: 'POST',
        credentials: 'include',
      });
    } catch (error) {
      console.error('Logout failed:', error);
    } finally {
      setUser(null);
    }
  };

  const value: AuthContextType = {
    user,
    isLoading,
    login,
    signup,
    logout,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}