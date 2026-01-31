'use client';

import { useState, useEffect } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import { useQuery } from '@tanstack/react-query';
import { Search, MessageSquare, Users, Briefcase, Archive, UserCheck, Star, Plus, MoreVertical } from 'lucide-react';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import { setConversations, setActiveConversation, resetUnread } from '@/store/features/conversationSlice';
import { toggleFolderAssignment } from '@/store/features/folderSlice';
import { logout } from '@/store/features/authSlice';
import { conversationApi } from '@/lib/conversationApi';
import { Input } from '@/components/ui/input';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { ContextMenu, ContextMenuContent, ContextMenuItem, ContextMenuTrigger, ContextMenuSeparator } from '@/components/ui/context-menu';
import { Conversation } from '@/types';
import { cn } from '@/lib/utils';
import { NewChatDialog } from './NewChatDialog';

interface MobileChatSidebarProps {
  onClose?: () => void;
}

interface Quote {
  quote: string;
  author: string;
}

export default function MobileChatSidebar({ onClose }: MobileChatSidebarProps) {
  const router = useRouter();
  const pathname = usePathname();
  const dispatch = useAppDispatch();
  const { token, user } = useAppSelector((state) => state.auth);
  const { conversations, activeConversationId } = useAppSelector((state) => state.conversation);
  const { assignments } = useAppSelector((state) => state.folders);
  const [searchQuery, setSearchQuery] = useState('');
  const [dailyQuote, setDailyQuote] = useState<Quote | null>(null);

  // Fetch conversations
  const { data, isLoading, isError } = useQuery({
    queryKey: ['conversations'],
    queryFn: () => conversationApi.getConversations(),
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

    if (activeFolderId === 'archive') {
      return isArchived;
    }

    if (isArchived) {
      return false;
    }

    if (activeFolderId) {
      return assignments[activeFolderId]?.includes(conv.id);
    }

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

  const handleToggleFolder = (folderId: string, conversationId: string) => {
    dispatch(toggleFolderAssignment({ folderId, conversationId }));
  };

  const getFormattedDate = () => {
    const now = new Date();
    const day = String(now.getDate()).padStart(2, '0');
    const month = String(now.getMonth() + 1).padStart(2, '0');
    const weekday = now.toLocaleDateString('en-US', { weekday: 'short' }).toUpperCase();
    return { day, month, weekday };
  };

  const { day, month, weekday } = getFormattedDate();

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

  const handleLogout = () => {
    dispatch(logout());
    router.push('/login');
  };

  // Fetch daily quote (local collection)
  useEffect(() => {
    const dayOfYear = Math.floor((new Date() - new Date(new Date().getFullYear(), 0, 0)) / (1000 * 60 * 60 * 24));
    const dailyQuotes = [
      { quote: "The only way to do great work is to love what you do.", author: "Steve Jobs" },
      { quote: "Success is not final, failure is not fatal: it is the courage to continue that counts.", author: "Winston Churchill" },
      { quote: "In the middle of every difficulty lies opportunity.", author: "Albert Einstein" },
      { quote: "Be the change you wish to see in the world.", author: "Mahatma Gandhi" },
      { quote: "The best way to predict the future is to create it.", author: "Peter Drucker" },
      { quote: "Life is what happens when you're busy making other plans.", author: "John Lennon" },
      { quote: "The purpose of our lives is to be happy.", author: "Dalai Lama" },
      { quote: "Stay hungry, stay foolish.", author: "Steve Jobs" },
      { quote: "Connection is the energy that exists between people when they feel seen, heard, and valued.", author: "Brené Brown" },
      { quote: "Talk to someone today, you might just make their day.", author: "Unknown" },
      { quote: "Every conversation is an opportunity to inspire.", author: "Unknown" },
      { quote: "Kind words can be short and easy to speak, but their echoes are truly endless.", author: "Mother Teresa" },
      { quote: "The more you praise and celebrate your life, the more there is in life to celebrate.", author: "Oprah Winfrey" }
    ];
    setDailyQuote(dailyQuotes[dayOfYear % dailyQuotes.length]);
  }, []);

  return (
    <div className="flex flex-col h-full w-full bg-gradient-to-b from-[#7678ed] via-[#6d6fe0] to-[#9ca3af] relative">
      {/* Background Header Content - Floats above gradient */}
      <div className="relative z-10 px-5 pt-8 pb-4">
        {/* Header Container with Greeting+Date and Menu */}
        <div className="flex flex-row items-start justify-between">
          {/* Greeting and Date Widget */}
          <div className="flex flex-col gap-2">
            <p className="text-white/90 text-base font-medium">
              {user ? `Hi, ${user.username}` : 'Hello'}
            </p>
            <span className="text-white text-5xl font-normal tracking-wide">{day}.{month} {weekday}</span>
          </div>

          {/* Menu Widget */}
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <button className="h-10 w-10 rounded-full bg-white/20 hover:bg-white/30 backdrop-blur-sm flex items-center justify-center text-white transition-all active:scale-95">
                <MoreVertical className="w-5 h-5" />
              </button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-48 bg-white border-0 shadow-2xl rounded-2xl p-1 mt-1">
              <div className="px-4 py-3 border-b border-[#7678ed]/10">
                <p className="text-sm font-semibold text-[#202022]">{user?.username}</p>
                <p className="text-xs text-[#202022]/50 truncate">{user?.email}</p>
              </div>
              <DropdownMenuSeparator className="bg-[#7678ed]/10 my-1" />
              <DropdownMenuItem
                className="focus:bg-[#7678ed]/10 focus:text-[#202022] cursor-pointer rounded-xl"
                onClick={() => router.push('/dashboard/profile')}
              >
                Profile
              </DropdownMenuItem>
              <DropdownMenuItem
                className="focus:bg-[#ff7a55]/10 focus:text-[#ff7a55] cursor-pointer rounded-xl"
                onClick={handleLogout}
              >
                Logout
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>

      {/* Daily Quote Section */}
      {dailyQuote && (
        <div className="relative z-10 px-5 pb-6">
          <div className="bg-white/10 backdrop-blur-sm rounded-2xl p-4 border border-white/20">
            <p className="text-white/95 text-sm italic leading-relaxed">"{dailyQuote.quote}"</p>
            <p className="text-white/60 text-xs mt-2 text-right">— {dailyQuote.author}</p>
          </div>
        </div>
      )}

      {/* Floating Card Content - Rounded top, floats on gradient */}
      <div className="flex-1 flex flex-col bg-[#f9fafc] rounded-t-3xl relative z-20 shadow-2xl overflow-hidden min-h-0 mt-2">
        {/* Search Section */}
        <div className="flex-shrink-0 px-4 pt-4 pb-3">
          <div className="relative">
            <Search className="absolute left-3.5 top-1/2 -translate-y-1/2 h-[18px] w-[18px] text-[#7678ed]/40" />
            <Input
              type="text"
              placeholder="Search conversations..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="pl-10 pr-4 h-11 bg-white/80 border border-gray-200/50 text-[#202022] text-[15px] placeholder:text-[#202022]/35 focus:bg-white focus:border-gray-300 focus:ring-0 focus:ring-offset-0 focus-visible:ring-0 rounded-xl shadow-sm transition-all"
            />
          </div>
        </div>

        {/* Filter Tabs */}
        <div className="flex-shrink-0 px-4 pb-3">
          <div className="flex gap-2">
            <button
              onClick={() => router.push('/dashboard')}
              className={cn(
                "flex-1 relative px-2 py-2 rounded-lg text-[13px] font-medium transition-all",
                !activeFolderId
                  ? "text-[#8a8cf5]"
                  : "text-[#202022]/45 hover:text-[#202022]/60"
              )}
            >
              All
              {!activeFolderId && <span className="absolute bottom-0 left-1/2 -translate-x-1/2 w-6 h-0.5 bg-[#8a8cf5] rounded-full"></span>}
            </button>
            <button
              onClick={() => router.push('/dashboard/work')}
              className={cn(
                "flex-1 relative px-2 py-2 rounded-lg text-[13px] font-medium transition-all",
                activeFolderId === 'work'
                  ? "text-[#8a8cf5]"
                  : "text-[#202022]/45 hover:text-[#202022]/60"
              )}
            >
              Work
              {activeFolderId === 'work' && <span className="absolute bottom-0 left-1/2 -translate-x-1/2 w-6 h-0.5 bg-[#8a8cf5] rounded-full"></span>}
            </button>
            <button
              onClick={() => router.push('/dashboard/friends')}
              className={cn(
                "flex-1 relative px-2 py-2 rounded-lg text-[13px] font-medium transition-all",
                activeFolderId === 'friends'
                  ? "text-[#8a8cf5]"
                  : "text-[#202022]/45 hover:text-[#202022]/60"
              )}
            >
              Friends
              {activeFolderId === 'friends' && <span className="absolute bottom-0 left-1/2 -translate-x-1/2 w-6 h-0.5 bg-[#8a8cf5] rounded-full"></span>}
            </button>
            <button
              onClick={() => router.push('/dashboard/archive')}
              className={cn(
                "flex-1 relative px-2 py-2 rounded-lg text-[13px] font-medium transition-all",
                activeFolderId === 'archive'
                  ? "text-[#8a8cf5]"
                  : "text-[#202022]/45 hover:text-[#202022]/60"
              )}
            >
              Archive
              {activeFolderId === 'archive' && <span className="absolute bottom-0 left-1/2 -translate-x-1/2 w-6 h-0.5 bg-[#8a8cf5] rounded-full"></span>}
            </button>
          </div>
        </div>

        {/* Conversation List - Part of card, no individual cards */}
        <div className="flex-1 overflow-y-auto px-0 pb-24">
          {isLoading && (
            <div className="flex items-center justify-center py-12">
              <div className="animate-spin rounded-full h-8 w-8 border-2 border-[#7678ed] border-t-transparent"></div>
            </div>
          )}

          {isError && (
            <div className="p-5 text-center">
              <p className="text-sm text-[#ff7a55]">Failed to load conversations</p>
            </div>
          )}

          {!isLoading && !isError && filteredConversations.length === 0 && (
            <div className="flex flex-col items-center justify-center py-12 px-4">
              <div className="w-20 h-20 rounded-2xl bg-[#7678ed]/10 flex items-center justify-center mb-5">
                {activeFolderId ? <Briefcase className="w-10 h-10 text-[#7678ed]" /> : <MessageSquare className="w-10 h-10 text-[#7678ed]" />}
              </div>
              <p className="text-base text-[#202022]/60 text-center font-medium">
                {searchQuery
                  ? 'No conversations found'
                  : (activeFolderId ? 'This folder is empty' : 'No conversations yet')}
              </p>
              {!activeFolderId && (
                  <p className="text-sm text-[#202022]/40 text-center mt-2">
                    Tap the + button to start a new chat
                  </p>
              )}
              {activeFolderId && (
                  <p className="text-sm text-[#202022]/40 text-center mt-2">
                    Long press on a chat in "All Chats" to add it here
                  </p>
              )}
            </div>
          )}

          {!isLoading && !isError && filteredConversations.map((conv, index) => (
            <ContextMenu key={conv.id}>
              <ContextMenuTrigger>
                  <div className="w-full">
                    <button
                        onClick={() => handleConversationClick(conv)}
                        className={cn(
                          "w-full flex items-center gap-3 px-4 py-2.5 transition-all duration-200",
                          activeConversationId === conv.id ? "bg-[#8a8cf5]/12" : "hover:bg-white/50"
                        )}
                    >
                        <div className="relative flex-shrink-0">
                        {conv.type === 'GROUP' ? (
                            <div className="w-11 h-11 rounded-full bg-gradient-to-br from-[#ff7a55] to-[#e66a47] flex items-center justify-center text-white">
                            <Users className="w-5 h-5" />
                            </div>
                        ) : (
                            <div className="w-11 h-11 rounded-full bg-gradient-to-br from-[#8a8cf5] to-[#7678ed] flex items-center justify-center text-white font-semibold text-sm">
                            {conv.target_name.charAt(0).toUpperCase()}
                            </div>
                        )}
                        {conv.is_online && conv.type === 'DM' && (
                            <div className="absolute bottom-0 right-0 w-2.5 h-2.5 bg-green-500 border-2 border-white rounded-full"></div>
                        )}
                        {conv.unread_count > 0 && (
                            <div className="absolute -top-0.5 -right-0.5 min-w-[16px] h-4.5 bg-[#ff7a55] rounded-full flex items-center justify-center px-1">
                            <span className="text-[9px] font-bold text-white">{conv.unread_count > 99 ? '99+' : conv.unread_count}</span>
                            </div>
                        )}
                        </div>

                        <div className="flex-1 min-w-0 text-left">
                        <div className="flex items-center justify-between mb-0.5">
                            <h3 className="text-sm font-medium text-[#202022] truncate pr-2">{conv.target_name}</h3>
                            <div className="flex items-center gap-1 flex-shrink-0">
                              <span className="text-[11px] text-[#202022]/30">
                              {formatTimestamp(conv.last_message_at)}
                              </span>
                              {assignments['work']?.includes(conv.id) && (
                                <Star className="w-3 h-3 text-[#8a8cf5] fill-[#8a8cf5]" />
                              )}
                            </div>
                        </div>
                        <p className={cn("text-xs truncate leading-relaxed", conv.unread_count > 0 ? "text-[#202022]/75 font-medium" : "text-[#202022]/40")}>
                            {truncateMessage(conv.last_message, 40)}
                        </p>
                        </div>
                    </button>
                    {index !== filteredConversations.length - 1 && (
                      <div className="mx-4 h-px bg-gray-200/60"></div>
                    )}
                  </div>
              </ContextMenuTrigger>
              <ContextMenuContent className="w-48 bg-white border-0 shadow-2xl rounded-2xl p-1">
                  <ContextMenuItem
                      className="focus:bg-[#7678ed]/10 focus:text-[#202022] cursor-pointer rounded-xl"
                      onClick={() => handleToggleFolder('work', conv.id)}
                  >
                      {assignments['work']?.includes(conv.id) ? (
                          <span className="text-[#7678ed]">Remove from Work</span>
                      ) : (
                          <span className="flex items-center gap-2"><Briefcase className="w-4 h-4" /> Add to Work</span>
                      )}
                  </ContextMenuItem>
                  <ContextMenuItem
                      className="focus:bg-[#7678ed]/10 focus:text-[#202022] cursor-pointer rounded-xl"
                      onClick={() => handleToggleFolder('friends', conv.id)}
                  >
                       {assignments['friends']?.includes(conv.id) ? (
                          <span className="text-green-600">Remove from Friends</span>
                       ) : (
                          <span className="flex items-center gap-2"><UserCheck className="w-4 h-4" /> Add to Friends</span>
                       )}
                  </ContextMenuItem>
                  <DropdownMenuSeparator className="bg-[#7678ed]/10 my-1" />
                  <ContextMenuItem
                      className="focus:bg-[#7678ed]/10 focus:text-[#202022] cursor-pointer rounded-xl"
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

        {/* Floating New Chat Button */}
        <div className="absolute bottom-6 right-6 z-30">
          <NewChatDialog
            trigger={
              <button className="h-16 w-16 rounded-full bg-gradient-to-br from-[#7678ed] to-[#5a5cd9] hover:from-[#6567d9] hover:to-[#4e50cd] text-white shadow-xl shadow-[#7678ed]/30 transition-all hover:scale-105 active:scale-95 flex items-center justify-center">
                <Plus className="w-7 h-7" />
              </button>
            }
          />
        </div>
      </div>
    </div>
  );
}
