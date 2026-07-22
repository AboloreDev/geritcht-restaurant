export const TIME_SLOTS = [
  "18:00:00",
  "20:00:00",
  "21:00:00",
  "22:00:00",
  "23:00:00",
  "02:00:00",
] as const;

// "18:00:00" -> "6:00 PM"
export function formatTimeSlot(slot: string): string {
  const [hourStr, minuteStr] = slot.split(":");
  const hour = parseInt(hourStr, 10);
  const period = hour >= 12 ? "PM" : "AM";
  const displayHour = hour % 12 === 0 ? 12 : hour % 12;
  return `${displayHour}:${minuteStr} ${period}`;
}
