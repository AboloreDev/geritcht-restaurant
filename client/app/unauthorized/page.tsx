"use client";

import Link from "next/link";
import { Button } from "@/components/ui/button";
import { ShieldMinus } from "@mynaui/icons-react";

export default function UnauthorizedPage() {
  return (
    <main className="flex min-h-screen items-center justify-center bg-[url('/assets/bg.png')] bg-cover bg-center bg-fixed px-6">
      <div className="w-full max-w-md rounded-3xl bg-[#fefae0] p-10 text-center shadow-lg">
        <div className="mx-auto flex h-20 w-20 items-center justify-center rounded-full bg-red-100">
          <ShieldMinus className="h-10 w-10 text-red-600" />
        </div>

        <h1 className="mt-6 font-serif text-3xl font-semibold">
          Access Denied
        </h1>

        <p className="mt-3 text-sm leading-6 text-muted-foreground">
          You don&apos;t have permission to access this page. Please sign in
          with an account that has the required permissions.
        </p>

        <div className="mt-8 flex flex-col gap-3">
          <Link href="/login">
            <Button className="w-full">Go to Login</Button>
          </Link>

          <Link href="/">
            <Button variant="outline" className="w-full">
              Back to Home
            </Button>
          </Link>
        </div>
      </div>
    </main>
  );
}
