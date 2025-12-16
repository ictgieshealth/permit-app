import { Domain } from "./domain";
import { Reference } from "./reference";
import { User } from "./user";

export interface Project {
  id: number;
  domain_id: number;
  name: string;
  code: string;
  description: string | null;
  status: boolean;
  project_status_id: number;
  started_at: string | null;
  finished_at: string | null;
  created_at: string;
  updated_at: string;
  domain?: Domain;
  project_status?: Reference;
  users?: User[];
}

export interface ProjectCreateRequest {
  domain_id: number;
  name: string;
  code?: string;
  description?: string;
  status?: boolean;
  project_status_id?: number;
  user_ids?: number[];
}

export interface ProjectUpdateRequest {
  name?: string;
  code?: string;
  description?: string;
  status?: boolean;
  project_status_id?: number;
  user_ids?: number[];
}

export interface ProjectStatusChangeRequest {
  status_id: number;
}

export interface ProjectListRequest {
  domain_id?: number;
  project_status_id?: number;
  name?: string;
  code?: string;
  status?: boolean;
  page?: number;
  limit?: number;
}
