"use client";

import { useGetMenusQuery } from "@/app/state/api/menuApi";
import { RootState, useAppSelector } from "@/app/state/redux";
import React from "react";

const Menu = () => {
  const filters = useAppSelector((state: RootState) => state.menu);

  const { data, isLoading } = useGetMenusQuery({
    page: filters.page,
    limit: filters.limit,
    category_id: filters.categoryId,
    query: filters.query,
    min_price: filters.minPrice,
    max_price: filters.maxPrice,
    sort_by: filters.sortBy,
    order: filters.order,
  });

  console.log(data);
  return (
    <section className="bg-[url('/assets/bg.png')] bg-cover bg-center bg-no-repeat text-primary">
      <div className="max-w-6xl mx-auto px-3 py-8">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
          {data?.data.map((item) => (
            <div
              key={item.id}
              className="bg-white/30 backdrop-blur-sm rounded-xl p-4 border border-white/50"
            >
              <div className="h-48 overflow-hidden rounded-lg mb-4">
                <img
                  src="/assets/food.jpg"
                  alt={item.name}
                  className="w-full h-full object-cover"
                />
              </div>
              <h3 className="text-xl font-bold mb-2">{item.name}</h3>
              <p className="mb-4">Some description about the food item.</p>
              <div className="flex justify-between items-center">
                <span className="font-bold text-lg">Ksh 250</span>
                <button className="bg-secondary hover:bg-secondary/90 text-white px-4 py-2 rounded-lg transition duration-300">
                  Add to Cart
                </button>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
};

export default Menu;
