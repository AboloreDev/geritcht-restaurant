import { TableStatus } from "@/app/state/types/tableTypes";

export function tableStatusStyle(status: TableStatus | string) {
  switch (status) {
    case TableStatus.Available:
      return "bg-green-600 hover:bg-green-600";
    case TableStatus.Occupied:
      return "bg-red-500 hover:bg-red-500";
    case TableStatus.Reserved:
      return "bg-amber-500 hover:bg-amber-500";
    default:
      return "bg-muted text-muted-foreground";
  }
}
