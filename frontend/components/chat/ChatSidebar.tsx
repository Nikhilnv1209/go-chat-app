'use client';

import { useEffect, useState } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import { useQuery } from '@tanstack/react-query';
import { Search, MessageSquare, Users, LogOut, Settings, X, Briefcase, Archive, UserCheck } from 'lucide-react';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import { setConversations, setActiveConversation, resetUnread } from '@/store/features/conversationSlice';
import { toggleFolderAssignment } from '@/store/features/folderSlice';
import { conversationApi } from '@/lib/conversationApi';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { ContextMenu, ContextMenuContent, ContextMenuItem, ContextMenuTrigger, ContextMenuSeparator, ContextMenuSub, ContextMenuSubTrigger, ContextMenuSubContent } from '@/components/ui/context-menu';
import { Conversation } from '@/types';
import { logout } from '@/store/features/authSlice';

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
    refetchInterval: 30000, // Refresh every 30 seconds
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
    <aside className={`flex flex-col w-full h-full bg-slate-950/50 backdrop-blur-xl md:bg-transparent md:backdrop-blur-none fixed md:relative z-30 md:z-0 transform transition-transform duration-300 ease-in-out ${
      isOpen ? 'translate-x-0' : '-translate-x-full md:translate-x-0'
    }`}>
      {/* Header */}
      <div className="flex-shrink-0 p-4 border-b border-white/[0.1]">
        {/* Mobile Header Controls */}
        <div className="md:hidden flex items-center justify-end mb-4">
          <div className="flex items-center gap-2">
            {/* Mobile Close Button */}
            <Button
              variant="ghost"
              size="icon"
              onClick={onClose}
              className="h-9 w-9 text-slate-400 hover:text-white hover:bg-white/[0.05]"
            >
              <X className="w-4 h-4" />
            </Button>
            <Button
              variant="ghost"
              size="icon"
              onClick={() => router.push('/dashboard/profile')}
              className="h-9 w-9 text-slate-400 hover:text-white hover:bg-white/[0.05]"
            >
              <Settings className="w-4 h-4" />
            </Button>
            <Button
              variant="ghost"
              size="icon"
              onClick={handleLogout}
              className="h-9 w-9 text-slate-400 hover:text-red-400 hover:bg-red-500/10"
            >
              <LogOut className="w-4 h-4" />
            </Button>
          </div>
        </div>

        {/* Search */}
        <div className="relative">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-slate-500" />
          <Input
            type="text"
            placeholder="Search conversations..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="pl-10 h-10 bg-white/[0.05] border-white/[0.1] text-white placeholder:text-slate-500 focus:border-indigo-500/50 focus:ring-indigo-500/20 rounded-xl"
          />
        </div>
      </div>

      {/* User Profile Mini - Mobile Only */}
      {user && (
        <div className="flex-shrink-0 p-4 border-b border-white/[0.1] bg-gradient-to-r from-indigo-500/10 via-purple-500/10 to-pink-500/10 md:hidden">
          <div className="flex items-center gap-3">
            <div className="relative">
              <div className="w-10 h-10 rounded-full bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center text-white font-semibold">
                {user.username.charAt(0).toUpperCase()}
              </div>
              <div className="absolute bottom-0 right-0 w-3 h-3 bg-green-500 border-2 border-slate-950 rounded-full"></div>
            </div>
            <div className="flex-1 min-w-0">
              <p className="text-sm font-semibold text-white truncate">{user.username}</p>
              <p className="text-xs text-slate-400 truncate">{user.email}</p>
            </div>
          </div>
        </div>
      )}

      {/* Conversation List */}
      <div className="flex-1 overflow-y-auto">
        {isLoading && (
          <div className="flex items-center justify-center py-12">
            <div className="animate-spin rounded-full h-8 w-8 border-2 border-indigo-500 border-t-transparent"></div>
          </div>
        )}

        {isError && (
          <div className="p-4 text-center">
            <p className="text-sm text-red-400">Failed to load conversations</p>
          </div>
        )}

        {!isLoading && !isError && filteredConversations.length === 0 && (
          <div className="flex flex-col items-center justify-center py-12 px-4">
            <div className="w-16 h-16 rounded-full bg-white/[0.05] flex items-center justify-center mb-4">
              {activeFolderId ? <Briefcase className="w-8 h-8 text-slate-500" /> : <MessageSquare className="w-8 h-8 text-slate-500" />}
            </div>
            <p className="text-sm text-slate-400 text-center">
              {searchQuery
                ? 'No conversations found'
                : (activeFolderId ? 'This folder is empty' : 'No conversations yet')}
            </p>
            {!activeFolderId && (
                <p className="text-xs text-slate-500 text-center mt-1">
                Start a new chat to get started
                </p>
            )}
            {activeFolderId && (
                <p className="text-xs text-slate-500 text-center mt-1">
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
                    className={`w-full p-4 flex items-start gap-3 hover:bg-white/[0.05] transition-colors border-b border-white/[0.05] ${
                    activeConversationId === conv.id ? 'bg-white/[0.08]' : ''
                    }`}
                >
                    {/* Avatar */}
                    <div className="relative flex-shrink-0">
                    {conv.type === 'GROUP' ? (
                        <div className="w-12 h-12 rounded-full bg-gradient-to-br from-purple-500 to-pink-600 flex items-center justify-center text-white">
                        <Users className="w-6 h-6" />
                        </div>
                    ) : (
                        <div className="w-12 h-12 rounded-full bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center text-white font-semibold text-sm">
                        {conv.target_name.charAt(0).toUpperCase()}
                        </div>
                    )}
                    {conv.is_online && conv.type === 'DM' && (
                        <div className="absolute bottom-0 right-0 w-3.5 h-3.5 bg-green-500 border-2 border-slate-950 rounded-full"></div>
                    )}
                    {conv.unread_count > 0 && (
                        <div className="absolute -top-1 -right-1 min-w-[20px] h-5 bg-indigo-500 rounded-full flex items-center justify-center px-1.5">
                        <span className="text-xs font-bold text-white">{conv.unread_count > 99 ? '99+' : conv.unread_count}</span>
                        </div>
                    )}
                    </div>

                    {/* Content */}
                    <div className="flex-1 min-w-0 text-left">
                    <div className="flex items-center justify-between mb-1">
                        <h3 className="text-sm font-semibold text-white truncate">{conv.target_name}</h3>
                        <span className="text-xs text-slate-500 flex-shrink-0 ml-2">
                        {formatTimestamp(conv.last_message_at)}
                        </span>
                    </div>
                    <p className={`text-xs truncate ${conv.unread_count > 0 ? 'text-slate-300 font-medium' : 'text-slate-500'}`}>
                        {truncateMessage(conv.last_message)}
                    </p>
                    {/* Folder Badges (only in All Chats view) */}
                    {!activeFolderId && (
                        <div className="flex gap-1 mt-1">
                            {assignments['work']?.includes(conv.id) && (
                                <span className="inline-flex items-center px-1.5 py-0.5 rounded text-[10px] font-medium bg-blue-500/20 text-blue-300">Work</span>
                            )}
                            {assignments['friends']?.includes(conv.id) && (
                                <span className="inline-flex items-center px-1.5 py-0.5 rounded text-[10px] font-medium bg-green-500/20 text-green-300">Friends</span>
                            )}
                        </div>
                    )}
                    </div>
                </button>
            </ContextMenuTrigger>
            <ContextMenuContent className="w-48 bg-slate-900 border-white/10 text-white">
                <ContextMenuItem
                    className="focus:bg-white/10 focus:text-white"
                    onClick={() => handleToggleFolder('work', conv.id)}
                >
                    {assignments['work']?.includes(conv.id) ? (
                        <span className="text-blue-400">Remove from Work</span>
                    ) : (
                        <span className="flex items-center gap-2"><Briefcase className="w-4 h-4" /> Add to Work</span>
                    )}
                </ContextMenuItem>
                <ContextMenuItem
                    className="focus:bg-white/10 focus:text-white"
                    onClick={() => handleToggleFolder('friends', conv.id)}
                >
                     {assignments['friends']?.includes(conv.id) ? (
                        <span className="text-green-400">Remove from Friends</span>
                    ) : (
                        <span className="flex items-center gap-2"><UserCheck className="w-4 h-4" /> Add to Friends</span>
                    )}
                </ContextMenuItem>
                <ContextMenuSeparator className="bg-white/10" />
                <ContextMenuItem
                    className="focus:bg-white/10 focus:text-white"
                    onClick={() => handleToggleFolder('archive', conv.id)}
                >
                    {assignments['archive']?.includes(conv.id) ? (
                        <span className="text-amber-400">Unarchive</span>
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
