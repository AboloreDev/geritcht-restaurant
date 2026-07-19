"use client";

import { useEffect, useRef, useState } from "react";
import { AnimatePresence, motion } from "framer-motion";
import { useMediaQuery } from "@/app/utils/useMediaQuery";
import { useDebounce } from "@/app/utils/useDebounce";
import { searchAcrossAllCategories } from "@/app/state/slices/menuSlice";
import { Search, X } from "@mynaui/icons-react";
import { Input } from "@/components/ui/input";
import { useAppDispatch } from "@/app/state/redux";
import { useSearchMenuQuery } from "@/app/state/api/menuApi";

function formatNaira(amount: number) {
  return new Intl.NumberFormat("en-NG", {
    style: "currency",
    currency: "NGN",
    maximumFractionDigits: 0,
  }).format(amount);
}

export function MenuSearch() {
  const dispatch = useAppDispatch();
  const isMobile = useMediaQuery("(max-width: 767px)");

  // Search input state
  const [value, setValue] = useState("");
  const [open, setOpen] = useState(false);

  const wrapRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  const debouncedValue = useDebounce(value, 500);
  const hasQuery = debouncedValue.trim().length > 0;

  const { data, isFetching } = useSearchMenuQuery(
    { q: debouncedValue },
    { skip: !hasQuery },
  );

  const results = hasQuery ? (data?.data ?? []) : [];

  useEffect(() => {
    if (isMobile) return;

    function onClickOutside(e: MouseEvent) {
      if (wrapRef.current && !wrapRef.current.contains(e.target as Node)) {
        setOpen(false);
      }
    }

    document.addEventListener("mousedown", onClickOutside);
    return () => document.removeEventListener("mousedown", onClickOutside);
  }, [isMobile]);

  useEffect(() => {
    if (isMobile && open) {
      document.body.style.overflow = "hidden";
      return () => {
        document.body.style.overflow = "";
      };
    }
  }, [isMobile, open]);

  useEffect(() => {
    if (open) {
      setTimeout(() => inputRef.current?.focus(), isMobile ? 250 : 100);
    }
  }, [open, isMobile]);

  function handleSelect(name: string) {
    dispatch(searchAcrossAllCategories(name));
    setValue(name);
    setOpen(false);
    document
      .getElementById("menu-grid")
      ?.scrollIntoView({ behavior: "smooth", block: "start" });
  }

  function close() {
    setOpen(false);
  }

  const resultsList = (
    <>
      {isFetching && (
        <div className="p-4 text-sm text-muted-foreground">Searching…</div>
      )}
      {!isFetching && results.length === 0 && (
        <div className="p-4 text-sm text-muted-foreground">
          No dishes found.
        </div>
      )}
      {results.map((item) => (
        <button
          key={item.id}
          onClick={() => handleSelect(item.name)}
          className="flex w-full justify-between p-4 text-left text-sm hover:bg-muted"
        >
          <span>{item.name}</span>
          <span className="text-muted-foreground">
            {formatNaira(item.price)}
          </span>
        </button>
      ))}
    </>
  );

  return (
    <div ref={wrapRef} className="relative">
      <button
        aria-label="Search"
        onClick={() => setOpen(true)}
        className="flex h-10 w-10 items-center justify-center rounded-full border"
      >
        <Search className="h-[18px] w-[18px] cursor-pointer text-white" />
      </button>

      {/* DESKTOP: inline expanding input + anchored dropdown */}
      {!isMobile && (
        <>
          {open && (
            <Input
              ref={inputRef}
              value={value}
              onChange={(e) => setValue(e.target.value)}
              placeholder="Search dishes…"
              className="absolute right-0 top-0 h-10 text-white  border-white w-64 rounded-full"
            />
          )}
          {open && hasQuery && (
            <div className="absolute right-0 top-12 z-20 w-64 h-40 overflow-y-auto text-black bg-[#fefae0]  rounded-md shadow-md">
              {resultsList}
            </div>
          )}
        </>
      )}

      {/* MOBILE: full-screen overlay */}
      {isMobile && (
        <AnimatePresence>
          {open && (
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              className="fixed inset-0 z-50 bg-[#fefae0]"
            >
              <div className="flex items-center gap-3 border-b px-4 py-3">
                <Search className="h-5 w-5 cursor-pointer shrink-0 text-white" />
                <Input
                  ref={inputRef}
                  value={value}
                  onChange={(e) => setValue(e.target.value)}
                  placeholder="Search dishes…"
                  className="flex-1 bg-transparent text-base outline-none"
                />
                <button
                  aria-label="Close search"
                  onClick={close}
                  className="shrink-0 p-1 cursor-pointer"
                >
                  <X className="h-5 w-5" />
                </button>
              </div>

              <div className="overflow-y-auto">
                {hasQuery ? (
                  resultsList
                ) : (
                  <div className="p-4 text-sm text-muted-foreground">
                    Start typing to search the menu.
                  </div>
                )}
              </div>
            </motion.div>
          )}
        </AnimatePresence>
      )}
    </div>
  );
}
