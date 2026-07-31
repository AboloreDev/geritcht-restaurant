"use client";

import { useAuth } from "@/app/hooks/isAuthenticated";
import React from "react";

const DashboardHome = () => {
  const { user } = useAuth();

  console.log(user);
  return <div></div>;
};

export default DashboardHome;
