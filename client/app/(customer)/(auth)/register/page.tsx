"use client";

import Link from "next/link";
import AuthLayout from "../layout";
import { RegisterForm } from "./components/RegisterForm";

export default function RegisterPage() {
  return (
    <AuthLayout
      imagePanel={
        <div className="relative h-full w-full bg-[url('/assets/auth-2.jpg')] bg-cover bg-center">
          <div className=" inset-0 bg-gradient-to-t from-black/60 via-black/10 to-transparent" />
          <div className="absolute bottom-8 left-8 right-8 text-white">
            <p className="font-serif text-2xl font-medium">Geritcht</p>
            <p className="mt-1 text-sm text-white/80">
              Where exceptional cuisine meets effortless dining.
            </p>
          </div>
        </div>
      }
    >
      <div className="px-6 py-4">
        <h1 className="mt-2 font-serif text-2xl font-semibold md:text-3xl">
          Create your account
        </h1>

        <div className="mt-4">
          <RegisterForm />
        </div>

        <p className="mt-2 text-center text-sm">
          Already have an account?{" "}
          <Link href="/login" className="font-medium hover:underline">
            Sign In
          </Link>
        </p>
      </div>
    </AuthLayout>
  );
}
