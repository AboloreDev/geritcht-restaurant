"use client";

import { useParams, usePathname, useRouter } from "next/navigation";
import { motion } from "framer-motion";
import { useState } from "react";
import { Clock3, Plus, Minus, ChevronRight } from "@mynaui/icons-react";
import { useGetSingleMenuQuery } from "@/app/state/api/menuApi";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import Link from "next/link";
import Image from "next/image";
import { useAuth } from "@/app/hooks/isAuthenticated";
import { formatNaira } from "@/app/utils/formatNaira";
import { RelatedItems } from "./components/RelatedItems";
import { useAddToCartMutation } from "@/app/state/api/cartApi";
import { toast } from "sonner";
import { getApiError } from "@/app/utils/apiError";

export default function SingleMenu() {
  const { id } = useParams();
  const menuID = id as string;
  const router = useRouter();
  const { isAuthenticated } = useAuth();
  const pathname = usePathname();
  const [addToCart, { isLoading: isAdding }] = useAddToCartMutation();

  const { data, isLoading, isError } = useGetSingleMenuQuery({ id: menuID });
  const menu = data?.data;

  const [activeImage, setActiveImage] = useState(0);
  const [quantity, setQuantity] = useState(1);
  const [specialInstructions, setSpecialInstructions] = useState("");

  if (isLoading) return <SingleMenuSkeleton />;

  if (isError || !menu) {
    return (
      <div className="flex min-h-screen flex-col items-center justify-center gap-3 text-center">
        <p className="text-lg font-medium text-primary-deep">
          This dish couldn&apos;t be found.
        </p>
        <Link href="/menu" className="text-sm text-primary hover:underline">
          Back to menu
        </Link>
      </div>
    );
  }

  // gallery: primary image first, then any secondary images
  const gallery = menu.images?.length
    ? [...menu.images]
        .sort((a, b) => (a.is_primary ? -1 : 1))
        .map((img) => img.alt_text)
    : menu.image_url
      ? [menu.image_url]
      : [];

  const allergens = menu.allergens ?? [];
  const dietaryTags = menu.dietary_tags ?? [];

  async function handleAddToCart() {
    if (!isAuthenticated) {
      router.push(`/login?redirect=${encodeURIComponent(pathname)}`);
      return;
    }

    try {
      const response = await addToCart({
        // @ts-expect-error "<>"
        menu_item_id: menu.id,
        quantity,
        special_instructions: specialInstructions || undefined,
      }).unwrap();
      toast.success(response.message);
      setQuantity(1);
      setSpecialInstructions("");
    } catch (err) {
      console.error(err);
      toast.error(getApiError(err));
    }
  }

  return (
    <div className="min-h-screen bg-[url('/assets/bg.png')] bg-cover bg-center bg-fixed">
      <div className="mx-auto max-w-5xl px-6 py-8">
        {/* breadcrumb */}
        <div className="flex items-center gap-2 text-sm text-text-secondary">
          <Link href="/menu" className="hover:text-primary-deep">
            Menu
          </Link>
          <ChevronRight className="h-4 w-4" />
          <span className="text-primary">{menu.name}</span>
        </div>

        <div className="mt-8 grid gap-10 md:grid-cols-2">
          {/* image gallery */}
          <div>
            <motion.div
              key={activeImage}
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              transition={{ duration: 0.25 }}
              className="relative aspect-square w-full overflow-hidden rounded-2xl bg-[#fefae0]"
            >
              {gallery[activeImage] ? (
                <Image
                  src={gallery[activeImage]}
                  alt={menu.name}
                  fill
                  className="object-cover"
                  sizes="(min-width: 768px) 480px, 100vw"
                />
              ) : (
                <div className="flex h-full items-center justify-center text-muted-foreground">
                  No image available
                </div>
              )}
            </motion.div>

            {gallery.length > 1 && (
              <div className="mt-3 flex gap-2">
                {gallery.map((src, i) => (
                  <button
                    key={i}
                    onClick={() => setActiveImage(i)}
                    className={`relative h-16 w-16 shrink-0 overflow-hidden rounded-lg border-2 transition-colors ${
                      activeImage === i
                        ? "border-primary"
                        : "border-transparent"
                    }`}
                  >
                    <Image
                      src={src}
                      alt=""
                      fill
                      className="object-cover"
                      sizes="64px"
                    />
                  </button>
                ))}
              </div>
            )}
          </div>

          {/* details */}
          <div className="flex flex-col">
            <span className="text-sm text-primary">{menu.category?.name}</span>
            <h1 className="mt-1 text-primary-deep font-serif text-3xl font-semibold">
              {menu.name}
            </h1>

            <div className="mt-3 flex items-center gap-4 text-sm text-primary">
              <span className="flex items-center gap-1">
                <Clock3 className="h-4 w-4" />
                Ready in {menu.prep_time_minutes} min
              </span>
              {menu.spice_level > 0 && (
                <span>{"🌶️".repeat(menu.spice_level)}</span>
              )}
            </div>

            <p className="mt-4 text-primary">{menu.description}</p>

            {(dietaryTags.length > 0 || allergens.length > 0) && (
              <div className="mt-4 flex flex-wrap gap-1.5 text-primary-deep">
                {dietaryTags.map((t) => (
                  <Badge
                    key={`d-${t.id}`}
                    variant="secondary"
                    className="text-xs font-normal"
                  >
                    {t.name}
                  </Badge>
                ))}
                {allergens.map((a) => (
                  <Badge
                    key={`a-${a.id}`}
                    variant="destructive"
                    className="text-xs font-normal"
                  >
                    {a.name}
                  </Badge>
                ))}
              </div>
            )}

            <div className="mt-6 text-2xl font-semibold text-primary-deep">
              {formatNaira(menu.price)}
            </div>

            {/* special instructions — this page has room for a real textarea, unlike the card */}
            <div className="mt-4">
              <label className="text-sm font-medium text-primary-deep">
                Special instructions{" "}
                <span className="text-primary/60">(optional)</span>
              </label>
              <textarea
                value={specialInstructions}
                onChange={(e) => setSpecialInstructions(e.target.value)}
                placeholder="No onions, extra spicy, etc."
                rows={2}
                className="mt-1.5 w-full focus:outline-none rounded-lg border border-[#fefae0] bg-[#fefae0] px-3 py-2 text-sm"
              />
            </div>

            {/* quantity + add to cart */}
            <div className="mt-6 flex items-center gap-4 text-primary-deep">
              <div className="flex items-center rounded-full border border-[#fefae0]">
                <button
                  onClick={() => setQuantity((q) => Math.max(1, q - 1))}
                  className="flex h-10 w-10 items-center justify-center"
                  aria-label="Decrease quantity"
                >
                  <Minus className="h-4 w-4" />
                </button>
                <span className="w-8 text-center text-sm">{quantity}</span>
                <button
                  onClick={() => setQuantity((q) => q + 1)}
                  className="flex h-10 w-10 items-center justify-center"
                  aria-label="Increase quantity"
                >
                  <Plus className="h-4 w-4" />
                </button>
              </div>

              <Button
                className="flex-1 text-black"
                disabled={!menu.is_available || isAdding}
                onClick={handleAddToCart}
              >
                {isAdding
                  ? "Adding…"
                  : menu.is_available
                    ? "Add to cart"
                    : "Sold out"}
              </Button>
            </div>
          </div>
        </div>

        {menu.category_id && (
          <RelatedItems categoryId={menu.category_id} currentItemId={menu.id} />
        )}
      </div>
    </div>
  );
}

function SingleMenuSkeleton() {
  return (
    <div className="mx-auto max-w-5xl px-6 py-8">
      <div className="grid gap-10 md:grid-cols-2">
        <div className="aspect-square animate-pulse rounded-2xl bg-primary" />
        <div className="space-y-4">
          <div className="h-4 w-24 animate-pulse rounded bg-primary" />
          <div className="h-8 w-2/3 animate-pulse rounded bg-primary" />
          <div className="h-20 animate-pulse rounded bg-primary" />
          <div className="h-10 w-32 animate-pulse rounded bg-primary" />
        </div>
      </div>
    </div>
  );
}
