"use client";
import { useState, useEffect } from "react";
import { useRouter, useParams } from "next/navigation";
import Button from "@/components/ui/button/Button";
import Input from "@/components/form/input/InputField";
import Label from "@/components/form/Label";
import Link from "next/link";
import { Domain } from "@/types/domain";
import { divisionService } from "@/services/division.service";
import { domainService } from "@/services/domain.service";

export default function EditDivisionPage() {
  const router = useRouter();
  const params = useParams();
  const divisionId = parseInt(params.id as string);

  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const [domains, setDomains] = useState<Domain[]>([]);
  const [formData, setFormData] = useState({
    domain_id: "",
    code: "",
    name: "",
  });

  useEffect(() => {
    loadDivision();
    loadDomains();
  }, [divisionId]);

  const loadDivision = async () => {
    try {
      const division = await divisionService.getById(divisionId);
      setFormData({
        domain_id: division.domain_id.toString(),
        code: division.code,
        name: division.name,
      });
      setLoading(false);
    } catch (err: any) {
      setError(err.message || "Failed to load division");
      setLoading(false);
    }
  };

  const loadDomains = async () => {
    try {
      const response = await domainService.getAll({ limit: 100, is_active: true });
      setDomains(response.data);
    } catch (err) {
      console.error("Failed to load domains:", err);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: value });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    if (!formData.domain_id) {
      setError("Please select a domain");
      return;
    }

    setSaving(true);
    try {
      await divisionService.update(divisionId, {
        domain_id: parseInt(formData.domain_id),
        code: formData.code,
        name: formData.name,
      });

      router.push("/divisions");
    } catch (err: any) {
      setError(err.message || "Failed to update division");
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
          <Link href="/divisions" className="text-gray-500 hover:text-gray-700 dark:text-gray-400">
            ← Back
          </Link>
        </div>
        <h1 className="text-2xl font-bold text-gray-800 dark:text-white">Edit Division</h1>
        <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
          Update division information
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
              Division Code <span className="text-error-500">*</span>
            </Label>
            <Input
              name="code"
              type="text"
              value={formData.code}
              onChange={handleChange}
              required
              placeholder="Enter division code"
            />
          </div>

          <div>
            <Label>
              Division Name <span className="text-error-500">*</span>
            </Label>
            <Input
              name="name"
              type="text"
              value={formData.name}
              onChange={handleChange}
              required
              placeholder="Enter division name"
            />
          </div>

          <div className="flex gap-3 pt-4">
            <Button type="submit" disabled={saving}>
              {saving ? "Saving..." : "Update Division"}
            </Button>
            <Link href="/divisions">
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
