"use client";

import { cn } from "@/lib/utils";
import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import Link from "next/link";
import Image from "next/image";
import MobileMenu from "./MobileMenu";

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
        "fixed inset-x-0 top-0 z-50 transition-all duration-500",
        scrolled
          ? "border-b border-border shadow-lg backdrop-blur-xl"
          : "bg-transparent",
      )}
    >
      <div className="mx-auto flex h-20 max-w-7xl items-center justify-between px-6 lg:px-10">
        {/* Logo */}

        <Link href="/">
          <motion.div whileHover={{ scale: 1.04 }}>
            <Image
              src="/assets/gericht.png"
              alt="Gericht"
              width={180}
              height={55}
              priority
              className="h-auto w-auto object-contain"
            />
          </motion.div>
        </Link>

        {/* Desktop Navigation */}

        <nav className="hidden lg:block">
          <ul className="flex items-center gap-12">
            {navLinks.map((item) => (
              <li key={item.href}>
                <a
                  href={item.href}
                  className="group relative py-2 text-sm font-medium text-text-primary transition-colors duration-300 hover:text-primary"
                >
                  {item.label}

                  <span className="absolute bottom-0 left-0 h-[2px] w-full origin-left scale-x-0 rounded-full bg-primary transition-transform duration-300 ease-out group-hover:scale-x-100" />
                </a>
              </li>
            ))}
          </ul>
        </nav>

        {/* Right */}

        <div className="hidden items-center gap-3 lg:flex">
          <Link
            href="/login"
            className="bg-[#1b2021] cursor-pointer rounded-2xl px-5 py-2 text-base font-medium text-primary transition-all duration-300 hover:text-primary-deep"
          >
            Log In
          </Link>

          <Link
            href="#"
            className="bg-[#1b2021] cursor-pointer rounded-2xl px-6 py-2 text-base font-semibold text-primary shadow-md transition-all duration-300 hover:shadow-xl"
          >
            Book Table
          </Link>
        </div>

        {/* Mobile */}

        <div className="lg:hidden">
          <MobileMenu open={open} onOpenChange={setOpen} navLinks={navLinks} />
        </div>
      </div>
    </motion.header>
  );
}
