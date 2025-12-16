import Link from "next/link";
import { Button } from "@/components/ui/button";
import { MessageSquare, Zap, Shield, ArrowRight } from "lucide-react";

export default function Home() {
  return (
    <div className="min-h-[100dvh] relative flex flex-col overflow-x-hidden">
      {/* Fixed Background Layer */}
      <div className="fixed inset-0 h-[100lvh] w-full overflow-hidden pointer-events-none bg-slate-950">
        {/* Animated Background */}
        <div className="absolute inset-0 bg-[radial-gradient(ellipse_80%_80%_at_50%_-20%,rgba(120,119,198,0.3),rgba(255,255,255,0))]" />

        {/* Floating Orbs */}
        <div
          className="absolute top-[-20%] left-[-10%] w-[500px] h-[500px] rounded-full pointer-events-none animate-float"
          style={{
            background: 'radial-gradient(circle, rgba(99,102,241,0.3) 0%, rgba(99,102,241,0) 70%)',
          }}
        />
        <div
          className="absolute bottom-[-20%] right-[-10%] w-[600px] h-[600px] rounded-full pointer-events-none animate-float-reverse"
          style={{
            background: 'radial-gradient(circle, rgba(168,85,247,0.25) 0%, rgba(168,85,247,0) 70%)',
          }}
        />
      </div>

      {/* Content */}
      <div className="relative z-10 flex flex-col items-center justify-center flex-grow min-h-[100dvh] px-4 py-12">
        {/* Hero Section */}
        <div className="text-center max-w-3xl mx-auto animate-enter">
          {/* Badge */}
          <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-white/[0.05] border border-white/[0.1] mb-8">
            <Zap className="w-4 h-4 text-yellow-400" />
            <span className="text-sm text-slate-300">Real-time messaging</span>
          </div>

          {/* Heading */}
          <h1 className="text-4xl sm:text-5xl md:text-7xl font-bold text-white mb-4 sm:mb-6 tracking-tight">
            Chat with{" "}
            <span className="bg-gradient-to-r from-indigo-400 via-purple-400 to-pink-400 bg-clip-text text-transparent">
              anyone
            </span>
            <br />anywhere
          </h1>

          <p className="text-base sm:text-lg md:text-xl text-slate-400 mb-8 sm:mb-12 max-w-xl mx-auto px-4">
            Experience seamless communication with instant messaging, read receipts,
            and typing indicators — all in real-time.
          </p>

          {/* CTA Buttons */}
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Button
              asChild
              size="lg"
              className="h-14 px-8 text-base font-medium bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 text-white border-0 shadow-lg shadow-indigo-500/25 transition-all hover:shadow-indigo-500/40 relative overflow-hidden group"
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
              className="h-14 px-8 text-base font-medium bg-white/[0.05] border-white/[0.1] text-white hover:bg-white/[0.1] hover:text-white transition-all"
            >
              <Link href="/login">Sign In</Link>
            </Button>
          </div>
        </div>

        {/* Features */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mt-12 md:mt-24 max-w-4xl mx-auto">
          <div className="p-6 rounded-2xl bg-white/[0.03] border border-white/[0.05] hover:bg-white/[0.05] transition-colors">
            <div className="w-12 h-12 rounded-xl bg-indigo-500/10 flex items-center justify-center mb-4">
              <MessageSquare className="w-6 h-6 text-indigo-400" />
            </div>
            <h3 className="text-lg font-semibold text-white mb-2">Instant Messages</h3>
            <p className="text-sm text-slate-400">Send and receive messages instantly with WebSocket technology.</p>
          </div>

          <div className="p-6 rounded-2xl bg-white/[0.03] border border-white/[0.05] hover:bg-white/[0.05] transition-colors">
            <div className="w-12 h-12 rounded-xl bg-purple-500/10 flex items-center justify-center mb-4">
              <Zap className="w-6 h-6 text-purple-400" />
            </div>
            <h3 className="text-lg font-semibold text-white mb-2">Real-time Updates</h3>
            <p className="text-sm text-slate-400">See typing indicators and read receipts as they happen.</p>
          </div>

          <div className="p-6 rounded-2xl bg-white/[0.03] border border-white/[0.05] hover:bg-white/[0.05] transition-colors">
            <div className="w-12 h-12 rounded-xl bg-pink-500/10 flex items-center justify-center mb-4">
              <Shield className="w-6 h-6 text-pink-400" />
            </div>
            <h3 className="text-lg font-semibold text-white mb-2">Secure & Private</h3>
            <p className="text-sm text-slate-400">Your conversations are protected with JWT authentication.</p>
          </div>
        </div>
      </div>
    </div>
  );
}
