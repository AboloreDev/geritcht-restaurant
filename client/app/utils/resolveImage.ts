import { Menu } from "../state/types/menuTypes";

export function resolveImageSrc(menu: Menu): string | null {
  if (menu.image_url) return menu.image_url;
  if (!menu.images?.length) return null;
  const primary = menu.images.find((img) => img.is_primary);
  return (primary ?? menu.images[0]).alt_text || null;
}
