"use client";

import { Suspense, useEffect, useState } from "react";
import { useSearchParams, useRouter } from "next/navigation";
import { useLazyVerifyPaymentQuery } from "@/app/state/api/paymentApi";

const MAX_POLL_ATTEMPTS = 5;
const POLL_INTERVAL_MS = 2500;

function CheckoutCallbackContent() {
  const searchParams = useSearchParams();
  const router = useRouter();
  const reference = searchParams.get("reference");

  const [verify, { data }] = useLazyVerifyPaymentQuery();
  const [attempts, setAttempts] = useState(0);
  const [outcome, setOutcome] = useState<
    "checking" | "confirmed" | "pending" | "failed"
  >("checking");

  useEffect(() => {
    if (!reference) {
      // eslint-disable-next-line react-hooks/set-state-in-effect
      setOutcome("failed");
      return;
    }

    let cancelled = false;

    async function poll() {
      const result = await verify(reference!)
        .unwrap()
        .catch(() => null);

      if (cancelled) return;

      if (result?.data.status && result.data.status !== "success") {
        setOutcome("confirmed");
        return;
      }

      if (result?.data.status === "failed") {
        setOutcome("failed");
        return;
      }

      setAttempts((a) => {
        const next = a + 1;
        if (next >= MAX_POLL_ATTEMPTS) {
          setOutcome("pending");
        } else {
          setTimeout(poll, POLL_INTERVAL_MS);
        }
        return next;
      });
    }

    poll();
    return () => {
      cancelled = true;
    };
  }, [reference, verify]);

  if (outcome === "checking") {
    return (
      <div className="flex min-h-screen flex-col items-center justify-center gap-3">
        <div className="h-8 w-8 animate-spin rounded-full border-2 border-primary border-t-transparent" />
        <p className="text-sm text-muted-foreground">
          Confirming your payment…
        </p>
      </div>
    );
  }

  if (outcome === "confirmed") {
    return (
      <div className="flex min-h-screen flex-col items-center justify-center gap-4 text-center">
        <p className="text-xl font-semibold">Order confirmed 🎉</p>
        <button
          onClick={() => router.push("/account/orders")}
          className="text-sm text-primary hover:underline"
        >
          View your orders
        </button>
      </div>
    );
  }

  if (outcome === "pending") {
    return (
      <div className="flex min-h-screen flex-col items-center justify-center gap-3 text-center">
        <p className="text-lg font-medium">Still confirming your payment…</p>
        <p className="text-sm text-muted-foreground">
          This can take a moment. Check your orders shortly.
        </p>
        <button
          onClick={() => router.push("/account/orders")}
          className="text-sm text-primary hover:underline"
        >
          View your orders
        </button>
      </div>
    );
  }

  return (
    <div className="flex min-h-screen flex-col items-center justify-center gap-3 text-center">
      <p className="text-lg font-medium text-destructive">
        Payment couldn&apos;t be confirmed.
      </p>
      <button
        onClick={() => router.push("/menu")}
        className="text-sm text-primary hover:underline"
      >
        Back to menu
      </button>
    </div>
  );
}

export default function CheckoutCallbackPage() {
  return (
    <Suspense fallback={null}>
      <CheckoutCallbackContent />
    </Suspense>
  );
}
