export const TABLE_STATUS_CONFIG = {
  available: {
    label: "Available",
    dot: "bg-emerald-500",
    card: "border-border bg-[#fefae0] hover:border-primary hover:shadow-md cursor-pointer",
    clickable: true,
  },
  confirmed: {
    label: "Booked",
    dot: "bg-amber-500",
    card: "border-border bg-[#fefae0] opacity-60 cursor-not-allowed",
    clickable: false,
  },
} as const;

export type TableStatus = keyof typeof TABLE_STATUS_CONFIG;

export function getStatusConfig(status: string) {
  return (
    TABLE_STATUS_CONFIG[status as TableStatus] ?? TABLE_STATUS_CONFIG.confirmed
  );
}
