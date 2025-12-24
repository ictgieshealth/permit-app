"use client";
import { useState, useRef, DragEvent } from "react";
import { TrashBinIcon } from "@/icons";

interface FileWithPreview {
  file: File;
  preview?: string;
  id: string;
}

interface FileUploadProps {
  files: File[];
  onChange: (files: File[]) => void;
  accept?: string;
  maxSize?: number; // in MB
  maxFiles?: number;
  multiple?: boolean;
}

export default function FileUpload({
  files,
  onChange,
  accept = "image/*,.pdf,.doc,.docx,.xls,.xlsx",
  maxSize = 10, // 10MB default
  maxFiles = 5,
  multiple = true,
}: FileUploadProps) {
  const [isDragging, setIsDragging] = useState(false);
  const [error, setError] = useState("");
  const [filesWithPreview, setFilesWithPreview] = useState<FileWithPreview[]>(
    []
  );
  const fileInputRef = useRef<HTMLInputElement>(null);

  const formatFileSize = (bytes: number): string => {
    if (bytes === 0) return "0 Bytes";
    const k = 1024;
    const sizes = ["Bytes", "KB", "MB", "GB"];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + " " + sizes[i];
  };

  const getFileIcon = (file: File): string => {
    const type = file.type;
    if (type.startsWith("image/")) return "🖼️";
    if (type.includes("pdf")) return "📄";
    if (type.includes("word") || type.includes("document")) return "📝";
    if (type.includes("sheet") || type.includes("excel")) return "📊";
    return "📎";
  };

  const validateFiles = (newFiles: File[]): File[] => {
    setError("");
    const validFiles: File[] = [];

    // Check max files
    if (files.length + newFiles.length > maxFiles) {
      setError(`Maximum ${maxFiles} files allowed`);
      return [];
    }

    for (const file of newFiles) {
      // Check file size
      if (file.size > maxSize * 1024 * 1024) {
        setError(`File ${file.name} exceeds ${maxSize}MB`);
        continue;
      }

      // Check duplicates
      const isDuplicate = files.some(
        (existingFile) =>
          existingFile.name === file.name && existingFile.size === file.size
      );
      if (isDuplicate) {
        setError(`File ${file.name} already added`);
        continue;
      }

      validFiles.push(file);
    }

    return validFiles;
  };

  const createPreview = (file: File): Promise<string> => {
    return new Promise((resolve) => {
      if (file.type.startsWith("image/")) {
        const reader = new FileReader();
        reader.onload = (e) => {
          resolve(e.target?.result as string);
        };
        reader.readAsDataURL(file);
      } else {
        resolve("");
      }
    });
  };

  const handleFiles = async (newFiles: FileList | null) => {
    if (!newFiles || newFiles.length === 0) return;

    const fileArray = Array.from(newFiles);
    const validFiles = validateFiles(fileArray);

    if (validFiles.length > 0) {
      const newFilesWithPreview: FileWithPreview[] = await Promise.all(
        validFiles.map(async (file) => ({
          file,
          preview: await createPreview(file),
          id: `${file.name}-${Date.now()}-${Math.random()}`,
        }))
      );

      setFilesWithPreview([...filesWithPreview, ...newFilesWithPreview]);
      onChange([...files, ...validFiles]);
    }
  };

  const handleDragOver = (e: DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    setIsDragging(true);
  };

  const handleDragLeave = (e: DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    setIsDragging(false);
  };

  const handleDrop = (e: DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    setIsDragging(false);
    handleFiles(e.dataTransfer.files);
  };

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    handleFiles(e.target.files);
  };

  const handleRemoveFile = (index: number) => {
    const newFiles = files.filter((_, i) => i !== index);
    const newFilesWithPreview = filesWithPreview.filter((_, i) => i !== index);
    setFilesWithPreview(newFilesWithPreview);
    onChange(newFiles);
  };

  const handleBrowseClick = () => {
    fileInputRef.current?.click();
  };

  return (
    <div className="w-full">
      {/* Drop Zone */}
      <div
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
        onDrop={handleDrop}
        onClick={handleBrowseClick}
        className={`relative border-2 border-dashed rounded-lg p-8 text-center cursor-pointer transition-colors ${
          isDragging
            ? "border-blue-500 bg-blue-50 dark:bg-blue-900/20"
            : "border-gray-300 dark:border-gray-600 hover:border-blue-400 dark:hover:border-blue-500"
        }`}
      >
        <input
          ref={fileInputRef}
          type="file"
          multiple={multiple}
          onChange={handleFileSelect}
          accept={accept}
          className="hidden"
        />

        <div className="space-y-2">
          <div className="text-4xl">📁</div>
          <div className="text-sm text-gray-600 dark:text-gray-400">
            <span className="font-semibold text-blue-600 dark:text-blue-400">
              Click to browse
            </span>{" "}
            or drag and drop files here
          </div>
          <div className="text-xs text-gray-500 dark:text-gray-500">
            Max {maxSize}MB per file • Up to {maxFiles} files
          </div>
        </div>
      </div>

      {/* Error Message */}
      {error && (
        <div className="mt-2 p-2 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded text-sm text-red-600 dark:text-red-400">
          {error}
        </div>
      )}

      {/* File List */}
      {filesWithPreview.length > 0 && (
        <div className="mt-4 space-y-2">
          <div className="text-sm font-medium text-gray-700 dark:text-gray-300">
            Selected Files ({filesWithPreview.length})
          </div>
          <div className="space-y-2">
            {filesWithPreview.map((item, index) => (
              <div
                key={item.id}
                className="flex items-center gap-3 p-3 bg-gray-50 dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg group hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
              >
                {/* Preview/Icon */}
                <div className="flex-shrink-0">
                  {item.preview ? (
                    <img
                      src={item.preview}
                      alt={item.file.name}
                      className="w-12 h-12 object-cover rounded"
                    />
                  ) : (
                    <div className="w-12 h-12 flex items-center justify-center text-2xl bg-white dark:bg-gray-900 rounded border border-gray-200 dark:border-gray-600">
                      {getFileIcon(item.file)}
                    </div>
                  )}
                </div>

                {/* File Info */}
                <div className="flex-1 min-w-0">
                  <div className="text-sm font-medium text-gray-900 dark:text-white truncate">
                    {item.file.name}
                  </div>
                  <div className="text-xs text-gray-500 dark:text-gray-400">
                    {formatFileSize(item.file.size)}
                  </div>
                </div>

                {/* Remove Button */}
                <button
                  type="button"
                  onClick={(e) => {
                    e.stopPropagation();
                    handleRemoveFile(index);
                  }}
                  className="flex-shrink-0 p-2 text-red-600 hover:text-red-700 dark:text-red-400 dark:hover:text-red-300 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-lg transition-colors"
                  title="Remove file"
                >
                  <TrashBinIcon className="w-4 h-4" />
                </button>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Summary */}
      {filesWithPreview.length > 0 && (
        <div className="mt-2 text-xs text-gray-500 dark:text-gray-400">
          Total size:{" "}
          {formatFileSize(files.reduce((acc, file) => acc + file.size, 0))}
        </div>
      )}
    </div>
  );
}
