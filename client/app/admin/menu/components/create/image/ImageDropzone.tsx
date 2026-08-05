"use client";

import { ImagePlus } from "@hugeicons/core-free-icons";
import { HugeiconsIcon } from "@hugeicons/react";

interface Props {
  images: File[];
  onFilesSelected: (files: File[]) => void;
}

export default function ImageDropzone({ images, onFilesSelected }: Props) {
  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = Array.from(e.target.files ?? []);

    if (!files.length) return;

    onFilesSelected(files);

    e.target.value = "";
  };

  return (
    <label className="flex cursor-pointer flex-col items-center justify-center rounded-3xl border-2 border-dashed border-muted-foreground/20 bg-muted/30 px-6 py-12 transition hover:border-primary hover:bg-primary/5">
      <HugeiconsIcon
        icon={ImagePlus}
        className="mb-4 h-12 w-12 text-muted-foreground"
      />

      <h3 className="font-semibold">Upload menu images</h3>

      <p className="mt-2 text-center text-sm text-muted-foreground">
        Drag & drop or click to browse
      </p>

      <p className="mt-1 text-xs text-muted-foreground">PNG • JPG • WEBP</p>

      <input
        hidden
        multiple
        accept="image/*"
        type="file"
        onChange={handleChange}
      />
    </label>
  );
}
