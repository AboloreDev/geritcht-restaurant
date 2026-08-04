"use client";

import * as React from "react";
import { motion } from "framer-motion";
import { HugeiconsIcon } from "@hugeicons/react";
import {
  Add01Icon,
  Delete02Icon,
  Edit02Icon,
  Search01Icon,
  Tag01Icon,
} from "@hugeicons/core-free-icons";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetFooter,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";

type Label = { id: number; name: string };

type LabelManagementDashboardProps = {
  title: string;
  singular: string;
  description: string;
  items: Label[];
  isLoading: boolean;
  accentClass: string;
  onCreate: (name: string) => Promise<void>;
  onUpdate: (id: number, name: string) => Promise<void>;
  onDelete: (id: number) => Promise<void>;
};

export function LabelManagementDashboard({
  title,
  singular,
  description,
  items,
  isLoading,
  accentClass,
  onCreate,
  onUpdate,
  onDelete,
}: LabelManagementDashboardProps) {
  const [query, setQuery] = React.useState("");
  const [sheetOpen, setSheetOpen] = React.useState(false);
  const [selected, setSelected] = React.useState<Label | null>(null);
  const [pendingDelete, setPendingDelete] = React.useState<Label | null>(null);
  const [name, setName] = React.useState("");
  const [isSaving, setIsSaving] = React.useState(false);
  const [isDeleting, setIsDeleting] = React.useState(false);

  const visibleItems = items.filter((item) =>
    item.name.toLowerCase().includes(query.toLowerCase()),
  );

  const openCreate = () => {
    setSelected(null);
    setName("");
    setSheetOpen(true);
  };
  const openEdit = (item: Label) => {
    setSelected(item);
    setName(item.name);
    setSheetOpen(true);
  };
  const closeSheet = () => {
    if (!isSaving) setSheetOpen(false);
  };

  const save = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!name.trim()) return;
    setIsSaving(true);
    try {
      if (selected) await onUpdate(selected.id, name.trim());
      else await onCreate(name.trim());
      setSheetOpen(false);
    } finally {
      setIsSaving(false);
    }
  };

  const confirmDelete = async () => {
    if (!pendingDelete) return;
    setIsDeleting(true);
    try {
      await onDelete(pendingDelete.id);
      setPendingDelete(null);
    } finally {
      setIsDeleting(false);
    }
  };

  return (
    <div className="flex h-screen flex-col overflow-hidden p-4">
      <main className="mx-auto flex w-full max-w-6xl flex-1 flex-col overflow-y-auto px-1 pb-8 sm:px-4">
        <motion.section
          initial={{ opacity: 0, y: 14 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.35 }}
          className="mt-4 rounded-3xl bg-[#faedcd] p-5 sm:p-7"
        >
          <div className="flex flex-col gap-5 sm:flex-row sm:items-start sm:justify-between">
            <div>
              <div
                className={`mb-4 grid size-11 place-items-center rounded-2xl ${accentClass}`}
              >
                <HugeiconsIcon icon={Tag01Icon} strokeWidth={2} />
              </div>
              <h1 className="text-2xl font-semibold">{title}</h1>
              <p className="mt-2 max-w-xl text-sm text-slate-600">
                {description}
              </p>
            </div>
            <Button className="h-10 rounded-2xl" onClick={openCreate}>
              <HugeiconsIcon icon={Add01Icon} strokeWidth={2} />
              Add {singular}
            </Button>
          </div>
          <div className="mt-7 grid gap-3 sm:grid-cols-2">
            <StatCard
              label={`Total ${title.toLowerCase()}`}
              value={items.length}
            />
            <StatCard label="Matching search" value={visibleItems.length} />
          </div>
        </motion.section>

        <motion.section
          initial={{ opacity: 0, y: 14 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.35, delay: 0.08 }}
          className="mt-6 rounded-3xl bg-[#faedcd] p-4 shadow-sm sm:p-6"
        >
          <div className="relative max-w-md">
            <HugeiconsIcon
              icon={Search01Icon}
              strokeWidth={2}
              className="absolute top-1/2 left-3 size-4 -translate-y-1/2 text-muted-foreground"
            />
            <Input
              value={query}
              onChange={(event) => setQuery(event.target.value)}
              className="h-10 rounded-xl pl-10"
              placeholder={`Search ${title.toLowerCase()}`}
            />
          </div>
          <div className="mt-5 overflow-hidden rounded-2xl ">
            <div className="hidden grid-cols-[1fr_80px_120px] items-center gap-4 border-b border-border bg-muted/40 px-5 py-3 text-xs font-medium uppercase tracking-wide text-muted-foreground md:grid">
              <span>Name</span>
              <span>ID</span>
              <span className="text-right">Actions</span>
            </div>

            <div className="divide-y divide-border">
              {isLoading
                ? Array.from({ length: 6 }).map((_, index) => (
                    <div
                      key={index}
                      className="mx-4 my-3 h-16 animate-pulse rounded-xl bg-muted"
                    />
                  ))
                : visibleItems.map((item, index) => (
                    <motion.article
                      key={item.id}
                      initial={{ opacity: 0, y: 8 }}
                      animate={{ opacity: 1, y: 0 }}
                      transition={{
                        duration: 0.25,
                        delay: Math.min(index * 0.03, 0.2),
                      }}
                      className="grid gap-4 px-5 py-4 md:grid-cols-[1fr_80px_120px] md:items-center"
                    >
                      {/* Name */}
                      <div className="flex items-center gap-3">
                        <div
                          className={`grid h-10 w-10 place-items-center rounded-xl ${accentClass}`}
                        >
                          <HugeiconsIcon icon={Tag01Icon} strokeWidth={2} />
                        </div>

                        <div className="min-w-0">
                          <p className="truncate font-medium">{item.name}</p>
                        </div>
                      </div>

                      {/* ID */}
                      <span className="text-sm text-muted-foreground">
                        #{item.id}
                      </span>

                      {/* Actions */}
                      <div className="flex justify-end gap-2">
                        <Button
                          variant="outline"
                          size="sm"
                          className="rounded-2xl border-none "
                          onClick={() => openEdit(item)}
                        >
                          <HugeiconsIcon icon={Edit02Icon} strokeWidth={2} />
                          Edit
                        </Button>

                        <Button
                          variant="destructive"
                          size="sm"
                          className="rounded-2xl border-none bg-red-500 text-white hover:bg-red/70"
                          onClick={() => setPendingDelete(item)}
                        >
                          <HugeiconsIcon icon={Delete02Icon} strokeWidth={2} />
                          Delete
                        </Button>
                      </div>
                    </motion.article>
                  ))}

              {!isLoading && visibleItems.length === 0 && (
                <div className="px-5 py-12 text-center text-sm text-muted-foreground">
                  No items found.
                </div>
              )}
            </div>
          </div>
          {!isLoading && visibleItems.length === 0 && (
            <div className="py-14 text-center text-sm text-muted-foreground">
              No {title.toLowerCase()} found.
            </div>
          )}
        </motion.section>
      </main>

      <Sheet open={sheetOpen} onOpenChange={(open) => open || closeSheet()}>
        <SheetContent
          side="right"
          className="w-full bg-[#fefae0] sm:w-1/2 sm:max-w-none"
        >
          <motion.div
            initial={{ opacity: 0, x: 18 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ duration: 0.32 }}
            className="flex h-full flex-col"
          >
            <SheetHeader>
              <SheetTitle>
                {selected ? `Edit ${singular}` : `Add ${singular}`}
              </SheetTitle>
              <SheetDescription className="text-slate-400">
                Use a clear name that your kitchen team can recognise.
              </SheetDescription>
            </SheetHeader>
            <form onSubmit={save} className="flex flex-1 flex-col">
              <div className="px-6 pt-3">
                <label htmlFor="label-name" className="text-sm font-medium">
                  Name
                </label>
                <Input
                  id="label-name"
                  autoFocus
                  value={name}
                  onChange={(event) => setName(event.target.value)}
                  placeholder={`e.g. ${singular === "allergen" ? "Gluten" : "Vegetarian"}`}
                  className="mt-2 h-10 rounded-xl"
                />
              </div>
              <SheetFooter className="">
                <Button
                  type="button"
                  variant="outline"
                  onClick={closeSheet}
                  disabled={isSaving}
                  className="bg-white text-black hover:bg-white/70 border-none"
                >
                  Cancel
                </Button>
                <Button type="submit" disabled={isSaving || !name.trim()}>
                  {isSaving
                    ? "Saving…"
                    : selected
                      ? "Save changes"
                      : `Add ${singular}`}
                </Button>
              </SheetFooter>
            </form>
          </motion.div>
        </SheetContent>
      </Sheet>

      <AlertDialog
        open={!!pendingDelete}
        onOpenChange={(open) => !open && !isDeleting && setPendingDelete(null)}
      >
        <AlertDialogContent className="bg-[#fefae0]">
          <AlertDialogHeader>
            <AlertDialogTitle>Delete {pendingDelete?.name}?</AlertDialogTitle>
            <AlertDialogDescription>
              This will permanently remove this {singular}. This action cannot
              be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel
              disabled={isDeleting}
              className="bg-white text-black border-none"
            >
              Cancel
            </AlertDialogCancel>
            <AlertDialogAction
              variant="destructive"
              onClick={confirmDelete}
              disabled={isDeleting}
              className="bg-red-500 text-white border-none"
            >
              {isDeleting ? "Deleting…" : `Delete ${singular}`}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}

function StatCard({ label, value }: { label: string; value: number }) {
  return (
    <div className="rounded-2xl bg-white/70 p-4">
      <p className="text-xs text-slate-600">{label}</p>
      <p className="mt-1 text-2xl font-semibold">{value}</p>
    </div>
  );
}
