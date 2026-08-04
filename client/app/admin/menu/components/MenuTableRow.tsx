"use client";

import Link from "next/link";
import { motion } from "framer-motion";
import Image from "next/image";
import { TableCell, TableRow } from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Edit02Icon, Delete02Icon } from "@hugeicons/core-free-icons";
import { HugeiconsIcon } from "@hugeicons/react";
import { formatNaira } from "@/app/utils/formatNaira";
import { Menu } from "@/app/state/types/menuTypes";
import { resolveImageSrc } from "@/app/utils/resolveImage";

interface Props {
  menu: Menu;
  onEdit?: (menu: Menu) => void;
  onDelete?: (menu: Menu) => void;
}

export default function MenuTableRow({ menu, onEdit, onDelete }: Props) {
  const imageSrc = resolveImageSrc(menu);
  return (
    <TableRow
      // @ts-expect-error "shadcn issues"
      asChild
      className="cursor-pointer transition-colors hover:bg-[#fefae0]"
    >
      <Link href={`/admin/menu/${menu.id}`}>
        <TableCell>
          <motion.div
            whileHover={{
              scale: 1.05,
            }}
            transition={{
              duration: 0.2,
            }}
            className="relative h-14 w-14 overflow-hidden rounded-xl border"
          >
            <Image
              src={imageSrc ?? "/"}
              alt={menu.name}
              fill
              className="object-cover"
            />
          </motion.div>
        </TableCell>

        {/* Name */}

        <TableCell className="max-w-sm">
          <p className="font-semibold">{menu.name}</p>

          <p className="mt-1 line-clamp-2 text-sm text-muted-foreground">
            {menu.description || "No description"}
          </p>
        </TableCell>

        {/* Category */}

        <TableCell>
          <Badge variant="secondary" className="rounded-full">
            {menu.category.name}
          </Badge>
        </TableCell>

        {/* Price */}

        <TableCell className="text-right font-semibold">
          {formatNaira(menu.price)}
        </TableCell>

        {/* Status */}

        <TableCell>
          <Badge
            className={
              menu.is_available
                ? "bg-green-600 hover:bg-green-600"
                : "bg-red-500 hover:bg-red-500"
            }
          >
            {menu.is_available ? "Available" : "Unavailable"}
          </Badge>
        </TableCell>

        {/* Actions */}

        <TableCell>
          <div className="flex justify-end gap-2">
            <Button
              size="icon"
              variant="ghost"
              onClick={(e) => {
                e.preventDefault();
                e.stopPropagation();
                onEdit?.(menu);
              }}
            >
              <HugeiconsIcon icon={Edit02Icon} strokeWidth={2} />
            </Button>

            <Button
              size="icon"
              variant="ghost"
              className="text-destructive hover:text-destructive"
              onClick={(e) => {
                e.preventDefault();
                e.stopPropagation();
                onDelete?.(menu);
              }}
            >
              <HugeiconsIcon icon={Delete02Icon} strokeWidth={2} />
            </Button>
          </div>
        </TableCell>
      </Link>
    </TableRow>
  );
}
