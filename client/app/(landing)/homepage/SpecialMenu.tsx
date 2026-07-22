"use client";

import Image from "next/image";
import { motion } from "framer-motion";

import { Button } from "@/components/ui/button";
import Subheading from "./SubHeading";
import ItemList from "./ItemList";
import Link from "next/link";

const fadeUp = {
  hidden: { opacity: 0, y: 40 },
  show: {
    opacity: 1,
    y: 0,
    transition: {
      duration: 0.7,
      ease: [0.22, 1, 0.36, 1],
    },
  },
};

export const chefSpecials = [
  {
    title: "Grilled Atlantic Salmon",
    price: "₦28,500",
    tags: "Salmon • Lemon Butter • Gluten Free • Contains Fish",
  },
  {
    title: "Prime Ribeye Steak",
    price: "₦34,000",
    tags: "400g • Garlic Herb Butter • Medium Rare Recommended",
  },
  {
    title: "Seafood Linguine",
    price: "₦24,500",
    tags: "Shrimp • Mussels • Parmesan • Contains Shellfish & Dairy",
  },
  {
    title: "Chicken Supreme",
    price: "₦19,800",
    tags: "Roasted Chicken • Cream Sauce • Contains Dairy",
  },
  {
    title: "Truffle Mushroom Risotto",
    price: "₦18,500",
    tags: "Vegetarian • Parmesan • Contains Dairy",
  },
];

export const signatureDrinks = [
  {
    title: "Tropical Sunset",
    price: "₦8,500",
    tags: "Pineapple • Passion Fruit • Fresh Mint • Non-Alcoholic",
  },
  {
    title: "Lagos Mule",
    price: "₦11,500",
    tags: "Vodka • Ginger Beer • Lime",
  },
  {
    title: "Berry Bliss Mojito",
    price: "₦10,500",
    tags: "Mixed Berries • Mint • White Rum",
  },
  {
    title: "Smoked Old Fashioned",
    price: "₦14,000",
    tags: "Bourbon • Bitters • Orange Peel",
  },
  {
    title: "Golden Espresso Martini",
    price: "₦13,500",
    tags: "Vodka • Espresso • Coffee Liqueur • Contains Caffeine",
  },
];

export default function SpecialMenu() {
  return (
    <section id="menu" className="py-24 lg:py-36">
      <div className="mx-auto flex max-w-7xl flex-col items-center px-6 lg:px-10">
        {/* Heading */}

        <motion.div
          // @ts-expect-error "<>"
          variants={fadeUp}
          initial="hidden"
          whileInView="show"
          viewport={{ once: true }}
          className="mb-16 text-center"
        >
          <Subheading>Curated for Every Occasion</Subheading>

          <h2 className="mt-4 font-display text-2xl text-primary lg:text-4xl">
            Today&apos;s Special
          </h2>
        </motion.div>

        {/* Menu */}

        <div className="grid w-full items-start gap-16 lg:grid-cols-[1fr_auto_1fr]">
          {/* Wine */}

          <motion.div
            // @ts-expect-error "<>"
            variants={fadeUp}
            initial="hidden"
            whileInView="show"
            viewport={{ once: true }}
            className="space-y-8"
          >
            <h3 className="font-heading text-xl text-white">Wine & Beer</h3>

            <div className="space-y-6">
              {chefSpecials.map((specials, index) => (
                <ItemList
                  key={`${specials.title}-${index}`}
                  title={specials.title}
                  price={specials.price}
                  tags={specials.tags}
                />
              ))}
            </div>
          </motion.div>

          {/* Image */}

          <motion.div
            initial={{ opacity: 0, scale: 0.9 }}
            whileInView={{ opacity: 1, scale: 1 }}
            viewport={{ once: true }}
            transition={{
              duration: 0.8,
              ease: [0.22, 1, 0.36, 1],
            }}
            whileHover={{
              y: -8,
            }}
            className="flex justify-center"
          >
            <Image
              src={"/assets/menu.png"}
              alt="Today's Special Menu"
              width={420}
              height={650}
              className="h-auto w-full max-w-sm object-contain"
            />
          </motion.div>

          {/* Cocktails */}

          <motion.div
            // @ts-expect-error "<>"
            variants={fadeUp}
            initial="hidden"
            whileInView="show"
            viewport={{ once: true }}
            transition={{ delay: 0.15 }}
            className="space-y-8"
          >
            <h3 className="font-heading text-xl text-white">
              Signature Cocktails
            </h3>

            <div className="space-y-6">
              {signatureDrinks.map((drinks, index) => (
                <ItemList
                  key={`${drinks.title}-${index}`}
                  title={drinks.title}
                  price={drinks.price}
                  tags={drinks.tags}
                />
              ))}
            </div>
          </motion.div>
        </div>

        {/* CTA */}

        <motion.div
          initial={{ opacity: 0, y: 25 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ delay: 0.35 }}
          className="mt-20"
        >
          <Link
            href={"/menu"}
            className="rounded-full bg-primary text-black px-4 py-4"
          >
            View Full Menu
          </Link>
        </motion.div>
      </div>
    </section>
  );
}
