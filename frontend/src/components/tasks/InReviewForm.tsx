"use client";
import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import Button from "@/components/ui/button/Button";
import Label from "@/components/form/Label";
import Link from "next/link";
import { Task, TaskInReviewRequest } from "@/types/task";
import { taskService } from "@/services/task.service";

interface InReviewFormProps {
  code: string;
}

export default function InReviewForm({ code }: InReviewFormProps) {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [loadingData, setLoadingData] = useState(true);
  const [error, setError] = useState("");
  const [task, setTask] = useState<Task | null>(null);
  const [files, setFiles] = useState<File[]>([]);

  const [formData, setFormData] = useState({
    description_before: "",
    description_after: "",
  });

  useEffect(() => {
    loadTaskData();
  }, [code]);

  const loadTaskData = async () => {
    setLoadingData(true);
    try {
      const taskData = await taskService.getByCode(code);
      setTask(taskData);
    } catch (err: any) {
      setError(err.message || "Failed to load task data");
    } finally {
      setLoadingData(false);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: value });
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files) {
      setFiles(Array.from(e.target.files));
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    if (!task) return;

    setLoading(true);
    try {
      const data: TaskInReviewRequest = {
        description_before: formData.description_before,
        description_after: formData.description_after,
      };

      await taskService.inReview(task.id, data, files);
      router.push("/tasks");
    } catch (err: any) {
      setError(err.message || "Failed to submit task for review");
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
          Submit Task for Review
        </h1>
        <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
          Task: {task?.code} - {task?.title}
        </p>
      </div>

      <div className="p-6 bg-white border border-gray-200 rounded-lg dark:bg-gray-dark dark:border-gray-800">
        <form onSubmit={handleSubmit} className="space-y-6">
          {error && (
            <div className="p-3 text-sm border rounded-lg bg-error-50 border-error-200 text-error-700 dark:bg-error-900/20">
              {error}
            </div>
          )}

          <div>
            <Label>
              Description Before <span className="text-error-500">*</span>
            </Label>
            <textarea
              name="description_before"
              value={formData.description_before}
              onChange={handleChange}
              required
              rows={5}
              placeholder="Describe the state before the changes"
              className="w-full px-3 py-2 border border-gray-300 rounded-lg dark:bg-gray-800 dark:border-gray-700 dark:text-white"
            />
          </div>

          <div>
            <Label>
              Description After <span className="text-error-500">*</span>
            </Label>
            <textarea
              name="description_after"
              value={formData.description_after}
              onChange={handleChange}
              required
              rows={5}
              placeholder="Describe the state after the changes"
              className="w-full px-3 py-2 border border-gray-300 rounded-lg dark:bg-gray-800 dark:border-gray-700 dark:text-white"
            />
          </div>

          <div>
            <Label>
              Attachments (Before/After screenshots)
            </Label>
            <input
              type="file"
              multiple
              onChange={handleFileChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg dark:bg-gray-800 dark:border-gray-700 dark:text-white"
              accept="image/*,.pdf"
            />
            {files.length > 0 && (
              <p className="mt-2 text-sm text-gray-500">
                {files.length} file(s) selected
              </p>
            )}
          </div>

          <div className="flex gap-3 pt-4">
            <Button type="submit" disabled={loading}>
              {loading ? "Submitting..." : "Submit for Review"}
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
