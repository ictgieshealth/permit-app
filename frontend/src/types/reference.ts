export interface ReferenceCategory {
  id: number;
  module_id: number;
  name: string;
  description?: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
  module?: {
    id: number;
    name: string;
  };
}

export interface Reference {
  id: number;
  reference_category_id: number;
  name: string;
  description?: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
  reference_category?: ReferenceCategory;
}

export interface ReferenceCategoryListRequest {
  module_id?: number;
  name?: string;
  is_active?: boolean;
  page?: number;
  limit?: number;
}

export interface ReferenceCategoryRequest {
  module_id: number;
  name: string;
  description?: string;
  is_active?: boolean;
}

export interface ReferenceListRequest {
  reference_category_id?: number;
  module_id?: number;
  name?: string;
  is_active?: boolean;
  page?: number;
  limit?: number;
}

export interface ReferenceRequest {
  reference_category_id: number;
  name: string;
  description?: string;
  is_active?: boolean;
}
