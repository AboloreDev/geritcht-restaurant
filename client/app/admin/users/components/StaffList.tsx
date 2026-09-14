"use client";

import { motion, AnimatePresence } from "framer-motion";
import Link from "next/link";
import { RootState, useAppDispatch, useAppSelector } from "@/app/state/redux";
import { setStaffPage } from "@/app/state/slices/userSlice";
import { useGetAllStaffQuery } from "@/app/state/api/userApi";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";

export default function StaffList() {
  const dispatch = useAppDispatch();
  const page = useAppSelector((s: RootState) => s.user.staff.page);

  const { data, isLoading, isFetching } = useGetAllStaffQuery({
    page,
    limit: 15,
  });
  const staff = data?.data ?? [];
  const hasMore = data ? page < data.meta.total_pages : false;

  return (
    <div>
      {isLoading ? (
        <div className="mt-6 space-y-3">
          {Array.from({ length: 6 }).map((_, i) => (
            <div key={i} className="h-16 animate-pulse rounded-xl bg-muted" />
          ))}
        </div>
      ) : staff.length === 0 ? (
        <p className="mt-16 text-center text-sm text-muted-foreground">
          No staff members yet.
        </p>
      ) : (
        <>
          <div className="mt-6 space-y-3">
            <AnimatePresence initial={false}>
              {staff.map((member, i) => (
                <motion.div
                  key={member.id}
                  initial={{ opacity: 0, y: 8 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ duration: 0.25, delay: i * 0.02 }}
                >
                  <Link
                    href={`#`}
                    className="flex items-center justify-between gap-4 rounded-2xl bg-[#faedcd] p-4 hover:bg-[#fefae0]"
                  >
                    <div>
                      <p className="font-medium">
                        {member.first_name} {member.last_name}
                      </p>
                      <p className="text-sm text-muted-foreground">
                        {member.email}
                      </p>
                    </div>
                    <Badge
                      className={
                        member.is_active ? "bg-green-600" : "bg-gray-400"
                      }
                    >
                      {member.is_active ? "Active" : "Inactive"}
                    </Badge>
                  </Link>
                </motion.div>
              ))}
            </AnimatePresence>
          </div>

          {hasMore && (
            <div className="mt-6 flex justify-center">
              <Button
                variant="outline"
                disabled={isFetching}
                onClick={() => dispatch(setStaffPage(page + 1))}
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
