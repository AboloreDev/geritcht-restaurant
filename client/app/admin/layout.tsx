import { AppSidebar } from "@/components/app-sidebar";
import { ProtectedRoute } from "@/components/code/ProtectedMenu";

import {
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/ui/sidebar";
import React from "react";

interface AdminDashboardProps {
  children: React.ReactNode;
}

const AdminDashboardLayout = ({ children }: AdminDashboardProps) => {
  return (
    <ProtectedRoute allowedRoles={["admin"]}>
      <SidebarProvider
        className="bg-white"
        style={
          {
            "--sidebar-width": "calc(var(--spacing) * 72)",
            "--header-height": "calc(var(--spacing) * 12)",
          } as React.CSSProperties
        }
      >
        <AppSidebar variant="inset" className="bg-primary-deep" />
        <SidebarInset className="bg-white">
          <div className="flex flex-1 flex-col">
            <div className="@container/main flex flex-1 flex-col gap-2">
              {children}
            </div>
          </div>
        </SidebarInset>
      </SidebarProvider>
    </ProtectedRoute>
  );
};

export default AdminDashboardLayout;
