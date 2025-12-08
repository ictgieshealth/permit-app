export interface Division {
  id: number;
  domain_id: number;
  name: string;
  code: string;
  description: string;
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
}
