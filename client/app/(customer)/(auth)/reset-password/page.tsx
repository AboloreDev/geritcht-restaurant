import { Suspense } from "react";
import Link from "next/link";
import { ArrowLeft } from "@mynaui/icons-react";
import { ResetPasswordForm } from "./components/ResetPasswordForm";
import AuthLayout from "../layout";

export default function ResetPasswordPage() {
  return (
    <AuthLayout
      imagePanel={
        <div className="relative h-full w-full bg-[url('/assets/auth-5.jpg')] bg-cover bg-center">
          <div className="absolute inset-0 bg-gradient-to-t from-black/60 via-black/10 to-transparent" />

          <div className="absolute bottom-8 left-8 right-8 text-white">
            <p className="font-serif text-2xl font-medium">Geritcht</p>
            <p className="mt-1 text-sm text-white/80">
              A secure password keeps your account and reservations protected.
            </p>
          </div>
        </div>
      }
    >
      <div className="flex h-full flex-col p-6 md:p-8">
        <Link
          href="/verify-reset-otp"
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

          <h1 className="mt-8 font-serif text-3xl font-semibold text-black">
            Set a New Password
          </h1>

          <p className="mt-3 max-w-md text-sm leading-6 text-black/60">
            Create a strong password to secure your account. Make sure it is at
            least 8 characters long and includes a mix of letters, numbers, and
            symbols.
          </p>
        </div>

        <div className="mt-8">
          <Suspense fallback={null}>
            <ResetPasswordForm />
          </Suspense>
        </div>
      </div>
    </AuthLayout>
  );
}
