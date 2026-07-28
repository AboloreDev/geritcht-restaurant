import { ProtectedRoute } from "@/components/code/ProtectedMenu";
import { OrderContents } from "./components/OrderContents";

const UserOrderPage = () => {
  return (
    <ProtectedRoute>
      <OrderContents />
    </ProtectedRoute>
  );
};

export default UserOrderPage;
