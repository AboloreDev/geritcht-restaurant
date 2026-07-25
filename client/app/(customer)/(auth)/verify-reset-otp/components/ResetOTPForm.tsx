"use client";

import { useEffect, useState } from "react";
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
import { useVerifyResetOTPMutation } from "@/app/state/api/authApi";
import { OtpFormValues, otpSchema } from "@/schema/authSchema";
import { getApiError } from "@/app/utils/apiError";
import { toast } from "sonner";

const RESEND_COOLDOWN_SECONDS = 30;

export function ResetOTPForm() {
  const router = useRouter();
  const searchParams = useSearchParams();

  const [verifyOtp, { isLoading, error }] = useVerifyResetOTPMutation();
  //   const [resendOtp, { isLoading: isResending }] = useResendOtpMutation();

  const [cooldown, setCooldown] = useState(RESEND_COOLDOWN_SECONDS);
  const [resendMessage, setResendMessage] = useState<string | null>(null);

  const { control, handleSubmit, watch } = useForm<OtpFormValues>({
    resolver: zodResolver(
      // @ts-expect-error "<>"
      otpSchema,
    ),
    defaultValues: { token: "" },
  });

  const otpValue = watch("token");

  // countdown ticks every second, only while > 0
  useEffect(() => {
    if (cooldown <= 0) return;
    const timer = setInterval(() => {
      setCooldown((c) => c - 1);
    }, 1000);
    return () => clearInterval(timer);
  }, [cooldown]);

  async function onSubmit(values: OtpFormValues) {
    try {
      const response = await verifyOtp(values).unwrap();
      console.log(response);
      toast.success(response.message);
      router.push("/reset-passsword");
    } catch (err) {
      console.error(err);
      toast.error(getApiError(err));
    }
  }

  //   async function handleResend() {
  //     try {
  //       await resendOtp({ email }).unwrap();
  //       setCooldown(RESEND_COOLDOWN_SECONDS);
  //       setResendMessage("A new code has been sent.");
  //       setTimeout(() => setResendMessage(null), 4000);
  //     } catch {
  //       setResendMessage("Couldn't resend the code. Try again shortly.");
  //     }
  //   }

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
                  // native paste support: pasting a 6-digit code fills all slots
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

        <Button
          type="submit"
          className="w-full bg-black text-primary hover:bg-black/80"
          disabled={isLoading || otpValue.length !== 6}
        >
          {isLoading ? "Verifying…" : "Verify"}
        </Button>

        <div className="text-center text-sm text-white/70">
          {cooldown > 0 ? (
            <span>
              Resend code in{" "}
              <span className="font-medium text-white">{cooldown}s</span>
            </span>
          ) : (
            <button
              type="button"
              //   onClick={handleResend}
              //   disabled={isResending}
              className="font-medium text-white hover:underline disabled:opacity-50"
            >
              {/* {isResending ? "Sending…" : "Resend code"} */}
            </button>
          )}
        </div>

        {resendMessage && (
          <p className="text-center text-xs text-white/60">{resendMessage}</p>
        )}
      </FieldGroup>
    </form>
  );
}
