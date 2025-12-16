import { Domain } from "./domain";
import { Role } from "./role";

export interface UserDomainRole {
  user_id: number;
  domain_id: number;
  role_id: number;
  is_default: boolean;
  domain?: Domain;
  role?: Role;
}

export interface User {
  id: number;
  username: string;
  email: string;
  full_name: string;
  phone_number: string;
  nip: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
  domain_roles?: UserDomainRole[];
}

export interface DomainRoleRequest {
  domain_id: number;
  role_id: number;
  is_default: boolean;
}

export interface UserRequest {
  username: string;
  email: string;
  password: string;
  full_name: string;
  phone_number?: string;
  nip?: string;
  is_active?: boolean;
  domain_roles: DomainRoleRequest[];
}

export interface UserUpdateRequest {
  username?: string;
  email?: string;
  password?: string;
  full_name?: string;
  phone_number?: string;
  nip?: string;
  is_active?: boolean;
  domain_roles?: DomainRoleRequest[];
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
