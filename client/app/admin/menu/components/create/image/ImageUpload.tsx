"use client";

import { useState } from "react";
import ImageDropzone from "./ImageDropzone";
import { AnimatePresence, motion } from "framer-motion";
import ImagePreviewCard from "./ImagePreviewCard";
import { useUploadMenuImageMutation } from "@/app/state/api/menuApi";
import { toast } from "sonner";
import { getApiError } from "@/app/utils/apiError";
import { Button } from "@/components/ui/button";
import { RootState, useAppDispatch, useAppSelector } from "@/app/state/redux";
import { setUploadProgress } from "@/app/state/slices/menuSlice";

interface Props {
  onFinish: () => void;
}

export default function MenuImageUpload({ onFinish }: Props) {
  const dispatch = useAppDispatch();

  // menuId and uploadProgress come from the slice — genuinely shared
  // UI state that other parts of the create flow also care about.
  // images itself CANNOT live in Redux: File objects aren't
  // serializable, and RTK's middleware will warn/error if raw File
  // instances are put into the store. This stays local component state.
  const { menuId, uploadProgress } = useAppSelector(
    (state: RootState) => state.menu.create,
  );

  const [uploadMenuImage, { isLoading }] = useUploadMenuImageMutation();
  const [images, setImages] = useState<File[]>([]);

  const addImages = (files: File[]) => {
    setImages((prev) => {
      const merged = [...prev, ...files];
      return merged.slice(0, 4);
    });
  };

  const uploadImages = async () => {
    if (!images.length) {
      toast.error("Please upload at least one image.");
      return;
    }

    // guards against menuId being null if this step somehow renders
    // before step 1 finishes creating the menu item
    if (!menuId) {
      toast.error("Something went wrong — no menu item to attach images to.");
      return;
    }

    try {
      for (const [index, image] of images.entries()) {
        dispatch(setUploadProgress(index + 1));
        await uploadMenuImage({
          id: menuId,
          image,
          is_primary: index === 0,
        }).unwrap();
      }
      toast.success("Menu created successfully.");
      setImages([]);
      dispatch(setUploadProgress(0));
      onFinish();
    } catch (error) {
      toast.error(getApiError(error));
    }
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

      <div className="mt-8 flex justify-end gap-3 pt-6">
        <Button
          variant="outline"
          onClick={() => setImages([])}
          disabled={isLoading}
          className="border-none bg-red-500 text-white hover:bg-red-600 hover:text-white"
        >
          Clear
        </Button>

        <Button
          disabled={images.length === 0 || isLoading}
          onClick={uploadImages}
        >
          {isLoading
            ? `Uploading ${uploadProgress}/${images.length}`
            : "Finish"}
        </Button>
      </div>
    </div>
  );
}
