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

  // Logic to show sidebar on mobile:
  // - If manually opened via hamburger (isSidebarOpen)
  // - OR if we are NOT in a chat sesssion (default view on mobile is list)
  // BUT: The current layout code uses `showChatSidebar` class logic which is messy.
  // Let's rely on standard logic:
  // Desktop: Always show sidebar (collapsed or expanded)
  // Mobile:
  //   - If isChatOpen -> Hide sidebar (unless manually toggled?) No, manually toggling in chat usually opens drawer.
  //   - If !isChatOpen -> Show sidebar (full width)

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
    // ... loading ...
    return <div className="min-h-screen bg-slate-950 flex items-center justify-center"><div className="animate-spin rounded-full h-8 w-8 border-2 border-indigo-500 border-t-transparent"></div></div>;
  }

  if (!isAuthenticated) return null;

  return (
    <div className="flex h-screen overflow-hidden bg-slate-950 ios-safe-area no-bounce">
      {/* Navigation Rail - Desktop Only */}
      <NavigationRail />

      {/* Main Layout Area */}
      <div className="flex-1 flex overflow-hidden relative">
        {/* Mobile Sidebar Overlay (Only when manually opened) */}
        {isSidebarOpen && (
          <div
            className="md:hidden fixed inset-0 bg-black/60 z-40 backdrop-blur-sm"
            onClick={() => setIsSidebarOpen(false)}
          />
        )}

        {/* Sidebar Structure */}
        <div
          className={cn(
            // Base layout
            "h-full border-r border-white/[0.05] bg-slate-950 relative z-30 transition-all duration-300 ease-in-out flex-shrink-0",

            // Mobile Behavior:
            // 1. If sidebar is manually OPEN, it slides in as a drawer (absolute/fixed behavior handled by inner or class)
            // 2. If we are on LIST view (!isChatOpen), it acts as the main page content (w-full).
            // 3. If we are on CHAT view (isChatOpen), it is hidden unless manually opened.

            // Actually, blending "Drawer" and "Page" logic is tricky.
            // Let's use absolute positioning for the Drawer mode on mobile.

            "md:static md:flex", // Desktop: static flex item

            isSidebarOpen
                ? "fixed inset-y-0 left-0 w-80 shadow-2xl translate-x-0"  // Open Drawer (Mobile)
                : (isChatOpen ? "hidden" : "flex w-full"), // Closed Drawer: Hidden if chatting, Full if list view

             // Desktop Widths
             "md:translate-x-0", // Always visible on desktop (reset transform)
             !isSidebarCollapsed ? "md:w-80 lg:w-96" : "md:w-[0px] md:border-r-0 md:overflow-hidden" // Collapse logic
          )}
        >
             <ChatSidebar isOpen={isSidebarOpen} onClose={() => setIsSidebarOpen(false)} />
        </div>

        {/* Mobile Header (Visible only when in Chat List mode on Mobile) */}
        {!isChatOpen && (
            <div className="md:hidden fixed top-0 left-0 right-0 z-20 bg-slate-950/80 backdrop-blur-md border-b border-white/10">
            <div className="flex items-center justify-between h-14 px-4">
                <Button variant="ghost" size="icon" className="text-slate-400 opacity-0 cursor-default">
                    <Menu className="w-5 h-5" />
                </Button>
                <h1 className="text-lg font-semibold text-white flex items-center gap-2">
                <MessageSquare className="w-5 h-5 text-indigo-400" />
                Chat
                </h1>
                <Button
                variant="ghost"
                size="icon"
                onClick={() => router.push('/dashboard/profile')}
                className="text-slate-400 hover:text-white hover:bg-white/5"
                >
                <Settings className="w-5 h-5" />
                </Button>
            </div>
            </div>
        )}

        {/* Main Content Area */}
        <main className={cn(
            "flex-1 flex flex-col overflow-hidden w-full relative z-0 bg-slate-950",
            // On Mobile:
            // If Chat is Open -> Show
            // If List is Open -> Hide (because List takes full width)
            isChatOpen ? "flex" : "hidden md:flex"
        )}>
          {/* spacer for header if needed, but usually handled by pages */}
          {/* In list view, header is in sidebar (above). In chat view, chat page has header. */}
          {children}
        </main>
      </div>
    </div>
  );
}
