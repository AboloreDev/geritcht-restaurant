"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";

import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet";
import { Menu, ShoppingBag } from "@mynaui/icons-react";
import Image from "next/image";
import { Button } from "@/components/ui/button";
import { openBookingModal } from "@/app/state/slices/reservationSlice";
import { useAppDispatch } from "@/app/state/redux";
import { useLogoutMutation } from "@/app/state/api/baseApi";
import { useAuth } from "@/app/hooks/isAuthenticated";

interface MobileMenuProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  navLinks: {
    label: string;
    href: string;
  }[];
}

const accountLinks = [
  { label: "My Orders", href: "/account/orders" },
  { label: "My Reservations", href: "/account/reservations" },
  { label: "Account Details", href: "/account/settings" },
];

export default function MobileMenu({
  open,
  onOpenChange,
  navLinks,
}: MobileMenuProps) {
  const dispatch = useAppDispatch();
  const router = useRouter();
  const { isAuthenticated, user } = useAuth();
  const [logout] = useLogoutMutation();

  // hydration-safety guard, same pattern as the desktop nav
  const [hasMounted, setHasMounted] = useState(false);
  useEffect(() => {
    // eslint-disable-next-line react-hooks/set-state-in-effect
    setHasMounted(true);
  }, []);

  const handleOpenBookingModal = () => {
    dispatch(openBookingModal());
    onOpenChange(false);
  };

  async function handleLogout() {
    onOpenChange(false);
    await logout();
    router.push("/");
  }

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetTrigger>
        <SheetTitle
          className="flex h-11 w-11 cursor-pointer items-center justify-center text-primary transition-all duration-300 hover:border-primary-deep hover:text-primary-deep"
          aria-label="Open Menu"
        >
          <Menu size={30} />
        </SheetTitle>
      </SheetTrigger>

      <SheetContent
        side="right"
        className="border-l border-border bg-background px-8"
      >
        <SheetHeader className="mb-10">
          <SheetTitle className="font-heading text-3xl text-primary-deep">
            <Image
              src="/assets/gericht.png"
              alt="Gericht"
              width={100}
              height={55}
              priority
              className="h-auto w-auto object-contain"
            />
          </SheetTitle>
        </SheetHeader>

        <nav>
          <ul className="space-y-8">
            {navLinks.map((item) => (
              <li key={item.href}>
                <Link
                  href={item.href}
                  onClick={() => onOpenChange(false)}
                  className="group relative text-xl font-medium text-text-primary transition-colors duration-300 hover:text-primary-deep"
                >
                  {item.label}

                  <span className="absolute -bottom-2 left-0 h-[2px] w-0 bg-primary-deep transition-all duration-300 group-hover:w-full" />
                </Link>
              </li>
            ))}
          </ul>

          <div className="mt-14 flex flex-col gap-4">
            {hasMounted && isAuthenticated ? (
              <button
                onClick={handleLogout}
                className="rounded-full border border-destructive py-3 text-center text-destructive transition-all duration-300 hover:bg-destructive hover:text-white"
              >
                Log out
              </button>
            ) : (
              <Link
                href="/login"
                onClick={() => onOpenChange(false)}
                className="rounded-full border border-primary-deep py-3 text-center text-text-primary transition-all duration-300 hover:bg-primary-deep hover:text-black"
              >
                Log In
              </Link>
            )}

            <Button
              onClick={handleOpenBookingModal}
              className="rounded-full bg-primary-deep py-3 text-center font-semibold text-black transition-all duration-300 hover:bg-primary"
            >
              Book Table
            </Button>
          </div>
        </nav>
      </SheetContent>
    </Sheet>
  );
}
