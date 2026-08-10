"use client";

import { Header } from "@/components/code/Header";
import { Separator } from "@/components/ui/separator";
import { Button } from "@/components/ui/button";
import { Add01Icon } from "@hugeicons/core-free-icons";
import { HugeiconsIcon } from "@hugeicons/react";

import { useGetAllTablesQuery } from "@/app/state/api/tableApi";
import { RootState, useAppDispatch, useAppSelector } from "@/app/state/redux";
import { openCreateTableSheet } from "@/app/state/slices/tableSlice";
import TableGrid from "./components/TableGrid";
import CreateTableSheet from "./components/CreateTableSheet";
import EditTableSheet from "./components/EditTableSheet";
import DeleteTableDialog from "./components/DeleteTableDialog";

export default function TablesPage() {
  const dispatch = useAppDispatch();
  const createOpen = useAppSelector((s: RootState) => s.table.create.open);
  const editOpen = useAppSelector((s: RootState) => s.table.edit.open);

  const { data, isLoading, isFetching } = useGetAllTablesQuery();
  const tables = data?.data ?? [];

  return (
    <div className="flex min-h-screen flex-col space-y-4 overflow-y-auto p-4">
      <Header
        title="🪑 Tables"
        subTitle="Manage seating, capacity, and table status."
      />

      <div className="flex flex-col space-y-5 rounded-2xl bg-[#faedcd] py-4 px-6">
        <div className="flex items-center justify-between">
          <p className="text-xl text-muted-foreground">
            {isFetching && !isLoading ? "Updating…" : `${tables.length} tables`}
          </p>
          <Button onClick={() => dispatch(openCreateTableSheet())}>
            <HugeiconsIcon icon={Add01Icon} strokeWidth={2} />
            Add Table
          </Button>
        </div>

        <TableGrid tables={tables} isLoading={isLoading} />
      </div>

      {createOpen && <CreateTableSheet />}
      {editOpen && <EditTableSheet />}
      <DeleteTableDialog />
    </div>
  );
}
