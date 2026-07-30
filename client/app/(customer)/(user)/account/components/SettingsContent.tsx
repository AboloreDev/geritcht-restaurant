"use client";

import { ArrowLeft } from "@mynaui/icons-react";
import Link from "next/link";
import { ProfileSection } from "./Profile";
import { PasswordSection } from "./PasswordSection";
import { DangerZoneSection } from "./DangerZoneSection";

export function SettingsContent() {
  return (
    <div className="min-h-screen bg-[url('/assets/bg.png')] bg-cover bg-center bg-fixed">
      <div className="mx-auto max-w-2xl px-6 py-10">
        <h1 className="font-serif text-3xl text-primary font-semibold">
          Account Settings
        </h1>
        <p className="mt-2 text-sm text-primary-deep">
          Manage your profile, security, and account.
        </p>

        <Link
          href="/"
          className="mt-4 inline-flex items-center gap-1.5 text-sm text-primary hover:text-foreground"
        >
          <ArrowLeft className="h-4 w-4" />
          Back
        </Link>

        <div className="mt-6">
          <ProfileSection />
          <PasswordSection />
          <DangerZoneSection />
        </div>
      </div>
    </div>
  );
}
