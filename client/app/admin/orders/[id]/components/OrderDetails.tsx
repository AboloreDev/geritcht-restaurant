"use client";

import {
  useAdminCancelOrderMutation,
  useGetOrderByIdQuery,
} from "@/app/state/api/orderApi";
import { getApiError } from "@/app/utils/apiError";
import { useParams } from "next/navigation";
import { useEffect, useState } from "react";
import { toast } from "sonner";
import { OrderDetailSkeleton } from "./OrderDetailsSkeleton";
import Link from "next/link";
import { ChevronRight, Radio } from "@mynaui/icons-react";
import { motion } from "framer-motion";
import { statusStyle } from "@/app/utils/orderStatusHelpers";
import { Button } from "@/components/ui/button";
import Image from "next/image";
import { formatNaira } from "@/app/utils/formatNaira";

export default function OrderDetailAdmin() {
  const { id } = useParams();
  const orderId = Number(id);

  const { data, isLoading, isError, refetch } = useGetOrderByIdQuery({
    id: orderId,
  });
  const [cancelOrder, { isLoading: isCancelling }] =
    useAdminCancelOrderMutation();

  // live status via websocket — assumes NEXT_PUBLIC_WS_URL is set,
  // and the backend pushes { order_id, status } on this channel
  const [liveStatus, setLiveStatus] = useState<string | null>(null);
  const [wsConnected, setWsConnected] = useState(false);

  useEffect(() => {
    const httpUrl = process.env.NEXT_PUBLIC_API_BASE_URL!;

    const wsUrl = httpUrl
      .replace("http://", "ws://")
      .replace("https://", "wss://");

    console.log("Connecting to websocket for order", wsUrl, orderId);

    const ws = new WebSocket(`${wsUrl}/ws/orders/${orderId}`);
    console.log("Connecting to websocket for order", ws);

    ws.onopen = () => setWsConnected(true);
    ws.onclose = () => setWsConnected(false);
    ws.onmessage = (event) => {
      try {
        const payload = JSON.parse(event.data);
        if (payload.status) {
          setLiveStatus(payload.status);
          refetch(); // pull full order again so items/totals stay in sync too
        }
      } catch {
        // ignore malformed frames
      }
    };

    return () => ws.close();
  }, [orderId, refetch]);

  async function handleCancel() {
    try {
      await cancelOrder({ id: orderId }).unwrap();
      toast.success("Order cancelled");
    } catch (err) {
      toast.error(getApiError(err));
    }
  }

  if (isLoading) return <OrderDetailSkeleton />;

  if (isError || !data?.data) {
    return (
      <div className="flex min-h-[60vh] flex-col items-center justify-center gap-3 text-center">
        <p className="text-lg font-medium">
          This order couldn&apos;t be found.
        </p>
        <Link
          href="/admin/orders"
          className="text-sm text-primary hover:underline"
        >
          Back to orders
        </Link>
      </div>
    );
  }

  const order = data.data;
  const currentStatus = liveStatus ?? order.status;
  const canCancel =
    currentStatus !== "cancelled" && currentStatus !== "completed";

  return (
    <div className="p-6">
      {/* breadcrumb */}
      <div className="flex items-center gap-2 text-sm text-muted-foreground">
        <Link href="/admin/orders" className="hover:text-foreground">
          Orders
        </Link>
        <ChevronRight className="h-4 w-4" />
        <span className="text-foreground">Order #{order.id}</span>
      </div>

      {/* header */}
      <div className="mt-4 flex flex-wrap items-start justify-between gap-4">
        <div>
          <h1 className="font-serif text-2xl font-semibold">
            Order #{order.id}
          </h1>
          <p className="mt-1 text-sm text-muted-foreground">
            {new Date(order.created_at).toLocaleString("en-NG", {
              dateStyle: "medium",
              timeStyle: "short",
            })}
          </p>
        </div>

        <div className="flex items-center gap-3">
          <span className="flex items-center gap-1.5 text-xs text-muted-foreground">
            <span
              className={`h-1.5 w-1.5 rounded-full ${wsConnected ? "bg-emerald-500" : "bg-muted-foreground"}`}
            />
            {wsConnected ? "Live" : "Offline"}
          </span>

          <motion.span
            key={currentStatus}
            initial={{ scale: 0.9, opacity: 0 }}
            animate={{ scale: 1, opacity: 1 }}
            className={`rounded-full px-3 py-1 text-xs font-medium capitalize ${statusStyle(currentStatus)}`}
          >
            {currentStatus}
          </motion.span>

          {canCancel && (
            <Button
              size="sm"
              variant="outline"
              disabled={isCancelling}
              onClick={handleCancel}
              className="text-white bg-red-500 hover:bg-red-400 border-none"
            >
              {isCancelling ? "Cancelling…" : "Cancel order"}
            </Button>
          )}
        </div>
      </div>

      {/* live status indicator note */}
      <div className="mt-4 flex items-center gap-2 rounded-lg border bg-muted/40 p-3 text-xs text-muted-foreground">
        <Radio className="h-3.5 w-3.5" />
        Status updates automatically via the order&apos;s real-time channel — no
        manual refresh needed.
      </div>

      <div className="mt-6 grid gap-6 lg:grid-cols-3">
        {/* items + notes — main column */}
        <div className="lg:col-span-2 space-y-4">
          <div className="rounded-2xl bg-[#faedcd]">
            {order.order_items.map((item, i) => (
              <div
                key={item.id}
                className={`flex gap-3 p-4 ${i !== order.order_items.length - 1 ? "border-b" : ""}`}
              >
                <div className="relative h-14 w-14 shrink-0 overflow-hidden rounded-lg bg-muted">
                  {item.menu_item.images?.[0] && (
                    <Image
                      src={
                        item.menu_item.images[0].alt_text ||
                        item.menu_item.images[0].url
                      }
                      alt={item.menu_item.name}
                      fill
                      className="object-cover"
                      sizes="56px"
                    />
                  )}
                </div>
                <div className="flex flex-1 items-start justify-between">
                  <div>
                    <p className="text-sm font-medium">{item.menu_item.name}</p>
                    <p className="mt-0.5 text-xs text-muted-foreground">
                      Qty {item.quantity} · {formatNaira(item.price)} each
                    </p>
                    {item.special_instructions && (
                      <p className="mt-1 text-xs italic text-muted-foreground">
                        &quot;{item.special_instructions}&quot;
                      </p>
                    )}
                  </div>
                  <span className="text-sm font-medium">
                    {formatNaira(item.price * item.quantity)}
                  </span>
                </div>
              </div>
            ))}
          </div>

          {order.notes && (
            <div className="rounded-2xl bg-[#faedcd] p-4">
              <p className="text-xs font-medium text-muted-foreground">Notes</p>
              <p className="mt-1 text-sm">{order.notes}</p>
            </div>
          )}
        </div>

        {/* customer + payment — side column */}
        <div className="space-y-4">
          <div className="rounded-2xl bg-[#faedcd] p-4">
            <p className="text-xs font-medium text-muted-foreground">
              Customer
            </p>
            <p className="mt-1 text-sm font-medium">
              {order.user?.first_name} {order.user?.last_name}
            </p>
            <p className="text-xs text-muted-foreground">{order.user?.email}</p>
            {order.user?.phone_number && (
              <p className="text-xs text-muted-foreground">
                {order.user.phone_number}
              </p>
            )}
          </div>

          <div className="rounded-2xl bg-[#faedcd] p-4">
            <p className="text-xs font-medium text-muted-foreground">Payment</p>
            <div className="mt-2 flex items-center justify-between text-sm">
              <span className="capitalize text-muted-foreground">
                {order.payment?.provider ?? "—"} · {order.payment_status}
              </span>
            </div>
            {order.payment?.reference && (
              <p className="mt-1 font-mono text-xs text-muted-foreground">
                {order.payment.reference}
              </p>
            )}
          </div>

          <div className="rounded-2xl bg-[#faedcd] p-4">
            <div className="flex items-center justify-between">
              <span className="text-sm font-medium">Total</span>
              <span className="text-lg font-semibold">
                {formatNaira(order.total_amount)}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
