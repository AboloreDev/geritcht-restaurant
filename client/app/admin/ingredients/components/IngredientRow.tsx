"use client";

import { motion } from "framer-motion";
import { Badge } from "@/components/ui/badge";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { HugeiconsIcon } from "@hugeicons/react";
import {
  MoreVerticalIcon,
  Edit02Icon,
  Delete02Icon,
} from "@hugeicons/core-free-icons";
import { Ingredient } from "@/app/state/types/ingredientTypes";
import { useAppDispatch } from "@/app/state/redux";
import {
  openEditIngredientSheet,
  openDeleteIngredientDialog,
} from "@/app/state/slices/ingredientSlice";
import StockLevelBar from "./StockLevelBar";

export default function IngredientRow({
  ingredient,
}: {
  ingredient: Ingredient;
}) {
  const dispatch = useAppDispatch();

  return (
    <motion.div
      layout
      initial={{ opacity: 0, y: 8 }}
      animate={{ opacity: 1, y: 0 }}
      exit={{ opacity: 0 }}
      transition={{ duration: 0.2 }}
      className="flex items-center gap-4 rounded-2xl bg-[#faedcd] p-4"
    >
      <div className="min-w-0 flex-1">
        <div className="flex items-center gap-2">
          <p className="font-serif font-medium">{ingredient.name}</p>
          {ingredient.is_low && (
            <Badge className="bg-red-500 text-white hover:bg-red-500">
              Low stock
            </Badge>
          )}
        </div>
      </div>

      <div className="w-48 shrink-0">
        <StockLevelBar
          current={ingredient.current_stock}
          threshold={ingredient.min_threshold}
          unit={ingredient.unit}
        />
      </div>

      <DropdownMenu>
        <DropdownMenuTrigger className="shrink-0 cursor-pointer p-1">
          <HugeiconsIcon icon={MoreVerticalIcon} strokeWidth={2} />
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end" className="bg-white">
          <DropdownMenuItem
            onClick={() => dispatch(openEditIngredientSheet(ingredient.id))}
            className="cursor-pointer"
          >
            <HugeiconsIcon icon={Edit02Icon} strokeWidth={2} />
            Edit
          </DropdownMenuItem>
          <DropdownMenuItem
            className="cursor-pointer text-red-600"
            onClick={() => dispatch(openDeleteIngredientDialog(ingredient.id))}
          >
            <HugeiconsIcon icon={Delete02Icon} strokeWidth={2} />
            Delete
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </motion.div>
  );
}
