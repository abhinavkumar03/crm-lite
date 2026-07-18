/**
 * Lightweight in-memory TTL cache for rarely-changing metadata
 * (modules, fields, validation schemas). Avoids refetching on every
 * navigation between Forms / Tables / Import / Export.
 */

const DEFAULT_TTL_MS = 60_000;

type Entry<T> = {
  value: T;
  expires: number;
};

const store = new Map<string, Entry<unknown>>();

export async function cached<T>(
  key: string,
  loader: () => Promise<T>,
  ttlMs: number = DEFAULT_TTL_MS
): Promise<T> {
  const hit = store.get(key);
  if (hit && hit.expires > Date.now()) {
    return hit.value as T;
  }
  const value = await loader();
  store.set(key, { value, expires: Date.now() + ttlMs });
  return value;
}

/** Drop all metadata entries (call after Settings mutations). */
export function invalidateMetadataCache(prefix?: string): void {
  if (!prefix) {
    store.clear();
    return;
  }
  for (const key of store.keys()) {
    if (key.startsWith(prefix)) {
      store.delete(key);
    }
  }
}
