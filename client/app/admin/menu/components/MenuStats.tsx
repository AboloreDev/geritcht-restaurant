"use client";

import { Menu } from "@/app/state/types/menuTypes";
import { SummaryCard } from "@/components/code/SummaryCard";

interface MenuStatsProps {
  menus: Menu[];
  total: number;
}

export default function MenuStats({ menus, total }: MenuStatsProps) {
  const availableMenus = menus.filter((menu) => menu.is_available).length;
  const totalCategories = new Set(menus.map((menu) => menu.category?.id)).size;
  const averagePrice =
    total > 0
      ? Math.round(menus.reduce((sum, menu) => sum + menu.price, 0) / total)
      : 0;

  return (
    <section className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
      <SummaryCard label="Menu Items" value={total} detail="Total menu items" />

      <SummaryCard
        label="Available"
        value={availableMenus}
        detail="Currently available"
      />

      <SummaryCard
        label="Categories"
        value={totalCategories}
        detail="Assigned categories"
      />

      <SummaryCard
        label="Average Price"
        value={averagePrice}
        detail="Average menu price"
      />
    </section>
  );
}
