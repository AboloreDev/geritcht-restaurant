"use client";

import * as React from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { HugeiconsIcon } from "@hugeicons/react";
import { ChevronDownIcon } from "@hugeicons/core-free-icons";
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/components/ui/collapsible";
import {
  SidebarGroup,
  SidebarGroupContent,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSub,
  SidebarMenuSubButton,
  SidebarMenuSubItem,
} from "@/components/ui/sidebar";

type NavItem = {
  title: string;
  url?: string;
  icon?: React.ReactNode;
  items?: { title: string; url: string }[];
};

const MENU_MANAGEMENT_STATE_KEY = "admin-menu-management-open";

function isRouteActive(pathname: string, url: string) {
  return pathname === url || pathname.startsWith(`${url}/`);
}

function subscribeToMenuManagementState(onStoreChange: () => void) {
  window.addEventListener("storage", onStoreChange);
  return () => window.removeEventListener("storage", onStoreChange);
}

function getMenuManagementState() {
  return window.localStorage.getItem(MENU_MANAGEMENT_STATE_KEY) !== "false";
}

export function NavMain({ items }: { items: NavItem[] }) {
  const pathname = usePathname();

  return (
    <SidebarGroup>
      <SidebarGroupContent>
        <SidebarMenu className="space-y-1">
          {items.map((item) =>
            item.items ? (
              <MenuManagementItem
                key={item.title}
                item={item}
                pathname={pathname}
              />
            ) : (
              <SidebarMenuItem key={item.title}>
                <SidebarMenuButton
                  isActive={
                    item.url ? isRouteActive(pathname, item.url) : false
                  }
                  tooltip={item.title}
                  render={<Link href={item.url ?? "#"} />}
                  className="h-11 rounded-xl px-3 transition-colors duration-200 hover:bg-primary hover:text-primary-foreground data-active:bg-primary data-active:text-primary-foreground"
                >
                  {item.icon && <span className="shrink-0">{item.icon}</span>}
                  <span>{item.title}</span>
                </SidebarMenuButton>
              </SidebarMenuItem>
            ),
          )}
        </SidebarMenu>
      </SidebarGroupContent>
    </SidebarGroup>
  );
}

function MenuManagementItem({
  item,
  pathname,
}: {
  item: NavItem;
  pathname: string;
}) {
  const hasActiveChild =
    item.items?.some((child) => isRouteActive(pathname, child.url)) ?? false;
  const savedOpen = React.useSyncExternalStore(
    subscribeToMenuManagementState,
    getMenuManagementState,
    () => true,
  );
  const [openOverride, setOpenOverride] = React.useState<boolean | null>(null);
  const open = hasActiveChild || (openOverride ?? savedOpen);

  const handleOpenChange = (nextOpen: boolean) => {
    setOpenOverride(nextOpen);
    window.localStorage.setItem(MENU_MANAGEMENT_STATE_KEY, String(nextOpen));
  };

  return (
    <Collapsible
      open={open}
      onOpenChange={handleOpenChange}
      className="group/collapsible"
    >
      <SidebarMenuItem>
        <SidebarMenuButton
          render={<CollapsibleTrigger />}
          isActive={hasActiveChild}
          tooltip={item.title}
          className="h-11 rounded-xl px-3 transition-colors duration-200 hover:bg-primary hover:text-primary-foreground data-active:bg-primary data-active:text-primary-foreground"
        >
          {item.icon && <span className="shrink-0">{item.icon}</span>}
          <span>{item.title}</span>
          <HugeiconsIcon
            icon={ChevronDownIcon}
            strokeWidth={2}
            className="ml-auto transition-transform duration-200 group-data-[state=open]/collapsible:rotate-180"
          />
        </SidebarMenuButton>
        <CollapsibleContent>
          <SidebarMenuSub className="mt-1">
            {item.items?.map((child) => (
              <SidebarMenuSubItem key={child.title}>
                <SidebarMenuSubButton
                  isActive={isRouteActive(pathname, child.url)}
                  render={<Link href={child.url} />}
                  className="h-9 rounded-lg px-3 data-active:bg-primary/15 data-active:font-medium data-active:text-primary"
                >
                  <span>{child.title}</span>
                </SidebarMenuSubButton>
              </SidebarMenuSubItem>
            ))}
          </SidebarMenuSub>
        </CollapsibleContent>
      </SidebarMenuItem>
    </Collapsible>
  );
}
