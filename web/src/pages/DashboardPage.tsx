import React from 'react';
import { Button } from '@/components/ui/button';
import { useAuth } from '@/contexts/AuthContext';

export function DashboardPage() {
  const { user, logout } = useAuth();

  const handleLogout = () => {
    logout();
  };

  return (
    <div className="min-h-screen bg-background p-8">
      <div className="max-w-2xl mx-auto">
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-3xl font-bold">Dashboard</h1>
          <Button onClick={handleLogout}>
            Sign out
          </Button>
        </div>
        
        <div className="space-y-4">
          <h2 className="text-xl">Welcome to OpenWaitlist!</h2>
          <p>Email: {user?.email}</p>
          <p>User ID: {user?.id}</p>
        </div>
      </div>
    </div>
  );
}