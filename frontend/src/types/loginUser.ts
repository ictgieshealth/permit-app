import { Domain } from "./domain";
import { Role } from "./role";
import { User, UserDomainRole } from "./user";

export interface LoginRequest {
  username: string;
  password: string;
  domain_id?: number;
}

export interface LoginResponse {
  token: string;
  user: User;
  current_domain: Domain;
  current_role: Role;
  domains: UserDomainRole[]; // All domains user has access to
}

export interface SwitchDomainRequest {
  domain_id: number;
}

export interface SwitchDomainResponse {
  token: string;
  current_domain: Domain;
  current_role: Role;
}