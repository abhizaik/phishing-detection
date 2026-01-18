export function isValidUrl(url: string): boolean {
  try {
    new URL(url);
    return true;
  } catch {
    return false;
  }
}

export function formatUrl(url: string): string {
  // Add protocol if missing
  if (!url.startsWith('http://') && !url.startsWith('https://')) {
    return `http://${url}`;
  }
  return url;
}

export function getDomainFromUrl(url: string): string {
  try {
    const urlObj = new URL(url);
    return urlObj.hostname;
  } catch {
    return '';
  }
}

export function formatDate(dateString: string): string {
  try {
    const date = new Date(dateString);
    return date.toLocaleDateString();
  } catch {
    return dateString;
  }
}

/**
 * Formats a URL for safe sharing by stripping schema and replacing dots with [.]
 * Example: "https://example.com" -> "example[.]com"
 */
export function formatUrlForShare(url: string): string {
  try {
    // Remove schema (http://, https://)
    let formatted = url.replace(/^https?:\/\//i, '');
    // Remove trailing slash
    formatted = formatted.replace(/\/$/, '');
    // Replace dots with [.]
    formatted = formatted.replace(/\./g, '[.]');
    return formatted;
  } catch {
    // If parsing fails, just replace dots
    return url.replace(/\./g, '[.]');
  }
}
