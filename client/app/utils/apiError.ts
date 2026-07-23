export function getApiError(error: any): string {
  if (!error) return "Something went wrong.";

  if (error.data?.message) return error.data.message;
  if (error.data?.error) return error.data.error;
  if (error.message) return error.message;

  return "Something went wrong.";
}
