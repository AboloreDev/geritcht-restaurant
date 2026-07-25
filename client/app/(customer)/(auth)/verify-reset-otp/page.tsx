import Link from "next/link";
import React from "react";
import AuthLayout from "../layout";
import { ResetOTPForm } from "./components/ResetOTPForm";
import { ArrowLeft } from "@mynaui/icons-react";

const VerifyResetOTP = () => {
  return (
    <AuthLayout
      imagePanel={
        <div className="relative h-full w-full bg-[url('/assets/auth-4.jpg')] bg-cover bg-center">
          <div className="absolute inset-0 bg-gradient-to-t from-black/60 via-black/10 to-transparent" />

          <div className="absolute bottom-8 left-8 right-8 text-white">
            <p className="font-serif text-2xl font-medium">Geritcht</p>
            <p className="mt-1 text-sm text-white/80">
              You&apos;re just one step away from securing your account.
            </p>
          </div>
        </div>
      }
    >
      <div className="flex h-full flex-col p-6 md:p-8">
        <Link
          href="/forgot-password"
          className="inline-flex items-center gap-1.5 text-sm text-black/60 transition-colors hover:text-black"
        >
          <ArrowLeft className="h-4 w-4" />
          Back
        </Link>

        <div className="mt-8">
          <Link
            href="/"
            className="block text-center text-4xl font-medium text-black"
          >
            Geritcht
          </Link>

          <h1 className="mt-8 font-serif text-3xl font-semibold">
            Verify Reset Code
          </h1>

          <p className="mt-3 max-w-md text-sm leading-6 text-black/60">
            We&apos;ve sent a verification code to your email address. Enter the
            code below to verify your identity before creating a new password.
          </p>
        </div>

        <div className="mt-8">
          <ResetOTPForm />
        </div>
      </div>
    </AuthLayout>
  );
};

export default VerifyResetOTP;
