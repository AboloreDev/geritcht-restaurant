"use client";

import Image from "next/image";
import { AnimatePresence, motion } from "framer-motion";
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";
import { Button } from "@/components/ui/button";
import {
  Plus,
  Minus,
  Trash,
  ShoppingBag,
  TrashSolid,
} from "@mynaui/icons-react";
import { useAppDispatch, useAppSelector, RootState } from "@/app/state/redux";
import { closeCartDrawer } from "@/app/state/slices/cartSlice";
import {
  useGetCartQuery,
  useUpdateCartItemMutation,
  useDeleteCartItemMutation,
  useClearCartMutation,
} from "@/app/state/api/cartApi";
import { formatNaira } from "@/app/utils/formatNaira";
import { toast } from "sonner";
import { getApiError } from "@/app/utils/apiError";
import Link from "next/link";
import { useState } from "react";
import { useCreateTakoutOrderMutation } from "@/app/state/api/orderApi";
import { useInitializePaymentMutation } from "@/app/state/api/paymentApi";

export function CartDrawer() {
  const dispatch = useAppDispatch();
  const isOpen = useAppSelector((s: RootState) => s.cart.isDrawerOpen);

  const { data, isLoading } = useGetCartQuery(undefined, { skip: !isOpen });
  const [updateItem] = useUpdateCartItemMutation();
  const [deleteItem] = useDeleteCartItemMutation();
  const [clearCart, { isLoading: isClearing }] = useClearCartMutation();

  const cart = data?.data;
  const items = cart?.cart_items ?? [];

  console.log(items);

  // Order
  const [step, setStep] = useState<"cart" | "notes">("cart");
  const [notes, setNotes] = useState("");
  const [paymentUrl, setPaymentUrl] = useState<string | null>(null);

  const [createOrder, { isLoading: isCreatingOrder }] =
    useCreateTakoutOrderMutation();
  const [initializePayment, { isLoading: isInitializing }] =
    useInitializePaymentMutation();

  async function handleUpdateQuantity(itemId: number, quantity: number) {
    if (quantity < 1) return;
    try {
      await updateItem({ itemId, body: { quantity } }).unwrap();
    } catch (err) {
      toast.error(getApiError(err));
    }
  }

  async function handleRemove(itemId: number) {
    try {
      await deleteItem(itemId).unwrap();
    } catch (err) {
      toast.error(getApiError(err));
    }
  }

  async function handleClear() {
    try {
      await clearCart().unwrap();
      toast.success("Cart cleared");
    } catch (err) {
      toast.error(getApiError(err));
    }
  }

  async function handlePlaceOrder() {
    try {
      const orderRes = await createOrder({
        notes: notes || undefined,
      }).unwrap();
      const orderId = orderRes.data.id;

      const paymentRes = await initializePayment({
        order_id: orderId,
      }).unwrap();
      const url = paymentRes.data.authorization_url;

      setPaymentUrl(url);
      window.location.href = url;
    } catch (err) {
      toast.error(getApiError(err));
    }
  }

  return (
    <Sheet
      open={isOpen}
      onOpenChange={(v) => !v && dispatch(closeCartDrawer())}
    >
      <SheetContent
        side="right"
        className="flex w-full flex-col border-l bg-primary p-0 sm:max-w-md"
      >
        {/* Header */}
        <div className="border-b px-6 py-3">
          <div className="flex items-center justify-between">
            <div>
              <h2 className="font-serif text-2xl font-semibold">Your Order</h2>

              <p className="mt-1 text-sm text-black">
                {cart?.item_count ?? 0} item
                {(cart?.item_count ?? 0) !== 1 && "s"}
              </p>
            </div>
          </div>
        </div>

        {isLoading ? (
          <div className="flex-1 space-y-4 p-6">
            {Array.from({ length: 4 }).map((_, i) => (
              <div
                key={i}
                className="h-28 animate-pulse rounded-2xl bg-muted"
              />
            ))}
          </div>
        ) : items.length === 0 ? (
          <div className="flex flex-1 flex-col items-center justify-center px-8 text-center">
            <ShoppingBag className="mb-5 h-12 w-12 text-primary" />

            <h3 className="font-serif text-xl font-semibold">
              Your cart is empty
            </h3>

            <p className="mt-2 text-sm text-muted-foreground">
              Looks like you haven&apos;t added anything yet.
            </p>

            <Link href="/menu" onClick={() => dispatch(closeCartDrawer())}>
              <Button className="mt-6 w-full rounded-full px-8">
                Browse Menu
              </Button>
            </Link>
          </div>
        ) : (
          <>
            {/* Items */}
            <div className="flex-1 space-y-4 overflow-y-auto p-6">
              <AnimatePresence>
                {items.map((item) => (
                  <motion.div
                    key={item.id}
                    layout
                    initial={{ opacity: 0, y: 15 }}
                    animate={{ opacity: 1, y: 0 }}
                    exit={{ opacity: 0, x: 50 }}
                    className="rounded-2xl bg-[#fefae0] p-2 shadow-sm"
                  >
                    <div className="flex gap-4">
                      <div className="relative h-20 w-20 overflow-hidden rounded-xl">
                        <Image
                          src={item.menu_item.images?.alt_text}
                          alt={
                            item.menu_item.images?.url || item.menu_item.name
                          }
                          fill
                          className="object-cover"
                        />
                      </div>
                      <div className="flex flex-1 flex-col">
                        <div className="flex justify-between">
                          <h4 className="font-medium">{item.menu_item.name}</h4>

                          <span className="font-semibold">
                            {formatNaira(item.subtotal)}
                          </span>
                        </div>

                        {item.special_instructions && (
                          <p className="mt-1 text-xs italic text-muted-foreground">
                            {item.special_instructions}
                          </p>
                        )}

                        <div className="mt-4 flex items-center justify-between">
                          <div className="flex items-center rounded-full border-2 bg-muted px-2 py-1">
                            <button
                              onClick={() =>
                                handleUpdateQuantity(item.id, item.quantity - 1)
                              }
                              className="rounded-full cursor-pointer p-1 hover:bg-white"
                            >
                              <Minus size={15} />
                            </button>

                            <span className="mx-3 min-w-[20px] text-center text-sm font-semibold">
                              {item.quantity}
                            </span>

                            <button
                              onClick={() =>
                                handleUpdateQuantity(item.id, item.quantity + 1)
                              }
                              className="rounded-full cursor-pointer p-1 hover:bg-white"
                            >
                              <Plus size={15} />
                            </button>
                          </div>

                          <button
                            onClick={() => handleRemove(item.id)}
                            className="text-sm cursor-pointer text-red-500 transition hover:text-red-600"
                          >
                            <TrashSolid />
                          </button>
                        </div>
                      </div>
                    </div>
                  </motion.div>
                ))}
              </AnimatePresence>
            </div>

            {/* Footer */}
            <div className="space-y-5 p-6 shadow-[0_-8px_30px_rgba(0,0,0,.04)]">
              {step === "cart" ? (
                <div>
                  <div className="flex items-center justify-between">
                    <span className="text-muted-foreground">Subtotal</span>

                    <span className="text-xl font-bold">
                      {formatNaira(cart?.total ?? 0)}
                    </span>
                  </div>

                  <p className="mt-1 text-xs text-muted-foreground">
                    Taxes and delivery will be calculated during checkout.
                  </p>

                  <div className="mt-5  flex gap-3">
                    <Button
                      variant="outline"
                      className="flex-1 bg-[#fefae0] border-none"
                      disabled={isClearing}
                      onClick={handleClear}
                    >
                      {isClearing ? "Clearing..." : "Clear cart"}
                    </Button>

                    <Button
                      className="flex-1 bg-black text-primary cursor-pointer"
                      size="lg"
                      onClick={() => setStep("notes")}
                    >
                      Proceed
                    </Button>
                  </div>
                </div>
              ) : (
                <div>
                  <button
                    onClick={() => setStep("cart")}
                    className="mb-3 text-xs text-muted-foreground hover:text-foreground"
                  >
                    ← Back to cart
                  </button>

                  <label className="text-sm font-medium">
                    Order notes{" "}
                    <span className="text-muted-foreground">(optional)</span>
                  </label>

                  <textarea
                    value={notes}
                    onChange={(e) => setNotes(e.target.value)}
                    placeholder="Any special requests for this order…"
                    rows={3}
                    className="mt-2 w-full rounded-lg border bg-[#fefae0] text-black px-3 py-2 text-sm outline-none"
                  />

                  <Button
                    className="mt-4 w-full bg-black text-primary"
                    size="lg"
                    disabled={isCreatingOrder || isInitializing}
                    onClick={handlePlaceOrder}
                  >
                    {isCreatingOrder || isInitializing
                      ? "Preparing payment..."
                      : "Continue to payment"}
                  </Button>

                  {paymentUrl && (
                    <p className="mt-2 text-center text-xs text-muted-foreground">
                      Redirecting to payment...{" "}
                      <button
                        onClick={() =>
                          navigator.clipboard.writeText(paymentUrl)
                        }
                        className="underline"
                      >
                        copy link
                      </button>
                    </p>
                  )}
                </div>
              )}
            </div>
          </>
        )}
      </SheetContent>
    </Sheet>
  );
}
