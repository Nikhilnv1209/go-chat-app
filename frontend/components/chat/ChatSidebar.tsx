'use client';

import { useEffect, useState } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import { useQuery } from '@tanstack/react-query';
import { Search, MessageSquare, Users, LogOut, Settings, X, Briefcase, Archive, UserCheck, Star } from 'lucide-react';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import { setConversations, setActiveConversation, resetUnread } from '@/store/features/conversationSlice';
import { toggleFolderAssignment } from '@/store/features/folderSlice';
import { conversationApi } from '@/lib/conversationApi';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { ContextMenu, ContextMenuContent, ContextMenuItem, ContextMenuTrigger, ContextMenuSeparator, ContextMenuSub, ContextMenuSubTrigger, ContextMenuSubContent } from '@/components/ui/context-menu';
import { Conversation } from '@/types';
import { logout } from '@/store/features/authSlice';
import { cn } from '@/lib/utils';

interface ChatSidebarProps {
  isOpen?: boolean;
  onClose?: () => void;
}

export default function ChatSidebar({ isOpen, onClose }: ChatSidebarProps = {}) {
  const router = useRouter();
  const pathname = usePathname();
  const dispatch = useAppDispatch();
  const { token, user } = useAppSelector((state) => state.auth);
  const { conversations, activeConversationId } = useAppSelector((state) => state.conversation);
  const { assignments } = useAppSelector((state) => state.folders);
  const [searchQuery, setSearchQuery] = useState('');

  // Fetch conversations
  const { data, isLoading, isError } = useQuery({
    queryKey: ['conversations'],
    queryFn: () => conversationApi.getConversations(token!),
    enabled: !!token,
  });

  useEffect(() => {
    if (data) {
      dispatch(setConversations(data));
    }
  }, [data, dispatch]);

  // Determine active folder from URL
  let activeFolderId: string | null = null;
  if (pathname.includes('/work')) activeFolderId = 'work';
  else if (pathname.includes('/friends')) activeFolderId = 'friends';
  else if (pathname.includes('/archive')) activeFolderId = 'archive';

  const filteredConversations = conversations.filter((conv) => {
    const matchesSearch = conv.target_name.toLowerCase().includes(searchQuery.toLowerCase());
    if (!matchesSearch) return false;

    const isArchived = assignments['archive']?.includes(conv.id);

    // If viewing Archive, ONLY show archived chats
    if (activeFolderId === 'archive') {
      return isArchived;
    }

    // For all other views (All Chats, Work, Friends), HIDE archived chats
    if (isArchived) {
      return false;
    }

    if (activeFolderId) {
      return assignments[activeFolderId]?.includes(conv.id);
    }

    // Default "All Chats" view
    return true;
  });

  const handleConversationClick = (conv: Conversation) => {
    dispatch(setActiveConversation(conv.id));
    dispatch(resetUnread(conv.id));
    router.push(`/dashboard/chat/${conv.type.toLowerCase()}/${conv.target_id}`);
    if (onClose) {
      onClose();
    }
  };

  const handleLogout = () => {
    dispatch(logout());
    router.push('/login');
  };

  const handleToggleFolder = (folderId: string, conversationId: string) => {
    dispatch(toggleFolderAssignment({ folderId, conversationId }));
  };

  const formatTimestamp = (timestamp: string) => {
    const date = new Date(timestamp);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffMins = Math.floor(diffMs / 60000);
    const diffHours = Math.floor(diffMins / 60);
    const diffDays = Math.floor(diffHours / 24);

    if (diffMins < 1) return 'Just now';
    if (diffMins < 60) return `${diffMins}m`;
    if (diffHours < 24) return `${diffHours}h`;
    if (diffDays < 7) return `${diffDays}d`;
    return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
  };

  const truncateMessage = (message: string | null, maxLength: number = 35) => {
    if (!message) return 'No messages yet';
    return message.length > maxLength ? `${message.slice(0, maxLength)}...` : message;
  };

  const getFolderLabel = (id: string) => {
    switch(id) {
        case 'work': return 'Work';
        case 'friends': return 'Friends';
        case 'archive': return 'Archive';
        default: return 'Folder';
    }
  };

  return (
    <aside className={cn(
        "flex flex-col w-full h-full transition-transform duration-300 ease-in-out",
        "bg-white"
    )}>
      {/* Header */}
      <div className="flex-shrink-0 p-4 border-b border-[#7678ed]/10">
        {/* Mobile Header Controls */}
        <div className="md:hidden flex items-center justify-end mb-4">
          <div className="flex items-center gap-2">
            {/* Mobile Close Button */}
            <Button
              variant="ghost"
              size="icon"
              onClick={onClose}
              className="h-9 w-9 text-[#202022]/50 hover:text-[#202022] hover:bg-[#7678ed]/10"
            >
              <X className="w-4 h-4" />
            </Button>
            <Button
              variant="ghost"
              size="icon"
              onClick={() => router.push('/dashboard/profile')}
              className="h-9 w-9 text-[#202022]/50 hover:text-[#202022] hover:bg-[#7678ed]/10"
            >
              <Settings className="w-4 h-4" />
            </Button>
            <Button
              variant="ghost"
              size="icon"
              onClick={handleLogout}
              className="h-9 w-9 text-[#202022]/50 hover:text-[#ff7a55] hover:bg-[#ff7a55]/10"
            >
              <LogOut className="w-4 h-4" />
            </Button>
          </div>
        </div>

        {/* Search */}
        <div className="relative">
          <Search className="absolute left-3.5 top-1/2 -translate-y-1/2 h-[18px] w-[18px] text-[#7678ed]/50" />
          <Input
            type="text"
            placeholder="Search"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="pl-10 pr-4 h-11 bg-[#e8e8f5] border-0 text-[#202022] text-[15px] placeholder:text-[#202022]/40 focus:bg-[#dcdcf0] focus:ring-0 focus:ring-offset-0 focus-visible:ring-0 rounded-[12px] transition-colors"
          />
        </div>
      </div>

      {/* User Profile Mini - Mobile Only */}
      {user && (
        <div className="flex-shrink-0 p-4 border-b border-[#7678ed]/10 bg-gradient-to-r from-[#7678ed]/10 via-[#7678ed]/5 to-[#ff7a55]/5 md:hidden">
          <div className="flex items-center gap-3">
            <div className="relative">
              <div className="w-10 h-10 rounded-full bg-gradient-to-br from-[#7678ed] to-[#5a5cd9] flex items-center justify-center text-white font-semibold">
                {user.username.charAt(0).toUpperCase()}
              </div>
              <div className="absolute bottom-0 right-0 w-3 h-3 bg-green-500 border-2 border-white rounded-full"></div>
            </div>
            <div className="flex-1 min-w-0">
              <p className="text-sm font-semibold text-[#202022] truncate">{user.username}</p>
              <p className="text-xs text-[#202022]/50 truncate">{user.email}</p>
            </div>
          </div>
        </div>
      )}

      {/* Conversation List */}
      <div className="flex-1 overflow-y-auto">
        {isLoading && (
          <div className="flex items-center justify-center py-12">
            <div className="animate-spin rounded-full h-8 w-8 border-2 border-[#7678ed] border-t-transparent"></div>
          </div>
        )}

        {isError && (
          <div className="p-4 text-center">
            <p className="text-sm text-[#ff7a55]">Failed to load conversations</p>
          </div>
        )}

        {!isLoading && !isError && filteredConversations.length === 0 && (
          <div className="flex flex-col items-center justify-center py-12 px-4">
            <div className="w-16 h-16 rounded-full bg-[#7678ed]/10 flex items-center justify-center mb-4">
              {activeFolderId ? <Briefcase className="w-8 h-8 text-[#7678ed]" /> : <MessageSquare className="w-8 h-8 text-[#7678ed]" />}
            </div>
            <p className="text-sm text-[#202022]/60 text-center">
              {searchQuery
                ? 'No conversations found'
                : (activeFolderId ? 'This folder is empty' : 'No conversations yet')}
            </p>
            {!activeFolderId && (
                <p className="text-xs text-[#202022]/40 text-center mt-1">
                Start a new chat to get started
                </p>
            )}
            {activeFolderId && (
                <p className="text-xs text-[#202022]/40 text-center mt-1">
                Right click on a chat in "All Chats" to add it here
                </p>
            )}
          </div>
        )}

        {!isLoading && !isError && filteredConversations.map((conv) => (
          <ContextMenu key={conv.id}>
            <ContextMenuTrigger>
                <button
                    onClick={() => handleConversationClick(conv)}
                    className={`w-full p-4 flex items-start gap-3 hover:bg-[#7678ed]/5 transition-colors border-b border-[#7678ed]/5 ${
                    activeConversationId === conv.id ? 'bg-[#7678ed]/10' : ''
                    }`}
                >
                    {/* Avatar */}
                    <div className="relative flex-shrink-0">
                    {conv.type === 'GROUP' ? (
                        <div className="w-12 h-12 rounded-full bg-gradient-to-br from-[#ff7a55] to-[#e66a47] flex items-center justify-center text-white">
                        <Users className="w-6 h-6" />
                        </div>
                    ) : (
                        <div className="w-12 h-12 rounded-full bg-gradient-to-br from-[#7678ed] to-[#5a5cd9] flex items-center justify-center text-white font-semibold text-sm">
                        {conv.target_name.charAt(0).toUpperCase()}
                        </div>
                    )}
                    {conv.is_online && conv.type === 'DM' && (
                        <div className="absolute bottom-0 right-0 w-3.5 h-3.5 bg-green-500 border-2 border-white rounded-full"></div>
                    )}
                    {conv.unread_count > 0 && (
                        <div className="absolute -top-1 -right-1 min-w-[20px] h-5 bg-[#ff7a55] rounded-full flex items-center justify-center px-1.5">
                        <span className="text-xs font-bold text-white">{conv.unread_count > 99 ? '99+' : conv.unread_count}</span>
                        </div>
                    )}
                    </div>

                    {/* Content */}
                    <div className="flex-1 min-w-0 text-left">
                    <div className="flex items-center justify-between mb-1">
                        <h3 className="text-sm font-semibold text-[#202022] truncate">{conv.target_name}</h3>
                        <div className="flex items-center gap-1.5">
                          <span className="text-xs text-[#202022]/40 flex-shrink-0">
                          {formatTimestamp(conv.last_message_at)}
                          </span>
                          {/* Pin indicator - using Star for now */}
                          {assignments['work']?.includes(conv.id) && (
                            <Star className="w-3 h-3 text-[#7678ed] fill-[#7678ed]" />
                          )}
                        </div>
                    </div>
                    <p className={`text-xs truncate ${conv.unread_count > 0 ? 'text-[#202022] font-medium' : 'text-[#202022]/50'}`}>
                        {truncateMessage(conv.last_message)}
                    </p>
                    {/* Folder Badges (only in All Chats view) */}
                    {!activeFolderId && (
                        <div className="flex gap-1 mt-1.5">
                            {assignments['work']?.includes(conv.id) && (
                                <span className="inline-flex items-center px-1.5 py-0.5 rounded text-[10px] font-medium bg-[#7678ed]/20 text-[#7678ed]">Work</span>
                            )}
                            {assignments['friends']?.includes(conv.id) && (
                                <span className="inline-flex items-center px-1.5 py-0.5 rounded text-[10px] font-medium bg-green-500/20 text-green-600">Friends</span>
                            )}
                        </div>
                    )}
                    </div>
                </button>
            </ContextMenuTrigger>
            <ContextMenuContent className="w-48 bg-white border-[#7678ed]/20 text-[#202022] shadow-lg">
                <ContextMenuItem
                    className="focus:bg-[#7678ed]/10 focus:text-[#202022]"
                    onClick={() => handleToggleFolder('work', conv.id)}
                >
                    {assignments['work']?.includes(conv.id) ? (
                        <span className="text-[#7678ed]">Remove from Work</span>
                    ) : (
                        <span className="flex items-center gap-2"><Briefcase className="w-4 h-4" /> Add to Work</span>
                    )}
                </ContextMenuItem>
                <ContextMenuItem
                    className="focus:bg-[#7678ed]/10 focus:text-[#202022]"
                    onClick={() => handleToggleFolder('friends', conv.id)}
                >
                     {assignments['friends']?.includes(conv.id) ? (
                        <span className="text-green-600">Remove from Friends</span>
                    ) : (
                        <span className="flex items-center gap-2"><UserCheck className="w-4 h-4" /> Add to Friends</span>
                    )}
                </ContextMenuItem>
                <ContextMenuSeparator className="bg-[#7678ed]/10" />
                <ContextMenuItem
                    className="focus:bg-[#7678ed]/10 focus:text-[#202022]"
                    onClick={() => handleToggleFolder('archive', conv.id)}
                >
                    {assignments['archive']?.includes(conv.id) ? (
                        <span className="text-[#ff7a55]">Unarchive</span>
                    ) : (
                        <span className="flex items-center gap-2"><Archive className="w-4 h-4" /> Archive</span>
                    )}
                </ContextMenuItem>
            </ContextMenuContent>
          </ContextMenu>
        ))}
      </div>
    </aside>
  );
}
