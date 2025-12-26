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
      {/* Fixed Background Layer - Lavender with wavy pattern effect */}
      <div className="fixed inset-0 h-[100lvh] w-full overflow-hidden pointer-events-none bg-[#d8d8ec]">
        {/* Wavy Pattern Background Layers */}
        <div
          className="absolute inset-0 opacity-30"
          style={{
            backgroundImage: `url("data:image/svg+xml,%3Csvg viewBox='0 0 1200 800' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath fill='%23ff7a55' fill-opacity='0.12' d='M0,192L48,176C96,160,192,128,288,144C384,160,480,224,576,245.3C672,267,768,245,864,213.3C960,181,1056,139,1152,144C1248,149,1344,203,1392,229.3L1440,256L1440,320L1392,320C1344,320,1248,320,1152,320C1056,320,960,320,864,320C768,320,672,320,576,320C480,320,384,320,288,320C192,320,96,320,48,320L0,320Z'%3E%3C/path%3E%3C/svg%3E")`,
            backgroundSize: 'cover',
            backgroundPosition: 'bottom',
          }}
        />

        {/* Floating Orbs - Coral/Orange Theme */}
        <div
          className="absolute top-[-20%] left-[-10%] w-[500px] h-[500px] rounded-full pointer-events-none animate-float"
          style={{
            background: 'radial-gradient(circle, rgba(255,122,85,0.25) 0%, rgba(255,122,85,0) 70%)',
          }}
        />
        <div
          className="absolute bottom-[-20%] right-[-10%] w-[600px] h-[600px] rounded-full pointer-events-none animate-float-reverse"
          style={{
            background: 'radial-gradient(circle, rgba(118,120,237,0.2) 0%, rgba(118,120,237,0) 70%)',
          }}
        />
      </div>

      <div className="relative z-10 flex items-center justify-center flex-grow min-h-[100dvh] px-4 py-8">
      <div className="w-full max-w-md md:max-w-xl">
        {/* White Card */}
        <div className="relative bg-white rounded-2xl p-8 md:p-10 shadow-xl shadow-[#ff7a55]/10 border border-[#ff7a55]/10">
          {/* Subtle glow effect */}
          <div className="absolute inset-0 rounded-2xl bg-gradient-to-r from-[#ff7a55]/5 via-transparent to-[#7678ed]/5 blur-xl -z-10" />

          {/* Header */}
          <div className="text-center mb-6 sm:mb-8">
            <div className="inline-flex items-center justify-center w-14 h-14 sm:w-16 sm:h-16 rounded-2xl bg-gradient-to-br from-[#ff7a55] to-[#e66a47] mb-4 sm:mb-6 shadow-lg shadow-[#ff7a55]/30">
              <Rocket className="w-7 h-7 sm:w-8 sm:h-8 text-white" />
            </div>
            <h1 className="text-2xl sm:text-3xl font-bold text-[#202022] mb-2">Get started</h1>
            <p className="text-sm sm:text-base text-[#202022]/60">Create your account in seconds</p>
          </div>

          {/* Form */}
          <form onSubmit={handleSubmit} className="space-y-5">
            <div className="space-y-3">
              <label htmlFor="username" className="block text-sm font-medium text-[#202022]">Username</label>
              <div className="relative group">
                <User className="absolute left-4 top-1/2 -translate-y-1/2 h-5 w-5 text-[#202022]/40 group-focus-within:text-[#ff7a55] transition-colors" />
                <Input
                  id="username"
                  type="text"
                  placeholder="johndoe"
                  className={`h-12 pl-12 bg-[#f9fafc] border-[#ff7a55]/20 text-[#202022] placeholder:text-[#202022]/40 focus:border-[#ff7a55] focus:ring-[#ff7a55]/20 transition-all rounded-xl ${errors.username ? 'border-[#ff7a55]' : ''}`}
                  value={username}
                  onChange={(e) => { setUsername(e.target.value); setErrors(prev => ({ ...prev, username: undefined })); }}
                />
              </div>
              {errors.username && (
                <p className="text-xs text-[#ff7a55] mt-1 animate-enter">{errors.username}</p>
              )}
            </div>

            <div className="space-y-3">
              <label htmlFor="email" className="block text-sm font-medium text-[#202022]">Email</label>
              <div className="relative group">
                <Mail className="absolute left-4 top-1/2 -translate-y-1/2 h-5 w-5 text-[#202022]/40 group-focus-within:text-[#ff7a55] transition-colors" />
                <Input
                  id="email"
                  type="email"
                  placeholder="you@example.com"
                  className={`h-12 pl-12 bg-[#f9fafc] border-[#ff7a55]/20 text-[#202022] placeholder:text-[#202022]/40 focus:border-[#ff7a55] focus:ring-[#ff7a55]/20 transition-all rounded-xl ${errors.email ? 'border-[#ff7a55]' : ''}`}
                  value={email}
                  onChange={(e) => { setEmail(e.target.value); setErrors(prev => ({ ...prev, email: undefined })); }}
                />
              </div>
              {errors.email && (
                <p className="text-xs text-[#ff7a55] mt-1 animate-enter">{errors.email}</p>
              )}
            </div>

            <div className="space-y-3">
              <label htmlFor="password" className="block text-sm font-medium text-[#202022]">Password</label>
              <div className="relative group">
                <Lock className="absolute left-4 top-1/2 -translate-y-1/2 h-5 w-5 text-[#202022]/40 group-focus-within:text-[#ff7a55] transition-colors" />
                <Input
                  id="password"
                  type="password"
                  placeholder="••••••••"
                  className={`h-12 pl-12 bg-[#f9fafc] border-[#ff7a55]/20 text-[#202022] placeholder:text-[#202022]/40 focus:border-[#ff7a55] focus:ring-[#ff7a55]/20 transition-all rounded-xl ${errors.password ? 'border-[#ff7a55]' : ''}`}
                  value={password}
                  onChange={(e) => { setPassword(e.target.value); setErrors(prev => ({ ...prev, password: undefined })); }}
                />
              </div>
              {errors.password && (
                <p className="text-xs text-[#ff7a55] mt-1 animate-enter">{errors.password}</p>
              )}
            </div>

            {registerMutation.isError && (
              <div className="p-4 text-sm text-[#ff7a55] bg-[#ff7a55]/10 border border-[#ff7a55]/20 rounded-xl animate-enter">
                {(() => {
                  const axiosError = registerMutation.error as AxiosError<{ error: { code: string; message: string } }>;
                  return axiosError.response?.data?.error?.message || 'Failed to create account';
                })()}
              </div>
            )}

            <Button
              className="w-full h-12 text-base font-medium bg-[#ff7a55] hover:bg-[#e66a47] text-white border-0 shadow-lg shadow-[#ff7a55]/25 transition-all hover:shadow-[#ff7a55]/40 rounded-xl relative overflow-hidden group"
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
            <p className="text-sm text-[#202022]/60">
              Already have an account?{' '}
              <Link href="/login" className="text-[#7678ed] font-semibold hover:text-[#5a5cd9] transition-colors">
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
