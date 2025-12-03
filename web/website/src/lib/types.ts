export interface ApiResponse<T> {
  data?: T;
  error?: string;
}

export interface AnalyzeResult {
  url: string;
  domain: string;
  result: {
    risk_score: number;
    trust_score: number;
    final_score: number;
    verdict: string;
    reasons?: {
      neutral_reasons?: string[];
      good_reasons?: string[];
      bad_reasons?: string[];
    };
  };
  analysis?: any;
  features?: any;
  infrastructure?: any;
  domain_info?: any;
  performance?: any;
  incomplete?: boolean;
  errors?: any;
}

export interface ScreenshotResponse {
  status: string;
  msg: string;
  file: string;
}

export interface RankResponse {
  rank: number;
}

export interface IpCheckResponse {
  uses_ip: boolean;
}

export interface IpResolveResponse {
  ip_addresses: string[];
}

export interface LengthResponse {
  too_long: boolean;
}

export interface DepthResponse {
  too_deep: boolean;
}

export interface HstsResponse {
  has_hsts: boolean;
}

export interface RedirectsResponse {
  redirect_count: number;
}

export interface PunycodeResponse {
  uses_punycode: boolean;
}

export interface TrustedTldResponse {
  is_trusted_tld: boolean;
  is_icann: boolean;
}

export interface RiskyTldResponse {
  is_risky_tld: boolean;
  is_icann: boolean;
}

export interface UrlShortenerResponse {
  is_url_shortener: boolean;
}

export interface StatusCodeResponse {
  status_code: number;
}

export interface WhoisResponse {
  domain: string;
  age: number;
  created_at: string;
  expires_at: string;
  registrar: string;
  raw_data: any;
}

export interface PhishingCheckResult {
  rank?: RankResponse;
  ipCheck?: IpCheckResponse;
  ipResolve?: IpResolveResponse;
  length?: LengthResponse;
  depth?: DepthResponse;
  hsts?: HstsResponse;
  redirects?: RedirectsResponse;
  punycode?: PunycodeResponse;
  trustedTld?: TrustedTldResponse;
  riskyTld?: RiskyTldResponse;
  urlShortener?: UrlShortenerResponse;
  statusCode?: StatusCodeResponse;
  whois?: WhoisResponse;
  [key: string]: any; // Allow string indexing
}

export interface CheckStatus {
  loading: boolean;
  error?: string;
  completed: boolean;
}
