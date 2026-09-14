"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { motion, AnimatePresence } from "framer-motion";
import { HugeiconsIcon } from "@hugeicons/react";
import { Search01Icon, UserAdd01Icon } from "@hugeicons/core-free-icons";
import { cn } from "@/lib/utils";
import { useDebounce } from "@/app/hooks/useDebounce";
import { RootState, useAppDispatch, useAppSelector } from "@/app/state/redux";
import {
  setUsersPage,
  setUsersSearch,
  openPromoteDialog,
} from "@/app/state/slices/userSlice";
import { useGetAllUsersQuery } from "@/app/state/api/userApi";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";

export default function UsersList() {
  const dispatch = useAppDispatch();
  const { page, query } = useAppSelector((s: RootState) => s.user.users);

  const [searchInput, setSearchInput] = useState(query);
  const debouncedSearch = useDebounce(searchInput, 400);

  useEffect(() => {
    dispatch(setUsersSearch(debouncedSearch));
  }, [debouncedSearch, dispatch]);

  const { data, isLoading, isFetching } = useGetAllUsersQuery({
    page,
    limit: 15,
    q: query || undefined,
  });

  const users = data?.data ?? [];
  const hasMore = data ? page < data.meta.total_pages : false;
  const isRefetching = isFetching && !isLoading && page === 1;

  return (
    <div>
      <div className="mt-6 flex items-center gap-3 rounded-2xl  px-6 py-3">
        <div className="relative flex-1">
          <HugeiconsIcon
            icon={Search01Icon}
            size={16}
            className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground"
          />
          <input
            type="text"
            value={searchInput}
            onChange={(e) => setSearchInput(e.target.value)}
            placeholder="Search by name or email…"
            className="w-full rounded-full border border-black py-1.5 pl-9 pr-3 text-sm"
          />
        </div>
        {isRefetching && (
          <span className="text-xs text-muted-foreground">Updating…</span>
        )}
      </div>

      {isLoading ? (
        <div className="mt-6 space-y-3">
          {Array.from({ length: 8 }).map((_, i) => (
            <div key={i} className="h-16 animate-pulse rounded-xl bg-muted" />
          ))}
        </div>
      ) : users.length === 0 ? (
        <p className="mt-16 text-center text-sm text-muted-foreground">
          No users found.
        </p>
      ) : (
        <>
          <div
            className={cn(
              "mt-6 space-y-3 transition-opacity",
              isRefetching && "opacity-50",
            )}
          >
            <AnimatePresence initial={false}>
              {users.map((user, i) => (
                <motion.div
                  key={user.id}
                  initial={{ opacity: 0, y: 8 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ duration: 0.25, delay: i * 0.02 }}
                  className="flex items-center justify-between gap-4 rounded-2xl bg-[#faedcd] p-4"
                >
                  <Link href={`#`} className="min-w-0 flex-1">
                    <p className="font-medium text-sm">
                      {user.first_name} {user.last_name}
                    </p>
                    <p className="text-xs text-slate-500">{user.email}</p>
                  </Link>

                  <div className="flex shrink-0 items-center gap-2">
                    <Badge
                      variant={user.role === "admin" ? "default" : "secondary"}
                      className="capitalize"
                    >
                      {user.role}
                    </Badge>
                    <Badge
                      className={
                        user.is_active
                          ? "bg-green-600 text-white"
                          : "bg-gray-400"
                      }
                    >
                      {user.is_active ? "Active" : "Inactive"}
                    </Badge>

                    {user.role === "user" && (
                      <Button
                        size="sm"
                        variant="outline"
                        onClick={() => dispatch(openPromoteDialog(user.id))}
                      >
                        <HugeiconsIcon icon={UserAdd01Icon} size={14} />
                        Promote
                      </Button>
                    )}
                  </div>
                </motion.div>
              ))}
            </AnimatePresence>
          </div>

          {hasMore && (
            <div className="mt-6 flex justify-center">
              <Button
                variant="outline"
                disabled={isFetching}
                onClick={() => dispatch(setUsersPage(page + 1))}
              >
                {isFetching && page > 1 ? "Loading…" : "Load more"}
              </Button>
            </div>
          )}
        </>
      )}
    </div>
  );
}
