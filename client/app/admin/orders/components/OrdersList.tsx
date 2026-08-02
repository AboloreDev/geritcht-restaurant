"use client";

import React from "react";
import Link from "next/link";
import { AnimatePresence, motion } from "framer-motion";

import { useAppDispatch } from "@/app/state/redux";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { formatNaira } from "@/app/utils/formatNaira";
import { Clipboard } from "@mynaui/icons-react";
import { statusStyle } from "@/app/utils/orderStatusHelpers";
import { Order } from "@/app/state/types/orderTypes";
import { setPage } from "@/app/state/slices/orderSlice";
import { useAdminCancelOrderMutation } from "@/app/state/api/orderApi";
import { toast } from "sonner";
import { getApiError } from "@/app/utils/apiError";

interface OrdersListProps {
  orders: Order[];
  isLoading: boolean;
  isFetching: boolean;
  hasMore: boolean;
  page: number;
}

const OrdersList = ({
  orders,
  isLoading,
  isFetching,
  hasMore,
  page,
}: OrdersListProps) => {
  const dispatch = useAppDispatch();
  const [cancellingId, setCancellingId] = React.useState<number | null>(null);
  const [cancelOrder, { isLoading: isCancelling }] =
    useAdminCancelOrderMutation();
  const isRefetching = isFetching && !isLoading && page === 1;

  async function handleCancel(e: React.MouseEvent, id: number) {
    e.preventDefault();
    e.stopPropagation();
    setCancellingId(id);
    try {
      const response = await cancelOrder({ id }).unwrap();
      toast.success(response.message);
    } catch (err) {
      toast.error(getApiError(err));
    } finally {
      setCancellingId(null);
    }
  }

  /* Loading state */
  if (isLoading) {
    return (
      <div className="mt-6 space-y-3">
        {Array.from({ length: 8 }).map((_, i) => (
          <div key={i} className="h-20 animate-pulse rounded-xl bg-[#d4a373]" />
        ))}
      </div>
    );
  }

  /* Empty state */
  if (!isFetching && orders.length === 0) {
    return (
      <div className="mt-16 flex flex-col items-center gap-3 text-center">
        <Clipboard className="h-10 w-10 text-muted-foreground" />
        <p className="text-sm text-muted-foreground">
          No orders match this filter.
        </p>
      </div>
    );
  }

  /*  Content  */
  return (
    <>
      <div
        className={cn(
          "mt-6 space-y-3 transition-opacity",
          isRefetching && "opacity-50",
        )}
      >
        <AnimatePresence initial={false}>
          {orders.map((order, index) => (
            <motion.div
              key={order.id}
              initial={{ opacity: 0, y: 8 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{
                duration: 0.25,
                delay: index * 0.02,
              }}
            >
              <Link
                href={`/admin/orders/${order.id}`}
                className="flex items-center justify-between gap-4 rounded-2xl bg-[#faedcd] p-4 transition-colors hover:bg-[#fefae0]"
              >
                {/* Left */}
                <div className="min-w-0 flex-1">
                  <div className="flex items-center gap-2">
                    <p className="font-medium">Order #{order.id}</p>

                    <span className="text-xs capitalize text-muted-foreground">
                      • {order.type.replace("_", " ")}
                    </span>
                  </div>

                  <p className="mt-1 text-sm text-muted-foreground">
                    {order.user?.first_name ?? "Guest"} •{" "}
                    {order.order_items.length} item
                    {order.order_items.length !== 1 && "s"} •{" "}
                    {formatNaira(order.total_amount)}
                  </p>

                  <p className="mt-0.5 text-xs text-muted-foreground">
                    {new Date(order.created_at).toLocaleString("en-NG", {
                      dateStyle: "medium",
                      timeStyle: "short",
                    })}
                  </p>
                </div>

                {/* Right */}
                <div className="flex shrink-0 items-center gap-3">
                  <span className="text-xs capitalize text-muted-foreground">
                    {order.payment_status}
                  </span>

                  <span
                    className={cn(
                      "rounded-full px-2.5 py-1 text-xs font-medium capitalize",
                      statusStyle(order.status),
                    )}
                  >
                    {order.status}
                  </span>

                  {order.status !== "cancelled" &&
                    order.status !== "completed" && (
                      <Button
                        size="xs"
                        variant="outline"
                        disabled={isCancelling && cancellingId === order.id}
                        onClick={(e) => handleCancel(e, order.id)}
                        className="bg-red-500 text-white hover:bg-destructive/10 border-none"
                      >
                        {isCancelling && cancellingId === order.id
                          ? "..."
                          : "Cancel"}
                      </Button>
                    )}
                </div>
              </Link>
            </motion.div>
          ))}
        </AnimatePresence>
      </div>

      {hasMore && (
        <div className="mt-6 flex justify-center">
          <Button
            variant="outline"
            disabled={isFetching}
            onClick={() => dispatch(setPage(page + 1))}
            className="border-none"
          >
            {isFetching && page > 1 ? "Loading..." : "Load more"}
          </Button>
        </div>
      )}
    </>
  );
};

export default OrdersList;
