"use client";

import { useState } from "react";
import ImageDropzone from "./ImageDropzone";
import { AnimatePresence, motion } from "framer-motion";
import ImagePreviewCard from "./ImagePreviewCard";

interface Props {
  menuId: number;
}

export default function MenuImageUpload({ menuId }: Props) {
  const [images, setImages] = useState<File[]>([]);

  const addImages = (files: File[]) => {
    setImages((prev) => {
      const merged = [...prev, ...files];

      // only keep the first four
      return merged.slice(0, 4);
    });
  };

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-xl font-semibold">Upload Menu Images</h2>

        <p className="mt-1 text-sm text-muted-foreground">
          Upload up to 4 images. The first image becomes the primary image.
        </p>
      </div>

      <ImageDropzone images={images} onFilesSelected={addImages} />

      {/* preview cards go here in step 2 */}
      <AnimatePresence>
        {images.length > 0 && (
          <motion.div
            layout
            className="grid gap-5 sm:grid-cols-2 lg:grid-cols-4"
          >
            {images.map((image, index) => (
              <ImagePreviewCard
                key={`${image.name}-${index}`}
                file={image}
                index={index}
                onRemove={() =>
                  setImages((prev) => prev.filter((_, i) => i !== index))
                }
              />
            ))}
          </motion.div>
        )}
      </AnimatePresence>

      <div className="text-sm text-muted-foreground">
        {images.length}/4 selected
      </div>
    </div>
  );
}
