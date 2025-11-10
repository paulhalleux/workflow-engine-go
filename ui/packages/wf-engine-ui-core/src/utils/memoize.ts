// eslint-disable-next-line @typescript-eslint/no-explicit-any
type FunctionLike = (...args: any[]) => any;

/**
 * Memoize a function with no arguments.
 * @param fn - The function to memoize.
 * @returns An object containing the memoized function and a method to clear the cache.
 */
export const memoize = <T extends FunctionLike>(
  fn: T,
): {
  (...args: Parameters<T>): ReturnType<T>;
  clearCache: () => void;
} => {
  let cache: ReturnType<T> | null = null;
  return Object.assign(
    (...args: Parameters<T>) => {
      if (cache === null) {
        cache = fn(...args);
      }
      return cache;
    },
    {
      clearCache: () => {
        cache = null;
      },
    },
  );
};

/**
 * Memoize an asynchronous function with no arguments.
 * @param fn - The asynchronous function to memoize.
 * @returns An object containing the memoized function and a method to clear the cache.
 */
export const memoizeAsync = <T extends FunctionLike>(
  fn: T,
): {
  (...args: Parameters<T>): Promise<ReturnType<T>>;
  clearCache: () => void;
} => {
  let cache: Promise<ReturnType<T>> | null = null;
  return Object.assign(
    async (...args: Parameters<T>) => {
      if (cache === null) {
        cache = await fn(...args);
      }
      return cache;
    },
    {
      clearCache: () => {
        cache = null;
      },
    },
  );
};
