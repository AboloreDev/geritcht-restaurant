"use client";

import Image from "next/image";
import { AnimatePresence, motion } from "framer-motion";
import { HugeiconsIcon } from "@hugeicons/react";
import { Cancel01Icon } from "@hugeicons/core-free-icons";
import { Badge } from "@/components/ui/badge";
import { MenuImage } from "@/app/state/types/menuTypes";

type Props = {
  open: boolean;
  image: MenuImage | null;
  onOpenChange: (open: boolean) => void;
};

export default function ImageLightbox({ open, image, onOpenChange }: Props) {
  return (
    <AnimatePresence>
      {open && image && (
        <>
          {/* backdrop — click to close */}
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            transition={{ duration: 0.15 }}
            className="fixed inset-0 z-50 bg-black/80 backdrop-blur-sm"
            onClick={() => onOpenChange(false)}
          />

          {/* image, centered, click-safe (doesn't close on image click itself) */}
          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            exit={{ opacity: 0, scale: 0.95 }}
            transition={{ duration: 0.2, ease: [0.16, 1, 0.3, 1] }}
            className="fixed inset-0 z-50 flex items-center justify-center p-6"
            onClick={() => onOpenChange(false)}
          >
            <div
              className="relative max-h-[85vh] w-full max-w-2xl overflow-hidden rounded-2xl bg-black"
              onClick={(e) => e.stopPropagation()}
            >
              <div className="relative aspect-square w-full">
                <Image
                  src={image.url}
                  alt={image.alt_text || "Menu image"}
                  fill
                  className="object-contain"
                  sizes="672px"
                />
              </div>

              {image.is_primary && (
                <Badge className="absolute left-4 top-4 rounded-full bg-amber-500 text-white hover:bg-amber-500">
                  ⭐ Primary
                </Badge>
              )}

              <button
                onClick={() => onOpenChange(false)}
                aria-label="Close preview"
                className="absolute right-4 top-4 flex h-9 w-9 items-center justify-center rounded-full bg-black/50 text-white hover:bg-black/70"
              >
                <HugeiconsIcon icon={Cancel01Icon} strokeWidth={2} size={18} />
              </button>
            </div>
          </motion.div>
        </>
      )}
    </AnimatePresence>
  );
}
