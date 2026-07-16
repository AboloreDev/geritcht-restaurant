"use client";

import { motion } from "framer-motion";
import { Mail } from "@mynaui/icons-react";

import SubHeading from "./SubHeading";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

export default function Newsletter() {
  return (
    <motion.section
      initial={{ opacity: 0, y: 70 }}
      whileInView={{ opacity: 1, y: 0 }}
      viewport={{ once: true }}
      transition={{
        duration: 0.8,
        ease: [0.22, 1, 0.36, 1],
      }}
      className="mx-auto -mb-28 max-w-5xl px-6 lg:px-10"
    >
      <div className="rounded-[32px] p-8 shadow-xl backdrop-blur-xl lg:p-14">
        <div className="text-center">
          <SubHeading>Newsletter</SubHeading>

          <h2 className="mt-4 font-display text-2xl text-primary lg:text-4xl">
            Stay Connected
          </h2>

          <p className="mx-auto mt-5 max-w-2xl leading-8 text-text-secondary">
            Be the first to discover seasonal menus, exclusive chef&apos;s
            specials, private dining events, and special offers curated
            exclusively for our guests.
          </p>
        </div>

        <form className="mx-auto mt-12 flex max-w-3xl flex-col gap-4 sm:flex-row">
          <div className="relative flex-1">
            <Mail className="absolute left-5 top-1/2 h-5 w-5 -translate-y-1/2 text-text-muted" />

            <Input
              type="email"
              placeholder="Enter your email address"
              className="h-14 rounded-full border-border bg-background pl-14 text-base placeholder:text-text-muted focus-visible:ring-primary"
            />
          </div>

          <Button
            type="submit"
            className="h-14 rounded-full px-10 text-base font-semibold"
          >
            Subscribe
          </Button>
        </form>

        <p className="mt-6 text-center text-sm text-text-muted">
          No spam. Just exclusive dining experiences, seasonal updates, and
          special invitations.
        </p>
      </div>
    </motion.section>
  );
}
