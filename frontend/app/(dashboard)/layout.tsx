'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAppSelector } from '@/store/hooks';

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const { isAuthenticated, isLoading } = useAppSelector((state) => state.auth);
  const router = useRouter();

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.push('/login');
    }
  }, [isLoading, isAuthenticated, router]);

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
      </div>
    );
  }

  if (!isAuthenticated) {
    return null; // Will redirect via useEffect
  }

  // Basic layout structure for now - Sidebar will be added in Phase 3
  return (
    <div className="flex h-screen overflow-hidden bg-background">
      {/* Sidebar Placeholder */}
      <aside className="hidden md:flex w-80 flex-col border-r bg-muted/10">
        <div className="p-4 border-b h-14 flex items-center">
          <span className="font-semibold">Sidebar (Phase 3)</span>
        </div>
      </aside>

      {/* Main Content Area */}
      <main className="flex-1 flex flex-col overflow-hidden">
        {children}
      </main>
    </div>
  );
}
