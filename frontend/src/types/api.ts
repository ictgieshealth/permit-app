export interface ApiResponse<T> {
  success: boolean;
  message: string;
  data: T;
  meta?: PageMeta;
  trace_id?: string;
}

export interface PageMeta {
  page: number;
  limit: number;
  total: number;
}

export interface ApiError {
  success: false;
  error: {
    code: string;
    message: string;
    details?: any;
  };
  trace_id?: string;
}
