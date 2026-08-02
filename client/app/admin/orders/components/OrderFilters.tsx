"use client";

import { useDebounce } from "@/app/hooks/useDebounce";
import { useGetAllOrdersQuery } from "@/app/state/api/orderApi";
import { RootState, useAppDispatch, useAppSelector } from "@/app/state/redux";
import {
  resetOrdersFilters,
  setFilterDate,
  setFilterStatus,
  setFilterType,
} from "@/app/state/slices/orderSlice";
import { STATUS_OPTIONS, STATUS_TYPES } from "@/app/utils/orderStatusHelpers";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Spinner } from "@/components/ui/spinner";
import { useEffect, useState } from "react";

interface OrderFiltersProps {
  isLoading: boolean;
  isFetching: boolean;
}

const OrderFilters = ({ isLoading, isFetching }: OrderFiltersProps) => {
  const dispatch = useAppDispatch();
  const { page, filterDate, filterStatus, filterType } = useAppSelector(
    (state: RootState) => state.order,
  );

  const hasActiveFilters = Boolean(filterDate || filterStatus || filterType);

  const isRefetchingFilters = isFetching && !isLoading && (page ?? 1) === 1;

  const [dateInput, setDateInput] = useState(filterDate ?? "");
  const debouncedDate = useDebounce(dateInput, 1000);

  useEffect(() => {
    dispatch(setFilterDate(debouncedDate || undefined));
  }, [debouncedDate, dispatch]);

  return (
    <div>
      {/* filter bar */}
      <div className="cursor-pointer flex z-40 flex-wrap rounded-2xl px-6 py-3 items-center gap-3">
        <input
          type="date"
          value={dateInput}
          onChange={(e) => setDateInput(e.target.value)}
          className="cursor-pointer border border-black rounded-full px-3 py-1.5 text-sm"
        />

        <Select
          value={filterStatus ?? "all"}
          onValueChange={(v) =>
            //   @ts-expect-error "<>"
            dispatch(setFilterStatus(v === "all" ? undefined : v))
          }
        >
          <SelectTrigger className="w-40 cursor-pointer">
            <SelectValue placeholder="All statuses" className="" />
          </SelectTrigger>
          <SelectContent className="bg-[#fefae0] cursor-pointer">
            <SelectItem value="all" className="">
              All statuses
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

        <Select
          value={filterType ?? "all"}
          onValueChange={(v) =>
            //   @ts-expect-error "<>"
            dispatch(setFilterType(v === "all" ? undefined : v))
          }
        >
          <SelectTrigger className="w-40 cursor-pointer">
            <SelectValue placeholder="All Types" className="" />
          </SelectTrigger>
          <SelectContent className="bg-[#fefae0] cursor-pointer">
            <SelectItem value="all" className="">
              All Types
            </SelectItem>
            {STATUS_TYPES.map((opt) => (
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
            onClick={() => dispatch(resetOrdersFilters())}
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
    </div>
  );
};

export default OrderFilters;
