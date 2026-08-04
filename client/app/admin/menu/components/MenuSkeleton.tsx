"use client";

import { Table, TableBody, TableCell, TableRow } from "@/components/ui/table";

export default function MenuSkeleton() {
  return (
    <div className="overflow-hidden rounded-2xl border bg-[#faedcd]">
      <Table>
        <TableBody>
          {Array.from({ length: 8 }).map((_, index) => (
            <TableRow key={index}>
              <TableCell>
                <div className="h-14 w-14 animate-pulse rounded-xl bg-[#d4a373]" />
              </TableCell>

              <TableCell>
                <div className="space-y-2">
                  <div className="h-4 w-40 animate-pulse rounded bg-[#d4a373]" />
                  <div className="h-3 w-60 animate-pulse rounded bg-[#d4a373]/70" />
                </div>
              </TableCell>

              <TableCell>
                <div className="h-4 w-24 animate-pulse rounded bg-[#d4a373]" />
              </TableCell>

              <TableCell className="text-right">
                <div className="ml-auto h-4 w-20 animate-pulse rounded bg-[#d4a373]" />
              </TableCell>

              <TableCell>
                <div className="h-7 w-24 animate-pulse rounded-full bg-[#d4a373]" />
              </TableCell>

              <TableCell>
                <div className="flex justify-end gap-2">
                  <div className="h-9 w-16 animate-pulse rounded-lg bg-[#d4a373]" />
                  <div className="h-9 w-16 animate-pulse rounded-lg bg-[#d4a373]" />
                </div>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
}
