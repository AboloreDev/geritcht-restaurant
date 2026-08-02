"use client";

import OrderSearch from "./components/OrderSearch";
import OrderFilters from "./components/OrderFilters";
import { Header } from "@/components/code/Header";
import { useGetAllOrdersQuery } from "@/app/state/api/orderApi";
import { RootState, useAppDispatch, useAppSelector } from "@/app/state/redux";
import OrdersList from "./components/OrdersList";
import { useEffect, useState } from "react";
import { Order } from "@/app/state/types/orderTypes";
import { appendOrders, setOrders } from "@/app/state/slices/orderSlice";

const OrderPage = () => {
  const dispatch = useAppDispatch();
  const { page, pageSize, orders, filterDate, filterStatus, filterType } =
    useAppSelector((state: RootState) => state.order);

  const { data, isLoading, isFetching } = useGetAllOrdersQuery({
    page,
    page_size: pageSize,
    date: filterDate || undefined,
    status: filterStatus || undefined,
    type: filterType || undefined,
  });

  useEffect(() => {
    if (!data) return;

    if (page === 1) {
      dispatch(setOrders(data.data.orders));
    } else {
      dispatch(appendOrders(data.data.orders));
    }
  }, [data, page, dispatch]);

  const hasMore = data ? data.data.page < data.data.total_pages : false;

  return (
    <div className="p-4 flex flex-col space-y-4 h-screen overflow-y-auto">
      <Header title="Orders" subTitle="View and Manage all orders" />
      <div className="bg-[#faedcd] flex flex-col md:flex-row items-center rounded-2xl">
        <OrderFilters isLoading={isLoading} isFetching={isFetching} />
        <OrderSearch />
      </div>

      <OrdersList
        orders={orders}
        isLoading={isLoading}
        isFetching={isFetching}
        hasMore={hasMore}
        page={page ?? 1}
      />
    </div>
  );
};

export default OrderPage;
