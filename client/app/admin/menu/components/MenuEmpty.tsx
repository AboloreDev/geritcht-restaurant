"use client";

import { Button } from "@/components/ui/button";
import { Bowl } from "@mynaui/icons-react";

interface MenuEmptyStateProps {
  onCreate?: () => void;
}

export default function MenuEmptyState({ onCreate }: MenuEmptyStateProps) {
  return (
    <div className="flex min-h-[350px] flex-col items-center justify-center rounded-2xl border bg-[#faedcd] p-8 text-center">
      <div className="mb-5 rounded-full bg-[#fefae0] p-5">
        <Bowl className="h-10 w-10 text-[#bc6c25]" />
      </div>

      <h3 className="text-xl font-semibold">No menu items yet</h3>

      <p className="mt-2 max-w-md text-sm text-muted-foreground">
        Your restaurant doesn&apos;t have any menu items yet. Create your first
        dish to start serving customers.
      </p>

      {onCreate && (
        <Button onClick={onCreate} className="mt-6 rounded-xl">
          Add First Menu Item
        </Button>
      )}
    </div>
  );
}
