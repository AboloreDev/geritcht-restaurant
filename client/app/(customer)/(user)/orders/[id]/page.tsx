// src/app/account/orders/[id]/OrderDetailsPage.tsx
"use client";

import { useParams } from "next/navigation";
import Image from "next/image";
import Link from "next/link";
import { ChevronRight } from "@mynaui/icons-react";
import { useGetOrderByIdQuery } from "@/app/state/api/orderApi";
import { formatNaira } from "@/app/utils/formatNaira";
import { statusStyle } from "@/app/utils/orderStatusHelpers";

export default function OrderDetailsPage() {
  const { id } = useParams();
  const { data, isLoading, isError } = useGetOrderByIdQuery({ id: Number(id) });

  console.log(data);
  if (isLoading) return <OrderDetailsSkeleton />;

  if (isError || !data?.data) {
    return (
      <div className="flex min-h-screen flex-col items-center justify-center gap-3 text-center">
        <p className="text-lg font-medium">
          This order couldn&apos;t be found.
        </p>
        <Link
          href="/account/orders"
          className="text-sm text-primary hover:underline"
        >
          Back to orders
        </Link>
      </div>
    );
  }

  const order = data.data;

  return (
    <div className="min-h-screen bg-[url('/assets/bg.png')] bg-cover bg-center bg-fixed">
      <div className="mx-auto max-w-2xl px-6 py-10">
        <div className="flex items-center gap-2 text-sm text-muted-foreground">
          <Link href="/account/orders" className="hover:text-foreground">
            My Orders
          </Link>
          <ChevronRight className="h-4 w-4" />
          <span className="text-foreground">Order #{order.id}</span>
        </div>

        <div className="mt-6 flex items-start justify-between gap-3">
          <div>
            <h1 className="font-serif text-2xl font-semibold">
              Order Number: #{order.id}
            </h1>
            <p className="mt-1 text-sm text-muted-foreground">
              {new Date(order.created_at).toLocaleString("en-NG", {
                dateStyle: "medium",
                timeStyle: "short",
              })}
            </p>
          </div>
          <span
            className={`shrink-0 rounded-full px-3 py-1 text-xs font-medium capitalize ${statusStyle(order.status)}`}
          >
            {order.status}
          </span>
        </div>

        {/* items */}
        <div className="mt-8 rounded-xl border bg-background">
          {order.order_items.map((item, i) => {
            return (
              <div
                key={item.id}
                className={`flex gap-3 p-4 ${i !== order.order_items.length - 1 ? "border-b" : ""}`}
              >
                <div className="relative h-14 w-14 shrink-0 overflow-hidden rounded-lg bg-muted">
                  {item.menu_item.images?.length ? (
                    <Image
                      src={item.menu_item.images[0].alt_text}
                      alt={item.menu_item.images[0].url || item.menu_item.name}
                      fill
                      className="object-cover"
                    />
                  ) : (
                    <div>{item.menu_item.image_url}</div>
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
            );
          })}
        </div>

        {/* notes */}
        {order.notes && (
          <div className="mt-4 rounded-xl border bg-background p-4">
            <p className="text-xs font-medium text-muted-foreground">Notes</p>
            <p className="mt-1 text-sm">{order.notes}</p>
          </div>
        )}

        {/* payment */}
        <div className="mt-4 rounded-xl border bg-background p-4">
          <p className="text-xs font-medium text-muted-foreground">Payment</p>
          <div className="mt-2 flex items-center justify-between text-sm">
            <span className="capitalize text-muted-foreground">
              {order.payment?.provider ?? "—"} · {order.payment_status}
            </span>
            {order.payment?.reference && (
              <span className="font-mono text-xs text-muted-foreground">
                {order.payment.reference}
              </span>
            )}
          </div>
        </div>

        {/* total */}
        <div className="mt-6 flex items-center justify-between rounded-xl border bg-background p-4">
          <span className="text-sm font-medium">Total</span>
          <span className="text-lg font-semibold">
            {formatNaira(order.total_amount)}
          </span>
        </div>
      </div>
    </div>
  );
}

function OrderDetailsSkeleton() {
  return (
    <div className="mx-auto max-w-2xl px-6 py-10">
      <div className="h-4 w-32 animate-pulse rounded bg-muted" />
      <div className="mt-6 h-8 w-48 animate-pulse rounded bg-muted" />
      <div className="mt-8 space-y-3">
        {Array.from({ length: 2 }).map((_, i) => (
          <div key={i} className="h-20 animate-pulse rounded-xl bg-muted" />
        ))}
      </div>
    </div>
  );
}
