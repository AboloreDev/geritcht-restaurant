"use client";

import Image from "next/image";
import { motion } from "framer-motion";
import { HugeiconsIcon } from "@hugeicons/react";
import { Delete02Icon, Image01Icon } from "@hugeicons/core-free-icons";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { useEffect, useState } from "react";

interface Props {
  file: File;
  index: number;
  onRemove: () => void;
}

export default function ImagePreviewCard({ file, index, onRemove }: Props) {
  const [preview, setPreview] = useState("");

  useEffect(() => {
    const url = URL.createObjectURL(file);

    setPreview(url);

    return () => URL.revokeObjectURL(url);
  }, [file]);

  const size = `${(file.size / 1024 / 1024).toFixed(2)} MB`;

  return (
    <motion.div
      layout
      initial={{
        opacity: 0,
        scale: 0.92,
        y: 20,
      }}
      animate={{
        opacity: 1,
        scale: 1,
        y: 0,
      }}
      exit={{
        opacity: 0,
        scale: 0.92,
      }}
      transition={{
        duration: 0.25,
      }}
      className="overflow-hidden rounded-3xl border bg-white shadow-sm"
    >
      <div className="relative aspect-square overflow-hidden">
        <Image src={preview} alt={file.name} fill className="object-cover" />

        {index === 0 && (
          <Badge className="absolute left-3 top-3 bg-amber-500 hover:bg-amber-500">
            ⭐ Primary
          </Badge>
        )}
      </div>

      <div className="space-y-3 p-4">
        <div className="flex items-start gap-3">
          <HugeiconsIcon icon={Image01Icon} strokeWidth={2} />

          <div className="min-w-0 flex-1">
            <p className="truncate text-sm font-medium">{file.name}</p>

            <p className="text-xs text-muted-foreground">{size}</p>
          </div>

          <Button
            size="icon-sm"
            variant="ghost"
            className="text-destructive"
            onClick={onRemove}
          >
            <HugeiconsIcon icon={Delete02Icon} strokeWidth={2} />
          </Button>
        </div>
      </div>
    </motion.div>
  );
}
