"use client";
import { useState, useEffect } from "react";
import { useRouter, useParams } from "next/navigation";
import Button from "@/components/ui/button/Button";
import Input from "@/components/form/input/InputField";
import Label from "@/components/form/Label";
import Link from "next/link";
import FileUpload from "@/components/form/FileUpload";
import { Domain } from "@/types/domain";
import { Division } from "@/types/division";
import { PermitType } from "@/types/permitType";
import { User } from "@/types/user";
import { Permit } from "@/types/permit";
import { permitService } from "@/services/permit.service";
import { domainService } from "@/services/domain.service";
import { divisionService } from "@/services/division.service";
import { permitTypeService } from "@/services/permitType.service";
import { userService } from "@/services/user.service";

export default function EditPermitPage() {
  const router = useRouter();
  const params = useParams();
  const permitId = parseInt(params.id as string);

  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const [uploadedFile, setUploadedFile] = useState<File | null>(null);
  const [currentPermit, setCurrentPermit] = useState<Permit | null>(null);
  const [domains, setDomains] = useState<Domain[]>([]);
  const [divisions, setDivisions] = useState<Division[]>([]);
  const [permitTypes, setPermitTypes] = useState<PermitType[]>([]);
  const [users, setUsers] = useState<User[]>([]);
  const [formData, setFormData] = useState({
    domain_id: "",
    division_id: "",
    permit_type_id: "",
    name: "",
    application_type: "",
    permit_no: "",
    effective_date: "",
    expiry_date: "",
    effective_term: "",
    responsible_person_id: "",
    responsible_doc_person_id: "",
    doc_name: "",
    doc_number: "",
    status: "active",
  });

  useEffect(() => {
    loadPermit();
    loadDomains();
    loadPermitTypes();
    loadUsers();
    return () => {
      // Cleanup blob URLs when component unmounts
      if (previewUrl && previewUrl.startsWith('blob:')) {
        URL.revokeObjectURL(previewUrl);
      }
      if (newFilePreview && newFilePreview.startsWith('blob:')) {
        URL.revokeObjectURL(newFilePreview);
      }
    };
  }, [permitId]);

  useEffect(() => {
    if (formData.domain_id) {
      loadDivisions(parseInt(formData.domain_id));
    } else {
      setDivisions([]);
    }
  }, [formData.domain_id]);

  const loadPermit = async () => {
    try {
      const permit = await permitService.getById(permitId);
      setCurrentPermit(permit);
      if (permit.doc_file_name && isPreviewable(permit.doc_file_type)) {
        // Load preview blob with authorization
        const blob = await permitService.getPreviewBlob(permitId);
        const blobUrl = URL.createObjectURL(blob);
        setPreviewUrl(blobUrl);
        setShowPreview(true);
      }
      setFormData({
        domain_id: permit.domain_id.toString(),
        division_id: permit.division_id ? permit.division_id.toString() : "",
        permit_type_id: permit.permit_type_id.toString(),
        name: permit.name,
        application_type: permit.application_type,
        permit_no: permit.permit_no,
        effective_date: permit.effective_date.split("T")[0],
        expiry_date: permit.expiry_date.split("T")[0],
        effective_term: permit.effective_term || "",
        responsible_person_id: permit.responsible_person_id ? permit.responsible_person_id.toString() : "",
        responsible_doc_person_id: permit.responsible_doc_person_id ? permit.responsible_doc_person_id.toString() : "",
        doc_name: permit.doc_name || "",
        doc_number: permit.doc_number || "",
        status: permit.status,
      });
      setLoading(false);
    } catch (err: any) {
      setError(err.message || "Failed to load permit");
      setLoading(false);
    }
  };

  const isPreviewable = (fileType?: string | null) => {
    if (!fileType) return false;
    const previewableTypes = [
      "application/pdf",
      "image/jpeg",
      "image/jpg",
      "image/png",
      "image/gif",
      "image/webp"
    ];
    return previewableTypes.includes(fileType.toLowerCase());
  };

  const loadDomains = async () => {
    try {
      const response = await domainService.getAll({ limit: 100, is_active: true });
      setDomains(response.data);
    } catch (err) {
      console.error("Failed to load domains:", err);
    }
  };

  const loadDivisions = async (domainId: number) => {
    try {
      const response = await divisionService.getAll({ limit: 100, domain_id: domainId });
      setDivisions(response.data);
    } catch (err) {
      console.error("Failed to load divisions:", err);
    }
  };

  const loadPermitTypes = async () => {
    try {
      const response = await permitTypeService.getAll({ limit: 100 });
      setPermitTypes(response.data);
    } catch (err) {
      console.error("Failed to load permit types:", err);
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
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>
  ) => {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: value });
  };

  const handleDownloadDocument = async () => {
    try {
      await permitService.downloadDocument(permitId);
    } catch (err: any) {
      setError(err.message || "Failed to download document");
    }
  };

  const formatFileSize = (bytes?: number | null) => {
    if (!bytes) return "-";
    if (bytes < 1024) return bytes + " B";
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + " KB";
    return (bytes / (1024 * 1024)).toFixed(2) + " MB";
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    if (!formData.domain_id) {
      setError("Please select a domain");
      return;
    }

    if (!formData.permit_type_id) {
      setError("Please select a permit type");
      return;
    }

    setSaving(true);
    try {
      await permitService.update(permitId, {
        domain_id: parseInt(formData.domain_id),
        division_id: formData.division_id ? parseInt(formData.division_id) : undefined,
        permit_type_id: parseInt(formData.permit_type_id),
        name: formData.name,
        application_type: formData.application_type,
        permit_no: formData.permit_no,
        effective_date: formData.effective_date,
        expiry_date: formData.expiry_date,
        effective_term: formData.effective_term || undefined,
        responsible_person_id: formData.responsible_person_id ? parseInt(formData.responsible_person_id) : undefined,
        responsible_doc_person_id: formData.responsible_doc_person_id ? parseInt(formData.responsible_doc_person_id) : undefined,
        doc_name: formData.doc_name || undefined,
        doc_number: formData.doc_number || undefined,
        status: formData.status,
      }, uploadedFile || undefined);

      router.push("/permits");
    } catch (err: any) {
      setError(err.message || "Failed to update permit");
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-gray-500">Loading...</div>
      </div>
    );
  }

  return (
    <div className="max-w-4xl p-6 mx-auto">
      <div className="mb-6">
        <div className="flex items-center gap-2 mb-2">
          <Link
            href="/permits"
            className="text-gray-500 hover:text-gray-700 dark:text-gray-400"
          >
            ← Back
          </Link>
        </div>
        <h1 className="text-2xl font-bold text-gray-800 dark:text-white">
          Edit Permit
        </h1>
        <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
          Update permit information
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
            <div>
              <Label>
                Domain <span className="text-error-500">*</span>
              </Label>
              <select
                name="domain_id"
                value={formData.domain_id}
                onChange={handleChange}
                required
                className="w-full px-3 py-2 border border-gray-300 rounded-lg dark:bg-gray-800 dark:border-gray-700 dark:text-white"
              >
                <option value="">Select a domain</option>
                {domains.map((domain) => (
                  <option key={domain.id} value={domain.id}>
                    {domain.name} ({domain.code})
                  </option>
                ))}
              </select>
            </div>

            <div>
              <Label>
                Category <span className="text-error-500">*</span>
              </Label>
              <select
                name="permit_type_id"
                value={formData.permit_type_id}
                onChange={handleChange}
                required
                className="w-full px-3 py-2 border border-gray-300 rounded-lg dark:bg-gray-800 dark:border-gray-700 dark:text-white"
              >
                <option value="">Select a category</option>
                {permitTypes.map((type) => (
                  <option key={type.id} value={type.id}>
                    {type.name}
                  </option>
                ))}
              </select>
            </div>
          </div>

          <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
            <div>
              <Label>
                Permit Name <span className="text-error-500">*</span>
              </Label>
              <Input
                name="name"
                type="text"
                value={formData.name}
                onChange={handleChange}
                required
                placeholder="Enter permit name (equipment/service/competency/operational)"
              />
            </div>

            <div>
              <Label>Division</Label>
              <select
                name="division_id"
                value={formData.division_id}
                onChange={handleChange}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg dark:bg-gray-800 dark:border-gray-700 dark:text-white"
                disabled={!formData.domain_id}
              >
                <option value="">Select a division</option>
                {divisions.map((division) => (
                  <option key={division.id} value={division.id}>
                    {division.name} ({division.code})
                  </option>
                ))}
              </select>
            </div>
          </div>

          <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
            <div>
              <Label>
                Permit Number <span className="text-error-500">*</span>
              </Label>
              <Input
                name="permit_no"
                type="text"
                value={formData.permit_no}
                onChange={handleChange}
                required
                placeholder="Enter permit number"
              />
            </div>

            <div>
              <Label>
                Application Type <span className="text-error-500">*</span>
              </Label>
              <Input
                name="application_type"
                type="text"
                value={formData.application_type}
                onChange={handleChange}
                required
                placeholder="Enter application type"
              />
            </div>
          </div>

          <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
            <div>
              <Label>
                Effective Date <span className="text-error-500">*</span>
              </Label>
              <Input
                name="effective_date"
                type="date"
                value={formData.effective_date}
                onChange={handleChange}
                required
              />
            </div>

            <div>
              <Label>
                Expiry Date <span className="text-error-500">*</span>
              </Label>
              <Input
                name="expiry_date"
                type="date"
                value={formData.expiry_date}
                onChange={handleChange}
                required
              />
            </div>
          </div>

          <div>
            <Label>Effective Term</Label>
            <Input
              name="effective_term"
              type="text"
              value={formData.effective_term}
              onChange={handleChange}
              placeholder="Enter effective term (e.g., 1 year)"
            />
          </div>

          <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
            <div>
              <Label>Permit Responsible Person</Label>
              <select
                name="responsible_person_id"
                value={formData.responsible_person_id}
                onChange={handleChange}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg dark:bg-gray-800 dark:border-gray-700 dark:text-white"
              >
                <option value="">Select a person</option>
                {users.map((user) => (
                  <option key={user.id} value={user.id}>
                    {user.full_name} ({user.email})
                  </option>
                ))}
              </select>
            </div>

            <div>
              <Label>Document Responsible Person</Label>
              <select
                name="responsible_doc_person_id"
                value={formData.responsible_doc_person_id}
                onChange={handleChange}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg dark:bg-gray-800 dark:border-gray-700 dark:text-white"
              >
                <option value="">Select a person</option>
                {users.map((user) => (
                  <option key={user.id} value={user.id}>
                    {user.full_name} ({user.email})
                  </option>
                ))}
              </select>
            </div>
          </div>

          <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
            <div>
              <Label>Document Name</Label>
              <Input
                name="doc_name"
                type="text"
                value={formData.doc_name}
                onChange={handleChange}
                placeholder="Enter document name"
              />
            </div>

            <div>
              <Label>Document Number</Label>
              <Input
                name="doc_number"
                type="text"
                value={formData.doc_number}
                onChange={handleChange}
                placeholder="Enter document number"
              />
            </div>
          </div>

          {/* Current Document Info */}
          {currentPermit?.doc_file_name && (
            <div className="p-4 border border-gray-200 rounded-lg bg-gray-50 dark:bg-gray-800 dark:border-gray-700">
              <h4 className="mb-3 text-sm font-medium text-gray-700 dark:text-gray-300">
                Current Document
              </h4>
              <div className="grid grid-cols-1 gap-3 md:grid-cols-3">
                <div>
                  <label className="block mb-1 text-xs font-medium text-gray-600 dark:text-gray-400">
                    File Name
                  </label>
                  <p className="text-sm text-gray-900 break-all dark:text-white">
                    {currentPermit.doc_file_name}
                  </p>
                </div>
                <div>
                  <label className="block mb-1 text-xs font-medium text-gray-600 dark:text-gray-400">
                    File Size
                  </label>
                  <p className="text-sm text-gray-900 dark:text-white">
                    {formatFileSize(currentPermit.doc_file_size)}
                  </p>
                </div>
                <div>
                  <label className="block mb-1 text-xs font-medium text-gray-600 dark:text-gray-400">
                    File Type
                  </label>
                  <p className="text-sm text-gray-900 dark:text-white">
                    {currentPermit.doc_file_type || "-"}
                  </p>
                </div>
              </div>
              <div className="mt-3">
                <Button
                  type="button"
                  onClick={handleDownloadDocument}
                  className="text-sm bg-blue-600 hover:bg-blue-700"
                >
                  📥 Download Current Document
                </Button>
              </div>
            </div>
          )}

          {/* Upload New Document */}
          <div>
            <Label>
              {currentPermit?.doc_file_name ? "Replace Document (Optional)" : "Upload Document"}
            </Label>
            <FileUpload
              onFileSelect={(file) => setUploadedFile(file)}
              disabled={saving}
            />
            <p className="mt-1 text-xs text-gray-500">
              {currentPermit?.doc_file_name 
                ? "Upload a new file to replace the current document. Leave empty to keep existing."
                : "Allowed file types: PDF, DOC, DOCX, JPG, PNG (max 10MB)"
              }
            </p>
            {uploadedFile && (
              <div className="p-3 mt-2 border border-green-200 rounded-lg bg-green-50 dark:bg-green-900/20 dark:border-green-800">
                <p className="text-sm text-green-700 dark:text-green-400">
                  ✓ New file selected: <strong>{uploadedFile.name}</strong> ({formatFileSize(uploadedFile.size)})
                </p>
              </div>
            )}
          </div>

          <div>
            <Label>
              Status <span className="text-error-500">*</span>
            </Label>
            <select
              name="status"
              value={formData.status}
              onChange={handleChange}
              required
              className="w-full px-3 py-2 border border-gray-300 rounded-lg dark:bg-gray-800 dark:border-gray-700 dark:text-white"
            >
              <option value="active">Active</option>
              <option value="expired">Expired</option>
              <option value="inactive">Inactive</option>
            </select>
          </div>

          <div className="flex gap-3 pt-4">
            <Button type="submit" disabled={saving}>
              {saving ? "Saving..." : "Update Permit"}
            </Button>
            <Link href="/permits">
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
