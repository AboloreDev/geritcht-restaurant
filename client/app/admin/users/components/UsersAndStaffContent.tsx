"use client";

import { RootState, useAppDispatch, useAppSelector } from "@/app/state/redux";
import { setActiveTab } from "@/app/state/slices/userSlice";
import { cn } from "@/lib/utils";
import { motion, AnimatePresence } from "framer-motion";
import UsersList from "./UsersList";
import StaffList from "./StaffList";
import PromoteToStaffDialog from "./PromoteUserDialog";

export default function UsersAndStaffContent() {
  const dispatch = useAppDispatch();
  const activeTab = useAppSelector((s: RootState) => s.user.activeTab);

  return (
    <div className="p-6">
      <h1 className="font-serif text-2xl font-semibold">People</h1>
      <p className="mt-1 text-sm text-muted-foreground">
        Manage registered users and staff members.
      </p>
      {/* tab switcher — pill style, matches the rest of the admin UI's
          rounded/segmented control language */}
      <div className="mt-6 inline-flex items-center gap-1 rounded-full  p-1">
        {(["users", "staff"] as const).map((tab) => (
          <button
            key={tab}
            onClick={() => dispatch(setActiveTab(tab))}
            className={cn(
              "relative rounded-full px-5 py-1.5 text-sm font-medium capitalize transition-colors",
              activeTab === tab
                ? "text-white"
                : "text-muted-foreground hover:text-foreground",
            )}
          >
            {activeTab === tab && (
              <motion.span
                layoutId="active-people-tab"
                className="absolute inset-0 rounded-full bg-primary"
                transition={{ type: "spring", stiffness: 400, damping: 32 }}
              />
            )}
            <span className="relative z-10">{tab}</span>
          </button>
        ))}
      </div>
      <AnimatePresence mode="wait">
        {activeTab === "users" ? (
          <motion.div
            key="users"
            initial={{ opacity: 0, y: 8 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -8 }}
            transition={{ duration: 0.2 }}
          >
            <UsersList />
          </motion.div>
        ) : (
          <motion.div
            key="staff"
            initial={{ opacity: 0, y: 8 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -8 }}
            transition={{ duration: 0.2 }}
          >
            <StaffList />
          </motion.div>
        )}
      </AnimatePresence>
      <PromoteToStaffDialog />
    </div>
  );
}
