'use client';

import { useEffect, useState } from 'react';
import { useAppSelector } from '@/store/hooks';
import MobileChatSidebar from './MobileChatSidebar';
import DesktopChatSidebar from './DesktopChatSidebar';
import { cn } from '@/lib/utils';

interface ChatSidebarProps {
  isOpen?: boolean;
  onClose?: () => void;
}

export default function ChatSidebar({ isOpen, onClose }: ChatSidebarProps = {}) {
  const { isSidebarCollapsed } = useAppSelector((state) => state.ui);

  return (
    <aside className={cn(
        "flex flex-col w-full h-full transition-all duration-300 ease-in-out relative",
        "overflow-hidden"
    )}>
      {/* Mobile Sidebar */}
      <div className="md:hidden flex w-full h-full">
        <MobileChatSidebar onClose={onClose} />
      </div>

      {/* Desktop Sidebar */}
      <div className="hidden md:flex w-full h-full">
        <DesktopChatSidebar />
      </div>
    </aside>
  );
}
