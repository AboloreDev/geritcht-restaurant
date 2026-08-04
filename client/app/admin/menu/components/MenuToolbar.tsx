"use client";

import { HugeiconsIcon } from "@hugeicons/react";
import { Add01Icon, Search01Icon } from "@hugeicons/core-free-icons";
import { Button } from "@/components/ui/button";
import { AdminMenuSearch } from "./AdminMenuSearch";

export function MenuToolbar({ onCreate }: { onCreate: () => void }) {
  return (
    <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
      <div className="relative w-full sm:max-w-sm">
        <AdminMenuSearch />
      </div>
      <Button className="h-10 rounded-2xl" onClick={onCreate}>
        <HugeiconsIcon icon={Add01Icon} strokeWidth={2} />
        Add menu item
      </Button>
    </div>
  );
}
