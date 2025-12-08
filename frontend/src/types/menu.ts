export interface Menu {
  id: number;
  name: string;
  path: string;
  icon?: string;
  parent_id?: number | null;
  order_index: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
  roles?: {
    id: number;
    name: string;
    code: string;
  }[];
  children?: Menu[];
}

export interface MenuRole {
  id: number;
  menu_id: number;
  role_id: number;
  created_at: string;
  updated_at: string;
}
