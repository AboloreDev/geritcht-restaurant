"use client";

import Image from "next/image";
import { motion } from "framer-motion";

import SubHeading from "./SubHeading";
import AwardCard from "./AwardCard";

export const awards = [
  {
    imgUrl: "/assets/award01.png",
    title: "Signature Dining Experience",
    subtitle:
      "Crafting unforgettable moments through elegant cuisine, impeccable service, and refined ambience.",
  },
  {
    imgUrl: "/assets/award02.png",
    title: "Seasonal Chef's Selection",
    subtitle:
      "A celebration of locally sourced ingredients transformed into contemporary culinary masterpieces.",
  },
  {
    imgUrl: "/assets/award03.png",
    title: "Guest Favourite",
    subtitle:
      "Consistently praised for exceptional hospitality, beautifully plated dishes, and memorable experiences.",
  },
  {
    imgUrl: "/assets/award05.png",
    title: "Curated Wine Collection",
    subtitle:
      "An expertly selected collection of international wines paired to complement every signature dish.",
  },
];

const container = {
  hidden: {},
  show: {
    transition: {
      staggerChildren: 0.15,
    },
  },
};

const fadeUp = {
  hidden: {
    opacity: 0,
    y: 40,
  },
  show: {
    opacity: 1,
    y: 0,
    transition: {
      duration: 0.7,
      ease: [0.22, 1, 0.36, 1],
    },
  },
};

export default function Awards() {
  return (
    <section className="py-24 lg:py-36">
      <div className="mx-auto grid max-w-7xl items-center gap-20 px-6 lg:grid-cols-2 lg:px-10">
        {/* Left */}

        <motion.div
          variants={container}
          initial="hidden"
          whileInView="show"
          viewport={{ once: true }}
        >
          <motion.div
            // @ts-expect-error "<>"
            variants={fadeUp}
          >
            <SubHeading>Awards & Recognition</SubHeading>

            <h2 className="mt-4 font-display text-2xl text-primary lg:text-4xl">
              Our Laurels
            </h2>
          </motion.div>

          <motion.div
            variants={container}
            className="mt-14 grid gap-8 sm:grid-cols-2"
          >
            {awards.map((award) => (
              <AwardCard key={award.title} {...award} />
            ))}
          </motion.div>
        </motion.div>

        {/* Right */}

        <motion.div
          // @ts-expect-error "<>"
          variants={fadeUp}
          initial="hidden"
          whileInView="show"
          viewport={{ once: true }}
          className="flex justify-center"
        >
          <Image
            src="/assets/laurels.png"
            alt="Restaurant Awards"
            width={520}
            height={700}
            className="h-auto w-full max-w-md object-contain"
          />
        </motion.div>
      </div>
    </section>
  );
}
