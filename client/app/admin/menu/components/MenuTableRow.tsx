"use client";

import Image from "next/image";
import { useRouter } from "next/navigation";
import { motion } from "framer-motion";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Switch } from "@/components/ui/switch";
import { TableCell, TableRow } from "@/components/ui/table";
import { Delete02Icon, Edit02Icon } from "@hugeicons/core-free-icons";
import { HugeiconsIcon } from "@hugeicons/react";
import { Menu } from "@/app/state/types/menuTypes";
import { formatNaira } from "@/app/utils/formatNaira";
import { resolveImageSrc } from "@/app/utils/resolveImage";
import { useAppDispatch } from "@/app/state/redux";
import {
  openEditSheet,
  openDeleteMenuDialog,
} from "@/app/state/slices/menuSlice";
import { useToggleMenuAvailabilityMutation } from "@/app/state/api/menuApi";
import { toast } from "sonner";
import { getApiError } from "@/app/utils/apiError";

interface Props {
  menu: Menu;
}

export default function MenuTableRow({ menu }: Props) {
  const dispatch = useAppDispatch();
  const router = useRouter();
  const [toggleAvailability, { isLoading: isToggling }] =
    useToggleMenuAvailabilityMutation();

  const image = resolveImageSrc(menu) ?? "/placeholder-food.png";

  function handleEdit(e: React.MouseEvent) {
    e.preventDefault();
    e.stopPropagation();
    dispatch(openEditSheet(menu.id));
  }

  function handleDelete(e: React.MouseEvent) {
    e.preventDefault();
    e.stopPropagation();
    dispatch(openDeleteMenuDialog(menu.id));
  }

  async function handleToggle() {
    try {
      const response = await toggleAvailability({
        id: menu.id,
        body: { is_available: !menu.is_available },
      }).unwrap();
      toast.success(response.message);
    } catch (error) {
      toast.error(getApiError(error));
    }
  }

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

      {/* Menu name + truncated description, click-through via row's own onClick */}
      <TableCell className="align-middle">
        <div className="max-w-md space-y-1">
          <p className="truncate font-semibold">{menu.name}</p>
          <p className="line-clamp-1 text-sm text-muted-foreground">
            {menu.description || "No description"}
          </p>
        </div>
      </TableCell>

      {/* Category */}
      <TableCell className="align-middle">
        <Badge variant="outline" className="rounded-full bg-white">
          {menu.category.name}
        </Badge>
      </TableCell>

      {/* Price */}
      <TableCell className="text-right align-middle font-semibold whitespace-nowrap">
        {formatNaira(menu.price)}
      </TableCell>

      {/* Availability — single source of truth, replaces the old redundant Status badge */}
      <TableCell
        className="text-center align-middle"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="flex flex-col items-center gap-1">
          <Switch
            checked={menu.is_available}
            disabled={isToggling}
            onCheckedChange={handleToggle}
            className="data-[state=checked]:bg-green-600! data-[state=unchecked]:bg-gray-300! [&>span]:bg-white! [&>span]:shadow-md! transition-colors"
          />
          <span
            className={`text-xs font-medium ${menu.is_available ? "text-green-600" : "text-red-500"}`}
          >
            {menu.is_available ? "Available" : "Unavailable"}
          </span>
        </div>
      </TableCell>

      {/* Actions */}
      <TableCell className="align-middle">
        <div className="flex justify-end gap-1">
          <Button size="icon" variant="ghost" onClick={handleEdit}>
            <HugeiconsIcon icon={Edit02Icon} strokeWidth={2} />
          </Button>

          <Button
            size="icon"
            variant="ghost"
            className="text-destructive hover:text-destructive"
            onClick={handleDelete}
          >
            <HugeiconsIcon icon={Delete02Icon} strokeWidth={2} />
          </Button>
        </div>
      </TableCell>
    </TableRow>
  );
}
