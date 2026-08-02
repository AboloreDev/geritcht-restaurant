"use client";

import { motion, AnimatePresence } from "framer-motion";
import { ArrowLeft, CalendarCheck, Spinner } from "@mynaui/icons-react";
import { useAppDispatch, useAppSelector, RootState } from "@/app/state/redux";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import Link from "next/link";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useEffect, useState } from "react";
import { useDebounce } from "@/app/hooks/useDebounce";
import { useGetAllUserTakeoutOrdersQuery } from "@/app/state/api/orderApi";
import { Order } from "@/app/state/types/orderTypes";
import {
  STATUS_OPTIONS,
  STATUS_TYPES,
  statusStyle,
} from "@/app/utils/orderStatusHelpers";
import {
  resetOrdersFilters,
  setFilterDate,
  setFilterStatus,
  setFilterType,
  setPage,
} from "@/app/state/slices/orderSlice";

export function OrderContents() {
  const dispatch = useAppDispatch();
  const { page, pageSize, filterDate, filterStatus, filterType } =
    useAppSelector((state: RootState) => state.order);

  const { data, isLoading, isFetching } = useGetAllUserTakeoutOrdersQuery({
    page,
    page_size: pageSize,
    date: filterDate || undefined,
    status: filterStatus || undefined,
    type: filterType || undefined,
  });

  const orders = data?.data.orders ?? [];
  const hasMore = data ? data.data.page < data.data.total_pages : false;
  const hasActiveFilters = Boolean(filterDate || filterStatus || filterType);

  const isRefetchingFilters = isFetching && !isLoading && (page ?? 1) === 1;

  const [dateInput, setDateInput] = useState(filterDate ?? "");
  const debouncedDate = useDebounce(dateInput, 1000);

  useEffect(() => {
    dispatch(setFilterDate(debouncedDate || undefined));
  }, [debouncedDate, dispatch]);

  return (
    <div className="min-h-screen bg-[url('/assets/bg.png')] bg-cover bg-center bg-fixed">
      <div className="mx-auto max-w-3xl px-6 py-10">
        <h1 className="text-3xl text-primary font-semibold">My Orders</h1>

        <p className="mt-2 text-sm text-primary-deep">
          Your current and previous orders.
        </p>

        <Link
          href="/"
          className="mt-4 inline-flex items-center gap-1.5 text-sm text-primary transition-colors hover:text-primary-deep"
        >
          <ArrowLeft className="h-4 w-4" />
          Back
        </Link>

        {/* filter bar */}
        <div className="mt-6 cursor-pointer flex z-40 flex-wrap bg-[#fefae0] rounded-2xl px-6 py-3 items-center gap-3">
          <input
            type="date"
            value={dateInput}
            onChange={(e) => setDateInput(e.target.value)}
            className="rounded-lg cursor-pointer border px-3 py-1.5 text-sm"
          />

          <Select
            value={filterStatus ?? "all"}
            onValueChange={(v) =>
              //   @ts-expect-error "<>"
              dispatch(setFilterStatus(v === "all" ? undefined : v))
            }
          >
            <SelectTrigger className="w-40 cursor-pointer">
              <SelectValue placeholder="All statuses" className="" />
            </SelectTrigger>
            <SelectContent className="bg-[#fefae0] cursor-pointer">
              <SelectItem value="all" className="">
                All statuses
              </SelectItem>
              {STATUS_OPTIONS.map((opt) => (
                <SelectItem
                  className="cursor-pointer"
                  key={opt.value}
                  value={opt.value}
                >
                  {opt.label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>

          <Select
            value={filterType ?? "all"}
            onValueChange={(v) =>
              // @ts-expect-error "<>"
              dispatch(setFilterType(v === "all" ? undefined : v))
            }
          >
            <SelectTrigger className="w-40 cursor-pointer">
              <SelectValue placeholder="All Types" className="" />
            </SelectTrigger>
            <SelectContent className="bg-[#fefae0] cursor-pointer">
              <SelectItem value="all" className="">
                All Types
              </SelectItem>
              {STATUS_TYPES.map((opt) => (
                <SelectItem
                  className="cursor-pointer"
                  key={opt.value}
                  value={opt.value}
                >
                  {opt.label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>

          {hasActiveFilters && (
            <button
              onClick={() => dispatch(resetOrdersFilters())}
              className="text-xs text-muted-foreground cursor-pointer hover:text-foreground"
            >
              Clear filters
            </button>
          )}

          {isRefetchingFilters && (
            <span className="flex items-center gap-1.5 text-xs text-muted-foreground">
              <Spinner className="h-3.5 w-3.5 animate-spin" />
              Updating…
            </span>
          )}
        </div>

        {isLoading ? (
          <div className="mt-8 space-y-4">
            {Array.from({ length: 5 }).map((_, i) => (
              <div
                key={i}
                className="h-28 animate-pulse rounded-xl bg-primary-deep"
              />
            ))}
          </div>
        ) : orders.length === 0 && !isFetching ? (
          <div className="mt-16 flex flex-col items-center gap-3 text-center">
            <CalendarCheck className="h-10 w-10 text-primary" />
            <p className="text-sm text-primary">
              You haven&apos;t made any orders yet.
            </p>
            <Link
              href="/menu"
              className="text-sm font-medium text-primary-deep hover:underline"
            >
              Browse the menu
            </Link>
          </div>
        ) : (
          <>
            <div
              className={cn(
                "mt-8 space-y-4 transition-opacity",
                isRefetchingFilters && "opacity-50",
              )}
            >
              <AnimatePresence initial={false}>
                {orders.map((order: Order, i: number) => (
                  <motion.div
                    key={order.id}
                    initial={{ opacity: 0, y: 12 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{
                      duration: 0.3,
                      delay: i * 0.03,
                      ease: [0.16, 1, 0.3, 1],
                    }}
                    className="bg-[#fefae0] rounded-2xl p-5"
                  >
                    <Link
                      href={`/orders/${order.id}`}
                      className="flex items-start justify-between gap-3"
                    >
                      <div>
                        <p className="font-medium">Order No: #{order.id}</p>

                        <p className="mt-1 text-sm text-muted-foreground">
                          Date:{" "}
                          {new Date(order.created_at).toLocaleDateString()}
                        </p>

                        <p className="mt-1 text-sm text-muted-foreground">
                          {order.order_items.length}{" "}
                          {order.order_items.length === 1 ? "item" : "items"}
                        </p>

                        <p className="mt-1 text-sm text-muted-foreground">
                          Total: ₦{order.total_amount.toLocaleString()}
                        </p>

                        <p className="mt-1 text-md text-muted-foreground capitalize">
                          Payment Status: {order.payment_status}
                        </p>

                        {order.notes && (
                          <p className="mt-2 text-sm italic text-muted-foreground">
                            &quot;{order.notes}&quot;
                          </p>
                        )}
                      </div>

                      <span
                        className={`shrink-0 rounded-full px-3 py-1 text-xs font-medium capitalize ${statusStyle(
                          order.status,
                        )}`}
                      >
                        {order.status}
                      </span>
                    </Link>
                  </motion.div>
                ))}
              </AnimatePresence>
            </div>

            {hasMore && (
              <div className="mt-8 flex justify-center">
                <Button
                  variant="outline"
                  disabled={isFetching}
                  onClick={() => dispatch(setPage((page ?? 1) + 1))}
                >
                  {isFetching && (page ?? 1) > 1 ? "Loading…" : "Load more"}
                </Button>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
}
