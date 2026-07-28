import { ProtectedRoute } from "@/components/code/ProtectedMenu";
import React from "react";
import { MyReservationsContent } from "./components/MyReservationContents";

const MyReservations = () => {
  return (
    <ProtectedRoute>
      <MyReservationsContent />
    </ProtectedRoute>
  );
};

export default MyReservations;
