"use client";

import { useEffect, useState } from "react";

export interface AuthUser {
  id: number;
  first_name: string;
  last_name: string;
  email: string;
}

function readUser(): AuthUser | null {
  if (typeof window === "undefined") return null;
  const raw = localStorage.getItem("user");
  if (!raw) return null;
  try {
    return JSON.parse(raw) as AuthUser;
  } catch {
    return null;
  }
}

export function useAuth() {
  const [isAuthenticated, setIsAuthenticated] = useState(() => {
    if (typeof window === "undefined") return false;
    return Boolean(localStorage.getItem("accessToken"));
  });
  const [user, setUser] = useState<AuthUser | null>(() => readUser());

  useEffect(() => {
    function handleStorageChange(e: StorageEvent) {
      if (e.key === "accessToken") {
        setIsAuthenticated(Boolean(e.newValue));
      }
      if (e.key === "user") {
        setUser(readUser());
      }
    }

    window.addEventListener("storage", handleStorageChange);
    return () => window.removeEventListener("storage", handleStorageChange);
  }, []);

  return { isAuthenticated, user };
}
