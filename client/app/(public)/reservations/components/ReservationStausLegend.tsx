import { TABLE_STATUS_CONFIG } from "@/app/utils/tableStatusConfig";

export function StatusLegend() {
  return (
    <div className="mt-10 flex flex-wrap items-center justify-center gap-6 border-t pt-6 text-sm text-primary-deep">
      {Object.values(TABLE_STATUS_CONFIG).map((status) => (
        <div key={status.label} className="flex items-center gap-2">
          <span className={`h-5 w-5 rounded-full ${status.dot}`} />
          {status.label}
        </div>
      ))}
    </div>
  );
}
