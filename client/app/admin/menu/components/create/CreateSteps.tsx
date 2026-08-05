type StepProps = {
  children: React.ReactNode;
  active?: boolean;
  completed?: boolean;
};

export function CreateStep({ children, active, completed }: StepProps) {
  return (
    <div
      className={[
        "flex h-10 w-10 items-center justify-center rounded-full text-sm font-semibold transition-all",
        completed && "bg-green-600 text-white",
        active && !completed && "bg-black text-white",
        !active && !completed && "bg-muted text-muted-foreground",
      ]
        .filter(Boolean)
        .join(" ")}
    >
      {children}
    </div>
  );
}
