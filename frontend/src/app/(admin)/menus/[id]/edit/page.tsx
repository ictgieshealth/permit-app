"use client";
import { useState, useEffect } from "react";
import { useRouter, useParams } from "next/navigation";
import Button from "@/components/ui/button/Button";
import Input from "@/components/form/input/InputField";
import Label from "@/components/form/Label";
import Link from "next/link";
import { menuService } from "@/services/menu.service";
import { roleService } from "@/services/role.service";
import { Role } from "@/types/role";
import { Menu } from "@/types/menu";

export default function EditMenuPage() {
  const router = useRouter();
  const params = useParams();
  const id = parseInt(params.id as string);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [roles, setRoles] = useState<Role[]>([]);
  const [parentMenus, setParentMenus] = useState<Menu[]>([]);
  const [formData, setFormData] = useState({
    name: "",
    path: "",
    icon: "",
    parent_id: "",
    order_index: "0",
    role_ids: [] as number[],
  });

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      setLoading(true);
      const [menuData, rolesData, menusData] = await Promise.all([
        menuService.getById(id),
        roleService.getAll({ limit: 100, is_active: true }),
        menuService.getAll({ limit: 100, is_active: true }),
      ]);

      setFormData({
        name: menuData.name,
        path: menuData.path,
        icon: menuData.icon || "",
        parent_id: menuData.parent_id?.toString() || "",
        order_index: menuData.order_index.toString(),
        role_ids: menuData.roles?.map((r) => r.id) || [],
      });
      setRoles(rolesData.data);
      setParentMenus(menusData.data.filter((m) => m.id !== id));
    } catch (err: any) {
      setError(err.message || "Failed to load menu data");
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: value });
  };

  const handleRoleToggle = (roleId: number) => {
    setFormData((prev) => ({
      ...prev,
      role_ids: prev.role_ids.includes(roleId)
        ? prev.role_ids.filter((id) => id !== roleId)
        : [...prev.role_ids, roleId],
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    if (formData.role_ids.length === 0) {
      setError("Please select at least one role");
      return;
    }

    setSubmitting(true);
    try {
      await menuService.update(id, {
        name: formData.name,
        path: formData.path,
        icon: formData.icon || undefined,
        parent_id: formData.parent_id ? parseInt(formData.parent_id) : null,
        order_index: parseInt(formData.order_index),
        role_ids: formData.role_ids,
      });

      router.push("/menus");
    } catch (err: any) {
      setError(err.message || "Failed to update menu");
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) {
    return (
      <div className="max-w-2xl p-6 mx-auto">
        <div className="text-center text-gray-500">Loading...</div>
      </div>
    );
  }

  return (
    <div className="max-w-2xl p-6 mx-auto">
      <div className="mb-6">
        <div className="flex items-center gap-2 mb-2">
          <Link href="/menus" className="text-gray-500 hover:text-gray-700 dark:text-gray-400">
            ← Back
          </Link>
        </div>
        <h1 className="text-2xl font-bold text-gray-800 dark:text-white">Edit Menu</h1>
        <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">Update menu information</p>
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
              Menu Name <span className="text-error-500">*</span>
            </Label>
            <Input
              name="name"
              type="text"
              value={formData.name}
              onChange={handleChange}
              required
              placeholder="Enter menu name"
            />
          </div>

          <div>
            <Label>
              Path <span className="text-error-500">*</span>
            </Label>
            <Input
              name="path"
              type="text"
              value={formData.path}
              onChange={handleChange}
              required
              placeholder="e.g., /domains, /users"
            />
          </div>

          <div>
            <Label>Icon Name</Label>
            <Input
              name="icon"
              type="text"
              value={formData.icon}
              onChange={handleChange}
              placeholder="e.g., GridIcon, UserCircleIcon"
            />
          </div>

          <div>
            <Label>Parent Menu</Label>
            <select
              name="parent_id"
              value={formData.parent_id}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg dark:bg-gray-800 dark:border-gray-700 dark:text-white"
            >
              <option value="">None (Top Level)</option>
              {parentMenus.map((menu) => (
                <option key={menu.id} value={menu.id}>
                  {menu.name}
                </option>
              ))}
            </select>
          </div>

          <div>
            <Label>
              Order Index <span className="text-error-500">*</span>
            </Label>
            <Input
              name="order_index"
              type="number"
              value={formData.order_index}
              onChange={handleChange}
              required
              min="0"
            />
          </div>

          <div>
            <Label>
              Allowed Roles <span className="text-error-500">*</span>
            </Label>
            <div className="p-4 mt-2 space-y-2 border border-gray-300 rounded-lg dark:border-gray-700">
              {roles.map((role) => (
                <label
                  key={role.id}
                  className="flex items-center gap-2 cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-800 p-2 rounded"
                >
                  <input
                    type="checkbox"
                    checked={formData.role_ids.includes(role.id)}
                    onChange={() => handleRoleToggle(role.id)}
                    className="w-4 h-4 text-brand-500 border-gray-300 rounded focus:ring-brand-500"
                  />
                  <span className="text-sm text-gray-700 dark:text-gray-300">
                    {role.name}
                  </span>
                </label>
              ))}
            </div>
          </div>

          <div className="flex gap-3 pt-4">
            <Button type="submit" disabled={submitting}>
              {submitting ? "Updating..." : "Update Menu"}
            </Button>
            <Link href="/menus">
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
