export interface Task {
  id: number;
  domain_id: number;
  project_id: number;
  code: string;
  title: string;
  description: string;
  description_before?: string;
  description_after?: string;
  reason?: string;
  revision?: string;
  status: boolean;
  status_id?: number;
  priority_id?: number;
  type_id?: number;
  stack_id?: number;
  assigned_id?: number;
  created_by: number;
  updated_by?: number;
  approved_by?: number;
  completed_by?: number;
  done_by?: number;
  approval_status_id?: number;
  start_date?: string;
  due_date?: string;
  completed_date?: string;
  approval_date?: string;
  done_at?: string;
  created_at: string;
  updated_at: string;
  project?: {
    id: number;
    code: string;
    name: string;
  };
  status_task?: Reference;
  priority?: Reference;
  type?: Reference;
  stack?: Reference;
  assignee?: UserBasic;
  creator?: UserBasic;
  approval_status?: Reference;
  task_files?: TaskFile[];
  approval_tasks?: ApprovalTask[];
}

export interface TaskFile {
  id: number;
  task_id: number;
  file_name: string;
  file_path: string;
  file_size?: string;
  file_type?: string;
  task_file_type?: number;
  status: boolean;
  created_at: string;
  updated_at: string;
}

export interface ApprovalTask {
  id: number;
  task_id: number;
  sequence: number;
  approved_by?: number;
  approval_status_id?: number;
  approval_date?: string;
  note?: string;
  status: boolean;
  created_at: string;
  updated_at: string;
  approver?: UserBasic;
  approved_by_user?: UserBasic;
  approval_status?: Reference;
}

export interface Reference {
  id: number;
  name: string;
}

export interface UserBasic {
  id: number;
  username: string;
  full_name: string;
  email: string;
}

export interface TaskListRequest {
  search?: string;
  project_id?: number;
  status_id?: number;
  approval_status_id?: number;
  assigned_id?: number;
  start_date?: string;
  end_date?: string;
  page?: number;
  limit?: number;
}

export interface TaskCreateRequest {
  project_id: number;
  title: string;
  description: string;
  priority_id: number;
  stack_id: number;
  assigned_id?: number;
  due_date?: string;
}

export interface TaskUpdateRequest {
  project_id: number;
  title: string;
  description: string;
  description_before?: string;
  description_after?: string;
  priority_id: number;
  stack_id: number;
  type_id?: number;
  assigned_id?: number;
  due_date?: string;
}

export interface TaskChangeStatusRequest {
  status_id: number;
}

export interface TaskChangeTypeRequest {
  type_id: number;
}

export interface TaskInReviewRequest {
  description_before: string;
  description_after: string;
}

export interface TaskReasonRequest {
  reason: string;
}

export interface TaskRevisionRequest {
  revision: string;
}

export interface ApprovalRequest {
  note?: string;
}
