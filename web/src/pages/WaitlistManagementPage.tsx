import React, { useState, useEffect } from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import { Settings, Users } from 'lucide-react';
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

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Checkbox } from '@/components/ui/checkbox';
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog';
import { LoadingSpinner } from '@/components/shared/LoadingSpinner';
import { ErrorMessage } from '@/components/shared/ErrorMessage';

interface Waitlist {
  id: number;
  name: string;
  slug: string;
  is_public: boolean;
  show_vendor_branding: boolean;
  created_at: string;
}

interface Signup {
  id: number;
  email: string;
  first_name?: string;
  last_name?: string;
  created_at: string;
  position: number;
}

export function WaitlistManagementPage() {
  const { slug } = useParams<{ slug: string }>();
  const navigate = useNavigate();
  const [waitlist, setWaitlist] = useState<Waitlist | null>(null);
  const [signups, setSignups] = useState<Signup[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState('');
  const [activeTab, setActiveTab] = useState<'signups' | 'settings'>('signups');
  const [activeSignupTab, setActiveSignupTab] = useState<'all' | 'offboarded' | 'import-export'>('all');
  const [isUpdatingSettings, setIsUpdatingSettings] = useState(false);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const [settingsFormData, setSettingsFormData] = useState({
    name: '',
    is_public: false,
    show_vendor_branding: false,
  });

  useEffect(() => {
    if (!slug) return;
    fetchWaitlist();
  }, [slug]);

  const fetchWaitlist = async () => {
    if (!slug) return;
    
    try {
      setLoading(true);
      setError(null);
      
      const response = await fetch(`/api/v1/waitlists/${slug}`, {
        credentials: 'include',
      });
      
      if (!response.ok) {
        if (response.status === 404) {
          throw new Error('Waitlist not found');
        }
        throw new Error('Failed to fetch waitlist');
      }
      
      const data = await response.json();
      setWaitlist(data);
      
      // Populate settings form data
      setSettingsFormData({
        name: data.name,
        is_public: data.is_public,
        show_vendor_branding: data.show_vendor_branding,
      });
      
      // TODO: Fetch signups for this waitlist
      // For now, using placeholder data
      setSignups([
        {
          id: 1,
          email: 'john@example.com',
          first_name: 'John',
          last_name: 'Doe',
          created_at: '2024-01-15T10:30:00Z',
          position: 1,
        },
        {
          id: 2,
          email: 'jane@example.com',
          first_name: 'Jane',
          last_name: 'Smith',
          created_at: '2024-01-16T14:20:00Z',
          position: 2,
        },
        {
          id: 3,
          email: 'bob@example.com',
          first_name: 'Bob',
          last_name: 'Johnson',
          created_at: '2024-01-17T09:15:00Z',
          position: 3,
        },
      ]);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred');
    } finally {
      setLoading(false);
    }
  };

  const filteredSignups = signups.filter(signup =>
    signup.email.toLowerCase().includes(searchQuery.toLowerCase()) ||
    `${signup.first_name || ''} ${signup.last_name || ''}`.toLowerCase().includes(searchQuery.toLowerCase())
  );

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const handleUpdateSettings = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!waitlist) return;

    setIsUpdatingSettings(true);
    try {
      const response = await fetch(`/api/v1/waitlists/${waitlist.slug}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify(settingsFormData),
      });

      if (!response.ok) {
        throw new Error('Failed to update waitlist settings');
      }

      const updatedWaitlist = await response.json();
      setWaitlist(updatedWaitlist);
      
      // Show success message or toast (for now just log)
      console.log('Settings updated successfully');
    } catch (error) {
      console.error('Settings update error:', error);
      alert(error instanceof Error ? error.message : 'Failed to update settings');
    } finally {
      setIsUpdatingSettings(false);
    }
  };

  const handleDeleteWaitlist = async () => {
    if (!waitlist) return;

    try {
      const response = await fetch(`/api/v1/waitlists/${waitlist.slug}`, {
        method: 'DELETE',
        credentials: 'include',
      });

      if (!response.ok) {
        if (response.status === 404) {
          throw new Error('Waitlist not found');
        } else if (response.status === 403) {
          throw new Error('You do not have permission to delete this waitlist');
        }
        throw new Error('Failed to delete waitlist');
      }

      // Navigate back to dashboard after successful deletion
      navigate('/dashboard');
    } catch (error) {
      console.error('Delete error:', error);
      alert(error instanceof Error ? error.message : 'Failed to delete waitlist');
    } finally {
      setIsDeleteDialogOpen(false);
    }
  };

  if (loading) {
    return (
      <SidebarProvider>
        <AppSidebar />
        <SidebarInset className="flex-1">
          <div className="flex items-center justify-center h-full">
            <LoadingSpinner />
          </div>
        </SidebarInset>
      </SidebarProvider>
    );
  }

  if (error || !waitlist) {
    return (
      <SidebarProvider>
        <AppSidebar />
        <SidebarInset className="flex-1">
          <div className="flex items-center justify-center h-full">
            <ErrorMessage message={error || 'Waitlist not found'} />
          </div>
        </SidebarInset>
      </SidebarProvider>
    );
  }

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
                  <BreadcrumbLink asChild>
                    <Link to="/dashboard">Dashboard</Link>
                  </BreadcrumbLink>
                </BreadcrumbItem>
                <BreadcrumbSeparator className="hidden md:block" />
                <BreadcrumbItem className="hidden md:block">
                  <BreadcrumbLink asChild>
                    <Link to="/dashboard">Waitlists</Link>
                  </BreadcrumbLink>
                </BreadcrumbItem>
                <BreadcrumbSeparator className="hidden md:block" />
                <BreadcrumbItem>
                  <BreadcrumbPage>{waitlist.name}</BreadcrumbPage>
                </BreadcrumbItem>
              </BreadcrumbList>
            </Breadcrumb>
          </div>
        </header>
        <main className="flex-1 overflow-auto">
          <div className="h-full px-6 py-6">
            {/* Header */}
            <div className="mb-6">
              <div>
                <h1 className="text-3xl font-bold tracking-tight">{waitlist.name}</h1>
                <p className="text-muted-foreground">
                  Created on {formatDate(waitlist.created_at)}
                </p>
              </div>
            </div>

            {/* Main Content with Custom Tabs */}
            <div className="rounded-xl border bg-card text-card-foreground shadow">
              {/* Main Tab Navigation */}
              <div className="border-b px-6">
                <div className="flex space-x-8">
                  <button
                    onClick={() => setActiveTab('signups')}
                    className={`flex items-center gap-2 text-sm px-0 py-4 border-b-2 transition-colors ${
                      activeTab === 'signups'
                        ? 'text-primary border-primary'
                        : 'text-muted-foreground border-transparent hover:text-foreground'
                    }`}
                  >
                    <Users className="w-4 h-4" />
                    Signups
                  </button>
                  <button
                    onClick={() => setActiveTab('settings')}
                    className={`flex items-center gap-2 text-sm px-0 py-4 border-b-2 transition-colors ${
                      activeTab === 'settings'
                        ? 'text-primary border-primary'
                        : 'text-muted-foreground border-transparent hover:text-foreground'
                    }`}
                  >
                    <Settings className="w-4 h-4" />
                    Settings
                  </button>
                </div>
              </div>
              
              {/* Tab Content */}
              <div className="p-6">
                {activeTab === 'signups' && (
                  <div className="space-y-6">
                    {/* Secondary Tab Navigation for Signups */}
                    <div className="flex items-center justify-between">
                      <div className="flex space-x-1 bg-muted/50 rounded-md p-1">
                        <button
                          onClick={() => setActiveSignupTab('all')}
                          className={`px-3 py-1.5 text-sm rounded-sm transition-colors ${
                            activeSignupTab === 'all'
                              ? 'bg-background text-foreground shadow-sm'
                              : 'text-muted-foreground hover:text-foreground'
                          }`}
                        >
                          All Signups
                        </button>
                        <button
                          onClick={() => setActiveSignupTab('offboarded')}
                          className={`px-3 py-1.5 text-sm rounded-sm transition-colors ${
                            activeSignupTab === 'offboarded'
                              ? 'bg-background text-foreground shadow-sm'
                              : 'text-muted-foreground hover:text-foreground'
                          }`}
                        >
                          Offboarded
                        </button>
                        <button
                          onClick={() => setActiveSignupTab('import-export')}
                          className={`px-3 py-1.5 text-sm rounded-sm transition-colors ${
                            activeSignupTab === 'import-export'
                              ? 'bg-background text-foreground shadow-sm'
                              : 'text-muted-foreground hover:text-foreground'
                          }`}
                        >
                          Import/Export
                        </button>
                      </div>
                      
                      {activeSignupTab === 'all' && (
                        <Input
                          placeholder="Search signups..."
                          value={searchQuery}
                          onChange={(e) => setSearchQuery(e.target.value)}
                          className="w-64"
                        />
                      )}
                    </div>

                    {/* Signup Tab Content */}
                    {activeSignupTab === 'all' && (
                      <div className="rounded-md border">
                        <Table>
                          <TableHeader>
                            <TableRow>
                              <TableHead className="w-16">Position</TableHead>
                              <TableHead>Name</TableHead>
                              <TableHead>Email</TableHead>
                              <TableHead>Joined</TableHead>
                            </TableRow>
                          </TableHeader>
                          <TableBody>
                            {filteredSignups.length === 0 ? (
                              <TableRow>
                                <TableCell colSpan={4} className="text-center py-8">
                                  {searchQuery ? 'No signups match your search.' : 'No signups yet.'}
                                </TableCell>
                              </TableRow>
                            ) : (
                              filteredSignups.map((signup) => (
                                <TableRow key={signup.id}>
                                  <TableCell className="font-mono text-sm">
                                    #{signup.position}
                                  </TableCell>
                                  <TableCell>
                                    {signup.first_name && signup.last_name
                                      ? `${signup.first_name} ${signup.last_name}`
                                      : '-'}
                                  </TableCell>
                                  <TableCell>{signup.email}</TableCell>
                                  <TableCell className="text-sm text-muted-foreground">
                                    {formatDate(signup.created_at)}
                                  </TableCell>
                                </TableRow>
                              ))
                            )}
                          </TableBody>
                        </Table>
                      </div>
                    )}

                    {activeSignupTab === 'offboarded' && (
                      <div className="text-center py-12">
                        <Users className="mx-auto h-12 w-12 text-muted-foreground" />
                        <h4 className="mt-4 text-lg font-medium">No Offboarded Signups</h4>
                        <p className="text-sm text-muted-foreground">
                          Offboarded signups will appear here when available.
                        </p>
                      </div>
                    )}

                    {activeSignupTab === 'import-export' && (
                      <div className="text-center py-12">
                        <Settings className="mx-auto h-12 w-12 text-muted-foreground" />
                        <h4 className="mt-4 text-lg font-medium">Import/Export Tools</h4>
                        <p className="text-sm text-muted-foreground">
                          Import and export functionality will be available here.
                        </p>
                      </div>
                    )}
                  </div>
                )}
                
                {activeTab === 'settings' && (
                  <div className="space-y-6">
                    <div>
                      <h3 className="text-lg font-medium">Settings</h3>
                      <p className="text-sm text-muted-foreground">
                        Manage your waitlist settings and configuration
                      </p>
                    </div>
                    
                    <form onSubmit={handleUpdateSettings} className="space-y-6">
                      <Card>
                        <CardHeader>
                          <CardTitle>General Settings</CardTitle>
                          <CardDescription>
                            Basic configuration for your waitlist
                          </CardDescription>
                        </CardHeader>
                        <CardContent className="space-y-4">
                          <div className="space-y-2">
                            <Label htmlFor="waitlist-name">Waitlist Name</Label>
                            <Input
                              id="waitlist-name"
                              value={settingsFormData.name}
                              onChange={(e) => setSettingsFormData(prev => ({ ...prev, name: e.target.value }))}
                              placeholder="Enter waitlist name"
                              disabled={isUpdatingSettings}
                            />
                          </div>
                          
                          <div className="flex items-center space-x-2">
                            <Checkbox
                              id="is-public"
                              checked={settingsFormData.is_public}
                              onCheckedChange={(checked) => 
                                setSettingsFormData(prev => ({ ...prev, is_public: checked === true }))
                              }
                              disabled={isUpdatingSettings}
                            />
                            <Label htmlFor="is-public" className="text-sm font-normal">
                              Make this waitlist public
                            </Label>
                          </div>
                          
                          <div className="flex items-center space-x-2">
                            <Checkbox
                              id="show-branding"
                              checked={settingsFormData.show_vendor_branding}
                              onCheckedChange={(checked) => 
                                setSettingsFormData(prev => ({ ...prev, show_vendor_branding: checked === true }))
                              }
                              disabled={isUpdatingSettings}
                            />
                            <Label htmlFor="show-branding" className="text-sm font-normal">
                              Show OpenWaitlist branding
                            </Label>
                          </div>
                        </CardContent>
                      </Card>
                      
                      <Card className="border-destructive/50">
                        <CardHeader>
                          <CardTitle className="text-destructive">Danger Zone</CardTitle>
                          <CardDescription>
                            Irreversible actions that will permanently affect your waitlist
                          </CardDescription>
                        </CardHeader>
                        <CardContent>
                          <div className="flex items-center justify-between">
                            <div>
                              <h4 className="font-medium">Delete Waitlist</h4>
                              <p className="text-sm text-muted-foreground">
                                Permanently delete this waitlist and all its data. This action cannot be undone.
                              </p>
                            </div>
                            <Button
                              variant="destructive"
                              onClick={() => setIsDeleteDialogOpen(true)}
                              disabled={isUpdatingSettings}
                            >
                              Delete Waitlist
                            </Button>
                          </div>
                        </CardContent>
                      </Card>
                      
                      <div className="flex justify-end">
                        <Button type="submit" disabled={isUpdatingSettings}>
                          {isUpdatingSettings ? 'Updating...' : 'Update Settings'}
                        </Button>
                      </div>
                    </form>
                  </div>
                )}
              </div>
            </div>
          </div>
        </main>
      </SidebarInset>
      
      {/* Delete Confirmation Dialog */}
      <AlertDialog open={isDeleteDialogOpen} onOpenChange={setIsDeleteDialogOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete Waitlist</AlertDialogTitle>
            <AlertDialogDescription>
              Are you sure you want to delete "{waitlist?.name}"? This action cannot be undone and will permanently delete all signups and data associated with this waitlist.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction onClick={handleDeleteWaitlist} className="bg-red-600 hover:bg-red-700">
              Delete Waitlist
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </SidebarProvider>
  );
}