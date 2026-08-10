import { Menu } from "../state/types/menuTypes";

export function resolveImageSrc(menu: Menu): string | null {
  if (!menu.images || menu.images.length === 0) {
    return null;
  }

  const primary = menu.images.find((img) => img.url) ?? menu.images[0];

  return primary.url;
}
