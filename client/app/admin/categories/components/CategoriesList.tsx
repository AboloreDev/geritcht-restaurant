"use client";

import { MenuCategory } from "@/app/state/types/menuTypes";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { HugeiconsIcon } from "@hugeicons/react";
import { Folder } from "@mynaui/icons-react";
import {
  Delete02Icon,
  Edit02Icon,
  MoreHorizontalIcon,
} from "@hugeicons/core-free-icons";
import { getCategoryColor } from "@/app/utils/categoryColor";

interface Categories {
  category: MenuCategory;
  onEdit: (category: MenuCategory) => void;
  onDelete: (category: MenuCategory) => void;
}

const CategoriesList = ({ category, onEdit, onDelete }: Categories) => {
  return (
    <article
      key={category.id}
      id={`category-${category.id}`}
      tabIndex={-1}
      className="grid gap-4 px-4 py-4 md:grid-cols-[minmax(220px,1fr)_120px_80px] md:items-center md:px-5"
    >
      {/* Category */}
      <div className="flex min-w-0 items-center gap-3">
        <div
          className={`grid h-10 w-10 shrink-0 place-items-center rounded-2xl ${getCategoryColor(category.id)}`}
        >
          <Folder strokeWidth={2} />
        </div>

        <div className="min-w-0">
          <p className="truncate font-medium">{category.name}</p>
          <p className="mt-0.5 truncate text-xs text-muted-foreground">
            {category.description}
          </p>
        </div>
      </div>

      {/* Status */}
      <div>
        <Badge
          variant={category.is_active ? "default" : "outline"}
          className={
            category.is_active
              ? "bg-emerald-600 text-white"
              : "bg-red-600 text-white"
          }
        >
          {category.is_active ? "Active" : "In-Active"}
        </Badge>
      </div>

      {/* Actions */}
      <div className="flex justify-end">
        <DropdownMenu>
          <DropdownMenuTrigger
            render={
              <Button
                variant="ghost"
                size="icon-sm"
                aria-label={`Actions for ${category.name}`}
              />
            }
          >
            <HugeiconsIcon icon={MoreHorizontalIcon} strokeWidth={2} />
          </DropdownMenuTrigger>

          <DropdownMenuContent align="end" className="min-w-44 bg-white">
            <DropdownMenuItem
              onClick={() => onEdit(category)}
              className="cursor-pointer"
            >
              <HugeiconsIcon icon={Edit02Icon} strokeWidth={2} />
              Edit category
            </DropdownMenuItem>

            <DropdownMenuSeparator />

            <DropdownMenuItem
              variant="destructive"
              className="cursor-pointer"
              onClick={() => onDelete(category)}
            >
              <HugeiconsIcon icon={Delete02Icon} strokeWidth={2} />
              Delete
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </article>
  );
};

export default CategoriesList;
