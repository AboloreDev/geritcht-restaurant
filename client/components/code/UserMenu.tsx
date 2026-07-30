// src/app/components/UserMenu.tsx
"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { motion, AnimatePresence } from "framer-motion";
import {
  User,
  Ticket,
  CalendarCheck,
  CogThree,
  Logout,
} from "@mynaui/icons-react";
import { useLogoutMutation } from "@/app/state/api/baseApi";

export function UserMenu({ name, email }: { name: string; email: string }) {
  const router = useRouter();
  const [open, setOpen] = useState(false);
  const [logout] = useLogoutMutation();

  async function handleLogout() {
    setOpen(false);
    await logout();
    router.push("/");
  }

  return (
    <div className="relative">
      <button
        onClick={() => setOpen((v) => !v)}
        aria-label="Account menu"
        className="flex h-10 w-10 items-center  cursor-pointer justify-center rounded-full bg-primary"
      >
        <User className="h-[18px] w-[18px]" />
      </button>

      {open && (
        <div className="fixed inset-0 z-10" onClick={() => setOpen(false)} />
      )}

      <AnimatePresence>
        {open && (
          <motion.div
            initial={{ opacity: 0, y: -6, scale: 0.97 }}
            animate={{ opacity: 1, y: 0, scale: 1 }}
            exit={{ opacity: 0, y: -6, scale: 0.97 }}
            transition={{ duration: 0.15, ease: [0.16, 1, 0.3, 1] }}
            className="absolute right-0 top-12 z-20 w-64 rounded-2xl bg-[#fefae0] p-1.5 shadow-lg"
          >
            <div className="px-3 py-2.5">
              <p className="truncate text-sm font-medium">{name}</p>
              <p className="truncate text-xs text-muted-foreground">{email}</p>
            </div>

            <div className="my-1 h-px bg-border" />

            <MenuLink
              href="/orders"
              icon={Ticket}
              label="My Orders"
              onClick={() => setOpen(false)}
            />
            <MenuLink
              href="/reservation"
              icon={CalendarCheck}
              label="My Reservations"
              onClick={() => setOpen(false)}
            />
            <MenuLink
              href="/account"
              icon={CogThree}
              label="Account Details"
              onClick={() => setOpen(false)}
            />

            <div className="my-1 h-px bg-border" />

            <button
              onClick={handleLogout}
              className="flex w-full cursor-pointer hover:bg-black hover:text-white items-center gap-2.5 rounded-lg px-3 py-2 text-left text-sm text-destructive hover:bg-destructive/10"
            >
              <Logout className="h-4 w-4" />
              Log out
            </button>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}

function MenuLink({
  href,
  icon: Icon,
  label,
  onClick,
}: {
  href: string;
  icon: React.ComponentType<{ className?: string }>;
  label: string;
  onClick: () => void;
}) {
  return (
    <Link
      href={href}
      onClick={onClick}
      className="flex items-center gap-2.5 rounded-lg px-3 py-2 text-sm hover:bg-black hover:text-white"
    >
      <Icon className="h-4 w-4 text-muted-foreground" />
      {label}
    </Link>
  );
}
