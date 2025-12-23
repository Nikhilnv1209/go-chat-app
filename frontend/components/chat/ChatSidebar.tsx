'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { useQuery } from '@tanstack/react-query';
import { Search, MessageSquare, Users, LogOut, Settings, X } from 'lucide-react';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import { setConversations, setActiveConversation, resetUnread } from '@/store/features/conversationSlice';
import { conversationApi } from '@/lib/conversationApi';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Conversation } from '@/types';
import { logout } from '@/store/features/authSlice';

interface ChatSidebarProps {
  isOpen?: boolean;
  onClose?: () => void;
}

export default function ChatSidebar({ isOpen, onClose }: ChatSidebarProps = {}) {
  const router = useRouter();
  const dispatch = useAppDispatch();
  const { token, user } = useAppSelector((state) => state.auth);
  const { conversations, activeConversationId } = useAppSelector((state) => state.conversation);
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

  const filteredConversations = conversations.filter((conv) =>
    conv.target_name.toLowerCase().includes(searchQuery.toLowerCase())
  );

  const handleConversationClick = (conv: Conversation) => {
    dispatch(setActiveConversation(conv.id));
    dispatch(resetUnread(conv.id));
    router.push(`/c/${conv.target_id}`);
    // Close sidebar on mobile after clicking a conversation
    if (onClose) {
      onClose();
    }
  };

  const handleLogout = () => {
    dispatch(logout());
    router.push('/login');
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

  return (
    <aside className={`flex flex-col w-full md:w-80 lg:w-96 border-r border-white/[0.1] bg-slate-950/50 backdrop-blur-xl fixed md:relative h-full z-30 md:z-0 transform transition-transform duration-300 ease-in-out ${
      isOpen ? 'translate-x-0' : '-translate-x-full md:translate-x-0'
    }`}>
      {/* Header */}
      <div className="flex-shrink-0 p-4 border-b border-white/[0.1]">
        <div className="flex items-center justify-between mb-4">
          <div className="hidden md:flex items-center gap-2 min-w-0">
            <MessageSquare className="w-5 h-5 text-indigo-400 flex-shrink-0" />
            <h2 className="text-xl font-bold text-white whitespace-nowrap">Chats</h2>
          </div>
          <div className="flex items-center gap-2 md:hidden">
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
              <MessageSquare className="w-8 h-8 text-slate-500" />
            </div>
            <p className="text-sm text-slate-400 text-center">
              {searchQuery ? 'No conversations found' : 'No conversations yet'}
            </p>
            <p className="text-xs text-slate-500 text-center mt-1">
              Start a new chat to get started
            </p>
          </div>
        )}

        {!isLoading && !isError && filteredConversations.map((conv) => (
          <button
            key={conv.id}
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
            </div>
          </button>
        ))}
      </div>
    </aside>
  );
}
