"use client";

import Link from "next/link";
import { useGetAllOrdersQuery } from "@/app/state/api/orderApi";
import { formatNaira } from "@/app/utils/formatNaira";
import { statusStyle } from "@/app/utils/orderStatusHelpers";
import { DashboardListCard } from "./DashboardListCard";
import { Order } from "@/app/state/types/orderTypes";

export function PendingOrdersList() {
  const { data, isLoading, isFetching } = useGetAllOrdersQuery({
    status: "pending",
    page: 1,
    page_size: 5,
  });

  const orders = data?.data.orders ?? [];

  return (
    <DashboardListCard
      title="Pending Orders"
      viewAllHref="/admin/orders?status=pending"
      isLoading={isLoading}
      isEmpty={orders.length === 0}
      emptyMessage="No pending orders."
      isFetching={isFetching}
    >
      {orders.map((order: Order) => (
        <Link
          key={order.id}
          href={`/admin/orders/${order.id}`}
          className="flex items-center justify-between rounded-lg p-3 text-sm hover:bg-muted"
        >
          <div>
            <p className="font-medium">Order #{order.id}</p>
            <p className="text-xs text-muted-foreground">
              {order.order_items.length} item
              {order.order_items.length !== 1 ? "s" : ""} ·{" "}
              {formatNaira(order.total_amount)}
            </p>
          </div>
          <span
            className={`shrink-0 rounded-full px-2.5 py-1 text-[10px] font-medium capitalize ${statusStyle(order.status)}`}
          >
            {order.status}
          </span>
        </Link>
      ))}
    </DashboardListCard>
  );
}
