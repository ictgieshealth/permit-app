export interface PermitType {
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
}
