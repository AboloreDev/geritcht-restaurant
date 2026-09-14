import { PaymentStatus } from "@/app/state/types/paymentTypes";

export const STATUS_OPTIONS = [
  { value: "paid", label: "Paid" },
  { value: "pending", label: "Pending" },
  { value: "unpaid", label: "Unpaid" },
  { value: "failed", label: "Failed" },
  { value: "refunded", label: "Refunded" },
];

export function statusStyle(status: string) {
  switch (status) {
    case "paid":
      return "bg-green-600 hover:bg-green-600 text-white";
    case "pending":
      return "bg-amber-500 hover:bg-amber-500 text-white";
    case "unpaid":
      return "bg-gray-400 hover:bg-gray-400 text-white";
    case "failed":
      return "bg-red-500 hover:bg-red-500 text-white";
    case "refunded":
      return "bg-blue-500 hover:bg-blue-500 text-white";
    default:
      return "bg-muted text-muted-foreground";
  }
}
