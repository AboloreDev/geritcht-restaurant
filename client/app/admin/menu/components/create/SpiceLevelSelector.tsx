"use client";

import { Controller, useFormContext } from "react-hook-form";
import { motion } from "framer-motion";
import { Flame } from "@mynaui/icons-react";

const labels = ["Not Spicy", "Mild", "Medium", "Hot", "Extra Hot", "Extreme"];

export default function SpiceLevelSelector() {
  const { control } = useFormContext();

  return (
    <Controller
      name="spice_level"
      control={control}
      render={({ field }) => (
        <div>
          <div className="mb-3 flex items-center justify-between">
            <p className="font-medium">{labels[field.value]}</p>

            <span className="text-sm text-muted-foreground">
              {field.value}/5
            </span>
          </div>

          <div className="flex gap-3">
            {Array.from({ length: 5 }).map((_, index) => {
              const level = index + 1;

              const active = level <= field.value;

              return (
                <motion.button
                  key={level}
                  type="button"
                  whileHover={{
                    scale: 1.15,
                  }}
                  whileTap={{
                    scale: 0.9,
                  }}
                  onClick={() => field.onChange(level)}
                  className="rounded-xl p-2 cursor-pointer transition-colors"
                >
                  <Flame
                    size={30}
                    className={
                      active ? "fill-red-500 text-red-500" : "text-gray-300"
                    }
                  />
                </motion.button>
              );
            })}
          </div>

          <button
            type="button"
            onClick={() => field.onChange(0)}
            className="mt-4 text-sm text-red-500 cursor-pointer hover:underline"
          >
            Clear spice level
          </button>
        </div>
      )}
    />
  );
}
