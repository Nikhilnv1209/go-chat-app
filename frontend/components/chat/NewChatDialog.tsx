'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { Search, UserPlus, X } from 'lucide-react';
import { useDebounce } from '@/hooks/use-debounce';
import { conversationApi } from '@/lib/conversationApi';
import { useAppSelector } from '@/store/hooks';
import { User } from '@/types';
import { Button } from '@/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription, // Added DialogDescription
  DialogTrigger,
} from '@/components/ui/dialog';
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from '@/components/ui/command'; // Check if these are exported from command.tsx
import { Avatar, AvatarFallback } from '@/components/ui/avatar';

export function NewChatDialog() {
  const router = useRouter();
  const [open, setOpen] = useState(false);
  const [query, setQuery] = useState('');
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(false);
  const { token, user: currentUser } = useAppSelector((state) => state.auth);

  const debouncedQuery = useDebounce(query, 300);

  useEffect(() => {
    async function search() {
      if (!debouncedQuery.trim() || !token) {
        setUsers([]);
        return;
      }

      setLoading(true);
      try {
        const results = await conversationApi.searchUsers(token, debouncedQuery);
        setUsers(results);
      } catch (error) {
        console.error('Failed to search users:', error);
      } finally {
        setLoading(false);
      }
    }

    search();
  }, [debouncedQuery, token]);

  const handleStartChat = (targetUser: User) => {
    // Navigate to DM URL: /dashboard/chat/dm/{targetId}
    setOpen(false);
    router.push(`/dashboard/chat/dm/${targetUser.id}`);
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button
            className="w-full justify-start gap-2 bg-[#7678ed] hover:bg-[#6567d9] text-white shadow-sm transition-all"
        >
          <UserPlus className="w-4 h-4" />
          <span>New Chat</span>
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px] bg-white p-0 gap-0 overflow-hidden border-0 shadow-2xl">
        <div className="p-6 pb-2 border-b border-[#7678ed]/10 bg-[#f9fafc]">
            <DialogHeader>
            <DialogTitle className="text-[#202022] flex items-center gap-2">
                <UserPlus className="w-5 h-5 text-[#7678ed]" />
                New Conversation
            </DialogTitle>
             <DialogDescription className="text-[#202022]/60 mt-1.5">
                Search for people to start a chat with.
            </DialogDescription>
            </DialogHeader>
        </div>

        <div className="p-2">
             <Command className="rounded-lg border-0 shadow-none">
                <CommandInput
                    placeholder="Search by username or email..."
                    value={query}
                    onValueChange={setQuery}
                    className="border-0 focus:ring-0 text-[#202022] placeholder:text-[#202022]/40"
                />
                <CommandList className="max-h-[300px] overflow-y-auto custom-scrollbar p-2">
                    {loading && (
                        <div className="py-6 text-center text-sm text-[#202022]/50 flex items-center justify-center gap-2">
                            <div className="animate-spin rounded-full h-4 w-4 border-2 border-[#7678ed] border-t-transparent"></div>
                            Searching...
                        </div>
                    )}

                    {!loading && users.length === 0 && query && (
                        <div className="py-8 text-center">
                            <p className="text-sm text-[#202022]/50">No users found</p>
                        </div>
                    )}

                    {!loading && !query && (
                        <div className="py-8 text-center">
                             <p className="text-sm text-[#202022]/40">Type to find people</p>
                        </div>
                    )}

                    <div className="space-y-1">
                        {!loading && users.map((u) => (
                            <div
                                key={u.id}
                                onClick={() => handleStartChat(u)}
                                className="flex items-center gap-3 p-3 rounded-lg hover:bg-[#7678ed]/5 cursor-pointer transition-colors group"
                            >
                                <Avatar className="h-10 w-10 border border-[#7678ed]/10">
                                    <AvatarFallback className="bg-gradient-to-br from-[#7678ed] to-[#5a5cd9] text-white text-sm font-semibold group-hover:from-[#6567d9] group-hover:to-[#4e50cd]">
                                        {u.username.charAt(0).toUpperCase()}
                                    </AvatarFallback>
                                </Avatar>
                                <div className="flex-1 min-w-0">
                                    <h4 className="text-sm font-medium text-[#202022] truncate group-hover:text-[#7678ed] transition-colors">
                                        {u.username}
                                    </h4>
                                    <p className="text-xs text-[#202022]/50 truncate">
                                        {u.email}
                                    </p>
                                </div>
                                {u.is_online && (
                                     <div className="w-2 h-2 rounded-full bg-green-500 shadow-sm"></div>
                                )}
                            </div>
                        ))}
                    </div>
                </CommandList>
             </Command>
        </div>
      </DialogContent>
    </Dialog>
  );
}
