-- Migration for Permit Management System
-- Created: 2025-11-25
-- Updated: 2025-12-03 - Added domains table
-- Updated: 2025-12-06 - Added roles and users table
-- Updated: 2025-12-08 - Added menus and menu_roles tables
-- Updated: 2025-12-15 - Restructured user-role-domain for multi-app support

-- Drop tables if exists (for clean migration)
DROP TABLE IF EXISTS approval_tasks CASCADE;
DROP TABLE IF EXISTS task_files CASCADE;
DROP TABLE IF EXISTS tasks CASCADE;
DROP TABLE IF EXISTS notifications CASCADE;
DROP TABLE IF EXISTS permits CASCADE;
DROP TABLE IF EXISTS permit_types CASCADE;
DROP TABLE IF EXISTS divisions CASCADE;
DROP TABLE IF EXISTS user_projects CASCADE;
DROP TABLE IF EXISTS projects CASCADE;
DROP TABLE IF EXISTS "references" CASCADE;
DROP TABLE IF EXISTS reference_categories CASCADE;
DROP TABLE IF EXISTS modules CASCADE;
DROP TABLE IF EXISTS menu_roles CASCADE;
DROP TABLE IF EXISTS menus CASCADE;
DROP TABLE IF EXISTS user_domain_roles CASCADE;
DROP TABLE IF EXISTS user_domains CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS roles CASCADE;
DROP TABLE IF EXISTS domains CASCADE;

-- Create Domains table
CREATE TABLE domains (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create Modules table for categorizing reference data
CREATE TABLE modules (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create Reference Categories table
CREATE TABLE reference_categories (
    id BIGSERIAL PRIMARY KEY,
    module_id BIGINT NOT NULL,
    name VARCHAR(100) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (module_id) REFERENCES modules(id) ON DELETE RESTRICT
);

-- Create References table
CREATE TABLE "references" (
    id BIGSERIAL PRIMARY KEY,
    reference_category_id BIGINT NOT NULL,
    name VARCHAR(100) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (reference_category_id) REFERENCES reference_categories(id) ON DELETE CASCADE
);

-- Create Roles table with category for multi-app support
CREATE TABLE roles (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(50) NOT NULL,
    category VARCHAR(50) NOT NULL, -- 'Permit' or 'Ticketing'
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create Users table (without role_id)
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    phone_number VARCHAR(20),
    nip VARCHAR(50),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create User-Domain-Role junction table (replaces user_domains)
-- Users can have different roles in different domains
CREATE TABLE user_domain_roles (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    domain_id BIGINT NOT NULL,
    role_id BIGINT NOT NULL,
    is_default BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (domain_id) REFERENCES domains(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    UNIQUE(user_id, domain_id, role_id)
);

-- Create Menus table for dynamic menu management
CREATE TABLE menus (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    path VARCHAR(255) NOT NULL UNIQUE,
    icon VARCHAR(50),
    parent_id BIGINT,
    order_index INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (parent_id) REFERENCES menus(id) ON DELETE CASCADE
);

-- Create Menu-Role junction table for role-based menu access control
CREATE TABLE menu_roles (
    id BIGSERIAL PRIMARY KEY,
    menu_id BIGINT NOT NULL,
    role_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (menu_id) REFERENCES menus(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    UNIQUE(menu_id, role_id)
);

-- Create Divisions table
CREATE TABLE divisions (
    id BIGSERIAL PRIMARY KEY,
    domain_id BIGINT NOT NULL,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (domain_id) REFERENCES domains(id) ON DELETE RESTRICT,
    UNIQUE(domain_id, code)
);

-- Create Permit Types table
CREATE TABLE permit_types (
    id BIGSERIAL PRIMARY KEY,
    division_id BIGINT,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    risk_point VARCHAR(50),
    default_application_type VARCHAR(100),
    default_validity_period VARCHAR(100),
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (division_id) REFERENCES divisions(id) ON DELETE SET NULL
);

-- Create Permits table
CREATE TABLE permits (
    id BIGSERIAL PRIMARY KEY,
    domain_id BIGINT NOT NULL,
    division_id BIGINT,
    permit_type_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    application_type VARCHAR(100) NOT NULL,
    permit_no VARCHAR(100) NOT NULL,
    effective_date DATE NOT NULL,
    expiry_date DATE NOT NULL,
    effective_term VARCHAR(100),
    responsible_person_id BIGINT,
    responsible_doc_person_id BIGINT,
    doc_name VARCHAR(255),
    doc_number VARCHAR(100),
    doc_file_name VARCHAR(255),
    doc_file_path VARCHAR(500),
    doc_file_size BIGINT,
    doc_file_type VARCHAR(100),
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (domain_id) REFERENCES domains(id) ON DELETE RESTRICT,
    FOREIGN KEY (division_id) REFERENCES divisions(id) ON DELETE SET NULL,
    FOREIGN KEY (permit_type_id) REFERENCES permit_types(id) ON DELETE RESTRICT,
    FOREIGN KEY (responsible_person_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (responsible_doc_person_id) REFERENCES users(id) ON DELETE SET NULL,
    UNIQUE(domain_id, permit_no)
);

CREATE TABLE IF NOT EXISTS notifications (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    permit_id BIGINT NOT NULL,
    type VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    read_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_notifications_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_notifications_permit FOREIGN KEY (permit_id) REFERENCES permits(id) ON DELETE CASCADE
);
-- Create Projects table
CREATE TABLE projects (
    id BIGSERIAL PRIMARY KEY,
    domain_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(100) NOT NULL,
    description TEXT,
    status BOOLEAN DEFAULT true,
    project_status_id BIGINT NOT NULL,
    started_at TIMESTAMP,
    finished_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (domain_id) REFERENCES domains(id) ON DELETE RESTRICT,
    FOREIGN KEY (project_status_id) REFERENCES "references"(id) ON DELETE RESTRICT,
    UNIQUE(domain_id, code)
);

-- Create User-Project junction table
-- Users can be assigned to projects (domain is already in project relation)
CREATE TABLE user_projects (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    project_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    UNIQUE(user_id, project_id)
);

-- Create Tasks table for task management module
CREATE TABLE tasks (
    id BIGSERIAL PRIMARY KEY,
    domain_id BIGINT NOT NULL,
    project_id BIGINT NOT NULL,
    code VARCHAR(50) NOT NULL UNIQUE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    description_before TEXT,
    description_after TEXT,
    reason TEXT,
    revision TEXT,
    status BOOLEAN DEFAULT true,
    status_id BIGINT,
    priority_id BIGINT,
    type_id BIGINT,
    stack_id BIGINT,
    assigned_id BIGINT,
    created_by BIGINT NOT NULL,
    updated_by BIGINT,
    approved_by BIGINT,
    completed_by BIGINT,
    done_by BIGINT,
    approval_status_id BIGINT,
    start_date TIMESTAMP WITH TIME ZONE,
    due_date TIMESTAMP WITH TIME ZONE,
    completed_date TIMESTAMP WITH TIME ZONE,
    approval_date TIMESTAMP WITH TIME ZONE,
    done_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    FOREIGN KEY (domain_id) REFERENCES domains(id) ON DELETE RESTRICT,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (status_id) REFERENCES "references"(id) ON DELETE SET NULL,
    FOREIGN KEY (priority_id) REFERENCES "references"(id) ON DELETE SET NULL,
    FOREIGN KEY (type_id) REFERENCES "references"(id) ON DELETE SET NULL,
    FOREIGN KEY (stack_id) REFERENCES "references"(id) ON DELETE SET NULL,
    FOREIGN KEY (assigned_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE RESTRICT,
    FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (approved_by) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (completed_by) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (done_by) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (approval_status_id) REFERENCES "references"(id) ON DELETE SET NULL,
    UNIQUE(domain_id, code)
);

-- Create Task Files table for task attachments
CREATE TABLE task_files (
    id BIGSERIAL PRIMARY KEY,
    task_id BIGINT NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_size VARCHAR(50),
    file_type VARCHAR(100),
    task_file_type INT,
    status BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
);

-- Create Approval Tasks table for task approval workflow
CREATE TABLE approval_tasks (
    id BIGSERIAL PRIMARY KEY,
    task_id BIGINT NOT NULL,
    sequence SMALLINT NOT NULL,
    approved_by BIGINT,
    approval_status_id BIGINT,
    approval_date TIMESTAMP WITH TIME ZONE,
    note VARCHAR(500),
    status BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    FOREIGN KEY (approved_by) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (approval_status_id) REFERENCES "references"(id) ON DELETE SET NULL,
    UNIQUE(task_id, sequence)
);

-- Create indexes for better query performance
CREATE INDEX idx_domains_code ON domains(code);
CREATE INDEX idx_domains_is_active ON domains(is_active);
CREATE INDEX idx_modules_code ON modules(code);
CREATE INDEX idx_modules_is_active ON modules(is_active);
CREATE INDEX idx_reference_categories_module_id ON reference_categories(module_id);
CREATE INDEX idx_reference_categories_is_active ON reference_categories(is_active);
CREATE INDEX idx_references_reference_category_id ON "references"(reference_category_id);
CREATE INDEX idx_references_is_active ON "references"(is_active);
CREATE INDEX idx_roles_code ON roles(code);
CREATE INDEX idx_roles_name ON roles(name);
CREATE INDEX idx_roles_category ON roles(category);
CREATE INDEX idx_user_projects_user_id ON user_projects(user_id);
CREATE INDEX idx_user_projects_project_id ON user_projects(project_id);
CREATE INDEX idx_tasks_domain_id ON tasks(domain_id);
CREATE INDEX idx_tasks_project_id ON tasks(project_id);
CREATE INDEX idx_tasks_code ON tasks(code);
CREATE INDEX idx_tasks_status_id ON tasks(status_id);
CREATE INDEX idx_tasks_priority_id ON tasks(priority_id);
CREATE INDEX idx_tasks_type_id ON tasks(type_id);
CREATE INDEX idx_tasks_assigned_id ON tasks(assigned_id);
CREATE INDEX idx_tasks_created_by ON tasks(created_by);
CREATE INDEX idx_tasks_approval_status_id ON tasks(approval_status_id);
CREATE INDEX idx_tasks_deleted_at ON tasks(deleted_at);
CREATE INDEX idx_task_files_task_id ON task_files(task_id);
CREATE INDEX idx_task_files_task_file_type ON task_files(task_file_type);
CREATE INDEX idx_approval_tasks_task_id ON approval_tasks(task_id);
CREATE INDEX idx_approval_tasks_approved_by ON approval_tasks(approved_by);
CREATE INDEX idx_approval_tasks_approval_status_id ON approval_tasks(approval_status_id);
CREATE INDEX idx_approval_tasks_sequence ON approval_tasks(sequence);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_projects_domain_id ON projects(domain_id);
CREATE INDEX idx_projects_code ON projects(code);
CREATE INDEX idx_projects_project_status_id ON projects(project_status_id);
CREATE INDEX idx_projects_status ON projects(status);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_is_active ON users(is_active);
CREATE INDEX idx_user_domain_roles_user_id ON user_domain_roles(user_id);
CREATE INDEX idx_user_domain_roles_domain_id ON user_domain_roles(domain_id);
CREATE INDEX idx_user_domain_roles_role_id ON user_domain_roles(role_id);
CREATE INDEX idx_user_domain_roles_is_default ON user_domain_roles(is_default);
CREATE INDEX idx_menus_parent_id ON menus(parent_id);
CREATE INDEX idx_menus_is_active ON menus(is_active);
CREATE INDEX idx_menus_order_index ON menus(order_index);
CREATE INDEX idx_menus_path ON menus(path);
CREATE INDEX idx_menu_roles_menu_id ON menu_roles(menu_id);
CREATE INDEX idx_menu_roles_role_id ON menu_roles(role_id);
CREATE INDEX idx_divisions_domain_id ON divisions(domain_id);
CREATE INDEX idx_divisions_code ON divisions(code);
CREATE INDEX idx_permit_types_division_id ON permit_types(division_id);
CREATE INDEX idx_permit_types_name ON permit_types(name);
CREATE INDEX idx_permit_types_code ON permit_types(code);
CREATE INDEX idx_permits_domain_id ON permits(domain_id);
CREATE INDEX idx_permits_division_id ON permits(division_id);
CREATE INDEX idx_permits_permit_type_id ON permits(permit_type_id);
CREATE INDEX idx_permits_name ON permits(name);
CREATE INDEX idx_permits_permit_no ON permits(permit_no);
CREATE INDEX idx_permits_responsible_person_id ON permits(responsible_person_id);
CREATE INDEX idx_permits_responsible_doc_person_id ON permits(responsible_doc_person_id);
CREATE INDEX idx_permits_status ON permits(status);
CREATE INDEX idx_permits_effective_date ON permits(effective_date);
CREATE INDEX idx_permits_expiry_date ON permits(expiry_date);
CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_permit_id ON notifications(permit_id);
CREATE INDEX IF NOT EXISTS idx_notifications_is_read ON notifications(is_read);
CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_notifications_type ON notifications(type);

-- Insert sample data (optional)

-- Sample Domains
INSERT INTO domains (code, name, description, is_active) VALUES
('RS-TRIADIPA', 'RS TRIADIPA', 'RS TRIADIPA', true),
('RS-BHC-SUMENEP', 'RS BHC SUMENEP', 'RS BHC SUMENEP', true);

-- Sample Modules
INSERT INTO modules (code, name, is_active) VALUES
('task', 'Task Management', true),
('project', 'Project Management', true),
('permit', 'Permit Management', true);

-- Sample Reference Categories
INSERT INTO reference_categories (id, module_id, name, is_active) VALUES
(1, 1, 'Task Status', true),
(2, 1, 'Task Priority', true),
(3, 1, 'Task Type', true),
(4, 1, 'Task Stack', true),
(5, 1, 'Task Approval Status', true),
(6, 2, 'Project Status', true),
(7, 1, 'Task File Type', true);

-- Sample References for Task Management
INSERT INTO "references" (id, reference_category_id, name, is_active) VALUES
-- Task Status (category_id = 1)
(1, 1, 'To Do', true),
(2, 1, 'On Hold', true),
(3, 1, 'On Progress', true),
(4, 1, 'Done', true),
(37, 1, 'In Review', true),
(39, 1, 'Revision', true),
-- Task Priority (category_id = 2)
(5, 2, 'Low', true),
(6, 2, 'Medium', true),
(7, 2, 'High', true),
-- Task Type (category_id = 3)
(24, 3, 'Maintenance', true),
(25, 3, 'Development', true),
-- Task Stack (category_id = 4)
(12, 4, 'BE', true),
(13, 4, 'Web', true),
(14, 4, 'Mobile - Console', true),
(15, 4, 'Mobile - User', true),
(16, 4, 'Dev - Ops', true),
-- Task Approval Status (category_id = 5)
(20, 5, 'Waiting', true),
(21, 5, 'Reject', true),
(22, 5, 'Approve', true),
(23, 5, 'Pending Manager', true),
-- Project Status (category_id = 6)
(26, 6, 'Pending', true),
(27, 6, 'On Hold', true),
(28, 6, 'On Progress', true),
(29, 6, 'Done', true),
-- Task File Type (category_id = 7)
(30, 7, 'Create', true),
(31, 7, 'Before', true),
(32, 7, 'After', true),
(38, 7, 'File Revision', true);

-- Sample Roles with category (Permit and Ticketing)
INSERT INTO roles (code, name, category, description) VALUES
-- Permit roles
('ADMIN', 'Admin', 'Permit', 'System administrator with full access to permit management'),
('PERMIT_MANAGER', 'Permit Manager', 'Permit', 'Manager for permit operations'),
('PERMIT_EMPLOYEE', 'Permit Employee', 'Permit', 'Employee handling permit tasks'),
-- Ticketing roles
('TICKETING_DEVELOPER', 'Ticketing Developer', 'Ticketing', 'Developer for ticketing system'),
('TICKETING_PIC', 'Ticketing PIC', 'Ticketing', 'Client PIC handling ticketing tasks'),
('TICKETING_MANAGER', 'Ticketing Manager', 'Ticketing', 'Client Manager for ticketing operations'),
('TICKETING_HEAD_OF_UNIT', 'Ticketing Head of Unit', 'Ticketing', 'Client Head of Unit for ticketing system');

-- Sample Users (password: "password123" hashed with bcrypt cost 10)
INSERT INTO users (username, email, password, full_name, is_active) VALUES
('admin', 'admin@example.com', '$2a$10$JhNQe7yZL1btTXQ9q27g4Ooz09HnlGwYATSfDFVRJJYahHGiOUDxy', 'System Administrator', true),
('manager1', 'manager1@example.com', '$2a$10$JhNQe7yZL1btTXQ9q27g4Ooz09HnlGwYATSfDFVRJJYahHGiOUDxy', 'Permit Manager One', true),
('employee1', 'employee1@example.com', '$2a$10$JhNQe7yZL1btTXQ9q27g4Ooz09HnlGwYATSfDFVRJJYahHGiOUDxy', 'Permit Employee One', true),
('ticketing_pic', 'ticketing.pic@example.com', '$2a$10$JhNQe7yZL1btTXQ9q27g4Ooz09HnlGwYATSfDFVRJJYahHGiOUDxy', 'Ticketing PIC', true),
('ticketing_dev', 'ticketing.dev@example.com', '$2a$10$JhNQe7yZL1btTXQ9q27g4Ooz09HnlGwYATSfDFVRJJYahHGiOUDxy', 'Ticketing Developer', true),
('ticketing_manager', 'ticketing.manager@example.com', '$2a$10$JhNQe7yZL1btTXQ9q27g4Ooz09HnlGwYATSfDFVRJJYahHGiOUDxy', 'Ticketing Manager One', true),
('ticketing_headunit', 'ticketing.headunit@example.com', '$2a$10$JhNQe7yZL1btTXQ9q27g4Ooz09HnlGwYATSfDFVRJJYahHGiOUDxy', 'Ticketing Head of Unit', true);

-- Sample User-Domain-Role relationships (users can have different roles in different domains)
INSERT INTO user_domain_roles (user_id, domain_id, role_id, is_default) VALUES
-- Admin has Admin role in COMPANY_A (default) and COMPANY_B
(1, 1, 1, true),   -- admin -> COMPANY_A -> Admin (default)
(1, 2, 1, false),  -- admin -> COMPANY_B -> Admin
-- Permit Manager has Permit Manager role in COMPANY_A
(2, 1, 2, true),   -- manager1 -> COMPANY_A -> Permit Manager
-- Permit Employee has Permit Employee role in COMPANY_A and COMPANY_B
(3, 1, 3, true),   -- employee1 -> COMPANY_A -> Permit Employee (default)
(3, 2, 3, false),  -- employee1 -> COMPANY_B -> Permit Employee
-- Ticketing PIC has Ticketing PIC role in COMPANY_B
(4, 1, 5, true),   -- ticketing_pic -> COMPANY_A -> Ticketing PIC
(4, 2, 5, false),   -- ticketing_pic -> COMPANY_B -> Ticketing PIC
-- Ticketing Developer has Ticketing Developer role in COMPANY_B
(5, 1, 4, false),   -- ticketing_dev -> COMPANY_A -> Ticketing Developer
(5, 2, 4, false),   -- ticketing_dev -> COMPANY_B -> Ticketing Developer
-- Ticketing Manager has Ticketing Manager role in COMPANY_A and COMPANY_B
(6, 1, 6, true),   -- ticketing_manager -> COMPANY_A -> Ticketing Manager (default)
(6, 2, 6, false),  -- ticketing_manager -> COMPANY_B -> Ticketing Manager
-- Ticketing Head of Unit has Ticketing Head of Unit role in COMPANY_A
(7, 1, 7, true),   -- ticketing_headunit -> COMPANY_A -> Ticketing Head of Unit
(7, 2, 7, false);  -- ticketing_headunit -> COMPANY_B -> Ticketing Head of Unit

-- Sample Menus
INSERT INTO menus (name, path, icon, parent_id, order_index, is_active) VALUES
('Dashboard', '/', 'GridIcon', NULL, 1, true),
('Permits', '/permits', 'ListIcon', NULL, 2, true),
('Tasks', '/tasks', 'ListIcon', NULL, 3, true),
('Task Requests', '/tasks/task-requests', 'PageIcon', NULL, 4, true),
('Projects', '/projects', 'PageIcon', NULL, 5, true),
('Domains', '/domains', 'PageIcon', NULL, 6, true),
('Divisions', '/divisions', 'BoxCubeIcon', NULL, 7, true),
('Permit Types', '/permit-types', 'TableIcon', NULL, 8, true),
('Roles', '/roles', 'UserCircleIcon', NULL, 9, true),
('Users', '/users', 'UserCircleIcon', NULL, 10, true),
('Menus', '/menus', 'PlugInIcon', NULL, 11, true);

-- Assign all menus to admin role (role_id = 1)
INSERT INTO menu_roles (menu_id, role_id) VALUES
(1, 1), -- Dashboard -> Admin
(2, 1), -- Permits -> Admin
(3, 1), -- Tasks -> Admin
(4, 1), -- Task Requests -> Admin
(5, 1), -- Projects -> Admin
(6, 1), -- Domains -> Admin
(7, 1), -- Divisions -> Admin
(8, 1), -- Permit Types -> Admin
(9, 1), -- Roles -> Admin
(10, 1), -- Users -> Admin
(11, 1); -- Menus -> Admin

-- Assign limited menus to manager role (role_id = 2)
INSERT INTO menu_roles (menu_id, role_id) VALUES
(1, 2), -- Dashboard -> Manager
(2, 2), -- Permits -> Manager
(3, 2), -- Tasks -> Manager
(4, 2), -- Task Requests -> Manager
(5, 2), -- Projects -> Manager
(7, 2), -- Divisions -> Manager
(8, 2); -- Permit Types -> Manager

-- Assign basic menus to employee role (role_id = 3)
INSERT INTO menu_roles (menu_id, role_id) VALUES
(1, 3), -- Dashboard -> Employee
(2, 3); -- Permits -> Employee

-- role_id = 4 (Ticketing Developer)
INSERT INTO menu_roles (menu_id, role_id) VALUES
(1, 4), -- Dashboard -> Ticketing Developer
(3, 4), -- Tasks -> Ticketing Developer
(4, 4), -- Task Requests -> Ticketing Developer
(5, 4); -- Projects -> Ticketing Developer

-- role_id = 5 (Ticketing PIC)
INSERT INTO menu_roles (menu_id, role_id) VALUES
(1, 5), -- Dashboard -> Ticketing PIC
(3, 5), -- Tasks -> Ticketing PIC
(4, 5); -- Task Requests -> Ticketing PIC

-- role_id = 6 (Ticketing Manager)
INSERT INTO menu_roles (menu_id, role_id) VALUES
(1, 6), -- Dashboard -> Ticketing Manager
(3, 6), -- Tasks -> Ticketing Manager
(4, 6), -- Task Requests -> Ticketing Manager
(5, 6); -- Projects -> Ticketing Manager

-- role_id = 7 (Ticketing Head of Unit)
INSERT INTO menu_roles (menu_id, role_id) VALUES
(1, 7), -- Dashboard -> Ticketing Head of Unit
(3, 7), -- Tasks -> Ticketing Head of Unit
(4, 7), -- Task Requests -> Ticketing Head of Unit
(5, 7); -- Projects -> Ticketing Head of Unit

-- Sample Divisions
INSERT INTO divisions (domain_id, code, name) VALUES
(1, 'IT', 'Information Technology'),
(1, 'HR', 'Human Resources'),
(1, 'OPS', 'Operations'),
(2, 'FIN', 'Finance'),
(2, 'SALES', 'Sales and Marketing');

-- Sample Permit Types
INSERT INTO permit_types (division_id, code, name, risk_point, default_application_type, default_validity_period, notes) VALUES
(1, 'FAC', 'Facility', 'High', 'New Application', '1 year', 'Permit for facility and infrastructure'),
(1, 'CLIN', 'Clinical', 'High', 'New Application', '2 years', 'Permit for clinical equipment and services'),
(2, 'KOMP', 'Kompetensi', 'Medium', 'New Application', '1 year', 'Competency certificates for staff'),
(3, 'OPER', 'Operasional', 'High', 'New Application', '3 years', 'Operational permits for business activities'),
(4, 'FIN', 'Financial', 'Low', 'New Application', '2 years', 'Financial and regulatory permits');
-- Sample Projects (using project_status_id = 26 for 'Pending' from references)
INSERT INTO projects (domain_id, name, code, description, status, project_status_id) VALUES
(1, 'RS TRIA DIPA', 'rs-tria-dipa', 'Hospital Management System for RS TRIA DIPA', true, 26),
(2, 'RS BHC SUMENEP', 'rs-bhc-sumenep', 'Hospital Management System for RS BHC SUMENEP', true, 26);

-- Sample User-Project assignments
-- Admin is assigned to both projects
INSERT INTO user_projects (user_id, project_id) VALUES
(1, 1),  -- admin -> RS TRIA DIPA project
(1, 2),  -- admin -> RS BHC SUMENEP project
-- Permit Manager assigned to project 1
(2, 1),  -- manager1 -> RS TRIA DIPA project
-- Permit Employee assigned to both projects
(3, 1, 1),  -- employee1 -> RS-TRIADIPA domain -> RS TRIA DIPA project
(3, 2, 2),  -- employee1 -> RS-BHC-SUMENEP domain -> RS BHC SUMENEP project
-- Ticketing Developer assigned to project 2
(4, 2, 2),  -- ticketing_dev -> RS-BHC-SUMENEP domain -> RS BHC SUMENEP project
-- Ticketing Manager assigned to both projects
(5, 1, 1),  -- ticketing_manager -> RS-TRIADIPA domain -> RS TRIA DIPA project
(5, 2, 2);  -- ticketing_manager -> RS-BHC-SUMENEP domain -> RS BHC SUMENEP project

-- Sample Permits
INSERT INTO permits (domain_id, division_id, permit_type_id, name, application_type, permit_no, effective_date, expiry_date, effective_term, responsible_person_id, responsible_doc_person_id, doc_name, doc_number, status) VALUES
(1, 1, 1, 'Building Construction', 'New Application', 'PRM-2025-001', '2025-01-01', '2026-01-01', '1 year', 1, 2, 'Building Construction Permit', 'DOC-2025-001', 'active'),
(1, 1, 2, 'Medical Laboratory', 'New Application', 'PRM-2025-002', '2025-01-15', '2027-01-15', '2 years', 2, 1, 'Clinical Lab Permit', 'DOC-2025-002', 'active'),
(1, 2, 3, 'Nurse Certification', 'Renewal', 'PRM-2025-003', '2025-02-01', '2026-02-01', '1 year', 3, 2, 'Competency Certificate', 'DOC-2025-003', 'active'),
(2, 4, 4, 'Surgery Room Operations', 'New Application', 'PRM-2025-004', '2025-03-01', '2028-03-01', '3 years', 4, 4, 'Operational Permit', 'DOC-2025-004', 'active'),
(2, 5, 5, 'Financial System', 'New Application', 'PRM-2025-005', '2025-04-01', '2027-04-01', '2 years', 4, 1, 'Financial License', 'DOC-2025-005', 'active');

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Triggers for auto-updating updated_at
CREATE TRIGGER update_domains_updated_at BEFORE UPDATE ON domains
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_modules_updated_at BEFORE UPDATE ON modules
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_reference_categories_updated_at BEFORE UPDATE ON reference_categories
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_references_updated_at BEFORE UPDATE ON "references"
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_roles_updated_at BEFORE UPDATE ON roles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_domain_roles_updated_at BEFORE UPDATE ON user_domain_roles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_menus_updated_at BEFORE UPDATE ON menus
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_menu_roles_updated_at BEFORE UPDATE ON menu_roles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_divisions_updated_at BEFORE UPDATE ON divisions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_permit_types_updated_at BEFORE UPDATE ON permit_types
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_permits_updated_at BEFORE UPDATE ON permits
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_projects_updated_at BEFORE UPDATE ON projects
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_projects_updated_at BEFORE UPDATE ON user_projects
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Add comments to tables
COMMENT ON TABLE domains IS 'Stores company/organization domains for multi-tenancy';
COMMENT ON TABLE modules IS 'Stores modules for categorizing reference data (Task, Project, Permit)';
COMMENT ON TABLE reference_categories IS 'Stores reference categories grouped by module';
COMMENT ON TABLE "references" IS 'Stores reference data values for various categories';
COMMENT ON TABLE roles IS 'Stores user roles for access control with category for multi-app support (Permit/Ticketing)';
COMMENT ON TABLE users IS 'Stores user accounts with authentication (roles assigned through user_domain_roles)';
COMMENT ON TABLE user_domain_roles IS 'Junction table for users, domains, and roles - allows users to have different roles in different domains';
COMMENT ON TABLE menus IS 'Stores menu items for dynamic navigation';
COMMENT ON TABLE menu_roles IS 'Junction table for role-based menu access control';
COMMENT ON TABLE divisions IS 'Stores company divisions/departments';
COMMENT ON TABLE permit_types IS 'Stores types of permits available';
COMMENT ON TABLE permits IS 'Stores individual permit records';

COMMENT ON COLUMN domains.code IS 'Unique code for the domain/company';
COMMENT ON COLUMN domains.name IS 'Full name of the domain/company';
COMMENT ON COLUMN domains.description IS 'Description of the domain/company';
COMMENT ON COLUMN domains.is_active IS 'Whether the domain is active or not';

COMMENT ON COLUMN roles.code IS 'Unique code for the role (e.g., ADMIN, PERMIT_MANAGER)';
COMMENT ON COLUMN roles.name IS 'Display name of the role';
COMMENT ON COLUMN roles.category IS 'Category of the role: Permit or Ticketing';
COMMENT ON COLUMN roles.description IS 'Description of the role and its permissions';

COMMENT ON COLUMN users.username IS 'Username for login (globally unique)';
COMMENT ON COLUMN users.email IS 'Email address (globally unique)';
COMMENT ON COLUMN users.password IS 'Hashed password';
COMMENT ON COLUMN users.full_name IS 'Full name of the user';
COMMENT ON COLUMN users.is_active IS 'Whether the user account is active';
COMMENT ON COLUMN users.phone_number IS 'User phone number';
COMMENT ON COLUMN users.nip IS 'Nomor Induk Pegawai (Employee ID Number)';

COMMENT ON COLUMN user_domain_roles.user_id IS 'Reference to the user';
COMMENT ON COLUMN user_domain_roles.domain_id IS 'Reference to the domain';
COMMENT ON COLUMN user_domain_roles.role_id IS 'Reference to the role';
COMMENT ON COLUMN user_domain_roles.is_default IS 'Whether this is the default domain-role combination for the user';

COMMENT ON COLUMN menus.name IS 'Display name of the menu item';
COMMENT ON COLUMN menus.path IS 'URL path for the menu item (unique)';
COMMENT ON COLUMN menus.icon IS 'Icon name for the menu item';
COMMENT ON COLUMN menus.parent_id IS 'Reference to parent menu for hierarchical structure';
COMMENT ON COLUMN menus.order_index IS 'Display order of menu items';
COMMENT ON COLUMN menus.is_active IS 'Whether the menu item is active';

COMMENT ON COLUMN menu_roles.menu_id IS 'Reference to the menu item';
COMMENT ON COLUMN menu_roles.role_id IS 'Reference to the role that can access this menu';

COMMENT ON COLUMN divisions.domain_id IS 'Reference to the domain/company';
COMMENT ON COLUMN divisions.code IS 'Code for the division (unique within domain)';
COMMENT ON COLUMN divisions.name IS 'Full name of the division';

COMMENT ON COLUMN permit_types.division_id IS 'Reference to the division that manages this permit type';
COMMENT ON COLUMN permit_types.name IS 'Category of the permit (e.g., Facility, Clinical, Kompetensi, Operasional)';
COMMENT ON COLUMN permit_types.risk_point IS 'Risk level associated with this permit type';
COMMENT ON COLUMN permit_types.default_application_type IS 'Default type of application (e.g., New Application, Renewal)';
COMMENT ON COLUMN permit_types.default_validity_period IS 'Default validity period for this permit type';

COMMENT ON COLUMN permits.domain_id IS 'Reference to the domain/company';
COMMENT ON COLUMN permits.division_id IS 'Reference to the division that manages this permit';
COMMENT ON COLUMN permits.permit_type_id IS 'Reference to the type of permit (category)';
COMMENT ON COLUMN permits.name IS 'Name of permit (equipment/service/competency/operational name)';
COMMENT ON COLUMN permits.application_type IS 'Type of application (New Application, Renewal, Amendment)';
COMMENT ON COLUMN permits.permit_no IS 'Permit number (unique within domain)';
COMMENT ON COLUMN permits.effective_date IS 'Date when the permit becomes effective';
COMMENT ON COLUMN permits.expiry_date IS 'Date when the permit expires';
COMMENT ON COLUMN permits.effective_term IS 'Duration of the permit validity';
COMMENT ON COLUMN permits.responsible_person_id IS 'Reference to user responsible for this permit';
COMMENT ON COLUMN permits.responsible_doc_person_id IS 'Reference to user responsible for permit documentation';
COMMENT ON COLUMN permits.status IS 'Current status of the permit (active, expired, revoked)';
