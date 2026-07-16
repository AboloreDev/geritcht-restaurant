"use client";

import React, { useRef, useState } from "react";
import { motion } from "framer-motion";
import { Play, Pause } from "@mynaui/icons-react";

export default function VideoIntro() {
  const [playing, setPlaying] = useState(false);
  const vidRef = useRef<HTMLVideoElement>(null);

  function handleVideo() {
    setPlaying((previousPlayVideo) => !previousPlayVideo);

    if (playing) {
      vidRef.current?.pause();
    } else {
      vidRef.current?.play();
    }
  }

  return (
    <section className="relative h-screen overflow-hidden">
      {/* Video */}

      <video
        ref={vidRef}
        className="absolute inset-0 h-full w-full object-cover"
        loop
        muted
        playsInline
        preload="metadata"
      >
        <source src="/assets/meal.mp4" type="video/mp4" />
      </video>

      {/* Gradient Overlay */}

      <div className="absolute inset-0 bg-black/55" />

      <div className="absolute inset-0 bg-gradient-to-t from-background via-black/30 to-black/40" />

      {/* Content */}

      <div className="relative z-10 flex h-full flex-col items-center justify-center px-6 text-center">
        <motion.h2
          initial={{ opacity: 0, y: 30 }}
          whileInView={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.7 }}
          className="font-display text-2xl text-white md:text-4xl"
        >
          Experience Every Moment
        </motion.h2>

        <motion.p
          initial={{ opacity: 0, y: 35 }}
          whileInView={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.15 }}
          className="mt-6 max-w-2xl text-md leading-8 text-neutral-300"
        >
          Watch our chefs transform the finest ingredients into unforgettable
          culinary experiences, crafted with passion and served with elegance.
        </motion.p>

        {/* Play Button */}

        <motion.button
          whileHover={{
            scale: 1.08,
          }}
          whileTap={{
            scale: 0.92,
          }}
          onClick={handleVideo}
          className="mt-12 flex h-20 w-20 items-center justify-center rounded-full border border-primary/40 bg-white/10 backdrop-blur-xl transition-all hover:border-primary hover:bg-primary/15"
        >
          {playing ? (
            <Pause className="h-8 w-8 text-primary" />
          ) : (
            <Play className="ml-1 h-8 w-8 text-primary" />
          )}
        </motion.button>
      </div>
    </section>
  );
}
