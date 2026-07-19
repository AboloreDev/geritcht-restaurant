"use client";

import Subheading from "@/app/(landing)/homepage/SubHeading";
import { ChevronRight } from "@mynaui/icons-react";
import { motion } from "framer-motion";
import Link from "next/link";
import { MenuSearch } from "./MenuSearch";
import { SortMenu } from "./SortMenu";

export default function MenuHeader() {
  return (
    <section className="py-6">
      <div className="mx-auto max-w-7xl px-6 text-center">
        <Subheading className="text-xl">Curated For Every Occasion</Subheading>

        <h1 className="mt-2 text-2xl font-semibold text-primary md:text-4xl">
          Explore Our Menu
        </h1>

        <div className="mt-10 flex items-center justify-between">
          <div className="w-10" aria-hidden />

          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.45 }}
            className="flex items-center justify-center gap-2 text-sm text-text-secondary"
          >
            <Link
              href="/"
              className="transition-colors hover:text-primary-deep"
            >
              Home
            </Link>

            <ChevronRight className="h-4 w-4" />

            <span className="text-primary">Menu</span>
          </motion.div>

          <SortMenu />

          <MenuSearch />
        </div>
      </div>
    </section>
  );
}
