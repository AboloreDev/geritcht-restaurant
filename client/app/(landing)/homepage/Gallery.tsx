"use client";

import { useRef } from "react";
import Image from "next/image";
import { motion } from "framer-motion";
import { Instagram } from "@mynaui/icons-react";

import { Button } from "@/components/ui/button";
import SubHeading from "./SubHeading";

const galleryImages = [
  "/assets/gallery01.png",
  "/assets/gallery02.png",
  "/assets/gallery03.png",
  "/assets/gallery04.png",
];

// Duplicate the array so the loop feels seamless
const loopedImages = [...galleryImages, ...galleryImages];

export default function Gallery() {
  const trackRef = useRef<HTMLDivElement>(null);

  return (
    <section className="overflow-hidden py-24 lg:py-36">
      <div className="mx-auto flex max-w-7xl flex-col gap-16 px-6 lg:flex-row lg:items-center lg:px-10">
        {/* Left */}

        <motion.div
          initial={{ opacity: 0, x: -60 }}
          whileInView={{ opacity: 1, x: 0 }}
          viewport={{ once: true }}
          transition={{
            duration: 0.8,
            ease: [0.22, 1, 0.36, 1],
          }}
          className="max-w-md shrink-0"
        >
          <SubHeading>Instagram</SubHeading>

          <h2 className="mt-4 font-display text-2xl text-primary lg:text-4xl">
            Photo Gallery
          </h2>

          <p className="mt-6 leading-8 text-text-secondary">
            Step inside the world of Gericht. Discover beautifully plated
            signature dishes, handcrafted cocktails, elegant interiors, and
            unforgettable moments shared by our guests.
          </p>

          <Button className="mt-8">Follow Us</Button>
        </motion.div>

        {/* Gallery */}

        <div className="relative flex-1 overflow-hidden">
          {/* Fade */}

          <div className="pointer-events-none absolute left-0 top-0 z-20 h-full w-16 bg-gradient-to-r from-background to-transparent" />

          <div className="pointer-events-none absolute right-0 top-0 z-20 h-full w-16 bg-gradient-to-l from-background to-transparent" />

          {/* Auto-scrolling track */}

          <div ref={trackRef} className="group flex w-max gap-6">
            <motion.div
              className="flex gap-6"
              animate={{ x: ["0%", "-50%"] }}
              transition={{
                duration: 25,
                ease: "linear",
                repeat: Infinity,
              }}
              style={{ willChange: "transform" }}
              whileHover={{ animationPlayState: "paused" }}
            >
              {loopedImages.map((image, i) => (
                <motion.div
                  key={`${image}-${i}`}
                  whileHover={{ y: -10 }}
                  className="group/card relative h-[460px] w-[320px] shrink-0 overflow-hidden rounded-3xl"
                >
                  <Image
                    src={image}
                    alt="Gallery"
                    fill
                    draggable={false}
                    className="object-cover transition duration-700 group-hover/card:scale-110"
                  />

                  <div className="absolute inset-0 flex items-center justify-center bg-black/20 opacity-0 transition duration-500 group-hover/card:opacity-100">
                    <Instagram className="h-10 w-10 text-white" />
                  </div>
                </motion.div>
              ))}
            </motion.div>
          </div>
        </div>
      </div>
    </section>
  );
}
