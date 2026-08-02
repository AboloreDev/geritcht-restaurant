export function getWeekRange(weeksAgo: number) {
  const now = new Date();
  const dayOfWeek = now.getDay();
  const startOfThisWeek = new Date(now);
  startOfThisWeek.setDate(now.getDate() - dayOfWeek);
  startOfThisWeek.setHours(0, 0, 0, 0);

  const start = new Date(startOfThisWeek);
  start.setDate(start.getDate() - weeksAgo * 7);

  const end = new Date(start);
  end.setDate(end.getDate() + 7);

  return { start, end };
}

export function todayDateString(): string {
  return new Date().toISOString().split("T")[0];
}
