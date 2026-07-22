"use client";

import { useCheckAvailabilityQuery } from "@/app/state/api/reservationsApi";
import { TableAvailability } from "@/app/state/types/reservationTypes";
import { formatTimeSlot } from "@/app/utils/timeSlots";
import { useRouter, useSearchParams } from "next/navigation";
import React, { useState } from "react";
import { ConfirmReservationModal } from "./ConfirmReservationModal";
import { StatusLegend } from "./ReservationStausLegend";
import { TableCard } from "./TableCard";
import { motion } from "framer-motion";
import Link from "next/link";
import { ChevronRight } from "@mynaui/icons-react";

const Reservations = () => {
  const searchParams = useSearchParams();
  const router = useRouter();

  const date = searchParams.get("date") ?? "";
  const timeSlot = searchParams.get("time_slot") ?? "";
  const partySize = Number(searchParams.get("party_size") ?? 0);

  const [selectedTable, setSelectedTable] = useState<TableAvailability | null>(
    null,
  );

  const { data, isLoading, isError } = useCheckAvailabilityQuery(
    { date, time_slot: timeSlot, party_size: partySize },
    { skip: !date || !timeSlot || !partySize },
  );

  if (!date || !timeSlot || !partySize) {
    return (
      <div className="flex min-h-screen flex-col items-center justify-center gap-3 text-center px-6">
        <p className="text-lg font-medium">No table criteria found.</p>
        <button
          onClick={() => router.push("/")}
          className="text-sm text-primary hover:underline"
        >
          Go back and search for a table
        </button>
      </div>
    );
  }

  const tables = data?.data.tables ?? [];

  return (
    <div className="min-h-screen bg-[url('/assets/bg.png')] bg-cover bg-center bg-fixed">
      <div className="mx-auto max-w-5xl px-6 py-10">
        <div className="text-center">
          <h1 className="font-serif text-3xl text-primary font-semibold">
            Available tables
          </h1>
          <p className="mt-2 text-sm text-primary-deep">
            {partySize} {partySize === 1 ? "guest" : "guests"} · {date} ·{" "}
            {formatTimeSlot(timeSlot)}
          </p>
        </div>

        {isLoading && (
          <div className="mt-10 grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4">
            {Array.from({ length: 8 }).map((_, i) => (
              <div
                key={i}
                className="h-32 animate-pulse rounded-xl bg-primary-deep"
              />
            ))}
          </div>
        )}

        <div className="mt-5 flex items-center justify-center">
          <div className="w-10" aria-hidden />

          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.45 }}
            className="flex items-center justify-center gap-2 text-sm text-text-secondary"
          >
            <Link
              href="/"
              className="transition-colors hover:text-primary-deep"
            >
              Home
            </Link>

            <ChevronRight className="h-4 w-4" />

            <span className="text-primary">Tables</span>
          </motion.div>
        </div>

        {isError && (
          <p className="mt-10 text-center text-sm text-destructive">
            Couldn&apos;t load tables. Please try again.
          </p>
        )}

        {!isLoading && !isError && tables.length === 0 && (
          <p className="mt-10 text-center text-sm text-primary-deep">
            No tables found for this search.
          </p>
        )}

        {!isLoading && tables.length > 0 && (
          <>
            <div className="mt-10 grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4">
              {tables.map((table) => (
                <TableCard
                  key={table.id}
                  table={table}
                  onSelect={setSelectedTable}
                />
              ))}
            </div>
            <StatusLegend />
          </>
        )}
      </div>

      <ConfirmReservationModal
        table={selectedTable}
        date={date}
        timeSlot={timeSlot}
        partySize={partySize}
        onClose={() => setSelectedTable(null)}
      />
    </div>
  );
};

export default Reservations;
