import React from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { AppSidebar } from '@/components/app-sidebar';
import { SidebarInset, SidebarProvider, SidebarTrigger } from '@/components/ui/sidebar';
import { Separator } from '@/components/ui/separator';
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb';

export function DashboardPage() {
  const { user } = useAuth();

  return (
    <SidebarProvider>
      <AppSidebar />
      <SidebarInset className="flex-1">
        <header className="flex h-16 shrink-0 items-center gap-2 border-b transition-[width,height] ease-linear group-has-[[data-collapsible=icon]]/sidebar-wrapper:h-12">
          <div className="flex items-center gap-2 px-6 w-full">
            <SidebarTrigger className="-ml-1" />
            <Separator orientation="vertical" className="mr-2 h-4" />
            <Breadcrumb>
              <BreadcrumbList>
                <BreadcrumbItem className="hidden md:block">
                  <BreadcrumbLink href="/dashboard">
                    Dashboard
                  </BreadcrumbLink>
                </BreadcrumbItem>
                <BreadcrumbSeparator className="hidden md:block" />
                <BreadcrumbItem>
                  <BreadcrumbPage>Overview</BreadcrumbPage>
                </BreadcrumbItem>
              </BreadcrumbList>
            </Breadcrumb>
          </div>
        </header>
        <main className="flex-1 overflow-auto">
          <div className="h-full px-6 py-6">
            {/* User Info Header */}
            <div className="mb-6">
              <h1 className="text-2xl font-semibold mb-2">Welcome to OpenWaitlist!</h1>
              <div className="flex items-center gap-6 text-sm text-muted-foreground">
                <span><strong>Email:</strong> {user?.email}</span>
                <span><strong>ID:</strong> {user?.id}</span>
                <span><strong>Display Name:</strong> {user?.display_name || 'Not set'}</span>
              </div>
            </div>

            {/* Main Content Area */}
            <div className="rounded-xl border bg-card text-card-foreground shadow">
              
            </div>
          </div>
        </main>
      </SidebarInset>
    </SidebarProvider>
  );
}