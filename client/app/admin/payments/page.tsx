import PaymentFilters from "./components/PaymentFilters";
import PaymentsList from "./components/PaymentList";

export default function PaymentsContent() {
  return (
    <div className="p-6">
      <h1 className="font-serif text-2xl font-semibold">Payments</h1>
      <p className="mt-1 text-sm text-muted-foreground">
        View and track all payment transactions.
      </p>

      <div className="mt-6">
        <PaymentFilters />
      </div>

      <PaymentsList />
    </div>
  );
}
