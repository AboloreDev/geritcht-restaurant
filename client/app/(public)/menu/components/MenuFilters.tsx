import React from "react";
import { SortMenu } from "./SortMenu";
import { PriceRangeFilter } from "./PriceRangeFilter";
import { MenuSearch } from "./MenuSearch";

const MenuFilters = () => {
  return (
    <div className="mx-auto max-w-7xl px-6 gap-10 text-center flex justify-center items-center">
      <div className="flex gap-5">
        <SortMenu />
        <PriceRangeFilter />
      </div>

      <MenuSearch />
    </div>
  );
};

export default MenuFilters;
