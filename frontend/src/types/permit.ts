export interface Permit {
  id: number;
  domain_id: number;
  division_id?: number | null;
  permit_type_id: number;
  name: string;
  application_type: string;
  permit_no: string;
  effective_date: string;
  expiry_date: string;
  effective_term?: string | null;
  responsible_person_id?: number | null;
  responsible_doc_person_id?: number | null;
  doc_name?: string | null;
  doc_number?: string | null;
  status: string;
  created_at: string;
  updated_at: string;
  domain?: {
    id: number;
    code: string;
    name: string;
    description: string;
    is_active: boolean;
    created_at: string;
    updated_at: string;
  };
  division?: {
    id: number;
    domain_id: number;
    code: string;
    name: string;
    created_at: string;
    updated_at: string;
  };
  permit_type?: {
    id: number;
    division_id?: number | null;
    name: string;
    risk_point?: string | null;
    default_application_type?: string | null;
    default_validity_period?: string | null;
    notes?: string | null;
    created_at: string;
    updated_at: string;
    division?: {
      id: number;
      name: string;
      code: string;
      description: string;
      created_at: string;
      updated_at: string;
    };
  };
  responsible_person?: {
    id: number;
    role_id: number;
    username: string;
    email: string;
    full_name: string;
    is_active: boolean;
  };
  responsible_doc_person?: {
    id: number;
    role_id: number;
    username: string;
    email: string;
    full_name: string;
    is_active: boolean;
  };
}

export interface PermitCreateRequest {
  domain_id: number;
  division_id?: number;
  permit_type_id: number;
  name: string;
  application_type: string;
  permit_no: string;
  effective_date: string;
  expiry_date: string;
  effective_term?: string;
  responsible_person_id?: number;
  responsible_doc_person_id?: number;
  doc_name?: string;
  doc_number?: string;
  status: string;
}

export interface PermitUpdateRequest {
  domain_id: number;
  division_id?: number;
  permit_type_id: number;
  name: string;
  application_type: string;
  permit_no: string;
  effective_date: string;
  expiry_date: string;
  effective_term?: string;
  responsible_person_id?: number;
  responsible_doc_person_id?: number;
  doc_name?: string;
  doc_number?: string;
  status: string;
}

export interface PermitListRequest {
  domain_id?: number;
  division_id?: number;
  permit_type_id?: number;
  name?: string;
  application_type?: string;
  permit_no?: string;
  responsible_person?: string;
  status?: string;
  page?: number;
  limit?: number;
}
