import React, { useState, useRef, useEffect } from 'react';
import { Send, Paperclip, Mic } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Textarea } from '@/components/ui/textarea';

interface ChatInputProps {
  onSendMessage: (content: string) => void;
  onTyping?: (isTyping: boolean) => void;
  isLoading?: boolean;
}

export default function ChatInput({ onSendMessage, onTyping, isLoading = false }: ChatInputProps) {
  const [message, setMessage] = useState('');
  const textareaRef = useRef<HTMLTextAreaElement>(null);
  const typingTimeoutRef = useRef<NodeJS.Timeout | null>(null);

  const handleSend = () => {
    if (message.trim() && !isLoading) {
      onSendMessage(message);
      setMessage('');

      // Stop typing immediately when sent
      if (onTyping && typingTimeoutRef.current) {
          clearTimeout(typingTimeoutRef.current);
          onTyping(false);
          typingTimeoutRef.current = null;
      }

      // Reset height
      if (textareaRef.current) {
        textareaRef.current.style.height = 'auto';
      }
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  const handleInput = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setMessage(e.target.value);

    // Typing Indicator Logic
    if (onTyping) {
        if (!typingTimeoutRef.current) {
            onTyping(true);
        } else {
            clearTimeout(typingTimeoutRef.current);
        }

        typingTimeoutRef.current = setTimeout(() => {
            onTyping(false);
            typingTimeoutRef.current = null;
        }, 1500); // Stop typing after 1.5s of inactivity
    }

    // Auto-resize
    if (textareaRef.current) {
        textareaRef.current.style.height = 'auto';
        textareaRef.current.style.height = `${textareaRef.current.scrollHeight}px`;
    }
  };

  return (
    <div className="p-4 border-t border-[#7678ed]/10 bg-white">
      <div className="flex items-end gap-2 bg-[#f9fafc] p-2 rounded-2xl border border-[#7678ed]/10 focus-within:border-[#7678ed]/30 transition-colors">
        {/* Attachment Button */}
        <Button
          variant="ghost"
          size="icon"
          className="h-9 w-9 shrink-0 text-[#202022]/40 hover:text-[#7678ed] hover:bg-[#7678ed]/10 rounded-xl"
        >
          <Paperclip className="w-5 h-5" />
        </Button>

        <Textarea
          ref={textareaRef}
          value={message}
          onChange={handleInput}
          onKeyDown={handleKeyDown}
          placeholder="Type a message..."
          className="min-h-[20px] max-h-[120px] bg-transparent border-0 focus-visible:ring-0 focus-visible:ring-offset-0 focus:outline-none text-[#202022] placeholder:text-[#202022]/40 resize-none py-2.5 px-1 shadow-none"
          rows={1}
        />

        {/* Mic Button */}
        <Button
          variant="ghost"
          size="icon"
          className="h-9 w-9 shrink-0 text-[#202022]/40 hover:text-[#7678ed] hover:bg-[#7678ed]/10 rounded-xl"
        >
          <Mic className="w-5 h-5" />
        </Button>

        {/* Send Button */}
        <Button
          size="icon"
          onClick={handleSend}
          disabled={!message.trim() || isLoading}
          className="h-10 w-10 shrink-0 bg-[#7678ed] hover:bg-[#5a5cd9] text-white rounded-xl transition-all disabled:opacity-50 disabled:cursor-not-allowed shadow-md shadow-[#7678ed]/20"
        >
          <Send className="w-5 h-5" />
        </Button>
      </div>
    </div>
  );
}
