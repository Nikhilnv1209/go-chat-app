'use client';

import { useEffect, useState } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import { useAppSelector } from '@/store/hooks';
import { MessageSquare, Menu, Settings } from 'lucide-react';
import ChatSidebar from '@/components/chat/ChatSidebar';
import NavigationRail from '@/components/dashboard/NavigationRail';
import { Button } from '@/components/ui/button';

import { useSocketConnection } from '@/hooks/useSocketConnection';
import { cn } from '@/lib/utils';

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  useSocketConnection(); // Init WebSocket

  const { isAuthenticated, isLoading } = useAppSelector((state) => state.auth);
  const { isSidebarCollapsed } = useAppSelector((state) => state.ui);
  const router = useRouter();
  const pathname = usePathname();
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);

  // Determine if we are in a chat session
  const isChatOpen = pathname.includes('/dashboard/chat/');

  const toggleSidebar = () => {
    setIsSidebarOpen(!isSidebarOpen);
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
    <div className="flex h-screen overflow-hidden bg-[#d8d8ec] ios-safe-area no-bounce">
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
            "h-full border-r border-[#7678ed]/10 bg-white relative z-30 transition-all duration-300 ease-in-out flex-shrink-0 shadow-lg shadow-[#7678ed]/5",

            "md:static md:flex",

            isSidebarOpen
                ? "fixed inset-y-0 left-0 w-80 shadow-2xl translate-x-0"
                : (isChatOpen ? "hidden" : "flex w-full"),

             // Desktop Widths
             "md:translate-x-0",
             !isSidebarCollapsed ? "md:w-80 lg:w-96" : "md:w-[0px] md:border-r-0 md:overflow-hidden"
          )}
        >
             <ChatSidebar isOpen={isSidebarOpen} onClose={() => setIsSidebarOpen(false)} />
        </div>

        {/* Mobile Header (Visible only when in Chat List mode on Mobile) */}
        {!isChatOpen && (
            <div className="md:hidden fixed top-0 left-0 right-0 z-20 bg-white/90 backdrop-blur-md border-b border-[#7678ed]/10">
            <div className="flex items-center justify-between h-14 px-4">
                <Button variant="ghost" size="icon" className="text-[#202022]/50 opacity-0 cursor-default">
                    <Menu className="w-5 h-5" />
                </Button>
                <h1 className="text-lg font-semibold text-[#202022] flex items-center gap-2">
                <MessageSquare className="w-5 h-5 text-[#7678ed]" />
                Chat
                </h1>
                <Button
                variant="ghost"
                size="icon"
                onClick={() => router.push('/dashboard/profile')}
                className="text-[#202022]/50 hover:text-[#202022] hover:bg-[#7678ed]/10"
                >
                <Settings className="w-5 h-5" />
                </Button>
            </div>
            </div>
        )}

        {/* Main Content Area */}
        <main className={cn(
            "flex-1 flex flex-col overflow-hidden w-full relative z-0 bg-[#f9fafc]",
            isChatOpen ? "flex" : "hidden md:flex"
        )}>
          {children}
        </main>
      </div>
    </div>
  );
}
