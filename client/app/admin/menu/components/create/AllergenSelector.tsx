"use client";

import { Controller, useFormContext } from "react-hook-form";
import { motion } from "framer-motion";
import { Check } from "@mynaui/icons-react";
import { Spinner } from "@/components/ui/spinner";
import { useGetAllergensQuery } from "@/app/state/api/allergenApi";

export default function AllergenSelector() {
  const { control } = useFormContext();
  const { data, isLoading } = useGetAllergensQuery();

  const allergens = data?.data ?? [];

  if (isLoading) {
    <div className="animate-spin w-10">
      <Spinner />
    </div>;
  }

  return (
    <Controller
      control={control}
      name="allergen_ids"
      render={({ field }) => {
        const selected = field.value ?? [];

        const toggle = (id: number) => {
          if (selected.includes(id)) {
            field.onChange(selected.filter((item: number) => item !== id));
            return;
          }

          field.onChange([...selected, id]);
        };

        return (
          <div className="mt-3 grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
            {allergens.map((item) => {
              const active = selected.includes(item.id);

              return (
                <motion.button
                  whileTap={{ scale: 0.96 }}
                  whileHover={{ scale: 1.02 }}
                  key={item.id}
                  type="button"
                  onClick={() => toggle(item.id)}
                  className={`flex items-center justify-between rounded-2xl p-2 text-left transition-all cursor-pointer ${
                    active
                      ? "border-black bg-black text-white"
                      : "border-border bg-white hover:border-black"
                  }`}
                >
                  <span className="font-medium">{item.name}</span>

                  <motion.div
                    initial={false}
                    animate={{
                      scale: active ? 1 : 0,
                      rotate: active ? 0 : -90,
                    }}
                    transition={{
                      duration: 0.15,
                    }}
                  >
                    <Check size={18} />
                  </motion.div>
                </motion.button>
              );
            })}
          </div>
        );
      }}
    />
  );
}
