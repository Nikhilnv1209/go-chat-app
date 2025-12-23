'use client';

import { useEffect, useState } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import { useAppSelector } from '@/store/hooks';
import { MessageSquare, Menu, Settings } from 'lucide-react';
import ChatSidebar from '@/components/chat/ChatSidebar';
import NavigationRail from '@/components/dashboard/NavigationRail';
import { Button } from '@/components/ui/button';

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const { isAuthenticated, isLoading } = useAppSelector((state) => state.auth);
  const router = useRouter();
  const pathname = usePathname();
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);

  // Hide ChatSidebar on Profile page to give it full width
  const showChatSidebar = !pathname.includes('/profile');

  const toggleSidebar = () => {
    setIsSidebarOpen(!isSidebarOpen);
  };

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.push('/login');
    }
  }, [isLoading, isAuthenticated, router]);

  if (isLoading) {
    return (
      <div className="min-h-[100dvh] relative flex flex-col overflow-x-hidden bg-slate-950">
        <div className="fixed inset-0 h-[100lvh] w-full overflow-hidden pointer-events-none bg-slate-950">
          <div className="absolute inset-0 bg-[radial-gradient(ellipse_80%_80%_at_50%_-20%,rgba(120,119,198,0.3),rgba(255,255,255,0))]" />
        </div>
        <div className="relative z-10 flex items-center justify-center flex-grow min-h-[100dvh]">
          <div className="animate-spin rounded-full h-12 w-12 border-2 border-indigo-500 border-t-transparent"></div>
        </div>
      </div>
    );
  }

  if (!isAuthenticated) {
    return null; // Will redirect via useEffect
  }

  return (
    <div className="flex h-screen overflow-hidden bg-slate-950 ios-safe-area no-bounce">
      {/* Navigation Rail - Desktop Only */}
      <NavigationRail />

      {/* Main Layout Area */}
      <div className="flex-1 flex overflow-hidden relative">
        {/* Mobile Sidebar Overlay */}
        {isSidebarOpen && (
          <div
            className="md:hidden fixed inset-0 bg-black/50 z-20"
            onClick={() => setIsSidebarOpen(false)}
          />
        )}

        {/* Sidebar (Chat List) */}
        <div className={showChatSidebar ? "" : "md:hidden"}>
           <ChatSidebar isOpen={isSidebarOpen} onClose={() => setIsSidebarOpen(false)} />
        </div>

        {/* Mobile Header */}
        <div className="md:hidden fixed top-0 left-0 right-0 z-20 bg-slate-950/80 backdrop-blur-sm border-b border-white/10">
          <div className="flex items-center justify-between h-14 px-4">
            <Button
              variant="ghost"
              size="icon"
              onClick={toggleSidebar}
              className="text-slate-400 hover:text-white hover:bg-white/5"
            >
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

        {/* Main Content Area */}
        <main className="flex-1 flex flex-col overflow-hidden w-full relative z-0">
          {/* Mobile spacer */}
          <div className="md:hidden h-14 flex-shrink-0"></div>
          {children}
        </main>
      </div>
    </div>
  );
}
