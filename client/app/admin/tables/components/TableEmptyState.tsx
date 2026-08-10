import { HugeiconsIcon } from "@hugeicons/react";
import { RestaurantTableIcon } from "@hugeicons/core-free-icons";

export default function TableEmptyState() {
  return (
    <div className="flex flex-col items-center gap-3 py-16 text-center">
      <HugeiconsIcon
        icon={RestaurantTableIcon}
        strokeWidth={1.5}
        size={40}
        className="text-muted-foreground"
      />
      <p className="text-sm text-muted-foreground">
        No tables yet. Add your first one to get started.
      </p>
    </div>
  );
}
