import type {
  RankResponse,
  IpCheckResponse,
  IpResolveResponse,
  LengthResponse,
  DepthResponse,
  HstsResponse,
  RedirectsResponse,
  PunycodeResponse,
  TrustedTldResponse,
  RiskyTldResponse,
  UrlShortenerResponse,
  StatusCodeResponse,
  WhoisResponse,
  ApiResponse,
} from './types';
import { env } from '$env/dynamic/public';

const PUBLIC_BASE_URL = env.PUBLIC_BASE_URL || 'http://localhost:8080/api/v1';

async function makeRequest<T>(endpoint: string, url: string): Promise<ApiResponse<T>> {
  try {
    // Ensure the URL is properly formatted before encoding
    const formattedUrl = url.trim();
    const encodedUrl = encodeURIComponent(formattedUrl);
    const fullUrl = `${PUBLIC_BASE_URL}${endpoint}?url=${encodedUrl}`;

    console.log(`Making request to: ${fullUrl}`);

    const response = await fetch(fullUrl);

    console.log(`Response status: ${response.status}`);

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ error: `HTTP ${response.status}` }));
      console.error(`API Error:`, errorData);
      return { error: errorData.error || `HTTP ${response.status}` };
    }

    const data = await response.json();
    console.log(`API Success for ${endpoint}:`, data);
    return { data };
  } catch (error) {
    console.error(`API Error for ${endpoint}:`, error);
    return { error: error instanceof Error ? error.message : 'Network error' };
  }
}

export const api = {
  async analyze(url: string): Promise<ApiResponse<any>> {
    return makeRequest<any>('/analyze', url);
  },
  async screenshot(url: string): Promise<ApiResponse<any>> {
    return makeRequest<any>('/screenshot', url);
  },
  async getRank(url: string): Promise<ApiResponse<RankResponse>> {
    return makeRequest<RankResponse>('/rank', url);
  },

  async getIpCheck(url: string): Promise<ApiResponse<IpCheckResponse>> {
    return makeRequest<IpCheckResponse>('/ip/check', url);
  },

  async getIpResolve(url: string): Promise<ApiResponse<IpResolveResponse>> {
    return makeRequest<IpResolveResponse>('/ip/resolve', url);
  },

  async getLength(url: string): Promise<ApiResponse<LengthResponse>> {
    return makeRequest<LengthResponse>('/length', url);
  },

  async getDepth(url: string): Promise<ApiResponse<DepthResponse>> {
    return makeRequest<DepthResponse>('/depth', url);
  },

  async getHsts(url: string): Promise<ApiResponse<HstsResponse>> {
    return makeRequest<HstsResponse>('/hsts', url);
  },

  async getRedirects(url: string): Promise<ApiResponse<RedirectsResponse>> {
    return makeRequest<RedirectsResponse>('/redirects', url);
  },

  async getPunycode(url: string): Promise<ApiResponse<PunycodeResponse>> {
    return makeRequest<PunycodeResponse>('/punycode', url);
  },

  async getTrustedTld(url: string): Promise<ApiResponse<TrustedTldResponse>> {
    return makeRequest<TrustedTldResponse>('/trusted-tld', url);
  },

  async getRiskyTld(url: string): Promise<ApiResponse<RiskyTldResponse>> {
    return makeRequest<RiskyTldResponse>('/risky-tld', url);
  },

  async getUrlShortener(url: string): Promise<ApiResponse<UrlShortenerResponse>> {
    return makeRequest<UrlShortenerResponse>('/url-shortener', url);
  },

  async getStatusCode(url: string): Promise<ApiResponse<StatusCodeResponse>> {
    return makeRequest<StatusCodeResponse>('/status-code', url);
  },

  async getWhois(url: string): Promise<ApiResponse<WhoisResponse>> {
    return makeRequest<WhoisResponse>('/whois', url);
  },
};
