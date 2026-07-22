"use client";

import { useRouter } from "next/navigation";
import Link from "next/link";
import {
  Field,
  FieldGroup,
  FieldLabel,
  FieldError,
} from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Controller, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { LoginFormData, loginSchema } from "@/schema/authSchema";
import { useAppDispatch, useAppSelector } from "@/app/state/redux";
import { EyeOff } from "@hugeicons/core-free-icons";
import { Eye, EyeOffSolid } from "@mynaui/icons-react";
import { setShowPassword } from "@/app/state/slices/globalSlice";

export function LoginForm() {
  const router = useRouter();
  const showPassword = useAppSelector((state) => state.global.showPassword);
  const dispatch = useAppDispatch();
  //   const [login, { isLoading, error }] = useLoginMutation();

  const { control, handleSubmit } = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  });

  async function onSubmit(data: LoginFormData) {
    // try {
    //   await login(data).unwrap();
    //   router.push("/menu");
    // } catch {
    //   // error state is already surfaced via the `error` variable below
    // }
  }

  const handlePasswordToggle = () => {
    dispatch(setShowPassword());
  };

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

        <Controller
          name="password"
          control={control}
          render={({ field, fieldState }) => (
            <Field data-invalid={fieldState.invalid}>
              <div className="flex items-center justify-between">
                <FieldLabel htmlFor={field.name}>Password</FieldLabel>
                <Link
                  href="/forgot-password"
                  className="text-xs hover:underline"
                >
                  Forgot password?
                </Link>
              </div>

              <div className="relative">
                <Input
                  {...field}
                  id={field.name}
                  type={showPassword ? "text" : "password"}
                  placeholder="••••••••"
                  aria-invalid={fieldState.invalid}
                  autoComplete="current-password"
                  className="pr-10"
                />

                <button
                  type="button"
                  onClick={handlePasswordToggle}
                  className="absolute right-3 top-1/2 cursor-pointer -translate-y-1/2 text-muted-foreground hover:text-foreground"
                  aria-label={showPassword ? "Hide password" : "Show password"}
                >
                  {showPassword ? <EyeOffSolid size={18} /> : <Eye size={18} />}
                </button>
              </div>

              {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
            </Field>
          )}
        />
        {/* {error && (
          <p className="text-sm text-red-400">
            {"data" in error &&
            typeof error.data === "object" &&
            error.data &&
            "message" in error.data
              ? String((error.data as { message: string }).message)
              : "Couldn't log you in. Check your details and try again."}
          </p>
        )}

        <Button type="submit" className="w-full" disabled={isLoading}>
          {isLoading ? "Logging in…" : "Log in"}
        </Button> */}
      </FieldGroup>
    </form>
  );
}
