export const STATUS_STYLES: Record<string, string> = {
  confirmed: "bg-emerald-500/10 text-emerald-600",
  pending: "bg-amber-500/10 text-amber-600",
  cancelled: "bg-red-500/10 text-red-600",
  completed: "bg-blue-500/10 text-blue-600",
  checked_in: "bg-blue-500/10 text-blue-600",
  no_show: "bg-gray-500/10 text-gray-600",
};

export function statusStyle(status: string) {
  return STATUS_STYLES[status] ?? "bg-muted text-muted-foreground";
}

export const STATUS_OPTIONS = [
  { value: "pending", label: "Pending" },
  { value: "confirmed", label: "Confirmed" },
  { value: "completed", label: "Completed" },
  { value: "cancelled", label: "Cancelled" },
  { value: "checked_in", label: "Checked In" },
  { value: "no_show", label: "No Show" },
];
