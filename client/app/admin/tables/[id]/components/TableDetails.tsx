"use client";

import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import { HugeiconsIcon } from "@hugeicons/react";
import {
  ArrowLeft01Icon,
  Edit02Icon,
  Delete02Icon,
  UserGroupIcon,
  Location01Icon,
  QrCodeIcon,
} from "@hugeicons/core-free-icons";

import { useGetTableQuery } from "@/app/state/api/tableApi";
import { useAppDispatch } from "@/app/state/redux";
import {
  openEditTableSheet,
  openDeleteTableDialog,
} from "@/app/state/slices/tableSlice";
import { tableStatusStyle } from "@/app/utils/tableStatusHelpers";
import { statusStyle as reservationStatusStyle } from "@/app/utils/reservationStatusHelper";
import { statusStyle as orderStatusStyle } from "@/app/utils/orderStatusHelpers";
import { formatTimeSlot } from "@/app/utils/timeSlots";
import { formatNaira } from "@/app/utils/formatNaira";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { TableDetailSkeleton } from "./TableSkeleton";

export default function TableDetailContent() {
  const { id } = useParams();
  const router = useRouter();
  const dispatch = useAppDispatch();

  const { data, isLoading, isError } = useGetTableQuery({ id: Number(id) });

  console.log(data);

  if (isLoading) return <TableDetailSkeleton />;

  if (isError || !data?.data) {
    return (
      <div className="flex min-h-[60vh] flex-col items-center justify-center gap-3 text-center">
        <p className="text-lg font-medium">
          This table couldn&apos;t be found.
        </p>
        <Link
          href="/admin/tables"
          className="text-sm text-primary hover:underline"
        >
          Back to tables
        </Link>
      </div>
    );
  }

  const table = data.data;

  return (
    <div className="p-6">
      <button
        onClick={() => router.push("/admin/tables")}
        className="flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground"
      >
        <HugeiconsIcon icon={ArrowLeft01Icon} size={16} />
        Back to tables
      </button>

      <div className="mt-4 flex flex-wrap items-start justify-between gap-4">
        <div>
          <div className="flex items-center gap-3">
            <h1 className="font-serif text-2xl font-semibold">{table.name}</h1>
            <Badge className={`text-white ${tableStatusStyle(table.status)}`}>
              {table.status}
            </Badge>
          </div>
          <p className="mt-1 text-sm text-muted-foreground">#{table.id}</p>
        </div>

        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            onClick={() => dispatch(openEditTableSheet(table.id))}
          >
            <HugeiconsIcon icon={Edit02Icon} strokeWidth={2} size={16} />
            Edit
          </Button>
          <Button
            variant="outline"
            className="bg-red-500 hover:bg-red-600 text-white"
            onClick={() => dispatch(openDeleteTableDialog(table.id))}
          >
            <HugeiconsIcon icon={Delete02Icon} strokeWidth={2} size={16} />
            Delete
          </Button>
        </div>
      </div>

      <div className="mt-6 grid gap-4 sm:grid-cols-3">
        <div className="rounded-2xl  bg-[#faedcd] p-4">
          <div className="flex items-center gap-1.5 text-xs font-medium text-muted-foreground">
            <HugeiconsIcon icon={UserGroupIcon} size={14} />
            Capacity
          </div>
          <p className="mt-1 text-lg font-semibold">{table.capacity} seats</p>
        </div>

        <div className="rounded-2xl  bg-[#faedcd] p-4">
          <div className="flex items-center gap-1.5 text-xs font-medium text-muted-foreground">
            <HugeiconsIcon icon={Location01Icon} size={14} />
            Location
          </div>
          <p className="mt-1 text-lg font-semibold">{table.location || "—"}</p>
        </div>

        {table.qr_code && (
          <div className="rounded-2xl  bg-[#faedcd] p-4">
            <div className="flex items-center gap-1.5 text-xs font-medium text-muted-foreground">
              <HugeiconsIcon icon={QrCodeIcon} size={14} />
              QR Code
            </div>
            <p className="mt-1 truncate text-sm">{table.qr_code}</p>
          </div>
        )}
      </div>

      {/* current reservation, if any */}
      {table.current_reservation && (
        <div className="mt-4 rounded-2xl  bg-[#faedcd] p-4">
          <div className="flex items-center justify-between">
            <p className="text-sm font-medium">Current Reservation</p>
            <span
              className={`rounded-full px-2.5 py-1 text-xs font-medium capitalize ${reservationStatusStyle(table.current_reservation.status)}`}
            >
              {table.current_reservation.status.replace("_", " ")}
            </span>
          </div>
          <p className="mt-2 text-sm text-muted-foreground">
            {table.current_reservation.date} ·{" "}
            {formatTimeSlot(table.current_reservation.time_slot)} ·{" "}
            {table.current_reservation.party_size} guest
            {table.current_reservation.party_size !== 1 ? "s" : ""}
          </p>
          {table.current_reservation.user && (
            <p className="mt-1 text-sm text-muted-foreground">
              {table.current_reservation.user.first_name}{" "}
              {table.current_reservation.user.last_name}
            </p>
          )}
          {table.current_reservation.special_requests && (
            <p className="mt-2 text-xs italic text-muted-foreground">
              &quot;{table.current_reservation.special_requests}&quot;
            </p>
          )}
          <Link
            href={`/admin/reservations/${table.current_reservation.id}`}
            className="mt-3 inline-block text-xs font-medium text-primary hover:underline"
          >
            View reservation →
          </Link>
        </div>
      )}

      {/* current order, if any */}
      {table.current_order && (
        <div className="mt-4 rounded-2xl  bg-[#faedcd] p-4">
          <div className="flex items-center justify-between">
            <p className="text-sm font-medium">Current Order</p>
            <span
              className={`rounded-full px-2.5 py-1 text-xs font-medium capitalize ${orderStatusStyle(table.current_order.status)}`}
            >
              {table.current_order.status}
            </span>
          </div>
          <p className="mt-2 text-sm text-muted-foreground">
            {table.current_order.order_items.length} item
            {table.current_order.order_items.length !== 1 ? "s" : ""} ·{" "}
            {formatNaira(table.current_order.total_amount)}
          </p>
          <Link
            href={`/admin/orders/${table.current_order.id}`}
            className="mt-3 inline-block text-xs font-medium text-primary hover:underline"
          >
            View order →
          </Link>
        </div>
      )}

      {!table.current_reservation && !table.current_order && (
        <div className="mt-4 rounded-2xl  bg-[#faedcd] p-6 text-center text-sm text-muted-foreground">
          No active reservation or order for this table right now.
        </div>
      )}
    </div>
  );
}
