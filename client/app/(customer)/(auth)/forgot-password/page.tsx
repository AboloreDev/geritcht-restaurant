import React from "react";
import AuthLayout from "../layout";
import { ForgotPasswordForm } from "./components/ForgotPassword";
import Link from "next/link";
import { ArrowLeft } from "@mynaui/icons-react";

const ForgotPassword = () => {
  return (
    <AuthLayout
      imagePanel={
        <div className="relative h-full w-full bg-[url('/assets/auth-3.jpg')] bg-cover bg-center">
          <div className=" inset-0 bg-gradient-to-t from-black/60 via-black/10 to-transparent" />
          <div className="absolute bottom-8 left-8 right-8 text-white">
            <p className="font-serif text-2xl font-medium">Geritcht</p>
            <p className="mt-1 text-sm text-white/80"></p>
          </div>
        </div>
      }
    >
      <div className="px-6">
        <Link
          href="/login"
          className="inline-flex items-center gap-1.5 text-sm text-black/60 hover:text-black"
        >
          <ArrowLeft className="h-4 w-4" />
          Back to login
        </Link>

        <Link
          href="/"
          className="mt-6 block text-4xl text-center text-black font-medium"
        >
          Geritcht
        </Link>

        <h1 className="mt-6 font-serif text-2xl font-semibold md:text-3xl">
          Forgot Password
        </h1>
        <p className="mt-2 text-sm text-black/60">
          No worries — enter the email you signed up with and we&apos;ll send
          you a link to reset your password.
        </p>

        <div className="mt-5">
          <ForgotPasswordForm />
        </div>
      </div>
    </AuthLayout>
  );
};

export default ForgotPassword;
