"use client";

import { useDebounce } from "@/app/hooks/useDebounce";
import { useMediaQuery } from "@/app/hooks/useMediaQuery";
import { useSearchOrderQuery } from "@/app/state/api/orderApi";
import { Order } from "@/app/state/types/orderTypes";
import { formatNaira } from "@/app/utils/formatNaira";
import { Input } from "@/components/ui/input";
import { Search, X } from "@mynaui/icons-react";
import { AnimatePresence, motion } from "framer-motion";
import React, { useEffect, useRef, useState } from "react";

const OrderSearch = () => {
  const isMobile = useMediaQuery("(max-width: 767px)");

  const [value, setValue] = useState("");
  const [open, setOpen] = useState(false);

  const wrapRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  const debouncedValue = useDebounce(value, 500);
  const hasQuery = debouncedValue.trim().length > 0;

  const { data, isFetching } = useSearchOrderQuery(
    { q: debouncedValue },
    { skip: !hasQuery },
  );

  const results = hasQuery ? (data?.data ?? []) : [];

  useEffect(() => {
    if (isMobile) return;

    function onClickOutside(e: MouseEvent) {
      if (wrapRef.current && !wrapRef.current.contains(e.target as Node)) {
        setOpen(false);
      }
    }

    document.addEventListener("mousedown", onClickOutside);

    return () => document.removeEventListener("mousedown", onClickOutside);
  }, [isMobile]);

  useEffect(() => {
    if (isMobile && open) {
      document.body.style.overflow = "hidden";

      return () => {
        document.body.style.overflow = "";
      };
    }
  }, [isMobile, open]);

  useEffect(() => {
    if (open) {
      setTimeout(() => inputRef.current?.focus(), isMobile ? 250 : 100);
    }
  }, [open, isMobile]);

  function handleSelect(order: Order) {
    setValue(`Order #${order.id}`);
    setOpen(false);

    document.getElementById("order-grid")?.scrollIntoView({
      behavior: "smooth",
      block: "start",
    });
  }

  function close() {
    setOpen(false);
  }

  const resultsList = (
    <>
      {isFetching && (
        <div className="flex items-center justify-center py-8 text-sm text-muted-foreground">
          Searching orders...
        </div>
      )}

      {!isFetching && results.length === 0 && (
        <div className="flex items-center justify-center py-8 text-sm text-muted-foreground">
          No matching orders found.
        </div>
      )}

      {results.map((item: Order) => (
        <button
          key={item.id}
          onClick={() => handleSelect(item)}
          className="w-full border-b p-4 text-left transition-colors hover:bg-muted last:border-b-0"
        >
          <div className="flex items-start justify-between gap-4">
            <div className="min-w-0 flex-1">
              {/* Order */}
              <div className="flex flex-wrap items-center gap-2">
                <span className="font-semibold text-sm">Order #{item.id}</span>

                <span className="rounded-full bg-primary/10 px-2 py-0.5 text-[10px] font-medium capitalize text-primary">
                  {item.status}
                </span>

                <span className="rounded-full bg-muted px-2 py-0.5 text-[10px] capitalize">
                  {item.type}
                </span>
              </div>

              {/* Customer */}
              {item.user && (
                <p className="mt-2 truncate text-sm font-medium">
                  {item.user.first_name} {item.user.last_name}
                </p>
              )}

              {/* Email */}
              {item.user?.email && (
                <p className="truncate text-xs text-muted-foreground">
                  {item.user.email}
                </p>
              )}

              {/* Payment Reference */}
              {item.payment?.reference && (
                <p className="mt-1 truncate font-mono text-xs text-muted-foreground">
                  Ref: {item.payment.reference}
                </p>
              )}

              {/* Notes */}
              {item.notes && (
                <p className="mt-1 line-clamp-1 text-xs italic text-muted-foreground">
                  {item.notes}
                </p>
              )}
            </div>

            <div className="shrink-0 text-right">
              <p className="font-semibold">{formatNaira(item.total_amount)}</p>

              <p className="mt-1 text-xs text-muted-foreground">
                {new Date(item.created_at).toLocaleDateString()}
              </p>
            </div>
          </div>
        </button>
      ))}
    </>
  );

  return (
    <div ref={wrapRef} className="relative">
      <button
        aria-label="Search Orders"
        onClick={() => setOpen(true)}
        className="flex h-10 w-10 cursor-pointer items-center justify-center rounded-full transition-colors hover:bg-muted"
      >
        <Search className="h-[18px] w-[18px]" />
      </button>

      {/* Desktop */}
      {!isMobile && (
        <>
          {open && (
            <Input
              ref={inputRef}
              value={value}
              onChange={(e) => setValue(e.target.value)}
              placeholder="Search by customer, email, order ID, payment reference, status or notes..."
              className="absolute left-0 top-0 h-10 w-72 rounded-full border-black bg-white text-black"
            />
          )}

          {open && hasQuery && (
            <div className="absolute left-0 top-12 z-30 max-h-96 w-[450px] overflow-y-auto rounded-xl border bg-[#fefae0] shadow-xl">
              {resultsList}
            </div>
          )}
        </>
      )}

      {/* Mobile */}
      {isMobile && (
        <AnimatePresence>
          {open && (
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              className="fixed inset-0 z-50 bg-[#fefae0]"
            >
              <div className="flex items-center gap-3 border-b px-4 py-3">
                <Search className="h-5 w-5 shrink-0" />

                <Input
                  ref={inputRef}
                  value={value}
                  onChange={(e) => setValue(e.target.value)}
                  placeholder="Search by customer, email, order ID, payment reference..."
                  className="flex-1 border-0 bg-transparent shadow-none focus-visible:ring-0"
                />

                <button
                  aria-label="Close search"
                  onClick={close}
                  className="rounded-md p-1 transition-colors hover:bg-muted"
                >
                  <X className="h-5 w-5" />
                </button>
              </div>

              <div className="overflow-y-auto">
                {hasQuery ? (
                  resultsList
                ) : (
                  <div className="p-6 text-center text-sm text-muted-foreground">
                    Search using customer name, email, payment reference, order
                    status, notes, or order ID.
                  </div>
                )}
              </div>
            </motion.div>
          )}
        </AnimatePresence>
      )}
    </div>
  );
};

export default OrderSearch;
