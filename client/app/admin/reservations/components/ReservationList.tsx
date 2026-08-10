"use client";

import React from "react";
import Link from "next/link";
import { AnimatePresence, motion } from "framer-motion";

import { useAppDispatch } from "@/app/state/redux";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { CalendarCheck } from "@mynaui/icons-react";
import { statusStyle } from "@/app/utils/reservationStatusHelper";
import { formatTimeSlot } from "@/app/utils/timeSlots";
import {
  Reservation,
  ReservationListResponse,
  ReservationResponse,
} from "@/app/state/types/reservationTypes";
import { setPage } from "@/app/state/slices/reservationSlice";
import {
  useCheckInReservationMutation,
  useAdminCancelReservationMutation,
} from "@/app/state/api/reservationsApi";
import { toast } from "sonner";
import { getApiError } from "@/app/utils/apiError";

interface ReservationsListProps {
  reservations: Reservation[];
  isLoading: boolean;
  isFetching: boolean;
  hasMore: boolean;
  page: number;
}

const ReservationsList = ({
  reservations,
  isLoading,
  isFetching,
  hasMore,
  page,
}: ReservationsListProps) => {
  const dispatch = useAppDispatch();
  const [actioningId, setActioningId] = React.useState<number | null>(null);

  const [checkIn, { isLoading: isCheckingIn }] =
    useCheckInReservationMutation();
  const [cancelReservation, { isLoading: isCancelling }] =
    useAdminCancelReservationMutation();

  const isRefetching = isFetching && !isLoading && page === 1;

  async function handleAction(
    e: React.MouseEvent,
    id: number,
    action: "checkin" | "cancel",
  ) {
    e.preventDefault();
    e.stopPropagation();
    setActioningId(id);
    try {
      const mutation = action === "checkin" ? checkIn : cancelReservation;
      const response = await mutation({ id }).unwrap();
      toast.success(response.message);
    } catch (err) {
      toast.error(getApiError(err));
    } finally {
      setActioningId(null);
    }
  }

  /* Loading state */
  if (isLoading) {
    return (
      <div className="mt-3 space-y-3   p-4">
        {Array.from({ length: 8 }).map((_, i) => (
          <div
            key={i}
            className="h-20 animate-pulse rounded-xl bg-[#d4a373] px-6"
          />
        ))}
      </div>
    );
  }

  /* Empty state */
  if (!isFetching && reservations.length === 0) {
    return (
      <div className="mt-16 flex flex-col items-center gap-3  p-4 text-center">
        <CalendarCheck className="h-10 w-10 text-muted-foreground" />
        <p className="text-sm text-muted-foreground">
          No reservations match this filter.
        </p>
      </div>
    );
  }

  /* Content */
  return (
    <>
      <div
        className={cn(
          "mt-3 space-y-3 p-4 transition-opacity",
          isRefetching && "opacity-50",
        )}
      >
        <AnimatePresence initial={false}>
          {reservations.map((reservation, index) => {
            const isDone =
              reservation.status === "cancelled" ||
              reservation.status === "no_show" ||
              reservation.status === "completed";
            const isActioning = actioningId === reservation.id;

            return (
              <motion.div
                key={reservation.id}
                initial={{ opacity: 0, y: 8 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.25, delay: index * 0.02 }}
              >
                <Link
                  href={`/admin/reservations/${reservation.id}`}
                  className="flex items-center justify-between gap-4 rounded-2xl bg-[#faedcd] p-4 transition-colors hover:bg-[#fefae0]"
                >
                  {/* Left */}
                  <div className="min-w-0 flex-1">
                    <div className="flex items-center gap-2">
                      <p className="font-medium">
                        {reservation.date} ·{" "}
                        {formatTimeSlot(reservation.time_slot)}
                      </p>

                      <span className="text-xs capitalize text-muted-foreground">
                        •{" "}
                        {reservation.table?.name ??
                          `Table ${reservation.table_id}`}
                      </span>
                    </div>

                    <p className="mt-1 text-sm text-muted-foreground">
                      {reservation.user?.first_name ?? "Guest"} •{" "}
                      {reservation.party_size} guest
                      {reservation.party_size !== 1 && "s"}
                    </p>

                    {reservation.special_requests && (
                      <p className="mt-0.5 text-xs italic text-muted-foreground">
                        &quot;{reservation.special_requests}&quot;
                      </p>
                    )}
                  </div>

                  {/* Right */}
                  <div className="flex shrink-0 items-center gap-2">
                    <span
                      className={cn(
                        "rounded-full px-2.5 py-1 text-xs font-medium capitalize",
                        statusStyle(reservation.status),
                      )}
                    >
                      {reservation.status.replace("_", " ")}
                    </span>

                    {/* no-show isn't a manual admin action — it's set
                        automatically by a background worker once a
                        reservation's time passes without check-in, so
                        only check-in and cancel are offered here */}
                    {!isDone && (
                      <>
                        {!reservation.checked_in_at && (
                          <Button
                            size="xs"
                            variant="outline"
                            disabled={isCheckingIn && isActioning}
                            onClick={(e) =>
                              handleAction(e, reservation.id, "checkin")
                            }
                            className="border-none bg-emerald-600 text-white hover:bg-emerald-700"
                          >
                            {isCheckingIn && isActioning ? "..." : "Check in"}
                          </Button>
                        )}

                        <Button
                          size="xs"
                          variant="outline"
                          disabled={isCancelling && isActioning}
                          onClick={(e) =>
                            handleAction(e, reservation.id, "cancel")
                          }
                          className="border-none bg-red-500 text-white hover:bg-red-600"
                        >
                          {isCancelling && isActioning ? "..." : "Cancel"}
                        </Button>
                      </>
                    )}
                  </div>
                </Link>
              </motion.div>
            );
          })}
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

export default ReservationsList;
