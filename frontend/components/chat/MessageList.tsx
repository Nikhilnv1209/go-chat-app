import React, { useEffect, useRef } from 'react';
import { Message } from '@/types';
import { useAppSelector } from '@/store/hooks';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { cn } from '@/lib/utils';
import { format } from 'date-fns';

interface MessageListProps {
  messages: Message[];
  isLoading: boolean;
  conversationType?: 'DM' | 'GROUP';
  targetName?: string;
}

export default function MessageList({ messages, isLoading, conversationType, targetName }: MessageListProps) {
  const { user } = useAppSelector((state) => state.auth);
  const scrollRef = useRef<HTMLDivElement>(null);
  const bottomRef = useRef<HTMLDivElement>(null);

  // Auto-scroll to bottom when messages change
  useEffect(() => {
    if (bottomRef.current) {
      bottomRef.current.scrollIntoView({ behavior: 'smooth' });
    }
  }, [messages]);

  if (isLoading) {
    return (
      <div className="flex-1 flex items-center justify-center">
        <div className="animate-spin rounded-full h-8 w-8 border-2 border-indigo-500 border-t-transparent"></div>
      </div>
    );
  }

  if (messages.length === 0) {
    return (
      <div className="flex-1 flex flex-col items-center justify-center text-slate-500">
        <p>No messages yet.</p>
        <p className="text-sm">Say hello! 👋</p>
      </div>
    );
  }

  return (
    <ScrollArea className="flex-1 p-4" ref={scrollRef}>
      <div className="space-y-4 pb-4">
        {messages.map((message) => {
          const isMe = message.sender_id === user?.id;

          // Determine avatar fallback and name
          let avatarLetter = '?';
          let displayName = '';

          if (!isMe) {
              if (message.sender?.username) {
                  displayName = message.sender.username;
                  avatarLetter = displayName.charAt(0).toUpperCase();
              } else if (conversationType === 'DM' && targetName) {
                  displayName = targetName;
                  avatarLetter = displayName.charAt(0).toUpperCase();
              }
          }

          return (
            <div
              key={message.id}
              className={cn(
                "flex w-full gap-2",
                isMe ? "justify-end" : "justify-start"
              )}
            >
              {!isMe && (
                 <div className="flex flex-col gap-1">
                    <Avatar className="h-8 w-8 mt-1">
                        <AvatarFallback className="bg-gradient-to-br from-indigo-500 to-purple-600 text-white text-xs">
                            {avatarLetter}
                        </AvatarFallback>
                    </Avatar>
                </div>
              )}

              <div className="flex flex-col max-w-[70%]">
                 {/* Show sender name in groups if it's not me */}
                 {!isMe && conversationType === 'GROUP' && displayName && (
                     <span className="text-[10px] text-slate-400 ml-1 mb-1">{displayName}</span>
                 )}

                  <div
                    className={cn(
                      "rounded-2xl px-4 py-2 text-sm shadow-sm",
                      isMe
                        ? "bg-indigo-600 text-white rounded-br-none"
                        : "bg-white/10 text-white rounded-bl-none"
                    )}
                  >
                    <p className="whitespace-pre-wrap break-words">{message.content}</p>
                    <div className={cn(
                        "text-[10px] mt-1 text-right opacity-70",
                        isMe ? "text-indigo-200" : "text-slate-400"
                    )}>
                      {format(new Date(message.created_at), 'h:mm a')}
                    </div>
                  </div>
              </div>
            </div>
          );
        })}
        <div ref={bottomRef} />
      </div>
    </ScrollArea>
  );
}
