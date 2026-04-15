export interface URL {
  long_url: string;
  short_code: string;
  clicks: number;
}

export interface URLsResponse {
  urls: URL[];
}

export interface ShortenRequest {
  url: string;
}

export interface ShortenResponse {
  short_url: string;
  code: string;
}