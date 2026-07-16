"use client";

import Image from "next/image";
import { motion } from "framer-motion";

import SubHeading from "./SubHeading";
import { Button } from "@/components/ui/button";
import { Mail, MapPin, Clock3, Telephone } from "@mynaui/icons-react";

const fadeLeft = {
  hidden: { opacity: 0, x: -60 },
  show: {
    opacity: 1,
    x: 0,
    transition: {
      duration: 0.8,
      ease: [0.22, 1, 0.36, 1],
    },
  },
};

const fadeRight = {
  hidden: { opacity: 0, x: 60 },
  show: {
    opacity: 1,
    x: 0,
    transition: {
      duration: 0.8,
      ease: [0.22, 1, 0.36, 1],
    },
  },
};

export default function Contact() {
  return (
    <section id="contact" className="py-24 lg:py-36">
      <div className="mx-auto grid max-w-7xl items-center gap-20 px-6 lg:grid-cols-2 lg:px-10">
        {/* Left */}

        <motion.div
          // @ts-expect-error "<>"
          variants={fadeLeft}
          initial="hidden"
          whileInView="show"
          viewport={{ once: true }}
        >
          <SubHeading>Visit Us</SubHeading>

          <h2 className="mt-4 font-display text-2xl text-primary lg:text-4xl">
            Find Us
          </h2>

          <p className="mt-6 max-w-lg leading-8 text-text-secondary">
            Whether you&apos;re planning a quiet dinner, celebrating a special
            occasion, or simply craving exceptional cuisine, we&apos;d love to
            welcome you to Gericht.
          </p>

          <div className="mt-10 space-y-6">
            <div className="flex items-start gap-4">
              <MapPin className="mt-1 h-5 w-5 text-primary" />

              <div>
                <h4 className="font-semibold text-primary">Address</h4>
                <p className="text-text-secondary">
                  12 Admiralty Way, Lekki Phase 1,
                  <br />
                  Lagos, Nigeria.
                </p>
              </div>
            </div>

            <div className="flex items-start gap-4">
              <Clock3 className="mt-1 h-5 w-5 text-primary" />

              <div>
                <h4 className="font-semibold text-primary">Opening Hours</h4>

                <p className="text-text-secondary">
                  Monday – Friday
                  <br />
                  11:00 AM – 11:00 PM
                </p>

                <p className="mt-2 text-text-secondary">
                  Saturday – Sunday
                  <br />
                  10:00 AM – Midnight
                </p>
              </div>
            </div>

            <div className="flex items-center gap-4">
              <Telephone className="h-5 w-5 text-primary" />

              <p className="text-text-secondary">+234 812 345 6789</p>
            </div>

            <div className="flex items-center gap-4">
              <Mail className="h-5 w-5 text-primary" />

              <p className="text-text-secondary">reservations@gericht.ng</p>
            </div>
          </div>

          <Button className="mt-10 rounded-full px-8">Reserve a Table</Button>
        </motion.div>

        {/* Right */}

        <motion.div
          // @ts-expect-error "<>"
          variants={fadeRight}
          initial="hidden"
          whileInView="show"
          viewport={{ once: true }}
          className="flex justify-center"
        >
          <div className="overflow-hidden rounded-3xl border border-border bg-surface shadow-lg">
            <Image
              src="/assets/findus.png"
              alt="Gericht Restaurant"
              width={560}
              height={700}
              className="h-auto w-full object-cover transition duration-700 hover:scale-105"
              priority
            />
          </div>
        </motion.div>
      </div>
    </section>
  );
}
