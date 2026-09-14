import { ProtectedRoute } from "@/components/code/ProtectedMenu";
import PaymentHistoryContent from "./components/PaymentHistoryContent";

export default function Page() {
  return (
    <ProtectedRoute>
      <PaymentHistoryContent />
    </ProtectedRoute>
  );
}
