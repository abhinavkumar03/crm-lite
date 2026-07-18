// apiErrorMessage extracts a human-readable message from an axios error's
// standard { success, message } response envelope, falling back to a default.
export function apiErrorMessage(err: unknown, fallback: string): string {
  if (
    typeof err === "object" &&
    err !== null &&
    "response" in err &&
    typeof (err as { response?: unknown }).response === "object"
  ) {
    const resp = (err as { response?: { data?: { message?: string } } })
      .response;
    if (resp?.data?.message) return resp.data.message;
  }
  return fallback;
}
