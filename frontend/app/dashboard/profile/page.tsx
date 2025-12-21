'use client';

import { ArrowLeft } from 'lucide-react';
import { useRouter } from 'next/navigation';
import UserProfile from '@/components/chat/UserProfile';
import { Button } from '@/components/ui/button';

export default function ProfilePage() {
  const router = useRouter();

  return (
    <div className="h-screen bg-slate-950 flex flex-col ios-safe-area no-bounce">
      {/* Header - Hidden on mobile since mobile header is handled by dashboard layout */}
      <header className="hidden md:block border-b border-white/10 bg-slate-950/80 backdrop-blur-sm">
        <div className="relative h-14 px-6">
          <Button
            variant="ghost"
            onClick={() => router.push('/dashboard')}
            className="absolute left-6 top-1/2 -translate-y-1/2 text-slate-400 hover:text-white hover:bg-white/5 h-9 px-4 transition-all"
          >
            <ArrowLeft className="w-4 h-4 mr-2" />
            Back
          </Button>
          <div className="flex items-center justify-center h-full">
            <h1 className="text-lg font-semibold text-white">Profile</h1>
          </div>
          <div className="w-12"></div> {/* Spacer for centering */}
        </div>
      </header>

      {/* Main Content */}
      <div className="flex-1 overflow-hidden">
        <UserProfile />
      </div>
    </div>
  );
}