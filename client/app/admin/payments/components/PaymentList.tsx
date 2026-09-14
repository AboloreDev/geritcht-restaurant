"use client";

import Link from "next/link";
import { motion, AnimatePresence } from "framer-motion";
import { HugeiconsIcon } from "@hugeicons/react";
import { CreditCardIcon, Loader } from "@hugeicons/core-free-icons";
import { cn } from "@/lib/utils";
import { RootState, useAppDispatch, useAppSelector } from "@/app/state/redux";
import { setPage } from "@/app/state/slices/paymentSlice";
import { useAdminGetAllPaymentQuery } from "@/app/state/api/paymentApi";
import { formatNaira } from "@/app/utils/formatNaira";
import { Button } from "@/components/ui/button";
import { statusStyle } from "@/app/utils/paymentStatus";
import { Spinner } from "@mynaui/icons-react";

export default function PaymentsList() {
  const dispatch = useAppDispatch();
  const { page, pageSize, filterStatus, query } = useAppSelector(
    (s: RootState) => s.payment,
  );

  const { data, isLoading, isFetching } = useAdminGetAllPaymentQuery({
    page,
    page_size: pageSize,
    status: filterStatus,
    reference: query || undefined,
  });

  const payments = data?.data ?? [];
  const hasMore = data ? page < data.meta.total_pages : false;

  // refetching because a filter changed — always resets to page 1,
  // so this never fires for a load-more (page > 1) fetch
  const isRefetching = isFetching && !isLoading && page === 1;

  if (isLoading) {
    return (
      <div className="mt-6 space-y-3">
        {Array.from({ length: 8 }).map((_, i) => (
          <div key={i} className="h-20 animate-pulse rounded-xl bg-[#d4a373]" />
        ))}
      </div>
    );
  }

  if (!isFetching && payments.length === 0) {
    return (
      <div className="mt-16 flex flex-col items-center gap-3 text-center">
        <HugeiconsIcon
          icon={CreditCardIcon}
          className="text-muted-foreground"
          size={40}
        />
        <p className="text-sm text-muted-foreground">
          No payments match this filter.
        </p>
      </div>
    );
  }

  return (
    <div className="relative">
      {/* status pill while refetching due to filter change */}
      <AnimatePresence>
        {isRefetching && (
          <motion.div
            initial={{ opacity: 0, y: -6 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -6 }}
            transition={{ duration: 0.15 }}
            className="mt-4 flex items-center justify-center gap-2 text-sm text-muted-foreground"
          >
            <Spinner className="h-4 w-4 animate-spin" />
            Updating…
          </motion.div>
        )}
      </AnimatePresence>

      <motion.div
        layout
        className={cn(
          "mt-6 space-y-3 transition-opacity duration-300",
          isRefetching && "opacity-50",
        )}
      >
        <AnimatePresence initial={false} mode="popLayout">
          {payments.map((payment, index) => (
            <motion.div
              key={payment.id}
              layout
              initial={{ opacity: 0, y: 8 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, scale: 0.98 }}
              transition={{
                duration: 0.25,
                delay: index * 0.02,
                ease: [0.16, 1, 0.3, 1],
              }}
            >
              <Link
                href={`/admin/payments/${payment.id}`}
                className="flex items-center justify-between gap-4 rounded-2xl bg-[#faedcd] p-4 transition-colors hover:bg-[#fefae0]"
              >
                <div className="min-w-0 flex-1">
                  <div className="flex items-center gap-2">
                    <p className="font-medium">Order #{payment.order_id}</p>
                    <span className="text-xs capitalize text-muted-foreground">
                      • {payment.provider}
                    </span>
                  </div>
                  <p className="mt-1 truncate font-mono text-xs text-muted-foreground">
                    {payment.reference}
                  </p>
                  <p className="mt-0.5 text-xs text-muted-foreground">
                    {new Date(payment.created_at).toLocaleString("en-NG", {
                      dateStyle: "medium",
                      timeStyle: "short",
                    })}
                  </p>
                </div>

                <div className="flex shrink-0 items-center gap-3">
                  <span className="font-semibold">
                    {formatNaira(payment.amount)}
                  </span>
                  <motion.span
                    key={payment.status}
                    initial={{ scale: 0.85, opacity: 0 }}
                    animate={{ scale: 1, opacity: 1 }}
                    transition={{ duration: 0.2 }}
                    className={cn(
                      "rounded-full px-2.5 py-1 text-xs font-medium capitalize",
                      statusStyle(payment.status),
                    )}
                  >
                    {payment.status}
                  </motion.span>
                </div>
              </Link>
            </motion.div>
          ))}
        </AnimatePresence>
      </motion.div>

      {hasMore && (
        <div className="mt-6 flex justify-center">
          <Button
            variant="outline"
            disabled={isFetching}
            onClick={() => dispatch(setPage(page + 1))}
          >
            {isFetching && page > 1 ? "Loading…" : "Load more"}
          </Button>
        </div>
      )}
    </div>
  );
}
