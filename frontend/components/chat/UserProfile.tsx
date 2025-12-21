'use client';

import { useRouter } from 'next/navigation';
import { User, Mail, Calendar, LogOut, Settings } from 'lucide-react';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import { logout } from '@/store/features/authSlice';
import { Button } from '@/components/ui/button';

export default function UserProfile() {
  const router = useRouter();
  const dispatch = useAppDispatch();
  const { user } = useAppSelector((state) => state.auth);

  const handleLogout = () => {
    dispatch(logout());
    router.push('/login');
  };

  const handleEditProfile = () => {
    // Navigate to edit profile page or show edit modal
    // For now, show a toast or alert that this feature is coming soon
    alert('Edit Profile feature coming soon!');
  };

  const handleSettings = () => {
    // Navigate to settings page or show settings modal
    // For now, show a toast or alert that this feature is coming soon
    alert('Settings page coming soon!');
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  };

  const formatLastSeen = (dateString: string) => {
    const date = new Date(dateString);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffMins = Math.floor(diffMs / 60000);
    const diffHours = Math.floor(diffMins / 60);
    const diffDays = Math.floor(diffHours / 24);

    if (diffMins < 5) return 'Active now';
    if (diffMins < 60) return `Last seen ${diffMins} minutes ago`;
    if (diffHours < 24) return `Last seen ${diffHours} hours ago`;
    if (diffDays < 7) return `Last seen ${diffDays} days ago`;
    return `Last seen ${formatDate(dateString)}`;
  };

  if (!user) {
    return (
      <div className="flex items-center justify-center p-8">
        <div className="animate-spin rounded-full h-8 w-8 border-2 border-indigo-500 border-t-transparent"></div>
      </div>
    );
  }

  return (
    <div className="flex flex-col h-full bg-slate-950">
      <div className="flex-1 overflow-y-auto scrollbar-modern">
        {/* Profile Info */}
        <div className="max-w-4xl mx-auto px-6 py-6 space-y-6">
        {/* Avatar Section */}
        <div className="flex flex-col items-center space-y-3 pb-4">
          <div className="relative">
            <div className="w-24 h-24 rounded-full bg-linear-to-br from-indigo-500 to-purple-600 flex items-center justify-center text-white font-bold text-3xl shadow-lg">
              {user.username.charAt(0).toUpperCase()}
            </div>
            <div className={`absolute bottom-1 right-1 w-6 h-6 border-2 border-slate-900 rounded-full shadow-md ${
              user.is_online ? 'bg-green-500' : 'bg-slate-600'
            }`}></div>
          </div>
          <div className="text-center">
            <h3 className="text-2xl font-bold text-white mb-1">{user.username}</h3>
            <p className={`text-sm ${
              user.is_online ? 'text-green-400' : 'text-slate-400'
            }`}>
              {user.is_online ? 'Online' : formatLastSeen(user.last_seen)}
            </p>
          </div>
        </div>

        {/* User Details */}
        <div className="space-y-3 pb-6">
          <div className="bg-white/3 border border-white/5 rounded-lg p-4 transition-all hover:bg-white/5">
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-md bg-indigo-500/20 flex items-center justify-center">
                <User className="w-5 h-5 text-indigo-400" />
              </div>
              <div className="flex-1">
                <p className="text-xs font-medium text-slate-400 mb-1">Username</p>
                <p className="text-base text-white font-medium">{user.username}</p>
              </div>
            </div>
          </div>

          <div className="bg-white/3 border border-white/5 rounded-lg p-4 transition-all hover:bg-white/5">
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-md bg-purple-500/20 flex items-center justify-center">
                <Mail className="w-5 h-5 text-purple-400" />
              </div>
              <div className="flex-1">
                <p className="text-xs font-medium text-slate-400 mb-1">Email Address</p>
                <p className="text-base text-white font-medium">{user.email}</p>
              </div>
            </div>
          </div>

          <div className="bg-white/3 border border-white/5 rounded-lg p-4 transition-all hover:bg-white/5">
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-md bg-blue-500/20 flex items-center justify-center">
                <Calendar className="w-5 h-5 text-blue-400" />
              </div>
              <div className="flex-1">
                <p className="text-xs font-medium text-slate-400 mb-1">Member Since</p>
                <p className="text-base text-white font-medium">{formatDate(user.created_at)}</p>
              </div>
            </div>
          </div>

          <div className="bg-white/3 border border-white/5 rounded-lg p-4 transition-all hover:bg-white/5">
            <div className="flex items-center gap-3">
              <div className={`w-10 h-10 rounded-md flex items-center justify-center ${
                user.is_online ? 'bg-green-500/20' : 'bg-slate-600/20'
              }`}>
                <div className={`w-5 h-5 rounded-full ${
                  user.is_online ? 'bg-green-500' : 'bg-slate-600'
                }`}></div>
              </div>
              <div className="flex-1">
                <p className="text-xs font-medium text-slate-400 mb-1">Status</p>
                <p className="text-base text-white font-medium">
                  {user.is_online ? 'Available' : 'Offline'}
                </p>
              </div>
            </div>
          </div>
        </div>

        </div>

        {/* Actions Section */}
        <div className="space-y-3 pb-8">
          <div className="flex flex-row gap-3 sm:gap-4 justify-center items-center">
            <Button
              variant="outline"
              onClick={handleEditProfile}
              className="w-32 h-10 px-4 sm:w-auto sm:h-11 sm:px-8 sm:min-w-32 text-sm sm:text-base bg-white/5 border-white/10 text-white hover:bg-white/10 hover:border-white/20 transition-all"
            >
              <Settings className="w-4 h-4 mr-2" />
              <span className="sm:hidden">Edit</span>
              <span className="hidden sm:inline">Edit Profile</span>
            </Button>

            <Button
              variant="destructive"
              onClick={handleLogout}
              className="w-32 h-10 px-4 sm:w-auto sm:h-11 sm:px-8 sm:min-w-32 text-sm sm:text-base bg-red-500/10 border border-red-500/20 text-red-400 hover:bg-red-500/20 hover:border-red-500/30 hover:text-red-300 transition-all"
            >
              <LogOut className="w-4 h-4 mr-2" />
              <span className="sm:hidden">Logout</span>
              <span className="hidden sm:inline">Logout</span>
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}