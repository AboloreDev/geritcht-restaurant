"use client";

import { useChangePasswordMutation } from "@/app/state/api/userApi";
import { getApiError } from "@/app/utils/apiError";
import { Button } from "@/components/ui/button";
import {
  Field,
  FieldError,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { passwordSchema, PasswordValues } from "@/schema/userSchema";
import { zodResolver } from "@hookform/resolvers/zod";
import { Controller, useForm } from "react-hook-form";
import { toast } from "sonner";

export function PasswordSection() {
  const [changePassword, { isLoading }] = useChangePasswordMutation();

  const { control, handleSubmit, reset } = useForm<PasswordValues>({
    // @ts-expect-error: zod resolver types
    resolver: zodResolver(passwordSchema),
    defaultValues: {
      current_password: "",
      new_password: "",
      confirm_password: "",
    },
  });

  async function onSubmit(values: PasswordValues) {
    try {
      await changePassword({
        current_password: values.current_password,
        new_password: values.new_password,
      }).unwrap();
      toast.success("Password changed");
      reset();
    } catch (err) {
      toast.error(getApiError(err));
    }
  }

  return (
    <div className="mt-6 rounded-2xl border bg-[#fefae0] p-5">
      <h2 className="font-serif text-lg font-medium">Security</h2>

      <form onSubmit={handleSubmit(onSubmit)} className="mt-4">
        <FieldGroup>
          <Controller
            name="current_password"
            control={control}
            render={({ field, fieldState }) => (
              <Field data-invalid={fieldState.invalid}>
                <FieldLabel htmlFor={field.name}>Current password</FieldLabel>
                <Input
                  {...field}
                  id={field.name}
                  type="password"
                  aria-invalid={fieldState.invalid}
                />
                {fieldState.invalid && (
                  <FieldError errors={[fieldState.error]} />
                )}
              </Field>
            )}
          />
          <Controller
            name="new_password"
            control={control}
            render={({ field, fieldState }) => (
              <Field data-invalid={fieldState.invalid}>
                <FieldLabel htmlFor={field.name}>New password</FieldLabel>
                <Input
                  {...field}
                  id={field.name}
                  type="password"
                  aria-invalid={fieldState.invalid}
                />
                {fieldState.invalid && (
                  <FieldError errors={[fieldState.error]} />
                )}
              </Field>
            )}
          />
          <Controller
            name="confirm_password"
            control={control}
            render={({ field, fieldState }) => (
              <Field data-invalid={fieldState.invalid}>
                <FieldLabel htmlFor={field.name}>
                  Confirm new password
                </FieldLabel>
                <Input
                  {...field}
                  id={field.name}
                  type="password"
                  aria-invalid={fieldState.invalid}
                />
                {fieldState.invalid && (
                  <FieldError errors={[fieldState.error]} />
                )}
              </Field>
            )}
          />

          <Button
            type="submit"
            disabled={isLoading}
            className="w-fit bg-black text-white"
          >
            {isLoading ? "Changing…" : "Change password"}
          </Button>
        </FieldGroup>
      </form>
    </div>
  );
}
