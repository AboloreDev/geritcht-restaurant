"use client";

import { useDeactivateAccountMutation } from "@/app/state/api/userApi";
import { getApiError } from "@/app/utils/apiError";
import { Button } from "@/components/ui/button";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { toast } from "sonner";

export function DangerZoneSection() {
  const router = useRouter();
  const [confirming, setConfirming] = useState(false);
  const [deactivate, { isLoading }] = useDeactivateAccountMutation();

  async function handleDeactivate() {
    try {
      await deactivate().unwrap();
      toast.success("Account deactivated");
      localStorage.removeItem("accessToken");
      localStorage.removeItem("user");
      router.push("/");
    } catch (err) {
      toast.error(getApiError(err));
    }
  }

  return (
    <div className="mt-6 rounded-2xl bg-[#fefae0] p-5">
      <h2 className="font-serif text-lg font-bold ">Danger zone</h2>
      <p className="mt-1 text-md text-muted-foreground">
        Deactivating your account will sign you out and disable access until
        reactivated.
      </p>

      {!confirming ? (
        <Button
          variant="outline"
          className="mt-4 text-white bg-red-600 hover:bg-destructive/10"
          onClick={() => setConfirming(true)}
        >
          Deactivate account
        </Button>
      ) : (
        <div className="mt-4 flex items-center gap-3">
          <Button variant="outline" onClick={() => setConfirming(false)}>
            Cancel
          </Button>
          <Button
            variant="destructive"
            disabled={isLoading}
            onClick={handleDeactivate}
            className="bg-black text-white"
          >
            {isLoading ? "Deactivating…" : "Yes, deactivate my account"}
          </Button>
        </div>
      )}
    </div>
  );
}
