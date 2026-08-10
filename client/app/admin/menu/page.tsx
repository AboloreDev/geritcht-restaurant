"use client";

import { Header } from "@/components/code/Header";
import { Separator } from "@/components/ui/separator";

import MenuStats from "./components/MenuStats";
import MenuTable from "./components/MenuTable";
import { MenuToolbar } from "./components/MenuToolbar";
import CreateMenuSheet from "./components/create/CreateMenuSheet";

import {
  useGetMenusAdminQuery,
  useGetMenusQuery,
} from "@/app/state/api/menuApi";

import { RootState, useAppDispatch, useAppSelector } from "@/app/state/redux";

import { openCreateSheet, setPage } from "@/app/state/slices/menuSlice";
import EditMenuSheet from "./components/edit/EditSheet";
import DeleteMenuDialog from "./components/delete/DeleteMenuDialog";

const Menu = () => {
  const dispatch = useAppDispatch();

  const {
    categoryId,
    page,
    limit,
    query,
    sortBy,
    sortOrder,
    maxPrice,
    minPrice,
    create,
    edit,
  } = useAppSelector((state: RootState) => state.menu);

  const {
    data: menu,
    isLoading,
    isFetching,
  } = useGetMenusAdminQuery({
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

  return (
    <div className="flex min-h-screen flex-col space-y-4 overflow-y-auto p-4">
      <Header
        title="🍽️ Menu"
        subTitle="Manage your restaurant menu, pricing, availability and images."
      />

      <MenuStats menus={menus} total={total} />

      <div className="h-[800px] overflow-y-auto">
        <div className="flex flex-col space-y-5 rounded-2xl bg-[#faedcd] p-4">
          <MenuToolbar onCreate={() => dispatch(openCreateSheet())} />

          <Separator />

          <MenuTable
            menus={menus}
            isLoading={isLoading}
            isFetching={isFetching}
            hasMore={hasMore}
            onLoadMore={() => dispatch(setPage(page + 1))}
          />
        </div>
      </div>

      {create.open && <CreateMenuSheet />}

      {edit.open && <EditMenuSheet />}

      <DeleteMenuDialog />
    </div>
  );
};

export default Menu;
