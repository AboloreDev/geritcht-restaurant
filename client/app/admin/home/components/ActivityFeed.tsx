"use client";

import { useGetAllOrdersQuery } from "@/app/state/api/orderApi";
import { useGetAllRservationsQuery } from "@/app/state/api/reservationsApi";
import { DashboardListCard } from "./DashboardListCard";
import { ReservationResponse } from "@/app/state/types/reservationTypes";
import { Order } from "@/app/state/types/orderTypes";
import { CalendarCheck, Clipboard } from "@mynaui/icons-react";

type ActivityItem = {
  id: string;
  type: "order" | "reservation";
  label: string;
  timestamp: string;
};

export function ActivityFeed() {
  const {
    data: ordersData,
    isLoading: isLoadingOrders,
    isFetching: isFetchingOrders,
  } = useGetAllOrdersQuery({
    page: 1,
    page_size: 5,
  });
  const {
    data: reservationsData,
    isLoading: isLoadingReservations,
    isFetching: isFetchingReservations,
  } = useGetAllRservationsQuery({
    status: "confirmed",
    page: 1,
    page_size: 5,
  });

  console.log("ordersData", ordersData);
  console.log("reservationsData", reservationsData);

  const isLoading =
    isLoadingOrders ||
    isLoadingReservations ||
    isFetchingOrders ||
    isFetchingReservations;

  // @ts-expect-error "type inference"
  const orderItems: ActivityItem[] = (ordersData?.data.orders ?? []).map(
    (o: Order) => ({
      id: `order-${o.id}`,
      type: "order",
      label: `New order #${o.id} placed`,
      timestamp: o.created_at,
    }),
  );

  const reservationItems: ActivityItem[] =
    // @ts-expect-error "type inference"
    (reservationsData?.data.reservations ?? []).map(
      (r: ReservationResponse) => ({
        id: `reservation-${r.id}`,
        type: "reservation",
        label: `New reservation for ${r.party_size} guests`,
        timestamp: r.created_at,
      }),
    );

  const activity = [...orderItems, ...reservationItems]
    .sort(
      (a, b) =>
        new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime(),
    )
    .slice(0, 8);

  return (
    <DashboardListCard
      title="Recent Activity"
      viewAllHref="/admin/activity"
      isLoading={isLoading}
      isEmpty={activity.length === 0}
      emptyMessage="Nothing's happened yet."
      isFetching={isLoading}
    >
      {activity.map((item) => (
        <div
          key={item.id}
          className="flex items-center gap-3 rounded-lg p-3 text-sm"
        >
          {item.type === "order" ? (
            <Clipboard className="h-4 w-4 shrink-0 text-muted-foreground" />
          ) : (
            <CalendarCheck className="h-4 w-4 shrink-0 text-muted-foreground" />
          )}
          <div className="min-w-0 flex-1">
            <p className="truncate">{item.label}</p>
            <p className="text-xs text-muted-foreground">
              {new Date(item.timestamp).toLocaleString("en-NG", {
                dateStyle: "medium",
                timeStyle: "short",
              })}
            </p>
          </div>
        </div>
      ))}
    </DashboardListCard>
  );
}
