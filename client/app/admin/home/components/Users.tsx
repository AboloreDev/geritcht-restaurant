"use client";

import Link from "next/link";
import { useGetAllUsersQuery } from "@/app/state/api/userApi";
import { DashboardListCard } from "./DashboardListCard";

export function RecentUsersList() {
  const { data, isLoading, isFetching } = useGetAllUsersQuery();

  const users = [...(data?.data ?? [])]
    .sort(
      (a, b) =>
        new Date(b.created_at).getTime() - new Date(a.created_at).getTime(),
    )
    .slice(0, 5);

  return (
    <DashboardListCard
      title="Recently Joined"
      viewAllHref="/admin/users"
      isLoading={isLoading}
      isEmpty={users.length === 0}
      emptyMessage="No users yet."
      isFetching={isFetching}
    >
      {users.map((user) => (
        <Link
          key={user.id}
          href={`/admin/users/${user.id}`}
          className="flex items-center justify-between rounded-lg p-3 text-sm hover:bg-muted"
        >
          <div>
            <p className="font-medium">
              {user.first_name} {user.last_name}
            </p>
            <p className="text-xs text-muted-foreground">{user.email}</p>
          </div>
          <span className="shrink-0 text-[10px] text-muted-foreground">
            {new Date(user.created_at).toLocaleDateString("en-NG", {
              month: "short",
              day: "numeric",
            })}
          </span>
        </Link>
      ))}
    </DashboardListCard>
  );
}
