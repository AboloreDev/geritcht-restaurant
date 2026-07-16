"use client";

import Link from "next/link";

import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet";
import { Menu } from "@mynaui/icons-react";

interface MobileMenuProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  navLinks: {
    label: string;
    href: string;
  }[];
}

export default function MobileMenu({
  open,
  onOpenChange,
  navLinks,
}: MobileMenuProps) {
  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetTrigger asChild>
        <button
          className="flex h-11 w-11 items-center justify-center rounded-full border border-border bg-surface text-text-primary transition-all duration-300 hover:border-primary-deep hover:text-primary-deep"
          aria-label="Open Menu"
        >
          <Menu size={22} />
        </button>
      </SheetTrigger>

      <SheetContent
        side="right"
        className="border-l border-border bg-background px-8"
      >
        <SheetHeader className="mb-10">
          <SheetTitle className="font-heading text-3xl text-primary-deep">
            Gericht
          </SheetTitle>
        </SheetHeader>

        <nav>
          <ul className="space-y-8">
            {navLinks.map((item) => (
              <li key={item.href}>
                <Link
                  href={item.href}
                  onClick={() => onOpenChange(false)}
                  className="group relative text-2xl font-medium text-text-primary transition-colors duration-300 hover:text-primary-deep"
                >
                  {item.label}

                  <span className="absolute -bottom-2 left-0 h-[2px] w-0 bg-primary-deep transition-all duration-300 group-hover:w-full" />
                </Link>
              </li>
            ))}
          </ul>

          <div className="mt-14 flex flex-col gap-4">
            <Link
              href="/login"
              onClick={() => onOpenChange(false)}
              className="rounded-full border border-primary-deep py-3 text-center text-text-primary transition-all duration-300 hover:bg-primary-deep hover:text-black"
            >
              Log In
            </Link>

            <Link
              href="#"
              onClick={() => onOpenChange(false)}
              className="rounded-full bg-primary-deep py-3 text-center font-semibold text-black transition-all duration-300 hover:bg-primary"
            >
              Book Table
            </Link>
          </div>
        </nav>
      </SheetContent>
    </Sheet>
  );
}
