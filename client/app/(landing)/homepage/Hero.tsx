"use client";

import Image from "next/image";
import { motion } from "framer-motion";

import { Button } from "@/components/ui/button";
import Subheading from "./SubHeading";
import Link from "next/link";

export default function Hero() {
  return (
    <header
      id="home"
      className="relative overflow-hidden bg-background pt-32 pb-20"
    >
      <div className="mx-auto grid min-h-[85vh] max-w-7xl items-center gap-16 px-6 lg:grid-cols-2 lg:px-10">
        {/* Left Content */}

        <motion.div
          initial={{ opacity: 0, x: -60 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ duration: 0.7 }}
          className="space-y-8"
        >
          <Subheading>Chase the New Flavor</Subheading>

          <h1 className="font-display text-6xl leading-none text-text-primary md:text-7xl xl:text-8xl">
            The Key to
            <span className="block text-primary-deep">Fine Dining</span>
          </h1>

          <p className="max-w-xl text-lg leading-8 text-text-secondary">
            Sit tellus lobortis sed senectus vivamus molestie. Condimentum
            volutpat morbi facilisis quam scelerisque sapien. Et, penatibus
            aliquam amet tellus.
          </p>

          <div className="flex items-center gap-4">
            <Link
              href={"/"}
              className="rounded-full bg-primary-deep px-6 py-3 text-base font-semibold text-black transition-all duration-300 hover:-translate-y-1 hover:bg-primary hover:shadow-xl"
            >
              Explore Menu
            </Link>

            <Link
              href={"/"}
              className="rounded-full border border-border px-6 py-3 text-text-primary hover:border-primary-deep hover:text-primary-deep"
            >
              Book Table
            </Link>
          </div>
        </motion.div>

        {/* Right Image */}

        <motion.div
          initial={{ opacity: 0, x: 60 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ duration: 0.8 }}
          className="relative flex justify-center"
        >
          <div className="absolute h-[420px] w-[420px] rounded-full bg-primary-deep/10 blur-3xl" />

          <Image
            src="/assets/welcome.png"
            alt="Fine dining"
            width={620}
            height={760}
            priority
            className="relative z-10 h-auto w-full max-w-xl object-contain"
          />
        </motion.div>
      </div>
    </header>
  );
}
