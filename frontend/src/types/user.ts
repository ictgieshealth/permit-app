import { Domain } from "./domain";
import { Role } from "./role";

export interface User {
  id: number;
  role_id: number;
  username: string;
  email: string;
  full_name: string;
  phone_number: string;
  nip: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
  role?: Role;
  domains?: Domain[];
}

export interface UserRequest {
  role_id: number;
  username: string;
  email: string;
  password: string;
  full_name: string;
  phone_number?: string;
  nip?: string;
  is_active?: boolean;
  domain_ids: number[];
}

export interface UserUpdateRequest {
  role_id?: number;
  username?: string;
  email?: string;
  password?: string;
  full_name?: string;
  phone_number?: string;
  nip?: string;
  is_active?: boolean;
  domain_ids?: number[];
}

export interface ChangePasswordRequest {
  old_password: string;
  new_password: string;
}

export interface UpdateProfileRequest {
  full_name?: string;
  email?: string;
  phone_number?: string;
  nip?: string;
}
