"use client";

import { motion, AnimatePresence } from "framer-motion";
import { ArrowLeft, CalendarCheck, Spinner } from "@mynaui/icons-react";
import { useAppDispatch, useAppSelector, RootState } from "@/app/state/redux";
import { cn } from "@/lib/utils";

import {
  resetReservationFilters,
  setFilterDate,
  setFilterStatus,
  setPage,
} from "@/app/state/slices/reservationSlice";
import { formatTimeSlot } from "@/app/utils/timeSlots";
import { Button } from "@/components/ui/button";
import Link from "next/link";
import { useGetAllUserRservationsQuery } from "@/app/state/api/reservationsApi";
import { ReservationResponse } from "@/app/state/types/reservationTypes";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useEffect, useState } from "react";
import { useDebounce } from "@/app/hooks/useDebounce";
import {
  STATUS_OPTIONS,
  statusStyle,
} from "@/app/utils/reservationStatusHelper";

export function MyReservationsContent() {
  const dispatch = useAppDispatch();
  const { page, pageSize, filterDate, filterStatus } = useAppSelector(
    (state: RootState) => state.reservation,
  );

  const { data, isLoading, isFetching } = useGetAllUserRservationsQuery({
    page,
    page_size: pageSize,
    date: filterDate || undefined,
    status: filterStatus || undefined,
  });

  // @ts-expect-error "<>"
  const reservations = data?.data.reservations ?? [];
  const hasMore = data ? (page ?? 1) < data.total_pages : false;
  const hasActiveFilters = Boolean(filterDate || filterStatus);

  // refetching due to a filter change, not the very first load —
  // filter changes always reset page to 1, so this won't fire for load-more
  const isRefetchingFilters = isFetching && !isLoading && (page ?? 1) === 1;

  const [dateInput, setDateInput] = useState(filterDate ?? "");
  const debouncedDate = useDebounce(dateInput, 400);

  useEffect(() => {
    dispatch(setFilterDate(debouncedDate || undefined));
  }, [debouncedDate, dispatch]);

  return (
    <div className="min-h-screen bg-[url('/assets/bg.png')] bg-cover bg-center bg-fixed">
      <div className="mx-auto max-w-3xl px-6 py-10">
        <h1 className="font-serif text-3xl text-primary font-semibold">
          My Reservations
        </h1>
        <p className="mt-2 text-sm text-primary-deep">
          Your upcoming and past table bookings.
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
              // @ts-expect-error "<>"
              dispatch(setFilterStatus(v === "all" ? undefined : v))
            }
          >
            <SelectTrigger className="w-40 cursor-pointer">
              <SelectValue placeholder="All" className="" />
            </SelectTrigger>
            <SelectContent className="bg-[#fefae0] cursor-pointer">
              <SelectItem value="all" className="">
                All
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

          {hasActiveFilters && (
            <button
              onClick={() => dispatch(resetReservationFilters())}
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
        ) : reservations.length === 0 && !isFetching ? (
          <div className="mt-16 flex flex-col items-center gap-3 text-center">
            <CalendarCheck className="h-10 w-10 text-primary" />
            <p className="text-sm text-primary">
              You haven&apos;t made any reservations yet.
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
                {reservations.map((r: ReservationResponse, i: number) => (
                  <motion.div
                    key={r.id}
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
                      href={`/reservation/${r.id}`}
                      className="flex items-start justify-between gap-3"
                    >
                      <div>
                        <p className="font-medium">
                          {r.date} · {formatTimeSlot(r.time_slot)}
                        </p>
                        <p className="mt-1 text-sm text-muted-foreground">
                          {r.party_size}{" "}
                          {r.party_size === 1 ? "guest" : "guests"}
                        </p>
                        <p className="mt-1 text-sm text-muted-foreground">
                          Table Name: {r.table.name} · Location:{" "}
                          {r.table.location}
                        </p>
                        {r.special_requests && (
                          <p className="mt-2 text-sm italic text-muted-foreground">
                            &quot;{r.special_requests}&quot;
                          </p>
                        )}
                      </div>

                      <span
                        className={`shrink-0 rounded-full px-3 py-1 text-xs font-medium capitalize ${statusStyle(r.status)}`}
                      >
                        {r.status}
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
