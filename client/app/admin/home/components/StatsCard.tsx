"use client";

import { useDashboardStats } from "@/app/hooks/useDashboardStats";
import { formatNaira } from "@/app/utils/formatNaira";
import {
  TrendingDown,
  TrendingUp,
  Clipboard,
  CalendarCheck,
  Pocket,
} from "@mynaui/icons-react";

function StatCard({
  label,
  value,
  icon: Icon,
  isLoading,
  trend,
}: {
  label: string;
  value: string;
  icon: React.ComponentType<{ className?: string }>;
  isLoading: boolean;
  trend?: { value: number; positive: boolean };
}) {
  return (
    <div className="rounded-2xl bg-[#faedcd] p-5">
      <div className="flex items-center justify-between">
        <span className="text-sm text-muted-foreground">{label}</span>
        <Icon className="h-4 w-4 text-muted-foreground" />
      </div>

      {isLoading ? (
        <div className="mt-3 h-7 w-20 animate-pulse bg-[#fefae0]" />
      ) : (
        <p className="mt-2 text-2xl font-semibold">{value}</p>
      )}

      {trend && !isLoading && (
        <div
          className={`mt-1 flex items-center gap-1 text-xs ${trend.positive ? "text-emerald-600" : "text-red-500"}`}
        >
          {trend.positive ? (
            <TrendingUp className="h-3 w-3" />
          ) : (
            <TrendingDown className="h-3 w-3" />
          )}
          {Math.abs(trend.value)}% vs last week
        </div>
      )}
    </div>
  );
}

export function StatsCards() {
  const stats = useDashboardStats();

  return (
    <div className="grid grid-cols-2 gap-4 lg:grid-cols-4">
      <StatCard
        label="Pending Orders"
        value={String(stats.pendingOrdersCount)}
        icon={Clipboard}
        isLoading={stats.isLoading}
      />
      <StatCard
        label="Reservations Today"
        value={String(stats.reservationsTodayCount)}
        icon={CalendarCheck}
        isLoading={stats.isLoading}
      />
      <StatCard
        label="Revenue Today"
        value={formatNaira(stats.revenueToday)}
        icon={Pocket}
        isLoading={stats.isLoading}
      />
      <StatCard
        label="New Users"
        value={String(stats.thisWeekCount)}
        icon={TrendingUp}
        isLoading={stats.isLoading}
        trend={{
          value: stats.userGrowthPercent,
          positive: stats.userGrowthPercent >= 0,
        }}
      />
    </div>
  );
}
