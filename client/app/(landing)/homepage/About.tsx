"use client";

import Image from "next/image";
import { motion } from "framer-motion";
import Subheading from "./SubHeading";
import { Button } from "@/components/ui/button";

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

export default function About() {
  return (
    <section id="about" className="relative overflow-hidden py-24 lg:py-36">
      {/* Decorative Background */}
      <motion.div
        className="pointer-events-none absolute inset-0 flex items-center justify-center"
        animate={{
          y: [-8, 8, -8],
          rotate: [-1, 1, -1],
        }}
        transition={{
          duration: 12,
          repeat: Infinity,
          ease: "easeInOut",
        }}
      >
        <Image
          src="/assets/G.png"
          alt=""
          width={420}
          height={520}
          priority
          className="select-none opacity-5"
        />
      </motion.div>

      <div className="relative z-10 mx-auto grid max-w-7xl items-center gap-16 px-6 lg:grid-cols-[1fr_auto_1fr] lg:px-10">
        {/* About */}
        <motion.div
          // @ts-expect-error "<>"
          variants={fadeLeft}
          initial="hidden"
          whileInView="show"
          viewport={{ once: true, amount: 0.3 }}
          className="flex flex-col items-end text-right"
        >
          <Subheading className="text-5xl">About Us</Subheading>

          <p className="mt-6 max-w-md leading-8 text-text-secondary">
            At Gericht, every meal is a celebration of craftsmanship,
            creativity, and exceptional hospitality. We source the finest
            ingredients and transform them into unforgettable dining
            experiences.
          </p>

          <motion.div whileHover={{ y: -4 }} whileTap={{ scale: 0.96 }}>
            <Button className="mt-8">Discover More</Button>
          </motion.div>
        </motion.div>

        {/* Knife */}
        <motion.div
          initial={{ opacity: 0, scale: 0.9 }}
          whileInView={{ opacity: 1, scale: 1 }}
          viewport={{ once: true }}
          transition={{
            duration: 0.8,
            ease: [0.22, 1, 0.36, 1],
          }}
          animate={{
            y: [-12, 12, -12],
          }}
          className="hidden justify-center lg:flex"
        >
          <Image
            src="/assets/knife.png"
            alt="Chef's Knife"
            width={90}
            height={760}
            className="h-[650px] w-auto object-contain"
          />
        </motion.div>

        {/* History */}
        <motion.div
          // @ts-expect-error "<>"
          variants={fadeRight}
          initial="hidden"
          whileInView="show"
          viewport={{ once: true, amount: 0.3 }}
          className="flex flex-col items-start text-left"
        >
          <Subheading className="text-5xl">Our History</Subheading>

          <p className="mt-6 max-w-md leading-8 text-text-secondary">
            From a simple vision to redefine fine dining, Gericht has become a
            destination where timeless tradition meets modern culinary
            innovation. Every dish tells a story of passion, excellence, and
            authenticity.
          </p>

          <motion.div whileHover={{ y: -4 }} whileTap={{ scale: 0.96 }}>
            <Button className="mt-8">Our Journey</Button>
          </motion.div>
        </motion.div>
      </div>
    </section>
  );
}
