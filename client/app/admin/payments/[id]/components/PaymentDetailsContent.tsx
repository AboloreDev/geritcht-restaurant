"use client";

import { useState } from "react";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import { HugeiconsIcon } from "@hugeicons/react";
import {
  ArrowLeft01Icon,
  RefreshIcon,
  MoneyBag02Icon,
} from "@hugeicons/core-free-icons";

import {
  useGetPaymentDetailsQuery,
  useReVerifyPaymentMutation,
  useProcessRefundMutation,
} from "@/app/state/api/paymentApi";
import { statusStyle } from "@/app/utils/paymentStatus";
import { formatNaira } from "@/app/utils/formatNaira";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { Textarea } from "@/components/ui/textarea";
import { toast } from "sonner";
import { getApiError } from "@/app/utils/apiError";

export default function PaymentDetailContent() {
  const { id } = useParams();
  const router = useRouter();
  const paymentId = Number(id);

  const { data, isLoading, isError } = useGetPaymentDetailsQuery({
    id: paymentId,
  });

  const [reVerify, { isLoading: isVerifying }] = useReVerifyPaymentMutation();
  const [processRefund, { isLoading: isRefunding }] =
    useProcessRefundMutation();

  const [confirmingRefund, setConfirmingRefund] = useState(false);
  const [refundNotes, setRefundNotes] = useState("");

  if (isLoading) return <PaymentDetailSkeleton />;

  // @ts-expect-error "<>"
  if (isError || !data?.data) {
    return (
      <div className="flex min-h-[60vh] flex-col items-center justify-center gap-3 text-center">
        <p className="text-lg font-medium">
          This payment couldn&apos;t be found.
        </p>
        <Link
          href="/admin/payments"
          className="text-sm text-primary hover:underline"
        >
          Back to payments
        </Link>
      </div>
    );
  }

  // @ts-expect-error "<>"
  const payment = data.data;
  console.log(payment);
  const canRefund = payment.status === "paid";
  const canReVerify =
    payment.status === "pending" || payment.status === "unpaid";

  async function handleReVerify() {
    try {
      const response = await reVerify(payment.reference).unwrap();
      toast.success(`Payment status: ${response.data.status}`);
    } catch (err) {
      toast.error(getApiError(err));
    }
  }

  async function handleRefund() {
    try {
      const response = await processRefund({
        orderId: payment.order_id,
        notes: refundNotes || undefined,
      }).unwrap();
      toast.success(response.message);
      setConfirmingRefund(false);
      setRefundNotes("");
    } catch (err) {
      toast.error(getApiError(err));
    }
  }

  return (
    <div className="p-6">
      <button
        onClick={() => router.push("/admin/payments")}
        className="flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground"
      >
        <HugeiconsIcon icon={ArrowLeft01Icon} size={16} />
        Back to payments
      </button>

      <div className="mt-4 flex flex-wrap items-start justify-between gap-4">
        <div>
          <div className="flex items-center gap-3">
            <h1 className="font-serif text-2xl font-semibold">
              {formatNaira(payment.amount)}
            </h1>
            <Badge className={statusStyle(payment.status)}>
              {payment.status}
            </Badge>
          </div>
          <p className="mt-1 font-mono text-sm text-muted-foreground">
            {payment.reference}
          </p>
        </div>

        <div className="flex items-center gap-2">
          {canReVerify && (
            <Button
              variant="outline"
              disabled={isVerifying}
              onClick={handleReVerify}
              className="border-none hover:text-slate-500"
            >
              <HugeiconsIcon icon={RefreshIcon} strokeWidth={2} size={16} />
              {isVerifying ? "Verifying…" : "Re-verify"}
            </Button>
          )}
          {canRefund && (
            <Button
              variant="outline"
              className="text-red-500 border-none"
              onClick={() => setConfirmingRefund(true)}
            >
              <HugeiconsIcon icon={MoneyBag02Icon} strokeWidth={2} size={16} />
              Refund
            </Button>
          )}
        </div>
      </div>

      <div className="mt-6 grid gap-4 sm:grid-cols-2">
        <div className="rounded-2xl bg-[#faedcd] p-4">
          <p className="text-xs font-medium text-muted-foreground">Order</p>
          <Link
            href={`/admin/orders/${payment.order_id}`}
            className="mt-1 block text-sm font-medium text-primary hover:underline"
          >
            Order #{payment.order_id} →
          </Link>
        </div>

        <div className="rounded-2xl bg-[#faedcd] p-4">
          <p className="text-xs font-medium text-muted-foreground">Provider</p>
          <p className="mt-1 text-sm font-medium capitalize">
            {payment.provider}
          </p>
          {payment.provider_reference && (
            <p className="mt-0.5 truncate font-mono text-xs text-muted-foreground">
              {payment.provider_reference}
            </p>
          )}
        </div>

        <div className="rounded-2xl bg-[#faedcd] p-4">
          <p className="text-xs font-medium text-muted-foreground">Currency</p>
          <p className="mt-1 text-sm font-medium">{payment.currency}</p>
        </div>

        <div className="rounded-2xl bg-[#faedcd] p-4">
          <p className="text-xs font-medium text-muted-foreground">Paid At</p>
          <p className="mt-1 text-sm font-medium">
            {payment.paid_at
              ? new Date(payment.paid_at).toLocaleString("en-NG", {
                  dateStyle: "medium",
                  timeStyle: "short",
                })
              : "Not yet paid"}
          </p>
        </div>

        {payment.failure_reason && (
          <div className="rounded-xl border border-red-200 bg-red-50 p-4 sm:col-span-2">
            <p className="text-xs font-medium text-red-700">Failure Reason</p>
            <p className="mt-1 text-sm text-red-900">
              {payment.failure_reason}
            </p>
          </div>
        )}

        <div className="rounded-2xl bg-[#faedcd] p-4 text-xs text-muted-foreground sm:col-span-2">
          Created{" "}
          {new Date(payment.created_at).toLocaleString("en-NG", {
            dateStyle: "medium",
            timeStyle: "short",
          })}
        </div>
      </div>

      {/* refund confirmation */}
      <AlertDialog open={confirmingRefund} onOpenChange={setConfirmingRefund}>
        <AlertDialogContent className="bg-[#faedcd] max-w-4xl">
          <AlertDialogHeader>
            <AlertDialogTitle>
              Refund {formatNaira(payment.amount)}?
            </AlertDialogTitle>
            <AlertDialogDescription>
              This will process an actual refund via Paystack and cancel the
              associated order. This cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>

          <div className="py-2">
            <Textarea
              value={refundNotes}
              onChange={(e) => setRefundNotes(e.target.value)}
              placeholder="Reason for refund (optional)"
              rows={3}
              className="placeholder:text-slate-500"
            />
          </div>

          <AlertDialogFooter>
            <AlertDialogCancel className="bg-red-500 text-white hover:bg-red-400">
              Cancel
            </AlertDialogCancel>
            <AlertDialogAction
              disabled={isRefunding}
              onClick={handleRefund}
              className="bg-white text-black hover:bg-white/70"
            >
              {isRefunding ? "Processing…" : "Confirm Refund"}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}

function PaymentDetailSkeleton() {
  return (
    <div className="p-6">
      <div className="h-4 w-32 animate-pulse rounded bg-[#faedcd]" />
      <div className="mt-4 h-8 w-48 animate-pulse rounded bg-[#faedcd]" />
      <div className="mt-6 grid gap-4 sm:grid-cols-2">
        {Array.from({ length: 4 }).map((_, i) => (
          <div key={i} className="h-20 animate-pulse rounded-xl bg-[#faedcd]" />
        ))}
      </div>
    </div>
  );
}
