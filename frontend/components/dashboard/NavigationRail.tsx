'use client';

import { usePathname, useRouter } from 'next/navigation';
import { MessageSquare, LogOut, Briefcase, Users2, Archive, User } from 'lucide-react';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import { logout } from '@/store/features/authSlice';
import { toggleSidebar, setSidebarCollapsed } from '@/store/features/uiSlice';
import { cn } from '@/lib/utils';
import { Button } from '@/components/ui/button';
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip';

export default function NavigationRail() {
  const router = useRouter();
  const pathname = usePathname();
  const dispatch = useAppDispatch();
  const { user } = useAppSelector((state) => state.auth);
  const { isSidebarCollapsed } = useAppSelector((state) => state.ui);

  const navItems = [
    { icon: MessageSquare, label: 'All chats', path: '/dashboard', activePath: '/dashboard', count: null },
    { icon: Briefcase, label: 'Work', path: '/dashboard/work', activePath: '/dashboard/work', count: null },
    { icon: Users2, label: 'Friends', path: '/dashboard/friends', activePath: '/dashboard/friends', count: null },
    { icon: Archive, label: 'Archive', path: '/dashboard/archive', activePath: '/dashboard/archive', count: null },
  ];

  const handleLogout = () => {
    dispatch(logout());
    router.push('/login');
  };

  return (
    <div className="hidden md:flex flex-col items-center w-[72px] py-6 bg-[#202022] h-full flex-shrink-0 z-50">
       {/* Logo */}
       <div className="mb-8">
         <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-[#7678ed] to-[#5a5cd9] flex items-center justify-center shadow-lg shadow-[#7678ed]/30">
           <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" className="text-white transform scale-75">
             <path d="M12 2L2 22H22L12 2Z" fill="currentColor" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
           </svg>
         </div>
       </div>

       {/* Nav Items */}
       <div className="flex-1 flex flex-col gap-3 w-full px-3">
         <TooltipProvider delayDuration={0}>
           {navItems.map((item) => {
             const isActive = pathname === item.path || (item.path !== '/dashboard' && pathname.startsWith(item.path));

             return (
               <Tooltip key={item.path}>
                 <TooltipTrigger asChild>
                   <button
                     onClick={() => {
                        if (isActive) {
                            dispatch(toggleSidebar());
                        } else {
                            if (isSidebarCollapsed) dispatch(setSidebarCollapsed(false));
                            router.push(item.path);
                        }
                     }}
                     className={cn(
                       "relative group flex flex-col items-center justify-center w-full aspect-square rounded-xl transition-all duration-200",
                       isActive
                         ? "bg-[#7678ed] text-white shadow-lg shadow-[#7678ed]/30"
                         : "text-white/60 hover:text-white hover:bg-white/10"
                     )}
                   >
                     <item.icon className={cn("w-5 h-5", isActive ? "text-white" : "text-current")} />

                     {/* Unread Badge - Mockup for now */}
                     {item.count && (
                       <span className="absolute top-2 right-2 w-2 h-2 rounded-full bg-[#ff7a55] ring-2 ring-[#202022]" />
                     )}
                   </button>
                 </TooltipTrigger>
                 <TooltipContent side="right" className="bg-[#202022] border-white/10 text-white">
                   <p>{item.label}</p>
                 </TooltipContent>
               </Tooltip>
             );
           })}
         </TooltipProvider>
       </div>

       {/* Bottom Actions */}
       <div className="flex flex-col gap-3 items-center w-full px-3 mt-auto">
         <TooltipProvider delayDuration={0}>
            {/* Profile */}
            <Tooltip>
                <TooltipTrigger asChild>
                    <button
                    onClick={() => router.push('/dashboard/profile')}
                    className={cn(
                        "relative group flex flex-col items-center justify-center w-full aspect-square rounded-xl transition-all duration-200",
                        pathname.includes('profile')
                            ? "bg-white/10 text-white"
                            : "text-white/60 hover:text-white hover:bg-white/10"
                    )}
                    >
                    <div className="w-8 h-8 rounded-full bg-gradient-to-br from-[#ff7a55] to-[#e66a47] flex items-center justify-center text-xs text-white font-bold ring-2 ring-[#202022] group-hover:ring-white/20 transition-all">
                        {user?.username?.[0]?.toUpperCase() || <User className="w-4 h-4" />}
                    </div>
                    </button>
                </TooltipTrigger>
                <TooltipContent side="right" className="bg-[#202022] border-white/10 text-white">
                  <p>Profile</p>
                </TooltipContent>
            </Tooltip>

            {/* Logout */}
            <Tooltip>
                <TooltipTrigger asChild>
                    <button
                    onClick={handleLogout}
                    className="group flex flex-col items-center justify-center w-full aspect-square rounded-xl text-white/60 hover:text-[#ff7a55] hover:bg-[#ff7a55]/10 transition-all duration-200"
                    >
                    <LogOut className="w-4 h-4" />
                    </button>
                </TooltipTrigger>
                <TooltipContent side="right" className="bg-[#202022] border-white/10 text-white">
                  <p>Logout</p>
                </TooltipContent>
            </Tooltip>
         </TooltipProvider>
       </div>
    </div>
  );
}
