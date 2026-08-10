"use client";

import { motion, AnimatePresence } from "framer-motion";
import { Table } from "@/app/state/types/tableTypes";

import TableEmptyState from "./TableEmptyState";
import TableListSkeleton from "./TableListSkeleton";
import TableRow from "./TableCard";

export default function TableList({
  tables,
  isLoading,
}: {
  tables: Table[];
  isLoading: boolean;
}) {
  if (isLoading) return <TableListSkeleton />;
  if (!tables.length) return <TableEmptyState />;

  return (
    <motion.div layout className="space-y-3">
      <AnimatePresence>
        {tables.map((table) => (
          <TableRow key={table.id} table={table} />
        ))}
      </AnimatePresence>
    </motion.div>
  );
}
