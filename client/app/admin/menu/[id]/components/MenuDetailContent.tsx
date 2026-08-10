"use client";

import { useParams, useRouter } from "next/navigation";
import Image from "next/image";
import Link from "next/link";
import { motion } from "framer-motion";
import { HugeiconsIcon } from "@hugeicons/react";
import {
  ArrowLeft01Icon,
  Edit02Icon,
  Delete02Icon,
  Clock01Icon,
  FireIcon,
} from "@hugeicons/core-free-icons";

import { useGetSingleMenuQuery } from "@/app/state/api/menuApi";
import { useAppDispatch } from "@/app/state/redux";
import { formatNaira } from "@/app/utils/formatNaira";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { useState } from "react";
import {
  openDeleteMenuDialog,
  openEditSheet,
} from "@/app/state/slices/menuSlice";

export default function MenuDetailContent() {
  const { id } = useParams();
  const router = useRouter();
  const dispatch = useAppDispatch();
  const [activeImage, setActiveImage] = useState(0);

  const { data, isLoading, isError } = useGetSingleMenuQuery({
    id: String(id),
  });

  if (isLoading) return <MenuDetailSkeleton />;

  if (isError || !data?.data) {
    return (
      <div className="flex min-h-[60vh] flex-col items-center justify-center gap-3 text-center">
        <p className="text-lg font-medium">
          This menu item couldn&apos;t be found.
        </p>
        <Link
          href="/admin/menu"
          className="text-sm text-primary hover:underline"
        >
          Back to menu
        </Link>
      </div>
    );
  }

  const menu = data.data;
  const gallery = menu.images?.length
    ? [...menu.images]
        .sort((a, b) => (a.is_primary ? -1 : 1))
        .map((img) => img.url)
    : menu.image_url
      ? [menu.image_url]
      : [];

  return (
    <div className="p-6">
      {/* breadcrumb / back */}
      <button
        onClick={() => router.push("/admin/menu")}
        className="flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground"
      >
        <HugeiconsIcon icon={ArrowLeft01Icon} size={16} />
        Back to menu
      </button>

      {/* header */}
      <div className="mt-4 flex flex-wrap items-start justify-between gap-4">
        <div>
          <div className="flex items-center gap-3">
            <h1 className="font-serif text-2xl font-semibold">{menu.name}</h1>
            <Badge
              className={
                menu.is_available
                  ? "bg-green-600 hover:bg-green-600"
                  : "bg-red-500 hover:bg-red-500"
              }
            >
              {menu.is_available ? "Available" : "Unavailable"}
            </Badge>
          </div>
          <p className="mt-1 text-sm text-muted-foreground">
            #{menu.id} · Added{" "}
            {new Date(menu.created_at).toLocaleDateString("en-NG", {
              dateStyle: "medium",
            })}
          </p>
        </div>

        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            onClick={() => dispatch(openEditSheet(menu.id))}
          >
            <HugeiconsIcon icon={Edit02Icon} strokeWidth={2} size={16} />
            Edit
          </Button>
          <Button
            variant="outline"
            className="bg-red-500 hover:bg-red-600 border-none text-white"
            onClick={() => dispatch(openDeleteMenuDialog(menu.id))}
          >
            <HugeiconsIcon icon={Delete02Icon} strokeWidth={2} size={16} />
            Delete
          </Button>
        </div>
      </div>

      <div className="mt-6 grid gap-6 lg:grid-cols-3">
        {/* images — main column */}
        <div className="lg:col-span-2 space-y-4">
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            className="relative aspect-square w-full overflow-hidden rounded-2xl bg-muted"
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
              <div className="flex h-full items-center justify-center text-sm text-muted-foreground">
                No image uploaded
              </div>
            )}
          </motion.div>

          {gallery.length > 0 && (
            <div className="flex gap-3">
              {gallery.map((img, i) => (
                <div
                  key={i}
                  onClick={() => setActiveImage(i)}
                  className="relative h-20 w-20 overflow-hidden rounded-lg bg-muted"
                >
                  <Image
                    src={img}
                    alt={menu.name}
                    fill
                    className="object-cover"
                    sizes="80px"
                  />
                </div>
              ))}
            </div>
          )}

          {/* description */}
          <div className="rounded-2xl  bg-[#faedcd] p-4">
            <p className="text-xs font-medium text-muted-foreground">
              Description
            </p>
            <p className="mt-1 text-sm">
              {menu.description || "No description provided."}
            </p>
          </div>
        </div>

        {/* details — side column */}
        <div className="space-y-4">
          <div className="rounded-2xl  bg-[#faedcd] p-4">
            <p className="text-xs font-medium text-muted-foreground">Price</p>
            <p className="mt-1 text-xl font-semibold">
              {formatNaira(menu.price)}
            </p>
          </div>

          <div className="rounded-2xl  bg-[#faedcd] p-4">
            <p className="text-xs font-medium text-muted-foreground">
              Category
            </p>
            <Badge variant="outline" className="mt-1.5 rounded-full">
              {menu.category.name}
            </Badge>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className="rounded-2xl  bg-[#faedcd] p-4">
              <div className="flex items-center gap-1.5 text-xs font-medium text-muted-foreground">
                <HugeiconsIcon icon={Clock01Icon} size={14} />
                Prep time
              </div>
              <p className="mt-1 text-sm font-medium">
                {menu.prep_time_minutes} min
              </p>
            </div>

            <div className="rounded-2xl  bg-[#faedcd] p-4">
              <div className="flex items-center gap-1.5 text-xs font-medium text-muted-foreground">
                <HugeiconsIcon icon={FireIcon} size={14} />
                Spice level
              </div>
              <p className="mt-1 text-sm font-medium">
                {menu.spice_level > 0 ? "🌶️".repeat(menu.spice_level) : "None"}
              </p>
            </div>
          </div>

          {menu.allergens.length > 0 && (
            <div className="rounded-2xl  bg-[#faedcd] p-4">
              <p className="text-xs font-medium text-muted-foreground">
                Allergens
              </p>
              <div className="mt-2 flex flex-wrap gap-1.5">
                {menu.allergens.map((a) => (
                  <Badge
                    key={a.id}
                    variant="destructive"
                    className="text-xs font-normal"
                  >
                    {a.name}
                  </Badge>
                ))}
              </div>
            </div>
          )}

          {menu.dietary_tags.length > 0 && (
            <div className="rounded-2xl  bg-[#faedcd] p-4">
              <p className="text-xs font-medium text-muted-foreground">
                Dietary Tags
              </p>
              <div className="mt-2 flex flex-wrap gap-1.5">
                {menu.dietary_tags.map((t) => (
                  <Badge
                    key={t.id}
                    variant="secondary"
                    className="text-xs font-normal"
                  >
                    {t.name}
                  </Badge>
                ))}
              </div>
            </div>
          )}

          <div className="rounded-2xl  bg-[#faedcd] p-4 text-xs text-muted-foreground">
            <p>Display order: {menu.display_order}</p>
            <p className="mt-1">
              Last updated{" "}
              {new Date(menu.updated_at).toLocaleString("en-NG", {
                dateStyle: "medium",
                timeStyle: "short",
              })}
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}

function MenuDetailSkeleton() {
  return (
    <div className="p-6">
      <div className="h-4 w-24 animate-pulse rounded bg-[#fefae0]" />
      <div className="mt-4 h-8 w-64 animate-pulse rounded bg-[#fefae0]" />
      <div className="mt-6 grid gap-6 lg:grid-cols-3">
        <div className="lg:col-span-2 aspect-square animate-pulse rounded-2xl bg-[#fefae0]" />
        <div className="space-y-4">
          {Array.from({ length: 4 }).map((_, i) => (
            <div
              key={i}
              className="h-20 animate-pulse rounded-xl bg-[#fefae0]"
            />
          ))}
        </div>
      </div>
    </div>
  );
}
