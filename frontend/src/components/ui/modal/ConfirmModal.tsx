"use client";
import React from "react";
import { Modal } from "./index";
import Button from "@/components/ui/button/Button";

interface ConfirmModalProps {
  isOpen: boolean;
  onClose: () => void;
  onConfirm: () => void;
  title: string;
  message: string;
  confirmText?: string;
  cancelText?: string;
  type?: "danger" | "warning" | "info";
  loading?: boolean;
}

export const ConfirmModal: React.FC<ConfirmModalProps> = ({
  isOpen,
  onClose,
  onConfirm,
  title,
  message,
  confirmText = "Confirm",
  cancelText = "Cancel",
  type = "danger",
  loading = false,
}) => {
  const getIconColor = () => {
    switch (type) {
      case "danger":
        return "text-error-500";
      case "warning":
        return "text-warning-500";
      case "info":
        return "text-brand-500";
      default:
        return "text-error-500";
    }
  };

  const getConfirmButtonClass = () => {
    switch (type) {
      case "danger":
        return "bg-error-500 hover:bg-error-600";
      case "warning":
        return "bg-warning-500 hover:bg-warning-600";
      case "info":
        return "bg-brand-500 hover:bg-brand-600";
      default:
        return "bg-error-500 hover:bg-error-600";
    }
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose} className="max-w-md">
      <div className="p-6">
        {/* Icon */}
        <div className="flex justify-center mb-4">
          <div className={`flex items-center justify-center w-16 h-16 rounded-full ${type === 'danger' ? 'bg-error-100 dark:bg-error-900/20' : type === 'warning' ? 'bg-warning-100 dark:bg-warning-900/20' : 'bg-brand-100 dark:bg-brand-900/20'}`}>
            <svg
              className={`w-8 h-8 ${getIconColor()}`}
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
              />
            </svg>
          </div>
        </div>

        {/* Title */}
        <h3 className="mb-2 text-xl font-bold text-center text-gray-900 dark:text-white">
          {title}
        </h3>

        {/* Message */}
        <p className="mb-6 text-sm text-center text-gray-600 dark:text-gray-400">
          {message}
        </p>

        {/* Actions */}
        <div className="flex gap-3">
          <Button
            onClick={onClose}
            disabled={loading}
            className="flex-1 bg-gray-500 hover:bg-gray-600"
          >
            {cancelText}
          </Button>
          <Button
            onClick={onConfirm}
            disabled={loading}
            className={`flex-1 ${getConfirmButtonClass()}`}
          >
            {loading ? "Processing..." : confirmText}
          </Button>
        </div>
      </div>
    </Modal>
  );
};
