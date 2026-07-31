"use client";

import * as React from "react";

import { NavDocuments } from "@/components/nav-documents";
import { NavMain } from "@/components/nav-main";
import { NavSecondary } from "@/components/nav-secondary";
import { NavUser } from "@/components/nav-user";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import { HugeiconsIcon } from "@hugeicons/react";
import {
  Settings05Icon,
  Home11Icon,
  ShoppingCart02Icon,
  Folder03Icon,
  FileStackIcon,
  TableRoundIcon,
  AirplaneSeatIcon,
  InvoiceIcon,
  CashbackIcon,
  AutoConversationsIcon,
} from "@hugeicons/core-free-icons";
import { useAuth } from "@/app/hooks/isAuthenticated";
import Image from "next/image";

const data = {
  navMain: [
    {
      title: "Dashboard",
      url: "/admin/home",
      icon: <HugeiconsIcon icon={Home11Icon} strokeWidth={2} />,
    },
    {
      title: "Orders",
      url: "/admin/orders",
      icon: <HugeiconsIcon icon={ShoppingCart02Icon} strokeWidth={2} />,
    },
    {
      title: "Menu",
      url: "/admin/menu",
      icon: <HugeiconsIcon icon={Folder03Icon} strokeWidth={2} />,
    },
    {
      title: "Categories",
      url: "/admin/categories",
      icon: <HugeiconsIcon icon={FileStackIcon} strokeWidth={2} />,
    },
    {
      title: "Tables",
      url: "/admin/tables",
      icon: <HugeiconsIcon icon={TableRoundIcon} strokeWidth={2} />,
    },
    {
      title: "Reservations",
      url: "/admin/reservations",
      icon: <HugeiconsIcon icon={AirplaneSeatIcon} strokeWidth={2} />,
    },
  ],
  navSecondary: [
    {
      title: "Settings",
      url: "#",
      icon: <HugeiconsIcon icon={Settings05Icon} strokeWidth={2} />,
    },
  ],
  documents: [
    {
      name: "Payments",
      url: "/admin/payments",
      icon: <HugeiconsIcon icon={InvoiceIcon} strokeWidth={2} />,
    },
    {
      name: "Refunds",
      url: "/admin/refunds",
      icon: <HugeiconsIcon icon={CashbackIcon} strokeWidth={2} />,
    },
    {
      name: "Analytics",
      url: "/admin/analytics",
      icon: <HugeiconsIcon icon={AutoConversationsIcon} strokeWidth={2} />,
    },
  ],
};
export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  const { user } = useAuth();

  return (
    <Sidebar collapsible="offcanvas" {...props}>
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton
              className="data-[slot=sidebar-menu-button]:p-1.5!"
              render={<a href="/admin/home" />}
            >
              <Image
                src={"/assets/gericht.png"}
                alt="Logo Image"
                width={100}
                height={100}
              />
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent>
        <NavMain items={data.navMain} />
        <NavDocuments items={data.documents} />
        <NavSecondary items={data.navSecondary} className="mt-auto" />
      </SidebarContent>
      <SidebarFooter>
        {/* @ts-expect-error "<>" */}
        <NavUser user={user} />
      </SidebarFooter>
    </Sidebar>
  );
}
