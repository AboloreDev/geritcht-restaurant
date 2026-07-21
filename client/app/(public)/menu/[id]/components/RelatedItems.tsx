"use client";

import Link from "next/link";
import { motion } from "framer-motion";
import { Clock3 } from "@mynaui/icons-react";
import { useGetMenusQuery } from "@/app/state/api/menuApi";
import { resolveImageSrc } from "@/app/utils/resolveImage";
import Image from "next/image";
import { formatNaira } from "@/app/utils/formatNaira";

export function RelatedItems({
  categoryId,
  currentItemId,
}: {
  categoryId: number;
  currentItemId: number;
}) {
  // fetch a couple extra so filtering out the current item still
  // leaves a full row, rather than requesting exactly 5 and risking 4
  const { data, isLoading } = useGetMenusQuery({
    category_id: categoryId,
    limit: 6,
    page: 1,
  });

  const items = (data?.data ?? [])
    .filter((item) => item.id !== currentItemId)
    .slice(0, 5);

  if (isLoading) return <RelatedItemsSkeleton />;
  if (items.length === 0) return null;

  return (
    <section className="mt-16">
      <h2 className="font-serif text-xl text-primary font-medium">
        You might also like
      </h2>

      <div className="mt-5 flex gap-4 overflow-x-auto pb-2">
        {items.map((item, i) => {
          const imageSrc = resolveImageSrc(item);
          return (
            <motion.div
              key={item.id}
              initial={{ opacity: 0, y: 12 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{
                duration: 0.35,
                delay: i * 0.05,
                ease: [0.16, 1, 0.3, 1],
              }}
              className="shrink-0"
            >
              <Link
                href={`/menu/${item.id}`}
                className="block w-40 rounded-2xl border bg-[#fefae0] p-3 transition-shadow hover:shadow-md"
              >
                <div className="relative aspect-square w-full overflow-hidden rounded-lg bg-muted">
                  {imageSrc ? (
                    <Image
                      src={imageSrc}
                      alt={item.name}
                      className="h-full w-full object-cover"
                    />
                  ) : (
                    <div className="flex h-full items-center justify-center text-xs text-muted-foreground">
                      No image
                    </div>
                  )}
                </div>

                <p className="mt-2 truncate text-sm font-medium">{item.name}</p>

                <div className="mt-1 flex items-center justify-between">
                  <span className="text-sm text-muted-foreground">
                    {formatNaira(item.price)}
                  </span>
                  <span className="flex items-center gap-0.5 text-xs text-muted-foreground">
                    <Clock3 className="h-3 w-3" />
                    {item.prep_time_minutes}min
                  </span>
                </div>
              </Link>
            </motion.div>
          );
        })}
      </div>
    </section>
  );
}

function RelatedItemsSkeleton() {
  return (
    <section className="mt-16">
      <div className="h-6 w-40 animate-pulse rounded bg-primary" />
      <div className="mt-5 flex gap-4">
        {Array.from({ length: 4 }).map((_, i) => (
          <div
            key={i}
            className="h-52 w-40 shrink-0 animate-pulse rounded-xl bg-primary"
          />
        ))}
      </div>
    </section>
  );
}
