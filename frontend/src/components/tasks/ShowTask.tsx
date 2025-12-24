"use client";
import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import Button from "@/components/ui/button/Button";
import Link from "next/link";
import { Task } from "@/types/task";
import { taskService } from "@/services/task.service";
import Image from "next/image";
import {
  FiFile,
  FiDownload,
  FiUser,
  FiClock,
  FiCheckCircle,
  FiXCircle,
  FiAlertCircle,
} from "react-icons/fi";

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

  const getPriorityColor = (priority?: string) => {
    const colors: Record<string, string> = {
      Critical: "text-red-600 dark:text-red-400",
      High: "text-orange-600 dark:text-orange-400",
      Medium: "text-yellow-600 dark:text-yellow-400",
      Low: "text-green-600 dark:text-green-400",
    };
    return colors[priority || ""] || "text-gray-600 dark:text-gray-400";
  };

  const isImageFile = (fileType?: string) => {
    return fileType?.startsWith("image/");
  };

  const formatFileSize = (bytes: string) => {
    const size = parseInt(bytes);
    if (size < 1024) return `${size} B`;
    if (size < 1024 * 1024) return `${(size / 1024).toFixed(2)} KB`;
    return `${(size / (1024 * 1024)).toFixed(2)} MB`;
  };

  const formatDateTime = (date?: string | null) => {
    if (!date) return "-";
    return new Date(date).toLocaleString("id-ID", {
      day: "2-digit",
      month: "short",
      year: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  const getApprovalIcon = (statusName?: string) => {
    if (statusName === "Approve")
      return <FiCheckCircle className="w-5 h-5 text-green-600" />;
    if (statusName === "Reject")
      return <FiXCircle className="w-5 h-5 text-red-600" />;
    return <FiAlertCircle className="w-5 h-5 text-yellow-600" />;
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
    <div className="max-w-7xl p-6 mx-auto">
      {/* Header */}
      <div className="mb-6">
        <div className="flex items-start justify-between mb-4">
          <div className="flex-1">
            <div className="flex items-center gap-2 mb-3">
              <Link
                href="/tasks"
                className="text-sm text-gray-500 hover:text-gray-700 dark:text-gray-400 hover:underline"
              >
                ← Back to Tasks
              </Link>
            </div>
            <div className="flex items-center gap-3 mb-2">
              <h1 className="text-3xl font-bold text-gray-900 dark:text-white">
                {task.code}
              </h1>
              {getStatusBadge(task.status_task?.name)}
            </div>
            <h2 className="text-xl text-gray-700 dark:text-gray-300 mb-3">
              {task.title}
            </h2>
            <div className="flex items-center gap-4 text-sm text-gray-600 dark:text-gray-400">
              <span className="flex items-center gap-1">
                <FiUser className="w-4 h-4" />
                Created by:{" "}
                <span className="font-medium">
                  {task.creator?.full_name || "-"}
                </span>
              </span>
              <span className="flex items-center gap-1">
                <FiClock className="w-4 h-4" />
                {formatDateTime(task.created_at)}
              </span>
            </div>
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

          {/* Files with Preview */}
          {task.task_files && task.task_files.length > 0 && (
            <div className="p-6 bg-white border border-gray-200 rounded-lg dark:bg-gray-dark dark:border-gray-800">
              <h2 className="mb-4 text-lg font-semibold text-gray-800 dark:text-white flex items-center gap-2">
                <FiFile className="w-5 h-5" />
                Attachments ({task.task_files.length})
              </h2>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {task.task_files.map((file) => {
                  const fileUrl = `${process.env.NEXT_PUBLIC_STORAGE_URL}/${file.file_path}`;
                  const isImage = isImageFile(file.file_type);

                  return (
                    <div
                      key={file.id}
                      className="border border-gray-200 rounded-lg overflow-hidden hover:shadow-lg transition-shadow dark:border-gray-700"
                    >
                      {/* Preview */}
                      {isImage ? (
                        <div className="relative h-48 bg-gray-100 dark:bg-gray-800">
                          <Image
                            src={fileUrl}
                            alt={file.file_name}
                            fill
                            className="object-contain"
                            unoptimized
                          />
                        </div>
                      ) : (
                        <div className="h-48 bg-gray-100 dark:bg-gray-800 flex items-center justify-center">
                          <FiFile className="w-16 h-16 text-gray-400" />
                        </div>
                      )}

                      {/* File Info */}
                      <div className="p-4">
                        <p className="font-medium text-gray-800 dark:text-white truncate mb-1">
                          {file.file_name}
                        </p>
                        <div className="flex items-center justify-between">
                          <span className="text-sm text-gray-500 dark:text-gray-400">
                            {file.file_size
                              ? formatFileSize(file.file_size)
                              : "-"}
                          </span>
                          <a
                            href={fileUrl}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="flex items-center gap-1 text-sm text-blue-600 hover:text-blue-800 dark:text-blue-400"
                          >
                            <FiDownload className="w-4 h-4" />
                            Download
                          </a>
                        </div>
                      </div>
                    </div>
                  );
                })}
              </div>
            </div>
          )}

          {/* Approval History */}
          {task.approval_tasks && task.approval_tasks.length > 0 && (
            <div className="p-6 bg-white border border-gray-200 rounded-lg dark:bg-gray-dark dark:border-gray-800">
              <h2 className="mb-6 text-lg font-semibold text-gray-800 dark:text-white">
                Approval History
              </h2>
              <div className="relative space-y-6">
                {/* Timeline line */}
                <div className="absolute left-5 top-0 bottom-0 w-0.5 bg-gray-200 dark:bg-gray-700"></div>

                {task.approval_tasks
                  .sort((a, b) => a.sequence - b.sequence)
                  .map((approval) => (
                    <div key={approval.id} className="relative flex gap-4">
                      {/* Icon */}
                      <div className="relative z-10 flex items-center justify-center w-10 h-10 bg-white dark:bg-gray-dark rounded-full border-2 border-gray-200 dark:border-gray-700">
                        {getApprovalIcon(approval.approval_status?.name)}
                      </div>

                      {/* Content */}
                      <div className="flex-1 pb-6">
                        <div className="p-4 border border-gray-200 rounded-lg dark:border-gray-700 bg-gray-50 dark:bg-gray-800/50">
                          <div className="flex items-start justify-between mb-3">
                            <div className="flex-1">
                              <div className="flex items-center gap-2 mb-2">
                                <span className="text-sm font-semibold text-gray-500 dark:text-gray-400">
                                  Level {approval.sequence}
                                </span>
                                {getApprovalBadge(
                                  approval.approval_status?.name
                                )}
                              </div>

                              {/* Approver Info */}
                              {(approval.approved_by_user || approval.approver) && (
                                <div className="space-y-1 mb-3">
                                  <div className="flex items-center gap-2">
                                    <FiUser className="w-4 h-4 text-gray-500" />
                                    <span className="font-medium text-gray-900 dark:text-white">
                                      {(approval.approved_by_user || approval.approver)?.full_name || 'Unknown'}
                                    </span>
                                  </div>
                                  <div className="text-sm text-gray-600 dark:text-gray-400 ml-6">
                                    {(approval.approved_by_user || approval.approver)?.email || '-'}
                                  </div>
                                  <div className="text-sm text-gray-500 dark:text-gray-500 ml-6">
                                    @{(approval.approved_by_user || approval.approver)?.username || '-'}
                                  </div>
                                </div>
                              )}

                              {/* Date */}
                              {approval.approval_date && (
                                <div className="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400">
                                  <FiClock className="w-4 h-4" />
                                  {formatDateTime(approval.approval_date)}
                                </div>
                              )}
                            </div>
                          </div>

                          {/* Note */}
                          {approval.note && (
                            <div className="mt-3 p-3 bg-white dark:bg-gray-900 rounded border border-gray-200 dark:border-gray-700">
                              <p className="text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                                Note:
                              </p>
                              <p className="text-sm text-gray-600 dark:text-gray-400">
                                {approval.note}
                              </p>
                            </div>
                          )}
                        </div>
                      </div>
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
            <dl className="space-y-4">
              <div>
                <dt className="text-xs font-medium uppercase text-gray-500 dark:text-gray-400 mb-1">
                  Project
                </dt>
                <dd className="text-sm font-medium text-gray-900 dark:text-white">
                  {task.project?.name || "-"}
                </dd>
              </div>

              <div className="pt-3 border-t border-gray-200 dark:border-gray-700">
                <dt className="text-xs font-medium uppercase text-gray-500 dark:text-gray-400 mb-2">
                  Status
                </dt>
                <dd>{getStatusBadge(task.status_task?.name)}</dd>
              </div>

              <div className="pt-3 border-t border-gray-200 dark:border-gray-700">
                <dt className="text-xs font-medium uppercase text-gray-500 dark:text-gray-400 mb-2">
                  Approval Status
                </dt>
                <dd>{getApprovalBadge(task.approval_status?.name)}</dd>
              </div>

              <div className="pt-3 border-t border-gray-200 dark:border-gray-700">
                <dt className="text-xs font-medium uppercase text-gray-500 dark:text-gray-400 mb-1">
                  Priority
                </dt>
                <dd
                  className={`text-sm font-semibold ${getPriorityColor(task.priority?.name)}`}
                >
                  {task.priority?.name || "-"}
                </dd>
              </div>

              <div className="pt-3 border-t border-gray-200 dark:border-gray-700">
                <dt className="text-xs font-medium uppercase text-gray-500 dark:text-gray-400 mb-1">
                  Stack
                </dt>
                <dd className="text-sm text-gray-900 dark:text-white">
                  {task.stack?.name || "-"}
                </dd>
              </div>

              {task.type?.name && (
                <div className="pt-3 border-t border-gray-200 dark:border-gray-700">
                  <dt className="text-xs font-medium uppercase text-gray-500 dark:text-gray-400 mb-1">
                    Type
                  </dt>
                  <dd className="text-sm text-gray-900 dark:text-white">
                    {task.type.name}
                  </dd>
                </div>
              )}
            </dl>
          </div>

          {/* People */}
          <div className="p-6 bg-white border border-gray-200 rounded-lg dark:bg-gray-dark dark:border-gray-800">
            <h2 className="mb-4 text-lg font-semibold text-gray-800 dark:text-white">
              People
            </h2>
            <dl className="space-y-4">
              <div>
                <dt className="text-xs font-medium uppercase text-gray-500 dark:text-gray-400 mb-2">
                  Assigned To
                </dt>
                <dd className="flex items-center gap-2">
                  <div className="w-8 h-8 bg-blue-100 dark:bg-blue-900 rounded-full flex items-center justify-center">
                    <span className="text-sm font-semibold text-blue-600 dark:text-blue-300">
                      {task.assignee?.full_name?.charAt(0) || "?"}
                    </span>
                  </div>
                  <div>
                    <p className="text-sm font-medium text-gray-900 dark:text-white">
                      {task.assignee?.full_name || "-"}
                    </p>
                    <p className="text-xs text-gray-500 dark:text-gray-400">
                      {task.assignee?.email || ""}
                    </p>
                  </div>
                </dd>
              </div>

              <div className="pt-3 border-t border-gray-200 dark:border-gray-700">
                <dt className="text-xs font-medium uppercase text-gray-500 dark:text-gray-400 mb-2">
                  Created By
                </dt>
                <dd className="flex items-center gap-2">
                  <div className="w-8 h-8 bg-purple-100 dark:bg-purple-900 rounded-full flex items-center justify-center">
                    <span className="text-sm font-semibold text-purple-600 dark:text-purple-300">
                      {task.creator?.full_name?.charAt(0) || "?"}
                    </span>
                  </div>
                  <div>
                    <p className="text-sm font-medium text-gray-900 dark:text-white">
                      {task.creator?.full_name || "-"}
                    </p>
                    <p className="text-xs text-gray-500 dark:text-gray-400">
                      {task.creator?.email || ""}
                    </p>
                  </div>
                </dd>
              </div>
            </dl>
          </div>

          {/* Dates */}
          <div className="p-6 bg-white border border-gray-200 rounded-lg dark:bg-gray-dark dark:border-gray-800">
            <h2 className="mb-4 text-lg font-semibold text-gray-800 dark:text-white">
              Timeline
            </h2>
            <dl className="space-y-3">
              <div>
                <dt className="text-xs font-medium uppercase text-gray-500 dark:text-gray-400 mb-1">
                  Created At
                </dt>
                <dd className="text-sm text-gray-900 dark:text-white">
                  {formatDateTime(task.created_at)}
                </dd>
              </div>

              {task.due_date && (
                <div className="pt-3 border-t border-gray-200 dark:border-gray-700">
                  <dt className="text-xs font-medium uppercase text-gray-500 dark:text-gray-400 mb-1">
                    Due Date
                  </dt>
                  <dd className="text-sm font-medium text-red-600 dark:text-red-400">
                    {formatDateTime(task.due_date)}
                  </dd>
                </div>
              )}

              {task.approval_date && (
                <div className="pt-3 border-t border-gray-200 dark:border-gray-700">
                  <dt className="text-xs font-medium uppercase text-gray-500 dark:text-gray-400 mb-1">
                    Approval Date
                  </dt>
                  <dd className="text-sm text-gray-900 dark:text-white">
                    {formatDateTime(task.approval_date)}
                  </dd>
                </div>
              )}

              {task.updated_at && task.updated_at !== task.created_at && (
                <div className="pt-3 border-t border-gray-200 dark:border-gray-700">
                  <dt className="text-xs font-medium uppercase text-gray-500 dark:text-gray-400 mb-1">
                    Last Updated
                  </dt>
                  <dd className="text-sm text-gray-900 dark:text-white">
                    {formatDateTime(task.updated_at)}
                  </dd>
                </div>
              )}
            </dl>
          </div>
        </div>
      </div>
    </div>
  );
}
