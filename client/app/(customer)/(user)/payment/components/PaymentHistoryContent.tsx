"use client";

import Link from "next/link";
import { motion, AnimatePresence } from "framer-motion";
import { HugeiconsIcon } from "@hugeicons/react";
import { CreditCardIcon } from "@hugeicons/core-free-icons";
import { cn } from "@/lib/utils";
import { RootState, useAppDispatch, useAppSelector } from "@/app/state/redux";
import { setPage, setFilterStatus } from "@/app/state/slices/userPaymentSlice";
import { useGetAllUserPaymentQuery } from "@/app/state/api/paymentApi";
import { statusStyle, STATUS_OPTIONS } from "@/app/utils/paymentStatus";
import { formatNaira } from "@/app/utils/formatNaira";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { ArrowLeft } from "@mynaui/icons-react";

export default function PaymentHistoryContent() {
  const dispatch = useAppDispatch();
  const { page, pageSize, filterStatus } = useAppSelector(
    (s: RootState) => s.userPayment,
  );

  const { data, isLoading, isFetching } = useGetAllUserPaymentQuery({
    page,
    page_size: pageSize,
    status: filterStatus,
  });

  const payments = data?.data ?? [];
  const hasMore = data ? page < data.meta.total_pages : false;
  const isRefetching = isFetching && !isLoading && page === 1;

  return (
    <div className="min-h-screen bg-[url('/assets/bg.png')] bg-cover bg-center bg-fixed">
      <div className="mx-auto max-w-3xl px-6 py-10">
        <Link
          href="/"
          className="mt-4 inline-flex items-center gap-1.5 text-sm text-primary transition-colors hover:text-primary-deep"
        >
          <ArrowLeft className="h-4 w-4" />
          Back
        </Link>
        <h1 className="font-serif text-3xl font-semibold text-primary">
          Payment History
        </h1>
        <p className="mt-2 text-sm text-muted-foreground text-primary-deep">
          A record of all your past payments.
        </p>

        {/* <div className="mt-6 flex flex-wrap items-center gap-3">
          <Select
            value={filterStatus ?? "all"}
            onValueChange={(v) =>
              dispatch(setFilterStatus(v === "all" ? undefined : (v as any)))
            }
          >
            <SelectTrigger className="w-40 bg-white">
              <SelectValue placeholder="All statuses" />
            </SelectTrigger>
            <SelectContent className="bg-white">
              <SelectItem value="all">All statuses</SelectItem>
              {STATUS_OPTIONS.map((opt) => (
                <SelectItem
                  key={opt.value}
                  value={opt.value}
                  className="bg-white"
                >
                  {opt.label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>

          {isRefetching && (
            <span className="text-xs text-muted-foreground">Updating…</span>
          )}
        </div> */}

        {isLoading ? (
          <div className="mt-6 space-y-3">
            {Array.from({ length: 5 }).map((_, i) => (
              <div key={i} className="h-20 animate-pulse rounded-xl bg-muted" />
            ))}
          </div>
        ) : payments.length === 0 && !isFetching ? (
          <div className="mt-16 flex flex-col items-center gap-3 text-center">
            <HugeiconsIcon
              icon={CreditCardIcon}
              className="text-muted-foreground"
              size={40}
            />
            <p className="text-sm text-muted-foreground">No payments yet.</p>
          </div>
        ) : (
          <>
            <div
              className={cn(
                "mt-6 space-y-3 transition-opacity",
                isRefetching && "opacity-50",
              )}
            >
              <AnimatePresence initial={false}>
                {payments.map((payment, i) => (
                  <motion.div
                    key={payment.id}
                    initial={{ opacity: 0, y: 8 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ duration: 0.25, delay: i * 0.02 }}
                  >
                    <Link
                      href={`/orders/${payment.order_id}`}
                      className="flex items-center justify-between gap-4 rounded-2xl bg-primary-deep p-4 hover:bg-muted/50"
                    >
                      <div className="min-w-0 flex-1">
                        <p className="font-medium">Order #{payment.order_id}</p>
                        <p className="mt-0.5 truncate font-mono text-xs text-muted-foreground">
                          {payment.reference}
                        </p>
                        <p className="mt-0.5 text-xs text-muted-foreground">
                          {new Date(payment.created_at).toLocaleString(
                            "en-NG",
                            {
                              dateStyle: "medium",
                              timeStyle: "short",
                            },
                          )}
                        </p>
                      </div>

                      <div className="flex shrink-0 items-center gap-3">
                        <span className="font-semibold">
                          {formatNaira(payment.amount)}
                        </span>
                        <span
                          className={cn(
                            "rounded-full px-2.5 py-1 text-xs font-medium capitalize",
                            statusStyle(payment.status),
                          )}
                        >
                          {payment.status}
                        </span>
                      </div>
                    </Link>
                  </motion.div>
                ))}
              </AnimatePresence>
            </div>

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
          </>
        )}
      </div>
    </div>
  );
}
