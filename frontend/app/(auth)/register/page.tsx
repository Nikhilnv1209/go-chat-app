'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { useMutation } from '@tanstack/react-query';
import { useDispatch } from 'react-redux';
import { AxiosError } from 'axios';
import { User, Mail, Lock, ArrowRight, Rocket } from 'lucide-react';

import api from '@/lib/api';
import { setCredentials } from '@/store/features/authSlice';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';

export default function RegisterPage() {
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [errors, setErrors] = useState<{ username?: string; email?: string; password?: string }>({});
  const router = useRouter();
  const dispatch = useDispatch();

  const registerMutation = useMutation({
    mutationFn: async () => {
      const response = await api.post('/auth/register', { username, email, password });
      return response.data;
    },
    onSuccess: (data) => {
      dispatch(setCredentials({ user: data.user, token: data.token }));
      router.push('/dashboard');
    },
  });

  const validateForm = () => {
    const newErrors: { username?: string; email?: string; password?: string } = {};

    if (!username) {
      newErrors.username = 'Username is required';
    } else if (username.length < 3) {
      newErrors.username = 'Username must be at least 3 characters';
    }

    if (!email) {
      newErrors.email = 'Email is required';
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
      newErrors.email = 'Please enter a valid email address';
    }

    if (!password) {
      newErrors.password = 'Password is required';
    } else if (password.length < 6) {
      newErrors.password = 'Password must be at least 6 characters';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (validateForm()) {
      registerMutation.mutate();
    }
  };

  return (
    <div className="min-h-[100dvh] relative flex flex-col overflow-x-hidden">
      {/* Fixed Background Layer */}
      <div className="fixed inset-0 h-[100lvh] w-full overflow-hidden pointer-events-none bg-slate-950">
        {/* Animated Background */}
        <div className="absolute inset-0 bg-[radial-gradient(ellipse_80%_80%_at_50%_-20%,rgba(120,119,198,0.3),rgba(255,255,255,0))]" />

        {/* Floating Orbs - Purple/Pink Theme */}
        <div
          className="absolute top-[-20%] left-[-10%] w-[500px] h-[500px] rounded-full pointer-events-none animate-float"
          style={{
            background: 'radial-gradient(circle, rgba(168,85,247,0.3) 0%, rgba(168,85,247,0) 70%)',
          }}
        />
        <div
          className="absolute bottom-[-20%] right-[-10%] w-[600px] h-[600px] rounded-full pointer-events-none animate-float-reverse"
          style={{
            background: 'radial-gradient(circle, rgba(236,72,153,0.25) 0%, rgba(236,72,153,0) 70%)',
          }}
        />
      </div>

      <div className="relative z-10 flex items-center justify-center flex-grow min-h-[100dvh] px-4 py-8">
      <div className="w-full max-w-md md:max-w-xl">
        {/* Glass Card */}
        <div className="relative backdrop-blur-xl bg-white/[0.05] border border-white/[0.1] rounded-2xl p-8 md:p-10 shadow-2xl">
          {/* Glow effect */}
          <div className="absolute inset-0 rounded-2xl bg-gradient-to-r from-purple-500/10 via-pink-500/10 to-orange-500/10 blur-xl -z-10" />

          {/* Header */}
          <div className="text-center mb-6 sm:mb-8">
            <div className="inline-flex items-center justify-center w-14 h-14 sm:w-16 sm:h-16 rounded-2xl bg-gradient-to-br from-purple-500 to-pink-600 mb-4 sm:mb-6 shadow-lg shadow-purple-500/25">
              <Rocket className="w-7 h-7 sm:w-8 sm:h-8 text-white" />
            </div>
            <h1 className="text-2xl sm:text-3xl font-bold text-white mb-2">Get started</h1>
            <p className="text-sm sm:text-base text-slate-400">Create your account in seconds</p>
          </div>

          {/* Form */}
          <form onSubmit={handleSubmit} className="space-y-5">
            <div className="space-y-3">
              <label htmlFor="username" className="block text-sm font-medium text-slate-300">Username</label>
              <div className="relative group">
                <User className="absolute left-4 top-1/2 -translate-y-1/2 h-5 w-5 text-slate-500 group-focus-within:text-purple-400 transition-colors" />
                <Input
                  id="username"
                  type="text"
                  placeholder="johndoe"
                  className={`h-12 pl-12 bg-white/[0.05] border-white/[0.1] text-white placeholder:text-slate-500 focus:border-purple-500/50 focus:ring-purple-500/20 transition-all ${errors.username ? 'border-red-500/50' : ''}`}
                  value={username}
                  onChange={(e) => { setUsername(e.target.value); setErrors(prev => ({ ...prev, username: undefined })); }}
                />
              </div>
              {errors.username && (
                <p className="text-xs text-red-400 mt-1 animate-enter">{errors.username}</p>
              )}
            </div>

            <div className="space-y-3">
              <label htmlFor="email" className="block text-sm font-medium text-slate-300">Email</label>
              <div className="relative group">
                <Mail className="absolute left-4 top-1/2 -translate-y-1/2 h-5 w-5 text-slate-500 group-focus-within:text-purple-400 transition-colors" />
                <Input
                  id="email"
                  type="email"
                  placeholder="you@example.com"
                  className={`h-12 pl-12 bg-white/[0.05] border-white/[0.1] text-white placeholder:text-slate-500 focus:border-purple-500/50 focus:ring-purple-500/20 transition-all ${errors.email ? 'border-red-500/50' : ''}`}
                  value={email}
                  onChange={(e) => { setEmail(e.target.value); setErrors(prev => ({ ...prev, email: undefined })); }}
                />
              </div>
              {errors.email && (
                <p className="text-xs text-red-400 mt-1 animate-enter">{errors.email}</p>
              )}
            </div>

            <div className="space-y-3">
              <label htmlFor="password" className="block text-sm font-medium text-slate-300">Password</label>
              <div className="relative group">
                <Lock className="absolute left-4 top-1/2 -translate-y-1/2 h-5 w-5 text-slate-500 group-focus-within:text-purple-400 transition-colors" />
                <Input
                  id="password"
                  type="password"
                  placeholder="••••••••"
                  className={`h-12 pl-12 bg-white/[0.05] border-white/[0.1] text-white placeholder:text-slate-500 focus:border-purple-500/50 focus:ring-purple-500/20 transition-all ${errors.password ? 'border-red-500/50' : ''}`}
                  value={password}
                  onChange={(e) => { setPassword(e.target.value); setErrors(prev => ({ ...prev, password: undefined })); }}
                />
              </div>
              {errors.password && (
                <p className="text-xs text-red-400 mt-1 animate-enter">{errors.password}</p>
              )}
            </div>

            {registerMutation.isError && (
              <div className="p-4 text-sm text-red-300 bg-red-500/10 border border-red-500/20 rounded-xl animate-enter">
                {(() => {
                  const axiosError = registerMutation.error as AxiosError<{ error: { code: string; message: string } }>;
                  return axiosError.response?.data?.error?.message || 'Failed to create account';
                })()}
              </div>
            )}

            <Button
              className="w-full h-12 text-base font-medium bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-700 hover:to-pink-700 text-white border-0 shadow-lg shadow-purple-500/25 transition-all hover:shadow-purple-500/40 relative overflow-hidden group"
              type="submit"
              disabled={registerMutation.isPending}
            >
              {registerMutation.isPending ? (
                <span className="flex items-center gap-2">
                  <div className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                  Creating account...
                </span>
              ) : (
                <span className="flex items-center gap-2">
                  Create Account
                  <ArrowRight className="w-4 h-4" />
                </span>
              )}
            </Button>
          </form>

          {/* Footer */}
          <div className="mt-8 text-center">
            <p className="text-sm text-slate-400">
              Already have an account?{' '}
              <Link href="/login" className="text-purple-400 font-semibold hover:text-purple-300 transition-colors">
                Sign in
              </Link>
            </p>
          </div>
        </div>
      </div>
      </div>
    </div>
  );
}
