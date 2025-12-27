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

// Group consecutive messages from the same sender
function groupMessages(messages: Message[]) {
  const groups: { senderId: string; messages: Message[] }[] = [];

  messages.forEach((message) => {
    const lastGroup = groups[groups.length - 1];
    if (lastGroup && lastGroup.senderId === message.sender_id) {
      lastGroup.messages.push(message);
    } else {
      groups.push({ senderId: message.sender_id, messages: [message] });
    }
  });

  return groups;
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
      <div className="flex-1 h-full flex items-center justify-center bg-[#f9fafc]">
        <div className="animate-spin rounded-full h-8 w-8 border-2 border-[#7678ed] border-t-transparent"></div>
      </div>
    );
  }

  if (messages.length === 0) {
    return (
      <div className="flex-1 h-full flex flex-col items-center justify-center text-[#202022]/50 bg-[#f9fafc]">
        <p className="font-medium text-base">No messages yet</p>
        <p className="text-sm">Start the conversation! 👋</p>
      </div>
    );
  }

  const messageGroups = groupMessages(messages);

  return (
    <div
      className="flex-1 overflow-y-auto p-4 space-y-4 scroll-smooth"
      ref={scrollRef}
      style={{ scrollbarWidth: 'thin', scrollbarColor: '#7678ed40 transparent' }}
    >
      <div className="flex flex-col gap-4 pb-4">
        {messageGroups.map((group, groupIndex) => {
          const isMe = group.senderId === user?.id;
          const firstMessage = group.messages[0];

          let avatarLetter = '?';
          let displayName = '';

          if (isMe) {
            displayName = user?.username || 'You';
            avatarLetter = displayName.charAt(0).toUpperCase();
          } else {
            if (firstMessage.sender?.username) {
              displayName = firstMessage.sender.username;
              avatarLetter = displayName.charAt(0).toUpperCase();
            } else if (conversationType === 'DM' && targetName) {
              displayName = targetName;
              avatarLetter = displayName.charAt(0).toUpperCase();
            }
          }

          return (
            <div
              key={`group-${groupIndex}`}
              className={cn(
                "flex gap-3 w-full",
                isMe ? "flex-row-reverse" : "flex-row"
              )}
            >
              {/* Avatar */}
              <div className="flex flex-col justify-end pb-1 flex-shrink-0">
                <Avatar className="h-9 w-9 border border-[#7678ed]/10 shadow-sm">
                  <AvatarFallback className={cn(
                    "text-white text-xs font-bold",
                    isMe
                      ? "bg-gradient-to-br from-[#ff7a55] to-[#e66a47]"
                      : "bg-gradient-to-br from-[#7678ed] to-[#5a5cd9]"
                  )}>
                    {avatarLetter}
                  </AvatarFallback>
                </Avatar>
              </div>

              {/* Message bubbles container */}
              <div className={cn(
                "flex flex-col gap-1 max-w-[75%]",
                isMe ? "items-end" : "items-start"
              )}>
                {group.messages.map((message, msgIndex) => {
                  const isLast = msgIndex === group.messages.length - 1;

                  return (
                    <div
                      key={message.id}
                      className={cn(
                        "px-4 py-2.5 text-[14px] leading-relaxed shadow-sm transition-all",
                        isMe
                          ? "bg-[#7678ed] text-white"
                          : "bg-white text-[#202022] border border-[#7678ed]/5",
                        isMe ? (
                          isLast
                            ? "rounded-2xl rounded-br-sm"
                            : "rounded-2xl rounded-br-2xl"
                        ) : (
                          isLast
                            ? "rounded-2xl rounded-bl-sm"
                            : "rounded-2xl rounded-bl-2xl"
                        )
                      )}
                    >
                      {msgIndex === 0 && !isMe && displayName && (
                        <div className="text-[11px] font-bold text-[#7678ed] mb-1.5 uppercase tracking-wider">
                          {displayName}
                        </div>
                      )}

                      <p className="whitespace-pre-wrap break-words">{message.content}</p>

                      <div className={cn(
                        "text-[10px] mt-1.5 text-right font-medium",
                        isMe ? "text-white/80" : "text-[#202022]/40"
                      )}>
                        {format(new Date(message.created_at), 'HH:mm')}
                      </div>
                    </div>
                  );
                })}
              </div>
            </div>
          );
        })}
        <div ref={bottomRef} className="h-2" />
      </div>
    </div>
  );
}
