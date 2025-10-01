import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Label } from '@/components/ui/label';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
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
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog';
import { LoadingSpinner } from '@/components/shared/LoadingSpinner';
import { ErrorMessage } from '@/components/shared/ErrorMessage';
import { Search, Plus, MoreHorizontal, ExternalLink, Eye, Pencil, Trash2 } from 'lucide-react';

interface Waitlist {
  id: number;
  slug: string;
  name: string;
  owner_user_id: number;
  is_public: boolean;
  show_vendor_branding: boolean;
  created_at: string;
}

interface WaitlistsResponse {
  waitlists: Waitlist[];
  total: number;
}

interface CreateWaitlistForm {
  name: string;
  is_public: boolean;
  show_vendor_branding: boolean;
}

export function WaitlistTable() {
  const navigate = useNavigate();
  const [waitlists, setWaitlists] = useState<Waitlist[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [searchInput, setSearchInput] = useState('');
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [isCreating, setIsCreating] = useState(false);
  const [formData, setFormData] = useState<CreateWaitlistForm>({
    name: '',
    is_public: true,
    show_vendor_branding: false,
  });
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const [waitlistToDelete, setWaitlistToDelete] = useState<{slug: string, name: string} | null>(null);

  // Debounce search
  useEffect(() => {
    const timer = setTimeout(() => {
      setSearchTerm(searchInput);
    }, 300);

    return () => clearTimeout(timer);
  }, [searchInput]);

  // Fetch waitlists
  useEffect(() => {
    fetchWaitlists();
  }, [searchTerm]);

  const fetchWaitlists = async () => {
    try {
      setLoading(true);
      setError(null);
      
      const queryParams = new URLSearchParams();
      if (searchTerm) {
        queryParams.append('search', searchTerm);
      }
      
      const response = await fetch(`/api/v1/waitlists?${queryParams.toString()}`, {
        credentials: 'include',
      });

      if (!response.ok) {
        throw new Error(`Failed to fetch waitlists: ${response.statusText}`);
      }

      const data: WaitlistsResponse = await response.json();
      setWaitlists(data.waitlists || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An unknown error occurred');
    } finally {
      setLoading(false);
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    });
  };

  const handleCreateWaitlist = async () => {
    if (!formData.name.trim()) return;

    try {
      setIsCreating(true);
      const response = await fetch('/api/v1/waitlists', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify(formData),
      });

      if (!response.ok) {
        throw new Error(`Failed to create waitlist: ${response.statusText}`);
      }

      // Reset form and close dialog
      setFormData({
        name: '',
        is_public: true,
        show_vendor_branding: false,
      });
      setIsCreateDialogOpen(false);
      
      // Refresh the waitlists
      await fetchWaitlists();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create waitlist');
    } finally {
      setIsCreating(false);
    }
  };

  const handleViewWaitlist = (slug: string) => {
    navigate(`/dashboard/waitlists/${slug}`);
  };

  const handleOpenInNewTab = (slug: string) => {
    window.open(`/dashboard/waitlists/${slug}`, '_blank');
  };

  const handleDeleteWaitlist = (slug: string, name: string) => {
    setWaitlistToDelete({ slug, name });
    setIsDeleteDialogOpen(true);
  };

  const confirmDeleteWaitlist = async () => {
    if (!waitlistToDelete) return;

    try {
      const response = await fetch(`/api/v1/waitlists/${waitlistToDelete.slug}`, {
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

      // Refresh the waitlists after successful deletion
      fetchWaitlists();
      setIsDeleteDialogOpen(false);
      setWaitlistToDelete(null);
    } catch (error) {
      console.error('Delete error:', error);
      alert(error instanceof Error ? error.message : 'Failed to delete waitlist');
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center py-8">
        <LoadingSpinner />
      </div>
    );
  }

  if (error) {
    return (
      <div className="py-8">
        <ErrorMessage message={error} />
        <Button onClick={fetchWaitlists} className="mt-4">
          Retry
        </Button>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* Header Section */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-xl font-semibold">Your Waitlists</h2>
          <p className="text-sm text-muted-foreground">
            Manage your waitlists and track signups
          </p>
        </div>
        <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
          <DialogTrigger asChild>
            <Button className="gap-2">
              <Plus className="h-4 w-4" />
              Create Waitlist
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-[425px]">
            <DialogHeader>
              <DialogTitle>Create New Waitlist</DialogTitle>
              <DialogDescription>
                Create a new waitlist to start collecting signups from your audience.
              </DialogDescription>
            </DialogHeader>
            <div className="grid gap-4 py-4">
              <div className="grid gap-2">
                <Label htmlFor="name">Waitlist Name</Label>
                <Input
                  id="name"
                  placeholder="Enter waitlist name..."
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                />
              </div>
              <div className="flex items-center space-x-2">
                <Checkbox
                  id="is_public"
                  checked={formData.is_public}
                  onCheckedChange={(checked) => setFormData({ ...formData, is_public: !!checked })}
                />
                <Label htmlFor="is_public">Make this waitlist public</Label>
              </div>
              <div className="flex items-center space-x-2">
                <Checkbox
                  id="show_vendor_branding"
                  checked={formData.show_vendor_branding}
                  onCheckedChange={(checked) => setFormData({ ...formData, show_vendor_branding: !!checked })}
                />
                <Label htmlFor="show_vendor_branding">Show OpenWaitlist branding</Label>
              </div>
            </div>
            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => setIsCreateDialogOpen(false)}
              >
                Cancel
              </Button>
              <Button
                type="button"
                onClick={handleCreateWaitlist}
                disabled={!formData.name.trim() || isCreating}
              >
                {isCreating ? 'Creating...' : 'Create Waitlist'}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>

      {/* Search Section */}
      <div className="relative max-w-sm">
        <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
        <Input
          placeholder="Search waitlists..."
          value={searchInput}
          onChange={(e) => setSearchInput(e.target.value)}
          className="pl-9"
        />
      </div>

      {/* Table Section */}
      <div className="rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Name</TableHead>
              <TableHead>Created</TableHead>
              <TableHead className="text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {waitlists.length === 0 ? (
              <TableRow>
                <TableCell colSpan={3} className="text-center py-8">
                  <div className="flex flex-col items-center gap-2">
                    <div className="text-muted-foreground">
                      {searchTerm ? 'No waitlists found matching your search.' : 'No waitlists created yet.'}
                    </div>
                    {!searchTerm && (
                      <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
                        <DialogTrigger asChild>
                          <Button variant="outline" size="sm" className="gap-2">
                            <Plus className="h-4 w-4" />
                            Create your first waitlist
                          </Button>
                        </DialogTrigger>
                        <DialogContent className="sm:max-w-[425px]">
                          <DialogHeader>
                            <DialogTitle>Create New Waitlist</DialogTitle>
                            <DialogDescription>
                              Create a new waitlist to start collecting signups from your audience.
                            </DialogDescription>
                          </DialogHeader>
                          <div className="grid gap-4 py-4">
                            <div className="grid gap-2">
                              <Label htmlFor="name">Waitlist Name</Label>
                              <Input
                                id="name"
                                placeholder="Enter waitlist name..."
                                value={formData.name}
                                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                              />
                            </div>
                            <div className="flex items-center space-x-2">
                              <Checkbox
                                id="is_public_empty"
                                checked={formData.is_public}
                                onCheckedChange={(checked) => setFormData({ ...formData, is_public: !!checked })}
                              />
                              <Label htmlFor="is_public_empty">Make this waitlist public</Label>
                            </div>
                            <div className="flex items-center space-x-2">
                              <Checkbox
                                id="show_vendor_branding_empty"
                                checked={formData.show_vendor_branding}
                                onCheckedChange={(checked) => setFormData({ ...formData, show_vendor_branding: !!checked })}
                              />
                              <Label htmlFor="show_vendor_branding_empty">Show OpenWaitlist branding</Label>
                            </div>
                          </div>
                          <DialogFooter>
                            <Button
                              type="button"
                              variant="outline"
                              onClick={() => setIsCreateDialogOpen(false)}
                            >
                              Cancel
                            </Button>
                            <Button
                              type="button"
                              onClick={handleCreateWaitlist}
                              disabled={!formData.name.trim() || isCreating}
                            >
                              {isCreating ? 'Creating...' : 'Create Waitlist'}
                            </Button>
                          </DialogFooter>
                        </DialogContent>
                      </Dialog>
                    )}
                  </div>
                </TableCell>
              </TableRow>
            ) : (
              waitlists.map((waitlist) => (
                <TableRow 
                  key={waitlist.id}
                  className="cursor-pointer hover:bg-muted/50"
                  onClick={() => handleViewWaitlist(waitlist.slug)}
                >
                  <TableCell className="font-medium">{waitlist.name}</TableCell>
                  <TableCell className="text-sm text-muted-foreground">
                    {formatDate(waitlist.created_at)}
                  </TableCell>
                  <TableCell className="text-right">
                    <DropdownMenu>
                      <DropdownMenuTrigger asChild>
                        <Button 
                          variant="ghost" 
                          className="h-8 w-8 p-0"
                          onClick={(e) => e.stopPropagation()}
                        >
                          <span className="sr-only">Open menu</span>
                          <MoreHorizontal className="h-4 w-4" />
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end">
                        <DropdownMenuItem 
                          onClick={(e) => {
                            e.stopPropagation();
                            handleViewWaitlist(waitlist.slug);
                          }}
                        >
                          <Eye className="mr-2 h-4 w-4" />
                          View
                        </DropdownMenuItem>
                        <DropdownMenuItem 
                          onClick={(e) => {
                            e.stopPropagation();
                            handleOpenInNewTab(waitlist.slug);
                          }}
                        >
                          <ExternalLink className="mr-2 h-4 w-4" />
                          Open in new tab
                        </DropdownMenuItem>
                        <DropdownMenuItem onClick={(e) => e.stopPropagation()}>
                          <Pencil className="mr-2 h-4 w-4" />
                          Edit
                        </DropdownMenuItem>
                        <DropdownMenuItem 
                          className="text-red-600"
                          onClick={(e) => {
                            e.stopPropagation();
                            handleDeleteWaitlist(waitlist.slug, waitlist.name);
                          }}
                        >
                          <Trash2 className="mr-2 h-4 w-4" />
                          Delete
                        </DropdownMenuItem>
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>

      {/* Footer Stats */}
      {waitlists.length > 0 && (
        <div className="text-sm text-muted-foreground">
          Showing {waitlists.length} waitlist{waitlists.length !== 1 ? 's' : ''}
          {searchTerm && ` matching "${searchTerm}"`}
        </div>
      )}

      {/* Delete Confirmation Dialog */}
      <AlertDialog open={isDeleteDialogOpen} onOpenChange={setIsDeleteDialogOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete Waitlist</AlertDialogTitle>
            <AlertDialogDescription>
              Are you sure you want to delete "{waitlistToDelete?.name}"? This action cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction onClick={confirmDeleteWaitlist} className="bg-red-600 hover:bg-red-700">
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}