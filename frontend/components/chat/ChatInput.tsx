import React, { useState, useRef, useEffect } from 'react';
import { Send } from 'lucide-react';
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
    <div className="p-4 border-t border-white/[0.1] bg-slate-950/50 backdrop-blur-md">
      <div className="flex items-end gap-2 bg-white/[0.05] p-1.5 rounded-xl border border-white/[0.1] focus-within:border-indigo-500/50 transition-colors">
        <Textarea
          ref={textareaRef}
          value={message}
          onChange={handleInput}
          onKeyDown={handleKeyDown}
          placeholder="Type a message..."
          className="min-h-[20px] max-h-[120px] bg-transparent border-0 focus-visible:ring-0 focus-visible:ring-offset-0 text-white placeholder:text-slate-500 resize-none py-2.5 px-3"
          rows={1}
        />
        <Button
          size="icon"
          onClick={handleSend}
          disabled={!message.trim() || isLoading}
          className="h-10 w-10 shrink-0 bg-indigo-600 hover:bg-indigo-500 text-white rounded-lg transition-all disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <Send className="w-5 h-5" />
        </Button>
      </div>
      <div className="text-xs text-slate-500 mt-2 text-center opacity-0 group-focus-within:opacity-100 transition-opacity">
        Press Enter to send, Shift + Enter for new line
      </div>
    </div>
  );
}
