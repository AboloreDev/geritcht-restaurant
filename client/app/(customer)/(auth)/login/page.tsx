"use client";

import Link from "next/link";
import { LoginForm } from "./components/LoginForm";
import AuthLayout from "../layout";

export default function LoginPage() {
  return (
    <AuthLayout
      imagePanel={
        <div className="relative h-full w-full bg-[url('/assets/auth-1.jpg')] bg-cover bg-center">
          <div className=" inset-0 bg-gradient-to-t from-black/60 via-black/10 to-transparent" />
          <div className="absolute bottom-8 left-8 right-8 text-white">
            <p className="font-serif text-2xl font-medium">Geritcht</p>
            <p className="mt-1 text-sm text-white/80">
              Reserve your table, order ahead, come hungry.
            </p>
          </div>
        </div>
      }
    >
      <div className="p-6">
        <Link href="/" className="text-4xl text-center text-black font-medium">
          Geritcht
        </Link>

        <h1 className="mt-6 font-serif text-2xl font-semibold md:text-3xl">
          Welcome back
        </h1>
        <p className="mt-2 text-sm">
          Log in to manage your orders and reservations.
        </p>

        <div className="mt-8">
          <LoginForm />
        </div>

        <p className="mt-6 text-center text-sm">
          Don&apos;t have an account?{" "}
          <Link href="/register" className="font-medium hover:underline">
            Register
          </Link>
        </p>
      </div>
    </AuthLayout>
  );
}
