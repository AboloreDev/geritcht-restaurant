"use client";

import Image from "next/image";
import { motion } from "framer-motion";

interface AwardCardProps {
  imgUrl: string;
  title: string;
  subtitle: string;
}

export default function AwardCard({ imgUrl, title, subtitle }: AwardCardProps) {
  return (
    <motion.div
      whileHover={{
        y: -8,
        scale: 1.03,
      }}
      transition={{
        duration: 0.3,
      }}
      className="group flex gap-5 rounded-3xl border border-border p-3 transition-colors"
    >
      <div className="flex h-16 w-16 shrink-0 items-center justify-center rounded-2xl">
        <Image
          src={imgUrl}
          alt={title}
          width={46}
          height={46}
          className="object-contain transition-transform duration-300 group-hover:rotate-6 group-hover:scale-110"
        />
      </div>

      <div>
        <h3 className="font-display text-xl text-primary">{title}</h3>

        <p className="mt-2 text-sm leading-7 text-text-secondary">{subtitle}</p>
      </div>
    </motion.div>
  );
}
