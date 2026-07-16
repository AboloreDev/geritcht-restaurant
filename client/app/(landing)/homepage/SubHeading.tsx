"use client";

import Image from "next/image";
import { cn } from "@/lib/utils";

interface SubheadingProps {
  children: React.ReactNode;
  className?: string;
}

export default function Subheading({ children, className }: SubheadingProps) {
  return (
    <div className="space-y-3">
      <p
        className={cn(
          "font-display text-xl font-medium tracking-wide text-primary",
          className,
        )}
      >
        {children}
      </p>

      <Image
        src="/assets/spoon.svg"
        alt="Decorative spoon"
        width={42}
        height={10}
        priority
        className="h-auto w-10 text-primary"
      />
    </div>
  );
}
