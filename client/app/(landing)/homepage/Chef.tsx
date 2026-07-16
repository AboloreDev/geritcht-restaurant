"use client";

import Image from "next/image";
import { motion } from "framer-motion";

import Subheading from "./SubHeading";

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

export default function Chef() {
  return (
    <section className="relative overflow-hidden py-24 lg:py-36">
      <div className="mx-auto grid max-w-7xl items-center gap-16 px-6 lg:grid-cols-2 lg:px-10">
        {/* Chef Image */}

        <motion.div
          // @ts-expect-error "<>"
          variants={fadeLeft}
          initial="hidden"
          whileInView="show"
          viewport={{ once: true }}
          className="flex justify-center"
        >
          <Image
            src="/assets/chef.png"
            alt="Executive Chef"
            width={520}
            height={700}
            className="h-auto w-full max-w-md object-contain"
            priority
          />
        </motion.div>

        {/* Content */}

        <motion.div
          // @ts-expect-error "<>"
          variants={fadeRight}
          initial="hidden"
          whileInView="show"
          viewport={{ once: true }}
        >
          <Subheading>Chef&apos;s Philosophy</Subheading>

          <h2 className="mt-4 font-display text-xl text-primary lg:text-3xl">
            Cooking with Purpose,
            <br />
            Serving with Passion
          </h2>

          <div className="mt-10 space-y-8">
            <div className="flex items-start gap-4">
              <Image src="/assets/quote.png" alt="" width={42} height={42} />

              <p className="leading-8 text-text-secondary">
                Great food begins with respect—for the ingredients, for the
                craft, and for the people gathered around the table. Every dish
                we serve is prepared with precision, creativity, and genuine
                care.
              </p>
            </div>

            <p className="leading-8 text-text-secondary">
              At Gericht, we blend timeless culinary traditions with modern
              techniques to create memorable dining experiences. From the first
              bite to the final course, our goal is to make every visit feel
              personal, refined, and unforgettable.
            </p>
          </div>

          <div className="mt-12">
            <h3 className="font-heading text-xl text-primary">
              Chef Adewale Johnson
            </h3>

            <p className="mt-2 text-text-muted">Executive Chef & Founder</p>

            <Image
              src="/assets/sign.png"
              alt="Chef Signature"
              width={220}
              height={90}
              className="mt-8 h-auto w-44"
            />
          </div>
        </motion.div>
      </div>
    </section>
  );
}
