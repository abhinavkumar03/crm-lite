import {
  Condition,
  FormValues,
  VisibilityRule,
} from "../types";

export function isEmptyValue(value: unknown): boolean {
  if (value === null || value === undefined) return true;
  if (typeof value === "string") return value.trim() === "";
  if (Array.isArray(value)) return value.length === 0;
  return false;
}

export function evaluateCondition(
  condition: Condition,
  values: FormValues
): boolean {
  const actual = values[condition.field];

  switch (condition.operator) {
    case "equals":
      return String(actual ?? "") === String(condition.value ?? "");
    case "not_equals":
      return String(actual ?? "") !== String(condition.value ?? "");
    case "in":
      return (
        Array.isArray(condition.value) &&
        condition.value.map(String).includes(String(actual ?? ""))
      );
    case "not_in":
      return !(
        Array.isArray(condition.value) &&
        condition.value.map(String).includes(String(actual ?? ""))
      );
    case "empty":
      return isEmptyValue(actual);
    case "not_empty":
      return !isEmptyValue(actual);
    case "truthy":
      return Boolean(actual);
    case "falsy":
      return !actual;
    default:
      return true;
  }
}

// computeHiddenFields returns the set of field api_names that should be hidden
// given the current values. A "show" rule keeps its targets hidden until the
// condition is met; a "hide" rule hides its targets while the condition is met.
export function computeHiddenFields(
  rules: VisibilityRule[],
  values: FormValues
): Set<string> {
  const hidden = new Set<string>();

  for (const rule of rules) {
    const met = evaluateCondition(rule.when, values);

    for (const target of rule.targets) {
      if (rule.effect === "show" && !met) hidden.add(target);
      if (rule.effect === "hide" && met) hidden.add(target);
    }
  }

  return hidden;
}
