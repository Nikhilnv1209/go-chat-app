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
      <div className="flex-1 flex items-center justify-center">
        <div className="animate-spin rounded-full h-8 w-8 border-2 border-[#7678ed] border-t-transparent"></div>
      </div>
    );
  }

  if (messages.length === 0) {
    return (
      <div className="flex-1 flex flex-col items-center justify-center text-[#202022]/50">
        <p>No messages yet.</p>
        <p className="text-sm">Say hello! 👋</p>
      </div>
    );
  }

  const messageGroups = groupMessages(messages);

  return (
    <ScrollArea className="flex-1 p-4" ref={scrollRef}>
      <div className="space-y-4 pb-4">
        {messageGroups.map((group, groupIndex) => {
          const isMe = group.senderId === user?.id;
          const firstMessage = group.messages[0];

          // Determine avatar fallback and name
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
                "flex gap-2.5",
                isMe ? "justify-end" : "justify-start"
              )}
            >
              {/* Avatar - only for other users, aligned to bottom of last message */}
              {!isMe && (
                <div className="flex flex-col justify-end pb-1">
                  <Avatar className="h-9 w-9 flex-shrink-0">
                    <AvatarFallback className="bg-gradient-to-br from-[#7678ed] to-[#5a5cd9] text-white text-xs font-semibold">
                      {avatarLetter}
                    </AvatarFallback>
                  </Avatar>
                </div>
              )}

              {/* Message bubbles container */}
              <div className={cn(
                "flex flex-col gap-0.5 max-w-[65%]",
                isMe ? "items-end" : "items-start"
              )}>
                {/* Individual messages in the group */}
                {group.messages.map((message, msgIndex) => {
                  const isLast = msgIndex === group.messages.length - 1;

                  return (
                    <div
                      key={message.id}
                      className={cn(
                        "px-3 py-2 text-[13px] leading-relaxed",
                        isMe
                          ? "bg-[#7678ed] text-white"
                          : "bg-white text-[#202022]",
                        // Rounded corners - only last message has tail
                        isLast ? (
                          isMe
                            ? "rounded-[18px] rounded-br-[4px]" // My message tail bottom-right
                            : "rounded-[18px] rounded-bl-[4px]" // Their message tail bottom-left
                        ) : (
                          "rounded-[18px]" // Earlier messages fully rounded
                        ),
                        !isMe && "shadow-sm"
                      )}
                    >
                      {/* Sender name - INSIDE bubble, only for first message of group */}
                      {msgIndex === 0 && !isMe && displayName && (
                        <div className="text-[11px] font-semibold text-[#7678ed] mb-1">
                          {displayName}
                        </div>
                      )}

                      {/* Message content */}
                      <p className="whitespace-pre-wrap break-words">{message.content}</p>

                      {/* Time - INSIDE bubble at bottom right */}
                      <div className={cn(
                        "text-[10px] mt-1 text-right opacity-60",
                        isMe ? "text-white" : "text-[#202022]"
                      )}>
                        {format(new Date(message.created_at), 'HH:mm')}
                      </div>
                    </div>
                  );
                })}
              </div>

              {/* Avatar for own messages - aligned to bottom of last message */}
              {isMe && (
                <div className="flex flex-col justify-end pb-1">
                  <Avatar className="h-9 w-9 flex-shrink-0">
                    <AvatarFallback className="bg-gradient-to-br from-[#ff7a55] to-[#e66a47] text-white text-xs font-semibold">
                      {avatarLetter}
                    </AvatarFallback>
                  </Avatar>
                </div>
              )}
            </div>
          );
        })}
        <div ref={bottomRef} />
      </div>
    </ScrollArea>
  );
}
