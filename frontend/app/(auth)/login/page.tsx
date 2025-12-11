'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { useMutation } from '@tanstack/react-query';
import { useDispatch } from 'react-redux';
import { AxiosError } from 'axios';
import { Mail, Lock, ArrowRight, Sparkles } from 'lucide-react';

import api from '@/lib/api';
import { setCredentials } from '@/store/features/authSlice';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';

export default function LoginPage() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [errors, setErrors] = useState<{ email?: string; password?: string }>({});
  const router = useRouter();
  const dispatch = useDispatch();

  const loginMutation = useMutation({
    mutationFn: async () => {
      const response = await api.post('/auth/login', { email, password });
      return response.data;
    },
    onSuccess: (data) => {
      dispatch(setCredentials({ user: data.user, token: data.token }));
      router.push('/');
    },
  });

  const validateForm = () => {
    const newErrors: { email?: string; password?: string } = {};

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
      loginMutation.mutate();
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center px-4 py-8">
      <div className="w-full max-w-md md:max-w-xl">
        {/* Glass Card */}
        <div className="relative backdrop-blur-xl bg-white/[0.05] border border-white/[0.1] rounded-2xl p-8 md:p-10 shadow-2xl">
          {/* Glow effect */}
          <div className="absolute inset-0 rounded-2xl bg-gradient-to-r from-indigo-500/10 via-purple-500/10 to-pink-500/10 blur-xl -z-10" />

          {/* Header */}
          <div className="text-center mb-6 sm:mb-8">
            <div className="inline-flex items-center justify-center w-14 h-14 sm:w-16 sm:h-16 rounded-2xl bg-gradient-to-br from-indigo-500 to-purple-600 mb-4 sm:mb-6 shadow-lg shadow-indigo-500/25">
              <Sparkles className="w-7 h-7 sm:w-8 sm:h-8 text-white" />
            </div>
            <h1 className="text-2xl sm:text-3xl font-bold text-white mb-2">Welcome back</h1>
            <p className="text-sm sm:text-base text-slate-400">Sign in to continue to GoChat</p>
          </div>

          {/* Form */}
          <form onSubmit={handleSubmit} className="space-y-5">
            <div className="space-y-3">
              <label htmlFor="email" className="block text-sm font-medium text-slate-300">Email</label>
              <div className="relative group">
                <Mail className="absolute left-4 top-1/2 -translate-y-1/2 h-5 w-5 text-slate-500 group-focus-within:text-indigo-400 transition-colors" />
                <Input
                  id="email"
                  type="email"
                  placeholder="you@example.com"
                  className={`h-12 pl-12 bg-white/[0.05] border-white/[0.1] text-white placeholder:text-slate-500 focus:border-indigo-500/50 focus:ring-indigo-500/20 transition-all ${errors.email ? 'border-red-500/50' : ''}`}
                  value={email}
                  onChange={(e) => { setEmail(e.target.value); setErrors(prev => ({ ...prev, email: undefined })); }}
                />
              </div>
              {errors.email && (
                <p className="text-xs text-red-400 mt-1 animate-enter">{errors.email}</p>
              )}
            </div>

            <div className="space-y-3">
              <div className="flex items-center justify-between">
                <label htmlFor="password" className="block text-sm font-medium text-slate-300">Password</label>
                <span className="text-xs text-slate-500 cursor-not-allowed">Forgot password?</span>
              </div>
              <div className="relative group">
                <Lock className="absolute left-4 top-1/2 -translate-y-1/2 h-5 w-5 text-slate-500 group-focus-within:text-indigo-400 transition-colors" />
                <Input
                  id="password"
                  type="password"
                  placeholder="••••••••"
                  className={`h-12 pl-12 bg-white/[0.05] border-white/[0.1] text-white placeholder:text-slate-500 focus:border-indigo-500/50 focus:ring-indigo-500/20 transition-all ${errors.password ? 'border-red-500/50' : ''}`}
                  value={password}
                  onChange={(e) => { setPassword(e.target.value); setErrors(prev => ({ ...prev, password: undefined })); }}
                />
              </div>
              {errors.password && (
                <p className="text-xs text-red-400 mt-1 animate-enter">{errors.password}</p>
              )}
            </div>

            {loginMutation.isError && (
              <div className="p-4 text-sm text-red-300 bg-red-500/10 border border-red-500/20 rounded-xl animate-enter">
                {(loginMutation.error as AxiosError<{ error: string }>).response?.data?.error || 'Authentication failed'}
              </div>
            )}

            <Button
              className="w-full h-12 text-base font-medium bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 text-white border-0 shadow-lg shadow-indigo-500/25 transition-all hover:shadow-indigo-500/40 relative overflow-hidden group"
              type="submit"
              disabled={loginMutation.isPending}
            >
              {loginMutation.isPending ? (
                <span className="flex items-center gap-2">
                  <div className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                  Signing in...
                </span>
              ) : (
                <span className="flex items-center gap-2">
                  Sign In
                  <ArrowRight className="w-4 h-4" />
                </span>
              )}
            </Button>
          </form>

          {/* Footer */}
          <div className="mt-8 text-center">
            <p className="text-sm text-slate-400">
              Don&apos;t have an account?{' '}
              <Link href="/register" className="text-indigo-400 font-semibold hover:text-indigo-300 transition-colors">
                Create one
              </Link>
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
