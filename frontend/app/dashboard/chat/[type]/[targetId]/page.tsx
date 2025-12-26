'use client';

import { useEffect, useState, use } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { useQuery } from '@tanstack/react-query';
import { MoreVertical, Phone, Video, Search, ArrowLeft } from 'lucide-react';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import { setActiveConversation } from '@/store/features/conversationSlice';
import { conversationApi } from '@/lib/conversationApi';
import MessageList from '@/components/chat/MessageList';
import ChatInput from '@/components/chat/ChatInput';
import { Button } from '@/components/ui/button';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { cn } from '@/lib/utils';
import { socketService } from '@/lib/socketService';

export default function ChatPage() {
  const { type: typeParam, targetId: targetIdParam } = useParams();
  const type = (Array.isArray(typeParam) ? typeParam[0] : typeParam)?.toUpperCase() as 'DM' | 'GROUP';
  const targetId = Array.isArray(targetIdParam) ? targetIdParam[0] : targetIdParam;

  const router = useRouter();
  const dispatch = useAppDispatch();
  const { conversations } = useAppSelector((state) => state.conversation);
  const { token } = useAppSelector((state) => state.auth);

  // Find the current conversation details from the store
  const conversation = conversations.find(
    (c) => c.target_id === targetId && c.type === type
  );

  // Set active conversation in Redux
  useEffect(() => {
    if (conversation) {
      dispatch(setActiveConversation(conversation.id));
    } else {
        dispatch(setActiveConversation(null));
    }

    return () => {
        dispatch(setActiveConversation(null));
    };
  }, [conversation, dispatch, targetId, type]);

  // Fetch Messages
  const { data: messages, isLoading } = useQuery({
    queryKey: ['messages', targetId, type],
    queryFn: () => conversationApi.getMessages(token!, targetId!, type),
    enabled: !!token && !!targetId && !!type,
  });

  const handleSendMessage = (content: string) => {
    if (targetId && type) {
        socketService.sendMessage(content, targetId, type);
    }
  };

  const handleBack = () => {
      router.push('/dashboard');
  };

  if (!targetId || !type) {
      return null;
  }

  const handleTyping = (isTyping: boolean) => {
    if (targetId && type) {
        socketService.sendTyping(targetId, type, isTyping);
    }
  };

  return (
    <div className="flex flex-col h-full w-full bg-[#f9fafc] relative overflow-hidden">
      {/* Header */}
      <div className="h-16 border-b border-[#7678ed]/10 bg-white flex items-center justify-between px-4 sticky top-0 z-10 shadow-sm">
        <div className="flex items-center gap-3">
            <Button variant="ghost" size="icon" className="md:hidden -ml-2 text-[#202022]/50 hover:text-[#202022] hover:bg-[#7678ed]/10" onClick={handleBack}>
                <ArrowLeft className="w-5 h-5" />
            </Button>

          <Avatar className="h-10 w-10">
            <AvatarFallback className={cn(
                "text-white font-semibold",
                type === 'GROUP' ? "bg-gradient-to-br from-[#ff7a55] to-[#e66a47]" : "bg-gradient-to-br from-[#7678ed] to-[#5a5cd9]"
            )}>
              {conversation ? conversation.target_name.charAt(0).toUpperCase() : '?'}
            </AvatarFallback>
          </Avatar>
          <div>
            <h2 className="text-sm font-semibold text-[#202022]">
              {conversation
                ? conversation.target_name
                : (conversations.length === 0 ? 'Loading...' : 'Chat not found')}
            </h2>
            {conversation && type === 'DM' && (
               <div className="flex items-center gap-1.5">
                   <div className={cn("w-2 h-2 rounded-full", conversation.is_online ? "bg-green-500" : "bg-[#202022]/30")}></div>
                   <span className="text-xs text-[#202022]/50">{conversation.is_online ? 'Online' : 'Offline'}</span>
               </div>
            )}
            {conversation && type === 'GROUP' && (
                 <p className="text-xs text-[#202022]/50">
                     {conversation.member_count ? `${conversation.member_count} members` : 'Group Chat'}
                 </p>
            )}
          </div>
        </div>

        <div className="flex items-center gap-1">
          <Button variant="ghost" size="icon" className="text-[#202022]/50 hover:text-[#7678ed] hover:bg-[#7678ed]/10">
            <Search className="w-5 h-5" />
          </Button>
          <Button variant="ghost" size="icon" className="text-[#202022]/50 hover:text-[#7678ed] hover:bg-[#7678ed]/10">
             <MoreVertical className="w-5 h-5" />
          </Button>
        </div>
      </div>

      {/* Messages Area */}
      <MessageList
        messages={messages || []}
        isLoading={isLoading}
        conversationType={type}
        targetName={conversation?.target_name}
      />

      {/* Input Area */}
      <ChatInput
        onSendMessage={handleSendMessage}
        onTyping={handleTyping}
        isLoading={false}
      />
    </div>
  );
}
