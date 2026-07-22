"use client";

import React, { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { AnimatePresence, motion } from "framer-motion";
import { X } from "@mynaui/icons-react";
import { formatTimeSlot } from "@/app/utils/timeSlots";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { TableAvailability } from "@/app/state/types/reservationTypes";
import { useAuth } from "@/app/hooks/isAuthenticated";
import { useCreateReservationMutation } from "@/app/state/api/reservationsApi";

export function ConfirmReservationModal({
  table,
  date,
  timeSlot,
  partySize,
  onClose,
}: {
  table: TableAvailability | null;
  date: string;
  timeSlot: string;
  partySize: number;
  onClose: () => void;
}) {
  const router = useRouter();
  const { isAuthenticated } = useAuth();
  const [specialRequests, setSpecialRequests] = useState("");
  const [createReservation, { isLoading, isSuccess, error }] =
    useCreateReservationMutation();

  const isOpen = Boolean(table);

  useEffect(() => {
    if (isOpen) {
      document.body.style.overflow = "hidden";
      return () => {
        document.body.style.overflow = "";
      };
    }
  }, [isOpen]);

  function handleConfirm() {
    if (!table) return;

    if (!isAuthenticated) {
      router.push("/login");
      return;
    }

    createReservation({
      table_id: table.id,
      date,
      time_slot: timeSlot,
      party_size: partySize,
      special_requests: specialRequests || undefined,
    });
  }

  return (
    <AnimatePresence>
      {isOpen && table && (
        <>
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            transition={{ duration: 0.2 }}
            className="fixed inset-0 z-40 bg-/50 backdrop-blur-sm"
            onClick={onClose}
          />

          <motion.div
            initial={{ opacity: 0, scale: 0.95, y: 12 }}
            animate={{ opacity: 1, scale: 1, y: 0 }}
            exit={{ opacity: 0, scale: 0.95, y: 12 }}
            transition={{ duration: 0.2, ease: [0.16, 1, 0.3, 1] }}
            className="fixed left-1/2 top-1/2 z-50 w-full max-w-md -translate-x-1/2 -translate-y-1/2 rounded-2xl bg-background p-6 shadow-xl"
          >
            <div className="flex items-center justify-between">
              <h2 className="font-serif text-xl font-medium text-primary">
                {isSuccess ? "Reservation confirmed" : `Reserve ${table.name}`}
              </h2>
              <button
                onClick={onClose}
                aria-label="Close"
                className="rounded-full p-1 text-white cursor-pointer hover:bg-muted"
              >
                <X className="h-5 w-5 text-white" />
              </button>
            </div>

            {isSuccess ? (
              <div className="mt-6">
                <p className="text-sm text-primary">
                  {table.name} is booked for {partySize}{" "}
                  {partySize === 1 ? "guest" : "guests"} on {date} at{" "}
                  {formatTimeSlot(timeSlot)}.
                </p>
                <Button className="mt-6 w-full" onClick={onClose}>
                  Done
                </Button>
              </div>
            ) : (
              <>
                <div className="mt-4 space-y-1 text-sm text-primary">
                  <p>
                    {date} · {formatTimeSlot(timeSlot)}
                  </p>
                  <p>
                    {partySize} {partySize === 1 ? "guest" : "guests"} · Seats{" "}
                    {table.capacity} · {table.location}
                  </p>
                </div>

                <div className="mt-5">
                  <label className="text-sm text-primary-deep font-medium">
                    Special requests{" "}
                    <span className="text-primary-deep">(optional)</span>
                  </label>
                  <Textarea
                    key={table.id}
                    value={specialRequests}
                    onChange={(e) => setSpecialRequests(e.target.value)}
                    placeholder="Birthday, allergies, seating preference…"
                    className="mt-1.5 text-primary-deep"
                    rows={3}
                  />
                </div>

                {error && (
                  <p className="mt-3 text-sm text-primary-deep">
                    {"data" in error &&
                    typeof error.data === "object" &&
                    error.data &&
                    "message" in error.data
                      ? String((error.data as { message: string }).message)
                      : "Couldn't complete this reservation. Please try another table."}
                  </p>
                )}

                <div className="mt-6 flex gap-3">
                  <Button
                    variant="outline"
                    className="flex-1 bg-[#fefae0] hover:bg-[#fefae0]/80"
                    onClick={onClose}
                  >
                    Cancel
                  </Button>
                  <Button
                    className="flex-1"
                    disabled={isLoading}
                    onClick={handleConfirm}
                  >
                    {isLoading ? "Booking…" : "Create reservation"}
                  </Button>
                </div>
              </>
            )}
          </motion.div>
        </>
      )}
    </AnimatePresence>
  );
}
