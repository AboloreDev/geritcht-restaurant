"use client";

import { useEffect, useState } from "react";
import { useDebounce } from "@/app/hooks/useDebounce";
import { RootState, useAppDispatch, useAppSelector } from "@/app/state/redux";
import {
  setFilterStatus,
  setSearch,
  resetFilters,
} from "@/app/state/slices/paymentSlice";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { HugeiconsIcon } from "@hugeicons/react";
import { Search01Icon } from "@hugeicons/core-free-icons";
import { STATUS_OPTIONS } from "@/app/utils/paymentStatus";
import { useAdminGetAllPaymentQuery } from "@/app/state/api/paymentApi";

export default function PaymentFilters() {
  const dispatch = useAppDispatch();
  const { filterStatus, query, page, pageSize } = useAppSelector(
    (s: RootState) => s.payment,
  );

  const { data, isFetching, isLoading } = useAdminGetAllPaymentQuery({
    page,
    page_size: pageSize,
    status: filterStatus,
    reference: query || undefined,
  });

  const isRefetching = isFetching && !isLoading && page === 1;

  const [searchInput, setSearchInput] = useState(query);
  const debouncedSearch = useDebounce(searchInput, 400);

  useEffect(() => {
    dispatch(setSearch(debouncedSearch));
  }, [debouncedSearch, dispatch]);

  const hasActiveFilters = Boolean(filterStatus || query);

  return (
    <div className="flex flex-wrap items-center gap-3 rounded-2xl px-6 py-3">
      <div className="relative min-w-[200px] flex-1">
        <HugeiconsIcon
          icon={Search01Icon}
          size={16}
          className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground"
        />
        <input
          type="text"
          value={searchInput}
          onChange={(e) => setSearchInput(e.target.value)}
          placeholder="Search by reference…"
          className="w-full rounded-full border border-black py-1.5 pl-9 pr-3 text-sm"
        />
      </div>

      <Select
        value={filterStatus ?? "all"}
        onValueChange={(v) =>
          // @ts-expect-error "<>"
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
      {hasActiveFilters && (
        <button
          onClick={() => {
            setSearchInput("");
            dispatch(resetFilters());
          }}
          className="text-xs text-muted-foreground hover:text-foreground"
        >
          Clear filters
        </button>
      )}
    </div>
  );
}
