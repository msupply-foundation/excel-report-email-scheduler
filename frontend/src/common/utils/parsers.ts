export const parseOrDefault = <T>(value: string, defaultValue: T) => {
  try {
    return JSON.parse(value) as T;
  } catch {
    return defaultValue;
  }
};
