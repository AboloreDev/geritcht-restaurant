"use client";

import Reservations from "./components/Reservations";
import { Suspense } from "react";

const ReservationPage = () => {
  return (
    <Suspense fallback={<div>Loading…</div>}>
      <Reservations />
    </Suspense>
  );
};

export default ReservationPage;
