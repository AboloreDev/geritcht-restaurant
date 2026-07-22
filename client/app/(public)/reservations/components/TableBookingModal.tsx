"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { AnimatePresence, motion } from "framer-motion";
import { Minus, Plus, X } from "@mynaui/icons-react";
import { useAppDispatch, useAppSelector, RootState } from "@/app/state/redux";
import {
  closeBookingModal,
  setDate,
  setPartySize,
  setTimeSlot,
} from "@/app/state/slices/reservationSlice";
import { TIME_SLOTS, formatTimeSlot } from "@/app/utils/timeSlots";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Input } from "@/components/ui/input";

export function BookTableModal() {
  const dispatch = useAppDispatch();
  const router = useRouter();
  const { isModalOpen, date, timeSlot, partySize } = useAppSelector(
    (state: RootState) => state.reservation,
  );

  useEffect(() => {
    if (isModalOpen) {
      document.body.style.overflow = "hidden";
      return () => {
        document.body.style.overflow = "";
      };
    }
  }, [isModalOpen]);

  const isValid = Boolean(date && timeSlot && partySize > 0);

  function handleSubmit() {
    if (!isValid) return;

    const params = new URLSearchParams({
      date,
      time_slot: timeSlot,
      party_size: String(partySize),
    });

    dispatch(closeBookingModal());
    router.push(`/reservations?${params.toString()}`);
  }

  return (
    <AnimatePresence>
      {isModalOpen && (
        <>
          {/* overlay */}
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            transition={{ duration: 0.2 }}
            className="fixed inset-0 z-40 bg-black/50 backdrop-blur-sm"
            onClick={() => dispatch(closeBookingModal())}
          />

          {/* modal, centered */}
          <motion.div
            initial={{ opacity: 0, scale: 0.95, y: 12 }}
            animate={{ opacity: 1, scale: 1, y: 0 }}
            exit={{ opacity: 0, scale: 0.95, y: 12 }}
            transition={{ duration: 0.2, ease: [0.16, 1, 0.3, 1] }}
            className="fixed left-1/2 top-1/2 z-50 w-full max-w-md -translate-x-1/2 -translate-y-1/2 rounded-2xl bg-background p-6 shadow-xl"
          >
            <div className="flex items-center justify-between">
              <h2 className="font-serif text-xl text-primary font-medium">
                Book a table
              </h2>
              <button
                onClick={() => dispatch(closeBookingModal())}
                aria-label="Close"
                className="rounded-full cursor-pointer p-1 text-primary hover:bg-muted"
              >
                <X className="h-5 w-5" />
              </button>
            </div>

            <div className="mt-6 space-y-5 text-primary-deep">
              <div>
                <label className="text-sm font-medium">Date</label>
                <Input
                  type="date"
                  value={date}
                  min={new Date().toISOString().split("T")[0]}
                  onChange={(e) => dispatch(setDate(e.target.value))}
                  className="mt-1.5 w-full rounded-lg border px-3 py-2 text-sm"
                />
              </div>

              <div>
                <label className="text-sm font-medium">Time</label>
                <Select
                  value={timeSlot}
                  onValueChange={(value) => dispatch(setTimeSlot(value ?? ""))}
                >
                  <SelectTrigger className="w-full">
                    <SelectValue placeholder="Select a time" />
                  </SelectTrigger>
                  <SelectContent className="bg-white">
                    {TIME_SLOTS.map((slot) => (
                      <SelectItem
                        key={slot}
                        value={slot}
                        className="bg-[#fefae0] cursor-pointer z-100"
                      >
                        {formatTimeSlot(slot)}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              <div>
                <label className="text-sm font-medium">Party size</label>
                <div className="mt-1.5 flex w-fit items-center rounded-lg border border-amber-300">
                  <button
                    onClick={() =>
                      dispatch(setPartySize(Math.max(1, partySize - 1)))
                    }
                    className="flex h-10 w-10 items-center justify-center"
                    aria-label="Decrease party size"
                  >
                    <Minus className="h-4 w-4" />
                  </button>
                  <span className="w-10 text-center text-sm">{partySize}</span>
                  <button
                    onClick={() => dispatch(setPartySize(partySize + 1))}
                    className="flex h-10 w-10 items-center justify-center"
                    aria-label="Increase party size"
                  >
                    <Plus className="h-4 w-4" />
                  </button>
                </div>
              </div>
            </div>

            <Button
              className="mt-8 w-full"
              disabled={!isValid}
              onClick={handleSubmit}
            >
              Check availability
            </Button>
          </motion.div>
        </>
      )}
    </AnimatePresence>
  );
}
