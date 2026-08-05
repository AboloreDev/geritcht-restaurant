"use client";
import { Header } from "@/components/code/Header";
import MenuStats from "./components/MenuStats";
import { useGetMenusQuery } from "@/app/state/api/menuApi";
import { RootState, useAppDispatch, useAppSelector } from "@/app/state/redux";
import { MenuToolbar } from "./components/MenuToolbar";
import { Separator } from "@/components/ui/separator";
import MenuTable from "./components/MenuTable";
import { setPage } from "@/app/state/slices/menuSlice";
import { useState } from "react";
import CreateMenuSheet from "./components/create/CreateMenuSheet";

const Menu = () => {
  const dispatch = useAppDispatch();
  const [openMenuCreateSheet, setOpenMenuCreateSheet] = useState(false);
  const {
    categoryId,
    page,
    limit,
    query,
    sortBy,
    sortOrder,
    maxPrice,
    minPrice,
  } = useAppSelector((state: RootState) => state.menu);

  const {
    data: menu,
    isLoading,
    isFetching,
  } = useGetMenusQuery({
    category_id: categoryId,
    page,
    limit,
    query: query || undefined,
    sort_by: sortBy,
    sort_order: sortOrder,
    max_price: maxPrice,
    min_price: minPrice,
  });

  const menus = menu?.data ?? [];
  const total = menu?.meta.total ?? 0;
  const hasMore = menu ? menu.meta.page < menu.meta.total_pages : false;

  const handleCreate = () => {
    setOpenMenuCreateSheet(true);
  };

  return (
    <div className="flex flex-col p-4 space-y-4 overflow-y-auto min-h-screen">
      <Header
        title="🍽️ Menu"
        subTitle="Manage your restaurant menu, pricing, availability and images."
      />

      <MenuStats menus={menus} total={total} />

      <div className="overflow-y-auto h-[800px] ">
        <div className="bg-[#faedcd] rounded-2xl flex flex-col space-y-5 p-4">
          <MenuToolbar onCreate={handleCreate} />

          <Separator />

          <MenuTable
            menus={menus}
            isLoading={isLoading}
            isFetching={isFetching}
            hasMore={hasMore}
            onLoadMore={() => setPage(page + 1)}
          />
        </div>
      </div>

      <CreateMenuSheet
        open={openMenuCreateSheet}
        onOpenChange={() => setOpenMenuCreateSheet}
      />
    </div>
  );
};

export default Menu;
