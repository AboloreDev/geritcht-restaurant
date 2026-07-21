import { Menu } from "../state/types/menuTypes";

export function getVisibleBadges(menu: Menu, max = 2) {
  const tags = [
    ...(menu.dietary_tags ?? []).map((t) => ({
      id: `d-${t.id}`,
      name: t.name,
      kind: "diet" as const,
    })),
    ...(menu.allergens ?? []).map((a) => ({
      id: `a-${a.id}`,
      name: a.name,
      kind: "allergen" as const,
    })),
  ];
  return {
    visible: tags.slice(0, max),
    overflowCount: Math.max(tags.length - max, 0),
  };
}
