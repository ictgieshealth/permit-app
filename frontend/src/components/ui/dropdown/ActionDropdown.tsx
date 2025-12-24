"use client";
import { useState, useRef, useEffect } from "react";
import {
  EyeIcon,
  PencilIcon,
  TrashBinIcon,
  CheckCircleIcon,
  XCircleIcon,
  CodeIcon,
  ToolsIcon,
} from "@/icons";

interface ActionDropdownProps {
  data: any;
  index: number;
  totalRows: number;
  onView?: (data: any) => void;
  onEdit?: (data: any) => void;
  onDelete?: (data: any) => void;
  onApprove?: (data: any) => void;
  onReject?: (data: any) => void;
  onDevelopment?: (data: any) => void;
  onMaintenance?: (data: any) => void;
  hideView?: boolean;
  hideEdit?: boolean;
  hideDelete?: boolean;
  hideApprove?: boolean;
  hideReject?: boolean;
  hideDevelopment?: boolean;
  hideMaintenance?: boolean;
}

export default function ActionDropdown({
  data,
  index,
  totalRows,
  onView,
  onEdit,
  onDelete,
  onApprove,
  onReject,
  onDevelopment,
  onMaintenance,
  hideView = true,
  hideEdit = true,
  hideDelete = true,
  hideApprove = true,
  hideReject = true,
  hideDevelopment = true,
  hideMaintenance = true,
}: ActionDropdownProps) {
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

  const handleAction = (action: () => void) => {
    action();
    setIsOpen(false);
  };

  const hasTypeActions = !hideDevelopment || !hideMaintenance;

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
          {/* Actions Section */}
          {(!hideView || !hideEdit || !hideDelete) && (
            <>
              <div className="px-4 py-2 text-xs font-semibold text-gray-500 uppercase bg-gray-50 dark:bg-gray-900 dark:text-gray-400 rounded-t-lg">
                Actions
              </div>
              {/* View */}
              {!hideView && onView && (
                <button
                  onClick={() => handleAction(() => onView(data))}
                  className="flex items-center w-full px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700 gap-2"
                >
                  <EyeIcon className="w-4 h-4" />
                  <span>View</span>
                </button>
              )}

              {/* Edit */}
              {!hideEdit && onEdit && (
                <button
                  onClick={() => handleAction(() => onEdit(data))}
                  className="flex items-center w-full px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700 gap-2"
                >
                  <PencilIcon className="w-4 h-4" />
                  <span>Edit</span>
                </button>
              )}

              {/* Delete */}
              {!hideDelete && onDelete && (
                <button
                  onClick={() => handleAction(() => onDelete(data))}
                  className="flex items-center w-full px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700 gap-2"
                >
                  <TrashBinIcon className="w-4 h-4 text-red-600" />
                  <span>Delete</span>
                </button>
              )}
            </>
          )}

          {/* Approval Section */}
          {(!hideApprove || !hideReject) && (
            <>
              {(!hideView || !hideEdit || !hideDelete) && <hr className="my-1 border-gray-200 dark:border-gray-700" />}
              <div className="px-4 py-2 text-xs font-semibold text-gray-500 uppercase bg-gray-50 dark:bg-gray-900 dark:text-gray-400">
                Approval
              </div>
              {/* Approve */}
              {!hideApprove && onApprove && (
                <button
                  onClick={() => handleAction(() => onApprove(data))}
                  className="flex items-center w-full px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700 gap-2"
                >
                  <CheckCircleIcon className="w-4 h-4 text-green-600" />
                  <span>Approve</span>
                </button>
              )}

              {/* Reject */}
              {!hideReject && onReject && (
                <button
                  onClick={() => handleAction(() => onReject(data))}
                  className="flex items-center w-full px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700 gap-2"
                >
                  <XCircleIcon className="w-4 h-4 text-red-600" />
                  <span>Reject</span>
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
              {/* Development */}
              {!hideDevelopment && onDevelopment && (
                <button
                  onClick={() => handleAction(() => onDevelopment(data))}
                  className="flex items-center w-full px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700 gap-2"
                >
                  <CodeIcon className="w-4 h-4 text-blue-600" />
                  <span>Development</span>
                </button>
              )}

              {/* Maintenance */}
              {!hideMaintenance && onMaintenance && (
                <button
                  onClick={() => handleAction(() => onMaintenance(data))}
                  className="flex items-center w-full px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700 gap-2 rounded-b-lg"
                >
                  <ToolsIcon className="w-4 h-4 text-orange-600" />
                  <span>Maintenance</span>
                </button>
              )}
            </>
          )}
        </div>
      )}
    </div>
  );
}
