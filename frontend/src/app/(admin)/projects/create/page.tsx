"use client";
import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import Button from "@/components/ui/button/Button";
import Input from "@/components/form/input/InputField";
import Label from "@/components/form/Label";
import Link from "next/link";
import { User } from "@/types/user";
import { authService } from "@/services/auth.service";
import { projectService } from "@/services/project.service";
import { userService } from "@/services/user.service";
import { useUserDomains } from "@/hooks/useUserDomains";

export default function CreateProjectPage() {
  const router = useRouter();
  const { userDomains, hasMultipleDomains, singleDomainId } = useUserDomains();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [users, setUsers] = useState<User[]>([]);
  const [selectedUsers, setSelectedUsers] = useState<number[]>([]);
  const [formData, setFormData] = useState({
    domain_id: "",
    name: "",
    code: "",
    description: "",
    status: true,
  });

  useEffect(() => {
    loadUsers();

    // Auto-select domain if user only has one
    if (singleDomainId) {
      setFormData((prev) => ({
        ...prev,
        domain_id: singleDomainId.toString(),
      }));
    }
  }, [singleDomainId]);

  const loadUsers = async () => {
    try {
      const response = await userService.getAll({ limit: 100, is_active: true });
      setUsers(response.data);
    } catch (err) {
      console.error("Failed to load users:", err);
    }
  };

  const handleChange = (
    e: React.ChangeEvent<
      HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement
    >
  ) => {
    const { name, value, type } = e.target;
    if (type === "checkbox") {
      const checked = (e.target as HTMLInputElement).checked;
      setFormData({ ...formData, [name]: checked });
    } else {
      setFormData({ ...formData, [name]: value });
    }
  };

  const handleUserToggle = (userId: number) => {
    setSelectedUsers((prev) =>
      prev.includes(userId)
        ? prev.filter((id) => id !== userId)
        : [...prev, userId]
    );
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    if (!formData.domain_id) {
      setError("Please select a domain");
      return;
    }

    if (!formData.name) {
      setError("Please enter project name");
      return;
    }

    try {
      setLoading(true);
      await projectService.create({
        domain_id: parseInt(formData.domain_id),
        name: formData.name,
        code: formData.code || undefined,
        description: formData.description || undefined,
        status: formData.status,
        user_ids: selectedUsers,
      });

      router.push("/projects");
    } catch (err: any) {
      setError(
        err.response?.data?.message || "Failed to create project"
      );
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="p-6">
      <div className="max-w-4xl mx-auto">
        <div className="mb-6">
          <h1 className="text-2xl font-semibold text-gray-900 dark:text-white">
            Create New Project
          </h1>
          <p className="mt-1 text-sm text-gray-600 dark:text-gray-400">
            Add a new project to the system
          </p>
        </div>

        {error && (
          <div className="mb-4 p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg">
            <p className="text-red-800 dark:text-red-400">{error}</p>
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-6">
          <div className="bg-white dark:bg-gray-800 shadow-md rounded-lg p-6">
            <h2 className="text-lg font-medium text-gray-900 dark:text-white mb-4">
              Project Information
            </h2>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              {/* Domain */}
              <div className="md:col-span-2">
                <Label htmlFor="domain_id">
                  Domain
                </Label>
                <select
                  id="domain_id"
                  name="domain_id"
                  value={formData.domain_id}
                  onChange={handleChange}
                  className="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700 dark:text-white"
                  required
                  disabled={!hasMultipleDomains}
                >
                  <option value="">Select Domain</option>
                  {userDomains.map((domain) => (
                    <option key={domain.id} value={domain.id}>
                      {domain.name}
                    </option>
                  ))}
                </select>
              </div>

              {/* Name */}
              <div className="md:col-span-2">
                <Label htmlFor="name">
                  Project Name
                </Label>
                <Input
                  id="name"
                  name="name"
                  type="text"
                  value={formData.name}
                  onChange={handleChange}
                  placeholder="e.g., Hospital Management System"
                  required
                />
              </div>

              {/* Code */}
              <div className="md:col-span-2">
                <Label htmlFor="code">
                  Project Code
                </Label>
                <Input
                  id="code"
                  name="code"
                  type="text"
                  value={formData.code}
                  onChange={handleChange}
                  placeholder="Leave empty to auto-generate from name"
                />
                <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
                  If left empty, will be auto-generated from project name
                </p>
              </div>

              {/* Description */}
              <div className="md:col-span-2">
                <Label htmlFor="description">Description</Label>
                <textarea
                  id="description"
                  name="description"
                  value={formData.description}
                  onChange={handleChange}
                  rows={4}
                  className="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700 dark:text-white"
                  placeholder="Enter project description..."
                />
              </div>

              {/* Status */}
              <div className="md:col-span-2">
                <div className="flex items-center">
                  <input
                    id="status"
                    name="status"
                    type="checkbox"
                    checked={formData.status}
                    onChange={handleChange}
                    className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                  />
                  <label
                    htmlFor="status"
                    className="ml-2 block text-sm text-gray-900 dark:text-gray-300"
                  >
                    Active
                  </label>
                </div>
              </div>
            </div>
          </div>

          {/* Assign Users */}
          <div className="bg-white dark:bg-gray-800 shadow-md rounded-lg p-6">
            <h2 className="text-lg font-medium text-gray-900 dark:text-white mb-4">
              Assign Users
            </h2>
            <p className="text-sm text-gray-600 dark:text-gray-400 mb-4">
              Select users to assign to this project
            </p>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {users.map((user) => (
                <div
                  key={user.id}
                  className="flex items-center p-3 border border-gray-200 dark:border-gray-600 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700"
                >
                  <input
                    id={`user-${user.id}`}
                    type="checkbox"
                    checked={selectedUsers.includes(user.id)}
                    onChange={() => handleUserToggle(user.id)}
                    className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                  />
                  <label
                    htmlFor={`user-${user.id}`}
                    className="ml-3 block text-sm text-gray-900 dark:text-gray-300 cursor-pointer"
                  >
                    <div className="font-medium">{user.full_name}</div>
                    <div className="text-xs text-gray-500 dark:text-gray-400">
                      {user.email}
                    </div>
                  </label>
                </div>
              ))}
            </div>

            {users.length === 0 && (
              <p className="text-sm text-gray-500 dark:text-gray-400 text-center py-4">
                No users available
              </p>
            )}
          </div>

          {/* Form Actions */}
          <div className="flex justify-end gap-4">
            <Link href="/projects">
              <Button variant="outline" size="md" type="button">
                Cancel
              </Button>
            </Link>
            <Button
              variant="primary"
              size="md"
              type="submit"
              disabled={loading}
            >
              {loading ? "Creating..." : "Create Project"}
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
}
