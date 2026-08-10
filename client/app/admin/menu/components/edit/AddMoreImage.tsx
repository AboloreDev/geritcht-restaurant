// src/app/admin/menu/components/edit/AddMoreImages.tsx
"use client";

import { useState } from "react";
import ImageDropzone from "../create/image/ImageDropzone";
import { Button } from "@/components/ui/button";
import { useUploadMenuImageMutation } from "@/app/state/api/menuApi";
import { toast } from "sonner";
import { getApiError } from "@/app/utils/apiError";

export default function AddMoreImages({
  menuId,
  currentCount,
}: {
  menuId: number;
  currentCount: number;
}) {
  // File objects can't go in Redux (not serializable) — same
  // constraint as the create flow's image step, stays local
  const [images, setImages] = useState<File[]>([]);
  const [uploadMenuImage, { isLoading }] = useUploadMenuImageMutation();

  const remaining = 4 - currentCount;

  function addImages(files: File[]) {
    setImages((prev) => [...prev, ...files].slice(0, remaining));
  }

  async function handleUpload() {
    try {
      for (const image of images) {
        await uploadMenuImage({
          id: menuId,
          image,
          is_primary: false,
        }).unwrap();
      }
      toast.success("Images added");
      setImages([]);
    } catch (error) {
      toast.error(getApiError(error));
    }
  }

  if (remaining <= 0) return null;

  return (
    <div className="space-y-3">
      <ImageDropzone images={images} onFilesSelected={addImages} />
      {images.length > 0 && (
        <Button size="sm" disabled={isLoading} onClick={handleUpload}>
          {isLoading
            ? "Uploading…"
            : `Upload ${images.length} image${images.length !== 1 ? "s" : ""}`}
        </Button>
      )}
    </div>
  );
}
