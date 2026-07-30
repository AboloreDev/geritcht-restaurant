// src/app/account/reservations/[id]/ReservationDetailsPage.tsx
"use client";

import { useParams } from "next/navigation";
import Link from "next/link";
import { ChevronRight } from "@mynaui/icons-react";
import { useGetReservationByIdQuery } from "@/app/state/api/reservationsApi";
import { formatTimeSlot } from "@/app/utils/timeSlots";
import { statusStyle } from "@/app/utils/reservationStatusHelper";

export default function ReservationDetailsPage() {
  const { id } = useParams();
  const { data, isLoading, isError } = useGetReservationByIdQuery({
    id: Number(id),
  });

  if (isLoading) return <ReservationDetailsSkeleton />;

  // @ts-expect-error "<>"
  if (isError || !data?.data) {
    return (
      <div className="flex min-h-screen flex-col items-center justify-center gap-3 text-center">
        <p className="text-lg font-medium text-primary">
          This reservation couldn&apos;t be found.
        </p>
        <Link
          href="/reservation"
          className="text-sm text-primary hover:underline"
        >
          Back to reservations
        </Link>
      </div>
    );
  }

  // @ts-expect-error "<>"
  const reservation = data.data;

  return (
    <div className="min-h-screen bg-[url('/assets/bg.png')] bg-cover bg-center bg-fixed">
      <div className="mx-auto max-w-2xl px-6 py-10">
        <div className="flex items-center gap-2 text-sm text-primary">
          <Link href="/reservation" className="hover:text-foreground">
            My Reservations
          </Link>
          <ChevronRight className="h-4 w-4" />
          <span className="text-primary-deep">
            Reservation #{reservation.id}
          </span>
        </div>

        <div className="mt-6 flex items-start justify-between gap-3">
          <div>
            <h1 className="font-serif text-2xl font-semibold text-primary">
              {reservation.date} · {formatTimeSlot(reservation.time_slot)}
            </h1>
            <p className="mt-1 text-sm text-primary-deep">
              Booked{" "}
              {new Date(reservation.created_at).toLocaleString("en-NG", {
                dateStyle: "medium",
                timeStyle: "short",
              })}
            </p>
          </div>
          <span
            className={`shrink-0 rounded-full px-3 py-1 text-xs font-bold capitalize ${statusStyle(reservation.status)}`}
          >
            {reservation.status.replace("_", " ")}
          </span>
        </div>

        {/* table details */}
        <div className="mt-8 rounded-2xl border bg-[#fefae0] p-4">
          <p className="text-md font-medium text-muted-foreground">Table:</p>
          <div className="mt-2 flex items-center justify-between text-md">
            <span>{reservation.table.name}</span>
            <span className="text-muted-foreground">
              {reservation.table.location}
            </span>
          </div>
          <p className="mt-1 text-md text-muted-foreground">
            Seats {reservation.table.capacity}
          </p>
        </div>

        {/* party size */}
        <div className="mt-4 rounded-2xl border bg-[#fefae0] p-4">
          <p className="text-md font-medium text-muted-foreground">
            Party size
          </p>
          <p className="mt-1 text-sm">
            {reservation.party_size}{" "}
            {reservation.party_size === 1 ? "guest" : "guests"}
          </p>
        </div>

        {/* special requests */}
        {reservation.special_requests && (
          <div className="mt-4 rounded-2xl border bg-[#fefae0] p-4">
            <p className="text-md font-medium text-muted-foreground">
              Special requests
            </p>
            <p className="mt-1 text-sm italic">
              &quot;{reservation.special_requests}&quot;
            </p>
          </div>
        )}

        {/* check-in status */}
        <div className="mt-4 rounded-2xl border bg-[#fefae0] p-4">
          <p className="text-md font-medium text-muted-foreground">Check-in</p>
          <p className="mt-1 text-sm text-muted-foreground">
            {reservation.checked_in_at
              ? new Date(reservation.checked_in_at).toLocaleString("en-NG", {
                  dateStyle: "medium",
                  timeStyle: "short",
                })
              : "No checked in time"}
          </p>
        </div>
      </div>
    </div>
  );
}

function ReservationDetailsSkeleton() {
  return (
    <div className="mx-auto max-w-2xl px-6 py-10">
      <div className="h-4 w-32 animate-pulse rounded bg-muted" />
      <div className="mt-6 h-8 w-56 animate-pulse rounded bg-muted" />
      <div className="mt-8 space-y-3">
        {Array.from({ length: 5 }).map((_, i) => (
          <div
            key={i}
            className="h-20 animate-pulse rounded-xl bg-primary-deep"
          />
        ))}
      </div>
    </div>
  );
}
