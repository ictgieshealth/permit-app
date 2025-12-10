"use client";
import { useState, useEffect } from "react";
import { useRouter, useParams } from "next/navigation";
import Button from "@/components/ui/button/Button";
import Input from "@/components/form/input/InputField";
import Label from "@/components/form/Label";
import Link from "next/link";
import { Role } from "@/types/role";
import { Domain } from "domain";
import { authService } from "@/services/auth.service";
import { userService } from "@/services/user.service";
import { roleService } from "@/services/role.service";
import { domainService } from "@/services/domain.service";

export default function EditUserPage() {
  const router = useRouter();
  const params = useParams();
  const userId = parseInt(params.id as string);

  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const [roles, setRoles] = useState<Role[]>([]);
  const [domains, setDomains] = useState<Domain[]>([]);
  const [formData, setFormData] = useState({
    username: "",
    email: "",
    password: "",
    full_name: "",
    phone_number: "",
    nip: "",
    role_id: "",
    is_active: true,
    domain_ids: [] as number[],
  });

  useEffect(() => {
    if (!authService.hasRole(["admin"])) {
      router.push("/");
      return;
    }

    loadUser();
    loadRoles();
    loadDomains();
  }, [userId]);

  const loadUser = async () => {
    try {
      const user = await userService.getById(userId);
      setFormData({
        username: user.username,
        email: user.email,
        password: "",
        full_name: user.full_name,
        phone_number: user.phone_number || "",
        nip: user.nip || "",
        role_id: user.role_id.toString(),
        is_active: user.is_active,
        domain_ids: user.domains?.map((d) => d.id) || [],
      });
      setLoading(false);
    } catch (err: any) {
      setError(err.message || "Failed to load user");
      setLoading(false);
    }
  };

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
      setDomains(response.data);
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
      } else {
        const domainId = parseInt(value);
        setFormData({
          ...formData,
          domain_ids: checked
            ? [...formData.domain_ids, domainId]
            : formData.domain_ids.filter((id) => id !== domainId),
        });
      }
    } else {
      setFormData({ ...formData, [name]: value });
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    if (formData.domain_ids.length === 0) {
      setError("Please select at least one domain");
      return;
    }

    if (!formData.role_id) {
      setError("Please select a role");
      return;
    }

    setSaving(true);
    try {
      const updateData: any = {
        username: formData.username,
        email: formData.email,
        full_name: formData.full_name,
        phone_number: formData.phone_number,
        nip: formData.nip,
        role_id: parseInt(formData.role_id),
        is_active: formData.is_active,
        domain_ids: formData.domain_ids,
      };

      // Only include password if it's not empty
      if (formData.password) {
        updateData.password = formData.password;
      }

      await userService.update(userId, updateData);
      router.push("/users");
    } catch (err: any) {
      setError(err.message || "Failed to update user");
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-gray-500">Loading...</div>
      </div>
    );
  }

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
        <h1 className="text-2xl font-bold text-gray-800 dark:text-white">Edit User</h1>
        <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
          Update user information
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

          <div>
            <Label>Password</Label>
            <Input
              name="password"
              type="password"
              value={formData.password}
              onChange={handleChange}
              placeholder="Leave blank to keep current password"
            />
            <p className="mt-1 text-xs text-gray-500">
              Only fill this if you want to change the password
            </p>
          </div>

          <div>
            <Label>
              Role <span className="text-error-500">*</span>
            </Label>
            <select
              name="role_id"
              value={formData.role_id}
              onChange={handleChange}
              required
              className="w-full px-3 py-2 border border-gray-300 rounded-lg dark:bg-gray-800 dark:border-gray-700 dark:text-white"
            >
              <option value="">Select a role</option>
              {roles.map((role) => (
                <option key={role.id} value={role.id}>
                  {role.name}
                </option>
              ))}
            </select>
          </div>

          <div>
            <Label>
              Domains <span className="text-error-500">*</span>
            </Label>
            <div className="p-4 space-y-2 border border-gray-300 rounded-lg dark:border-gray-700">
              {domains.length === 0 ? (
                <p className="text-sm text-gray-500">No domains available</p>
              ) : (
                domains.map((domain: any) => (
                  <label
                    key={domain.id}
                    className="flex items-center gap-2 cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-800 p-2 rounded"
                  >
                    <input
                      type="checkbox"
                      value={domain.id}
                      checked={formData.domain_ids.includes(domain.id)}
                      onChange={handleChange}
                      className="w-4 h-4"
                    />
                    <span className="text-sm text-gray-700 dark:text-gray-300">
                      {domain.name} ({domain.code})
                    </span>
                  </label>
                ))
              )}
            </div>
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
            <Button type="submit" disabled={saving}>
              {saving ? "Saving..." : "Update User"}
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
