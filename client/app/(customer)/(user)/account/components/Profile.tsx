"use client";

import {
  useGetUserProfileQuery,
  useUpdateUserProfileMutation,
} from "@/app/state/api/userApi";
import { getApiError } from "@/app/utils/apiError";
import { Button } from "@/components/ui/button";
import {
  Field,
  FieldError,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { profileSchema, ProfileValues } from "@/schema/userSchema";
import { zodResolver } from "@hookform/resolvers/zod";
import { useEffect } from "react";
import { Controller, useForm } from "react-hook-form";
import { toast } from "sonner";

export function ProfileSection() {
  const { data, isLoading } = useGetUserProfileQuery();
  const [updateProfile, { isLoading: isSaving }] =
    useUpdateUserProfileMutation();

  const { control, handleSubmit, reset } = useForm<ProfileValues>({
    // @ts-expect-error: zod resolver types
    resolver: zodResolver(profileSchema),
    defaultValues: { first_name: "", last_name: "", phone_number: "" },
  });

  // populate the form once profile data actually arrives
  useEffect(() => {
    if (data?.data) {
      reset({
        first_name: data.data.first_name,
        last_name: data.data.last_name,
        phone_number: data.data.phone_number,
      });
    }
  }, [data, reset]);

  async function onSubmit(values: ProfileValues) {
    try {
      await updateProfile({
        first_name: values.first_name,
        last_name: values.last_name,
        phone_number: values.phone_number ?? "",
      }).unwrap();
      toast.success("Profile updated");
    } catch (err) {
      toast.error(getApiError(err));
    }
  }

  if (isLoading) {
    return <div className="h-48 animate-pulse rounded-xl bg-muted" />;
  }

  return (
    <div className="rounded-2xl border bg-[#fefae0] p-5">
      <h2 className="font-serif text-lg font-medium">Profile</h2>

      <form onSubmit={handleSubmit(onSubmit)} className="mt-4">
        <FieldGroup>
          <div className="grid grid-cols-2 gap-4">
            <Controller
              name="first_name"
              control={control}
              render={({ field, fieldState }) => (
                <Field data-invalid={fieldState.invalid}>
                  <FieldLabel htmlFor={field.name}>First name</FieldLabel>
                  <Input
                    {...field}
                    id={field.name}
                    aria-invalid={fieldState.invalid}
                  />
                  {fieldState.invalid && (
                    <FieldError errors={[fieldState.error]} />
                  )}
                </Field>
              )}
            />
            <Controller
              name="last_name"
              control={control}
              render={({ field, fieldState }) => (
                <Field data-invalid={fieldState.invalid}>
                  <FieldLabel htmlFor={field.name}>Last name</FieldLabel>
                  <Input
                    {...field}
                    id={field.name}
                    aria-invalid={fieldState.invalid}
                  />
                  {fieldState.invalid && (
                    <FieldError errors={[fieldState.error]} />
                  )}
                </Field>
              )}
            />
          </div>

          <Field>
            <FieldLabel>Email</FieldLabel>
            <Input
              value={data?.data.email}
              disabled
              className="opacity-60"
              readOnly
            />
          </Field>

          <Controller
            name="phone_number"
            control={control}
            render={({ field, fieldState }) => (
              <Field data-invalid={fieldState.invalid}>
                <FieldLabel htmlFor={field.name}>
                  Phone number{" "}
                  <span className="text-muted-foreground">(optional)</span>
                </FieldLabel>
                <Input
                  {...field}
                  id={field.name}
                  type="tel"
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
            disabled={isSaving}
            className="w-fit text-white bg-black"
          >
            {isSaving ? "Saving…" : "Save changes"}
          </Button>
        </FieldGroup>
      </form>
    </div>
  );
}
