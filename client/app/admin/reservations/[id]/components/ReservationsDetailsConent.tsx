"use client";

import { useState } from "react";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import { HugeiconsIcon } from "@hugeicons/react";
import {
  ArrowLeft01Icon,
  UserGroupIcon,
  Calendar01Icon,
  Clock01Icon,
  CheckmarkCircle01Icon,
  Cancel01Icon,
} from "@hugeicons/core-free-icons";

import { useGetReservationByIdQuery } from "@/app/state/api/reservationsApi";
import {
  useCheckInReservationMutation,
  useAdminCancelReservationMutation,
} from "@/app/state/api/reservationsApi";
import { statusStyle } from "@/app/utils/reservationStatusHelper";
import { formatTimeSlot } from "@/app/utils/timeSlots";
import { tableStatusStyle } from "@/app/utils/tableStatusHelpers";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";
import { getApiError } from "@/app/utils/apiError";
import { ReservationDetailSkeleton } from "./ReservationDetailsSkeleton";

export default function ReservationDetailContent() {
  const { id } = useParams();
  const router = useRouter();
  const reservationId = Number(id);

  const { data, isLoading, isError } = useGetReservationByIdQuery({
    id: reservationId,
  });

  const [checkIn, { isLoading: isCheckingIn }] =
    useCheckInReservationMutation();
  const [cancelReservation, { isLoading: isCancelling }] =
    useAdminCancelReservationMutation();
  const [actioning, setActioning] = useState<"checkin" | "cancel" | null>(null);

  if (isLoading) return <ReservationDetailSkeleton />;

  if (isError || !data?.data) {
    return (
      <div className="flex min-h-[60vh] flex-col items-center justify-center gap-3 text-center">
        <p className="text-lg font-medium">
          This reservation couldn&apos;t be found.
        </p>
        <Link
          href="/admin/reservations"
          className="text-sm text-primary hover:underline"
        >
          Back to reservations
        </Link>
      </div>
    );
  }

  const reservation = data.data;
  const isDone =
    reservation.status === "cancelled" ||
    reservation.status === "no_show" ||
    reservation.status === "completed";

  async function handleCheckIn() {
    setActioning("checkin");
    try {
      const response = await checkIn({ id: reservationId }).unwrap();
      toast.success(response.message);
    } catch (err) {
      toast.error(getApiError(err));
    } finally {
      setActioning(null);
    }
  }

  async function handleCancel() {
    setActioning("cancel");
    try {
      const response = await cancelReservation({ id: reservationId }).unwrap();
      toast.success(response.message);
    } catch (err) {
      toast.error(getApiError(err));
    } finally {
      setActioning(null);
    }
  }

  return (
    <div className="p-6">
      <button
        onClick={() => router.push("/admin/reservations")}
        className="flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground"
      >
        <HugeiconsIcon icon={ArrowLeft01Icon} size={16} />
        Back to reservations
      </button>

      {/* header */}
      <div className="mt-4 flex flex-wrap items-start justify-between gap-4">
        <div>
          <div className="flex items-center gap-3">
            <h1 className="font-serif text-2xl font-semibold">
              {reservation.date} · {formatTimeSlot(reservation.time_slot)}
            </h1>
            <span
              className={`rounded-full px-2.5 py-1 text-xs font-medium capitalize ${statusStyle(reservation.status)}`}
            >
              {reservation.status.replace("_", " ")}
            </span>
          </div>
          <p className="mt-1 text-sm text-muted-foreground">
            #{reservation.id} · Booked{" "}
            {new Date(reservation.created_at).toLocaleString("en-NG", {
              dateStyle: "medium",
              timeStyle: "short",
            })}
          </p>
        </div>

        {!isDone && (
          <div className="flex items-center gap-2">
            {!reservation.checked_in_at && (
              <Button
                variant="outline"
                disabled={isCheckingIn && actioning === "checkin"}
                onClick={handleCheckIn}
                className="text-emerald-600 hover:bg-emerald-50"
              >
                <HugeiconsIcon
                  icon={CheckmarkCircle01Icon}
                  strokeWidth={2}
                  size={16}
                />
                {isCheckingIn && actioning === "checkin"
                  ? "Checking in…"
                  : "Check in"}
              </Button>
            )}
            <Button
              variant="outline"
              disabled={isCancelling && actioning === "cancel"}
              onClick={handleCancel}
              className="bg-red-500 hover:bg-red-600 text-white border-none"
            >
              <HugeiconsIcon icon={Cancel01Icon} strokeWidth={2} size={16} />
              {isCancelling && actioning === "cancel"
                ? "Cancelling…"
                : "Cancel"}
            </Button>
          </div>
        )}
      </div>

      <div className="mt-6 grid gap-6 lg:grid-cols-3">
        {/* main column */}
        <div className="space-y-4 lg:col-span-2">
          <div className="rounded-2xl bg-[#faedcd] p-4">
            <p className="text-xs font-medium text-muted-foreground">
              Customer
            </p>
            {reservation.user ? (
              <>
                <p className="mt-1 text-sm font-medium">
                  {reservation.user.first_name} {reservation.user.last_name}
                </p>
                <p className="text-xs text-muted-foreground">
                  {reservation.user.email}
                </p>
                {reservation.user.phone_number && (
                  <p className="text-xs text-muted-foreground">
                    {reservation.user.phone_number}
                  </p>
                )}
              </>
            ) : (
              <p className="mt-1 text-sm text-muted-foreground">Guest</p>
            )}
          </div>

          {reservation.special_requests && (
            <div className="rounded-2xl bg-[#faedcd] p-4">
              <p className="text-xs font-medium text-muted-foreground">
                Special Requests
              </p>
              <p className="mt-1 text-sm italic">
                &quot;{reservation.special_requests}&quot;
              </p>
            </div>
          )}

          <div className="rounded-2xl bg-[#faedcd] p-4">
            <p className="text-xs font-medium text-muted-foreground">
              Check-in
            </p>
            <p className="mt-1 text-sm text-muted-foreground">
              {reservation.checked_in_at
                ? new Date(reservation.checked_in_at).toLocaleString("en-NG", {
                    dateStyle: "medium",
                    timeStyle: "short",
                  })
                : "Not checked in"}
            </p>
          </div>
        </div>

        {/* side column */}
        <div className="space-y-4">
          <div className="rounded-2xl bg-[#faedcd] p-4">
            <div className="flex items-center gap-1.5 text-xs font-medium text-muted-foreground">
              <HugeiconsIcon icon={UserGroupIcon} size={14} />
              Party size
            </div>
            <p className="mt-1 text-lg font-semibold">
              {reservation.party_size} guest
              {reservation.party_size !== 1 ? "s" : ""}
            </p>
          </div>

          <div className="rounded-2xl bg-[#faedcd] p-4">
            <div className="flex items-center gap-1.5 text-xs font-medium text-muted-foreground">
              <HugeiconsIcon icon={Calendar01Icon} size={14} />
              Date &amp; Time
            </div>
            <p className="mt-1 text-sm font-medium">{reservation.date}</p>
            <p className="text-sm text-muted-foreground">
              {formatTimeSlot(reservation.time_slot)}
            </p>
          </div>

          {/* table info + link to table detail */}
          <div className="rounded-2xl bg-[#faedcd] p-4">
            <p className="text-xs font-medium text-muted-foreground">Table</p>
            <div className="mt-2 flex items-center justify-between">
              <div>
                <p className="text-sm font-medium">{reservation.table.name}</p>
                <p className="text-xs text-muted-foreground">
                  Seats {reservation.table.capacity}
                  {reservation.table.location
                    ? ` · ${reservation.table.location}`
                    : ""}
                </p>
              </div>
              <Badge
                className={`text-white ${tableStatusStyle(reservation.table.status)}`}
              >
                {reservation.table.status}
              </Badge>
            </div>
            <Link
              href={`/admin/tables/${reservation.table.id}`}
              className="mt-3 inline-block text-xs font-medium text-primary hover:underline"
            >
              View table →
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}
