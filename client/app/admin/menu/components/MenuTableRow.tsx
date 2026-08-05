"use client";

import Image from "next/image";
import { useRouter } from "next/navigation";
import { motion } from "framer-motion";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { TableCell, TableRow } from "@/components/ui/table";
import { Delete02Icon, Edit02Icon } from "@hugeicons/core-free-icons";
import { HugeiconsIcon } from "@hugeicons/react";

import { Menu } from "@/app/state/types/menuTypes";
import { formatNaira } from "@/app/utils/formatNaira";
import { resolveImageSrc } from "@/app/utils/resolveImage";

interface Props {
  menu: Menu;
  onEdit?: (menu: Menu) => void;
  onDelete?: (menu: Menu) => void;
}

export default function MenuTableRow({ menu, onEdit, onDelete }: Props) {
  const router = useRouter();

  const image = resolveImageSrc(menu) ?? "/placeholder-food.png";

  return (
    <TableRow
      className="group h-24 cursor-pointer transition-colors hover:bg-[#fefae0]"
      onClick={() => router.push(`/admin/menu/${menu.id}`)}
    >
      {/* Image */}

      <TableCell className="align-middle">
        <motion.div
          whileHover={{ scale: 1.05 }}
          transition={{ duration: 0.2 }}
          className="relative size-16 overflow-hidden rounded-xl bg-muted"
        >
          <Image
            src={image}
            alt={menu.name}
            fill
            sizes="64px"
            className="object-cover"
          />
        </motion.div>
      </TableCell>

      {/* Name */}

      <TableCell className="align-middle">
        <div className="max-w-md space-y-1">
          <p className="truncate font-semibold">{menu.name}</p>

          <p className="line-clamp-2 text-sm text-muted-foreground">
            {menu.description || "No description"}
          </p>
        </div>
      </TableCell>

      {/* Category */}

      <TableCell className="align-middle">
        <Badge variant="outline" className="rounded-full">
          {menu.category.name}
        </Badge>
      </TableCell>

      {/* Price */}

      <TableCell className="text-right align-middle font-semibold whitespace-nowrap">
        {formatNaira(menu.price)}
      </TableCell>

      {/* Status */}

      <TableCell className="text-center align-middle">
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

      <TableCell className="align-middle">
        <div className="flex justify-end gap-1">
          <Button
            size="icon"
            variant="ghost"
            onClick={(e) => {
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
              e.stopPropagation();
              onDelete?.(menu);
            }}
          >
            <HugeiconsIcon icon={Delete02Icon} strokeWidth={2} />
          </Button>
        </div>
      </TableCell>
    </TableRow>
  );
}
