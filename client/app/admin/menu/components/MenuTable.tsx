"use client";

import {
  Table,
  TableBody,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { Menu } from "@/app/state/types/menuTypes";
import MenuEmptyState from "./MenuEmpty";
import MenuSkeleton from "./MenuSkeleton";
import MenuTableRow from "./MenuTableRow";

interface MenuTableProps {
  menus: Menu[];
  isLoading: boolean;
  isFetching: boolean;
  hasMore: boolean;
  onLoadMore: () => void;
}

export default function MenuTable({
  menus,
  isLoading,
  isFetching,
  hasMore,
  onLoadMore,
}: MenuTableProps) {
  if (isLoading) {
    return <MenuSkeleton />;
  }

  if (!menus.length) {
    return <MenuEmptyState />;
  }

  return (
    <div className="overflow-hidden rounded-2xl border bg-[#faedcd]">
      <Table>
        <TableHeader>
          <TableRow className="bg-[#fefae0] hover:bg-[#fefae0]">
            <TableHead className="w-20">Image</TableHead>

            <TableHead>Menu</TableHead>

            <TableHead>Category</TableHead>

            <TableHead className="text-right">Price</TableHead>

            <TableHead>Status</TableHead>

            <TableHead className="w-32 text-right">Actions</TableHead>
          </TableRow>
        </TableHeader>

        <TableBody>
          {menus.map((menu) => (
            <MenuTableRow key={menu.id} menu={menu} />
          ))}
        </TableBody>
      </Table>

      {hasMore && (
        <div className="flex justify-center border-t bg-white p-5">
          <Button variant="outline" disabled={isFetching} onClick={onLoadMore}>
            {isFetching ? "Loading..." : "Load More"}
          </Button>
        </div>
      )}
    </div>
  );
}
