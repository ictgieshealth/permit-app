"use client";
import { useState, useRef, useEffect } from "react";
import { Task } from "@/types/task";

interface TaskActionsDropdownProps {
  task: Task;
  index: number;
  totalRows: number;
  userRole: string | null;
  onTodo?: (task: Task) => void;
  onStart?: (task: Task) => void;
  onHold?: (task: Task) => void;
  onDone?: (task: Task) => void;
  onInReview?: (task: Task) => void;
  onShow?: (task: Task) => void;
  onAddReason?: (task: Task) => void;
  onRevision?: (task: Task) => void;
  onEdit?: (task: Task) => void;
  onDelete?: (task: Task) => void;
  onDevelopment?: (task: Task) => void;
  onMaintenance?: (task: Task) => void;
}

export default function TaskActionsDropdown({
  task,
  index,
  totalRows,
  userRole,
  onTodo,
  onStart,
  onHold,
  onDone,
  onInReview,
  onShow,
  onAddReason,
  onRevision,
  onEdit,
  onDelete,
  onDevelopment,
  onMaintenance,
}: TaskActionsDropdownProps) {
  const [isOpen, setIsOpen] = useState(false);
  const [position, setPosition] = useState({ top: 0, right: 0 });
  const dropdownRef = useRef<HTMLDivElement>(null);
  const buttonRef = useRef<HTMLButtonElement>(null);
  const isNearBottom = index >= totalRows - 2;

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        dropdownRef.current &&
        !dropdownRef.current.contains(event.target as Node) &&
        buttonRef.current &&
        !buttonRef.current.contains(event.target as Node)
      ) {
        setIsOpen(false);
      }
    };

    if (isOpen) {
      document.addEventListener("mousedown", handleClickOutside);
    }

    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, [isOpen]);

  useEffect(() => {
    if (isOpen && buttonRef.current) {
      const rect = buttonRef.current.getBoundingClientRect();
      setPosition({
        top: isNearBottom ? rect.bottom - 200 : rect.top,
        right: window.innerWidth - rect.left + 8,
      });
    }
  }, [isOpen, isNearBottom]);

  const getActionVisibility = () => {
    const isAdmin = userRole === "Admin";
    const isDev = userRole === "Ticketing Developer";
    const isManager = userRole === "Ticketing Manager";
    const isPIC = userRole === "Ticketing PIC";
    const isHeadOfUnit = userRole === "Ticketing Head of Unit";
    const isClient = !isAdmin && !isDev && !isManager && !isPIC && !isHeadOfUnit;

    const statusId = task.status_id;
    const typeName = task.type?.name;

    // Admin and Developer Actions
    if (isAdmin || isDev) {
      return {
        showTodo: statusId === 2, // From On Hold
        showStart: statusId === 1, // From To Do
        showOnHold: [1, 3, 37].includes(statusId || 0), // From To Do, On Progress, In Review
        showDone: statusId === 37, // From In Review
        showInReview: statusId === 3, // From On Progress
        showShow: true,
        showAddReason: statusId === 2, // On Hold
        showRevision: statusId === 37, // In Review
        showEdit: ![4, 37].includes(statusId || 0), // Not Done or In Review
        showDelete: statusId === 1, // To Do only
        showDevelopment: typeName === "Maintenance",
        showMaintenance: typeName === "Development",
      };
    }

    // Client Actions (PIC, Manager, Head of Unit)
    if (isClient || isPIC || isManager || isHeadOfUnit) {
      return {
        showTodo: statusId === 2, // From On Hold
        showStart: statusId === 1, // From To Do
        showOnHold: [1, 3, 37].includes(statusId || 0), // From To Do, On Progress, In Review
        showDone: statusId === 37, // From In Review
        showInReview: false,
        showShow: true,
        showAddReason: false,
        showRevision: statusId === 37, // In Review
        showEdit: ![4].includes(statusId || 0), // Not Done
        showDelete: false,
        showDevelopment: false,
        showMaintenance: false,
      };
    }

    return {
      showTodo: false,
      showStart: false,
      showOnHold: false,
      showDone: false,
      showInReview: false,
      showShow: false,
      showAddReason: false,
      showRevision: false,
      showEdit: false,
      showDelete: false,
      showDevelopment: false,
      showMaintenance: false,
    };
  };

  const visibility = getActionVisibility();

  const hasStatusActions =
    visibility.showTodo ||
    visibility.showStart ||
    visibility.showOnHold ||
    visibility.showDone ||
    visibility.showInReview;

  const hasSettingActions =
    visibility.showShow ||
    visibility.showAddReason ||
    visibility.showRevision ||
    visibility.showEdit ||
    visibility.showDelete;

  const hasTypeActions = visibility.showDevelopment || visibility.showMaintenance;

  if (!hasStatusActions && !hasSettingActions && !hasTypeActions) {
    return null;
  }

  const handleAction = (action: () => void) => {
    setIsOpen(false);
    action();
  };

  return (
    <div className="relative inline-block text-left">
      <button
        ref={buttonRef}
        onClick={() => setIsOpen(!isOpen)}
        className="inline-flex items-center justify-center w-8 h-8 text-gray-600 transition-colors rounded-lg hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-700"
        title="Actions"
      >
        <svg
          className="w-5 h-5"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z"
          />
        </svg>
      </button>

      {isOpen && (
        <div
          ref={dropdownRef}
          className="fixed w-56 bg-white dark:bg-gray-800 rounded-lg shadow-lg border border-gray-200 dark:border-gray-700 z-50"
          style={{
            top: `${position.top}px`,
            right: `${position.right}px`,
          }}
        >
          {/* Status Section */}
          {hasStatusActions && (
            <>
              <div className="px-4 py-2 text-xs font-semibold text-gray-500 uppercase bg-gray-50 dark:bg-gray-900 dark:text-gray-400 rounded-t-lg">
                Status
              </div>
              {visibility.showTodo && onTodo && (
                <button
                  onClick={() => handleAction(() => onTodo(task))}
                  className="w-full px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700 flex items-center gap-2"
                >
                  <span className="text-gray-600">☑️</span>
                  To Do
                </button>
              )}
              {visibility.showStart && onStart && (
                <button
                  onClick={() => handleAction(() => onStart(task))}
                  className="w-full px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700 flex items-center gap-2"
                >
                  <span className="text-green-600">▶️</span>
                  Start Task
                </button>
              )}
              {visibility.showOnHold && onHold && (
                <button
                  onClick={() => handleAction(() => onHold(task))}
                  className="w-full px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700 flex items-center gap-2"
                >
                  <span className="text-yellow-600">⏸️</span>
                  On Hold Task
                </button>
              )}
              {visibility.showInReview && onInReview && (
                <button
                  onClick={() => handleAction(() => onInReview(task))}
                  className="w-full px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700 flex items-center gap-2"
                >
                  <span className="text-blue-600">🔍</span>
                  In Review
                </button>
              )}
              {visibility.showDone && onDone && (
                <button
                  onClick={() => handleAction(() => onDone(task))}
                  className="w-full px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700 flex items-center gap-2"
                >
                  <span className="text-green-600">✅</span>
                  Done
                </button>
              )}
            </>
          )}

          {/* Setting Section */}
          {hasSettingActions && (
            <>
              {hasStatusActions && <hr className="my-1 border-gray-200 dark:border-gray-700" />}
              <div className="px-4 py-2 text-xs font-semibold text-gray-500 uppercase bg-gray-50 dark:bg-gray-900 dark:text-gray-400">
                Setting
              </div>
              {visibility.showShow && onShow && (
                <button
                  onClick={() => handleAction(() => onShow(task))}
                  className="w-full px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700 flex items-center gap-2"
                >
                  <span className="text-blue-600">👁️</span>
                  Show
                </button>
              )}
              {visibility.showAddReason && onAddReason && (
                <button
                  onClick={() => handleAction(() => onAddReason(task))}
                  className="w-full px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700 flex items-center gap-2"
                >
                  <span className="text-purple-600">➕</span>
                  Add Reason
                </button>
              )}
              {visibility.showRevision && onRevision && (
                <button
                  onClick={() => handleAction(() => onRevision(task))}
                  className="w-full px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700 flex items-center gap-2"
                >
                  <span className="text-orange-600">🔄</span>
                  Revision
                </button>
              )}
              {visibility.showEdit && onEdit && (
                <button
                  onClick={() => handleAction(() => onEdit(task))}
                  className="w-full px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700 flex items-center gap-2"
                >
                  <span className="text-indigo-600">✏️</span>
                  Edit
                </button>
              )}
              {visibility.showDelete && onDelete && (
                <button
                  onClick={() => handleAction(() => onDelete(task))}
                  className="w-full px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700 flex items-center gap-2"
                >
                  <span className="text-red-600">🗑️</span>
                  Delete
                </button>
              )}
            </>
          )}

          {/* Type Section */}
          {hasTypeActions && (
            <>
              <hr className="my-1 border-gray-200 dark:border-gray-700" />
              <div className="px-4 py-2 text-xs font-semibold text-gray-500 uppercase bg-gray-50 dark:bg-gray-900 dark:text-gray-400">
                Type
              </div>
              {visibility.showDevelopment && onDevelopment && (
                <button
                  onClick={() => handleAction(() => onDevelopment(task))}
                  className="w-full px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700 flex items-center gap-2"
                >
                  <span className="text-blue-600">💻</span>
                  Development
                </button>
              )}
              {visibility.showMaintenance && onMaintenance && (
                <button
                  onClick={() => handleAction(() => onMaintenance(task))}
                  className="w-full px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700 flex items-center gap-2 rounded-b-lg"
                >
                  <span className="text-yellow-600">🔧</span>
                  Maintenance
                </button>
              )}
            </>
          )}
        </div>
      )}
    </div>
  );
}
