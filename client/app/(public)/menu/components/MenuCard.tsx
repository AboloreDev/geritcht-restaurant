"use client";

import Image from "next/image";
import { motion } from "framer-motion";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Clock3, Plus, ImageRectangle } from "@mynaui/icons-react";
import { Menu } from "@/app/state/types/menuTypes";

function formatNaira(amount: number) {
  return new Intl.NumberFormat("en-NG", {
    style: "currency",
    currency: "NGN",
    maximumFractionDigits: 0,
  }).format(amount);
}

function resolveImageSrc(menu: Menu): string | null {
  if (menu.image_url) return menu.image_url;
  if (!menu.images?.length) return null;
  const primary = menu.images.find((img) => img.is_primary);
  return (primary ?? menu.images[0]).alt_text || null;
}

function getVisibleBadges(menu: Menu, max = 2) {
  const tags = [
    ...(menu.dietary_tags ?? []).map((t) => ({
      id: `d-${t.id}`,
      name: t.name,
      kind: "diet" as const,
    })),
    ...(menu.allergens ?? []).map((a) => ({
      id: `a-${a.id}`,
      name: a.name,
      kind: "allergen" as const,
    })),
  ];
  return {
    visible: tags.slice(0, max),
    overflowCount: Math.max(tags.length - max, 0),
  };
}

export function MenuCard({ menu }: { menu: Menu }) {
  const imageSrc = resolveImageSrc(menu);
  const { visible: badges, overflowCount } = getVisibleBadges(menu);

  return (
    <div className="relative mt-10 rounded-2xl bg-[#fefee3] p-4 pt-4 sm:mt-12 sm:p-5">
      {/* Floating Image */}
      <div className="absolute left-1/2 top-0 -translate-x-1/2 -translate-y-1/2">
        <motion.div
          whileHover={{ scale: 1.05 }}
          transition={{ type: "spring", stiffness: 300, damping: 20 }}
          className="relative flex h-20 w-20 items-center justify-center overflow-hidden rounded-full shadow-md ring-4 ring-white sm:h-24 sm:w-24"
        >
          {imageSrc ? (
            <Image
              src={imageSrc}
              alt={menu.name}
              fill
              sizes="(max-width: 640px) 80px, 96px"
              className="object-cover"
            />
          ) : (
            <ImageRectangle
              size={24}
              className="text-muted-foreground sm:size-7"
            />
          )}
        </motion.div>
      </div>

      {/* Content */}
      <div className="flex flex-col pt-12 sm:pt-14">
        {/* Title + Prep Time */}
        <div className="flex items-start justify-between gap-3">
          <h3 className="font-serif text-base font-semibold leading-tight sm:text-lg">
            {menu.name}
          </h3>

          <span className="flex shrink-0 items-center gap-1 text-xs text-muted-foreground">
            <Clock3 size={13} />
            {menu.prep_time_minutes} min
          </span>
        </div>

        {/* Description */}
        <p className="mt-2 line-clamp-2 text-sm leading-6 text-muted-foreground">
          {menu.description}
        </p>

        {/* Tags */}
        {badges.length > 0 && (
          <div className="mt-3 flex flex-wrap gap-1.5">
            {badges.map((badge) => (
              <Badge
                key={badge.id}
                variant={
                  badge.kind === "allergen" ? "destructive" : "secondary"
                }
                className="text-[10px] font-normal"
              >
                {badge.name}
              </Badge>
            ))}

            {overflowCount > 0 && (
              <Badge variant="outline" className="text-[10px] font-normal">
                +{overflowCount}
              </Badge>
            )}
          </div>
        )}

        {/* Footer */}
        <div className="mt-4 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
          <span className="text-lg font-semibold text-primary">
            {formatNaira(menu.price)}
          </span>

          <Button
            size="sm"
            disabled={!menu.is_available}
            className="w-full sm:w-auto"
          >
            <Plus size={14} />
            {menu.is_available ? "Add" : "Sold out"}
          </Button>
        </div>
      </div>
    </div>
  );
}
