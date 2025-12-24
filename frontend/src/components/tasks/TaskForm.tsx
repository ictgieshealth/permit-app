"use client";
import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import Button from "@/components/ui/button/Button";
import Input from "@/components/form/input/InputField";
import Label from "@/components/form/Label";
import Link from "next/link";
import { Task, TaskCreateRequest, TaskUpdateRequest } from "@/types/task";
import { taskService } from "@/services/task.service";
import { projectService } from "@/services/project.service";
import { referenceService } from "@/services/reference.service";
import { userService } from "@/services/user.service";
import FileUpload from "@/components/ui/FileUpload";

interface TaskFormProps {
  code?: string;
  isEdit?: boolean;
}

export default function TaskForm({ code, isEdit = false }: TaskFormProps) {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [loadingData, setLoadingData] = useState(isEdit);
  const [error, setError] = useState("");
  const [projects, setProjects] = useState<any[]>([]);
  const [priorities, setPriorities] = useState<any[]>([]);
  const [stacks, setStacks] = useState<any[]>([]);
  const [users, setUsers] = useState<any[]>([]);
  const [files, setFiles] = useState<File[]>([]);
  const [existingTask, setExistingTask] = useState<Task | null>(null);
  const [existingAttachments, setExistingAttachments] = useState<any[]>([]);
  const [deletedFileIds, setDeletedFileIds] = useState<number[]>([]);

  const [formData, setFormData] = useState({
    project_id: "",
    title: "",
    description: "",
    priority_id: "",
    stack_id: "",
    assigned_id: "",
    due_date: "",
  });

  useEffect(() => {
    loadProjects();
    loadOptions();
    if (isEdit && code) {
      loadTaskData();
    }
  }, [isEdit, code]);

  const loadProjects = async () => {
    try {
      const response = await projectService.getAll({ limit: 100 });
      setProjects(response.data);
    } catch (err) {
      console.error("Failed to load projects:", err);
    }
  };

  const loadOptions = async () => {
    try {
      // Load priorities from reference_category_id = 2 (Low, Medium, High)
      const prioritiesResponse = await referenceService.getByCategoryId(2);
      setPriorities(prioritiesResponse);

      // Load stacks from reference_category_id = 4 (Backend, Web, Mobile, DevOps, Database)
      const stacksResponse = await referenceService.getByCategoryId(4);
      setStacks(stacksResponse);

    } catch (err) {
      console.error("Failed to load options:", err);
    }
  };

  useEffect(() => {
    if (formData.project_id) {
      loadUsersByProjectId(parseInt(formData.project_id));
    } else {
      setUsers([]);
    }
  }, [formData.project_id]);

  const loadUsersByProjectId = async (projectId: number) => {
    try {
      const usersResponse = await projectService.getByProjectId(projectId);
      setUsers(usersResponse);
    } catch (err) {
      console.error("Failed to load users by project ID:", err);
    }
  };

  const loadTaskData = async () => {
    if (!code) return;

    setLoadingData(true);
    try {
      const task = await taskService.getByCode(code);
      setExistingTask(task);
      
      // Load existing attachments
      if (task.task_files && task.task_files.length > 0) {
        setExistingAttachments(task.task_files);
      }
      
      setFormData({
        project_id: task.project_id.toString(),
        title: task.title,
        description: task.description,
        priority_id: task.priority_id?.toString() || "",
        stack_id: task.stack_id?.toString() || "",
        assigned_id: task.assigned_id?.toString() || "",
        due_date: task.due_date
          ? new Date(task.due_date).toISOString().split("T")[0]
          : "",
      });
    } catch (err: any) {
      setError(err.message || "Failed to load task data");
    } finally {
      setLoadingData(false);
    }
  };

  const handleChange = (
    e: React.ChangeEvent<
      HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement
    >
  ) => {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: value });
  };

  const handleFilesChange = (newFiles: File[]) => {
    setFiles(newFiles);
  };

  const handleDeleteExistingFile = (fileId: number) => {
    setDeletedFileIds([...deletedFileIds, fileId]);
    setExistingAttachments(existingAttachments.filter(f => f.id !== fileId));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    if (!formData.project_id) {
      setError("Please select a project");
      return;
    }

    setLoading(true);
    try {
      if (isEdit && existingTask) {
        const updateData: TaskUpdateRequest = {
          project_id: parseInt(formData.project_id),
          title: formData.title,
          description: formData.description,
          priority_id: parseInt(formData.priority_id),
          stack_id: parseInt(formData.stack_id),
          assigned_id: formData.assigned_id
            ? parseInt(formData.assigned_id)
            : undefined,
          due_date: formData.due_date || undefined,
        };
        await taskService.update(existingTask.id, updateData, files, deletedFileIds);
      } else {
        const createData: TaskCreateRequest = {
          project_id: parseInt(formData.project_id),
          title: formData.title,
          description: formData.description,
          priority_id: parseInt(formData.priority_id),
          stack_id: parseInt(formData.stack_id),
          assigned_id: formData.assigned_id
            ? parseInt(formData.assigned_id)
            : undefined,
          due_date: formData.due_date || undefined,
        };
        await taskService.create(createData, files);
      }

      router.push("/tasks");
    } catch (err: any) {
      setError(err.message || `Failed to ${isEdit ? "update" : "create"} task`);
    } finally {
      setLoading(false);
    }
  };

  if (loadingData) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="w-12 h-12 border-4 border-blue-600 rounded-full border-t-transparent animate-spin"></div>
      </div>
    );
  }

  return (
    <div className="max-w-4xl p-6 mx-auto">
      <div className="mb-6">
        <div className="flex items-center gap-2 mb-2">
          <Link
            href="/tasks"
            className="text-gray-500 hover:text-gray-700 dark:text-gray-400"
          >
            ← Back
          </Link>
        </div>
        <h1 className="text-2xl font-bold text-gray-800 dark:text-white">
          {isEdit ? "Edit Task" : "Create New Task"}
        </h1>
        <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
          {isEdit ? `Edit task ${code}` : "Add a new task to the system"}
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
            <div className="md:col-span-2">
              <Label>
                Project <span className="text-error-500">*</span>
              </Label>
              <select
                name="project_id"
                value={formData.project_id}
                onChange={handleChange}
                required
                className="w-full px-3 py-2 border border-gray-300 rounded-lg dark:bg-gray-800 dark:border-gray-700 dark:text-white"
              >
                <option value="">Select a project</option>
                {projects.map((project) => (
                  <option key={project.id} value={project.id}>
                    {project.name} ({project.code})
                  </option>
                ))}
              </select>
            </div>

            <div className="md:col-span-2">
              <Label>
                Title <span className="text-error-500">*</span>
              </Label>
              <Input
                name="title"
                type="text"
                value={formData.title}
                onChange={handleChange}
                required
                placeholder="Enter task title"
              />
            </div>

            <div className="md:col-span-2">
              <Label>
                Description <span className="text-error-500">*</span>
              </Label>
              <textarea
                name="description"
                value={formData.description}
                onChange={handleChange}
                required
                rows={5}
                placeholder="Enter task description"
                className="w-full px-3 py-2 border border-gray-300 rounded-lg dark:bg-gray-800 dark:border-gray-700 dark:text-white"
              />
            </div>

            <div>
              <Label>
                Priority <span className="text-error-500">*</span>
              </Label>
              <select
                name="priority_id"
                value={formData.priority_id}
                onChange={handleChange}
                required
                className="w-full px-3 py-2 border border-gray-300 rounded-lg dark:bg-gray-800 dark:border-gray-700 dark:text-white"
              >
                <option value="">Select priority</option>
                {priorities.map((priority) => (
                  <option key={priority.id} value={priority.id}>
                    {priority.name}
                  </option>
                ))}
              </select>
            </div>

            <div>
              <Label>
                Stack <span className="text-error-500">*</span>
              </Label>
              <select
                name="stack_id"
                value={formData.stack_id}
                onChange={handleChange}
                required
                className="w-full px-3 py-2 border border-gray-300 rounded-lg dark:bg-gray-800 dark:border-gray-700 dark:text-white"
              >
                <option value="">Select stack</option>
                {stacks.map((stack) => (
                  <option key={stack.id} value={stack.id}>
                    {stack.name}
                  </option>
                ))}
              </select>
            </div>

            <div>
              <Label>Assigned To</Label>
              <select
                name="assigned_id"
                value={formData.assigned_id}
                onChange={handleChange}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg dark:bg-gray-800 dark:border-gray-700 dark:text-white"
              >
                <option value="">Select user (optional)</option>
                {users.map((user) => (
                  <option key={user.id} value={user.id}>
                    {user.full_name}
                  </option>
                ))}
              </select>
            </div>

            <div>
              <Label>Due Date</Label>
              <Input
                name="due_date"
                type="date"
                value={formData.due_date}
                onChange={handleChange}
              />
            </div>

            <div className="md:col-span-2">
              <Label>Attachments</Label>
              
              {/* Existing Attachments */}
              {isEdit && existingAttachments.length > 0 && (
                <div className="mb-4">
                  <div className="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Existing Attachments ({existingAttachments.length})
                  </div>
                  <div className="space-y-2">
                    {existingAttachments.map((file) => (
                      <div
                        key={file.id}
                        className="flex items-center gap-3 p-3 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg"
                      >
                        {/* File Icon */}
                        <div className="flex-shrink-0">
                          <div className="w-12 h-12 flex items-center justify-center text-2xl bg-white dark:bg-gray-900 rounded border border-blue-300 dark:border-blue-600">
                            {file.file_type?.startsWith("image/") ? "🖼️" : 
                             file.file_type?.includes("pdf") ? "📄" : 
                             file.file_type?.includes("word") || file.file_type?.includes("document") ? "📝" : 
                             file.file_type?.includes("sheet") || file.file_type?.includes("excel") ? "📊" : "📎"}
                          </div>
                        </div>

                        {/* File Info */}
                        <div className="flex-1 min-w-0">
                          <div className="text-sm font-medium text-gray-900 dark:text-white truncate">
                            {file.file_name}
                          </div>
                          <div className="text-xs text-gray-500 dark:text-gray-400">
                            {file.file_size || "Unknown size"}
                          </div>
                        </div>

                        {/* Action Buttons */}
                        <div className="flex gap-2">
                          <a
                            href={`${process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"}/${file.file_path}`}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="flex-shrink-0 px-3 py-1.5 text-xs font-medium text-blue-600 hover:text-blue-700 dark:text-blue-400 dark:hover:text-blue-300 bg-blue-100 hover:bg-blue-200 dark:bg-blue-900/30 dark:hover:bg-blue-900/50 rounded-lg transition-colors"
                          >
                            View
                          </a>
                          <button
                            type="button"
                            onClick={() => handleDeleteExistingFile(file.id)}
                            className="flex-shrink-0 px-3 py-1.5 text-xs font-medium text-red-600 hover:text-red-700 dark:text-red-400 dark:hover:text-red-300 bg-red-100 hover:bg-red-200 dark:bg-red-900/30 dark:hover:bg-red-900/50 rounded-lg transition-colors"
                          >
                            Delete
                          </button>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              )}
              
              {/* New File Upload */}
              <div>
                <div className="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  {isEdit ? "Add New Attachments" : "Upload Attachments"}
                </div>
                <FileUpload
                  files={files}
                  onChange={handleFilesChange}
                  accept="image/*,.pdf,.doc,.docx,.xls,.xlsx"
                  maxSize={10}
                  maxFiles={5}
                  multiple={true}
                />
              </div>
            </div>
          </div>

          <div className="flex gap-3 pt-4">
            <Button type="submit" disabled={loading}>
              {loading
                ? isEdit
                  ? "Updating..."
                  : "Creating..."
                : isEdit
                  ? "Update Task"
                  : "Create Task"}
            </Button>
            <Link href="/tasks">
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
