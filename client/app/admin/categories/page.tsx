"use client";

import * as React from "react";
import { HugeiconsIcon } from "@hugeicons/react";
import {
  Add01Icon,
  Delete02Icon,
  Edit02Icon,
  Folder01Icon,
  MoreHorizontalIcon,
  Search01Icon,
} from "@hugeicons/core-free-icons";
import { Header } from "@/components/code/Header";
import { Button } from "@/components/ui/button";
import { useGetCategoriesQuery } from "@/app/state/api/categoriesApi";
import { RootState, useAppSelector } from "@/app/state/redux";
import { MenuCategory } from "@/app/state/types/menuTypes";
import CategoriesList from "./components/CategoriesList";
import EditCategorySheet from "./components/EditCategorySheet";
import DeleteCategoryDialog from "./components/DeleteCategoryDialog";
import CreateCategorySheet from "./components/CreateCategorySheet";
import CategoriesSearch from "./components/CategoriesSearch";
import { SummaryCard } from "@/components/code/SummaryCard";

export default function CategoriesPage() {
  const { query, page, limit } = useAppSelector(
    (state: RootState) => state.category,
  );
  const { data } = useGetCategoriesQuery({
    page: page,
    limit: limit,
    query: query,
  });

  const categories = data?.data ?? [];

  const [editingCategory, setEditingCategory] =
    React.useState<MenuCategory | null>(null);
  const [deletingCategory, setDeletingCategory] =
    React.useState<MenuCategory | null>(null);
  const [createOpen, setCreateOpen] = React.useState(false);

  const openCreateDialog = () => {
    setCreateOpen(true);
  };

  return (
    <div className="flex h-screen flex-col overflow-hidden p-4">
      <Header
        title="Categories"
        subTitle="Organise the groups your guests browse on the menu."
      />

      <main className="flex-1 overflow-y-auto px-4 pb-8 sm:px-6 lg:px-8">
        <div className="mx-auto max-w-6xl space-y-6 pt-3">
          <section className="grid gap-4 sm:grid-cols-3">
            <SummaryCard
              label="Total categories"
              value={categories.length}
              detail="Across your menu"
            />
            <SummaryCard
              label="Visible to guests"
              value={categories.filter((category) => category.is_active).length}
              detail="Currently active"
            />
            <SummaryCard
              label="Menu items"
              value={data?.meta.total ?? 0}
              detail="Assigned to categories"
            />
          </section>

          <section className="rounded-3xl bg-[#faedcd] p-4 shadow-sm sm:p-6">
            <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
              <div>
                <h1 className="text-lg font-semibold text-foreground">
                  Menu categories
                </h1>
                <p className="mt-1 text-sm text-muted-foreground">
                  Create, reorder, and control which categories appear on your
                  menu.
                </p>
              </div>
              <Button className="h-10 rounded-2xl" onClick={openCreateDialog}>
                <HugeiconsIcon icon={Add01Icon} strokeWidth={2} />
                Add category
              </Button>
            </div>

            <div className="mt-4 flex items-center gap-3">
              <p>Search Categories:</p>
              <CategoriesSearch />
            </div>

            <div className="mt-5 overflow-hidden rounded-2xl ">
              {/* Header */}
              <div className="hidden grid-cols-[minmax(220px,1fr)_120px_80px] items-center gap-4 border-b border-border bg-muted/40 px-5 py-3 text-xs font-medium uppercase tracking-wide text-muted-foreground md:grid">
                <span>Category</span>
                <span>Status</span>
                <span className="text-right">Actions</span>
              </div>

              <div className="divide-y divide-border">
                {categories.map((category: MenuCategory) => (
                  <CategoriesList
                    key={category.id}
                    category={category}
                    onEdit={setEditingCategory}
                    onDelete={setDeletingCategory}
                  />
                ))}

                {categories.length === 0 && (
                  <div className="px-5 py-12 text-center text-sm text-muted-foreground">
                    No categories available.
                  </div>
                )}
              </div>
            </div>
          </section>
        </div>
      </main>

      {editingCategory && (
        <EditCategorySheet
          key={editingCategory.id}
          category={editingCategory}
          onOpenChange={(open) => !open && setEditingCategory(null)}
        />
      )}
      {deletingCategory && (
        <DeleteCategoryDialog
          category={deletingCategory}
          onOpenChange={(open) => !open && setDeletingCategory(null)}
        />
      )}
      <CreateCategorySheet open={createOpen} onOpenChange={setCreateOpen} />
    </div>
  );
}
