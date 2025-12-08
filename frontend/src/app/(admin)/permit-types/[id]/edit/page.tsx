"use client";
import { useState, useEffect } from "react";
import { useRouter, useParams } from "next/navigation";
import Button from "@/components/ui/button/Button";
import Input from "@/components/form/input/InputField";
import Label from "@/components/form/Label";
import Link from "next/link";
import { permitTypeService } from "@/services/permitType.service";

export default function EditPermitTypePage() {
  const router = useRouter();
  const params = useParams();
  const permitTypeId = parseInt(params.id as string);

  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const [formData, setFormData] = useState({
    code: "",
    name: "",
    description: "",
  });

  useEffect(() => {
    loadPermitType();
  }, [permitTypeId]);

  const loadPermitType = async () => {
    try {
      const permitType = await permitTypeService.getById(permitTypeId);
      setFormData({
        code: permitType.code,
        name: permitType.name,
        description: permitType.description || "",
      });
      setLoading(false);
    } catch (err: any) {
      setError(err.message || "Failed to load permit type");
      setLoading(false);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: value });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    setSaving(true);
    try {
      await permitTypeService.update(permitTypeId, {
        code: formData.code,
        name: formData.name,
        description: formData.description || undefined,
      });

      router.push("/permit-types");
    } catch (err: any) {
      setError(err.message || "Failed to update permit type");
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
    <div className="max-w-2xl p-6 mx-auto">
      <div className="mb-6">
        <div className="flex items-center gap-2 mb-2">
          <Link href="/permit-types" className="text-gray-500 hover:text-gray-700 dark:text-gray-400">
            ← Back
          </Link>
        </div>
        <h1 className="text-2xl font-bold text-gray-800 dark:text-white">Edit Permit Type</h1>
        <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
          Update permit type information
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
              Permit Type Code <span className="text-error-500">*</span>
            </Label>
            <Input
              name="code"
              type="text"
              value={formData.code}
              onChange={handleChange}
              required
              placeholder="Enter permit type code"
            />
          </div>

          <div>
            <Label>
              Permit Type Name <span className="text-error-500">*</span>
            </Label>
            <Input
              name="name"
              type="text"
              value={formData.name}
              onChange={handleChange}
              required
              placeholder="Enter permit type name"
            />
          </div>

          <div>
            <Label>Description</Label>
            <textarea
              name="description"
              value={formData.description}
              onChange={handleChange}
              rows={4}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg dark:bg-gray-800 dark:border-gray-700 dark:text-white"
              placeholder="Enter description (optional)"
            />
          </div>

          <div className="flex gap-3 pt-4">
            <Button type="submit" disabled={saving}>
              {saving ? "Saving..." : "Update Permit Type"}
            </Button>
            <Link href="/permit-types">
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
