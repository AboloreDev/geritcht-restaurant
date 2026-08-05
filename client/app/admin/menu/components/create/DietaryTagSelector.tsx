"use client";

import { Controller, useFormContext } from "react-hook-form";
import { motion } from "framer-motion";
import { Check } from "@mynaui/icons-react";
import { useGetDietaryTagsQuery } from "@/app/state/api/dietaryTagsApi";
import { Spinner } from "@/components/ui/spinner";

export default function DietaryTagSelector() {
  const { control } = useFormContext();
  const { data, isLoading } = useGetDietaryTagsQuery();

  const dietaryTags = data?.data ?? [];

  if (isLoading) {
    <div className="animate-spin w-10">
      <Spinner />
    </div>;
  }

  return (
    <Controller
      control={control}
      name="dietary_tag_ids"
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
            {dietaryTags.map((item) => {
              const active = selected.includes(item.id);

              return (
                <motion.button
                  key={item.id}
                  type="button"
                  whileHover={{ scale: 1.02 }}
                  whileTap={{ scale: 0.96 }}
                  onClick={() => toggle(item.id)}
                  className={`flex items-center justify-between rounded-2xl p-2 cursor-pointer text-left transition-all ${
                    active
                      ? "border-emerald-600 bg-emerald-600 text-white"
                      : "border-border bg-white hover:border-emerald-500"
                  }`}
                >
                  <span className="font-medium">{item.name}</span>

                  <motion.div
                    initial={false}
                    animate={{
                      scale: active ? 1 : 0,
                      rotate: active ? 0 : -90,
                    }}
                    transition={{ duration: 0.15 }}
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
