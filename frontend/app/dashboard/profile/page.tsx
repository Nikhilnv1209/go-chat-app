'use client';

import { ArrowLeft } from 'lucide-react';
import { useRouter } from 'next/navigation';
import UserProfile from '@/components/chat/UserProfile';
import { Button } from '@/components/ui/button';

export default function ProfilePage() {
  const router = useRouter();

  return (
    <div className="min-h-screen bg-[#f9fafc] flex flex-col">
      {/* Mobile Header */}
      <header className="md:hidden flex items-center gap-3 h-14 px-4 bg-white border-b border-[#7678ed]/10 sticky top-0 z-10">
        <Button
          variant="ghost"
          size="icon"
          onClick={() => router.push('/dashboard')}
          className="text-[#202022]/50 hover:text-[#202022] hover:bg-[#7678ed]/10 h-9 w-9"
        >
          <ArrowLeft className="w-5 h-5" />
        </Button>
        <h1 className="text-lg font-semibold text-[#202022]">Profile</h1>
      </header>

      {/* Desktop Header */}
      <header className="hidden md:block border-b border-[#7678ed]/10 bg-white">
        <div className="relative h-14 px-6">
          <Button
            variant="ghost"
            onClick={() => router.push('/dashboard')}
            className="absolute left-6 top-1/2 -translate-y-1/2 text-[#202022]/50 hover:text-[#202022] hover:bg-[#7678ed]/10 h-9 px-4 transition-all"
          >
            <ArrowLeft className="w-4 h-4 mr-2" />
            Back
          </Button>
          <div className="flex items-center justify-center h-full">
            <h1 className="text-lg font-semibold text-[#202022]">Profile</h1>
          </div>
          <div className="w-12"></div> {/* Spacer for centering */}
        </div>
      </header>

      {/* Main Content */}
      <div className="flex-1 overflow-y-auto">
        <UserProfile />
      </div>
    </div>
  );
}
