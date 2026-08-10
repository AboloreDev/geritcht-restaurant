"use client";

import { useAppSelector, RootState } from "@/app/state/redux";
import { useGetAllRservationsQuery } from "@/app/state/api/reservationsApi";

import ReservationFilters from "./ReservationFilters";
import ReservationsList from "./ReservationList";
import ReservationSearch from "./ReservationSearch";
import { Header } from "@/components/code/Header";

const ReservationContent = () => {
  const { page, pageSize, filterDate, filterStatus, filterTimeSlot } =
    useAppSelector((state: RootState) => state.reservation);

  const { data, isLoading, isFetching } = useGetAllRservationsQuery({
    page,
    page_size: pageSize,
    date: filterDate,
    status: filterStatus,
    time_slot: filterTimeSlot,
  });

  const reservations = data?.data.reservations ?? [];
  const hasMore = data ? page < data.total_pages : false;

  return (
    <div className="">
      <Header
        title="Reservations"
        subTitle="View and manage table reservatioons"
      />

      <div className="mt-3 flex flex-col">
        <ReservationFilters isLoading={isLoading} isFetching={isFetching} />
        <div className="flex items-center gap-3 px-6">
          <p>Search Reservations:</p>
          <ReservationSearch />
        </div>
      </div>

      <div id="reservation-list">
        <ReservationsList
          reservations={reservations}
          isLoading={isLoading}
          isFetching={isFetching}
          hasMore={hasMore}
          page={page ?? 1}
        />
      </div>
    </div>
  );
};

export default ReservationContent;
