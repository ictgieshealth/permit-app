"use client";
import { useState, useEffect } from "react";
import { useRouter, useParams } from "next/navigation";
import Button from "@/components/ui/button/Button";
import Input from "@/components/form/input/InputField";
import Label from "@/components/form/Label";
import Link from "next/link";
import { User } from "@/types/user";
import { Project } from "@/types/project";
import { projectService } from "@/services/project.service";
import { userService } from "@/services/user.service";

export default function EditProjectPage() {
  const router = useRouter();
  const params = useParams();
  const projectId = params.id as string;

  const [loading, setLoading] = useState(false);
  const [fetchingProject, setFetchingProject] = useState(true);
  const [error, setError] = useState("");
  const [users, setUsers] = useState<User[]>([]);
  const [selectedUsers, setSelectedUsers] = useState<number[]>([]);
  const [project, setProject] = useState<Project | null>(null);
  const [formData, setFormData] = useState({
    name: "",
    code: "",
    description: "",
    status: true,
  });

  useEffect(() => {
    loadProject();
    loadUsers();
  }, [projectId]);

  const loadProject = async () => {
    try {
      setFetchingProject(true);
      const data = await projectService.getById(parseInt(projectId));
      setProject(data);
      setFormData({
        name: data.name,
        code: data.code,
        description: data.description || "",
        status: data.status,
      });
      setSelectedUsers(data.users?.map((u) => u.id) || []);
    } catch (err: any) {
      setError(err.response?.data?.message || "Failed to load project");
    } finally {
      setFetchingProject(false);
    }
  };

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

    if (!formData.name) {
      setError("Please enter project name");
      return;
    }

    try {
      setLoading(true);
      await projectService.update(parseInt(projectId), {
        name: formData.name,
        code: formData.code || undefined,
        description: formData.description || undefined,
        status: formData.status,
        user_ids: selectedUsers,
      });

      router.push("/projects");
    } catch (err: any) {
      setError(err.response?.data?.message || "Failed to update project");
    } finally {
      setLoading(false);
    }
  };

  if (fetchingProject) {
    return (
      <div className="p-6 flex justify-center items-center h-screen">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  if (!project) {
    return (
      <div className="p-6">
        <div className="max-w-4xl mx-auto">
          <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4">
            <p className="text-red-800 dark:text-red-400">Project not found</p>
          </div>
          <div className="mt-4">
            <Link href="/projects">
              <Button variant="outline" size="md">
                Back to Projects
              </Button>
            </Link>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="p-6">
      <div className="max-w-4xl mx-auto">
        <div className="mb-6">
          <h1 className="text-2xl font-semibold text-gray-900 dark:text-white">
            Edit Project
          </h1>
          <p className="mt-1 text-sm text-gray-600 dark:text-gray-400">
            Update project information
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
              {/* Domain (Read-only) */}
              <div className="md:col-span-2">
                <Label htmlFor="domain">Domain</Label>
                <Input
                  id="domain"
                  type="text"
                  value={project.domain?.name || ""}
                  disabled
                  className="bg-gray-100 dark:bg-gray-700"
                />
                <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
                  Domain cannot be changed after creation
                </p>
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
                <Label htmlFor="code">Project Code</Label>
                <Input
                  id="code"
                  name="code"
                  type="text"
                  value={formData.code}
                  onChange={handleChange}
                  placeholder="Project code"
                />
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

              {/* Current Status */}
              <div className="md:col-span-2">
                <Label htmlFor="project_status">Project Status</Label>
                <Input
                  id="project_status"
                  type="text"
                  value={project.project_status?.name || ""}
                  disabled
                  className="bg-gray-100 dark:bg-gray-700"
                />
                <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
                  Use status action buttons in the list page to change project
                  status
                </p>
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
              {loading ? "Updating..." : "Update Project"}
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
}
