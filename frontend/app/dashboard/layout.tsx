'use client';

import { useEffect, useState } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import { useAppSelector } from '@/store/hooks';
import { MessageSquare, MoreVertical } from 'lucide-react';
import ChatSidebar from '@/components/chat/ChatSidebar';
import NavigationRail from '@/components/dashboard/NavigationRail';
import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';

import { useSocketConnection } from '@/hooks/useSocketConnection';
import { cn } from '@/lib/utils';
import { useAppDispatch } from '@/store/hooks';
import { logout } from '@/store/features/authSlice';

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  useSocketConnection(); // Init WebSocket

  const dispatch = useAppDispatch();
  const { isAuthenticated, isLoading, user } = useAppSelector((state) => state.auth);
  const { isSidebarCollapsed } = useAppSelector((state) => state.ui);
  const router = useRouter();
  const pathname = usePathname();
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);

  // Determine if we are in a chat session or profile page
  const isChatOpen = pathname.includes('/dashboard/chat/');
  const isProfilePage = pathname === '/dashboard/profile';
  const showMainContent = isChatOpen || isProfilePage;

  const toggleSidebar = () => {
    setIsSidebarOpen(!isSidebarOpen);
  };

  const handleLogout = () => {
    dispatch(logout());
    router.push('/login');
  };

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.push('/login');
    }
    // Close sidebar when navigating (especially on mobile)
    setIsSidebarOpen(false);
  }, [isLoading, isAuthenticated, router, pathname]);

  if (isLoading) {
    return (
      <div className="min-h-screen bg-[#d8d8ec] flex items-center justify-center">
        <div className="animate-spin rounded-full h-8 w-8 border-2 border-[#7678ed] border-t-transparent"></div>
      </div>
    );
  }

  if (!isAuthenticated) return null;

  return (
    <div className="flex h-screen overflow-hidden bg-[#d8d8ec]">
      {/* Navigation Rail - Desktop Only */}
      <NavigationRail />

      {/* Main Layout Area */}
      <div className="flex-1 flex overflow-hidden relative">
        {/* Mobile Sidebar Overlay (Only when manually opened) */}
        {isSidebarOpen && (
          <div
            className="md:hidden fixed inset-0 bg-[#202022]/40 z-40 backdrop-blur-sm"
            onClick={() => setIsSidebarOpen(false)}
          />
        )}

        {/* Sidebar Structure */}
        <div
          className={cn(
            // Base layout
            "h-full md:border-r md:border-[#7678ed]/10 md:bg-white relative z-30 transition-all duration-300 ease-in-out flex-shrink-0 shadow-lg shadow-[#7678ed]/5",

            "md:static md:flex",

            isSidebarOpen
                ? "fixed inset-y-0 left-0 w-80 shadow-2xl translate-x-0"
                : (isChatOpen || isProfilePage ? "hidden" : "flex flex-col w-full overflow-hidden md:overflow-visible"),

             // Desktop Widths
             "md:translate-x-0",
             !isSidebarCollapsed ? "md:w-80 lg:w-96" : "md:w-[0px] md:border-r-0 md:overflow-hidden"
          )}
        >
             <ChatSidebar isOpen={isSidebarOpen} onClose={() => setIsSidebarOpen(false)} />
        </div>

        {/* Main Content Area */}
        <main className={cn(
            "flex-1 flex flex-col overflow-hidden w-full relative z-0 bg-[#f9fafc]",
            showMainContent ? "flex" : "hidden md:flex"
        )}>
          {children}
        </main>
      </div>
    </div>
  );
}
