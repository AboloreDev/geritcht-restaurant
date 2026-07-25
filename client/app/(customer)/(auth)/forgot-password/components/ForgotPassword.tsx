"use client";

import { useState } from "react";
import { Controller, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import {
  Field,
  FieldGroup,
  FieldLabel,
  FieldError,
} from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { useForgotPasswordMutation } from "@/app/state/api/authApi";
import {
  forgotPasswordSchema,
  ForgotPasswordValues,
} from "@/schema/authSchema";
import { getApiError } from "@/app/utils/apiError";
import { toast } from "sonner";
import { useRouter } from "next/navigation";

export function ForgotPasswordForm() {
  const router = useRouter();
  const [submittedEmail, setSubmittedEmail] = useState<string | null>(null);
  const [forgotPassword, { isLoading, error }] = useForgotPasswordMutation();

  const { control, handleSubmit } = useForm<ForgotPasswordValues>({
    resolver: zodResolver(
      // @ts-expect-error "<>"
      forgotPasswordSchema,
    ),
    defaultValues: { email: "" },
  });

  async function onSubmit(values: ForgotPasswordValues) {
    try {
      const reposne = await forgotPassword(values).unwrap();
      console.log(reposne);
      setSubmittedEmail(values.email);
      toast.success(
        "Reset link sent successfully, you will get an otp if use is registered!",
      );
      router.push("/verify-reset-otp");
    } catch (err) {
      console.error(err);
      toast.error(getApiError(err));
    }
  }

  if (submittedEmail) {
    return (
      <div className="text-center">
        <p className="text-white/90">
          If an account exists for{" "}
          <span className="font-medium">{submittedEmail}</span>, we&apos;ve sent
          a link to reset your password.
        </p>
        <p className="mt-4 text-sm text-white/60">
          Didn&apos;t get it? Check your spam folder, or{" "}
          <button
            onClick={() => setSubmittedEmail(null)}
            className="font-medium text-white hover:underline"
          >
            try a different email
          </button>
          .
        </p>
      </div>
    );
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      <FieldGroup>
        <Controller
          name="email"
          control={control}
          render={({ field, fieldState }) => (
            <Field data-invalid={fieldState.invalid}>
              <FieldLabel htmlFor={field.name}>Email</FieldLabel>
              <Input
                {...field}
                id={field.name}
                type="email"
                placeholder="you@example.com"
                aria-invalid={fieldState.invalid}
                autoComplete="email"
              />
              {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
            </Field>
          )}
        />

        <Button
          type="submit"
          className="w-full bg-black text-primary hover:bg-black/80"
          disabled={isLoading}
        >
          {isLoading ? "Sending…" : "Send reset link"}
        </Button>
      </FieldGroup>
    </form>
  );
}
