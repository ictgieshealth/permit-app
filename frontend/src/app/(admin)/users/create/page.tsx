"use client";
import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import Button from "@/components/ui/button/Button";
import Input from "@/components/form/input/InputField";
import Label from "@/components/form/Label";
import Link from "next/link";
import { Role } from "@/types/role";
import { Domain } from "@/types/domain";
import { authService } from "@/services/auth.service";
import { roleService } from "@/services/role.service";
import { domainService } from "@/services/domain.service";
import { userService } from "@/services/user.service";

export default function CreateUserPage() {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [roles, setRoles] = useState<Role[]>([]);
  const [domains, setDomains] = useState<Domain[]>([]);
  const [formData, setFormData] = useState({
    username: "",
    email: "",
    password: "",
    confirm_password: "",
    full_name: "",
    phone_number: "",
    nip: "",
    is_active: true,
    domain_roles: [] as { domain_id: number; role_id: number; is_default: boolean }[],
  });

  useEffect(() => {
    if (!authService.hasRole(["ADMIN"])) {
      router.push("/");
      return;
    }

    loadRoles();
    loadDomains();
  }, []);

  const loadRoles = async () => {
    try {
      const response = await roleService.getAll({ limit: 100 });
      setRoles(response.data);
    } catch (err) {
      console.error("Failed to load roles:", err);
    }
  };

  const loadDomains = async () => {
    try {
      const response = await domainService.getAll({ limit: 100, is_active: true });
      setDomains(response.data as any as Domain[]);
    } catch (err) {
      console.error("Failed to load domains:", err);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value, type } = e.target;
    
    if (type === "checkbox") {
      const checked = (e.target as HTMLInputElement).checked;
      if (name === "is_active") {
        setFormData({ ...formData, [name]: checked });
      }
    } else {
      setFormData({ ...formData, [name]: value });
    }
  };

  const handleDomainRoleToggle = (domainId: number) => {
    const existingIndex = formData.domain_roles.findIndex(dr => dr.domain_id === domainId);
    if (existingIndex >= 0) {
      // Remove domain role
      setFormData({
        ...formData,
        domain_roles: formData.domain_roles.filter(dr => dr.domain_id !== domainId),
      });
    } else {
      // Add domain role (requires role selection later)
      setFormData({
        ...formData,
        domain_roles: [...formData.domain_roles, { domain_id: domainId, role_id: 0, is_default: formData.domain_roles.length === 0 }],
      });
    }
  };

  const handleRoleChange = (domainId: number, roleId: number) => {
    setFormData({
      ...formData,
      domain_roles: formData.domain_roles.map(dr =>
        dr.domain_id === domainId ? { ...dr, role_id: roleId } : dr
      ),
    });
  };

  const handleDefaultChange = (domainId: number) => {
    setFormData({
      ...formData,
      domain_roles: formData.domain_roles.map(dr => ({
        ...dr,
        is_default: dr.domain_id === domainId,
      })),
    });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    // Validation
    if (formData.password !== formData.confirm_password) {
      setError("Passwords do not match");
      return;
    }

    if (formData.domain_roles.length === 0) {
      setError("Please select at least one domain");
      return;
    }

    // Check if all domain roles have role selected
    const invalidDomainRole = formData.domain_roles.find(dr => !dr.role_id || dr.role_id === 0);
    if (invalidDomainRole) {
      setError("Please select a role for all selected domains");
      return;
    }

    // Check if at least one is default
    const hasDefault = formData.domain_roles.some(dr => dr.is_default);
    if (!hasDefault && formData.domain_roles.length > 0) {
      setError("Please set one domain-role combination as default");
      return;
    }

    setLoading(true);
    try {
      await userService.create({
        username: formData.username,
        email: formData.email,
        password: formData.password,
        full_name: formData.full_name,
        phone_number: formData.phone_number,
        nip: formData.nip,
        is_active: formData.is_active,
        domain_roles: formData.domain_roles,
      });

      router.push("/users");
    } catch (err: any) {
      setError(err.message || "Failed to create user");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-3xl p-6 mx-auto">
      <div className="mb-6">
        <div className="flex items-center gap-2 mb-2">
          <Link
            href="/users"
            className="text-gray-500 hover:text-gray-700 dark:text-gray-400"
          >
            ← Back
          </Link>
        </div>
        <h1 className="text-2xl font-bold text-gray-800 dark:text-white">Create New User</h1>
        <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
          Add a new user to the system
        </p>
      </div>

      <div className="p-6 bg-white border border-gray-200 rounded-lg dark:bg-gray-dark dark:border-gray-800">
        <form onSubmit={handleSubmit} className="space-y-6">
          {error && (
            <div className="p-3 text-sm border rounded-lg bg-error-50 border-error-200 text-error-700 dark:bg-error-900/20">
              {error}
            </div>
          )}

          <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
            <div>
              <Label>
                Username <span className="text-error-500">*</span>
              </Label>
              <Input
                name="username"
                type="text"
                value={formData.username}
                onChange={handleChange}
                required
                placeholder="Enter username"
              />
            </div>

            <div>
              <Label>
                Email <span className="text-error-500">*</span>
              </Label>
              <Input
                name="email"
                type="email"
                value={formData.email}
                onChange={handleChange}
                required
                placeholder="Enter email"
              />
            </div>
          </div>

          <div>
            <Label>
              Full Name <span className="text-error-500">*</span>
            </Label>
            <Input
              name="full_name"
              type="text"
              value={formData.full_name}
              onChange={handleChange}
              required
              placeholder="Enter full name"
            />
          </div>

          <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
            <div>
              <Label>Phone Number</Label>
              <Input
                name="phone_number"
                type="text"
                value={formData.phone_number}
                onChange={handleChange}
                placeholder="Enter phone number"
              />
            </div>

            <div>
              <Label>NIP</Label>
              <Input
                name="nip"
                type="text"
                value={formData.nip}
                onChange={handleChange}
                placeholder="Enter NIP"
              />
            </div>
          </div>

          <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
            <div>
              <Label>
                Password <span className="text-error-500">*</span>
              </Label>
              <Input
                name="password"
                type="password"
                value={formData.password}
                onChange={handleChange}
                required
                placeholder="Enter password"
              />
            </div>

            <div>
              <Label>
                Confirm Password <span className="text-error-500">*</span>
              </Label>
              <Input
                name="confirm_password"
                type="password"
                value={formData.confirm_password}
                onChange={handleChange}
                required
                placeholder="Confirm password"
              />
            </div>
          </div>

          <div>
            <Label>
              Domain & Role Assignments <span className="text-error-500">*</span>
            </Label>
            <div className="p-4 space-y-4 border border-gray-300 rounded-lg dark:border-gray-700">
              {domains.length === 0 ? (
                <p className="text-sm text-gray-500">No domains available</p>
              ) : (
                domains.map((domain) => {
                  const domainRole = formData.domain_roles.find(dr => dr.domain_id === domain.id);
                  const isSelected = !!domainRole;
                  
                  return (
                    <div key={domain.id} className="p-3 border border-gray-200 rounded-lg dark:border-gray-600">
                      <label className="flex items-center gap-2 cursor-pointer mb-2">
                        <input
                          type="checkbox"
                          checked={isSelected}
                          onChange={() => handleDomainRoleToggle(domain.id)}
                          className="w-4 h-4"
                        />
                        <span className="font-medium text-gray-700 dark:text-gray-300">
                          {domain.name} ({domain.code})
                        </span>
                      </label>
                      
                      {isSelected && (
                        <div className="ml-6 space-y-2">
                          <div>
                            <label className="block text-xs text-gray-600 dark:text-gray-400 mb-1">
                              Role for this domain:
                            </label>
                            <select
                              value={domainRole?.role_id || ""}
                              onChange={(e) => handleRoleChange(domain.id, parseInt(e.target.value))}
                              className="w-full px-2 py-1 text-sm border border-gray-300 rounded dark:bg-gray-800 dark:border-gray-600 dark:text-white"
                              required
                            >
                              <option value="">Select Role</option>
                              {roles.map((role) => (
                                <option key={role.id} value={role.id}>
                                  {role.name} ({role.category})
                                </option>
                              ))}
                            </select>
                          </div>
                          
                          <label className="flex items-center gap-2 cursor-pointer">
                            <input
                              type="radio"
                              name="default_domain"
                              checked={domainRole?.is_default || false}
                              onChange={() => handleDefaultChange(domain.id)}
                              className="w-4 h-4"
                            />
                            <span className="text-xs text-gray-600 dark:text-gray-400">
                              Set as default domain
                            </span>
                          </label>
                        </div>
                      )}
                    </div>
                  );
                })
              )}
            </div>
            <p className="mt-2 text-xs text-gray-500 dark:text-gray-400">
              Select domains and assign a role for each. Mark one as default for login.
            </p>
          </div>

          <div>
            <label className="flex items-center gap-2 cursor-pointer">
              <input
                type="checkbox"
                name="is_active"
                checked={formData.is_active}
                onChange={handleChange}
                className="w-4 h-4"
              />
              <span className="text-sm text-gray-700 dark:text-gray-300">Active</span>
            </label>
          </div>

          <div className="flex gap-3 pt-4">
            <Button type="submit" disabled={loading}>
              {loading ? "Creating..." : "Create User"}
            </Button>
            <Link href="/users">
              <Button type="button" className="bg-gray-500 hover:bg-gray-600">
                Cancel
              </Button>
            </Link>
          </div>
        </form>
      </div>
    </div>
  );
}
