"use client";

import { useEffect, useState } from "react";
import { useDebounce } from "@/app/hooks/useDebounce";
import { RootState, useAppDispatch, useAppSelector } from "@/app/state/redux";
import {
  resetReservationFilters,
  setFilterDate,
  setFilterStatus,
  setFilterTimeSlot,
} from "@/app/state/slices/reservationSlice";
import { STATUS_OPTIONS } from "@/app/utils/reservationStatusHelper";
import { TIME_SLOTS, formatTimeSlot } from "@/app/utils/timeSlots";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Spinner } from "@/components/ui/spinner";

interface ReservationFilterProps {
  isLoading: boolean;
  isFetching: boolean;
}

export default function ReservationFilters({
  isLoading,
  isFetching,
}: ReservationFilterProps) {
  const dispatch = useAppDispatch();
  const { page, filterDate, filterStatus, filterTimeSlot } = useAppSelector(
    (state: RootState) => state.reservation,
  );

  const hasActiveFilters = Boolean(
    filterDate || filterStatus || filterTimeSlot,
  );
  const isRefetchingFilters = isFetching && !isLoading && (page ?? 1) === 1;

  const [dateInput, setDateInput] = useState(filterDate ?? "");
  const debouncedDate = useDebounce(dateInput, 1000);

  useEffect(() => {
    dispatch(setFilterDate(debouncedDate || undefined));
  }, [debouncedDate, dispatch]);

  return (
    <div className="flex flex-wrap items-center gap-3 rounded-2xl px-6 py-3">
      <input
        type="date"
        value={dateInput}
        onChange={(e) => setDateInput(e.target.value)}
        className="cursor-pointer rounded-full border border-black px-3 py-1.5 text-sm"
      />

      <Select
        value={filterTimeSlot ?? "all"}
        onValueChange={(v) =>
          // @ts-expect-error "<>"
          dispatch(setFilterTimeSlot(v === "all" ? undefined : v))
        }
      >
        <SelectTrigger className="w-40 cursor-pointer">
          <SelectValue placeholder="All times" />
        </SelectTrigger>
        <SelectContent className="bg-[#fefae0]">
          <SelectItem value="all">All times</SelectItem>
          {TIME_SLOTS.map((slot) => (
            <SelectItem className="cursor-pointer" key={slot} value={slot}>
              {formatTimeSlot(slot)}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>

      <Select
        value={filterStatus ?? "all"}
        onValueChange={(v) =>
          // @ts-expect-error "<>"
          dispatch(setFilterStatus(v === "all" ? undefined : v))
        }
      >
        <SelectTrigger className="w-40 cursor-pointer">
          <SelectValue placeholder="All statuses" />
        </SelectTrigger>
        <SelectContent className="bg-[#fefae0]">
          <SelectItem value="all">All statuses</SelectItem>
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
          onClick={() => {
            setDateInput("");
            dispatch(resetReservationFilters());
          }}
          className="text-xs text-muted-foreground hover:text-foreground"
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
  );
}
