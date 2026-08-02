"use client";

import Link from "next/link";
import { ArrowRight } from "@mynaui/icons-react";

export function DashboardListCard({
  title,
  viewAllHref,
  isLoading,
  isEmpty,
  emptyMessage,
  children,
  isFetching,
}: {
  title: string;
  viewAllHref: string;
  isLoading: boolean;
  isEmpty: boolean;
  emptyMessage: string;
  children: React.ReactNode;
  isFetching: boolean;
}) {
  return (
    <div className="flex h-full flex-col rounded-2xl bg-[#faedcd] p-5">
      <div className="flex items-center justify-between">
        <h3 className="font-serif text-base font-medium">{title}</h3>
        <Link
          href={viewAllHref}
          className="flex items-center gap-1 text-xs text-muted-foreground hover:text-foreground"
        >
          View all <ArrowRight className="h-3 w-3" />
        </Link>
      </div>

      <div className="mt-4 flex-1 space-y-2">
        {isLoading || isFetching ? (
          Array.from({ length: 5 }).map((_, i) => (
            <div
              key={i}
              className="h-14 animate-pulse rounded-lg bg-[#fefae0]"
            />
          ))
        ) : isEmpty ? (
          <p className="py-6 text-center text-sm text-muted-foreground">
            {emptyMessage}
          </p>
        ) : (
          children
        )}
      </div>
    </div>
  );
}
