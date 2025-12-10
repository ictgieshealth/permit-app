-- Migration for Permit Management System
-- Created: 2025-11-25
-- Updated: 2025-12-03 - Added domains table
-- Updated: 2025-12-06 - Added roles and users table
-- Updated: 2025-12-08 - Added menus and menu_roles tables

-- Drop tables if exists (for clean migration)
DROP TABLE IF EXISTS permits CASCADE;
DROP TABLE IF EXISTS permit_types CASCADE;
DROP TABLE IF EXISTS divisions CASCADE;
DROP TABLE IF EXISTS menu_roles CASCADE;
DROP TABLE IF EXISTS menus CASCADE;
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

-- Create Roles table
CREATE TABLE roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create Users table
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    role_id BIGINT NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE RESTRICT
);

-- Create User-Domain junction table (many-to-many relationship)
CREATE TABLE user_domains (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    domain_id BIGINT NOT NULL,
    is_default BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (domain_id) REFERENCES domains(id) ON DELETE CASCADE,
    UNIQUE(user_id, domain_id)
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

-- Create indexes for better query performance
CREATE INDEX idx_domains_code ON domains(code);
CREATE INDEX idx_domains_is_active ON domains(is_active);
CREATE INDEX idx_roles_name ON roles(name);
CREATE INDEX idx_users_role_id ON users(role_id);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_is_active ON users(is_active);
CREATE INDEX idx_user_domains_user_id ON user_domains(user_id);
CREATE INDEX idx_user_domains_domain_id ON user_domains(domain_id);
CREATE INDEX idx_user_domains_is_default ON user_domains(is_default);
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

-- Insert sample data (optional)

-- Sample Domains
INSERT INTO domains (code, name, description, is_active) VALUES
('COMPANY_A', 'Company A', 'First company using the system', true),
('COMPANY_B', 'Company B', 'Second company using the system', true),
('COMPANY_C', 'Company C', 'Third company using the system', false);

-- Sample Roles
INSERT INTO roles (name) VALUES
('admin'),
('manager'),
('employee'),
('viewer');

-- Sample Users (password: "password123" hashed with bcrypt cost 10)
INSERT INTO users (role_id, username, email, password, full_name, is_active) VALUES
(1, 'admin', 'admin@example.com', '$2a$10$JhNQe7yZL1btTXQ9q27g4Ooz09HnlGwYATSfDFVRJJYahHGiOUDxy', 'System Administrator', true),
(2, 'manager1', 'manager1@example.com', '$2a$10$JhNQe7yZL1btTXQ9q27g4Ooz09HnlGwYATSfDFVRJJYahHGiOUDxy', 'Manager One', true),
(3, 'employee1', 'employee1@example.com', '$2a$10$JhNQe7yZL1btTXQ9q27g4Ooz09HnlGwYATSfDFVRJJYahHGiOUDxy', 'Employee One', true),
(1, 'admin2', 'admin2@example.com', '$2a$10$JhNQe7yZL1btTXQ9q27g4Ooz09HnlGwYATSfDFVRJJYahHGiOUDxy', 'Administrator Two', true);

-- Sample User-Domain relationships (users can have multiple domains)
INSERT INTO user_domains (user_id, domain_id, is_default) VALUES
(1, 1, true),   -- admin has access to COMPANY_A (default)
(1, 2, false),  -- admin also has access to COMPANY_B
(2, 1, true),   -- manager1 has access to COMPANY_A
(3, 1, true),   -- employee1 has access to COMPANY_A
(3, 2, false),  -- employee1 also has access to COMPANY_B
(4, 2, true);   -- admin2 has access to COMPANY_B

-- Sample Menus
INSERT INTO menus (name, path, icon, parent_id, order_index, is_active) VALUES
('Dashboard', '/', 'GridIcon', NULL, 1, true),
('Permits', '/permits', 'ListIcon', NULL, 2, true),
('Domains', '/domains', 'PageIcon', NULL, 3, true),
('Divisions', '/divisions', 'BoxCubeIcon', NULL, 4, true),
('Permit Types', '/permit-types', 'TableIcon', NULL, 5, true),
('Users', '/users', 'UserCircleIcon', NULL, 6, true),
('Menus', '/menus', 'PlugInIcon', NULL, 7, true);

-- Assign all menus to admin role (role_id = 1)
INSERT INTO menu_roles (menu_id, role_id) VALUES
(1, 1), -- Dashboard -> Admin
(2, 1), -- Permits -> Admin
(3, 1), -- Domains -> Admin
(4, 1), -- Divisions -> Admin
(5, 1), -- Permit Types -> Admin
(6, 1), -- Users -> Admin
(7, 1); -- Menus -> Admin

-- Assign limited menus to manager role (role_id = 2)
INSERT INTO menu_roles (menu_id, role_id) VALUES
(1, 2), -- Dashboard -> Manager
(2, 2), -- Permits -> Manager
(4, 2), -- Divisions -> Manager
(5, 2); -- Permit Types -> Manager

-- Assign basic menus to employee role (role_id = 3)
INSERT INTO menu_roles (menu_id, role_id) VALUES
(1, 3), -- Dashboard -> Employee
(2, 3); -- Permits -> Employee

-- Sample Divisions
INSERT INTO divisions (domain_id, code, name) VALUES
(1, 'IT', 'Information Technology'),
(1, 'HR', 'Human Resources'),
(1, 'OPS', 'Operations'),
(2, 'FIN', 'Finance'),
(2, 'SALES', 'Sales and Marketing');

-- Sample Permit Types
INSERT INTO permit_types (division_id, name, risk_point, default_application_type, default_validity_period, notes) VALUES
(1, 'Facility', 'High', 'New Application', '1 year', 'Permit for facility and infrastructure'),
(1, 'Clinical', 'High', 'New Application', '2 years', 'Permit for clinical equipment and services'),
(2, 'Kompetensi', 'Medium', 'New Application', '1 year', 'Competency certificates for staff'),
(3, 'Operasional', 'High', 'New Application', '3 years', 'Operational permits for business activities'),
(4, 'Financial', 'Low', 'New Application', '2 years', 'Financial and regulatory permits');

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

CREATE TRIGGER update_roles_updated_at BEFORE UPDATE ON roles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_domains_updated_at BEFORE UPDATE ON user_domains
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

-- Add comments to tables
COMMENT ON TABLE domains IS 'Stores company/organization domains for multi-tenancy';
COMMENT ON TABLE roles IS 'Stores user roles for access control';
COMMENT ON TABLE users IS 'Stores user accounts with authentication';
COMMENT ON TABLE user_domains IS 'Junction table for many-to-many relationship between users and domains';
COMMENT ON TABLE menus IS 'Stores menu items for dynamic navigation';
COMMENT ON TABLE menu_roles IS 'Junction table for role-based menu access control';
COMMENT ON TABLE divisions IS 'Stores company divisions/departments';
COMMENT ON TABLE permit_types IS 'Stores types of permits available';
COMMENT ON TABLE permits IS 'Stores individual permit records';

COMMENT ON COLUMN domains.code IS 'Unique code for the domain/company';
COMMENT ON COLUMN domains.name IS 'Full name of the domain/company';
COMMENT ON COLUMN domains.description IS 'Description of the domain/company';
COMMENT ON COLUMN domains.is_active IS 'Whether the domain is active or not';

COMMENT ON COLUMN roles.name IS 'Role name (e.g., admin, manager, employee)';

COMMENT ON COLUMN users.role_id IS 'Reference to the user role';
COMMENT ON COLUMN users.username IS 'Username for login (globally unique)';
COMMENT ON COLUMN users.email IS 'Email address (globally unique)';
COMMENT ON COLUMN users.password IS 'Hashed password';
COMMENT ON COLUMN users.full_name IS 'Full name of the user';
COMMENT ON COLUMN users.is_active IS 'Whether the user account is active';

COMMENT ON COLUMN user_domains.user_id IS 'Reference to the user';
COMMENT ON COLUMN user_domains.domain_id IS 'Reference to the domain';
COMMENT ON COLUMN user_domains.is_default IS 'Whether this is the default domain for the user';

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
