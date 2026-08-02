"use client";

import { SidebarTrigger } from "@/components/ui/sidebar";
import { Button } from "@/components/ui/button";
import { Bell, Search } from "@mynaui/icons-react";

interface DashboardHeaderProps {
  title: string;
  subTitle?: string;
  showSearch?: boolean;
  showNotifications?: boolean;
  unreadCount?: number;
  isNotificationSheetOpen?: boolean;
  setIsNotificationSheetOpen?: (open: boolean) => void;
}

export function Header({
  title,
  subTitle,
  showSearch = false,
  showNotifications = true,
  unreadCount = 0,
  isNotificationSheetOpen = false,
  setIsNotificationSheetOpen,
}: DashboardHeaderProps) {
  return (
    <header className="top-0 z-50 bg-white/90 px-4 py-3 backdrop-blur-sm sm:px-6 lg:px-8">
      <div className="flex items-center justify-between gap-2">
        {/* Left Section */}
        <div className="flex min-w-0 flex-1 items-center gap-2 sm:gap-3 md:gap-4">
          <SidebarTrigger className="shrink-0 xl:hidden lg:block" />

          <div className="flex min-w-0 flex-1 flex-col gap-0.5 sm:gap-1">
            <h2 className="truncate text-base font-semibold sm:text-lg md:text-xl lg:text-2xl">
              {title}
            </h2>

            {subTitle && (
              <p className="hidden truncate text-xs text-muted-foreground sm:block sm:text-sm">
                {subTitle}
              </p>
            )}
          </div>
        </div>

        {/* Right Section */}
        {showNotifications && (
          <div className="flex shrink-0 items-center gap-1.5 sm:gap-2 md:gap-3">
            {/* Search */}
            {showSearch && (
              <div className="hidden lg:block">
                <Search />
              </div>
            )}

            {/* Notifications */}
            <Button
              variant="ghost"
              size="icon"
              className="relative h-8 w-8 shrink-0 rounded-full border bg-white transition-colors hover:bg-gray-50 sm:h-9 sm:w-9 md:h-10 md:w-10"
              aria-label="Notifications"
              onClick={() =>
                setIsNotificationSheetOpen?.(!isNotificationSheetOpen)
              }
            >
              <Bell className="h-5 w-5 sm:h-4 sm:w-4 md:h-[18px] md:w-[18px]" />

              {unreadCount > 0 && (
                <span className="absolute -right-0.5 -top-0.5 flex h-3.5 w-3.5 items-center justify-center rounded-full bg-red-500 text-[9px] font-medium text-white sm:-right-1 sm:-top-1 sm:h-4 sm:w-4 sm:text-[10px]">
                  {unreadCount > 99 ? "99+" : unreadCount}
                </span>
              )}
            </Button>
          </div>
        )}
      </div>
    </header>
  );
}
