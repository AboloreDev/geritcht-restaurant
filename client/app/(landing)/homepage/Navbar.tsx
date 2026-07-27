"use client";

import { cn } from "@/lib/utils";
import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import Link from "next/link";
import Image from "next/image";
import { ShoppingBag } from "@mynaui/icons-react";
import MobileMenu from "./MobileMenu";
import { useAppDispatch } from "@/app/state/redux";
import { openBookingModal } from "@/app/state/slices/reservationSlice";
import { Button } from "@/components/ui/button";
import { useAuth } from "@/app/hooks/isAuthenticated";
import { UserMenu } from "@/components/code/UserMenu";
import { openCartDrawer } from "@/app/state/slices/cartSlice";
import { useGetCartQuery } from "@/app/state/api/cartApi";

const navLinks = [
  { label: "Home", href: "#home" },
  { label: "About", href: "#about" },
  { label: "Menu", href: "/menu" },
  { label: "Awards", href: "#awards" },
  { label: "Contact", href: "#contact" },
];

export default function Navbar() {
  const [open, setOpen] = useState(false);
  const [scrolled, setScrolled] = useState(false);
  const dispatch = useAppDispatch();
  const [hasMounted, setHasMounted] = useState(false);
  const { isAuthenticated, user } = useAuth();

  const { data: cartData } = useGetCartQuery(undefined, {
    skip: !isAuthenticated,
  });
  const cart = cartData?.data;

  useEffect(() => {
    // eslint-disable-next-line react-hooks/set-state-in-effect
    setHasMounted(true);
  }, []);

  useEffect(() => {
    const handleScroll = () => setScrolled(window.scrollY > 20);

    window.addEventListener("scroll", handleScroll);

    return () => window.removeEventListener("scroll", handleScroll);
  }, []);

  return (
    <motion.header
      initial={{ y: -80, opacity: 0 }}
      animate={{ y: 0, opacity: 1 }}
      transition={{
        duration: 0.6,
        ease: [0.22, 1, 0.36, 1],
      }}
      className={cn(
        "sticky inset-x-0 top-0 z-50 transition-all duration-500",
        scrolled
          ? "border-b border-white/10 bg-[#111315]/70 shadow-xl backdrop-blur-2xl"
          : "bg-transparent",
      )}
    >
      <div className="mx-auto flex h-20 max-w-7xl items-center justify-between px-6 lg:px-10">
        {/* Logo */}
        <Link href="/">
          <motion.div whileHover={{ scale: 1.03 }}>
            <Image
              src="/assets/gericht.png"
              alt="Gericht"
              width={180}
              height={55}
              priority
              className="h-28 w-28 md:w-auto md:h-auto object-contain"
            />
          </motion.div>
        </Link>

        {/* Desktop Navigation */}
        <nav className="hidden">
          <ul className="flex items-center gap-10">
            {navLinks.map((item) => (
              <li key={item.href}>
                <Link
                  href={item.href}
                  className="group relative py-2 text-sm font-medium text-text-primary transition-colors duration-300 hover:text-primary"
                >
                  {item.label}

                  <span className="absolute bottom-0 left-0 h-[2px] w-full origin-left scale-x-0 rounded-full bg-primary transition-transform duration-300 ease-out group-hover:scale-x-100" />
                </Link>
              </li>
            ))}
          </ul>
        </nav>

        {/* Right */}
        <div className="flex items-center gap-4">
          {hasMounted && isAuthenticated && user ? (
            <>
              <button
                aria-label="Cart"
                onClick={() => dispatch(openCartDrawer())}
                className="relative flex h-10 w-10 items-center justify-center rounded-2xl border border-primary/20 bg-[#1b2021] text-primary transition-all duration-300 hover:scale-105 hover:bg-primary hover:text-black"
              >
                <ShoppingBag className="h-[18px] w-[18px]" />

                {
                  // @ts-expect-error "<>"
                  cart?.item_count > 0 && (
                    <span className="absolute -right-1 -top-1 flex h-4 w-4 items-center justify-center rounded-full bg-primary text-[10px] text-black">
                      {cart?.item_count}
                    </span>
                  )
                }
              </button>

              <UserMenu name={user.first_name} email={user.email} />
            </>
          ) : (
            <Link
              href="/login"
              className="hidden md:block rounded-2xl border border-primary/20 bg-[#1b2021] px-5 py-2 text-sm font-medium text-primary transition-all duration-300 hover:bg-primary hover:text-black"
            >
              Log In
            </Link>
          )}

          <Button
            onClick={() => dispatch(openBookingModal())}
            className="hidden md:flex rounded-2xl bg-primary px-6 py-2 font-semibold text-black shadow-md transition-all duration-300 hover:scale-[1.02] hover:bg-primary/90 hover:shadow-xl"
          >
            Book Table
          </Button>
        </div>

        {/* Mobile */}
        <div className="md:hidden">
          <MobileMenu open={open} onOpenChange={setOpen} navLinks={navLinks} />
        </div>
      </div>
    </motion.header>
  );
}
