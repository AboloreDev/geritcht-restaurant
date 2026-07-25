"use client";

import { useState } from "react";
import { Controller, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useRouter, useSearchParams } from "next/navigation";
import { Eye, EyeOff, Check, X as XIcon } from "@mynaui/icons-react";
import {
  Field,
  FieldGroup,
  FieldLabel,
  FieldError,
} from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { useResetPasswordMutation } from "@/app/state/api/authApi";
import { resetPasswordSchema, ResetPasswordValues } from "@/schema/authSchema";
import { toast } from "sonner";
import { getApiError } from "@/app/utils/apiError";

function PasswordRule({ met, label }: { met: boolean; label: string }) {
  return (
    <div
      className={`flex items-center gap-1.5 text-xs ${met ? "text-emerald-400" : "text-white/50"}`}
    >
      {met ? (
        <Check className="h-3.5 w-3.5" />
      ) : (
        <XIcon className="h-3.5 w-3.5" />
      )}
      {label}
    </div>
  );
}

export function ResetPasswordForm() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const token = searchParams.get("token") ?? "";

  const [showPassword, setShowPassword] = useState(false);
  const [showConfirm, setShowConfirm] = useState(false);
  const [isDone, setIsDone] = useState(false);

  const [resetPassword, { isLoading, error }] = useResetPasswordMutation();

  const { control, handleSubmit, watch } = useForm<ResetPasswordValues>({
    resolver: zodResolver(
      // @ts-expect-error "<>"
      resetPasswordSchema,
    ),
    defaultValues: { password: "", confirmPassword: "" },
  });

  const password = watch("password");
  const rules = {
    length: password.length >= 8,
    uppercase: /[A-Z]/.test(password),
    number: /[0-9]/.test(password),
  };

  async function onSubmit(values: ResetPasswordValues) {
    try {
      const response = await resetPassword({
        token,
        new_password: values.password,
      }).unwrap();
      setIsDone(true);
      console.log(response);
      toast.success(response.message);
      router.push("/login");
    } catch (err) {
      console.error(err);
      toast.error(getApiError(err));
    }
  }

  if (!token) {
    return (
      <p className="text-center text-sm text-black">
        This reset link is invalid or has expired. Please request a new one from
        the{" "}
        <a
          href="/forgot-password"
          className="font-medium text-white hover:underline"
        >
          forgot password
        </a>{" "}
        page.
      </p>
    );
  }

  if (isDone) {
    return (
      <div className="text-center">
        <p className="text-white/90">Your password has been reset.</p>
        <Button className="mt-6 w-full" onClick={() => router.push("/login")}>
          Continue to login
        </Button>
      </div>
    );
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      <FieldGroup>
        <Controller
          name="password"
          control={control}
          render={({ field, fieldState }) => (
            <Field data-invalid={fieldState.invalid}>
              <FieldLabel htmlFor={field.name}>New password</FieldLabel>
              <div className="relative">
                <Input
                  {...field}
                  id={field.name}
                  type={showPassword ? "text" : "password"}
                  placeholder="••••••••"
                  aria-invalid={fieldState.invalid}
                  autoComplete="new-password"
                  className="pr-10"
                />
                <button
                  type="button"
                  onClick={() => setShowPassword((v) => !v)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-white/50 hover:text-white"
                  aria-label={showPassword ? "Hide password" : "Show password"}
                >
                  {showPassword ? (
                    <EyeOff className="h-4 w-4" />
                  ) : (
                    <Eye className="h-4 w-4" />
                  )}
                </button>
              </div>

              <div className="mt-2 flex flex-wrap gap-x-4 gap-y-1">
                <PasswordRule met={rules.length} label="8+ characters" />
                <PasswordRule met={rules.uppercase} label="One uppercase" />
                <PasswordRule met={rules.number} label="One number" />
              </div>

              {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
            </Field>
          )}
        />

        <Controller
          name="confirmPassword"
          control={control}
          render={({ field, fieldState }) => (
            <Field data-invalid={fieldState.invalid}>
              <FieldLabel htmlFor={field.name}>Confirm new password</FieldLabel>
              <div className="relative">
                <Input
                  {...field}
                  id={field.name}
                  type={showConfirm ? "text" : "password"}
                  placeholder="••••••••"
                  aria-invalid={fieldState.invalid}
                  autoComplete="new-password"
                  className="pr-10"
                />
                <button
                  type="button"
                  onClick={() => setShowConfirm((v) => !v)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-white/50 hover:text-white"
                  aria-label={showConfirm ? "Hide password" : "Show password"}
                >
                  {showConfirm ? (
                    <EyeOff className="h-4 w-4" />
                  ) : (
                    <Eye className="h-4 w-4" />
                  )}
                </button>
              </div>
              {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
            </Field>
          )}
        />

        <Button
          type="submit"
          className="w-full text-primary bg-black hover:bg-black/80"
          disabled={isLoading}
        >
          {isLoading ? "Resetting…" : "Reset password"}
        </Button>
      </FieldGroup>
    </form>
  );
}
