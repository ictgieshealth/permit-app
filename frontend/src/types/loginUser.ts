import { Domain } from "./domain";
import { User } from "./user";

export interface LoginRequest {
  username: string;
  password: string;
  domain_id?: number;
}

export interface LoginResponse {
  token: string;
  user: User;
  default_domain?: Domain;
}