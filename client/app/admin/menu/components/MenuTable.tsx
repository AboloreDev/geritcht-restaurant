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
import MenuTableRow from "./MenuTableRow";
import MenuSkeleton from "./MenuSkeleton";
import MenuEmptyState from "./MenuEmpty";

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
  if (isLoading) return <MenuSkeleton />;
  if (!menus.length) return <MenuEmptyState />;

  return (
    <div className="overflow-hidden rounded-2xl bg-[#faedcd]">
      <Table>
        <TableHeader>
          <TableRow className="h-14 bg-[#fefae0] hover:bg-[#fefae0]">
            <TableHead className="w-[90px]">Image</TableHead>
            <TableHead className="min-w-[300px]">Menu</TableHead>
            <TableHead className="w-[160px]">Category</TableHead>
            <TableHead className="w-[110px] text-right">Price</TableHead>
            <TableHead className="w-[130px] text-center">
              Availability
            </TableHead>
            <TableHead className="w-[110px] text-right">Actions</TableHead>
          </TableRow>
        </TableHeader>

        <TableBody>
          {menus.map((menu) => (
            <MenuTableRow key={menu.id} menu={menu} />
          ))}
        </TableBody>
      </Table>

      {hasMore && (
        <div className="border-t bg-white p-5">
          <div className="flex justify-center">
            <Button
              variant="outline"
              disabled={isFetching}
              onClick={onLoadMore}
            >
              {isFetching ? "Loading..." : "Load More"}
            </Button>
          </div>
        </div>
      )}
    </div>
  );
}
