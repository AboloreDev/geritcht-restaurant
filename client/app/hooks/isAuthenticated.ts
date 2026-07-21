"use client";

import { useEffect, useState } from "react";

export function useAuth() {
  const [isAuthenticated, setIsAuthenticated] = useState(() => {
    if (typeof window === "undefined") return false;
    return Boolean(localStorage.getItem("accessToken"));
  });

  useEffect(() => {
    function handleStorageChange(e: StorageEvent) {
      if (e.key === "accessToken") {
        setIsAuthenticated(Boolean(e.newValue));
      }
    }
    window.addEventListener("storage", handleStorageChange);
    return () => window.removeEventListener("storage", handleStorageChange);
  }, []);

  return { isAuthenticated };
}
