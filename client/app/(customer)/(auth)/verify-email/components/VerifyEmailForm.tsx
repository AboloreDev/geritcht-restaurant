"use client";

import { Controller, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useRouter, useSearchParams } from "next/navigation";
import {
  InputOTP,
  InputOTPGroup,
  InputOTPSlot,
  InputOTPSeparator,
} from "@/components/ui/input-otp";
import { Field, FieldGroup, FieldError } from "@/components/ui/field";
import { Button } from "@/components/ui/button";
import { useVerifyEmailMutation } from "@/app/state/api/authApi";
import { OtpFormValues, otpSchema } from "@/schema/authSchema";
import { toast } from "sonner";
import { getApiError } from "@/app/utils/apiError";

export function VerifyEmailForm() {
  const router = useRouter();

  const [verifyEmail, { isLoading, error }] = useVerifyEmailMutation();

  const { control, handleSubmit, watch } = useForm<OtpFormValues>({
    resolver: zodResolver(
      // @ts-expect-error "<>"
      otpSchema,
    ),
    defaultValues: { token: "" },
  });

  const otpValue = watch("token");

  async function onSubmit(values: OtpFormValues) {
    try {
      const response = await verifyEmail(values).unwrap();
      console.log(response);
      toast.success(response.message);
      router.push("/login");
    } catch (err) {
      console.error(err);
      toast.error(getApiError(err));
    }
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      <FieldGroup>
        <Controller
          name="token"
          control={control}
          render={({ field, fieldState }) => (
            <Field data-invalid={fieldState.invalid}>
              <div className="flex justify-center">
                <InputOTP
                  maxLength={6}
                  value={field.value}
                  onChange={field.onChange}
                  pattern="^[0-9]+$"
                >
                  <InputOTPGroup>
                    <InputOTPSlot index={0} />
                    <InputOTPSlot index={1} />
                    <InputOTPSlot index={2} />
                  </InputOTPGroup>
                  <InputOTPSeparator />
                  <InputOTPGroup>
                    <InputOTPSlot index={3} />
                    <InputOTPSlot index={4} />
                    <InputOTPSlot index={5} />
                  </InputOTPGroup>
                </InputOTP>
              </div>
              {fieldState.invalid && (
                <div className="mt-2 flex justify-center">
                  <FieldError errors={[fieldState.error]} />
                </div>
              )}
            </Field>
          )}
        />

        {error && (
          <p className="text-center text-sm text-red-400">
            {"data" in error &&
            typeof error.data === "object" &&
            error.data &&
            "message" in error.data
              ? String((error.data as { message: string }).message)
              : "That code didn't work. Please try again."}
          </p>
        )}

        <Button
          type="submit"
          className="w-full bg-black text-primary hover:bg-black/80"
          disabled={isLoading || otpValue.length !== 6}
        >
          {isLoading ? "Verifying…" : "Verify email"}
        </Button>
      </FieldGroup>
    </form>
  );
}
