"use client";

import { useGetAllOrdersQuery } from "../state/api/orderApi";
import { useGetAllRservationsQuery } from "../state/api/reservationsApi";
import { useGetAllUsersQuery } from "../state/api/userApi";
import { Order } from "../state/types/orderTypes";
import { getWeekRange, todayDateString } from "../utils/dateRanges";

export function useDashboardStats() {
  const today = todayDateString();

  // pending orders — only need the count, so page_size: 1
  const { data: pendingOrdersData, isLoading: isLoadingPending } =
    useGetAllOrdersQuery({ page: 1, page_size: 1 });

  // today's reservations — same count-only trick
  const { data: reservationsTodayData, isLoading: isLoadingReservations } =
    useGetAllRservationsQuery({ date: today, page: 1, page_size: 1 });

  // today's revenue — need actual order amounts, not just count,
  // so fetch a real page of today's orders and sum client-side
  const { data: todaysOrdersData, isLoading: isLoadingRevenue } =
    useGetAllOrdersQuery({ date: today, page: 1, page_size: 10000 });

  // all users — no date filtering available server-side, so we
  // fetch everything and bucket by created_at locally
  const { data: usersData, isLoading: isLoadingUsers } = useGetAllUsersQuery();

  // @ts-expect-error "type inference"
  const pendingOrdersCount = pendingOrdersData?.data.total ?? 0;
  // @ts-expect-error "type inference"
  const reservationsTodayCount = reservationsTodayData?.data.total ?? 0;

  // @ts-expect-error "type inference"
  const revenueToday = (todaysOrdersData?.data.orders ?? []).reduce(
    (sum: number, order: Order) => sum + order.total_amount,
    0,
  );

  const users = usersData?.data ?? [];
  const thisWeek = getWeekRange(0);
  const lastWeek = getWeekRange(1);

  const thisWeekCount = users.filter((u) => {
    const created = new Date(u.created_at);
    return created >= thisWeek.start && created < thisWeek.end;
  }).length;

  const lastWeekCount = users.filter((u) => {
    const created = new Date(u.created_at);
    return created >= lastWeek.start && created < lastWeek.end;
  }).length;

  const userGrowthPercent =
    lastWeekCount === 0
      ? thisWeekCount > 0
        ? 100
        : 0
      : Math.round(((thisWeekCount - lastWeekCount) / lastWeekCount) * 100);

  return {
    pendingOrdersCount,
    reservationsTodayCount,
    revenueToday,
    userGrowthPercent,
    thisWeekCount,
    isLoading:
      isLoadingPending ||
      isLoadingReservations ||
      isLoadingRevenue ||
      isLoadingUsers,
  };
}
