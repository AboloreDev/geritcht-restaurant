"use client";

import { CategoryTabs } from "../categories/CategoryTabs";
import MenuFilters from "./components/MenuFilters";
import { MenuGrid } from "./components/MenuGrid";
import MenuHeader from "./components/MenuHeader";

export default function MenuPage() {
  return (
    <main className="min-h-screen bg-[url('/assets/bg.png')] bg-cover bg-center bg-fixed">
      <MenuHeader />

      <MenuFilters />
      <section className="mx-auto max-w-7xl px-6 py-5 lg:px-10">
        <CategoryTabs />

        <div className="mt-20 space-y-24">
          {/* Menu items will be rendered here */}
          <MenuGrid />
        </div>
      </section>
    </main>
  );
}
