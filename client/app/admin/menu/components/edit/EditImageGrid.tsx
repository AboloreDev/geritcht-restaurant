"use client";

import { AnimatePresence, motion } from "framer-motion";
import ExistingImageCard from "./ExistingImageCard";
import DeleteImageDialog from "./DeleteImageDialog";
import { MenuImage } from "@/app/state/types/menuTypes";
import {
  openDeleteImageDialog,
  closeDeleteImageDialog,
  openLightbox,
  closeLightbox,
} from "@/app/state/slices/menuSlice";
import { useDeleteImageUploadMutation } from "@/app/state/api/menuApi";
import { RootState, useAppDispatch, useAppSelector } from "@/app/state/redux";
import { toast } from "sonner";
import { getApiError } from "@/app/utils/apiError";
import ImageLightbox from "./LightBox";

export default function EditImageGrid({
  menuId,
  images,
}: {
  menuId: number;
  images: MenuImage[];
}) {
  const dispatch = useAppDispatch();

  // dialog/lightbox state now lives in the slice — dispatched from
  // wherever a card is clicked, read back here to decide what to render
  const deleteImageId = useAppSelector(
    (s: RootState) => s.menu.dialogs.deleteImageId,
  );
  const lightbox = useAppSelector((s: RootState) => s.menu.lightbox);

  const [deleteImage, { isLoading: isDeleting }] =
    useDeleteImageUploadMutation();

  const imageForDelete = images.find((img) => img.id === deleteImageId) ?? null;
  const imageForLightbox =
    images.find((img) => img.id === lightbox.imageId) ?? null;

  async function handleConfirmDelete() {
    if (!deleteImageId) return;
    try {
      const response = await deleteImage({ id: deleteImageId }).unwrap();
      toast.success(response.message);
      dispatch(closeDeleteImageDialog());
    } catch (error) {
      toast.error(getApiError(error));
    }
  }

  return (
    <>
      <motion.div layout className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <AnimatePresence>
          {images.map((image) => (
            <ExistingImageCard
              key={image.id}
              image={image}
              imageCount={images.length}
              onDelete={(img) => dispatch(openDeleteImageDialog(img.id))}
              onPreview={(img) => dispatch(openLightbox(img.id))}
            />
          ))}
        </AnimatePresence>
      </motion.div>

      <DeleteImageDialog
        open={deleteImageId !== null}
        image={imageForDelete}
        loading={isDeleting}
        onOpenChange={(open) => !open && dispatch(closeDeleteImageDialog())}
        onConfirm={handleConfirmDelete}
      />

      <ImageLightbox
        open={lightbox.open}
        image={imageForLightbox}
        onOpenChange={(open) => !open && dispatch(closeLightbox())}
      />
    </>
  );
}
