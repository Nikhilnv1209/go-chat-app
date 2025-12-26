'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAppSelector } from '@/store/hooks';
import Link from 'next/link';
import { Button } from '@/components/ui/button';
import { ArrowRight, MessageSquare, ShieldCheck, Zap, Shield } from 'lucide-react';

export default function Home() {
  const { isAuthenticated, isLoading } = useAppSelector((state) => state.auth);
  const router = useRouter();

  useEffect(() => {
    if (!isLoading && isAuthenticated) {
      router.push('/dashboard');
    }
  }, [isLoading, isAuthenticated, router]);

  if (isLoading) {
    return (
      <div className="min-h-[100dvh] relative flex flex-col overflow-x-hidden">
        <div className="fixed inset-0 h-[100lvh] w-full overflow-hidden pointer-events-none bg-[#d8d8ec]"></div>
        <div className="relative z-10 flex items-center justify-center flex-grow min-h-[100dvh]">
          <div className="animate-spin rounded-full h-12 w-12 border-2 border-[#7678ed] border-t-transparent"></div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-[100dvh] relative flex flex-col overflow-x-hidden">
      {/* Fixed Background Layer */}
      <div className="fixed inset-0 h-[100lvh] w-full overflow-hidden pointer-events-none bg-[#d8d8ec]">
        {/* Wavy Pattern Background */}
        <div
          className="absolute inset-0 opacity-40"
          style={{
            backgroundImage: `url("data:image/svg+xml,%3Csvg viewBox='0 0 1200 800' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath fill='%237678ed' fill-opacity='0.12' d='M0,192L48,176C96,160,192,128,288,144C384,160,480,224,576,245.3C672,267,768,245,864,213.3C960,181,1056,139,1152,144C1248,149,1344,203,1392,229.3L1440,256L1440,320L1392,320C1344,320,1248,320,1152,320C1056,320,960,320,864,320C768,320,672,320,576,320C480,320,384,320,288,320C192,320,96,320,48,320L0,320Z'%3E%3C/path%3E%3C/svg%3E")`,
            backgroundSize: 'cover',
            backgroundPosition: 'bottom',
          }}
        />

        {/* Floating Orbs */}
        <div
          className="absolute top-[-20%] left-[-10%] w-[500px] h-[500px] rounded-full pointer-events-none animate-float"
          style={{
            background: 'radial-gradient(circle, rgba(118,120,237,0.25) 0%, rgba(118,120,237,0) 70%)',
          }}
        />
        <div
          className="absolute bottom-[-20%] right-[-10%] w-[600px] h-[600px] rounded-full pointer-events-none animate-float-reverse"
          style={{
            background: 'radial-gradient(circle, rgba(255,122,85,0.2) 0%, rgba(255,122,85,0) 70%)',
          }}
        />
      </div>

      {/* Content */}
      <div className="relative z-10 flex flex-col items-center justify-center flex-grow min-h-[100dvh] px-4 py-12">
        {/* Hero Section */}
        <div className="text-center max-w-3xl mx-auto animate-enter">
          {/* Badge */}
          <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-white shadow-sm border border-[#7678ed]/10 mb-8">
            <Zap className="w-4 h-4 text-[#ff7a55]" />
            <span className="text-sm text-[#202022]">Real-time messaging</span>
          </div>

          {/* Heading */}
          <h1 className="text-4xl sm:text-5xl md:text-7xl font-bold text-[#202022] mb-4 sm:mb-6 tracking-tight">
            Chat with{" "}
            <span className="bg-gradient-to-r from-[#7678ed] via-[#7678ed] to-[#ff7a55] bg-clip-text text-transparent">
              anyone
            </span>
            <br />anywhere
          </h1>

          <p className="text-base sm:text-lg md:text-xl text-[#202022]/60 mb-8 sm:mb-12 max-w-xl mx-auto px-4">
            Experience seamless communication with instant messaging, read receipts,
            and typing indicators — all in real-time.
          </p>

          {/* CTA Buttons */}
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Button
              asChild
              size="lg"
              className="h-14 px-8 text-base font-medium bg-[#7678ed] hover:bg-[#5a5cd9] text-white border-0 shadow-lg shadow-[#7678ed]/25 transition-all hover:shadow-[#7678ed]/40 relative overflow-hidden group"
            >
              <Link href="/register" className="flex items-center gap-2">
                Get Started Free
                <ArrowRight className="w-4 h-4" />
              </Link>
            </Button>
            <Button
              asChild
              variant="outline"
              size="lg"
              className="h-14 px-8 text-base font-medium bg-white border-[#7678ed]/20 text-[#202022] hover:bg-[#7678ed]/5 hover:border-[#7678ed]/30 transition-all"
            >
              <Link href="/login">Sign In</Link>
            </Button>
          </div>
        </div>

        {/* Features */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mt-12 md:mt-24 max-w-4xl mx-auto">
          <div className="p-6 rounded-2xl bg-white border border-[#7678ed]/10 hover:border-[#7678ed]/20 transition-colors shadow-sm">
            <div className="w-12 h-12 rounded-xl bg-[#7678ed]/10 flex items-center justify-center mb-4">
              <MessageSquare className="w-6 h-6 text-[#7678ed]" />
            </div>
            <h3 className="text-lg font-semibold text-[#202022] mb-2">Instant Messages</h3>
            <p className="text-sm text-[#202022]/60">Send and receive messages instantly with WebSocket technology.</p>
          </div>

          <div className="p-6 rounded-2xl bg-white border border-[#7678ed]/10 hover:border-[#7678ed]/20 transition-colors shadow-sm">
            <div className="w-12 h-12 rounded-xl bg-[#ff7a55]/10 flex items-center justify-center mb-4">
              <Zap className="w-6 h-6 text-[#ff7a55]" />
            </div>
            <h3 className="text-lg font-semibold text-[#202022] mb-2">Real-time Updates</h3>
            <p className="text-sm text-[#202022]/60">See typing indicators and read receipts as they happen.</p>
          </div>

          <div className="p-6 rounded-2xl bg-white border border-[#7678ed]/10 hover:border-[#7678ed]/20 transition-colors shadow-sm">
            <div className="w-12 h-12 rounded-xl bg-[#7678ed]/10 flex items-center justify-center mb-4">
              <Shield className="w-6 h-6 text-[#7678ed]" />
            </div>
            <h3 className="text-lg font-semibold text-[#202022] mb-2">Secure & Private</h3>
            <p className="text-sm text-[#202022]/60">Your conversations are protected with JWT authentication.</p>
          </div>
        </div>
      </div>
    </div>
  );
}
