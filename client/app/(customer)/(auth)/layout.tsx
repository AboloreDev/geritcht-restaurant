// src/app/(auth)/components/AuthLayout.tsx
"use client";

import { motion } from "framer-motion";
import type { ReactNode } from "react";

export default function AuthLayout({
  imagePanel,
  children,
}: {
  imagePanel?: ReactNode;
  children: ReactNode;
}) {
  return (
    <div className="relative flex min-h-screen w-full items-center justify-center bg-[url('/assets/bg.png')] bg-cover bg-center bg-fixed">
      <motion.div
        initial={{ opacity: 0, y: 16 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5, ease: [0.16, 1, 0.3, 1] }}
        className="flex h-[640px] w-full max-w-5xl overflow-hidden bg-primary rounded-3xl shadow-2xl"
      >
        {/* left panel — image, exactly half width, full card height */}
        {imagePanel && (
          <div className="relative hidden h-full w-full md:block">
            {imagePanel}
          </div>
        )}

        {/* right panel — form, exactly half width, own scroll only if content genuinely overflows */}
        <div className="flex h-full w-full flex-col justify-center">
          {children}
        </div>
      </motion.div>
    </div>
  );
}
