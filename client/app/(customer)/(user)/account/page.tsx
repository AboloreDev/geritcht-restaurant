import { ProtectedRoute } from "@/components/code/ProtectedMenu";
import React from "react";
import { SettingsContent } from "./components/SettingsContent";

const UserSettings = () => {
  return (
    <ProtectedRoute>
      <SettingsContent />
    </ProtectedRoute>
  );
};

export default UserSettings;
