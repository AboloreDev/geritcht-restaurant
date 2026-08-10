"use client";

import Image from "next/image";
import { motion } from "framer-motion";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Delete02Icon, Image01Icon } from "@hugeicons/core-free-icons";
import { HugeiconsIcon } from "@hugeicons/react";
import { useDeleteImageUploadMutation } from "@/app/state/api/menuApi";
import { MenuImage } from "@/app/state/types/menuTypes";
import { getApiError } from "@/app/utils/apiError";

type ExistingImageCardProps = {
  image: MenuImage;
  imageCount: number;
  onDelete: (image: MenuImage) => void;
  onPreview: (image: MenuImage) => void;
};

export default function ExistingImageCard({
  image,
  imageCount,
  onDelete,
  onPreview,
}: ExistingImageCardProps) {
  const [deleteImage, { isLoading }] = useDeleteImageUploadMutation();

  const handleDelete = async () => {
    if (imageCount === 1) return;

    try {
      const response = await deleteImage({
        id: image.id,
      }).unwrap();

      toast.success(response.message);
    } catch (error) {
      toast.error(getApiError(error));
    }
  };

  return (
    <motion.div
      layout
      initial={{
        opacity: 0,
        scale: 0.95,
      }}
      animate={{
        opacity: 1,
        scale: 1,
      }}
      exit={{
        opacity: 0,
        scale: 0.9,
      }}
      whileHover={{
        y: -3,
      }}
      transition={{
        duration: 0.2,
      }}
      className="group relative overflow-hidden rounded-2xl border bg-card"
    >
      {/* Image */}

      <div
        className="relative aspect-square cursor-pointer"
        onClick={() => onPreview(image)}
      >
        <Image
          src={image.url}
          alt={image.alt_text}
          fill
          className="object-cover transition-transform duration-300 group-hover:scale-105"
        />

        {/* Overlay */}

        <div className="absolute inset-0 flex items-end justify-between bg-gradient-to-t from-black/70 via-black/10 to-transparent p-3 opacity-0 transition-opacity duration-200 group-hover:opacity-100">
          <Button
            size="icon"
            variant="destructive"
            disabled={imageCount === 1}
            onClick={(e) => {
              e.stopPropagation();
              onDelete(image);
            }}
          >
            <HugeiconsIcon icon={Delete02Icon} strokeWidth={2} />
          </Button>
        </div>

        {/* Primary */}

        {image.is_primary && (
          <Badge className="absolute left-3 top-3 rounded-full bg-amber-500 text-white hover:bg-amber-500">
            ⭐ Primary
          </Badge>
        )}
      </div>

      {/* Footer */}

      <div className="flex items-center justify-between px-4 py-3">
        <div className="flex items-center gap-2 text-sm font-medium">
          <HugeiconsIcon icon={Image01Icon} strokeWidth={2} size={18} />

          <span className="truncate">{image.alt_text || "Menu Image"}</span>
        </div>

        <span className="text-xs text-muted-foreground">#{image.id}</span>
      </div>
    </motion.div>
  );
}
