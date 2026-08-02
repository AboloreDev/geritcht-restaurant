"use client";

import { useAuth } from "@/app/hooks/isAuthenticated";
import { StatsCards } from "./components/StatsCard";
import { Header } from "@/components/code/Header";
import { PendingOrdersList } from "./components/PendingOrders";
import { TodaysReservationsList } from "./components/Reservations";
import { RecentUsersList } from "./components/Users";
import { ActivityFeed } from "./components/ActivityFeed";
import { InventoryAlerts } from "./components/inventoryAlerts";

const DashboardHome = () => {
  const { user } = useAuth();

  return (
    <div className="p-4 flex flex-col space-y-6 h-screen overflow-y-auto">
      <Header title="Dashboard" subTitle="Here's what's happening today" />
      <StatsCards />

      <InventoryAlerts />

      <div className="grid grid-cols-1 gap-4 lg:grid-cols-3">
        <PendingOrdersList />
        <TodaysReservationsList />
        <RecentUsersList />
      </div>

      <div className="grid grid-cols-1 gap-4 lg:grid-cols-2">
        <ActivityFeed />
        {/* <WaitlistPanel /> */}
      </div>
    </div>
  );
};

export default DashboardHome;
