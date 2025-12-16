"use client";
import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import Button from "@/components/ui/button/Button";
import Link from "next/link";
import { Task } from "@/types/task";
import { taskService } from "@/services/task.service";

interface ShowTaskProps {
  code: string;
}

export default function ShowTask({ code }: ShowTaskProps) {
  const router = useRouter();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [task, setTask] = useState<Task | null>(null);

  useEffect(() => {
    loadTaskData();
  }, [code]);

  const loadTaskData = async () => {
    setLoading(true);
    try {
      const taskData = await taskService.getByCode(code);
      setTask(taskData);
    } catch (err: any) {
      setError(err.message || "Failed to load task data");
    } finally {
      setLoading(false);
    }
  };

  const getStatusBadge = (statusName?: string) => {
    const colors: Record<string, string> = {
      "To Do": "bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-300",
      "On Hold":
        "bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300",
      "On Progress":
        "bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300",
      Done: "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300",
      "In Review":
        "bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-300",
      Revision:
        "bg-orange-100 text-orange-800 dark:bg-orange-900 dark:text-orange-300",
    };
    const colorClass =
      colors[statusName || ""] ||
      "bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-300";
    return (
      <span
        className={`px-3 py-1 text-sm font-medium rounded-full ${colorClass}`}
      >
        {statusName || "-"}
      </span>
    );
  };

  const getApprovalBadge = (approvalName?: string) => {
    const colors: Record<string, string> = {
      Waiting:
        "bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300",
      Reject: "bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300",
      Approve:
        "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300",
      "Pending Manager":
        "bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300",
    };
    const colorClass =
      colors[approvalName || ""] ||
      "bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-300";
    return (
      <span
        className={`px-3 py-1 text-sm font-medium rounded-full ${colorClass}`}
      >
        {approvalName || "-"}
      </span>
    );
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="w-12 h-12 border-4 border-blue-600 rounded-full border-t-transparent animate-spin"></div>
      </div>
    );
  }

  if (error || !task) {
    return (
      <div className="p-6">
        <div className="p-4 border rounded-lg bg-error-50 border-error-200 text-error-700">
          {error || "Task not found"}
        </div>
        <Link href="/tasks">
          <Button className="mt-4">Back to Tasks</Button>
        </Link>
      </div>
    );
  }

  return (
    <div className="max-w-6xl p-6 mx-auto">
      <div className="mb-6">
        <div className="flex items-center justify-between mb-4">
          <div>
            <div className="flex items-center gap-2 mb-2">
              <Link
                href="/tasks"
                className="text-gray-500 hover:text-gray-700 dark:text-gray-400"
              >
                ← Back to Tasks
              </Link>
            </div>
            <h1 className="text-2xl font-bold text-gray-800 dark:text-white">
              {task.code}
            </h1>
            <p className="mt-1 text-lg text-gray-600 dark:text-gray-300">
              {task.title}
            </p>
          </div>
          <div className="flex gap-2">
            <Link href={`/tasks/edit/${task.code}`}>
              <Button variant="outline">Edit Task</Button>
            </Link>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
        {/* Main Content */}
        <div className="lg:col-span-2 space-y-6">
          {/* Description */}
          <div className="p-6 bg-white border border-gray-200 rounded-lg dark:bg-gray-dark dark:border-gray-800">
            <h2 className="mb-4 text-lg font-semibold text-gray-800 dark:text-white">
              Description
            </h2>
            <div className="text-gray-600 whitespace-pre-wrap dark:text-gray-300">
              {task.description}
            </div>
          </div>

          {/* Before/After if exists */}
          {(task.description_before || task.description_after) && (
            <div className="p-6 bg-white border border-gray-200 rounded-lg dark:bg-gray-dark dark:border-gray-800">
              <h2 className="mb-4 text-lg font-semibold text-gray-800 dark:text-white">
                Before & After
              </h2>
              {task.description_before && (
                <div className="mb-4">
                  <h3 className="mb-2 font-medium text-gray-700 dark:text-gray-300">
                    Before:
                  </h3>
                  <div className="p-3 text-gray-600 bg-gray-50 rounded dark:bg-gray-800 dark:text-gray-300">
                    {task.description_before}
                  </div>
                </div>
              )}
              {task.description_after && (
                <div>
                  <h3 className="mb-2 font-medium text-gray-700 dark:text-gray-300">
                    After:
                  </h3>
                  <div className="p-3 text-gray-600 bg-gray-50 rounded dark:bg-gray-800 dark:text-gray-300">
                    {task.description_after}
                  </div>
                </div>
              )}
            </div>
          )}

          {/* Reason/Revision if exists */}
          {(task.reason || task.revision) && (
            <div className="p-6 bg-white border border-gray-200 rounded-lg dark:bg-gray-dark dark:border-gray-800">
              {task.reason && (
                <div className="mb-4">
                  <h3 className="mb-2 font-medium text-red-700 dark:text-red-400">
                    Rejection Reason:
                  </h3>
                  <div className="p-3 bg-red-50 rounded dark:bg-red-900/20 text-red-700 dark:text-red-300">
                    {task.reason}
                  </div>
                </div>
              )}
              {task.revision && (
                <div>
                  <h3 className="mb-2 font-medium text-blue-700 dark:text-blue-400">
                    Revision Notes:
                  </h3>
                  <div className="p-3 bg-blue-50 rounded dark:bg-blue-900/20 text-blue-700 dark:text-blue-300">
                    {task.revision}
                  </div>
                </div>
              )}
            </div>
          )}

          {/* Files */}
          {task.task_files && task.task_files.length > 0 && (
            <div className="p-6 bg-white border border-gray-200 rounded-lg dark:bg-gray-dark dark:border-gray-800">
              <h2 className="mb-4 text-lg font-semibold text-gray-800 dark:text-white">
                Attachments
              </h2>
              <div className="space-y-2">
                {task.task_files.map((file) => (
                  <div
                    key={file.id}
                    className="flex items-center justify-between p-3 bg-gray-50 rounded dark:bg-gray-800"
                  >
                    <div>
                      <p className="font-medium text-gray-800 dark:text-white">
                        {file.file_name}
                      </p>
                      <p className="text-sm text-gray-500 dark:text-gray-400">
                        {file.file_size
                          ? `${(parseInt(file.file_size) / 1024).toFixed(2)} KB`
                          : ""}
                      </p>
                    </div>
                    <a
                      href={`${process.env.NEXT_PUBLIC_STORAGE_URL}/${file.file_path}`}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="text-blue-600 hover:text-blue-800 dark:text-blue-400"
                    >
                      Download
                    </a>
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Approval History */}
          {task.approval_tasks && task.approval_tasks.length > 0 && (
            <div className="p-6 bg-white border border-gray-200 rounded-lg dark:bg-gray-dark dark:border-gray-800">
              <h2 className="mb-4 text-lg font-semibold text-gray-800 dark:text-white">
                Approval History
              </h2>
              <div className="space-y-3">
                {task.approval_tasks.map((approval) => (
                  <div
                    key={approval.id}
                    className="p-4 border border-gray-200 rounded-lg dark:border-gray-700"
                  >
                    <div className="flex items-center justify-between mb-2">
                      <span className="font-medium text-gray-800 dark:text-white">
                        Sequence {approval.sequence}
                      </span>
                      {getApprovalBadge(approval.approval_status?.name)}
                    </div>
                    {approval.approver && (
                      <p className="text-sm text-gray-600 dark:text-gray-400">
                        Approver: {approval.approver.full_name}
                      </p>
                    )}
                    {approval.approval_date && (
                      <p className="text-sm text-gray-500 dark:text-gray-500">
                        Date:{" "}
                        {new Date(approval.approval_date).toLocaleString()}
                      </p>
                    )}
                    {approval.note && (
                      <p className="mt-2 text-sm text-gray-600 dark:text-gray-300">
                        Note: {approval.note}
                      </p>
                    )}
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>

        {/* Sidebar */}
        <div className="space-y-6">
          {/* Task Info */}
          <div className="p-6 bg-white border border-gray-200 rounded-lg dark:bg-gray-dark dark:border-gray-800">
            <h2 className="mb-4 text-lg font-semibold text-gray-800 dark:text-white">
              Task Information
            </h2>
            <dl className="space-y-3">
              <div>
                <dt className="text-sm font-medium text-gray-500 dark:text-gray-400">
                  Project
                </dt>
                <dd className="mt-1 text-sm text-gray-900 dark:text-white">
                  {task.project?.name || "-"}
                </dd>
              </div>
              <div>
                <dt className="text-sm font-medium text-gray-500 dark:text-gray-400">
                  Status
                </dt>
                <dd className="mt-1">{getStatusBadge(task.status_task?.name)}</dd>
              </div>
              <div>
                <dt className="text-sm font-medium text-gray-500 dark:text-gray-400">
                  Approval Status
                </dt>
                <dd className="mt-1">
                  {getApprovalBadge(task.approval_status?.name)}
                </dd>
              </div>
              <div>
                <dt className="text-sm font-medium text-gray-500 dark:text-gray-400">
                  Priority
                </dt>
                <dd className="mt-1 text-sm text-gray-900 dark:text-white">
                  {task.priority?.name || "-"}
                </dd>
              </div>
              <div>
                <dt className="text-sm font-medium text-gray-500 dark:text-gray-400">
                  Stack
                </dt>
                <dd className="mt-1 text-sm text-gray-900 dark:text-white">
                  {task.stack?.name || "-"}
                </dd>
              </div>
              <div>
                <dt className="text-sm font-medium text-gray-500 dark:text-gray-400">
                  Type
                </dt>
                <dd className="mt-1 text-sm text-gray-900 dark:text-white">
                  {task.type?.name || "-"}
                </dd>
              </div>
              <div>
                <dt className="text-sm font-medium text-gray-500 dark:text-gray-400">
                  Assigned To
                </dt>
                <dd className="mt-1 text-sm text-gray-900 dark:text-white">
                  {task.assignee?.full_name || "-"}
                </dd>
              </div>
              <div>
                <dt className="text-sm font-medium text-gray-500 dark:text-gray-400">
                  Created By
                </dt>
                <dd className="mt-1 text-sm text-gray-900 dark:text-white">
                  {task.creator?.full_name || "-"}
                </dd>
              </div>
              <div>
                <dt className="text-sm font-medium text-gray-500 dark:text-gray-400">
                  Due Date
                </dt>
                <dd className="mt-1 text-sm text-gray-900 dark:text-white">
                  {task.due_date
                    ? new Date(task.due_date).toLocaleDateString()
                    : "-"}
                </dd>
              </div>
              <div>
                <dt className="text-sm font-medium text-gray-500 dark:text-gray-400">
                  Created At
                </dt>
                <dd className="mt-1 text-sm text-gray-900 dark:text-white">
                  {new Date(task.created_at).toLocaleString()}
                </dd>
              </div>
            </dl>
          </div>
        </div>
      </div>
    </div>
  );
}
