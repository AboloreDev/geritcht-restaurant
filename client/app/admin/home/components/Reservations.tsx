"use client";

import Link from "next/link";
import { CheckCircle } from "@mynaui/icons-react";
import {
  useCheckInReservationMutation,
  useGetAllRservationsQuery,
} from "@/app/state/api/reservationsApi";
import { formatTimeSlot } from "@/app/utils/timeSlots";
import { statusStyle } from "@/app/utils/reservationStatusHelper";
import { DashboardListCard } from "./DashboardListCard";
import { toast } from "sonner";
import { getApiError } from "@/app/utils/apiError";
import { todayDateString } from "@/app/utils/dateRanges";
import { ReservationResponse } from "@/app/state/types/reservationTypes";

export function TodaysReservationsList() {
  const today = todayDateString();
  const { data, isLoading } = useGetAllRservationsQuery({
    date: today,
    page: 1,
    page_size: 5,
  });

  const [checkIn, { isLoading: isCheckingIn }] =
    useCheckInReservationMutation();

  const reservation = data?.data ?? [];
  const reservations = reservation.reservations ?? [];

  async function handleCheckIn(e: React.MouseEvent, id: number) {
    e.preventDefault();
    e.stopPropagation();
    try {
      const response = await checkIn({
        id,
      }).unwrap();
      toast.success(response.message || "Checked in successfully");
    } catch (err) {
      console.error(err);
      toast.error(getApiError(err));
    }
  }

  return (
    <DashboardListCard
      title="Today's Reservations"
      viewAllHref="/admin/reservations"
      isLoading={isLoading}
      isEmpty={reservations.length === 0}
      emptyMessage="No reservations today."
    >
      {reservations.map((r: ReservationResponse) => (
        <Link
          key={r.id}
          href={`/admin/reservations/${r.id}`}
          className="flex items-center justify-between gap-2 rounded-lg p-3 text-sm hover:bg-muted"
        >
          <div className="min-w-0">
            <p className="font-medium">
              {formatTimeSlot(r.time_slot)} · Table {r.table_id}
            </p>
            <p className="text-xs text-muted-foreground">
              {r.party_size} guest{r.party_size !== 1 ? "s" : ""}
            </p>
          </div>

          <div className="flex shrink-0 items-center gap-2">
            {!r.checked_in_at &&
              r.status !== "cancelled" &&
              r.status !== "no_show" && (
                <button
                  onClick={(e) => handleCheckIn(e, r.id)}
                  disabled={isCheckingIn}
                  aria-label="Check in"
                  className="flex h-7 w-7 items-center justify-center rounded-full border text-emerald-600 hover:bg-emerald-50"
                >
                  <CheckCircle className="h-4 w-4" />
                </button>
              )}
            <span
              className={`rounded-full px-2.5 py-1 text-[10px] font-medium capitalize ${statusStyle(r.status)}`}
            >
              {r.status}
            </span>
          </div>
        </Link>
      ))}
    </DashboardListCard>
  );
}
